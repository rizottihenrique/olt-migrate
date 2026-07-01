# OLT Migrate

**OLT Migrate** é uma plataforma corporativa projetada para automatizar e simplificar a migração de configurações de OLTs legadas (Optical Line Terminals) para novas infraestruturas.

A aplicação é capaz de ler arquivos de backup de OLTs antigas (como Fiberhome séries 5000 e 6000), mapear as portas lógicas e físicas (Slots e PONs) usando regras personalizadas, e gerar scripts de provisionamento prontos para serem inseridos nas novas OLTs (como Nokia ISAM FX 7360).

## 🚀 Tecnologias

- **Backend:** Go 1.22+ (Rotas nativas `net/http`, Parser de Regex)
- **Frontend:** Next.js 15, React, TailwindCSS v4
- **Arquitetura:** Padrão Adapter (Multi-Vendor), Design Minimalista no Frontend

## 📦 Como rodar localmente

### 1. Iniciar o Backend (Go API)
O backend processa os arquivos em memória usando Regex e o padrão Factory/Adapter.
```bash
cd backend
go mod tidy
go run cmd/main.go
```
A API ficará disponível em `http://localhost:8080/api/migrate`.

### 2. Iniciar o Frontend (Next.js)
A interface é uma aplicação web focada na experiência do usuário.
```bash
cd frontend
npm install
npm run dev
```
Acesse `http://localhost:3000` no seu navegador.

## ✨ Principais Funcionalidades

- **Multi-Vendor Engine**: Estrutura pronta para plugar novas marcas de OLT a qualquer momento (via `core.Parser` e `core.Generator`).
- **Resolução de Conflitos**: Aglutinação de várias PONs velhas em uma PON nova re-organizando os `OnuIDs` automaticamente para evitar colisões.
- **Detecção Inteligente (Bridge vs Router)**: Identifica automaticamente se a ONU era Bridge ou PPPoE/Router com base nos parâmetros capturados (como PPPoE User/Pass) e gera os comandos Nokia ISAM (`veip` vs `ethernet`) de acordo.
- **Interface Minimalista**: Sem distrações, upload rápido de arquivo e resultado na tela com 1 clique.
