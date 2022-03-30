package main

import (
	"fmt"
	"log"
	"os"

	poquer "github.com/Erikaa81/aprenda-go-com-testes/criando-uma-aplicacao"
)

const dbFileName = "game.db.json"

func main() {
	armazenamento, close, err := poquer.ArmazenamentoSistemaDeArquivoJogadorAPartirDeArquivo(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Vamos jogar poquer")
	fmt.Println("Digite {Nome} venceu para registrar uma vitoria")
	poquer.NovoCLI(armazenamento, os.Stdin, poquer.BlindAlerterFunc(poquer.StdOutAlerter)).JogarPoquer()
}
