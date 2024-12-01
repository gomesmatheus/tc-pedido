package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gomesmatheus/tc-pedido/domain/entity"
	"github.com/gomesmatheus/tc-pedido/usecase"
)

type PedidoHandler struct {
	pedidoUseCases usecase.PedidoUseCases
}

type PatchPedido struct {
	Status string `json:"status"`
}

func NewPedidoHandler(pedidoUseCases usecase.PedidoUseCases) *PedidoHandler {
	return &PedidoHandler{
		pedidoUseCases: pedidoUseCases,
	}
}

func (c *PedidoHandler) CriacaoPedidoRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			fmt.Println("Error parsing request body")
			w.WriteHeader(400)
			w.Write([]byte("400 bad request"))
			return
		}

		var pedido entity.Pedido
		err = json.Unmarshal(body, &pedido)
		if err != nil {
			fmt.Println("Error parsing request body")
			w.WriteHeader(400)
			w.Write([]byte("400 bad request"))
			return
		}
		fmt.Println(pedido)

		pedido, err = c.pedidoUseCases.CriarPedido(pedido)
		if err != nil {
			fmt.Println("Erro ao cadastrar o pedido", err)
			w.WriteHeader(500)
			w.Write([]byte("Erro ao cadastrar o pedido"))
			return
		}

		w.WriteHeader(201)
		w.Write([]byte(fmt.Sprintf("Pedido inserido com id %d", pedido.Id)))
	} else if r.Method == "GET" {
		pedidos, err := c.pedidoUseCases.RecuperarPedidos()
		if err != nil {
			fmt.Println("Erro ao recuperar pedidos", err)
			w.WriteHeader(500)
			w.Write([]byte("Erro ao recuperar pedidos"))
			return
		}
		response, _ := json.Marshal(pedidos)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(response)
	}

	return
}

func (c *PedidoHandler) AtualizarPedidoRoute(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(strings.Split(r.URL.Path, "/")[3], 10, 64)
	if err != nil {
		fmt.Println("Error parsing request body")
		w.WriteHeader(400)
		w.Write([]byte("400 bad request"))
		return
	}
	if r.Method == "PATCH" {
		var patchPedido PatchPedido
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			fmt.Println("Error parsing request body")
			w.WriteHeader(400)
			w.Write([]byte("400 bad request"))
			return
		}
		err = json.Unmarshal(body, &patchPedido)
		if err != nil {
			fmt.Println("Error parsing request body")
			w.WriteHeader(400)
			w.Write([]byte("400 bad request"))
			return
		}

		err = c.pedidoUseCases.AtualizarStatus(int(id), patchPedido.Status)
		if err != nil {
			fmt.Println("Erro ao atualizar o pedido", err)
			w.WriteHeader(500)
			w.Write([]byte("Erro ao atualizar o pedido"))
			return
		}

		w.WriteHeader(201)
		w.Write([]byte("Pedido atualizado"))
	}
}
