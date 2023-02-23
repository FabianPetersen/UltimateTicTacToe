package main

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
	"github.com/FabianPetersen/UltimateTicTacToe/mtd"
	"testing"
	"time"
)

func TestSpeedMTDF(t *testing.T) {
	start := time.Now()

	move := mtd.IterativeDeepening(&minimax.Node{
		State: game.Copy(),
	}, 10)

	fmt.Println(time.Since(start))
	fmt.Println(move)
}
