package mtd

import (
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
	"time"
)

const inf float64 = 100000

func mtdF(state *Game.Game, start *time.Time, maxDuration *time.Duration, f float64, d byte, maxPlayer Game.Player) (float64, byte, byte) {
	g := f
	lowerBound, upperBound := -inf, inf
	beta := -inf
	var bestMove byte = 253
	var bestBoard byte = 253
	var nBestMove, nBestBoard = byte(0), byte(0)
	for lowerBound < upperBound && time.Since(*start) < *maxDuration {
		if g == lowerBound {
			beta = g + 1
		} else {
			beta = g
		}

		g, nBestMove, nBestBoard = minimax.Search(state, beta-1, beta, d, maxPlayer, start, maxDuration)
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

func IterativeDeepeningTime(state *Game.Game, maxDepth byte, maxTime time.Duration) (byte, byte) {
	// Start the guess at the current heuristic
	var maxPlayer = Game.Player(state.Board[Game.PlayerBoardIndex] & 0x1)
	var firstGuess = state.HeuristicPlayer(maxPlayer)

	var bestMove byte = 255
	var bestBoard byte = 255
	var d byte = 0
	// Game.HeuristicStorage.Reset()
	// minimax.TranspositionTable.Reset()
	start := time.Now()
	for ; time.Since(start) < maxTime && d < maxDepth; d++ {
		firstGuess, bestMove, bestBoard = mtdF(state, &start, &maxTime, firstGuess, d, maxPlayer)
	}
	// fmt.Fprintf(os.Stderr, "Stored nodes, %d Depth %d \n", minimax.TranspositionTable.Count(), d)
	return bestMove, bestBoard
}
