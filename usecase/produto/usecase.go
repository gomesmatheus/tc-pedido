package produto_usecase

import (
	"errors"

	"github.com/gomesmatheus/tc-pedido/domain/entity"
	"github.com/gomesmatheus/tc-pedido/infraestructure/persistence"
)

type produtoUseCases struct {
	database persistence.ProdutoRepository
}

func NewProdutoUseCases(ProdutoRepository persistence.ProdutoRepository) *produtoUseCases {
	return &produtoUseCases{
		database: ProdutoRepository,
	}
}

func (usecase *produtoUseCases) CriarProduto(p entity.Produto) (entity.Produto, error) {
	if !isProdutoValido(p) {
		return p, errors.New("Produto inválido")
	}

	return usecase.database.CriarProduto(p)
}

func (usecase *produtoUseCases) RecuperarProdutos(categoriaId int) ([]entity.Produto, error) {
	return usecase.database.RecuperarProdutos(categoriaId)
}

func (usecase *produtoUseCases) AtualizarProduto(id int, p entity.Produto) error {
	if !isProdutoValido(p) {
		return errors.New("Produto inválido")
	}

	return usecase.database.AtualizarProduto(id, p)
}

func (usecase *produtoUseCases) DeletarProduto(id int) error {
	return usecase.database.DeletarProduto(id)
}

func isProdutoValido(p entity.Produto) bool {
	return p.Nome != "" && p.Preco != 0 && p.Descricao != "" && p.CategoriaId != 0 && p.TempoDePreparo != 0
}
