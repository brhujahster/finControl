package handlers

import (
	"fincontrol/models"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var emprestimosTmpl = template.Must(template.New("base.html").Funcs(funcMap).ParseFiles(
	"templates/base.html",
	"templates/emprestimos/lista.html",
))

var emprestimosFormTmpl = template.Must(template.New("base.html").Funcs(funcMap).ParseFiles(
	"templates/base.html",
	"templates/emprestimos/form.html",
))

func EmprestimosIndex(w http.ResponseWriter, r *http.Request) {
	emprestimos, err := models.ListarEmprestimos()
	if err != nil {
		renderErro(w, "Erro ao listar empréstimos", err.Error(), "/", http.StatusInternalServerError)
		return
	}
	emprestimosTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Emprestimos": emprestimos,
	})
}

func EmprestimosNovo(w http.ResponseWriter, r *http.Request) {
	cartoes, _ := models.ListarCartoes()
	terceiros, _ := models.ListarTerceiros()
	emprestimosFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Emprestimo": nil,
		"Titulo":     "Novo Empréstimo",
		"Action":     "/emprestimos",
		"Cartoes":    cartoes,
		"Terceiros":  terceiros,
	})
}

func EmprestimosSalvar(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderErro(w, "Formulário inválido", err.Error(), "/emprestimos", http.StatusBadRequest)
		return
	}

	emprestimo, err := emprestimoFromForm(r)
	if err != nil {
		cartoes, _ := models.ListarCartoes()
		terceiros, _ := models.ListarTerceiros()
		emprestimosFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
			"Emprestimo": nil,
			"Titulo":     "Novo Empréstimo",
			"Action":     "/emprestimos",
			"Cartoes":    cartoes,
			"Terceiros":  terceiros,
			"Erro":       "Dados inválidos: " + err.Error(),
		})
		return
	}

	diaVencimento := 0
	if emprestimo.CartaoID != nil {
		cartao, err := models.BuscarCartao(*emprestimo.CartaoID)
		if err == nil {
			diaVencimento = cartao.DiaVencimento
		}
	}

	if err := models.CriarEmprestimo(emprestimo, diaVencimento); err != nil {
		renderErro(w, "Erro ao salvar empréstimo", err.Error(), "/emprestimos", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/emprestimos", http.StatusSeeOther)
}

func EmprestimosEditar(w http.ResponseWriter, r *http.Request) {
	id, err := extrairID(r.URL.Path, "/emprestimos/", "/editar")
	if err != nil {
		renderErro(w, "ID inválido", "O identificador do empréstimo é inválido.", "/emprestimos", http.StatusBadRequest)
		return
	}

	emprestimo, err := models.BuscarEmprestimo(id)
	if err != nil {
		renderErro(w, "Empréstimo não encontrado", "O empréstimo solicitado não foi encontrado.", "/emprestimos", http.StatusNotFound)
		return
	}

	cartoes, _ := models.ListarCartoes()
	terceiros, _ := models.ListarTerceiros()
	emprestimosFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Emprestimo": emprestimo,
		"Titulo":     "Editar Empréstimo",
		"Action":     "/emprestimos/" + strconv.Itoa(id),
		"Cartoes":    cartoes,
		"Terceiros":  terceiros,
	})
}

func EmprestimosAtualizar(w http.ResponseWriter, r *http.Request) {
	id, err := extrairID(r.URL.Path, "/emprestimos/", "")
	if err != nil {
		renderErro(w, "ID inválido", "O identificador do empréstimo é inválido.", "/emprestimos", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		renderErro(w, "Formulário inválido", err.Error(), "/emprestimos", http.StatusBadRequest)
		return
	}

	emprestimo, err := emprestimoFromForm(r)
	if err != nil {
		cartoes, _ := models.ListarCartoes()
		terceiros, _ := models.ListarTerceiros()
		emprestimosFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
			"Emprestimo": nil,
			"Titulo":     "Editar Empréstimo",
			"Action":     "/emprestimos/" + strconv.Itoa(id),
			"Cartoes":    cartoes,
			"Terceiros":  terceiros,
			"Erro":       "Dados inválidos: " + err.Error(),
		})
		return
	}
	emprestimo.ID = id

	if err := models.AtualizarEmprestimo(emprestimo); err != nil {
		renderErro(w, "Erro ao atualizar empréstimo", err.Error(), "/emprestimos", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/emprestimos", http.StatusSeeOther)
}

func EmprestimosDeletar(w http.ResponseWriter, r *http.Request) {
	id, err := extrairID(r.URL.Path, "/emprestimos/", "/deletar")
	if err != nil {
		renderErro(w, "ID inválido", "O identificador do empréstimo é inválido.", "/emprestimos", http.StatusBadRequest)
		return
	}

	if err := models.DeletarEmprestimo(id); err != nil {
		renderErro(w, "Erro ao deletar empréstimo", err.Error(), "/emprestimos", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/emprestimos", http.StatusSeeOther)
}

func emprestimoFromForm(r *http.Request) (models.Emprestimo, error) {
	valorTotal, err := strconv.ParseFloat(strings.ReplaceAll(r.FormValue("valor_total"), ",", "."), 64)
	if err != nil {
		return models.Emprestimo{}, err
	}

	totalParcelas, err := strconv.Atoi(r.FormValue("total_parcelas"))
	if err != nil {
		return models.Emprestimo{}, err
	}

	dataInicio, err := time.Parse("2006-01-02", r.FormValue("data_inicio"))
	if err != nil {
		return models.Emprestimo{}, err
	}

	e := models.Emprestimo{
		Descricao:      r.FormValue("descricao"),
		ValorTotal:     valorTotal,
		TotalParcelas:  totalParcelas,
		DataInicio:     dataInicio,
		FormaPagamento: r.FormValue("forma_pagamento"),
	}

	if v := r.FormValue("taxa_juros"); v != "" {
		f, err := strconv.ParseFloat(strings.ReplaceAll(v, ",", "."), 64)
		if err != nil {
			return models.Emprestimo{}, err
		}
		e.TaxaJuros = &f
		tj := r.FormValue("tipo_juros")
		e.TipoJuros = &tj
	}

	if v := r.FormValue("cartao_id"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return models.Emprestimo{}, err
		}
		e.CartaoID = &n
	}

	if v := r.FormValue("terceiro_id"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return models.Emprestimo{}, err
		}
		e.TerceiroID = &n
	}

	if v := r.FormValue("observacao"); v != "" {
		e.Observacao = &v
	}

	return e, nil
}
