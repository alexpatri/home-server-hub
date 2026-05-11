// Package events implementa um broadcaster que consome eventos do daemon
// Docker, resolve container_id → aplicação registrada, e faz fan-out de
// atualizações de status para múltiplos subscribers.
package events

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	dockerevents "github.com/docker/docker/api/types/events"

	"home-server-hub/internal/docker"
	"home-server-hub/internal/models"
)

const (
	subscriberBufferSize = 16
	reconnectDelay       = 5 * time.Second
)

// StatusEvent é o payload que vai para cada subscriber.
type StatusEvent struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// Subscription representa um consumidor inscrito no broadcaster.
type Subscription struct {
	ch     chan StatusEvent
	closer func()
}

// Channel devolve o canal de leitura dos eventos.
func (s *Subscription) Channel() <-chan StatusEvent {
	return s.ch
}

// Close cancela a inscrição. Idempotente.
func (s *Subscription) Close() {
	s.closer()
}

// Broadcaster consome o stream de eventos do Docker e distribui atualizações
// de status para subscribers ativos.
type Broadcaster struct {
	docker *docker.Client
	repo   models.ApplicationRepository

	mu          sync.RWMutex
	subscribers map[chan StatusEvent]struct{}
}

// NewBroadcaster cria um broadcaster pronto para ser executado via Run.
func NewBroadcaster(docker *docker.Client, repo models.ApplicationRepository) *Broadcaster {
	return &Broadcaster{
		docker:      docker,
		repo:        repo,
		subscribers: make(map[chan StatusEvent]struct{}),
	}
}

// Subscribe registra um novo consumidor.
func (b *Broadcaster) Subscribe() *Subscription {
	ch := make(chan StatusEvent, subscriberBufferSize)
	b.mu.Lock()
	b.subscribers[ch] = struct{}{}
	b.mu.Unlock()

	var once sync.Once
	closer := func() {
		once.Do(func() {
			b.mu.Lock()
			delete(b.subscribers, ch)
			b.mu.Unlock()
			close(ch)
		})
	}
	return &Subscription{ch: ch, closer: closer}
}

// Run bloqueia consumindo eventos do daemon e fazendo fan-out.
// Reconecta automaticamente se o stream cair, até ctx ser cancelado.
func (b *Broadcaster) Run(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}
		if err := b.consume(ctx); err != nil && ctx.Err() == nil {
			log.Printf("docker events: stream caiu (%v); reconectando em %s", err, reconnectDelay)
			select {
			case <-ctx.Done():
				return
			case <-time.After(reconnectDelay):
			}
		}
	}
}

func (b *Broadcaster) consume(ctx context.Context) error {
	msgs, errs := b.docker.Events(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errs:
			if err == nil {
				return errors.New("docker events: canal de erro fechou inesperadamente")
			}
			return err
		case msg, ok := <-msgs:
			if !ok {
				return errors.New("docker events: canal de mensagens fechou inesperadamente")
			}
			b.handle(msg)
		}
	}
}

func (b *Broadcaster) handle(msg dockerevents.Message) {
	status, ok := actionToStatus(msg.Action)
	if !ok {
		return
	}
	containerID := msg.Actor.ID
	if containerID == "" {
		return
	}
	app, err := b.repo.FindByContainer(containerID)
	if err != nil {
		// App não cadastrada ou erro de DB — ignora silenciosamente exceto logs.
		if !errors.Is(err, models.ErrApplicationNotFound) {
			log.Printf("docker events: erro ao resolver container %s: %v", containerID, err)
		}
		return
	}
	b.dispatch(StatusEvent{ID: app.ID, Status: status})
}

func (b *Broadcaster) dispatch(ev StatusEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.subscribers {
		select {
		case ch <- ev:
		default:
			// Subscriber lento; o evento é descartado para este cliente.
			// Reconectar fará o frontend re-sincronizar via GET /applications.
		}
	}
}

func actionToStatus(action dockerevents.Action) (string, bool) {
	switch action {
	case dockerevents.ActionStart:
		return "running", true
	case dockerevents.ActionDie:
		return "stopped", true
	default:
		return "", false
	}
}
