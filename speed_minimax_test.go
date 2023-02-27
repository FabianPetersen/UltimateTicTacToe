package main

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
	"testing"
)

func TestSpeedMinimax(t *testing.T) {
	mm := minimax.NewMinimax()
	mm.Depth = 10
	move, _ := mm.Search(game)
	fmt.Println(move)
}
