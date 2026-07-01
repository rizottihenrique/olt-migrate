# Estratégia de Testes

Para garantir que a migração dos clientes seja perfeita (já que comandos errados derrubam a internet de milhares de usuários), a testabilidade é prioritária.

## 1. Testes do Backend (Go)

- **Testes Unitários de Parsing:** Serão escritos testes para o `parser`. Arquivos contendo trechos simulados (Mocks) da CLI Fiberhome devem ser processados e assertados (`assert.Equal`) contra a estrutura `models.ONU` esperada.
- **Testes Unitários de Geração:** A saída do gerador Nokia será testada verificando a presença das strings corretas para a criação de um serviço com base no modelo interno.
- **Ferramenta:** `go test ./...`

## 2. Testes de Mapeamento (Edge Cases)

Devem existir casos de testes para:
1. Conflitos de ONUIDs: Se eu migrar duas PONs Fiberhome para a mesma PON Nokia, e ambas têm a ONU ID 1, o sistema deve automaticamente remapear a segunda ONU ID para o próximo número livre (ex: ID 2).

## 3. Testes do Frontend (Next.js)

- **Testes de Componente:** Verificação do `<MappingUI />` para garantir que as linhas de mapeamento estão gerando o payload JSON correto para a API.
- **Ferramenta:** `Jest` + `React Testing Library`.
