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
        renderErro(w, "Erro ao listar cartões", err.Error(), "/", http.StatusInternalServerError)
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
        renderErro(w, "Formulário inválido", err.Error(), "/cartoes", http.StatusBadRequest)
        return
    }

    cartao, err := cartaoFromForm(r)
    if err != nil {
        cartoesFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
            "Cartao": nil,
            "Titulo": "Novo Cartão",
            "Action": "/cartoes",
            "Erro":   "Dados inválidos: " + err.Error(),
        })
        return
    }

    if err := models.CriarCartao(cartao); err != nil {
        renderErro(w, "Erro ao salvar cartão", err.Error(), "/cartoes", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/cartoes", http.StatusSeeOther)
}

func CartoesEditar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/cartoes/", "/editar")
    if err != nil {
        renderErro(w, "ID inválido", "O identificador do cartão é inválido.", "/cartoes", http.StatusBadRequest)
        return
    }

    cartao, err := models.BuscarCartao(id)
    if err != nil {
        renderErro(w, "Cartão não encontrado", "O cartão solicitado não foi encontrado.", "/cartoes", http.StatusNotFound)
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
        renderErro(w, "ID inválido", "O identificador do cartão é inválido.", "/cartoes", http.StatusBadRequest)
        return
    }

    if err := r.ParseForm(); err != nil {
        renderErro(w, "Formulário inválido", err.Error(), "/cartoes", http.StatusBadRequest)
        return
    }

    cartao, err := cartaoFromForm(r)
    if err != nil {
        cartoesFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
            "Cartao": nil,
            "Titulo": "Editar Cartão",
            "Action": "/cartoes/" + strconv.Itoa(id),
            "Erro":   "Dados inválidos: " + err.Error(),
        })
        return
    }
    cartao.ID = id

    if err := models.AtualizarCartao(cartao); err != nil {
        renderErro(w, "Erro ao atualizar cartão", err.Error(), "/cartoes", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/cartoes", http.StatusSeeOther)
}

func CartoesDeletar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/cartoes/", "/deletar")
    if err != nil {
        renderErro(w, "ID inválido", "O identificador do cartão é inválido.", "/cartoes", http.StatusBadRequest)
        return
    }

    temDespesas, err := models.CartaoTemDespesas(id)
    if err != nil {
        renderErro(w, "Erro ao verificar dependências", err.Error(), "/cartoes", http.StatusInternalServerError)
        return
    }
    if temDespesas {
        renderErro(w, "Cartão em uso", "Não é possível excluir este cartão pois existem despesas associadas a ele.", "/cartoes", http.StatusConflict)
        return
    }

    if err := models.DeletarCartao(id); err != nil {
        renderErro(w, "Erro ao deletar cartão", err.Error(), "/cartoes", http.StatusInternalServerError)
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
