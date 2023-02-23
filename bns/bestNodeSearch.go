package bns

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
)

func nextGuess(alpha float64, beta float64, subtreeCount int) float64 {
	return (alpha + beta) / 2
	/*
		if alpha <= 0 {
			beta = math.Min(beta, inf/2)
		}

		if beta >= 0 {
			alpha = math.Max(alpha, -inf/2)
		}

		count := float64(subtreeCount)
		guess := (alpha + beta) / 2 * (count - 1) / count // (count-1)/count*math.Abs(beta-math.Abs(alpha))
		if guess == alpha {
			return guess + 1
		} else if guess == beta {
			return guess - 1
		}
		return guess
	*/
}

const inf float64 = 2000

func BestNodeSearch(node *minimax.Node, test float64, depth byte) (int, float64) {
	children := []int{}
	for i := 0; i < node.State.Len(); i++ {
		children = append(children, i)
	}

	alpha, beta := -inf, inf
	bestMove := 0
	for alpha+0.25 < beta && len(children) > 1 {
		worthyChildren := []int{}

		for _, i := range children {
			bestVal, _ := minimax.NewNode(&node.State, i).Search(-test, -(test - 1), depth, node.State.CurrentPlayer)
			if bestVal >= test {
				bestMove = i
				worthyChildren = append(worthyChildren, i)
			}
		}

		if len(worthyChildren) > 1 {
			// Reduce nodes
			alpha = test
			children = worthyChildren

			// All are better
		} else {
			beta = test
		}
		test = nextGuess(alpha, beta, len(children))
	}
	return bestMove, test
}

func IterativeDeepening(node *minimax.Node, maxDepth byte) int {
	// Start the guess at the current heuristic
	var firstGuess float64 = node.State.HeuristicPlayer(node.State.CurrentPlayer)
	bestMove := 0
	var d byte = 0
	minimax.TranspositionTable.Reset()
	for ; d < maxDepth; d++ {
		bestMove, firstGuess = BestNodeSearch(node, firstGuess, d)
	}
	fmt.Printf("Stored nodes, %d Depth %d \n", minimax.TranspositionTable.Count(), maxDepth)
	return bestMove
}
