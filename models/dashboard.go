package models

import (
	"fincontrol/db"
	"fmt"
)

type TerceiroSaldo struct {
	ID             int
	Nome           string
	LimiteLiberado float64
	TotalCartao    float64
	TotalDinheiro  float64
	Total          float64
}

type EmprestimoAtivo struct {
	Descricao      string
	ValorParcela   float64
	ParcelaAtual   int
	TotalParcelas  int
	FormaPagamento string
	TerceiroNome   *string
}

type CartaoUso struct {
	ID            int
	Nome          string
	Limite        float64
	Utilizado     float64
	Disponivel    float64
	PctUso        float64
	Excedido      bool
}

func TerceirosSaldoPorMes(ano, mes int) ([]TerceiroSaldo, error) {
	mesStr := fmt.Sprintf("%04d-%02d", ano, mes)
	rows, err := db.DB.Query(`
		SELECT t.id, t.nome, t.limite_liberado,
		       SUM(CASE WHEN d.forma_pagamento = 'cartao' THEN d.valor_total ELSE 0 END),
		       SUM(CASE WHEN d.forma_pagamento != 'cartao' THEN d.valor_total ELSE 0 END)
		FROM terceiros t
		JOIN despesas d ON d.terceiro_id = t.id
		WHERE (d.fixa = 0 AND strftime('%Y-%m', d.data_compra) = ?)
		   OR (d.fixa = 1 AND strftime('%Y-%m', d.data_compra) <= ?)
		GROUP BY t.id, t.nome, t.limite_liberado
		ORDER BY t.nome ASC`, mesStr, mesStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []TerceiroSaldo
	for rows.Next() {
		var s TerceiroSaldo
		if err := rows.Scan(&s.ID, &s.Nome, &s.LimiteLiberado, &s.TotalCartao, &s.TotalDinheiro); err != nil {
			return nil, err
		}
		s.Total = s.TotalCartao + s.TotalDinheiro
		resultado = append(resultado, s)
	}
	return resultado, nil
}

func EmprestimosAtivosPorMes(ano, mes int) ([]EmprestimoAtivo, error) {
	mesStr := fmt.Sprintf("%04d-%02d", ano, mes)
	rows, err := db.DB.Query(`
		SELECT d.descricao, d.valor_total, d.parcela_atual, d.total_parcelas, d.forma_pagamento, t.nome
		FROM despesas d
		LEFT JOIN terceiros t ON t.id = d.terceiro_id
		WHERE d.categoria = 'Empréstimo'
		  AND strftime('%Y-%m', d.data_compra) = ?
		ORDER BY d.descricao ASC`, mesStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []EmprestimoAtivo
	for rows.Next() {
		var e EmprestimoAtivo
		var pa, tp *int
		if err := rows.Scan(&e.Descricao, &e.ValorParcela, &pa, &tp, &e.FormaPagamento, &e.TerceiroNome); err != nil {
			return nil, err
		}
		if pa != nil {
			e.ParcelaAtual = *pa
		}
		if tp != nil {
			e.TotalParcelas = *tp
		}
		resultado = append(resultado, e)
	}
	return resultado, nil
}

func CartoesPorUsoMes(ano, mes int) ([]CartaoUso, error) {
	mesStr := fmt.Sprintf("%04d-%02d", ano, mes)
	rows, err := db.DB.Query(`
		SELECT c.id, c.nome, c.limite,
		       COALESCE(SUM(d.valor_total), 0)
		FROM cartoes c
		LEFT JOIN despesas d ON d.cartao_id = c.id
		  AND d.forma_pagamento = 'cartao'
		  AND (
		        (d.fixa = 0 AND strftime('%Y-%m', d.data_compra) = ?)
		     OR (d.fixa = 1 AND strftime('%Y-%m', d.data_compra) <= ?)
		  )
		GROUP BY c.id, c.nome, c.limite
		ORDER BY c.nome ASC`, mesStr, mesStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []CartaoUso
	for rows.Next() {
		var c CartaoUso
		if err := rows.Scan(&c.ID, &c.Nome, &c.Limite, &c.Utilizado); err != nil {
			return nil, err
		}
		c.Disponivel = c.Limite - c.Utilizado
		c.Excedido = c.Utilizado > c.Limite
		if c.Limite > 0 {
			c.PctUso = c.Utilizado / c.Limite * 100
			if c.PctUso > 100 {
				c.PctUso = 100
			}
		}
		resultado = append(resultado, c)
	}
	return resultado, nil
}
