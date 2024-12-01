package persistence

import (
	"database/sql"
	"fmt"

	"github.com/gomesmatheus/tc-pedido/domain/entity"
	_ "github.com/mattn/go-sqlite3"
)

type ProdutoDbMock struct {
	Db *sql.DB
}

func (repo *ProdutoDbMock) CriarProduto(p entity.Produto) (entity.Produto, error) {
	_, err := repo.Db.Exec(
		"INSERT INTO produtos (categoria_id, nome, descricao, preco, tempo_de_preparo_minutos) VALUES (?, ?, ?, ?, ?)",
		p.CategoriaId, p.Nome, p.Descricao, p.Preco, p.TempoDePreparo,
	)
	if err != nil {
		return p, fmt.Errorf("erro ao inserir produto na base de dados: %v", err)
	}
	return p, nil
}

func (repo *ProdutoDbMock) RecuperarProdutos(categoriaId int) ([]entity.Produto, error) {
	var produtos []entity.Produto

	rows, err := repo.Db.Query(
		"SELECT id, categoria_id, nome, descricao, preco, tempo_de_preparo_minutos FROM produtos WHERE categoria_id = ?",
		categoriaId,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos por categoria_id (%d): %v", categoriaId, err)
	}
	defer rows.Close()

	fmt.Println(rows)
	for rows.Next() {
		var p entity.Produto
		var id sql.NullInt64
		var categoriaId2 sql.NullInt64
		var nome, descricao sql.NullString
		var preco sql.NullFloat64
		var tempoDePreparo sql.NullInt64

		if err := rows.Scan(&id, &categoriaId2, &nome, &descricao, &preco, &tempoDePreparo); err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("erro ao fazer scanning de produto: %v", err)
		}
		p.Id = int(id.Int64)
		p.CategoriaId = int(categoriaId2.Int64)
		p.Nome = nome.String
		p.Descricao = descricao.String
		p.Preco = float32(preco.Float64)
		p.TempoDePreparo = int(tempoDePreparo.Int64)
		produtos = append(produtos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar pelos produtos: %v", err)
	}

	return produtos, nil
}

func (repo *ProdutoDbMock) AtualizarProduto(id int, p entity.Produto) error {
	_, err := repo.Db.Exec(
		"UPDATE produtos SET categoria_id = ?, nome = ?, descricao = ?, preco = ?, tempo_de_preparo_minutos = ? WHERE id = ?",
		p.CategoriaId, p.Nome, p.Descricao, p.Preco, p.TempoDePreparo, id,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar produto na base de dados: %v", err)
	}
	return nil
}

func (repo *ProdutoDbMock) DeletarProduto(id int) error {
	_, err := repo.Db.Exec("DELETE FROM produtos WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("erro ao deletar produto da base de dados: %v", err)
	}
	return nil
}
