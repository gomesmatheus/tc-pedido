package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gomesmatheus/tc-pedido/domain/entity"
)

type mockPedidoUseCases struct {
	CreateResult    entity.Pedido
	FetchPedidos    []entity.Pedido
	UpdatedStatus   error
	FetchPedidosErr error
	CreateErr       error
}

func (m *mockPedidoUseCases) CriarPedido(p entity.Pedido) (entity.Pedido, error) {
	return m.CreateResult, m.CreateErr
}

func (m *mockPedidoUseCases) RecuperarPedidos() ([]entity.Pedido, error) {
	return m.FetchPedidos, m.FetchPedidosErr
}

func (m *mockPedidoUseCases) AtualizarStatus(id int, status string) error {
	return m.UpdatedStatus
}

func TestCriacaoPedidoRoute(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Successful POST",
			method:       "POST",
			body:         `{"cpf":12345,"status":"Pending","metodo_pagamento":"card"}`,
			expectedCode: 201,
			expectedBody: "Pedido inserido com id 1",
		},
		{
			name:         "POST with invalid JSON",
			method:       "POST",
			body:         "{invalid-json}",
			expectedCode: 400,
			expectedBody: "400 bad request",
		},
		{
			name:         "POST with internal error",
			method:       "POST",
			body:         `{"cpf":12345,"status":"Pending","metodo_pagamento":"card"}`,
			expectedCode: 500,
			expectedBody: "Erro ao cadastrar o pedido",
		},
		{
			name:         "Successful GET",
			method:       "GET",
			expectedCode: 200,
			expectedBody: `[{"id":1,"cpf":12345,"produtos":null,"status":"Pending","metodo_de_pagamento":"card","pagamento_aprovado":false}]`,
		},
		{
			name:         "GET with internal error",
			method:       "GET",
			expectedCode: 500,
			expectedBody: "Erro ao recuperar pedidos",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := &mockPedidoUseCases{
				CreateResult: entity.Pedido{Id: 1, Cpf: 12345, Status: "Pending", MetodoPagamento: "card"},
				FetchPedidos: []entity.Pedido{{Id: 1, Cpf: 12345, Status: "Pending", MetodoPagamento: "card"}},
				CreateErr:    nil,
			}

			// Simulate errors for specific test cases
			if test.name == "POST with internal error" {
				mockUsecase.CreateErr = fmt.Errorf("internal error")
			} else if test.name == "GET with internal error" {
				mockUsecase.FetchPedidosErr = fmt.Errorf("internal error")
			}

			handler := NewPedidoHandler(mockUsecase)

			req := httptest.NewRequest(test.method, "/pedidos", bytes.NewBuffer([]byte(test.body)))
			rec := httptest.NewRecorder()

			handler.CriacaoPedidoRoute(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			if resp.StatusCode != test.expectedCode {
				t.Errorf("Expected status code %d, got %d", test.expectedCode, resp.StatusCode)
			}

			if string(body) != test.expectedBody {
				t.Errorf("Expected body %q, got %q", test.expectedBody, string(body))
			}
		})
	}
}

func TestAtualizarPedidoRoute(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		body         string
		url          string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Successful PATCH",
			method:       "PATCH",
			body:         `{"status":"Completed"}`,
			url:          "/pedido/atualizar/1",
			expectedCode: 201,
			expectedBody: "Pedido atualizado",
		},
		{
			name:         "PATCH with invalid JSON",
			method:       "PATCH",
			body:         "{invalid-json}",
			url:          "/pedido/atualizar/1",
			expectedCode: 400,
			expectedBody: "400 bad request",
		},
		{
			name:         "PATCH with internal error",
			method:       "PATCH",
			body:         `{"status":"Completed"}`,
			url:          "/pedido/atualizar/1",
			expectedCode: 500,
			expectedBody: "Erro ao atualizar o pedido",
		},
		{
			name:         "PATCH with missing ID",
			method:       "PATCH",
			body:         `{"status":"Completed"}`,
			url:          "/pedido/atualizar/abc",
			expectedCode: 400,
			expectedBody: "400 bad request",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := &mockPedidoUseCases{
				UpdatedStatus: nil,
			}

			// Simulate errors for specific test cases
			if test.name == "PATCH with internal error" {
				mockUsecase.UpdatedStatus = fmt.Errorf("internal error")
			}

			handler := NewPedidoHandler(mockUsecase)

			req := httptest.NewRequest(test.method, test.url, bytes.NewBuffer([]byte(test.body)))
			rec := httptest.NewRecorder()

			handler.AtualizarPedidoRoute(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			if resp.StatusCode != test.expectedCode {
				t.Errorf("Expected status code %d, got %d", test.expectedCode, resp.StatusCode)
			}

			if string(body) != test.expectedBody {
				t.Errorf("Expected body %q, got %q", test.expectedBody, string(body))
			}
		})
	}
}
