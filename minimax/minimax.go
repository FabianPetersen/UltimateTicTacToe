package minimax

import (
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"math"
)

type Node struct {
	state *Game.Game
}

func NewNode(state *Game.Game, action int) *Node {
	cp := state.Copy()
	cp.ApplyActionModify(action)
	return &Node{
		state: cp,
	}
}

func (n *Node) Search(alpha float64, beta float64, depth int, maxPlayer Game.Player) (float64, int) {
	if depth == 0 || n.state.IsTerminal() {
		return n.state.Heuristic()[maxPlayer], 0
	}

	if n.state.CurrentPlayer == maxPlayer {
		value := math.Inf(-1)
		bestMove := 0
		for i := 0; i < n.state.Len(); i++ {
			searchValue, _ := NewNode(n.state, i).Search(alpha, beta, depth-1, maxPlayer)
			if searchValue >= value {
				value = searchValue
				bestMove = i
			}

			alpha = math.Max(alpha, value)
			if value >= beta {
				break
			}
		}
		return value, bestMove

	} else {
		value := math.Inf(1)
		bestMove := 0
		for i := 0; i < n.state.Len(); i++ {
			searchValue, _ := NewNode(n.state, i).Search(alpha, beta, depth-1, maxPlayer)
			if searchValue <= value {
				value = searchValue
				bestMove = i
			}
			beta = math.Min(beta, value)
			if value <= alpha {
				break
			}
		}
		return value, bestMove
	}
}

type Minimax struct {
	root  *Node
	depth int
}

func NewMinimax(state *Game.Game) *Minimax {
	return &Minimax{root: &Node{
		state: state,
	}, depth: 9}
}

func (minimax *Minimax) Search() int {
	_, bestMove := minimax.root.Search(math.Inf(-1), math.Inf(1), minimax.depth, minimax.root.state.CurrentPlayer)
	return bestMove
}
