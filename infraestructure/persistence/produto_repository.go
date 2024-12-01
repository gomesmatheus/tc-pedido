package persistence

import (
	"context"
	"fmt"

	"github.com/gomesmatheus/tc-pedido/domain/entity"
	"github.com/jackc/pgx/v5"
)

type ProdutoDbConnection struct {
	Db *pgx.Conn
}

func (repo *ProdutoDbConnection) CriarProduto(p entity.Produto) (entity.Produto, error) {
	_, err := repo.Db.Exec(context.Background(), "INSERT INTO produtos (categoria_id, nome, descricao, preco, tempo_de_preparo_minutos) VALUES ($1, $2, $3, $4, $5)", p.CategoriaId, p.Nome, p.Descricao, p.Preco, p.TempoDePreparo)
	if err != nil {
		fmt.Println("Erro ao inserir produto na base de dados", err)
	}
	return p, err
}

func (repo *ProdutoDbConnection) RecuperarProdutos(categoriaId int) ([]entity.Produto, error) {
	var produtos []entity.Produto
	rows, err := repo.Db.Query(context.Background(), "SELECT id, categoria_id, nome, descricao, preco, tempo_de_preparo_minutos FROM produtos WHERE categoria_id = $1", categoriaId)
	defer rows.Close()
	if err != nil {
		fmt.Println("Erro ao buscar por categoria_id", categoriaId)
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		var p entity.Produto
		if err = rows.Scan(&p.Id, &p.CategoriaId, &p.Nome, &p.Descricao, &p.Preco, &p.TempoDePreparo); err != nil {
			fmt.Println("Erro fazendo scanning de produto")
			fmt.Println(err)
			return nil, err
		}
		produtos = append(produtos, p)
	}

	return produtos, err
}

func (repo *ProdutoDbConnection) AtualizarProduto(id int, p entity.Produto) error {
	_, err := repo.Db.Exec(context.Background(), "UPDATE produtos set categoria_id = $1, nome = $2, descricao = $3, preco = $4, tempo_de_preparo_minutos = $5 WHERE id = $6", p.CategoriaId, p.Nome, p.Descricao, p.Preco, p.TempoDePreparo, id)
	if err != nil {
		fmt.Println("Erro ao atualizar produto na base de dados", err)
	}
	return err
}

func (repo *ProdutoDbConnection) DeletarProduto(id int) error {
	_, err := repo.Db.Exec(context.Background(), "DELETE FROM produtos WHERE id = $1", id)
	if err != nil {
		fmt.Println("Erro ao deletar produto da base de dados", err)
	}
	return err
}
