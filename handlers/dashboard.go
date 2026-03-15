package handlers

import (
	"fincontrol/models"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var dashboardTmpl = template.Must(template.New("base.html").Funcs(funcMap).ParseFiles(
	"templates/base.html",
	"templates/dashboard.html",
))

func Dashboard(w http.ResponseWriter, r *http.Request) {
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

	receitas, _ := models.ListarReceitasPorMes(ano, mes)
	despesas, _ := models.ListarDespesasPorMes(ano, mes)
	terceiros, _ := models.TerceirosSaldoPorMes(ano, mes)
	emprestimos, _ := models.EmprestimosAtivosPorMes(ano, mes)
	cartoesUso, _ := models.CartoesPorUsoMes(ano, mes)

	receitaTotal := 0.0
	for _, rc := range receitas {
		receitaTotal += rc.Valor
	}

	totalCartao := 0.0
	totalDinheiro := 0.0
	for _, d := range despesas {
		if d.FormaPagamento == "cartao" {
			totalCartao += d.ValorTotal
		} else {
			totalDinheiro += d.ValorTotal
		}
	}
	totalDespesas := totalCartao + totalDinheiro
	saldo := receitaTotal - totalDespesas

	percentual := 0.0
	if receitaTotal > 0 {
		percentual = totalDespesas / receitaTotal * 100
	}

	cartaoExcedido := false
	for _, c := range cartoesUso {
		if c.Excedido {
			cartaoExcedido = true
			break
		}
	}

	terceiroExcedido := false
	for _, t := range terceiros {
		if t.Total > t.LimiteLiberado {
			terceiroExcedido = true
			break
		}
	}

	mesAnterior := mes - 1
	anoAnterior := ano
	if mesAnterior < 1 {
		mesAnterior = 12
		anoAnterior--
	}
	mesProximo := mes + 1
	anoProximo := ano
	if mesProximo > 12 {
		mesProximo = 1
		anoProximo++
	}

	mesesNomes := []string{"", "Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho",
		"Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro"}

	pctCartao := 0.0
	pctDinheiro := 0.0
	if totalDespesas > 0 {
		pctCartao = totalCartao / totalDespesas * 100
		pctDinheiro = totalDinheiro / totalDespesas * 100
	}

	dashboardTmpl.ExecuteTemplate(w, "base", map[string]interface{}{
		"Ano":              ano,
		"Mes":              mes,
		"MesNome":          mesesNomes[mes],
		"ReceitaTotal":     receitaTotal,
		"TotalDespesas":    totalDespesas,
		"TotalCartao":      totalCartao,
		"TotalDinheiro":    totalDinheiro,
		"PctCartao":        pctCartao,
		"PctDinheiro":      pctDinheiro,
		"Saldo":            saldo,
		"Percentual":       percentual,
		"Terceiros":        terceiros,
		"Emprestimos":      emprestimos,
		"CartoesUso":       cartoesUso,
		"SaldoNegativo":    saldo < 0,
		"CartaoExcedido":   cartaoExcedido,
		"TerceiroExcedido": terceiroExcedido,
		"MesAnterior":      mesAnterior,
		"AnoAnterior":      anoAnterior,
		"MesProximo":       mesProximo,
		"AnoProximo":       anoProximo,
	})
}
