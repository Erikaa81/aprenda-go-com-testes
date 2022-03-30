package poquer_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	poquer "github.com/Erikaa81/aprenda-go-com-testes/criando-uma-aplicacao"
)

type scheduledAlert struct {
	at     time.Duration
	amount int
}
type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{at, amount})
}

var dummySpyAlerter = &SpyBlindAlerter{}

func TestCLI(t *testing.T) {

	t.Run("recorda vencedor chris digitado pelo usuario", func(t *testing.T) {
		in := strings.NewReader("Chris venceu\n")
		armazenamentoJogador := &poquer.EsbocoArmazenamentoJogador{}
		blindAlerter := &SpyBlindAlerter{}

		cli := poquer.NovoCLI(armazenamentoJogador, in, blindAlerter)
		cli.JogarPoquer()

		poquer.VerificaVitoriaJogador(t, armazenamentoJogador, "Chris")
	})

	t.Run("recorda vencedor cleo digitado pelo usuario", func(t *testing.T) {
		in := strings.NewReader("Cleo venceu\n")
		armazenamentoJogador := &poquer.EsbocoArmazenamentoJogador{}
		blindAlerter := &SpyBlindAlerter{}

		cli := poquer.NovoCLI(armazenamentoJogador, in, blindAlerter)
		cli.JogarPoquer()

		poquer.VerificaVitoriaJogador(t, armazenamentoJogador, "Cleo")
	})

	t.Run("agnedamento de impressao de valores cegos ", func(t *testing.T) {
		in := strings.NewReader("Chris venceu\n")
		armazenamentoJogador := &poquer.EsbocoArmazenamentoJogador{}
		blindAlerter := &SpyBlindAlerter{}

		cli := poquer.NovoCLI(armazenamentoJogador, in, blindAlerter)
		cli.JogarPoquer()

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {

				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)

			})

		}
	})
}

func assertScheduledAlert(t *testing.T, got, want scheduledAlert) {
	t.Helper()
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
