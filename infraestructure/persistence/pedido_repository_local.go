package persistence

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gomesmatheus/tc-pedido/domain/entity"
	_ "github.com/mattn/go-sqlite3"
)

type PedidoDbMock struct {
	Db *sql.DB
}

const (
	QUERY_PEDIDOS_SQLITE = `
        SELECT
            A.id,
            A.cliente_cpf,
            A.status,
            A.metodo_pagamento,
            A.pagamento_aprovado,
            B.produto_id,
            B.quantidade,
            B.observacao
        FROM pedidos A
        INNER JOIN produto_pedido B ON A.id = B.pedido_id;
    `
)

func (repo *PedidoDbMock) CriarPedido(p entity.Pedido) (entity.Pedido, error) {
	var idPedido int
	tx, err := repo.Db.Begin()
	if err != nil {
		return p, fmt.Errorf("error starting transaction: %v", err)
	}

	err = tx.QueryRow("INSERT INTO pedidos (cliente_cpf, status, data, metodo_pagamento) VALUES (?, ?, ?, ?) RETURNING id",
		p.Cpf, "Recebido", time.Now(), p.MetodoPagamento).Scan(&idPedido)
	if err != nil {
		tx.Rollback()
		return p, fmt.Errorf("error inserting pedido: %v", err)
	}

	for _, pp := range p.Produtos {
		_, err := tx.Exec("INSERT INTO produto_pedido (produto_id, pedido_id, quantidade, observacao) VALUES (?, ?, ?, ?)",
			pp.ProdutoId, idPedido, pp.Quantidade, pp.Observacao)
		if err != nil {
			tx.Rollback()
			return p, fmt.Errorf("error inserting produto_pedido: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return p, fmt.Errorf("error committing transaction: %v", err)
	}

	p.Id = idPedido
	return p, nil
}

func (repo *PedidoDbMock) RecuperarPedidos() ([]entity.Pedido, error) {
	var pedidos []entity.Pedido
	rows, err := repo.Db.Query(QUERY_PEDIDOS_SQLITE)
	if err != nil {
		return nil, fmt.Errorf("error querying pedidos: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r PedidoRow
		err := rows.Scan(&r.Id, &r.Cpf, &r.Status, &r.MetodoPagamento, &r.PagamentoAprovado, &r.ProdutoId, &r.Quantidade, &r.Observacao)
		if err != nil {
			return nil, fmt.Errorf("error scanning pedido: %v", err)
		}

		pedidoJaExiste := false
		for i, p := range pedidos {
			if p.Id == r.Id {
				pedidoJaExiste = true
				pedidos[i].Produtos = append(p.Produtos, entity.ProdutoPedido{
					ProdutoId:  r.ProdutoId,
					Quantidade: r.Quantidade,
					Observacao: r.Observacao,
				})
				break
			}
		}

		if !pedidoJaExiste {
			pedidos = append(pedidos, entity.Pedido{
				Id:                r.Id,
				Cpf:               r.Cpf,
				Status:            r.Status,
				MetodoPagamento:   r.MetodoPagamento,
				PagamentoAprovado: r.PagamentoAprovado,
				Produtos: []entity.ProdutoPedido{
					{
						ProdutoId:  r.ProdutoId,
						Quantidade: r.Quantidade,
						Observacao: r.Observacao,
					},
				},
			})
		}
	}

	return pedidos, nil
}

func (repo *PedidoDbMock) AtualizarStatus(idPedido int, status string) error {
	_, err := repo.Db.Exec("UPDATE pedidos SET status = ? WHERE id = ?", status, idPedido)
	if err != nil {
		return fmt.Errorf("error updating pedido status: %v", err)
	}
	return nil
}

func (repo *PedidoDbMock) AtualizarPagamento(idPedido int, pagamentoAprovado bool) error {
	_, err := repo.Db.Exec("UPDATE pedidos SET pagamento_aprovado = ? WHERE id = ?", pagamentoAprovado, idPedido)
	if err != nil {
		return fmt.Errorf("error updating pedido pagamento_aprovado: %v", err)
	}
	return nil
}
