package handlers

import (
    "fincontrol/models"
    "html/template"
    "net/http"
    "strconv"
    "strings"
    "time"
)

var funcMap = template.FuncMap{
    "deref": func(p *int) int {
        if p == nil {
            return 0
        }
        return *p
    },
    "derefStr": func(p *string) string {
        if p == nil {
            return ""
        }
        return *p
    },
}

var despesasTmpl = template.Must(template.New("base.html").Funcs(funcMap).ParseFiles(
    "templates/base.html",
    "templates/despesas/lista.html",
))

var despesasFormTmpl = template.Must(template.New("base.html").Funcs(funcMap).ParseFiles(
    "templates/base.html",
    "templates/despesas/form.html",
))

func DespesasIndex(w http.ResponseWriter, r *http.Request) {
    agora := time.Now()
    ano := agora.Year()
    mes := int(agora.Month())

    if v := r.URL.Query().Get("ano"); v != "" {
        if n, err := strconv.Atoi(v); err == nil {
            ano = n
        }
    }
    if v := r.URL.Query().Get("mes"); v != "" {
        if n, err := strconv.Atoi(v); err == nil {
            mes = n
        }
    }

    despesas, err := models.ListarDespesasPorMes(ano, mes)
    if err != nil {
        http.Error(w, "Erro ao listar despesas: "+err.Error(), http.StatusInternalServerError)
        return
    }

    totalCartao := 0.0
    totalDinheiro := 0.0
    for _, d := range despesas {
        if d.FormaPagamento == "cartao" {
            totalCartao += d.ValorTotal
        } else {
            totalDinheiro += d.ValorTotal
        }
    }

    despesasTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
        "Despesas":      despesas,
        "TotalCartao":   totalCartao,
        "TotalDinheiro": totalDinheiro,
        "Total":         totalCartao + totalDinheiro,
        "Ano":           ano,
        "Mes":           mes,
    })
}

func DespesasNova(w http.ResponseWriter, r *http.Request) {
    cartoes, _ := models.ListarCartoes()
    terceiros, _ := models.ListarTerceiros()
    despesasFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
        "Despesa":   nil,
        "Titulo":    "Nova Despesa",
        "Action":    "/despesas",
        "Cartoes":   cartoes,
        "Terceiros": terceiros,
    })
}

func DespesasSalvar(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Formulário inválido", http.StatusBadRequest)
        return
    }

    despesa, cartaoID, err := despesaFromForm(r)
    if err != nil {
        http.Error(w, "Dados inválidos: "+err.Error(), http.StatusBadRequest)
        return
    }

    diaVencimento := 0
    if cartaoID != nil {
        cartao, err := models.BuscarCartao(*cartaoID)
        if err == nil {
            diaVencimento = cartao.DiaVencimento
        }
    }

    if err := models.CriarDespesa(despesa, diaVencimento); err != nil {
        http.Error(w, "Erro ao salvar despesa: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/despesas", http.StatusSeeOther)
}

func DespesasEditar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/despesas/", "/editar")
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    despesa, err := models.BuscarDespesa(id)
    if err != nil {
        http.Error(w, "Despesa não encontrada", http.StatusNotFound)
        return
    }

    cartoes, _ := models.ListarCartoes()
    terceiros, _ := models.ListarTerceiros()
    despesasFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
        "Despesa":   despesa,
        "Titulo":    "Editar Despesa",
        "Action":    "/despesas/" + strconv.Itoa(id),
        "Cartoes":   cartoes,
        "Terceiros": terceiros,
    })
}

func DespesasAtualizar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/despesas/", "")
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Formulário inválido", http.StatusBadRequest)
        return
    }

    despesa, _, err := despesaFromForm(r)
    if err != nil {
        http.Error(w, "Dados inválidos: "+err.Error(), http.StatusBadRequest)
        return
    }
    despesa.ID = id

    if err := models.AtualizarDespesa(despesa); err != nil {
        http.Error(w, "Erro ao atualizar despesa: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/despesas", http.StatusSeeOther)
}

func DespesasDeletar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/despesas/", "/deletar")
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    if err := models.DeletarDespesa(id); err != nil {
        http.Error(w, "Erro ao deletar despesa: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/despesas", http.StatusSeeOther)
}

func despesaFromForm(r *http.Request) (models.Despesa, *int, error) {
    valorTotal, err := strconv.ParseFloat(strings.ReplaceAll(r.FormValue("valor_total"), ",", "."), 64)
    if err != nil {
        return models.Despesa{}, nil, err
    }

    dataCompra, err := time.Parse("2006-01-02", r.FormValue("data_compra"))
    if err != nil {
        return models.Despesa{}, nil, err
    }

    d := models.Despesa{
        Descricao:      r.FormValue("descricao"),
        ValorTotal:     valorTotal,
        DataCompra:     dataCompra,
        FormaPagamento: r.FormValue("forma_pagamento"),
        Parcelado:      r.FormValue("parcelado") == "1",
        Fixa:           r.FormValue("fixa") == "1",
    }

    var cartaoID *int
    if v := r.FormValue("cartao_id"); v != "" {
        n, err := strconv.Atoi(v)
        if err != nil {
            return models.Despesa{}, nil, err
        }
        d.CartaoID = &n
        cartaoID = &n
    }

    if v := r.FormValue("total_parcelas"); v != "" && d.Parcelado {
        n, err := strconv.Atoi(v)
        if err != nil {
            return models.Despesa{}, nil, err
        }
        d.TotalParcelas = &n
    }

    if v := r.FormValue("terceiro_id"); v != "" {
        n, err := strconv.Atoi(v)
        if err != nil {
            return models.Despesa{}, nil, err
        }
        d.TerceiroID = &n
    }

    if v := r.FormValue("categoria"); v != "" {
        d.Categoria = &v
    }
    if v := r.FormValue("observacao"); v != "" {
        d.Observacao = &v
    }

    return d, cartaoID, nil
}