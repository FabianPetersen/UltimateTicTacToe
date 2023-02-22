package main

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
	"testing"
)

func TestSpeedMinimax(t *testing.T) {
	mm := minimax.NewMinimax(game)
	mm.Depth = 10
	move := mm.Search()
	fmt.Println(move)
}
