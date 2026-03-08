package db

import "log"

func Migrate() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS receitas (
            id             INTEGER PRIMARY KEY AUTOINCREMENT,
            descricao      TEXT    NOT NULL,
            valor          REAL    NOT NULL,
            tipo           TEXT    NOT NULL CHECK(tipo IN ('recorrente', 'pontual')),
            mes_referencia DATE,
            data_inicio    DATE,
            data_fim       DATE,
            criado_em      DATETIME NOT NULL DEFAULT (datetime('now'))
        )`,
		`CREATE TABLE IF NOT EXISTS cartoes (
            id              INTEGER PRIMARY KEY AUTOINCREMENT,
            nome            TEXT    NOT NULL,
            dia_vencimento  INTEGER NOT NULL,
            limite          REAL    NOT NULL,
            criado_em       DATETIME NOT NULL DEFAULT (datetime('now'))
        )`,
		`CREATE TABLE IF NOT EXISTS terceiros (
            id               INTEGER PRIMARY KEY AUTOINCREMENT,
            nome             TEXT NOT NULL,
            limite_liberado  REAL NOT NULL,
            observacao       TEXT,
            criado_em        DATETIME NOT NULL DEFAULT (datetime('now'))
        )`,
		`CREATE TABLE IF NOT EXISTS despesas (
            id               INTEGER PRIMARY KEY AUTOINCREMENT,
            descricao        TEXT    NOT NULL,
            valor_total      REAL    NOT NULL,
            data_compra      DATE    NOT NULL,
            forma_pagamento  TEXT    NOT NULL CHECK(forma_pagamento IN ('cartao', 'dinheiro')),
            cartao_id        INTEGER REFERENCES cartoes(id),
            parcelado        BOOLEAN NOT NULL DEFAULT 0,
            total_parcelas   INTEGER,
            parcela_atual    INTEGER,
            fixa             BOOLEAN NOT NULL DEFAULT 0,
            terceiro_id      INTEGER REFERENCES terceiros(id),
            categoria        TEXT,
            observacao       TEXT,
            criado_em        DATETIME NOT NULL DEFAULT (datetime('now'))
        )`,
		`CREATE TABLE IF NOT EXISTS emprestimos (
            id               INTEGER PRIMARY KEY AUTOINCREMENT,
            descricao        TEXT    NOT NULL,
            valor_total      REAL    NOT NULL,
            total_parcelas   INTEGER NOT NULL,
            taxa_juros       REAL,
            tipo_juros       TEXT CHECK(tipo_juros IN ('simples', 'composto')),
            data_inicio      DATE    NOT NULL,
            forma_pagamento  TEXT    NOT NULL CHECK(forma_pagamento IN ('cartao', 'dinheiro')),
            cartao_id        INTEGER REFERENCES cartoes(id),
            terceiro_id      INTEGER REFERENCES terceiros(id),
            observacao       TEXT,
            criado_em        DATETIME NOT NULL DEFAULT (datetime('now'))
        )`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatalf("Erro ao executar migration: %v", err)
		}
	}

	log.Println("Migrations executadas com sucesso.")
}
