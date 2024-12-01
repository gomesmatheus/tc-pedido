package handlers

import (
	"bytes"
	"errors"
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
	tests := []struct {
		name         string
		body         string
		expectedCode int
		expectedBody string
		mockResponse func() *MockProdutoUseCases
	}{
		{
			name:         "Successful creation",
			body:         `{"nome":"Produto Teste"}`,
			expectedCode: http.StatusCreated,
			expectedBody: "Produto inserido",
			mockResponse: func() *MockProdutoUseCases {
				return &MockProdutoUseCases{
					CriarProdutoFn: func(produto entity.Produto) (entity.Produto, error) {
						return produto, nil
					},
				}
			},
		},
		{
			name:         "Error creating product",
			body:         `{"nome":"Produto Teste"}`,
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Erro ao cadastrar o produto",
			mockResponse: func() *MockProdutoUseCases {
				return &MockProdutoUseCases{
					CriarProdutoFn: func(produto entity.Produto) (entity.Produto, error) {
						return produto, errors.New("error")
					},
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := test.mockResponse()
			handler := NewProdutoHandler(mock)

			req := httptest.NewRequest("POST", "/produtos", bytes.NewBuffer([]byte(test.body)))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.CriacaoProdutoRoute(rr, req)

			resp := rr.Result()
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedCode {
				t.Errorf("expected status %d, got %d", test.expectedCode, resp.StatusCode)
			}

			responseBody, _ := ioutil.ReadAll(resp.Body)
			if string(responseBody) != test.expectedBody {
				t.Errorf("expected body %q, got %q", test.expectedBody, responseBody)
			}
		})
	}
}

func TestRecuperarProdutosRoute(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		expectedCode int
		expectedBody string
		mockResponse func() *MockProdutoUseCases
	}{
		{
			name:         "Successful retrieval",
			url:          "/produtos/1",
			expectedCode: http.StatusOK,
			expectedBody: `[{"id":1,"categoria_id":0,"nome":"Produto Teste","descricao":"","preco":0,"tempo_de_preparo":0}]`,
			mockResponse: func() *MockProdutoUseCases {
				return &MockProdutoUseCases{
					RecuperarProdutosFn: func(categoriaId int) ([]entity.Produto, error) {
						return []entity.Produto{{Id: 1, Nome: "Produto Teste"}}, nil
					},
				}
			},
		},
		{
			name:         "Error retrieving products",
			url:          "/produtos/1",
			expectedCode: http.StatusNotFound,
			expectedBody: "Erro ao recuperar produtos com categoria_id 1",
			mockResponse: func() *MockProdutoUseCases {
				return &MockProdutoUseCases{
					RecuperarProdutosFn: func(categoriaId int) ([]entity.Produto, error) {
						return nil, errors.New("error")
					},
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := test.mockResponse()
			handler := NewProdutoHandler(mock)

			req := httptest.NewRequest("GET", test.url, nil)
			rr := httptest.NewRecorder()

			handler.RecuperarProdutosRoute(rr, req)

			resp := rr.Result()
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedCode {
				t.Errorf("expected status %d, got %d", test.expectedCode, resp.StatusCode)
			}

			responseBody, _ := ioutil.ReadAll(resp.Body)
			if !strings.Contains(string(responseBody), test.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", test.expectedBody, responseBody)
			}
		})
	}
}

func TestAtualizarProdutoRoute(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		body         string
		expectedCode int
		expectedBody string
		mockResponse func() *MockProdutoUseCases
	}{
		{
			name:         "Successful update",
			url:          "/produtos/1",
			body:         `{"nome":"Produto Atualizado"}`,
			expectedCode: http.StatusOK,
			expectedBody: "",
			mockResponse: func() *MockProdutoUseCases {
				return &MockProdutoUseCases{
					AtualizarProdutoFn: func(id int, produto entity.Produto) error {
						return nil
					},
				}
			},
		},
		{
			name:         "Error updating product",
			url:          "/produtos/1",
			body:         `{"nome":"Produto Atualizado"}`,
			expectedCode: http.StatusInternalServerError,
			expectedBody: "500 Erro ao atualizar produto",
			mockResponse: func() *MockProdutoUseCases {
				return &MockProdutoUseCases{
					AtualizarProdutoFn: func(id int, produto entity.Produto) error {
						return errors.New("error")
					},
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := test.mockResponse()
			handler := NewProdutoHandler(mock)

			req := httptest.NewRequest("PUT", test.url, bytes.NewBuffer([]byte(test.body)))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.RecuperarProdutosRoute(rr, req)

			resp := rr.Result()
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedCode {
				t.Errorf("expected status %d, got %d", test.expectedCode, resp.StatusCode)
			}

			responseBody, _ := ioutil.ReadAll(resp.Body)
			if string(responseBody) != test.expectedBody {
				t.Errorf("expected body %q, got %q", test.expectedBody, responseBody)
			}
		})
	}
}

func TestDeletarProdutoRoute(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		expectedCode int
		expectedBody string
		mockResponse func() *MockProdutoUseCases
	}{
		{
			name:         "Successful deletion",
			url:          "/produtos/1",
			expectedCode: http.StatusOK,
			expectedBody: "",
			mockResponse: func() *MockProdutoUseCases {
				return &MockProdutoUseCases{
					DeletarProdutoFn: func(id int) error {
						return nil
					},
				}
			},
		},
		{
			name:         "Error deleting product",
			url:          "/produtos/1",
			expectedCode: http.StatusInternalServerError,
			expectedBody: "500 Erro ao deletar produto",
			mockResponse: func() *MockProdutoUseCases {
				return &MockProdutoUseCases{
					DeletarProdutoFn: func(id int) error {
						return errors.New("error")
					},
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := test.mockResponse()
			handler := NewProdutoHandler(mock)

			req := httptest.NewRequest("DELETE", test.url, nil)
			rr := httptest.NewRecorder()

			handler.RecuperarProdutosRoute(rr, req)

			resp := rr.Result()
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedCode {
				t.Errorf("expected status %d, got %d", test.expectedCode, resp.StatusCode)
			}

			responseBody, _ := ioutil.ReadAll(resp.Body)
			if string(responseBody) != test.expectedBody {
				t.Errorf("expected body %q, got %q", test.expectedBody, responseBody)
			}
		})
	}
}
