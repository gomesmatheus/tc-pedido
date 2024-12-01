package pedido_usecase

import (
	"github.com/gomesmatheus/tc-pedido/domain/entity"
	"github.com/gomesmatheus/tc-pedido/infraestructure/persistence"
)

type pedidoUseCases struct {
	database persistence.PedidoRepository
}

func NewPedidoUseCases(pedidoRepository persistence.PedidoRepository) *pedidoUseCases {
	return &pedidoUseCases{
		database: pedidoRepository,
	}
}

func (usecase *pedidoUseCases) CriarPedido(p entity.Pedido) (entity.Pedido, error) {
	return usecase.database.CriarPedido(p)
}
func (usecase *pedidoUseCases) RecuperarPedidos() ([]entity.Pedido, error) {
	return usecase.database.RecuperarPedidos()
}

func (usecase *pedidoUseCases) AtualizarStatus(id int, status string) error {
	return usecase.database.AtualizarStatus(id, status)
}
