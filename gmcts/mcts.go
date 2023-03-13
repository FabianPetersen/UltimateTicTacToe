package gmcts

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"os"
	"time"
)

const bestActionPolicy = ROBUST_CHILD

var nodePool = [700000]Node{}
var NodePoolIndex = 1

// MCTS contains functionality for the MCTS algorithm
type MCTS struct {
	game     *Game.Game
	gameCopy *Game.Game
	root     *Node
}

var ggCopy Game.Game

// NewMCTS returns a new MCTS wrapper
func NewMCTS(initial *Game.Game) *MCTS {
	NodePoolIndex = 1

	nodePool[0].parent = nil
	nodePool[0].nodeVisits = 1
	nodePool[0].nodeScore = 0
	nodePool[0].childrenCount = 0
	return &MCTS{
		game:     initial,
		gameCopy: &ggCopy,
		root:     &nodePool[0],
	}
}

var player Game.Player
var winningPlayer Game.Player
var node *Node = nil
var availableMoves byte = 0

func (m *MCTS) search() {
	// Selection
	node = m.root
	m.gameCopy.OverallBoard = m.game.OverallBoard
	player = Game.Player(m.gameCopy.Board[Game.PlayerBoardIndex] & 0x1)
	for i := 0; i < 10; i++ {
		m.gameCopy.Board[i] = m.game.Board[i]
	}

	for node.childrenCount > 0 {
		// Check children (tree policy)
		node = node.treePolicy()
		m.gameCopy.MakeMove(node.board, node.move)
	}

	// Expansion
	if !m.gameCopy.IsTerminal() {
		// Fill out the slice to make room for new items
		availableMoves = m.gameCopy.Len()
		if node.maxChildren < availableMoves {
			node.children = append(node.children, make([]*Node, availableMoves-node.maxChildren)...)
			node.maxChildren = availableMoves
		}

		// Iterate over all children
		node.childrenCount = 0
		m.gameCopy.GetMoves(func(board byte, move byte) bool {
			NodePoolIndex++
			node.children[node.childrenCount] = &nodePool[NodePoolIndex]
			node.children[node.childrenCount].parent = node
			node.children[node.childrenCount].move = move
			node.children[node.childrenCount].board = board
			node.children[node.childrenCount].nodeVisits = 1
			node.children[node.childrenCount].nodeScore = 0
			node.children[node.childrenCount].childrenCount = 0
			node.childrenCount++
			return false
		})

		node = node.treePolicy()
		m.gameCopy.MakeMove(node.board, node.move)
	}

	// Simulation
	m.gameCopy.MakeMoveRandUntilTerminal()

	// Backpropagation
	player = m.gameCopy.WinningPlayer()
	for node.parent != nil {
		if player == winningPlayer {
			node.nodeScore += 2
		} else if winningPlayer == 2 {
			node.nodeScore += 1
		}
		node.nodeVisits += 1
		node.nodeExploit = float32(node.nodeScore>>1) / float32(node.nodeVisits)
		node = node.parent
	}
}

func (t *MCTS) BestAction() (byte, byte) {
	var bestAction byte
	var bestBoard byte
	//Select the child with the highest winrate
	if bestActionPolicy == MAX_CHILD_SCORE {
		var bestWinRate float32 = 0
		player = Game.Player(t.game.Board[Game.PlayerBoardIndex] & 0x1)
		for i := byte(0); i < t.game.Len(); i++ {
			winRate := float32(t.root.children[i].nodeScore>>1) / float32(t.root.children[i].nodeVisits)
			if winRate >= bestWinRate {
				bestAction = t.root.children[i].move
				bestBoard = t.root.children[i].board
				bestWinRate = winRate
			}
		}
	} else if bestActionPolicy == ROBUST_CHILD {
		var mostVisists uint16 = 1
		for i := byte(0); i < t.game.Len(); i++ {
			if t.root.children[i].nodeVisits >= mostVisists {
				bestAction = t.root.children[i].move
				bestBoard = t.root.children[i].board
				mostVisists = t.root.children[i].nodeVisits
			}
		}
	}

	return bestAction, bestBoard
}

// SearchTime searches the tree for a specified time
func (t *MCTS) SearchTime(duration time.Duration) {
	var i int
	end := time.Now().Add(duration)
	for {
		if i&0x3F == 0 { // Check in every 128th iteration
			if time.Now().After(end) {
				fmt.Fprintf(os.Stderr, "Rounds %d\n", i)
				break
			}
		}
		t.search()
		i++
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
