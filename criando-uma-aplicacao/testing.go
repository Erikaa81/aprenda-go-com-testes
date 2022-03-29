package poquer

import (
	"testing"
)

type EsbocoArmazenamentoJogador struct {
	pontuacoes        map[string]int
	registrosVitorias []string
	liga              []Jogador
}

func (s *EsbocoArmazenamentoJogador) ObterPontuacaoDeJogador(nome string) int {
	pontuacao := s.pontuacoes[nome]
	return pontuacao
}

func (e *EsbocoArmazenamentoJogador) RegistrarVitoria(nome string) {
	e.registrosVitorias = append(e.registrosVitorias, nome)
}

func (s *EsbocoArmazenamentoJogador) ObterLiga() []Jogador {
	return s.liga
}

func (e *EsbocoArmazenamentoJogador) ObterPontuacaoJogador(nome string) int {
	pontuacao := e.pontuacoes[nome]
	return pontuacao
}

func VerificaVitoriaJogador(t *testing.T, armazenamento *EsbocoArmazenamentoJogador, vencedor string) {
	t.Helper()

	if len(armazenamento.registrosVitorias) != 1 {
		t.Fatalf("recebi %d chamadas de GravarVitoria, esperava %d", len(armazenamento.registrosVitorias), 1)
	}

	if armazenamento.registrosVitorias[0] != vencedor {
		t.Errorf("não armazenou o vencedor correto recebi '%s', esperava '%s'", armazenamento.registrosVitorias[0], vencedor)
	}
}
