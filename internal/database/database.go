package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"home-server-hub/internal/utils/config"
)

const schema = `
CREATE TABLE IF NOT EXISTS applications (
    id           TEXT PRIMARY KEY,
    name         TEXT NOT NULL,
    tags         TEXT NOT NULL DEFAULT '[]',
    container    TEXT NOT NULL DEFAULT '',
    ip           TEXT NOT NULL DEFAULT '',
    port         INTEGER NOT NULL DEFAULT 0,
    url          TEXT NOT NULL DEFAULT '',
    image_name   TEXT NOT NULL DEFAULT '',
    image_width  INTEGER NOT NULL DEFAULT 0,
    image_height INTEGER NOT NULL DEFAULT 0,
    image_size   INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_applications_container ON applications(container);
`

// Connect abre o arquivo SQLite, aplica PRAGMAs e garante o schema.
func Connect(dbConfig config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)&_pragma=busy_timeout(5000)",
		filepath.ToSlash(dbConfig.Path))

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// SQLite não lida bem com várias conexões abertas escrevendo em paralelo.
	db.SetMaxOpenConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	if _, err := db.ExecContext(ctx, schema); err != nil {
		_ = db.Close()
		return nil, err
	}

	log.Printf("Conectado ao SQLite em %s", dbConfig.Path)
	return db, nil
}

// Disconnect fecha a conexão com o banco.
func Disconnect(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Erro ao fechar SQLite: %v", err)
		return
	}
	log.Println("Desconectado do SQLite")
}
