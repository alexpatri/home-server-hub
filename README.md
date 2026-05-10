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

`PUT /applications/{id}` aceita `multipart/form-data` com os mesmos campos do `POST` (`name`, `port`, `url`, `image`); apenas os campos enviados são atualizados, os demais permanecem inalterados.

`DELETE /applications/{id}` responde `204 No Content` em sucesso e `404` se o ID não existir.

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
├── cmd/main.go             # Entry point
├── internal/
│   ├── api/                # Servidor Fiber e rotas
│   ├── config/             # Carregamento de configuração via env
│   ├── controllers/        # Handlers HTTP
│   ├── database/           # Conexão e bootstrap do SQLite
│   ├── docker/             # Cliente Docker (descoberta de containers)
│   ├── models/             # Tipos de domínio + interface do repositório
│   ├── repository/         # Implementação SQLite
│   ├── services/           # Regra de negócio
│   └── utils/              # Helpers (parse de imagem, network)
├── docker/Dockerfile       # Build multi-stage (Go puro, sem CGO)
├── docker-compose.yml
└── docs/                   # Swagger gerado
```

### Executando localmente

```bash
go run ./cmd/main.go --dotenv   # carrega .env se presente
```

Para regerar a documentação Swagger após mudanças nas anotações:

```bash
swag init -g cmd/main.go
```

### Backup

O banco é um único arquivo. Para um snapshot consistente (incluindo o WAL):

```bash
sqlite3 ./data/home_server_hub.db ".backup ./backup.db"
```

## Roadmap Futuro

- Monitoramento de recursos do servidor (CPU, memória, disco)
- Gerenciamento de containers (iniciar, parar, reiniciar)
- Visualização de logs
- Notificações de status

## Contribuição

Este projeto está em fase inicial de desenvolvimento. Contribuições são bem-vindas!
