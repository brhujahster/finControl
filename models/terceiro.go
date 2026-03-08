package models

import (
    "fincontrol/db"
    "time"
)

type Terceiro struct {
    ID             int
    Nome           string
    LimiteLiberado float64
    Observacao     *string
    CriadoEm      time.Time
}

func ListarTerceiros() ([]Terceiro, error) {
    rows, err := db.DB.Query(`SELECT id, nome, limite_liberado, observacao, criado_em FROM terceiros ORDER BY nome ASC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var terceiros []Terceiro
    for rows.Next() {
        var t Terceiro
        if err := rows.Scan(&t.ID, &t.Nome, &t.LimiteLiberado, &t.Observacao, &t.CriadoEm); err != nil {
            return nil, err
        }
        terceiros = append(terceiros, t)
    }
    return terceiros, nil
}

func BuscarTerceiro(id int) (*Terceiro, error) {
    row := db.DB.QueryRow(`SELECT id, nome, limite_liberado, observacao, criado_em FROM terceiros WHERE id = ?`, id)
    var t Terceiro
    if err := row.Scan(&t.ID, &t.Nome, &t.LimiteLiberado, &t.Observacao, &t.CriadoEm); err != nil {
        return nil, err
    }
    return &t, nil
}

func CriarTerceiro(t Terceiro) error {
    _, err := db.DB.Exec(
        `INSERT INTO terceiros (nome, limite_liberado, observacao) VALUES (?, ?, ?)`,
        t.Nome, t.LimiteLiberado, t.Observacao,
    )
    return err
}

func AtualizarTerceiro(t Terceiro) error {
    _, err := db.DB.Exec(
        `UPDATE terceiros SET nome=?, limite_liberado=?, observacao=? WHERE id=?`,
        t.Nome, t.LimiteLiberado, t.Observacao, t.ID,
    )
    return err
}

func DeletarTerceiro(id int) error {
    _, err := db.DB.Exec(`DELETE FROM terceiros WHERE id = ?`, id)
    return err
}