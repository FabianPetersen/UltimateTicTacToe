package main

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/gmcts"
	"math"
	"testing"
	"time"
	"unsafe"
)

func TestSpeed(t *testing.T) {
	start := time.Now()
	game.MakeMove(8, 8)

	var mcts *gmcts.MCTS
	for i := 0; i < 100; i++ {
		mcts = gmcts.NewMCTS(game)
		mcts.SearchRounds(25000)
	}

	fmt.Println(mcts.BestAction())
	fmt.Println(gmcts.NodePoolIndex)
	fmt.Println(time.Since(start))
}

func ln(x float32) float32 {
	var bx = *(*uint32)(unsafe.Pointer(&x))
	var ex = bx >> 23
	var t = float32(ex) - 127
	bx = 1065353216 | (bx & 8388607)
	x = *(*float32)(unsafe.Pointer(&bx))
	return -1.49278 + (2.11263+(-0.729104+0.10969*x)*x)*x + 0.6931471806*t
}

func TestLN(t *testing.T) {
	for i := float32(0); i < 100; i += 0.01 {
		fmt.Println(math.Log(float64(i)), ln(i))
	}
}
