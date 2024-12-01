package usecase

import "github.com/gomesmatheus/tc-pedido/domain/entity"

type ProdutoUseCases interface {
	CriarProduto(p entity.Produto) (entity.Produto, error)
	RecuperarProdutos(categoriaId int) ([]entity.Produto, error)
	AtualizarProduto(id int, p entity.Produto) error
	DeletarProduto(id int) error
}

type PedidoUseCases interface {
	CriarPedido(entity.Pedido) (entity.Pedido, error)
	RecuperarPedidos() ([]entity.Pedido, error)
	AtualizarStatus(id int, status string) error
}
