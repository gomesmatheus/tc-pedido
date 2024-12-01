package persistence

import (
	"database/sql"
	"testing"

	"github.com/gomesmatheus/tc-pedido/domain/entity"
	_ "github.com/mattn/go-sqlite3"
)

func setupProdutoTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	_, err = db.Exec(`
	CREATE TABLE produtos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		categoria_id INTEGER NOT NULL,
		nome TEXT NOT NULL,
		descricao TEXT,
		preco REAL NOT NULL,
		tempo_de_preparo_minutos INTEGER NOT NULL
	);
	`)
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	return db
}

func TestCriarProduto(t *testing.T) {
	db := setupProdutoTestDB(t)
	defer db.Close()

	repo := &ProdutoDbMock{Db: db}

	t.Run("Create a valid produto", func(t *testing.T) {
		produto := entity.Produto{
			CategoriaId:    1,
			Nome:           "Pizza Margherita",
			Descricao:      "Classic pizza with tomato, mozzarella, and basil",
			Preco:          29.99,
			TempoDePreparo: 15,
		}

		createdProduto, err := repo.CriarProduto(produto)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if createdProduto.Nome != "Pizza Margherita" {
			t.Errorf("expected produto name to be 'Pizza Margherita', got '%s'", createdProduto.Nome)
		}
	})
}

func TestRecuperarProdutos(t *testing.T) {
	db := setupProdutoTestDB(t)
	defer db.Close()

	repo := &ProdutoDbMock{Db: db}

	// Insert sample data
	_, err := db.Exec(`INSERT INTO produtos (categoria_id, nome, descricao, preco, tempo_de_preparo_minutos) VALUES (?, ?, ?, ?, ?)`,
		1, "Pizza Margherita", "Classic pizza with tomato, mozzarella, and basil", 29.99, 15)
	if err != nil {
		t.Fatalf("failed to insert sample produto: %v", err)
	}

	t.Run("Retrieve produtos by categoria_id", func(t *testing.T) {
		produtos, err := repo.RecuperarProdutos(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(produtos) != 1 {
			t.Errorf("expected 1 produto, got %d", len(produtos))
		}

		if produtos[0].Nome != "Pizza Margherita" {
			t.Errorf("expected produto name to be 'Pizza Margherita', got '%s'", produtos[0].Nome)
		}
	})
}

func TestAtualizarProduto(t *testing.T) {
	db := setupProdutoTestDB(t)
	defer db.Close()

	repo := &ProdutoDbMock{Db: db}

	// Insert sample data
	_, err := db.Exec(`INSERT INTO produtos (categoria_id, nome, descricao, preco, tempo_de_preparo_minutos) VALUES (?, ?, ?, ?, ?)`,
		1, "Pizza Margherita", "Classic pizza with tomato, mozzarella, and basil", 29.99, 15)
	if err != nil {
		t.Fatalf("failed to insert sample produto: %v", err)
	}

	t.Run("Update produto details", func(t *testing.T) {
		produto := entity.Produto{
			CategoriaId:    1,
			Nome:           "Pizza Pepperoni",
			Descricao:      "Spicy pepperoni pizza with cheese",
			Preco:          35,
			TempoDePreparo: 20,
		}

		err := repo.AtualizarProduto(1, produto)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var nome, descricao string
		var preco float64
		err = db.QueryRow(`SELECT nome, descricao, preco FROM produtos WHERE id = ?`, 1).Scan(&nome, &descricao, &preco)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if nome != "Pizza Pepperoni" {
			t.Errorf("expected produto name to be 'Pizza Pepperoni', got '%s'", nome)
		}
		if preco != 35 {
			t.Errorf("expected produto price to be 35, got %f", preco)
		}
	})
}

func TestDeletarProduto(t *testing.T) {
	db := setupProdutoTestDB(t)
	defer db.Close()

	repo := &ProdutoDbMock{Db: db}

	// Insert sample data
	_, err := db.Exec(`INSERT INTO produtos (categoria_id, nome, descricao, preco, tempo_de_preparo_minutos) VALUES (?, ?, ?, ?, ?)`,
		1, "Pizza Margherita", "Classic pizza with tomato, mozzarella, and basil", 29.99, 15)
	if err != nil {
		t.Fatalf("failed to insert sample produto: %v", err)
	}

	t.Run("Delete produto by ID", func(t *testing.T) {
		err := repo.DeletarProduto(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var count int
		err = db.QueryRow(`SELECT COUNT(*) FROM produtos WHERE id = ?`, 1).Scan(&count)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if count != 0 {
			t.Errorf("expected produto to be deleted, but found %d rows", count)
		}
	})
}
