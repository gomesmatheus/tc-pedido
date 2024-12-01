package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gomesmatheus/tc-pedido/domain/entity"
)

type MockProdutoUseCases struct {
	CriarProdutoFn      func(produto entity.Produto) (entity.Produto, error)
	RecuperarProdutosFn func(categoriaId int) ([]entity.Produto, error)
	AtualizarProdutoFn  func(id int, produto entity.Produto) error
	DeletarProdutoFn    func(id int) error
}

func (m *MockProdutoUseCases) CriarProduto(produto entity.Produto) (entity.Produto, error) {
	return m.CriarProdutoFn(produto)
}

func (m *MockProdutoUseCases) RecuperarProdutos(categoriaId int) ([]entity.Produto, error) {
	return m.RecuperarProdutosFn(categoriaId)
}

func (m *MockProdutoUseCases) AtualizarProduto(id int, produto entity.Produto) error {
	return m.AtualizarProdutoFn(id, produto)
}

func (m *MockProdutoUseCases) DeletarProduto(id int) error {
	return m.DeletarProdutoFn(id)
}

func TestCriacaoProdutoRoute(t *testing.T) {
	mock := &MockProdutoUseCases{
		CriarProdutoFn: func(produto entity.Produto) (entity.Produto, error) {
			return produto, nil
		},
	}
	handler := NewProdutoHandler(mock)

	t.Run("POST - Successful creation", func(t *testing.T) {
		produto := entity.Produto{Nome: "Produto Teste"}
		body, _ := json.Marshal(produto)

		req := httptest.NewRequest("POST", "/produtos", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.CriacaoProdutoRoute(rr, req)

		resp := rr.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}

		responseBody, _ := ioutil.ReadAll(resp.Body)
		if string(responseBody) != "Produto inserido" {
			t.Errorf("unexpected response: %s", responseBody)
		}
	})
}

func TestRecuperarProdutosRoute(t *testing.T) {
	mock := &MockProdutoUseCases{
		RecuperarProdutosFn: func(categoriaId int) ([]entity.Produto, error) {
			return []entity.Produto{
				{Id: 1, Nome: "Produto Teste"},
			}, nil
		},
	}
	handler := NewProdutoHandler(mock)

	t.Run("GET - Successful retrieval", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/produtos/1", nil)
		rr := httptest.NewRecorder()

		handler.RecuperarProdutosRoute(rr, req)

		resp := rr.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		responseBody, _ := ioutil.ReadAll(resp.Body)
		if !strings.Contains(string(responseBody), "Produto Teste") {
			t.Errorf("unexpected response: %s", responseBody)
		}
	})
}

func TestAtualizarProdutoRoute(t *testing.T) {
	mock := &MockProdutoUseCases{
		AtualizarProdutoFn: func(id int, produto entity.Produto) error {
			return nil
		},
	}
	handler := NewProdutoHandler(mock)

	t.Run("PUT - Successful update", func(t *testing.T) {
		produto := entity.Produto{Nome: "Produto Atualizado"}
		body, _ := json.Marshal(produto)

		req := httptest.NewRequest("PUT", "/produtos/1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.RecuperarProdutosRoute(rr, req)

		resp := rr.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})
}

func TestDeletarProdutoRoute(t *testing.T) {
	mock := &MockProdutoUseCases{
		DeletarProdutoFn: func(id int) error {
			return nil
		},
	}
	handler := NewProdutoHandler(mock)

	t.Run("DELETE - Successful deletion", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/produtos/1", nil)
		rr := httptest.NewRecorder()

		handler.RecuperarProdutosRoute(rr, req)

		resp := rr.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})
}
