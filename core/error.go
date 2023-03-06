package core

import (
	"math/rand"
	"time"
)

func Pkmn() string {
	rand.Seed(time.Now().Unix())
	pkmns := []string{
		"pikachu",
	}
	index := rand.Intn(len(pkmns))
	return pkmns[index]

}