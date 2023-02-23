package main

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
	"github.com/FabianPetersen/UltimateTicTacToe/mtd"
	"testing"
	"time"
)

func TestSpeedMTDF(t *testing.T) {
	start := time.Now()
	fmt.Println("populate board time", time.Since(start), len(Game.BoardHeuristicCacheP1))

	start = time.Now()

	move := mtd.IterativeDeepening(&minimax.Node{
		State: game.Copy(),
	}, 10)

	fmt.Println(time.Since(start))

	// Formatted string, such as "2h3m0.5s" or "4.503μs"
	fmt.Println(move)

}
