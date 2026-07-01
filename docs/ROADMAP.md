# Roadmap (Plano de Evolução)

Este documento descreve os passos e funcionalidades planejadas para o futuro da aplicação.

## Fase 1 (MVP - Produto Mínimo Viável) - Foco Atual
- [ ] Backend: Parser de backup da Fiberhome (AN5516-01 / Linhas 5k e 6k).
- [ ] Backend: Serviço de Mapeamento Básico (Slot/PON para Slot/PON).
- [ ] Backend: Gerador de comandos básicos para Nokia FX 7360 (Perfis básicos e cadastro de ONU/VLAN).
- [ ] Frontend: Dashboard para upload, interface de mapeamento e visualização de resultados.

## Fase 2 (Refinamento e Profiles)
- [ ] Identificação inteligente e replicação de `LineProfiles` e `ServiceProfiles`.
- [ ] Suporte a migração de configurações de multicast (IPTV) e voz (VoIP).
- [ ] Opção de validação pré-migração (Checar se as portas de destino suportam a quantidade de ONUs).

## Fase 3 (Integrações e Novos Vendors)
- [ ] Suporte à leitura de configurações de OLTs ZTE (C320/C300).
- [ ] Suporte à leitura de configurações de OLTs Huawei (MA5800/MA5608T).
- [ ] Banco de dados para log e retenção de arquivos.
- [ ] (Opcional) Integração via Telnet/SSH/NETCONF para injetar os comandos gerados automaticamente na OLT de destino.
