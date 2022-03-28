package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestObterJogadores(t *testing.T) {
	armazenamento := EsbocoArmazenamentoJogador{
		pontuacoes: map[string]int{"Maria": 20,
			"Pedro": 10},
		registrosVitorias: []string{},
	}

	servidor := NovoServidorJogador(&armazenamento)

	t.Run("retorna pontuacao de Maria", func(t *testing.T) {
		requisicao := novaRequisicaoObterPontuacao("Maria")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
		verificarCorpoRequisicao(t, resposta.Body.String(), "20")
	})

	t.Run("retorna pontuacao de Pedro", func(t *testing.T) {
		requisicao := novaRequisicaoObterPontuacao("Pedro")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
		verificarCorpoRequisicao(t, resposta.Body.String(), "10")
	})

	t.Run("retorna 404 para jogador não encontrado", func(t *testing.T) {
		requisicao := novaRequisicaoObterPontuacao("Jorge")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		recebido := resposta.Code
		esperado := http.StatusNotFound

		if recebido != esperado {
			t.Errorf("recebido status %d esperado %d", recebido, esperado)
		}
	})
}

func novaRequisicaoObterPontuacao(nome string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/jogadores/%s", nome), nil)
	return req
}

func verificarCorpoRequisicao(t *testing.T, recebido, esperado string) {
	t.Helper()
	if recebido != esperado {
		t.Errorf("corpo da requisição é inválido, recebido '%s' esperado '%s'", recebido, esperado)
	}
}

func verificarRespostaCodigoStatus(t *testing.T, recebido, esperado int) {
	t.Helper()
	if recebido != esperado {
		t.Errorf("não recebeu código de status HTTP esperado, recebido %d, esperado %d", recebido, esperado)
	}
}
func TestArmazenamentoVitorias(t *testing.T) {
	armazenamento := EsbocoArmazenamentoJogador{
		pontuacoes:        map[string]int{},
		registrosVitorias: []string{},
		liga:              []Jogador{},
	}
	servidor := NovoServidorJogador(&armazenamento)

	t.Run("registra vitorias na chamada ao método HTTP POST", func(t *testing.T) {
		requisicao := novaRequisicaoRegistrarVitoriaPost("Maria")
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusAccepted)

		if len(armazenamento.registrosVitorias) != 1 {
			t.Errorf("verifiquei %d chamadas a RegistrarVitoria, esperava %d", len(armazenamento.registrosVitorias), 3)
		}
	})
}

func novaRequisicaoRegistrarVitoriaPost(nome string) *http.Request {
	requisicao, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/jogadores/%s", nome), nil)
	return requisicao
}

func TestLiga(t *testing.T) {
	armazenamento := EsbocoArmazenamentoJogador{}
	servidor := NovoServidorJogador(&armazenamento)

	t.Run("retorna 200 em /liga", func(t *testing.T) {
		requisicao, _ := http.NewRequest(http.MethodGet, "/liga", nil)
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		var obtido []Jogador

		err := json.NewDecoder(resposta.Body).Decode(&obtido)

		if err != nil {
			t.Fatalf("Não foi possível fazer parse da resposta do servidor '%s' no slice de Jogador, '%v'", resposta.Body, err)
		}

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
	})

	t.Run("retorna a tabela da Liga como JSON", func(t *testing.T) {
		ligaEsperada := []Jogador{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		armazenamento := EsbocoArmazenamentoJogador{nil, nil, ligaEsperada}
		servidor := NovoServidorJogador(&armazenamento)

		requisicao := novaRequisicaoDeLiga()
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		obtido := obterLigaDaResposta(t, resposta.Body)
		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
		verificaLiga(t, obtido, ligaEsperada)
	})

}
func obterLigaDaResposta(t *testing.T, body io.Reader) (liga []Jogador) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&liga)

	if err != nil {
		t.Fatalf("Não foi possível fazer parse da resposta do servidor '%s' no slice de Jogador, '%v'", body, err)
	}

	return
}

func verificaLiga(t *testing.T, obtido, esperado []Jogador) {
	t.Helper()
	if !reflect.DeepEqual(obtido, esperado) {
		t.Errorf("obtido %v esperado %v", obtido, esperado)
	}
}

func novaRequisicaoDeLiga() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/liga", nil)
	return req
}

func TestaArmazenamentoDeSistemaDeArquivo(t *testing.T) {

	t.Run("liga de um leitor", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
            {"Nome": "Cleo", "Vitorias": 10},
            {"Nome": "Chris", "Vitorias": 33}]`)
		defer limpaBancoDeDados()

		armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{
			bancoDeDados: bancoDeDados,
			liga:         []Jogador{},
		}

		recebido := armazenamento.ObterLiga()

		esperado := []Jogador{
			{"Cleo", 10},
			{"Chris", 33},
		}

		defineLiga(t, recebido, esperado)

		recebido = armazenamento.ObterLiga()
		defineLiga(t, recebido, esperado)
	})

	t.Run("retorna pontuação do jogador", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
            {"Nome": "Cleo", "Vitorias": 10},
            {"Nome": "Chris", "Vitorias": 33}]`)
		defer limpaBancoDeDados()

		armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{
			bancoDeDados: bancoDeDados,
			liga:         []Jogador{},
		}

		recebido := armazenamento.ObterPontuacaoJogador("Chris")
		esperado := 33
		definePontuacao(t, recebido, esperado)
	})
	t.Run("armazena vitórias de um jogador existente", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}]`)
		defer limpaBancoDeDados()

		armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{
			bancoDeDados: bancoDeDados,
			liga:         []Jogador{},
		}

		armazenamento.RegistrarVitoria("Chris")

		recebido := armazenamento.ObterPontuacaoJogador("Chris")
		esperado := 34
		definePontuacao(t, recebido, esperado)
	})

	t.Run("armazena vitorias de novos jogadores", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}]`)

		defer limpaBancoDeDados()

		armazenamento := SistemaDeArquivoDeArmazenamentoDoJogador{
			bancoDeDados: bancoDeDados,
			liga:         []Jogador{},
		}

		armazenamento.RegistrarVitoria("Pepper")

		recebido := armazenamento.ObterPontuacaoJogador("Pepper")
		esperado := 1
		definePontuacao(t, recebido, esperado)
	})

}

func definePontuacao(t *testing.T, recebido, esperado int) {
	t.Helper()
	if !reflect.DeepEqual(recebido, esperado) {
		t.Errorf("recebido %v esperado %v", recebido, esperado)
	}
}

func defineLiga(t *testing.T, recebido, esperado []Jogador) {
	t.Helper()
	if !reflect.DeepEqual(recebido, esperado) {
		t.Errorf("recebido %v esperado %v", recebido, esperado)
	}

}

func criaArquivoTemporario(t *testing.T, dadoInicial string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	arquivotmp, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("não foi possivel escrever o arquivo temporário %v", err)
	}

	arquivotmp.Write([]byte(dadoInicial))

	removeArquivo := func() {
		arquivotmp.Close()
		os.Remove(arquivotmp.Name())
	}

	return arquivotmp, removeArquivo
}
