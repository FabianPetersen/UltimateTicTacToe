package main

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/gmcts"
	"testing"
	"time"
)

func TestTiming(t *testing.T) {
	game.MakeMove(8, 8)

	start := time.Now()
	mcts := gmcts.NewMCTS(game)
	mcts.SearchTime(97 * time.Millisecond)
	fmt.Println(time.Since(start))
	fmt.Println(gmcts.NodePoolIndex)
}
