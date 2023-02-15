package main

import (
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
	"testing"
)

func TestSpeedMinimax(t *testing.T) {
	game := Game.NewGame()
	mm := minimax.NewMinimax(game)
	mm.Depth = 10
	mm.Search()
}
