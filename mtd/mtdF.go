package mtd

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
	"time"
)

const inf float64 = 5000

func mtdF(state *Game.Game, f float64, d byte, maxPlayer Game.Player) (float64, byte, byte) {
	g := f
	lowerBound, upperBound := -inf, inf
	beta := -inf
	var bestMove byte = 253
	var bestBoard byte = 253
	var nBestMove, nBestBoard = byte(0), byte(0)
	for lowerBound < upperBound {
		if g == lowerBound {
			beta = g + 1
		} else {
			beta = g
		}

		g, nBestMove, nBestBoard = minimax.Search(state, beta-1, beta, d, maxPlayer)
		if nBestBoard < 200 && nBestMove < 200 {
			bestMove = nBestMove
			bestBoard = nBestBoard
		}

		if g < beta {
			upperBound = g
		} else {
			lowerBound = g
		}
	}

	return g, bestMove, bestBoard
}

func IterativeDeepening(state *Game.Game, maxDepth byte) (byte, byte) {
	// Start the guess at the current heuristic
	var maxPlayer = Game.Player(state.Board[Game.PlayerBoardIndex] & 0x1)
	var firstGuess float64 = state.HeuristicPlayer(maxPlayer)
	var bestMove byte = 0
	var bestBoard byte = 0
	var d byte = 0
	// Game.HeuristicStorage.Reset()
	minimax.TranspositionTable.Reset()
	for ; d < maxDepth; d++ {
		firstGuess, bestMove, bestBoard = mtdF(state, firstGuess, d, maxPlayer)
	}
	fmt.Printf("Stored nodes, %d Depth %d \n", minimax.TranspositionTable.Count(), maxDepth)
	return bestMove, bestBoard
}

func IterativeDeepeningTime(state *Game.Game, maxTime time.Duration) (byte, byte) {
	// Start the guess at the current heuristic
	var maxPlayer = Game.Player(state.Board[Game.PlayerBoardIndex] & 0x1)
	var firstGuess = state.HeuristicPlayer(maxPlayer)

	var bestMove byte = 255
	var bestBoard byte = 255
	var maxDepth byte = 25
	var d byte = 0
	// Game.HeuristicStorage.Reset()
	minimax.TranspositionTable.Reset()
	start := time.Now()
	for ; time.Since(start) < maxTime && d < maxDepth; d++ {
		firstGuess, bestMove, bestBoard = mtdF(state, firstGuess, d, maxPlayer)
	}
	fmt.Printf("Stored nodes, %d Depth %d \n", minimax.TranspositionTable.Count(), d)
	return bestMove, bestBoard
}
