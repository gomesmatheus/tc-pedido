package main

import (
	"fmt"
	"log"
	"net/http"

	handlers "github.com/gomesmatheus/tc-pedido/delivery/http/handler"
	"github.com/gomesmatheus/tc-pedido/infraestructure/database"
	pedido_usecase "github.com/gomesmatheus/tc-pedido/usecase/pedido"
	produto_usecase "github.com/gomesmatheus/tc-pedido/usecase/produto"
)

func main() {
	pedidoRepository, err := database.NewPedidoRepository()
	produtoRepository, err := database.NewProdutoRepository()

	// pedidoRepository, err := database.NewPedidoRepositoryLocal()
	// produtoRepository, err := database.NewProdutoRepositoryLocal()

	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	produtoUseCases := produto_usecase.NewProdutoUseCases(produtoRepository)
	produtoHandler := handlers.NewProdutoHandler(produtoUseCases)

	pedidoUseCases := pedido_usecase.NewPedidoUseCases(pedidoRepository)
	pedidoHandler := handlers.NewPedidoHandler(pedidoUseCases)

	http.HandleFunc("/produto", produtoHandler.CriacaoProdutoRoute)
	http.HandleFunc("/produto/", produtoHandler.RecuperarProdutosRoute)
	http.HandleFunc("/pedido", pedidoHandler.CriacaoPedidoRoute)
	http.HandleFunc("/pedido/atualizar/", pedidoHandler.AtualizarPedidoRoute)

	fmt.Println("Pedido ms running!")
	log.Fatal(http.ListenAndServe(":3333", nil))
}
