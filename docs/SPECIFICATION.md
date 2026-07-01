# Especificação do Sistema

Este documento detalha as especificações técnicas e regras de negócio para a ferramenta.

## 1. Regras de Negócio

### Leitura da Origem (Fiberhome)
1. O sistema deve suportar as sintaxes de exportação padrão do equipamento Fiberhome (Linha 5000/6000).
2. O sistema deve ignorar os comandos irrelevantes para a migração (como comandos de fan, ambiente, ou rotas de gerência), focando apenas no provisionamento das ONUs, LineProfiles, SrvProfiles e Service Ports/VLANs.

### Mapeamento
1. O sistema precisa permitir que o usuário defina que as ONUs do `Slot 1 / PON 1` da OLT Antiga irão fisicamente para o `Slot 3 / PON 2` da OLT Nova (Nokia).
2. Deve suportar mapeamento N para 1 (Duas portas PONs antigas sub-utilizadas podem ir para a mesma PON na Nokia, desde que não exceda o limite de slots lógicos - ex: 128 por porta).

### Geração para Nokia (FX 7360)
1. Os comandos gerados devem estar na sintaxe correta e mais atual da CLI Nokia (ISAM).
2. O arquivo gerado de saída deve estar formatado, separando blocos de provisionamento de profiles globais, do provisionamento porta a porta, finalizando com as associações de serviços (VLAN).

## 2. Requisitos Não Funcionais (NFRs)

- **Performance:** O parser deve conseguir ler um arquivo de 50MB (milhares de comandos Fiberhome) em menos de 2 segundos de processamento de CPU. Para isso, o Go será fundamental com sua alta concorrência.
- **Portabilidade:** A aplicação web deve ser responsiva, mas focada primordialmente para uso em monitores desktop de engenheiros de redes (1080p+).
