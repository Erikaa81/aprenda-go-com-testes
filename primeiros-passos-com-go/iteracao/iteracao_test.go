package iteracao

import (
	"fmt"
	"testing"
)

func BenchmarkRepetir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repetir("a")
	}
}

const quantidadeRepeticoes = 10

func Repetir(caractere string) string {
	var repeticoes string
	for i := 0; i < quantidadeRepeticoes; i++ {
		repeticoes += caractere
	}
	return repeticoes

}

func TestRepetir(t *testing.T) {
	repeticoes := Repetir("a")
	esperado := "aaaaaaaaaa"

	if repeticoes != esperado {
		t.Errorf("esperado '%s' mas obteve '%s'", esperado, repeticoes)
	}
}
func ExampleRepetir() {
	repeticoes := Repetir("b")

	fmt.Println(repeticoes)
	// Output: bbbbbbbbbb
}
