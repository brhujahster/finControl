package models

import (
    "fincontrol/db"
    "math"
    "time"
)

type Emprestimo struct {
    ID             int
    Descricao      string
    ValorTotal     float64
    TotalParcelas  int
    TaxaJuros      *float64
    TipoJuros      *string
    DataInicio     time.Time
    FormaPagamento string
    CartaoID       *int
    CartaoNome     *string
    TerceiroID     *int
    TerceiroNome   *string
    Observacao     *string
    CriadoEm      time.Time
}

func ListarEmprestimos() ([]Emprestimo, error) {
    rows, err := db.DB.Query(`
        SELECT e.id, e.descricao, e.valor_total, e.total_parcelas, e.taxa_juros, e.tipo_juros,
               e.data_inicio, e.forma_pagamento, e.cartao_id, c.nome, e.terceiro_id, t.nome,
               e.observacao, e.criado_em
        FROM emprestimos e
        LEFT JOIN cartoes c ON c.id = e.cartao_id
        LEFT JOIN terceiros t ON t.id = e.terceiro_id
        ORDER BY e.data_inicio DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    return scanEmprestimos(rows)
}

func BuscarEmprestimo(id int) (*Emprestimo, error) {
    row := db.DB.QueryRow(`
        SELECT e.id, e.descricao, e.valor_total, e.total_parcelas, e.taxa_juros, e.tipo_juros,
               e.data_inicio, e.forma_pagamento, e.cartao_id, c.nome, e.terceiro_id, t.nome,
               e.observacao, e.criado_em
        FROM emprestimos e
        LEFT JOIN cartoes c ON c.id = e.cartao_id
        LEFT JOIN terceiros t ON t.id = e.terceiro_id
        WHERE e.id = ?`, id)

    var e Emprestimo
    var dataInicio, criadoEm string
    err := row.Scan(
        &e.ID, &e.Descricao, &e.ValorTotal, &e.TotalParcelas, &e.TaxaJuros, &e.TipoJuros,
        &dataInicio, &e.FormaPagamento, &e.CartaoID, &e.CartaoNome, &e.TerceiroID, &e.TerceiroNome,
        &e.Observacao, &criadoEm,
    )
    if err != nil {
        return nil, err
    }
    if t, err := time.Parse("2006-01-02", dataInicio); err == nil {
        e.DataInicio = t
    }
    return &e, nil
}

func CriarEmprestimo(e Emprestimo, diaVencimento int) error {
    res, err := db.DB.Exec(`
        INSERT INTO emprestimos (descricao, valor_total, total_parcelas, taxa_juros, tipo_juros,
                                 data_inicio, forma_pagamento, cartao_id, terceiro_id, observacao)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        e.Descricao, e.ValorTotal, e.TotalParcelas, e.TaxaJuros, e.TipoJuros,
        e.DataInicio.Format("2006-01-02"), e.FormaPagamento, e.CartaoID, e.TerceiroID, e.Observacao,
    )
    if err != nil {
        return err
    }

    id, _ := res.LastInsertId()
    return gerarDespesasEmprestimo(e, int(id), diaVencimento)
}

func gerarDespesasEmprestimo(e Emprestimo, emprestimoID, diaVencimento int) error {
    parcelas := calcularParcelas(e)

    var competenciaBase time.Time
    if e.FormaPagamento == "cartao" {
        competenciaBase = e.DataInicio
        if e.DataInicio.Day() >= diaVencimento {
            competenciaBase = competenciaBase.AddDate(0, 1, 0)
        }
        competenciaBase = time.Date(competenciaBase.Year(), competenciaBase.Month(), 1, 0, 0, 0, 0, time.UTC)
    } else {
        competenciaBase = time.Date(e.DataInicio.Year(), e.DataInicio.Month(), e.DataInicio.Day(), 0, 0, 0, 0, time.UTC)
    }

    descBase := e.Descricao
    total := e.TotalParcelas

    for i, valorParcela := range parcelas {
        dataComp := competenciaBase.AddDate(0, i, 0)
        parcela := i + 1
        _, err := db.DB.Exec(`
            INSERT INTO despesas (descricao, valor_total, data_compra, forma_pagamento, cartao_id,
                                  parcelado, total_parcelas, parcela_atual, fixa, terceiro_id, categoria, observacao)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
            descBase, valorParcela, dataComp.Format("2006-01-02"), e.FormaPagamento,
            e.CartaoID, total > 1, total, parcela, false,
            e.TerceiroID, "Empréstimo",
            emprestimoID,
        )
        if err != nil {
            return err
        }
    }
    return nil
}

func calcularParcelas(e Emprestimo) []float64 {
    n := e.TotalParcelas
    parcelas := make([]float64, n)

    if e.TaxaJuros == nil || *e.TaxaJuros == 0 {
        vp := math.Round(e.ValorTotal/float64(n)*100) / 100
        for i := range parcelas {
            parcelas[i] = vp
        }
        return parcelas
    }

    taxa := *e.TaxaJuros / 100.0

    if e.TipoJuros != nil && *e.TipoJuros == "composto" {
        pmt := e.ValorTotal * taxa / (1 - math.Pow(1+taxa, -float64(n)))
        pmt = math.Round(pmt*100) / 100
        for i := range parcelas {
            parcelas[i] = pmt
        }
    } else {
        totalComJuros := e.ValorTotal * (1 + taxa*float64(n))
        vp := math.Round(totalComJuros/float64(n)*100) / 100
        for i := range parcelas {
            parcelas[i] = vp
        }
    }

    return parcelas
}

func AtualizarEmprestimo(e Emprestimo) error {
    _, err := db.DB.Exec(`
        UPDATE emprestimos SET descricao=?, valor_total=?, total_parcelas=?, taxa_juros=?, tipo_juros=?,
        data_inicio=?, forma_pagamento=?, cartao_id=?, terceiro_id=?, observacao=? WHERE id=?`,
        e.Descricao, e.ValorTotal, e.TotalParcelas, e.TaxaJuros, e.TipoJuros,
        e.DataInicio.Format("2006-01-02"), e.FormaPagamento, e.CartaoID, e.TerceiroID, e.Observacao, e.ID,
    )
    return err
}

func DeletarEmprestimo(id int) error {
    if _, err := db.DB.Exec(`DELETE FROM despesas WHERE categoria = 'Empréstimo' AND CAST(observacao AS INTEGER) = ?`, id); err != nil {
        return err
    }
    _, err := db.DB.Exec(`DELETE FROM emprestimos WHERE id = ?`, id)
    return err
}

func scanEmprestimos(rows interface {
    Scan(...interface{}) error
    Next() bool
    Close() error
}) ([]Emprestimo, error) {
    defer rows.Close()
    var lista []Emprestimo
    for rows.Next() {
        var e Emprestimo
        var dataInicio, criadoEm string
        if err := rows.Scan(
            &e.ID, &e.Descricao, &e.ValorTotal, &e.TotalParcelas, &e.TaxaJuros, &e.TipoJuros,
            &dataInicio, &e.FormaPagamento, &e.CartaoID, &e.CartaoNome, &e.TerceiroID, &e.TerceiroNome,
            &e.Observacao, &criadoEm,
        ); err != nil {
            return nil, err
        }
        if t, err := time.Parse("2006-01-02", dataInicio); err == nil {
            e.DataInicio = t
        }
        lista = append(lista, e)
    }
    return lista, nil
}
