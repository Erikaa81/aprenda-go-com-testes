package main

type ArmazenamentoJogadorEmMemoria struct {
	armazenamento map[string]int
}

func (s *ArmazenamentoJogadorEmMemoria) RegistrarVitoria(nome string) {
	s.armazenamento[nome]++
}

func (s *ArmazenamentoJogadorEmMemoria) ObterPontuacaoJogador(nome string) int {
	return s.armazenamento[nome]
}
