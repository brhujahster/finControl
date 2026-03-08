package models

import (
    "fincontrol/db"
    "time"
)

type Cartao struct {
    ID             int
    Nome           string
    DiaVencimento  int
    Limite         float64
    CriadoEm       time.Time
}

func ListarCartoes() ([]Cartao, error) {
    rows, err := db.DB.Query(`SELECT id, nome, dia_vencimento, limite, criado_em FROM cartoes ORDER BY nome ASC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cartoes []Cartao
    for rows.Next() {
        var c Cartao
        if err := rows.Scan(&c.ID, &c.Nome, &c.DiaVencimento, &c.Limite, &c.CriadoEm); err != nil {
            return nil, err
        }
        cartoes = append(cartoes, c)
    }
    return cartoes, nil
}

func BuscarCartao(id int) (*Cartao, error) {
    row := db.DB.QueryRow(`SELECT id, nome, dia_vencimento, limite, criado_em FROM cartoes WHERE id = ?`, id)
    var c Cartao
    if err := row.Scan(&c.ID, &c.Nome, &c.DiaVencimento, &c.Limite, &c.CriadoEm); err != nil {
        return nil, err
    }
    return &c, nil
}

func CriarCartao(c Cartao) error {
    _, err := db.DB.Exec(
        `INSERT INTO cartoes (nome, dia_vencimento, limite) VALUES (?, ?, ?)`,
        c.Nome, c.DiaVencimento, c.Limite,
    )
    return err
}

func AtualizarCartao(c Cartao) error {
    _, err := db.DB.Exec(
        `UPDATE cartoes SET nome=?, dia_vencimento=?, limite=? WHERE id=?`,
        c.Nome, c.DiaVencimento, c.Limite, c.ID,
    )
    return err
}

func DeletarCartao(id int) error {
    _, err := db.DB.Exec(`DELETE FROM cartoes WHERE id = ?`, id)
    return err
}