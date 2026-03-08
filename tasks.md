# Tasks – FinControl v1.0

Tarefas sequenciais para construção do sistema. Cada tarefa deve ser concluída antes de iniciar a próxima.

---

## Tarefa 1 – Estrutura do Projeto

- Criar o diretório raiz do projeto Go com `go mod init`
- Criar a estrutura de pastas: `db/`, `handlers/`, `models/`, `templates/`, `static/`
- Criar o arquivo `main.go` com servidor HTTP básico (net/http) escutando em `localhost:8080`
- Adicionar dependência do driver SQLite3 (`github.com/mattn/go-sqlite3`)

---

## Tarefa 2 – Banco de Dados: Conexão e Migrations

- Criar `db/db.go` com função de conexão ao SQLite
- Criar `db/migrations.go` com a função que executa os `CREATE TABLE IF NOT EXISTS` para todas as entidades
- Tabelas a criar nesta tarefa:
  - `receitas`
  - `cartoes`
  - `terceiros`
  - `despesas`
  - `emprestimos`
- Chamar a migration no startup do `main.go`

---

## Tarefa 3 – Model e CRUD: Receitas

- Criar `models/receita.go` com a struct `Receita` e todos os seus campos conforme o PRD
- Criar `handlers/receitas.go` com os handlers HTTP:
  - `GET /receitas` – listar receitas
  - `GET /receitas/nova` – formulário de criação
  - `POST /receitas` – salvar nova receita
  - `GET /receitas/{id}/editar` – formulário de edição
  - `POST /receitas/{id}` – atualizar receita
  - `POST /receitas/{id}/deletar` – remover receita
- Criar templates `templates/receitas/lista.html`, `form.html`
- Implementar regra: receita recorrente projeta nos meses futuros; pontual só no mês informado

---

## Tarefa 4 – Model e CRUD: Cartões de Crédito

- Criar `models/cartao.go` com a struct `Cartao`
- Criar `handlers/cartoes.go` com os handlers HTTP:
  - `GET /cartoes` – listar cartões
  - `GET /cartoes/novo` – formulário de criação
  - `POST /cartoes` – salvar novo cartão
  - `GET /cartoes/{id}/editar` – formulário de edição
  - `POST /cartoes/{id}` – atualizar cartão
  - `POST /cartoes/{id}/deletar` – remover cartão
- Criar templates `templates/cartoes/lista.html`, `form.html`

---

## Tarefa 5 – Model e CRUD: Terceiros

- Criar `models/terceiro.go` com a struct `Terceiro`
- Criar `handlers/terceiros.go` com os handlers HTTP:
  - `GET /terceiros` – listar terceiros
  - `GET /terceiros/novo` – formulário de criação
  - `POST /terceiros` – salvar novo terceiro
  - `GET /terceiros/{id}/editar` – formulário de edição
  - `POST /terceiros/{id}` – atualizar terceiro
  - `POST /terceiros/{id}/deletar` – remover terceiro
- Criar templates `templates/terceiros/lista.html`, `form.html`

---

## Tarefa 6 – Model e CRUD: Despesas

- Criar `models/despesa.go` com a struct `Despesa`
- Criar `handlers/despesas.go` com os handlers HTTP:
  - `GET /despesas` – listar despesas (filtro por mês)
  - `GET /despesas/nova` – formulário de criação
  - `POST /despesas` – salvar nova despesa
  - `GET /despesas/{id}/editar` – formulário de edição
  - `POST /despesas/{id}` – atualizar despesa
  - `POST /despesas/{id}/deletar` – remover despesa
- Criar templates `templates/despesas/lista.html`, `form.html`
- Implementar regras:
  - Despesas parceladas no cartão geram N parcelas com mês de vencimento calculado pelo `dia_vencimento` do cartão
  - Despesas fixas em dinheiro repetem todo mês no mesmo dia
  - Despesa pode ser vinculada a terceiro

---

## Tarefa 7 – Model e CRUD: Empréstimos

- Criar `models/emprestimo.go` com a struct `Emprestimo`
- Criar `handlers/emprestimos.go` com os handlers HTTP:
  - `GET /emprestimos` – listar empréstimos
  - `GET /emprestimos/novo` – formulário de criação
  - `POST /emprestimos` – salvar novo empréstimo
  - `GET /emprestimos/{id}/editar` – formulário de edição
  - `POST /emprestimos/{id}` – atualizar empréstimo
  - `POST /emprestimos/{id}/deletar` – remover empréstimo
- Criar templates `templates/emprestimos/lista.html`, `form.html`
- Implementar regras:
  - Calcular parcelas com juros simples ou composto conforme `tipo_juros`
  - Cada parcela aparece como despesa no mês correspondente
  - Pode vincular terceiro como responsável pelo pagamento

---

## Tarefa 8 – Layout Base e Navegação

- Criar `templates/layout.html` com estrutura HTML base (header, nav, footer)
- Criar menu de navegação com links para: Dashboard, Receitas, Cartões, Terceiros, Despesas, Empréstimos
- Criar `static/css/style.css` com estilos básicos responsivos
- Garantir que todos os templates das tarefas anteriores herdam o layout base

---

## Tarefa 9 – Dashboard Financeiro

- Criar `handlers/dashboard.go` com handler `GET /` que recebe parâmetro de mês/ano (padrão: mês atual)
- Criar `templates/dashboard.html`
- Implementar os componentes do dashboard:
  - **Resumo geral:** Receita Total, Total de Despesas, Saldo Disponível, % Comprometido
  - **Despesas por forma de pagamento:** separadas por cartão e dinheiro/débito
  - **Terceiros devedores do mês:** valor gasto por terceiro discriminado por cartão e dinheiro
  - **Proporção dos gastos** em relação à receita total
  - **Empréstimos ativos** no mês selecionado
- Implementar seletor de mês/ano para navegação entre períodos

---

## Tarefa 10 – Sistema de Alertas

- Implementar lógica de alertas no dashboard:
  - Alerta quando o limite de um cartão for ultrapassado pelas despesas do mês
  - Alerta quando o gasto de um terceiro ultrapassar o `limite_liberado`
  - Alerta quando o saldo disponível do mês for negativo
- Exibir alertas no topo do dashboard de forma visível

---

## Tarefa 11 – Validações e Tratamento de Erros

- Validar campos obrigatórios em todos os formulários (front-end com HTML5 e back-end nos handlers)
- Exibir mensagens de erro amigáveis ao usuário
- Garantir integridade referencial: impedir exclusão de cartão com despesas associadas, terceiro com despesas, etc.
- Tratar erros de banco de dados e exibir página de erro genérica

---

## Tarefa 12 – Testes e Ajustes Finais

- Testar todos os fluxos CRUD de cada módulo
- Testar cálculo de parcelas de despesas e empréstimos
- Testar projeção de receitas recorrentes em meses futuros
- Testar todos os alertas do dashboard
- Verificar responsividade da interface
- Ajustar performance para carregamento < 500ms
- Compilar binário final e verificar execução standalone com o SQLite embutido
