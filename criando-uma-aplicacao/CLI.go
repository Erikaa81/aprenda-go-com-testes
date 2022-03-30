package poquer

import (
	"bufio"
	"io"
	"strings"
	"time"
)

type CLI struct {
	armazenamentoJogador ArmazenamentoJogador
	in                   *bufio.Scanner
	alerter              BlindAlerter
}

type scheduledAlert struct {
	at     time.Duration
	amount int
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func NovoCLI(armazenamento ArmazenamentoJogador, in io.Reader, alerter BlindAlerter) *CLI {
	return &CLI{
		armazenamentoJogador: armazenamento,
		in:                   bufio.NewScanner(in),
		alerter:              alerter,
	}
}
func (cli *CLI) JogarPoquer() {
	cli.sheduleBlindAlerts()
	userInput := cli.readLine()
	cli.armazenamentoJogador.RegistrarVitoria(extrairVencedor(userInput))
}

func (cli *CLI) sheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		cli.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Minute
	}

}
func extrairVencedor(userInput string) string {
	return strings.Replace(userInput, " venceu", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
