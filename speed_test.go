package main

import (
	"github.com/FabianPetersen/UltimateTicTacToe/gmcts"
	"testing"
)

func TestSpeed(t *testing.T) {
	game.MakeMove(8, 8)
	mcts := gmcts.NewMCTS(game)
	mcts.SearchRounds(25000)
}
