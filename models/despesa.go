package models

import (
    "fincontrol/db"
    "fmt"
    "time"
)

type Despesa struct {
    ID             int
    Descricao      string
    ValorTotal     float64
    DataCompra     time.Time
    FormaPagamento string
    CartaoID       *int
    CartaoNome     *string
    Parcelado      bool
    TotalParcelas  *int
    ParcelaAtual   *int
    Fixa           bool
    TerceiroID     *int
    TerceiroNome   *string
    Categoria      *string
    Observacao     *string
    CriadoEm      time.Time
}

func ListarDespesasPorMes(ano, mes int) ([]Despesa, error) {
    mesStr := fmt.Sprintf("%04d-%02d", ano, mes)
    rows, err := db.DB.Query(`
        SELECT d.id, d.descricao, d.valor_total, d.data_compra, d.forma_pagamento,
               d.cartao_id, c.nome, d.parcelado, d.total_parcelas, d.parcela_atual,
               d.fixa, d.terceiro_id, t.nome, d.categoria, d.observacao, d.criado_em
        FROM despesas d
        LEFT JOIN cartoes c ON c.id = d.cartao_id
        LEFT JOIN terceiros t ON t.id = d.terceiro_id
        WHERE (d.fixa = 0 AND strftime('%Y-%m', d.data_compra) = ?)
           OR (d.fixa = 1 AND strftime('%Y-%m', d.data_compra) <= ?)
        ORDER BY d.data_compra ASC`, mesStr, mesStr)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    return scanDespesas(rows)
}

func ListarTodasDespesas() ([]Despesa, error) {
    rows, err := db.DB.Query(`
        SELECT d.id, d.descricao, d.valor_total, d.data_compra, d.forma_pagamento,
               d.cartao_id, c.nome, d.parcelado, d.total_parcelas, d.parcela_atual,
               d.fixa, d.terceiro_id, t.nome, d.categoria, d.observacao, d.criado_em
        FROM despesas d
        LEFT JOIN cartoes c ON c.id = d.cartao_id
        LEFT JOIN terceiros t ON t.id = d.terceiro_id
        ORDER BY d.data_compra DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    return scanDespesas(rows)
}

func BuscarDespesa(id int) (*Despesa, error) {
    row := db.DB.QueryRow(`
        SELECT d.id, d.descricao, d.valor_total, d.data_compra, d.forma_pagamento,
               d.cartao_id, c.nome, d.parcelado, d.total_parcelas, d.parcela_atual,
               d.fixa, d.terceiro_id, t.nome, d.categoria, d.observacao, d.criado_em
        FROM despesas d
        LEFT JOIN cartoes c ON c.id = d.cartao_id
        LEFT JOIN terceiros t ON t.id = d.terceiro_id
        WHERE d.id = ?`, id)

    despesas, err := scanDespesas(rowToRows(row))
    if err != nil || len(despesas) == 0 {
        return nil, err
    }
    return &despesas[0], nil
}

func CriarDespesa(d Despesa, cartaoDiaVencimento int) error {
    if d.Parcelado && d.TotalParcelas != nil && *d.TotalParcelas > 1 {
        return inserirParcelas(d, cartaoDiaVencimento)
    }
    _, err := db.DB.Exec(`
        INSERT INTO despesas (descricao, valor_total, data_compra, forma_pagamento, cartao_id,
                              parcelado, total_parcelas, parcela_atual, fixa, terceiro_id, categoria, observacao)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        d.Descricao, d.ValorTotal, d.DataCompra.Format("2006-01-02"), d.FormaPagamento,
        d.CartaoID, d.Parcelado, d.TotalParcelas, 1, d.Fixa,
        d.TerceiroID, d.Categoria, d.Observacao,
    )
    return err
}

func inserirParcelas(d Despesa, diaVencimento int) error {
    total := *d.TotalParcelas
    valorParcela := d.ValorTotal / float64(total)

    competencia := d.DataCompra
    if d.DataCompra.Day() >= diaVencimento {
        competencia = competencia.AddDate(0, 1, 0)
    }
    competencia = time.Date(competencia.Year(), competencia.Month(), 1, 0, 0, 0, 0, time.UTC)

    for i := 1; i <= total; i++ {
        mesComp := competencia.AddDate(0, i-1, 0)
        parcela := i
        _, err := db.DB.Exec(`
            INSERT INTO despesas (descricao, valor_total, data_compra, forma_pagamento, cartao_id,
                                  parcelado, total_parcelas, parcela_atual, fixa, terceiro_id, categoria, observacao)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
            d.Descricao, valorParcela, mesComp.Format("2006-01-02"), d.FormaPagamento,
            d.CartaoID, true, total, parcela, false,
            d.TerceiroID, d.Categoria, d.Observacao,
        )
        if err != nil {
            return err
        }
    }
    return nil
}

func AtualizarDespesa(d Despesa) error {
    _, err := db.DB.Exec(`
        UPDATE despesas SET descricao=?, valor_total=?, data_compra=?, forma_pagamento=?,
        cartao_id=?, fixa=?, terceiro_id=?, categoria=?, observacao=? WHERE id=?`,
        d.Descricao, d.ValorTotal, d.DataCompra.Format("2006-01-02"), d.FormaPagamento,
        d.CartaoID, d.Fixa, d.TerceiroID, d.Categoria, d.Observacao, d.ID,
    )
    return err
}

func DeletarDespesa(id int) error {
    _, err := db.DB.Exec(`DELETE FROM despesas WHERE id = ?`, id)
    return err
}

func scanDespesas(rows interface{ Scan(...interface{}) error; Next() bool; Close() error }) ([]Despesa, error) {
    defer rows.Close()
    var despesas []Despesa
    for rows.Next() {
        var d Despesa
        var dataCompra, criadoEm string
        if err := rows.Scan(
            &d.ID, &d.Descricao, &d.ValorTotal, &dataCompra, &d.FormaPagamento,
            &d.CartaoID, &d.CartaoNome, &d.Parcelado, &d.TotalParcelas, &d.ParcelaAtual,
            &d.Fixa, &d.TerceiroID, &d.TerceiroNome, &d.Categoria, &d.Observacao, &criadoEm,
        ); err != nil {
            return nil, err
        }
        if t, err := time.Parse("2006-01-02", dataCompra); err == nil {
            d.DataCompra = t
        }
        despesas = append(despesas, d)
    }
    return despesas, nil
}

type singleRow struct {
    row interface{ Scan(...interface{}) error }
    done bool
}

func rowToRows(row interface{ Scan(...interface{}) error }) *singleRow {
    return &singleRow{row: row}
}

func (r *singleRow) Next() bool {
    if r.done {
        return false
    }
    r.done = true
    return true
}
func (r *singleRow) Scan(dest ...interface{}) error { return r.row.Scan(dest...) }
func (r *singleRow) Close() error                   { return nil }