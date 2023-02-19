package main

import (
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
	"github.com/FabianPetersen/UltimateTicTacToe/mtd"
	"testing"
)

func TestSpeedMTDF(t *testing.T) {
	game := Game.NewGame()
	mtd.IterativeDeepening(&minimax.Node{
		State: game.Copy(),
	}, 10)
}
