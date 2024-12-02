package pedido_usecase

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gomesmatheus/tc-pedido/domain/entity"
)

type MockPedidoRepository struct {
	CreatePedidoMock       func(p entity.Pedido) (entity.Pedido, error)
	RecuperarPedidosMock   func() ([]entity.Pedido, error)
	AtualizarStatusMock    func(id int, status string) error
	AtualizarPagamentoMock func(id int, aprovado bool) error
}

func (m *MockPedidoRepository) CriarPedido(p entity.Pedido) (entity.Pedido, error) {
	return m.CreatePedidoMock(p)
}

func (m *MockPedidoRepository) RecuperarPedidos() ([]entity.Pedido, error) {
	return m.RecuperarPedidosMock()
}

func (m *MockPedidoRepository) AtualizarStatus(id int, status string) error {
	return m.AtualizarStatusMock(id, status)
}

func (m *MockPedidoRepository) AtualizarPagamento(id int, aprovado bool) error {
	return m.AtualizarPagamentoMock(id, aprovado)
}

func TestRecuperarPedidos(t *testing.T) {
	mockPedidos := []entity.Pedido{
		{Id: 1, Cpf: 12345678900, Status: "Recebido", MetodoPagamento: "Cartão"},
		{Id: 2, Cpf: 98765432100, Status: "Recebido", MetodoPagamento: "Boleto"},
	}

	mockRepo := &MockPedidoRepository{
		RecuperarPedidosMock: func() ([]entity.Pedido, error) {
			return mockPedidos, nil
		},
	}

	usecase := NewPedidoUseCases(mockRepo)

	t.Run("Should retrieve pedidos successfully", func(t *testing.T) {
		pedidos, err := usecase.RecuperarPedidos()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(pedidos, mockPedidos) {
			t.Errorf("expected %v, got %v", mockPedidos, pedidos)
		}
	})
}

func TestAtualizarStatus(t *testing.T) {
	mockRepo := &MockPedidoRepository{
		AtualizarStatusMock: func(id int, status string) error {
			if id == 2 {
				return errors.New("failed to update status")
			}
			return nil
		},
	}

	usecase := NewPedidoUseCases(mockRepo)

	t.Run("Should update status successfully", func(t *testing.T) {
		err := usecase.AtualizarStatus(1, "Concluído")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("Should fail to update status", func(t *testing.T) {
		err := usecase.AtualizarStatus(2, "Cancelado")
		if err == nil || err.Error() != "failed to update status" {
			t.Errorf("expected 'failed to update status', got %v", err)
		}
	})
}
