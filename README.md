# Home Server Hub

Um hub centralizado para gerenciar e acessar aplicações rodando em seu servidor doméstico.

## Visão Geral

Home Server Hub é uma aplicação simples que serve como ponto central para visualizar, acessar e gerenciar todas as aplicações rodando no seu servidor doméstico. Através de uma interface limpa e intuitiva, você pode facilmente ver o status de suas aplicações, acessá-las com um clique e gerenciar suas configurações.

## Funcionalidades

- **Dashboard centralizado**: Visualize todas as suas aplicações em um só lugar
- **Status em tempo real**: Verifique se suas aplicações estão funcionando corretamente
- **Acesso rápido**: Acesse qualquer aplicação com apenas um clique
- **Gerenciamento de aplicações**: Adicione, edite ou remova aplicações do hub
- **Descoberta automática**: Detecte automaticamente aplicações Docker rodando em seu servidor

## Arquitetura

### Stack Tecnológica

- **Backend**: Go com Fiber — API RESTful de alta performance e baixo consumo
- **Banco de Dados**: SQLite (arquivo único, driver Go puro `modernc.org/sqlite`) — zero dependência externa, backup é só copiar o arquivo
- **Armazenamento de imagens**: arquivos em disco em `IMAGES_DIR` (apenas metadados ficam no banco)
- **Integração**: Docker API — para descoberta e monitoramento de aplicações containerizadas

### Estrutura de Dados

#### Tabela `applications`

| Coluna         | Tipo    | Descrição                                |
| -------------- | ------- | ---------------------------------------- |
| `id`           | TEXT    | UUID v4, chave primária                  |
| `name`         | TEXT    | Nome da aplicação                        |
| `tags`         | TEXT    | Array de tags serializado como JSON      |
| `container`    | TEXT    | ID do container Docker associado         |
| `ip`           | TEXT    | IP do host                               |
| `port`         | INTEGER | Porta exposta da aplicação               |
| `url`          | TEXT    | URL personalizada (opcional)             |
| `image_name`   | TEXT    | Nome do arquivo de imagem original       |
| `image_width`  | INTEGER | Largura da imagem em pixels              |
| `image_height` | INTEGER | Altura da imagem em pixels               |
| `image_size`   | INTEGER | Tamanho do arquivo em bytes              |

O conteúdo binário das imagens é gravado em `<IMAGES_DIR>/<id>` e servido via endpoint próprio.

### API Endpoints

#### Aplicações

| Método | Rota                            | Descrição                                            |
| ------ | ------------------------------- | ---------------------------------------------------- |
| GET    | `/applications`                 | Lista todas as aplicações registradas                |
| POST   | `/applications`                 | Cria uma nova aplicação a partir de um container     |
| GET    | `/applications/discover`        | Descobre containers Docker disponíveis               |
| GET    | `/applications/{id}`            | Retorna uma aplicação específica                     |
| PUT    | `/applications/{id}`            | Atualiza parcialmente uma aplicação                  |
| DELETE | `/applications/{id}`            | Remove uma aplicação                                 |
| GET    | `/applications/{id}/image`      | Retorna o arquivo de imagem da aplicação             |
| POST   | `/applications/{id}/start`      | Inicia o container associado                         |
| POST   | `/applications/{id}/stop`       | Para o container associado (timeout de 10s)          |
| POST   | `/applications/{id}/restart`    | Reinicia o container associado                       |
| GET    | `/applications/events`          | Stream SSE de atualizações de status em tempo real   |

`PUT /applications/{id}` aceita `multipart/form-data` com os mesmos campos do `POST` (`name`, `port`, `url`, `image`); apenas os campos enviados são atualizados, os demais permanecem inalterados.

`DELETE /applications/{id}` responde `204 No Content` em sucesso e `404` se o ID não existir.

As três ações de container respondem `200 OK` com a aplicação atualizada (incluindo o novo `status`), `404` se o ID não existir, ou `500` com a mensagem do Docker em caso de falha (ex.: container já removido).

#### Stream de status em tempo real

