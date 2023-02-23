package mtd

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
)

const inf float64 = 10000

func mtdF(node *minimax.Node, f float64, d byte) (float64, int) {
	g := f
	lowerBound, upperBound := -inf, inf
	beta := -inf
	bestMove := 0
	for lowerBound < upperBound {
		if g == lowerBound {
			beta = g + 1
		} else {
			beta = g
		}

		g, bestMove = node.Search(beta-1, beta, d, node.State.CurrentPlayer)

		if g < beta {
			upperBound = g
		} else {
			lowerBound = g
		}
	}

	return g, bestMove
}

func IterativeDeepening(node *minimax.Node, maxDepth byte) int {
	// Start the guess at the current heuristic
	var firstGuess float64 = node.State.HeuristicPlayer(node.State.CurrentPlayer)
	bestMove := 0
	var d byte = 0
	// Game.HeuristicStorage.Reset()
	minimax.TranspositionTable.Reset()
	for ; d < maxDepth; d++ {
		firstGuess, bestMove = mtdF(node, firstGuess, d)
	}
	fmt.Printf("Stored nodes, %d Depth %d \n", minimax.TranspositionTable.Count(), maxDepth)
	return bestMove
}
