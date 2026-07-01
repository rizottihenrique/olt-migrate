# Changelog

Todas as mudanças notáveis deste projeto serão documentadas neste arquivo.

## [v1.0.0-alpha] - 2026-06-30

### Adicionado
- **Motor de Parser Fiberhome**: Implementação de leitura inteligente com Regex para abranger OLTs linha 5000 (comandos `set`) e linha 6000 (comandos baseados em contexto hierárquico `interface`).
- **Detecção de Autenticação**: Extração nativa de PPPoE/Senhas.
- **Mapeador Inteligente (`mapper`)**: Ferramenta de roteamento visual para realocar Slot X/Pon Y para Slot Z/Pon W.
- **Gerador Nokia ISAM**: Geração de CLI prontas (`configure equipment ont interface...`).
- **Templates Dinâmicos**: O sistema passa a injetar `veip` se o parser identificar parâmetros PPPoE (Router) e `ethernet`/pvid caso não identifique (Bridge).
- **Interface Minimalista Web**: Novo frontend Next.js 15 Tailwind v4 limpo, sem componentes falsos, estritamente funcional e utilitário para técnicos de rede.
- **Feature de Download**: Geração e download imediato do script `.txt` de migração.
- **Arquitetura Multi-Vendor**: Refatoração do backend para padrão de Interfaces (`core.Parser` e `core.Generator`), possibilitando integração rápida de novas marcas na plataforma além de Fiberhome e Nokia.

### Modificado
- Design Frontend: Migrado do layout noturno glassmorphism para um visual técnico corporativo fundo branco (Light mode forçado).
- Atualização do ícone dinâmico do Header e Favicon `.svg` nativos do App Router Next.js.
