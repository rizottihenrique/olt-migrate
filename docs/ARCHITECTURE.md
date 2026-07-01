# Arquitetura do Sistema

O projeto OLT Migrate foi construído com foco em **escalabilidade** e **desacoplamento**, garantindo que a inserção de novas tecnologias de rede no futuro exija esforço mínimo de codificação.

## Padrão Multi-Vendor (Adapter Pattern)

A necessidade de ler configurações de múltiplos equipamentos (Fiberhome, ZTE, Datacom) e escrever em múltiplos destinos (Nokia, Huawei) nos levou à adoção do padrão Adapter.

Toda a lógica central se encontra no pacote `internal/core`, onde residem duas interfaces vitais:

```go
type Parser interface {
	Parse(file io.Reader) ([]models.ONU, error)
}

type Generator interface {
	Generate(onus []models.ONU) string
}
```

### O Fluxo da Aplicação

1. **Upload via Frontend**: O usuário insere o arquivo de backup e define no form as opções `sourceVendor` e `destVendor`.
2. **Registry**: Em `core/registry.go`, a aplicação atua como uma fábrica (Factory), instanciando as estruturas corretas de acordo com a seleção.
   - *Exemplo*: Se `sourceVendor == "Fiberhome"`, ele invoca o `&parser.FiberhomeParser{}`.
3. **Parse (Extração)**: O parser específico varre o arquivo bruto linha a linha usando Expressões Regulares (Regex) e agrupa os dados de rede espalhados em um array genérico de `models.ONU`.
4. **Mapping (Transposição)**: O pacote `internal/mapper` cruza os ONUs extraídos com as regras de "De-Para" fornecidas pelo front-end (Ex: Slot 1 -> Slot 3). Nessa fase, IDs conflitantes (`OnuID`) são resolvidos matematicamente.
5. **Generate (Produção)**: A lista mapeada é enviada para o Generator específico (ex: `&generator.NokiaGenerator{}`), que itera sobre o array e cria as strings do CLI final.

## Frontend UI/UX

Desenvolvido em Next.js (App Router), o frontend prioriza o design utilitário minimalista. 
A comunicação é feita de forma assíncrona com o backend via Fetch API em POST multipart, enviando os arquivos de texto diretamente para o Go processar em memória, não havendo persistência em banco de dados por questões de LGPD e segurança de dados dos provedores.
