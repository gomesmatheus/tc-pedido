package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	createTables = `
    CREATE TABLE IF NOT EXISTS categoria_produtos (
        id SERIAL PRIMARY KEY,
        descricao VARCHAR(255) NOT NULL UNIQUE
    );

    INSERT INTO categoria_produtos (descricao) VALUES ('Lanche'), ('Acompanhamento'), ('Bebida'), ('Sobremesa') ON CONFLICT (descricao) DO NOTHING;
    SELECT * FROM categoria_produtos;

	CREATE TABLE IF NOT EXISTS produtos (
        id SERIAL PRIMARY KEY,
        categoria_id INTEGER NOT NULL,
        nome VARCHAR(255) NOT NULL UNIQUE,
        descricao VARCHAR(255) NOT NULL,
        preco FLOAT NOT NULL,
        tempo_de_preparo_minutos INTEGER NOT NULL,

        CONSTRAINT fk_categoria_id FOREIGN KEY(categoria_id) REFERENCES categoria_produtos(id)
    );

    CREATE TABLE IF NOT EXISTS pedidos (
        id SERIAL PRIMARY KEY,
        cliente_cpf BIGINT,
        status VARCHAR(255),
        data TIMESTAMP,
        metodo_pagamento VARCHAR(255),
        pagamento_aprovado BOOLEAN DEFAULT FALSE
    );

    CREATE TABLE IF NOT EXISTS produto_pedido (
        produto_id INTEGER NOT NULL,
        pedido_id INTEGER NOT NULL,
        quantidade INTEGER NOT NULL,
        observacao VARCHAR,

        PRIMARY KEY (produto_id, pedido_id),
        CONSTRAINT fk_produto FOREIGN KEY (produto_id) REFERENCES produtos(id),
        CONSTRAINT fk_pedido FOREIGN KEY (pedido_id) REFERENCES pedidos(id)
    );
    `
)

const maxRetries = 10
const retryInterval = 2 * time.Second

func NewPostgresDb(url string) (*pgx.Conn, error) {
	var db *pgx.Conn
	var err error

	for i := 0; i < maxRetries; i++ {
		config, err := pgx.ParseConfig(url)
		if err != nil {
			fmt.Println("Error parsing config:", err)
			return nil, err
		}

		db, err = pgx.ConnectConfig(context.Background(), config)
		if err == nil {
			break
		}
		fmt.Printf("Attempt %d: Error connecting to database: %v\n", i+1, err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL after %d attempts: %w", maxRetries, err)
	}

	if _, err := db.Exec(context.Background(), createTables); err != nil {
		fmt.Println("Error creating table Pedidos:", err)
		return nil, err
	}

	return db, nil
}
