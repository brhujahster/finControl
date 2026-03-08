# PRD – Serviço de Controle Orçamentário Pessoal

## 1. Visão Geral do Produto

**Nome do Produto:** FinControl  
**Versão:** 1.0  
**Stack Técnica:** Go + GoHTML + SQLite3  
**Tipo:** Aplicação Web Local (self-hosted)  
**Público-alvo:** Usuário individual que deseja controlar suas finanças pessoais, incluindo receitas, despesas, cartões de crédito, terceiros e empréstimos.

## 2. Objetivos do Produto

- Centralizar o controle financeiro pessoal em uma única aplicação leve e offline-first.
- Oferecer visibilidade clara sobre receitas, despesas, dívidas de terceiros e empréstimos.
- Permitir planejamento mensal com base em receitas recorrentes e pontuais.
- Facilitar a cobrança de terceiros que realizaram compras no cartão do usuário.

## 3. Escopo da Versão 1.0

### Incluído
- Cadastro de receitas (recorrentes e pontuais)
- Cadastro de cartões de crédito
- Cadastro de terceiros
- Cadastro de despesas (cartão parcelado, dinheiro, fixas)
- Cadastro de empréstimos
- Dashboard financeiro mensal

### Fora do Escopo (v1.0)
- Autenticação multi-usuário
- Integração com bancos/Open Finance
- Aplicativo mobile
- Exportação para PDF/Excel (planejado para v1.1)

## 4. Requisitos Funcionais

### 4.1 Módulo de Receitas

**Descrição:** O usuário deve ser capaz de registrar os valores que terá disponíveis em cada mês.

**Regras de Negócio:**
- Uma receita pode ser **recorrente** (se repete todo mês indefinidamente ou até uma data de encerramento) ou **pontual** (válida apenas para o mês/ano informado).
- Receitas recorrentes devem ser projetadas automaticamente nos meses futuros até a data de encerramento (ou indefinidamente).
- O valor total de receitas do mês é a soma de todas as receitas recorrentes ativas + receitas pontuais do mês.

**Campos:**

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `id` | INTEGER PK | Sim | Identificador único |
| `descricao` | TEXT | Sim | Ex: "Salário", "Freelance" |
| `valor` | REAL | Sim | Valor em reais |
| `tipo` | ENUM | Sim | `recorrente` ou `pontual` |
| `mes_referencia` | DATE | Condicional | Obrigatório se pontual (mês/ano) |
| `data_inicio` | DATE | Condicional | Obrigatório se recorrente |
| `data_fim` | DATE | Não | Data de encerramento da recorrência |
| `criado_em` | DATETIME | Sim | Timestamp de criação |

### 4.2 Módulo de Cartões de Crédito

**Descrição:** O usuário deve cadastrar seus cartões de crédito para associá-los às despesas.

**Regras de Negócio:**
- Um cartão possui limite total. O sistema deve calcular o **limite utilizado** com base nas despesas associadas ao cartão no mês de vencimento da fatura.
- O **dia de vencimento** define em qual mês a fatura será cobrada.
- O sistema deve alertar quando o limite do cartão for ultrapassado.

**Campos:**

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `id` | INTEGER PK | Sim |
| `nome` | TEXT | Sim |
| `dia_vencimento` | INTEGER | Sim |
| `limite` | REAL | Sim |
| `criado_em` | DATETIME | Sim |

### 4.3 Módulo de Terceiros

**Descrição:** Pessoas que realizaram compras no cartão do usuário e que devem ser cobradas posteriormente.

**Regras de Negócio:**
- Cada terceiro possui um **limite liberado**.
- O sistema deve calcular o **saldo devedor** do terceiro por mês.
- O sistema deve alertar quando o gasto do terceiro ultrapassar o limite liberado.
- Um terceiro pode ser associado a um empréstimo como responsável pelo pagamento.

**Campos:**

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `id` | INTEGER PK | Sim |
| `nome` | TEXT | Sim |
| `limite_liberado` | REAL | Sim |
| `observacao` | TEXT | Não |
| `criado_em` | DATETIME | Sim |

### 4.4 Módulo de Despesas

**Descrição:** Registro de todas as saídas financeiras do usuário.

**Regras de Negócio:**
- Uma despesa pode ser paga com **cartão de crédito** ou **dinheiro/débito**.
- Despesas no cartão podem ser **parceladas**.
- Uma despesa pode ser marcada como **fixa** (recorrente mensal).
- Uma despesa pode ser associada a um **terceiro**.
- O mês de competência de uma parcela de cartão é determinado pelo mês de vencimento da fatura.
- Despesas fixas em dinheiro repetem no mesmo dia todo mês.

**Campos:**

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `id` | INTEGER PK | Sim |
| `descricao` | TEXT | Sim |
| `valor_total` | REAL | Sim |
| `data_compra` | DATE | Sim |
| `forma_pagamento` | ENUM | Sim |
| `cartao_id` | INTEGER FK | Condicional |
| `parcelado` | BOOLEAN | Não |
| `total_parcelas` | INTEGER | Condicional |
| `parcela_atual` | INTEGER | Auto |
| `fixa` | BOOLEAN | Não |
| `terceiro_id` | INTEGER FK | Não |
| `categoria` | TEXT | Não |
| `observacao` | TEXT | Não |
| `criado_em` | DATETIME | Sim |

### 4.5 Módulo de Empréstimos

**Descrição:** Registro de empréstimos realizados pelo usuário.

**Regras de Negócio:**
- Um empréstimo possui valor total, número de parcelas e taxa de juros (opcional).
- Um terceiro pode ser designado como **responsável pelo pagamento**.
- Cada parcela do empréstimo deve aparecer como despesa no mês correspondente.
- O empréstimo pode ser associado a um cartão de crédito ou ser em dinheiro.

**Campos:**

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `id` | INTEGER PK | Sim |
| `descricao` | TEXT | Sim |
| `valor_total` | REAL | Sim |
| `total_parcelas` | INTEGER | Sim |
| `taxa_juros` | REAL | Não |
| `tipo_juros` | ENUM | Não |
| `data_inicio` | DATE | Sim |
| `forma_pagamento` | ENUM | Sim |
| `cartao_id` | INTEGER FK | Condicional |
| `terceiro_id` | INTEGER FK | Não |
| `observacao` | TEXT | Não |
| `criado_em` | DATETIME | Sim |

### 4.6 Dashboard Financeiro

**Descrição:** Visão consolidada das finanças do mês selecionado.

**Componentes:**
- Resumo geral (Receita Total, Total de Despesas, Saldo Disponível, % Comprometido).
- Despesas por forma de pagamento (por cartão e dinheiro).
- Terceiros devedores do mês, discriminando valores por cartão e dinheiro.
- Proporção dos gastos em relação à receita total.
- Empréstimos ativos no mês.
- Alertas de limite de cartão, limite de terceiro e saldo negativo.

## 5. Requisitos Não Funcionais

- Performance: carregamento das páginas < 500ms local.
- Portabilidade: binário único Go + SQLite.
- Usabilidade: interface simples e responsiva.
- Segurança: acesso local apenas (v1.0).

## 6. Arquitetura Técnica (Visão Geral)

Estrutura sugerida de pastas:

- `main.go`
- `db/` (conexão e migrations)
- `handlers/` (lógica HTTP)
- `models/` (modelos de dados)
- `templates/` (GoHTML)
- `static/` (CSS/JS)
