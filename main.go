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

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        http.Redirect(w, r, "/receitas", http.StatusSeeOther)
    })

    // Receitas
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

    // Cartões
    mux.HandleFunc("/cartoes", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            handlers.CartoesSalvar(w, r)
        } else {
            handlers.CartoesIndex(w, r)
        }
    })
    mux.HandleFunc("/cartoes/novo", handlers.CartoesNovo)
    mux.HandleFunc("/cartoes/", func(w http.ResponseWriter, r *http.Request) {
        switch {
        case strings.HasSuffix(r.URL.Path, "/editar"):
            handlers.CartoesEditar(w, r)
        case strings.HasSuffix(r.URL.Path, "/deletar") && r.Method == http.MethodPost:
            handlers.CartoesDeletar(w, r)
        case r.Method == http.MethodPost:
            handlers.CartoesAtualizar(w, r)
        default:
            http.NotFound(w, r)
        }
    })

    // Terceiros
    mux.HandleFunc("/terceiros", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            handlers.TerceirosSalvar(w, r)
        } else {
            handlers.TerceirosIndex(w, r)
        }
    })
    mux.HandleFunc("/terceiros/novo", handlers.TerceirosNovo)
    mux.HandleFunc("/terceiros/", func(w http.ResponseWriter, r *http.Request) {
        switch {
        case strings.HasSuffix(r.URL.Path, "/editar"):
            handlers.TerceirosEditar(w, r)
        case strings.HasSuffix(r.URL.Path, "/deletar") && r.Method == http.MethodPost:
            handlers.TerceirosDeletar(w, r)
        case r.Method == http.MethodPost:
            handlers.TerceirosAtualizar(w, r)
        default:
            http.NotFound(w, r)
        }
    })

    // Despesas
    mux.HandleFunc("/despesas", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            handlers.DespesasSalvar(w, r)
        } else {
            handlers.DespesasIndex(w, r)
        }
    })
    mux.HandleFunc("/despesas/nova", handlers.DespesasNova)
    mux.HandleFunc("/despesas/", func(w http.ResponseWriter, r *http.Request) {
        switch {
        case strings.HasSuffix(r.URL.Path, "/editar"):
            handlers.DespesasEditar(w, r)
        case strings.HasSuffix(r.URL.Path, "/deletar") && r.Method == http.MethodPost:
            handlers.DespesasDeletar(w, r)
        case r.Method == http.MethodPost:
            handlers.DespesasAtualizar(w, r)
        default:
            http.NotFound(w, r)
        }
    })

    // Empréstimos
    mux.HandleFunc("/emprestimos", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            handlers.EmprestimosSalvar(w, r)
        } else {
            handlers.EmprestimosIndex(w, r)
        }
    })
    mux.HandleFunc("/emprestimos/novo", handlers.EmprestimosNovo)
    mux.HandleFunc("/emprestimos/", func(w http.ResponseWriter, r *http.Request) {
        switch {
        case strings.HasSuffix(r.URL.Path, "/editar"):
            handlers.EmprestimosEditar(w, r)
        case strings.HasSuffix(r.URL.Path, "/deletar") && r.Method == http.MethodPost:
            handlers.EmprestimosDeletar(w, r)
        case r.Method == http.MethodPost:
            handlers.EmprestimosAtualizar(w, r)
        default:
            http.NotFound(w, r)
        }
    })

    addr := "localhost:8080"
    log.Printf("Servidor iniciado em http://%s\n", addr)
    log.Fatal(http.ListenAndServe(addr, mux))
}