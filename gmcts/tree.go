package gmcts

/*
import (
	"context"
	"time"
)

//SearchTime searches the tree for a specified time
//
//SearchTime will panic if the Game's ApplyAction
//method returns an error or if any game state's Hash()
//method returns a noncomparable value.
func (t *Tree) SearchTime(duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	t.SearchContext(ctx)
}

//SearchContext searches the tree using a given context
//
//SearchContext will panic if the Game's ApplyAction
//method returns an error or if any game state's Hash()
//method returns a noncomparable value.
func (t *Tree) SearchContext(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			t.Search()
		}
	}
}

//SearchRounds searches the tree for a specified number of rounds
//
//SearchRounds will panic if the Game's ApplyAction
//method returns an error or if any game state's Hash()
//method returns a noncomparable value.
func (t *Tree) SearchRounds(rounds int) {
	for i := 0; i < rounds; i++ {
		t.Search()
	}
}

//Search performs 1 round of the MCTS algorithm
func (t *Tree) Search() {
	t.current.runSimulation()
}

//Rounds returns the number of MCTS rounds were performed
//on this tree.
func (t Tree) Rounds() int {
	return t.current.nodeVisits
}

//Nodes returns the number of nodes created on this tree.
func (t Tree) Nodes() int {
	return len(t.gameStates)
}

func (t *Tree) bestAction() int {
	root := t.current

	var bestAction int
	//Select the child with the highest winrate
	if t.bestActionPolicy == MAX_CHILD_SCORE {
		bestWinRate := -1.0
		player := root.state.Player()
		for i := 0; i < root.actionCount; i++ {
			winRate := root.children[i].nodeScore[player] / root.childVisits[i]
			if winRate > bestWinRate {
				bestAction = i
				bestWinRate = winRate
			}
		}
	} else if t.bestActionPolicy == ROBUST_CHILD {
		mostVisists := -1.0
		for i := 0; i < root.actionCount; i++ {
			if root.childVisits[i] >= mostVisists {
				bestAction = i
				mostVisists = root.childVisits[i]
			}
		}
	}

	return bestAction
}
*/
