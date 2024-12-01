package produto_usecase

import (
	"errors"
	"testing"

	"github.com/gomesmatheus/tc-pedido/domain/entity"
)

type MockProdutoRepository struct {
	CriarProdutoMock      func(p entity.Produto) (entity.Produto, error)
	RecuperarProdutosMock func(categoriaId int) ([]entity.Produto, error)
	AtualizarProdutoMock  func(id int, p entity.Produto) error
	DeletarProdutoMock    func(id int) error
}

func (m *MockProdutoRepository) CriarProduto(p entity.Produto) (entity.Produto, error) {
	return m.CriarProdutoMock(p)
}

func (m *MockProdutoRepository) RecuperarProdutos(categoriaId int) ([]entity.Produto, error) {
	return m.RecuperarProdutosMock(categoriaId)
}

func (m *MockProdutoRepository) AtualizarProduto(id int, p entity.Produto) error {
	return m.AtualizarProdutoMock(id, p)
}

func (m *MockProdutoRepository) DeletarProduto(id int) error {
	return m.DeletarProdutoMock(id)
}

func TestProdutoUseCases(t *testing.T) {
	mockRepo := &MockProdutoRepository{
		CriarProdutoMock: func(p entity.Produto) (entity.Produto, error) {
			if p.Nome == "Invalid" {
				return p, errors.New("Produto inválido")
			}
			return p, nil
		},
		RecuperarProdutosMock: func(categoriaId int) ([]entity.Produto, error) {
			if categoriaId == 0 {
				return nil, errors.New("Categoria inválida")
			}
			return []entity.Produto{
				{Id: 1, Nome: "Produto1", CategoriaId: categoriaId},
				{Id: 2, Nome: "Produto2", CategoriaId: categoriaId},
			}, nil
		},
		AtualizarProdutoMock: func(id int, p entity.Produto) error {
			if id == 0 || p.Nome == "Invalid" {
				return errors.New("Produto inválido")
			}
			return nil
		},
		DeletarProdutoMock: func(id int) error {
			if id == 0 {
				return errors.New("ID inválido")
			}
			return nil
		},
	}

	usecase := NewProdutoUseCases(mockRepo)

	t.Run("CriarProduto - Valid Product", func(t *testing.T) {
		produto := entity.Produto{Nome: "Produto1", CategoriaId: 1, Preco: 10.0, Descricao: "Desc", TempoDePreparo: 15}
		_, err := usecase.CriarProduto(produto)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("CriarProduto - Invalid Product", func(t *testing.T) {
		produto := entity.Produto{Nome: "Invalid", CategoriaId: 0}
		_, err := usecase.CriarProduto(produto)
		if err == nil || err.Error() != "Produto inválido" {
			t.Errorf("expected 'Produto inválido', got %v", err)
		}
	})

	t.Run("RecuperarProdutos - Valid Category", func(t *testing.T) {
		_, err := usecase.RecuperarProdutos(1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("RecuperarProdutos - Invalid Category", func(t *testing.T) {
		_, err := usecase.RecuperarProdutos(0)
		if err == nil || err.Error() != "Categoria inválida" {
			t.Errorf("expected 'Categoria inválida', got %v", err)
		}
	})

	t.Run("AtualizarProduto - Valid Product", func(t *testing.T) {
		produto := entity.Produto{Nome: "Produto1", CategoriaId: 1, Preco: 10.0, Descricao: "Desc", TempoDePreparo: 15}
		err := usecase.AtualizarProduto(1, produto)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("AtualizarProduto - Invalid Product", func(t *testing.T) {
		produto := entity.Produto{Nome: "Invalid", CategoriaId: 0}
		err := usecase.AtualizarProduto(0, produto)
		if err == nil || err.Error() != "Produto inválido" {
			t.Errorf("expected 'Produto inválido', got %v", err)
		}
	})

	t.Run("DeletarProduto - Valid ID", func(t *testing.T) {
		err := usecase.DeletarProduto(1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("DeletarProduto - Invalid ID", func(t *testing.T) {
		err := usecase.DeletarProduto(0)
		if err == nil || err.Error() != "ID inválido" {
			t.Errorf("expected 'ID inválido', got %v", err)
		}
	})
}
