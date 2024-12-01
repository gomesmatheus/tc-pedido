package entity

type Produto struct {
	Id             int     `json:"id"`
	CategoriaId    int     `json:"categoria_id"`
	Nome           string  `json:"nome"`
	Descricao      string  `json:"descricao"`
	Preco          float32 `json:"preco"`
	TempoDePreparo int     `json:"tempo_de_preparo"`
}
