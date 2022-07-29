package main

import (
	"math/rand"
)

type monster struct {
	faceTo dir
	pos    point
	deck   []*card
	random *rand.Rand
}

func newMonster(random *rand.Rand) *monster {
	return &monster{
		faceTo: left,
		pos:    point{14, 9},
		deck:   newDeck(),
		random: random,
	}
}

type card struct {
	step  int
	kills int
}

func newDeck() []*card {
	return []*card{
		{5, 99},
		{7, 99},
		{7, 99},
		{8, 99},
		{8, 99},
		{10, 99},
		{20, 1},
		{20, 2},
	}
}
