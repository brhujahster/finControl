package models

import (
    "fincontrol/db"
    "time"
)

type Receita struct {
    ID            int
    Descricao     string
    Valor         float64
    Tipo          string
    MesReferencia *time.Time
    DataInicio    *time.Time
    DataFim       *time.Time
    CriadoEm      time.Time
}

func ListarReceitas() ([]Receita, error) {
    rows, err := db.DB.Query(`SELECT id, descricao, valor, tipo, mes_referencia, data_inicio, data_fim, criado_em FROM receitas ORDER BY criado_em DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var receitas []Receita
    for rows.Next() {
        var r Receita
        var mesRef, dataInicio, dataFim *string
        if err := rows.Scan(&r.ID, &r.Descricao, &r.Valor, &r.Tipo, &mesRef, &dataInicio, &dataFim, &r.CriadoEm); err != nil {
            return nil, err
        }
        r.MesReferencia = parseDate(mesRef)
        r.DataInicio = parseDate(dataInicio)
        r.DataFim = parseDate(dataFim)
        receitas = append(receitas, r)
    }
    return receitas, nil
}

func ListarReceitasPorMes(ano, mes int) ([]Receita, error) {
    alvo := time.Date(ano, time.Month(mes), 1, 0, 0, 0, 0, time.UTC)
    todas, err := ListarReceitas()
    if err != nil {
        return nil, err
    }

    var resultado []Receita
    for _, r := range todas {
        if r.Tipo == "pontual" {
            if r.MesReferencia != nil &&
                r.MesReferencia.Year() == ano &&
                int(r.MesReferencia.Month()) == mes {
                resultado = append(resultado, r)
            }
        } else {
            if r.DataInicio == nil {
                continue
            }
            inicio := time.Date(r.DataInicio.Year(), r.DataInicio.Month(), 1, 0, 0, 0, 0, time.UTC)
            if alvo.Before(inicio) {
                continue
            }
            if r.DataFim != nil {
                fim := time.Date(r.DataFim.Year(), r.DataFim.Month(), 1, 0, 0, 0, 0, time.UTC)
                if alvo.After(fim) {
                    continue
                }
            }
            resultado = append(resultado, r)
        }
    }
    return resultado, nil
}

func BuscarReceita(id int) (*Receita, error) {
    row := db.DB.QueryRow(`SELECT id, descricao, valor, tipo, mes_referencia, data_inicio, data_fim, criado_em FROM receitas WHERE id = ?`, id)
    var r Receita
    var mesRef, dataInicio, dataFim *string
    if err := row.Scan(&r.ID, &r.Descricao, &r.Valor, &r.Tipo, &mesRef, &dataInicio, &dataFim, &r.CriadoEm); err != nil {
        return nil, err
    }
    r.MesReferencia = parseDate(mesRef)
    r.DataInicio = parseDate(dataInicio)
    r.DataFim = parseDate(dataFim)
    return &r, nil
}

func CriarReceita(r Receita) error {
    _, err := db.DB.Exec(
        `INSERT INTO receitas (descricao, valor, tipo, mes_referencia, data_inicio, data_fim) VALUES (?, ?, ?, ?, ?, ?)`,
        r.Descricao, r.Valor, r.Tipo, formatDate(r.MesReferencia), formatDate(r.DataInicio), formatDate(r.DataFim),
    )
    return err
}

func AtualizarReceita(r Receita) error {
    _, err := db.DB.Exec(
        `UPDATE receitas SET descricao=?, valor=?, tipo=?, mes_referencia=?, data_inicio=?, data_fim=? WHERE id=?`,
        r.Descricao, r.Valor, r.Tipo, formatDate(r.MesReferencia), formatDate(r.DataInicio), formatDate(r.DataFim), r.ID,
    )
    return err
}

func DeletarReceita(id int) error {
    _, err := db.DB.Exec(`DELETE FROM receitas WHERE id = ?`, id)
    return err
}

func parseDate(s *string) *time.Time {
    if s == nil || *s == "" {
        return nil
    }
    for _, layout := range []string{"2006-01-02", "2006-01-02 15:04:05"} {
        if t, err := time.Parse(layout, *s); err == nil {
            return &t
        }
    }
    return nil
}

func formatDate(t *time.Time) *string {
    if t == nil {
        return nil
    }
    s := t.Format("2006-01-02")
    return &s
}