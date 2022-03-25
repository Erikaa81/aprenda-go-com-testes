package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func NovaLiga(rdr io.Reader) ([]Jogador, error) {
	var liga []Jogador
	err := json.NewDecoder(rdr).Decode(&liga)
	if err != nil {
		err = fmt.Errorf("Problema parseando a liga, %v", err)
	}

	return liga, err
}

type Liga []Jogador

func (l Liga) Find(nome string) *Jogador {
	for i, p := range l {
		if p.Nome == nome {
			return &l[i]
		}
	}
	return nil
}
