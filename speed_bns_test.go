package main

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"github.com/FabianPetersen/UltimateTicTacToe/bns"
	"testing"
)

var game = Game.NewGame()

func TestSpeedBNS(t *testing.T) {
	move := bns.IterativeDeepening(game, 10)
	fmt.Println(move)
}
