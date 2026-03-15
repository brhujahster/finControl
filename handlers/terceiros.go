package handlers

import (
    "fincontrol/models"
    "html/template"
    "net/http"
    "strconv"
    "strings"
)

var terceirosTmpl = template.Must(template.ParseFiles(
    "templates/base.html",
    "templates/terceiros/lista.html",
))

var terceirosFormTmpl = template.Must(template.ParseFiles(
    "templates/base.html",
    "templates/terceiros/form.html",
))

func TerceirosIndex(w http.ResponseWriter, r *http.Request) {
    terceiros, err := models.ListarTerceiros()
    if err != nil {
        renderErro(w, "Erro ao listar terceiros", err.Error(), "/", http.StatusInternalServerError)
        return
    }
    terceirosTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
        "Terceiros": terceiros,
    })
}

func TerceirosNovo(w http.ResponseWriter, r *http.Request) {
    terceirosFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
        "Terceiro": nil,
        "Titulo":   "Novo Terceiro",
        "Action":   "/terceiros",
    })
}

func TerceirosSalvar(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        renderErro(w, "Formulário inválido", err.Error(), "/terceiros", http.StatusBadRequest)
        return
    }

    terceiro, err := terceiroFromForm(r)
    if err != nil {
        terceirosFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
            "Terceiro": nil,
            "Titulo":   "Novo Terceiro",
            "Action":   "/terceiros",
            "Erro":     "Dados inválidos: " + err.Error(),
        })
        return
    }

    if err := models.CriarTerceiro(terceiro); err != nil {
        renderErro(w, "Erro ao salvar terceiro", err.Error(), "/terceiros", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/terceiros", http.StatusSeeOther)
}

func TerceirosEditar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/terceiros/", "/editar")
    if err != nil {
        renderErro(w, "ID inválido", "O identificador do terceiro é inválido.", "/terceiros", http.StatusBadRequest)
        return
    }

    terceiro, err := models.BuscarTerceiro(id)
    if err != nil {
        renderErro(w, "Terceiro não encontrado", "O terceiro solicitado não foi encontrado.", "/terceiros", http.StatusNotFound)
        return
    }

    terceirosFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
        "Terceiro": terceiro,
        "Titulo":   "Editar Terceiro",
        "Action":   "/terceiros/" + strconv.Itoa(id),
    })
}

func TerceirosAtualizar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/terceiros/", "")
    if err != nil {
        renderErro(w, "ID inválido", "O identificador do terceiro é inválido.", "/terceiros", http.StatusBadRequest)
        return
    }

    if err := r.ParseForm(); err != nil {
        renderErro(w, "Formulário inválido", err.Error(), "/terceiros", http.StatusBadRequest)
        return
    }

    terceiro, err := terceiroFromForm(r)
    if err != nil {
        terceirosFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
            "Terceiro": nil,
            "Titulo":   "Editar Terceiro",
            "Action":   "/terceiros/" + strconv.Itoa(id),
            "Erro":     "Dados inválidos: " + err.Error(),
        })
        return
    }
    terceiro.ID = id

    if err := models.AtualizarTerceiro(terceiro); err != nil {
        renderErro(w, "Erro ao atualizar terceiro", err.Error(), "/terceiros", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/terceiros", http.StatusSeeOther)
}

func TerceirosDeletar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/terceiros/", "/deletar")
    if err != nil {
        renderErro(w, "ID inválido", "O identificador do terceiro é inválido.", "/terceiros", http.StatusBadRequest)
        return
    }

    temDespesas, err := models.TerceiroTemDespesas(id)
    if err != nil {
        renderErro(w, "Erro ao verificar dependências", err.Error(), "/terceiros", http.StatusInternalServerError)
        return
    }
    if temDespesas {
        renderErro(w, "Terceiro em uso", "Não é possível excluir este terceiro pois existem despesas associadas a ele.", "/terceiros", http.StatusConflict)
        return
    }

    temEmprestimos, err := models.TerceiroTemEmprestimos(id)
    if err != nil {
        renderErro(w, "Erro ao verificar dependências", err.Error(), "/terceiros", http.StatusInternalServerError)
        return
    }
    if temEmprestimos {
        renderErro(w, "Terceiro em uso", "Não é possível excluir este terceiro pois existem empréstimos associados a ele.", "/terceiros", http.StatusConflict)
        return
    }

    if err := models.DeletarTerceiro(id); err != nil {
        renderErro(w, "Erro ao deletar terceiro", err.Error(), "/terceiros", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/terceiros", http.StatusSeeOther)
}

func terceiroFromForm(r *http.Request) (models.Terceiro, error) {
    limite, err := strconv.ParseFloat(strings.ReplaceAll(r.FormValue("limite_liberado"), ",", "."), 64)
    if err != nil {
        return models.Terceiro{}, err
    }

    t := models.Terceiro{
        Nome:           r.FormValue("nome"),
        LimiteLiberado: limite,
    }

    if obs := r.FormValue("observacao"); obs != "" {
        t.Observacao = &obs
    }

    return t, nil
}
