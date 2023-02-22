package main

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"github.com/FabianPetersen/UltimateTicTacToe/bns"
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
	"testing"
)

var game = Game.NewGame()

func TestSpeedBNS(t *testing.T) {
	move := bns.IterativeDeepening(&minimax.Node{
		State: game.Copy(),
	}, 10)
	fmt.Println(move)
}
