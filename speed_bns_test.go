package main

import (
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"github.com/FabianPetersen/UltimateTicTacToe/bns"
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
	"testing"
)

func TestSpeedBNS(t *testing.T) {
	game := Game.NewGame()
	bns.IterativeDeepening(&minimax.Node{
		State: game.Copy(),
	}, 10)
}
