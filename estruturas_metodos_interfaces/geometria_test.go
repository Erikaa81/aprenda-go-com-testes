package estruturasmetodosinterfaces

import (
	"math"
	"testing"
)

type Forma interface {
	Area() float64
}
type Retangulo struct {
	Largura float64
	Altura  float64
}
type Circulo struct {
	Raio float64
}
type Triangulo struct {
	Base   float64
	Altura float64
}
type Quadrado struct {
	Lado float64
}

func (r Retangulo) Area() float64 {
	return r.Largura * r.Altura
}

func (c Circulo) Area() float64 {
	return math.Pi * c.Raio * c.Raio
}

func (t Triangulo) Area() float64 {
	return (t.Base * t.Altura) * 0.5
}

func (q Quadrado) Area() float64 {
	return q.Lado * q.Lado
}

func TestArea(t *testing.T) {
	testesArea := []struct {
		nome    string
		forma   Forma
		temArea float64
	}{
		{nome: "Retângulo", forma: Retangulo{Largura: 12, Altura: 6}, temArea: 72.0},
		{nome: "Círculo", forma: Circulo{Raio: 10}, temArea: 314.1592653589793},
		{nome: "Triângulo", forma: Triangulo{Base: 12, Altura: 6}, temArea: 36.0},
		{nome: "Quadrado", forma: Quadrado{Lado: 8}, temArea: 64.0},
	}

	for _, tt := range testesArea {
		t.Run(tt.nome, func(t *testing.T) {
			resultado := tt.forma.Area()
			if resultado != tt.temArea {
				t.Errorf("%#v resultado %.2f, esperado %.2f", tt.forma, resultado, tt.temArea)
			}
		})
	}
}
