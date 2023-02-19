package bns

import (
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
)

/*
function nextGuess(α, β, subtreeCount) is
    return α + (β − α) × (subtreeCount − 1) / subtreeCount


function bns(node, α, β) is
    subtreeCount := number of children of node

    do
        test := nextGuess(α, β, subtreeCount)
        betterCount := 0
        for each child of node do
            bestVal := −alphabeta(child, −test, −(test − 1))
            if bestVal ≥ test then
                betterCount := betterCount + 1
                bestNode := child
        (update number of sub-trees that exceeds separation test value)
        (update alpha-beta range)
    while not (β − α < 2 or betterCount = 1)

    return bestNode
*/

func nextGuess(alpha float64, beta float64, subtreeCount int) float64 {
	count := float64(subtreeCount)
	return alpha + ((beta - alpha) * (count - 1) / count)
}

func BestNodeSearch(node *minimax.Node) int {
	trueCount := node.State.Len()
	subtreeCount := node.State.Len()

	h := node.State.HeuristicPlayer(node.State.CurrentPlayer)
	alpha, beta := -h, h
	var depth byte = minimax.GetDepth(node.State) - 2
	bestMove := 0
	for true {
		test := nextGuess(alpha, beta, subtreeCount)
		betterCount := 0
		for i := 0; i < trueCount; i++ {
			bestVal, _ := minimax.NewNode(node.State, i).Search(-test, -(test - 1), depth, node.State.CurrentPlayer)
			if bestVal >= test {
				betterCount += 1
				bestMove = i
			}
		}

		// All are better
		if betterCount == subtreeCount {
			// Reduce beta-value
			beta = -1
		}

		if betterCount > 1 {
			// Reduce nodes
			alpha = test
			subtreeCount = betterCount
		}

		if beta-alpha < 2 || betterCount == 1 {
			break
		}
	}

	return bestMove
}
