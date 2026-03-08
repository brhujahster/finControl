package main

import (
    "fincontrol/db"
    "fincontrol/handlers"
    "log"
    "net/http"
	"strings"
)

func main() {
    db.Connect()
    db.Migrate()

    mux := http.NewServeMux()

    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    mux.HandleFunc("/receitas", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            handlers.ReceitasSalvar(w, r)
        } else {
            handlers.ReceitasIndex(w, r)
        }
    })
    mux.HandleFunc("/receitas/nova", handlers.ReceitasNova)
    mux.HandleFunc("/receitas/", func(w http.ResponseWriter, r *http.Request) {
        switch {
        case strings.HasSuffix(r.URL.Path, "/editar"):
            handlers.ReceitasEditar(w, r)
        case strings.HasSuffix(r.URL.Path, "/deletar") && r.Method == http.MethodPost:
            handlers.ReceitasDeletar(w, r)
        case r.Method == http.MethodPost:
            handlers.ReceitasAtualizar(w, r)
        default:
            http.NotFound(w, r)
        }
    })

    addr := "localhost:8080"
    log.Printf("Servidor iniciado em http://%s\n", addr)
    log.Fatal(http.ListenAndServe(addr, mux))
}