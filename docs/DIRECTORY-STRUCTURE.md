# Estrutura de Diretórios

O projeto utiliza um padrão de monorepo separando os dois ecossistemas (Go e Node).

```
d:\Dev\olt-migrate
│
├── backend/                   # Aplicação Go (API REST e Lógica de Domínio)
│   ├── cmd/
│   │   └── main.go            # Ponto de entrada e servidor HTTP
│   ├── internal/
│   │   ├── core/              # Interfaces de Adapter e Factory (Parser/Generator)
│   │   ├── generator/         # Scripts de saída (ex: nokia.go)
│   │   ├── mapper/            # Transposição de portas e resolução de conflitos
│   │   ├── models/            # Structs de dados
│   │   └── parser/            # Processamento de arquivos legados (ex: fiberhome.go)
│   ├── go.mod
│   └── go.sum
│
├── frontend/                  # Aplicação Web (Next.js)
│   ├── src/
│   │   └── app/
│   │       ├── globals.css    # Estilização base do Tailwind (Light Theme Minimalist)
│   │       ├── layout.tsx     # Definição de rotas, fontes e Navbar
│   │       ├── page.tsx       # Componente principal do App de Migração
│   │       └── icon.svg       # Favicon dinâmico
│   ├── package.json
│   ├── tsconfig.json
│   └── next.config.mjs
│
└── docs/                      # Documentação técnica e de arquitetura
    ├── ARCHITECTURE.md
    ├── CHANGELOG.md
    ├── DIRECTORY-STRUCTURE.md
    ├── MODELS.md
    ├── README.md
    └── ref/                   # Backups de exemplo (5000 e 6000 series)
```
