package gmcts

import (
	"context"
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"time"
)

const bestActionPolicy = ROBUST_CHILD

// MCTS contains functionality for the MCTS algorithm
type MCTS struct {
	game *Game.Game
	root *Node
}

// NewMCTS returns a new MCTS wrapper
func NewMCTS(initial *Game.Game) *MCTS {
	return &MCTS{
		game: initial,
		root: &Node{
			nodeVisits: 1,
		},
	}
}

func (m *MCTS) search() {
	// Selection
	node := m.root
	g := m.game.Copy()
	var player Game.Player
	for node.children != nil {
		// Check children (tree policy)
		player = Game.Player(g.Board[Game.PlayerBoardIndex] & 0x1)
		node = node.treePolicy(&player)
		g.MakeMove(node.board, node.move)
	}

	// Expansion
	if !g.IsTerminal() {
		i := 0
		node.children = make([]*Node, g.Len())
		g.GetMoves(func(board byte, move byte) bool {
			node.children[i] = &Node{
				parent:     node,
				move:       move,
				board:      board,
				nodeVisits: 1,
			}
			i++
			return false
		})

		player = Game.Player(g.Board[Game.PlayerBoardIndex] & 0x1)
		node = node.treePolicy(&player)
		g.MakeMove(node.board, node.move)
	}

	// Simulation
	g.MakeMoveRandUntilTerminal()

	// Backpropagation
	winner := g.WinningPlayer()
	for node.parent != nil {
		if winner < 2 {
			node.nodeScore[winner] += 1
		} else {
			node.nodeScore[Game.Player1] += 0.5
			node.nodeScore[Game.Player2] += 0.5
		}
		node.nodeVisits += 1
		node = node.parent
	}
}

func (t *MCTS) bestAction() byte {
	var bestAction byte
	//Select the child with the highest winrate
	if bestActionPolicy == MAX_CHILD_SCORE {
		bestWinRate := -1.0
		player := Game.Player(t.game.Board[Game.PlayerBoardIndex] & 0x1)
		for i := byte(0); i < t.game.Len(); i++ {
			winRate := t.root.children[i].nodeScore[player] / float64(t.root.children[i].nodeVisits)
			if winRate > bestWinRate {
				bestAction = i
				bestWinRate = winRate
			}
		}
	} else if bestActionPolicy == ROBUST_CHILD {
		mostVisists := -1.0
		for i := byte(0); i < t.game.Len(); i++ {
			if float64(t.root.children[i].nodeVisits) >= mostVisists {
				bestAction = i
				mostVisists = float64(t.root.children[i].nodeVisits)
			}
		}
	}

	return bestAction
}

// SearchTime searches the tree for a specified time
func (t *MCTS) SearchTime(duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	t.SearchContext(ctx)
}

// SearchContext searches the tree using a given context
func (t *MCTS) SearchContext(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			t.search()
		}
	}
}

// SearchRounds searches the tree for a specified number of rounds
//
// SearchRounds will panic if the Game's ApplyAction
// method returns an error or if any game state's Hash()
// method returns a noncomparable value.
func (t *MCTS) SearchRounds(rounds int) {
	for i := 0; i < rounds; i++ {
		t.search()
	}
}
