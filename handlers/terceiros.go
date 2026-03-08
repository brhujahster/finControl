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
        http.Error(w, "Erro ao listar terceiros: "+err.Error(), http.StatusInternalServerError)
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
        http.Error(w, "Formulário inválido", http.StatusBadRequest)
        return
    }

    terceiro, err := terceiroFromForm(r)
    if err != nil {
        http.Error(w, "Dados inválidos: "+err.Error(), http.StatusBadRequest)
        return
    }

    if err := models.CriarTerceiro(terceiro); err != nil {
        http.Error(w, "Erro ao salvar terceiro: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/terceiros", http.StatusSeeOther)
}

func TerceirosEditar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/terceiros/", "/editar")
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    terceiro, err := models.BuscarTerceiro(id)
    if err != nil {
        http.Error(w, "Terceiro não encontrado", http.StatusNotFound)
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
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Formulário inválido", http.StatusBadRequest)
        return
    }

    terceiro, err := terceiroFromForm(r)
    if err != nil {
        http.Error(w, "Dados inválidos: "+err.Error(), http.StatusBadRequest)
        return
    }
    terceiro.ID = id

    if err := models.AtualizarTerceiro(terceiro); err != nil {
        http.Error(w, "Erro ao atualizar terceiro: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/terceiros", http.StatusSeeOther)
}

func TerceirosDeletar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/terceiros/", "/deletar")
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    if err := models.DeletarTerceiro(id); err != nil {
        http.Error(w, "Erro ao deletar terceiro: "+err.Error(), http.StatusInternalServerError)
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