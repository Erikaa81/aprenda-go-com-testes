package poquer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
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

type ServidorJogador struct {
	armazenamento ArmazenamentoJogador
	http.Handler
}
type SistemaDeArquivoDeArmazenamentoDoJogador struct {
	bancoDeDados *json.Encoder
	liga         Liga
}

type Reader interface {
	Read(p []byte) (n int, err error)
}

const jsonContentType = "application/json"

func NovoSistemaDeArquivoDeArmazenamentoDoJogador(arquivo *os.File) (*SistemaDeArquivoDeArmazenamentoDoJogador, error) {

	err := iniciaArquivoBDDeJogador(arquivo)

	if err != nil {
		return nil, fmt.Errorf("problema inicializando arquivo do jogador, %v", err)
	}

	liga, err := NovaLiga(arquivo)

	if err != nil {
		return nil, fmt.Errorf("problema carregando armazenamento de jogador do arquivo %s, %v", arquivo.Name(), err)
	}

	return &SistemaDeArquivoDeArmazenamentoDoJogador{
		bancoDeDados: json.NewEncoder(&fita{arquivo}),
		liga:         liga,
	}, nil
}

func ArmazenamentoSistemaDeArquivoJogadorAPartirDeArquivo(path string) (*SistemaDeArquivoDeArmazenamentoDoJogador, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("falha ao abrir %s %v", path, err)
	}

	closeFunc := func() {
		db.Close()
	}

	armazenamento, err := NovoSistemaDeArquivoDeArmazenamentoDoJogador(db)

	if err != nil {
		return nil, nil, fmt.Errorf("falha ao criar sistema de arquivos para armazenar jogadores, %v ", err)
	}

	return armazenamento, closeFunc, nil
}

func iniciaArquivoBDDeJogador(arquivo *os.File) error {
	arquivo.Seek(0, 0)

	info, err := arquivo.Stat()

	if err != nil {
		return fmt.Errorf("problema ao usar arquivo %s, %v", arquivo.Name(), err)
	}

	if info.Size() == 0 {
		arquivo.Write([]byte("[]"))
		arquivo.Seek(0, 0)
	}

	return nil
}

func (f *SistemaDeArquivoDeArmazenamentoDoJogador) ObterPontuacaoJogador(nome string) int {

	jogador := f.liga.Find(nome)

	if jogador != nil {
		return jogador.Vitorias
	}

	return 0
}

func (f *SistemaDeArquivoDeArmazenamentoDoJogador) RegistrarVitoria(nome string) {
	jogador := f.liga.Find(nome)

	if jogador != nil {
		jogador.Vitorias++
	} else {
		f.liga = append(f.liga, Jogador{nome, 1})
	}

	f.bancoDeDados.Encode(f.liga)
}

func (f *SistemaDeArquivoDeArmazenamentoDoJogador) ObterLiga() []Jogador {
	sort.Slice(f.liga, func(i, j int) bool {
		return f.liga[i].Vitorias > f.liga[j].Vitorias
	})
	return f.liga
}

func NovoServidorJogador(armazenamento ArmazenamentoJogador) *ServidorJogador {
	s := new(ServidorJogador)

	s.armazenamento = armazenamento

	roteador := http.NewServeMux()
	roteador.Handle("/liga", http.HandlerFunc(s.ManipulaLiga))
	roteador.Handle("/jogadores/", http.HandlerFunc(s.ManipulaJogadores))

	s.Handler = roteador

	return s
}
func (s *ServidorJogador) ManipulaLiga(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(s.armazenamento.ObterLiga())
}

func (s *ServidorJogador) ManipulaJogadores(w http.ResponseWriter, r *http.Request) {
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
