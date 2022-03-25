package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Jogador struct {
	Nome     string
	Vitorias int
}
type ArmazenamentoJogador interface {
	ObterPontuacaoJogador(nome string) int
	RegistrarVitoria(nome string)
	ObterLiga() []Jogador
}

type EsbocoArmazenamentoJogador struct {
	pontuacoes        map[string]int
	registrosVitorias []string
	liga              []Jogador
}
type ServidorJogador struct {
	armazenamento ArmazenamentoJogador
	http.Handler
}

type SistemaDeArquivoDeArmazenamentoDoJogador struct {
	bancoDeDados io.ReadWriteSeeker
}

func (f *SistemaDeArquivoDeArmazenamentoDoJogador) ObterLiga() []Jogador {
	f.bancoDeDados.Seek(0, 0)
	liga, _ := NovaLiga(f.bancoDeDados)
	return liga
}

func (e *EsbocoArmazenamentoJogador) ObterPontuacaoJogador(nome string) int {
	pontuacao := e.pontuacoes[nome]
	return pontuacao
}

func (e *EsbocoArmazenamentoJogador) RegistrarVitoria(nome string) {
	e.registrosVitorias = append(e.registrosVitorias, nome)
}

func NovoServidorJogador(armazenamento ArmazenamentoJogador) *ServidorJogador {
	s := new(ServidorJogador)

	s.armazenamento = armazenamento

	roteador := http.NewServeMux()
	roteador.Handle("/liga", http.HandlerFunc(s.manipulaLiga))
	roteador.Handle("/jogadores/", http.HandlerFunc(s.manipulaJogadores))

	s.Handler = roteador

	return s
}
func (s *ServidorJogador) manipulaLiga(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(s.armazenamento.ObterLiga())
	w.Header().Set("content-type", "application/json")
}

func (s *ServidorJogador) obterTabelaDaLiga() []Jogador {
	return []Jogador{
		{"Chris", 20},
	}

}

func (s *ServidorJogador) manipulaJogadores(w http.ResponseWriter, r *http.Request) {
	jogador := r.URL.Path[len("/jogadores/"):]

	switch r.Method {
	case http.MethodPost:
		s.registrarVitoria(w, jogador)
	case http.MethodGet:
		s.mostrarPontuacao(w, jogador)
	}
}

func (s *ServidorJogador) mostrarPontuacao(w http.ResponseWriter, jogador string) {
	pontuacao := s.armazenamento.ObterPontuacaoJogador(jogador)

	if pontuacao == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, pontuacao)
}

func (s *ServidorJogador) registrarVitoria(w http.ResponseWriter, jogador string) {
	s.armazenamento.RegistrarVitoria(jogador)
	w.WriteHeader(http.StatusAccepted)
}

func (s *EsbocoArmazenamentoJogador) ObterLiga() []Jogador {
	return s.liga
}

func (a *ArmazenamentoJogadorEmMemoria) ObterLiga() []Jogador {
	var liga []Jogador
	for nome, vitórias := range a.armazenamento {
		liga = append(liga, Jogador{nome, vitórias})
	}
	return liga
}

func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PegaLiga() []Jogador {
	var liga []Jogador
	json.NewDecoder(f.bancoDeDados).Decode(&liga)
	return liga
}

type Reader interface {
	Read(p []byte) (n int, err error)
}

type ReadSeeker interface {
	Reader
	Seeker
}
type Seeker interface {
	Seek(offset int64, whence int) (int64, error)
}

func (f *SistemaDeArquivoDeArmazenamentoDoJogador) PontuacaoJogador(nome string) int {

	jogador := f.PegaLiga().Find(nome)

	if jogador != nil {
		return jogador.Vitorias
	}

	return 0
}

func (f *SistemaDeArquivoDeArmazenamentoDoJogador) SalvaVitoria(nome string) {
	liga := f.PegaLiga()
	jogador := liga.Find(nome)

	if jogador != nil {
		jogador.Wins++
	} else {
		liga = append(liga, Jogador{nome, 1})
	}

	f.bancoDeDados.Seek(0, 0)
	json.NewEncoder(f.bancoDeDados).Encode(liga)
}
