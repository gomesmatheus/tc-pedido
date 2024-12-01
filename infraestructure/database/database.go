package database

import "github.com/gomesmatheus/tc-pedido/infraestructure/persistence"

func NewPedidoRepository() (persistence.PedidoRepository, error) {
	pgDb, _ := NewPostgresDb("postgres://postgres:123@pedido-db:5432/postgres")

	return &persistence.PedidoDbConnection{
		Db: pgDb,
	}, nil
}

func NewPedidoRepositoryLocal() (persistence.PedidoRepository, error) {
	db := NewSqliteDB()

	return &persistence.PedidoDbMock{
		Db: db,
	}, nil
}

func NewProdutoRepository() (persistence.ProdutoRepository, error) {
	pgDb, _ := NewPostgresDb("postgres://postgres:123@pedido-db:5432/postgres")

	return &persistence.ProdutoDbConnection{
		Db: pgDb,
	}, nil
}

func NewProdutoRepositoryLocal() (persistence.ProdutoRepository, error) {
	db := NewSqliteDB()

	return &persistence.ProdutoDbMock{
		Db: db,
	}, nil
}
