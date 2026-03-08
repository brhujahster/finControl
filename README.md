# FinControl

Sistema web local para controle orçamentário pessoal. Desenvolvido com Go, GoHTML e SQLite3.

## Tecnologias

- **Go** – servidor HTTP com `net/http`
- **SQLite** – banco de dados local via `modernc.org/sqlite` (sem necessidade de GCC)
- **GoHTML** – templates HTML server-side
- **Tailwind CSS** – estilização via CDN

## Funcionalidades

- Cadastro de **receitas** recorrentes e pontuais
- Cadastro de **cartões de crédito** com controle de limite
- Cadastro de **terceiros** com limite liberado e saldo devedor
- Cadastro de **despesas** (cartão parcelado, dinheiro, fixas)
- Cadastro de **empréstimos** com parcelas e taxas de juros
- **Dashboard** financeiro mensal consolidado

## Pré-requisitos

- [Go 1.21+](https://golang.org/dl/)

## Instalação e Execução

```bash
# Clone o repositório
git clone https://github.com/seu-usuario/fincontrol.git
cd fincontrol

# Instale as dependências
go mod tidy

# Execute o servidor
go run main.go

fincontrol/
├── main.go               # Servidor HTTP e rotas
├── db/
│   ├── db.go             # Conexão com SQLite
│   └── migrations.go     # Criação das tabelas
├── handlers/             # Lógica dos endpoints HTTP
├── models/               # Structs e queries do banco
├── templates/            # Templates GoHTML
│   ├── base.html
│   └── receitas/
│       ├── lista.html
│       └── form.html
└── static/               # Arquivos estáticos (CSS/JS)