`GET /applications/events` é uma conexão **Server-Sent Events**. O backend consome o stream de eventos do daemon Docker (`docker events`), resolve `container_id → application.id`, e emite uma linha por mudança de estado:

```
data: {"id":"<app-id>","status":"running"}

data: {"id":"<app-id>","status":"stopped"}

```

Funciona para mudanças causadas pelo próprio hub (via `/start`/`/stop`/`/restart`) **e** para mudanças externas (`docker stop` no terminal, Watchtower, etc.). Pings de keep-alive (`: ping`) são enviados a cada 25 segundos para evitar timeout em proxies. O fluxo recomendado do frontend é: carregar o estado inicial via `GET /applications`, depois abrir `new EventSource('/api/v1/applications/events')` e fazer patch incremental.

Todas as rotas estão sob o prefixo `/api/v1`.

#### Resposta de exemplo (GET /applications)

```json
{
    "applications": [
        {
            "id": "uuid da aplicação",
            "name": "nome da aplicação",
            "tags": ["array de tags"],
            "image": {
                "name": "icon.png",
                "height": 128,
                "width": 128,
                "size": 4321
            },
            "ip": "192.168.0.10",
            "port": 8080,
            "url": "https://app.exemplo.local",
            "status": "running"
        }
    ],
    "total": 1
}
```

A imagem em si é recuperada via `GET /api/v1/applications/{id}/image`.

## Uso

### Pré-requisitos

- Docker e Docker Compose instalados no servidor
- Go 1.23+ (apenas para compilar localmente, fora do container)

### Instalação

```bash
git clone <repo>
cd home-server-hub
docker compose up --build -d
```

A API ficará disponível em `http://localhost:8000` e o Swagger em `http://localhost:8000/docs`.

### Configuração

Variáveis de ambiente (com valores padrão):

| Variável      | Padrão                            | Descrição                                |
| ------------- | --------------------------------- | ---------------------------------------- |
| `SERVER_PORT` | `8000`                            | Porta HTTP do servidor                   |
| `SQLITE_PATH` | `/app/data/home_server_hub.db`    | Caminho do arquivo SQLite                |
| `IMAGES_DIR`  | `/app/data/images`                | Diretório onde imagens são armazenadas   |
| `DOCKER_HOST` | `unix:///var/run/docker.sock`     | Endpoint do daemon Docker                |

O `docker-compose.yml` monta `./data` em `/app/data`, então banco e imagens persistem entre reinicializações no host.

## Desenvolvimento

### Estrutura do Projeto

```
home-server-hub/
├── cmd/server/main.go      # Entry point do servidor
├── cmd/tui/main.go          # TUI (Bubble Tea) - experimental
├── internal/
│   ├── api/                # Servidor Fiber e rotas
│   ├── database/           # Conexão e bootstrap do SQLite
│   ├── events/             # Broadcaster de eventos SSE
│   ├── handlers/           # Handlers HTTP
│   ├── models/             # Tipos de domínio + interface do repositório
│   ├── repository/         # Implementação SQLite
│   ├── services/           # Regra de negócio
│   └── utils/              # Helpers
│       ├── config/         # Carregamento de configuração via env
│       ├── docker/         # Cliente Docker (descoberta de containers)
│       ├── image.go        # Parse de imagem
│       └── network/        # Utilitários de rede
├── docker/Dockerfile       # Build multi-stage (Go puro, sem CGO)
├── docker-compose.yml
└── docs/                   # Swagger gerado
```

### Executando localmente

```bash
go run ./cmd/server/main.go --dotenv   # carrega .env se presente
```

Para regerar a documentação Swagger após mudanças nas anotações:

```bash
swag init -g cmd/server/main.go
```

### Backup

O banco é um único arquivo. Para um snapshot consistente (incluindo o WAL):

```bash
sqlite3 ./data/home_server_hub.db ".backup ./backup.db"
```

## Roadmap Futuro

- Streaming de logs dos containers via SSE/WebSocket
- Métricas por container (uso de CPU e memória)
- Monitoramento de recursos do host (CPU, memória, disco)
- Cobertura de testes automatizados

## Contribuição

Este projeto está em fase inicial de desenvolvimento. Contribuições são bem-vindas!
