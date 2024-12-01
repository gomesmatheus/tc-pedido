package persistence

import "github.com/gomesmatheus/tc-pedido/domain/entity"

type ProdutoRepository interface {
	CriarProduto(p entity.Produto) (entity.Produto, error)
	RecuperarProdutos(categoriaId int) ([]entity.Produto, error)
	AtualizarProduto(id int, p entity.Produto) error
	DeletarProduto(id int) error
}

type PedidoRepository interface {
	CriarPedido(entity.Pedido) (entity.Pedido, error)
	RecuperarPedidos() ([]entity.Pedido, error)
	AtualizarStatus(id int, status string) error
	AtualizarPagamento(id int, status bool) error
}
