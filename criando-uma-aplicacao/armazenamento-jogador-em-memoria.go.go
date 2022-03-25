package main

func NovoArmazenamentoJogadorEmMemoria() *ArmazenamentoJogadorEmMemoria {
	return &ArmazenamentoJogadorEmMemoria{map[string]int{}}
}

type ArmazenamentoJogadorEmMemoria struct {
	armazenamento map[string]int
}

func (s *ArmazenamentoJogadorEmMemoria) RegistrarVitoria(nome string) {
	s.armazenamento[nome]++
}

func (s *ArmazenamentoJogadorEmMemoria) ObterPontuacaoJogador(nome string) int {
	return s.armazenamento[nome]
}
