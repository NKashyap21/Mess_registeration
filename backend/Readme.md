
mess-registration/
│── build/                   # Build artifacts (ignored in git usually)
│
│── config/                  # Configurations
│   ├── config.go
│   ├── db.go
│   └── env.go
│
│── internal/
│   ├── controller/          # Gin Handlers (API layer)
│   ├── services/            # Business logic
│   ├── repository/          # Database queries
│   ├── models/              # Structs (DB + API request/response)
│   ├── router/              # Route definitions
│   ├── middlewares/         # Auth, logging, recovery
│   └── utils/               # Helper functions
│
│── migrations/              # SQL migrations (or goose, migrate)
│   ├── 001_init.sql
│   └── 002_add_registration.sql
│
│── logs/                    # App logs
│
│── Dockerfile               # Go API docker image
│── docker-compose.yml       # API + Postgres setup
│── go.mod
│── go.sum
│── main.go                  # Entry point
│── Makefile                 # For build/test commands
│── README.md
