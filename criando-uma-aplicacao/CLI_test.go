package poquer_test

import (
	"strings"
	"testing"

	poquer "github.com/Erikaa81/aprenda-go-com-testes/criando-uma-aplicacao"
)

func TestCLI(t *testing.T) {

	t.Run("recorda vencedor chris digitado pelo usuario", func(t *testing.T) {
		in := strings.NewReader("Chris venceu\n")
		armazenamentoJogador := &poquer.EsbocoArmazenamentoJogador{}

		cli := poquer.NovoCLI(armazenamentoJogador, in)
		cli.JogarPoquer()

		poquer.VerificaVitoriaJogador(t, armazenamentoJogador, "Chris")
	})

	t.Run("recorda vencedor cleo digitado pelo usuario", func(t *testing.T) {
		in := strings.NewReader("Cleo venceu\n")
		armazenamentoJogador := &poquer.EsbocoArmazenamentoJogador{}

		cli := poquer.NovoCLI(armazenamentoJogador, in)
		cli.JogarPoquer()

		poquer.VerificaVitoriaJogador(t, armazenamentoJogador, "Cleo")
	})

}
