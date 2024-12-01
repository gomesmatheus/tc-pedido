package entity

type Pedido struct {
	Id                int             `json:"id"`
	Cpf               int64           `json:"cpf"`
	Produtos          []ProdutoPedido `json:"produtos"`
	Status            string          `json:"status"`
	MetodoPagamento   string          `json:"metodo_de_pagamento"`
	PagamentoAprovado bool            `json:"pagamento_aprovado"`
}

type ProdutoPedido struct {
	ProdutoId  int    `json:"produto_id"`
	Quantidade int    `json:"quantidade"`
	Observacao string `json:"observacao"`
}
