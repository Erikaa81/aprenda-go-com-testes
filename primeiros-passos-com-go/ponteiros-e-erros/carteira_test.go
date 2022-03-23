package ponteiroseerros

import (
	"errors"
	"fmt"
	"testing"
)

type Bitcoin int

type Stringer interface {
	String() string
}

type Carteira struct {
	saldo Bitcoin
}

func (c *Carteira) Depositar(quantidade Bitcoin) {
	c.saldo += quantidade
}

func (c *Carteira) Saldo() Bitcoin {
	return c.saldo
}

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}

var ErroSaldoInsuficiente = errors.New("não é possível retirar: saldo insuficiente")

func (c *Carteira) Retirar(quantidade Bitcoin) error {

	if quantidade > c.saldo {
		return ErroSaldoInsuficiente
	}

	c.saldo -= quantidade
	return nil
}
func TestCarteira(t *testing.T) {
	t.Run("Depositar", func(t *testing.T) {
		carteira := Carteira{}
		carteira.Depositar(Bitcoin(10))

		confirmaSaldo(t, carteira, Bitcoin(10))
	})

	t.Run("Retirar com saldo suficiente", func(t *testing.T) {
		carteira := Carteira{Bitcoin(20)}
		erro := carteira.Retirar(Bitcoin(10))

		confirmaSaldo(t, carteira, Bitcoin(10))
		confirmaErroInexistente(t, erro)
	})

	t.Run("Retirar com saldo insuficiente", func(t *testing.T) {
		saldoInicial := Bitcoin(20)
		carteira := Carteira{saldoInicial}
		erro := carteira.Retirar(Bitcoin(100))

		confirmaSaldo(t, carteira, saldoInicial)
		confirmaErro(t, erro, ErroSaldoInsuficiente)
	})
}

func confirmaSaldo(t *testing.T, carteira Carteira, esperado Bitcoin) {
	t.Helper()
	resultado := carteira.Saldo()

	if resultado != esperado {
		t.Errorf("resultado %s, esperado %s", resultado, esperado)
	}
}

func confirmaErroInexistente(t *testing.T, resultado error) {
	t.Helper()
	if resultado != nil {
		t.Fatal("erro inesperado recebido")
	}
}

func confirmaErro(t *testing.T, resultado error, esperado error) {
	t.Helper()
	if resultado == nil {
		t.Fatal("esperava um erro, mas nenhum ocorreu")
	}

	if resultado != esperado {
		t.Errorf("erro resultado %s, erro esperado %s", resultado, esperado)
	}
}
