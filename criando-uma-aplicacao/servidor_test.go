package poquer

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestGETJogadores(t *testing.T) {
	armazenamento := EsbocoArmazenamentoJogador{
		map[string]int{
			"Maria": 20,
			"Pedro": 10,
		},
		nil,
		nil,
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

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusNotFound)
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
		map[string]int{},
		nil,
		nil,
	}
	servidor := NovoServidorJogador(&armazenamento)

	t.Run("registra vitorias na chamada ao método HTTP POST", func(t *testing.T) {

		jogador := "Maria"

		requisicao := novaRequisicaoRegistrarVitoriaPost(jogador)
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusAccepted)

		if len(armazenamento.registrosVitorias) != 1 {
			t.Fatalf("verifiquei %d chamadas a RegistrarVitoria, esperava %d", len(armazenamento.registrosVitorias[0]), 3)
		}
	})
}

func novaRequisicaoRegistrarVitoriaPost(nome string) *http.Request {
	requisicao, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/jogadores/%s", nome), nil)
	return requisicao
}
func TestGravarVitoriasEBuscarEstasVitorias(t *testing.T) {
	bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, "")
	defer limpaBancoDeDados()
	armazenamento, err := NovoSistemaDeArquivoDeArmazenamentoDoJogador(bancoDeDados)

	defineSemErro(t, err)

	servidor := NovoServidorJogador(armazenamento)
	jogador := "Pepper"

	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoRegistrarVitoriaPost(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoRegistrarVitoriaPost(jogador))
	servidor.ServeHTTP(httptest.NewRecorder(), novaRequisicaoRegistrarVitoriaPost(jogador))

	t.Run("obter pontuação", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, novaRequisicaoObterPontuacao(jogador))
		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)

		verificarCorpoRequisicao(t, resposta.Body.String(), "3")
	})

	t.Run("obter liga", func(t *testing.T) {
		resposta := httptest.NewRecorder()
		servidor.ServeHTTP(resposta, novaRequisicaoObterDeLiga())
		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)

		obtido := obterLigaDaResposta(t, resposta.Body)
		esperado := []Jogador{
			{"Pepper", 3},
		}
		verificaLiga(t, obtido, esperado)
	})
	resposta := httptest.NewRecorder()
	servidor.ServeHTTP(resposta, novaRequisicaoObterPontuacao(jogador))
	verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)

	verificarCorpoRequisicao(t, resposta.Body.String(), "3")

}

func TestLiga(t *testing.T) {

	t.Run("retorna a tabela da Liga como JSON", func(t *testing.T) {
		ligaEsperada := []Jogador{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		armazenamento := EsbocoArmazenamentoJogador{nil, nil, ligaEsperada}
		servidor := NovoServidorJogador(&armazenamento)

		requisicao := novaRequisicaoObterDeLiga()
		resposta := httptest.NewRecorder()

		servidor.ServeHTTP(resposta, requisicao)

		recebido := obterLigaDaResposta(t, resposta.Body)

		verificarRespostaCodigoStatus(t, resposta.Code, http.StatusOK)
		defineLiga(t, recebido, ligaEsperada)
		defineTipoDeConteudo(t, resposta, jsonContentType)

	})

}

func defineTipoDeConteudo(t *testing.T, resposta *httptest.ResponseRecorder, esperado string) {
	t.Helper()
	if resposta.Header().Get("content-type") != esperado {
		t.Errorf("resposta nao tinha content-type de %s, recebido %v", esperado, resposta.HeaderMap)
	}
}

func verificaLiga(t *testing.T, obtido, esperado []Jogador) {
	t.Helper()
	if !reflect.DeepEqual(obtido, esperado) {
		t.Errorf("obtido %v esperado %v", obtido, esperado)
	}
}

func novaRequisicaoObterDeLiga() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/liga", nil)
	return req
}

func TestArmazenamentoDeSistemaDeArquivo(t *testing.T) {

	t.Run("liga ordenada", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
            {"Nome": "Cleo", "Vitorias": 10},
            {"Nome": "Chris", "Vitorias": 33}]`)
		defer limpaBancoDeDados()

		armazenamento, err := NovoSistemaDeArquivoDeArmazenamentoDoJogador(bancoDeDados)

		defineSemErro(t, err)

		recebido := armazenamento.ObterLiga()

		esperado := []Jogador{
			{"Chris", 33},
			{"Cleo", 10},
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

		armazenamento, err := NovoSistemaDeArquivoDeArmazenamentoDoJogador(bancoDeDados)

		defineSemErro(t, err)

		recebido := armazenamento.ObterPontuacaoJogador("Chris")
		esperado := 33
		definePontuacao(t, recebido, esperado)
	})

	t.Run("armazena vitórias de um jogador existente", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}]`)
		defer limpaBancoDeDados()

		armazenamento, err := NovoSistemaDeArquivoDeArmazenamentoDoJogador(bancoDeDados)

		defineSemErro(t, err)

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

		armazenamento, err := NovoSistemaDeArquivoDeArmazenamentoDoJogador(bancoDeDados)

		defineSemErro(t, err)

		armazenamento.RegistrarVitoria("Pepper")

		recebido := armazenamento.ObterPontuacaoJogador("Pepper")
		esperado := 1
		definePontuacao(t, recebido, esperado)
	})
	t.Run("funciona com um arquivo vazio", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, "")
		defer limpaBancoDeDados()

		_, err := NovoSistemaDeArquivoDeArmazenamentoDoJogador(bancoDeDados)

		defineSemErro(t, err)
	})
	t.Run("liga ordernada", func(t *testing.T) {
		bancoDeDados, limpaBancoDeDados := criaArquivoTemporario(t, `[
			{"Nome": "Cleo", "Vitorias": 10},
			{"Nome": "Chris", "Vitorias": 33}]`)
		defer limpaBancoDeDados()

		armazenamento, err := NovoSistemaDeArquivoDeArmazenamentoDoJogador(bancoDeDados)

		defineSemErro(t, err)

		recebido := armazenamento.ObterLiga()

		esperado := []Jogador{
			{"Chris", 33},
			{"Cleo", 10},
		}

		defineLiga(t, recebido, esperado)

		recebido = armazenamento.ObterLiga()
		defineLiga(t, recebido, esperado)
	})

}

func definePontuacao(t *testing.T, recebido, esperado int) {
	t.Helper()
	if !reflect.DeepEqual(recebido, esperado) {
		t.Errorf("recebido %v esperado %v", recebido, esperado)
	}
}

func obterLigaDaResposta(t *testing.T, body io.Reader) []Jogador {
	t.Helper()
	liga, err := NovaLiga(body)

	if err != nil {
		t.Fatalf("Nao foi possivel passar resposta do servidor '%s' em pedaco de jogador, '%v'", body, err)
	}

	return liga
}

func defineLiga(t *testing.T, recebido, esperado []Jogador) {
	t.Helper()
	if !reflect.DeepEqual(recebido, esperado) {
		t.Errorf("recebido %v esperado %v", recebido, esperado)
	}

}
func defineSemErro(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("não esperava um erro mas obteve um, %v", err)
	}
}

func criaArquivoTemporario(t *testing.T, dadoInicial string) (*os.File, func()) {
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
