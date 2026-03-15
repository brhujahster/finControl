package handlers

import (
    "html/template"
    "net/http"
)

var erroTmpl = template.Must(template.ParseFiles(
    "templates/base.html",
    "templates/error.html",
))

func renderErro(w http.ResponseWriter, titulo, mensagem, voltar string, status int) {
    w.WriteHeader(status)
    erroTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
        "Titulo":   titulo,
        "Mensagem": mensagem,
        "Voltar":   voltar,
    })
}
