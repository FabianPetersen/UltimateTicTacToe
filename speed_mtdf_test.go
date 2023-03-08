package main

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/mtd"
	"testing"
	"time"
)

func TestSpeedMTDF(t *testing.T) {
	start := time.Now()

	mtd.IterativeDeepeningTime(game, 10, time.Second*10)

	fmt.Println(time.Since(start))
}
