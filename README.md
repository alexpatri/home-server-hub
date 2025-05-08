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

- **Frontend**: Vue.js - Interface de usuário reativa e moderna
- **Backend**: Go com Fiber - API RESTful de alta performance e baixo consumo
- **Banco de Dados**: MongoDB - Armazenamento de dados flexível e escalável
- **Integração**: Docker API - Para descoberta e monitoramento de aplicações containerizadas

### Estrutura de Dados

#### Collection `applications`

```json
{
    "_id": "objectid único de cada aplicação",
    "name": "nome da aplicação",
    "tags": ["array com as tags da aplicação"],
    "image": {
        "name": "nome da imagem",
        "data": "arquivo da imagem",
        "height": "altura",
        "widht": "largura",
        "size": "tamanho"
    },
    "container": "nome ou id do container docker",
    "ip": "ip do host",
    "port": "porta da aplicação",
    "url": "url personalizada caso haja"
}
```

### API Endpoints

#### Aplicações

| Método | Rota                     | Descrição                                      |
| ------ | ------------------------ | ---------------------------------------------- |
| GET    | `/applications`          | Lista todas as aplicações registradas          |
| PUT    | `/applications/{id}`     | Atualiza os dados de uma aplicação específica  |
| POST   | `/applications/discover` | Descobre containers Docker e os registra       |
| DELETE | `/applications/{id}`     | Remove uma aplicação registrada *(planejado)*  |
| GET    | `/applications/{id}`     | Retorna uma aplicação específica *(planejado)* |


#### Resposta de exemplo (GET /applications)

```json
[
    {
        "id": "objectid único de cada aplicação",
        "name": "nome da aplicação",
        "tags": ["array com as tags da aplicação"],
        "image": "arquivo de imagem",
        "ip": "ip do host",
        "port": "porta da aplicação",
        "url": "url personalizada caso haja",
        "status": "status da aplicação de acordo com o status do container"
    }
]
```

## Uso

### Pré-requisitos

- Docker instalado no servidor
- MongoDB rodando
- Go 1.16+ (para compilar o backend)
- Node.js 14+ (para o frontend)

### Instalação

[Instruções de instalação serão adicionadas]

### Configuração

[Instruções de configuração serão adicionadas]

## Desenvolvimento

### Estrutura do Projeto

```
/home-server-hub
├── /backend            # API Go com Fiber
│   ├── /controllers    # Controladores da API
│   ├── /models         # Modelos de dados
│   ├── /services       # Serviços de negócio
│   └── /config         # Configurações
├── /frontend           # Aplicação Vue.js
│   ├── /src            # Código fonte
│   ├── /public         # Assets públicos
│   └── /components     # Componentes Vue
└── /docs               # Documentação
```

### Executando em desenvolvimento

[Instruções para execução em ambiente de desenvolvimento serão adicionadas]

## Roadmap Futuro

- Autenticação de usuários
- Monitoramento de recursos do servidor (CPU, memória, disco)
- Gerenciamento de containers (iniciar, parar, reiniciar)
- Visualização de logs
- Notificações de status
- Temas personalizáveis

## Contribuição

Este projeto está em fase inicial de desenvolvimento. Contribuições são bem-vindas!