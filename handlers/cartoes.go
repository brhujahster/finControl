package handlers

import (
    "fincontrol/models"
    "html/template"
    "net/http"
    "strconv"
    "strings"
)

var cartoesTmpl = template.Must(template.ParseFiles(
    "templates/base.html",
    "templates/cartoes/lista.html",
))

var cartoesFormTmpl = template.Must(template.ParseFiles(
    "templates/base.html",
    "templates/cartoes/form.html",
))

func CartoesIndex(w http.ResponseWriter, r *http.Request) {
    cartoes, err := models.ListarCartoes()
    if err != nil {
        http.Error(w, "Erro ao listar cartões: "+err.Error(), http.StatusInternalServerError)
        return
    }
    cartoesTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
        "Cartoes": cartoes,
    })
}

func CartoesNovo(w http.ResponseWriter, r *http.Request) {
    cartoesFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
        "Cartao": nil,
        "Titulo": "Novo Cartão",
        "Action": "/cartoes",
    })
}

func CartoesSalvar(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Formulário inválido", http.StatusBadRequest)
        return
    }

    cartao, err := cartaoFromForm(r)
    if err != nil {
        http.Error(w, "Dados inválidos: "+err.Error(), http.StatusBadRequest)
        return
    }

    if err := models.CriarCartao(cartao); err != nil {
        http.Error(w, "Erro ao salvar cartão: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/cartoes", http.StatusSeeOther)
}

func CartoesEditar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/cartoes/", "/editar")
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    cartao, err := models.BuscarCartao(id)
    if err != nil {
        http.Error(w, "Cartão não encontrado", http.StatusNotFound)
        return
    }

    cartoesFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
        "Cartao": cartao,
        "Titulo": "Editar Cartão",
        "Action": "/cartoes/" + strconv.Itoa(id),
    })
}

func CartoesAtualizar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/cartoes/", "")
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Formulário inválido", http.StatusBadRequest)
        return
    }

    cartao, err := cartaoFromForm(r)
    if err != nil {
        http.Error(w, "Dados inválidos: "+err.Error(), http.StatusBadRequest)
        return
    }
    cartao.ID = id

    if err := models.AtualizarCartao(cartao); err != nil {
        http.Error(w, "Erro ao atualizar cartão: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/cartoes", http.StatusSeeOther)
}

func CartoesDeletar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/cartoes/", "/deletar")
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    if err := models.DeletarCartao(id); err != nil {
        http.Error(w, "Erro ao deletar cartão: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/cartoes", http.StatusSeeOther)
}

func cartaoFromForm(r *http.Request) (models.Cartao, error) {
    limite, err := strconv.ParseFloat(strings.ReplaceAll(r.FormValue("limite"), ",", "."), 64)
    if err != nil {
        return models.Cartao{}, err
    }

    diaVencimento, err := strconv.Atoi(r.FormValue("dia_vencimento"))
    if err != nil {
        return models.Cartao{}, err
    }

    return models.Cartao{
        Nome:          r.FormValue("nome"),
        DiaVencimento: diaVencimento,
        Limite:        limite,
    }, nil
}