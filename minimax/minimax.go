package minimax

import (
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"math"
)

type Storage struct {
	nodeStore map[Game.GameHash]*Node
}

func (storage *Storage) Get(hash Game.GameHash) (*Node, bool) {
	node, exists := storage.nodeStore[hash]
	return node, exists
}

func (storage *Storage) Set(node *Node) {
	storage.nodeStore[node.State.Hash()] = node
}

func (storage *Storage) Reset() {
	storage.nodeStore = map[Game.GameHash]*Node{}
}

func NewStorage() Storage {
	return Storage{
		nodeStore: map[Game.GameHash]*Node{},
	}
}

var TranspositionTable = NewStorage()

type Flag byte

const (
	EXACT       = 0
	UPPER_BOUND = 1
	LOWER_BOUND = 2
)

type Node struct {
	State      *Game.Game
	lowerBound float64
	upperBound float64
	bestMove   int
	depth      byte
	cached     bool
	flag       Flag
}

func NewNode(state *Game.Game, action int) *Node {
	cp := state.Copy()
	cp.ApplyActionModify(int(action))

	// Return the old node if exists
	if oldNode, exists := TranspositionTable.Get(cp.Hash()); exists {
		oldNode.cached = true
		return oldNode
	}

	return &Node{
		State:      cp,
		lowerBound: math.Inf(-1),
		upperBound: math.Inf(1),
		cached:     false,
	}
}

func (n *Node) Search(alpha float64, beta float64, depth byte, maxPlayer Game.Player) (float64, int) {
	// Restore the values from the last node
	if n.cached && n.depth >= depth {
		if n.flag == EXACT {
			return n.lowerBound, n.bestMove
		} else if n.flag == LOWER_BOUND {
			alpha = math.Max(alpha, n.lowerBound)
		} else if n.flag == UPPER_BOUND {
			beta = math.Min(beta, n.upperBound)
		}

		if alpha >= beta {
			return n.lowerBound, n.bestMove
		}
	}

	var value float64 = 0
	var currentBestMove int = 0
	if depth == 0 || n.State.IsTerminal() {
		value = n.State.HeuristicPlayer(maxPlayer)

		// This is a max node
	} else if n.State.CurrentPlayer == maxPlayer {
		value = math.Inf(-1)
		length := n.State.Len()
		a := alpha
		for i := 0; value < beta && i < length; i++ {
			searchValue, _ := NewNode(n.State, i).Search(a, beta, depth-1, maxPlayer)
			if searchValue >= value {
				value = searchValue
				currentBestMove = i
			}

			a = math.Max(a, value)
		}

	} else {
		value = math.Inf(1)
		length := n.State.Len()
		b := beta
		for i := 0; value > alpha && i < length; i++ {
			searchValue, _ := NewNode(n.State, i).Search(alpha, b, depth-1, maxPlayer)
			if searchValue <= value {
				value = searchValue
				currentBestMove = i
			}
			b = math.Min(b, value)
		}
	}

	/* Traditional transposition table storing of bounds */
	/* Fail low result implies an upper bound */
	if value <= alpha {
		n.upperBound = value
		n.flag = UPPER_BOUND
	}
	/* Found an exact minimax value – will not occur if called with zero window */
	if value > alpha && value < beta {
		n.lowerBound = value
		n.upperBound = value
		n.bestMove = currentBestMove
		n.flag = EXACT
	}
	/* Fail high result implies a lower bound */
	if value >= beta {
		n.lowerBound = value
		n.bestMove = currentBestMove
		n.flag = LOWER_BOUND
	}
	n.depth = depth
	TranspositionTable.Set(n)
	return value, n.bestMove
}

type Minimax struct {
	root  *Node
	depth byte
}

func NewMinimax(state *Game.Game) *Minimax {
	return &Minimax{root: &Node{
		State: state,
	}, depth: 9}
}

func (minimax *Minimax) Search() int {
	TranspositionTable.Reset()
	_, bestMove := minimax.root.Search(math.Inf(-1), math.Inf(1), minimax.depth, minimax.root.State.CurrentPlayer)
	return bestMove
}
