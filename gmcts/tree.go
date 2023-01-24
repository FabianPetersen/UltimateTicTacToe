package gmcts

import (
	"context"
	"time"
)

//Search searches the tree for a specified time
//
//Search will panic if the Game's ApplyAction
//method returns an error or if any game state's Hash()
//method returns a noncomparable value.
func (t *Tree) Search(duration time.Duration) {
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
			t.search()
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
		t.search()
	}
}

//search performs 1 round of the MCTS algorithm
func (t *Tree) search() {
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

//MaxDepth returns the maximum depth of this tree.
//The value can be thought of as the amount of moves ahead
//this tree searched through.
func (t Tree) MaxDepth() int {
	maxDepth := 0
	for _, node := range t.gameStates {
		if node.state.turn > maxDepth {
			maxDepth = node.state.turn
		}
	}
	return maxDepth
}

func (t *Tree) bestAction() int {
	root := t.current

	//Select the child with the highest winrate
	var bestAction int
	bestWinRate := -1.0
	player := root.state.Player()
	for i := 0; i < root.actionCount; i++ {
		winRate := root.children[i].nodeScore[player] / root.childVisits[i]
		if winRate > bestWinRate {
			bestAction = i
			bestWinRate = winRate
		}
	}

	return bestAction
}
