package bns

import (
	"fmt"
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
	/*
		if alpha <= 0 {
			beta = math.Min(beta, inf/2)
		}

		if beta >= 0 {
			alpha = math.Max(alpha, -inf/2)
		}
	*/

	//count := float64(subtreeCount)
	guess := (alpha + beta) / 2 // * (count - 1) / count // (count-1)/count*math.Abs(beta-math.Abs(alpha))
	/*
		if guess == alpha {
			return guess + 1
		} else if guess == beta {
			return guess - 1
		}
	*/
	return guess
}

const inf float64 = 8000

/*
private ICollection<TMove> BestNodeSearch(TPosition position) // TODO: initial guess from previous iteration!
        {
            var alpha = -int.MaxValue;
            var beta = int.MaxValue;

            IList<TMove> candidates = rules.LegalMovesAt(position).ToList();

            while (alpha + 1 < beta && candidates.Count > 1)
            {
                int guess = NextGuess(alpha, beta, candidates.Count);

                var newCandidates = new List<TMove>();

                foreach (var move in candidates)
                {
                    int value = NullWindowTest(position, move, guess);

                    if (searchTreeManager.IsStopRequested())
                    {
                        break;
                    }
                    if (value >= guess)
                    {
                        newCandidates.Add(move);
                    }
                }

                if (searchTreeManager.IsStopRequested())
                {
                    break;
                }
                if (newCandidates.Count > 0)
                {
                    candidates = newCandidates;
                    alpha = guess;
                }
                else
                {
                    beta = guess;
                }
            }

            return candidates;
        }
*/

func BestNodeSearch(node *minimax.Node, test float64, depth byte) int {
	children := []int{}
	for i := 0; i < node.State.Len(); i++ {
		children = append(children, i)
	}

	alpha, beta := -inf, inf
	bestMove := 0
	for alpha+0.05 < beta && len(children) > 1 {
		worthyChildren := []int{}

		test = nextGuess(alpha, beta, len(children))
		for _, i := range children {
			bestVal, _ := minimax.NewNode(node.State, i).Search(-test, -(test - 1), depth, node.State.CurrentPlayer)
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
	}
	return bestMove
}

func IterativeDeepening(node *minimax.Node, maxDepth byte) int {
	// Start the guess at the current heuristic
	var firstGuess float64 = node.State.HeuristicPlayer(node.State.CurrentPlayer)
	bestMove := 0
	var d byte = 0
	minimax.TranspositionTable.Reset()
	for ; d < maxDepth; d++ {
		bestMove = BestNodeSearch(node, firstGuess, d)
	}
	fmt.Printf("Stored nodes, %d Depth %d \n", minimax.TranspositionTable.Count(), maxDepth)
	return bestMove
}
