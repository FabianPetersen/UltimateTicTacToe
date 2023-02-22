package main

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
	"testing"
)

func TestSpeedMinimaxIterative(t *testing.T) {
	mm := minimax.NewMinimax(game)
	mm.Depth = 10
	move := mm.SearchIterative()
	fmt.Println(move)
}
