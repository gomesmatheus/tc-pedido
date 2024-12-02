package pedido_usecase

import (
	"fmt"
	"log"
	"net/http"

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
	if !validarCliente(p.Cpf) {
		return p, fmt.Errorf("Client with CPF %d is not registered", p.Cpf)
	}

	return usecase.database.CriarPedido(p)
}
func (usecase *pedidoUseCases) RecuperarPedidos() ([]entity.Pedido, error) {
	return usecase.database.RecuperarPedidos()
}

func (usecase *pedidoUseCases) AtualizarStatus(id int, status string) error {
	return usecase.database.AtualizarStatus(id, status)
}

func validarCliente(cpf int64) bool {
	// Define the service URL and endpoint
	baseURL := "http://svc-cliente-app"
	port := 80
	endpoint := fmt.Sprintf("%s:%d/cliente/%d", baseURL, port, cpf)

	// Create an HTTP GET request
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Fatalf("Failed to call cliente-app: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	return resp.StatusCode == 200
}
