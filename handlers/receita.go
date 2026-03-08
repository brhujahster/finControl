package handlers

import (
    "fincontrol/models"
    "html/template"
    "net/http"
    "strconv"
    "strings"
    "time"
)

var receitasTmpl = template.Must(template.ParseFiles(
    "templates/base.html",
    "templates/receitas/lista.html",
))

var receitasFormTmpl = template.Must(template.ParseFiles(
    "templates/base.html",
    "templates/receitas/form.html",
))

func ReceitasIndex(w http.ResponseWriter, r *http.Request) {
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

    receitas, err := models.ListarReceitasPorMes(ano, mes)
    if err != nil {
        http.Error(w, "Erro ao listar receitas: "+err.Error(), http.StatusInternalServerError)
        return
    }

    total := 0.0
    for _, rc := range receitas {
        total += rc.Valor
    }

    data := map[string]interface{}{
        "Receitas": receitas,
        "Total":    total,
        "Ano":      ano,
        "Mes":      mes,
    }
    receitasTmpl.ExecuteTemplate(w, "base", data)
}

func ReceitasNova(w http.ResponseWriter, r *http.Request) {
    receitasFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
        "Receita": nil,
        "Titulo":  "Nova Receita",
        "Action":  "/receitas",
    })
}

func ReceitasSalvar(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Formulário inválido", http.StatusBadRequest)
        return
    }

    receita, err := receitaFromForm(r)
    if err != nil {
        http.Error(w, "Dados inválidos: "+err.Error(), http.StatusBadRequest)
        return
    }

    if err := models.CriarReceita(receita); err != nil {
        http.Error(w, "Erro ao salvar receita: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/receitas", http.StatusSeeOther)
}

func ReceitasEditar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/receitas/", "/editar")
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    receita, err := models.BuscarReceita(id)
    if err != nil {
        http.Error(w, "Receita não encontrada", http.StatusNotFound)
        return
    }

    receitasFormTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
        "Receita": receita,
        "Titulo":  "Editar Receita",
        "Action":  "/receitas/" + strconv.Itoa(id),
    })
}

func ReceitasAtualizar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/receitas/", "")
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    if err := r.ParseForm(); err != nil {
        http.Error(w, "Formulário inválido", http.StatusBadRequest)
        return
    }

    receita, err := receitaFromForm(r)
    if err != nil {
        http.Error(w, "Dados inválidos: "+err.Error(), http.StatusBadRequest)
        return
    }
    receita.ID = id

    if err := models.AtualizarReceita(receita); err != nil {
        http.Error(w, "Erro ao atualizar receita: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/receitas", http.StatusSeeOther)
}

func ReceitasDeletar(w http.ResponseWriter, r *http.Request) {
    id, err := extrairID(r.URL.Path, "/receitas/", "/deletar")
    if err != nil {
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }

    if err := models.DeletarReceita(id); err != nil {
        http.Error(w, "Erro ao deletar receita: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/receitas", http.StatusSeeOther)
}

func receitaFromForm(r *http.Request) (models.Receita, error) {
    valor, err := strconv.ParseFloat(strings.ReplaceAll(r.FormValue("valor"), ",", "."), 64)
    if err != nil {
        return models.Receita{}, err
    }

    receita := models.Receita{
        Descricao: r.FormValue("descricao"),
        Valor:     valor,
        Tipo:      r.FormValue("tipo"),
    }

    if v := r.FormValue("mes_referencia"); v != "" {
        t, err := time.Parse("2006-01-02", v)
        if err != nil {
            return models.Receita{}, err
        }
        receita.MesReferencia = &t
    }
    if v := r.FormValue("data_inicio"); v != "" {
        t, err := time.Parse("2006-01-02", v)
        if err != nil {
            return models.Receita{}, err
        }
        receita.DataInicio = &t
    }
    if v := r.FormValue("data_fim"); v != "" {
        t, err := time.Parse("2006-01-02", v)
        if err != nil {
            return models.Receita{}, err
        }
        receita.DataFim = &t
    }

    return receita, nil
}

func extrairID(path, prefix, suffix string) (int, error) {
    s := strings.TrimPrefix(path, prefix)
    s = strings.TrimSuffix(s, suffix)
    return strconv.Atoi(s)
}