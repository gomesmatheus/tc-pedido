package persistence

import (
	"database/sql"
	"testing"
	"time"

	"github.com/gomesmatheus/tc-pedido/domain/entity"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE pedidos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cliente_cpf BIGINT NOT NULL,
		status TEXT NOT NULL,
		data TIMESTAMP NOT NULL,
		metodo_pagamento TEXT NOT NULL,
		pagamento_aprovado BOOLEAN DEFAULT FALSE
	);

	CREATE TABLE produto_pedido (
		produto_id INTEGER NOT NULL,
		pedido_id INTEGER NOT NULL,
		quantidade INTEGER NOT NULL,
		observacao TEXT,
		PRIMARY KEY (produto_id, pedido_id)
	);
	`)
	if err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	return db
}

func TestCriarPedido(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &PedidoDbMock{Db: db}

	t.Run("Create a valid pedido", func(t *testing.T) {
		pedido := entity.Pedido{
			Cpf:             123456789,
			MetodoPagamento: "Cartão",
			Produtos: []entity.ProdutoPedido{
				{ProdutoId: 1, Quantidade: 2, Observacao: "Extra cheese"},
				{ProdutoId: 2, Quantidade: 1, Observacao: "No onions"},
			},
		}

		createdPedido, err := repo.CriarPedido(pedido)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if createdPedido.Id == 0 {
			t.Errorf("expected pedido ID to be generated, got 0")
		}
	})
}

func TestRecuperarPedidos(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &PedidoDbMock{Db: db}

	// Insert sample data
	_, err := db.Exec(`INSERT INTO pedidos (cliente_cpf, status, data, metodo_pagamento) VALUES (?, ?, ?, ?)`,
		123456789, "Recebido", time.Now(), "Cartão")
	if err != nil {
		t.Fatalf("failed to insert sample pedido: %v", err)
	}
	_, err = db.Exec(`INSERT INTO produto_pedido (produto_id, pedido_id, quantidade, observacao) VALUES (?, ?, ?, ?)`,
		1, 1, 2, "Extra cheese")
	if err != nil {
		t.Fatalf("failed to insert sample produto_pedido: %v", err)
	}

	t.Run("Retrieve pedidos", func(t *testing.T) {
		pedidos, err := repo.RecuperarPedidos()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(pedidos) != 1 {
			t.Errorf("expected 1 pedido, got %d", len(pedidos))
		}

		if pedidos[0].Produtos[0].ProdutoId != 1 {
			t.Errorf("expected produto ID 1, got %d", pedidos[0].Produtos[0].ProdutoId)
		}
	})
}

func TestAtualizarStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &PedidoDbMock{Db: db}

	// Insert sample data
	_, err := db.Exec(`INSERT INTO pedidos (cliente_cpf, status, data, metodo_pagamento) VALUES (?, ?, ?, ?)`,
		123456789, "Recebido", time.Now(), "Cartão")
	if err != nil {
		t.Fatalf("failed to insert sample pedido: %v", err)
	}

	t.Run("Update pedido status", func(t *testing.T) {
		err := repo.AtualizarStatus(1, "Em Andamento")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var status string
		err = db.QueryRow(`SELECT status FROM pedidos WHERE id = ?`, 1).Scan(&status)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if status != "Em Andamento" {
			t.Errorf("expected status 'Em Andamento', got '%s'", status)
		}
	})
}

func TestAtualizarPagamento(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &PedidoDbMock{Db: db}

	// Insert sample data
	_, err := db.Exec(`INSERT INTO pedidos (cliente_cpf, status, data, metodo_pagamento) VALUES (?, ?, ?, ?)`,
		123456789, "Recebido", time.Now(), "Cartão")
	if err != nil {
		t.Fatalf("failed to insert sample pedido: %v", err)
	}

	t.Run("Update pagamento_aprovado", func(t *testing.T) {
		err := repo.AtualizarPagamento(1, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var pagamentoAprovado bool
		err = db.QueryRow(`SELECT pagamento_aprovado FROM pedidos WHERE id = ?`, 1).Scan(&pagamentoAprovado)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !pagamentoAprovado {
			t.Errorf("expected pagamento_aprovado to be true, got false")
		}
	})
}
