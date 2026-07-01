# Segurança (Security Practices)

Como o OLT Migrate lida com dados sensíveis de infraestrutura de rede (configurações completas, mac addresses de clientes, senhas de PPPoE em alguns casos), as seguintes diretrizes de segurança são adotadas:

## 1. Retenção de Dados e Privacidade
Por padrão, a aplicação na Fase 1 opera em memória (stateless).
- Os arquivos de configuração de origem (`.txt`, `.dat`, `.xml`) **não são salvos em disco permanentemente**.
- O backend processa o arquivo no momento da requisição HTTP, extrai o modelo e devolve os comandos. A memória é liberada pelo Garbage Collector do Go.
- Não guardamos as senhas PPPoE dos clientes que possam vir exportadas, a menos que seja estritamente necessário para gerar o comando do destino.

## 2. Comunicação Segura
- Quando hospedado externamente, o acesso ao frontend e à API Go deve ser exclusivamente feito através de HTTPS (TLS 1.2+).
- CORS (Cross-Origin Resource Sharing) da API Go deve ser estritamente configurado para aceitar requisições apenas do domínio do frontend.

## 3. Sanitização de Entrada
- O parser em Go que lê os arquivos da OLT deve prever arquivos muito grandes e incluir limites de tamanho (Rate Limiting e Size Limiting no Body Parser) para evitar ataques de negação de serviço (DoS) de exaustão de memória.

## 4. Auditoria
- Na Fase 2 e 3 (quando introduzirmos banco de dados), senhas e dados confidenciais exportados pelas OLTs antigas devem ser criptografados no banco (ex: AES-256-GCM) antes de persistidos, para garantir o compliance com a LGPD.
