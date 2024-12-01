package database

import (
	"database/sql"
	"log"
)

const (
	createSqliteTables = `
        CREATE TABLE IF NOT EXISTS produtos (
            id INTEGER PRIMARY KEY AUTOINCREMENT, -- Use AUTOINCREMENT for auto-generated IDs
            categoria_id INTEGER NOT NULL,
            nome TEXT NOT NULL UNIQUE, -- Replace VARCHAR with TEXT (SQLite does not differentiate)
            descricao TEXT NOT NULL,
            preco REAL NOT NULL, -- Use REAL instead of FLOAT
            tempo_de_preparo_minutos INTEGER NOT NULL,
            FOREIGN KEY (categoria_id) REFERENCES categoria_produtos(id)
        );

        CREATE TABLE IF NOT EXISTS pedidos (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            cliente_cpf BIGINT,
            status TEXT,
            data TIMESTAMP,
            metodo_pagamento TEXT,
            pagamento_aprovado BOOLEAN DEFAULT 0, -- Use 0/1 for BOOLEAN
            FOREIGN KEY (cliente_cpf) REFERENCES clientes(cpf)
        );

        CREATE TABLE IF NOT EXISTS produto_pedido (
            produto_id INTEGER NOT NULL,
            pedido_id INTEGER NOT NULL,
            quantidade INTEGER NOT NULL,
            observacao TEXT,
            PRIMARY KEY (produto_id, pedido_id),
            FOREIGN KEY (produto_id) REFERENCES produtos(id),
            FOREIGN KEY (pedido_id) REFERENCES pedidos(id)
        );
    `
)

func NewSqliteDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	// db, err := sql.Open("sqlite3", "./local.db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	_, err = db.Exec(createSqliteTables)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	return db
}
