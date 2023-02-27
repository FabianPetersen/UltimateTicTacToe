package minimax

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"math"
)

var TranspositionTable = NewStorage()

type Flag byte

const (
	EXACT       Flag = 0
	UPPER_BOUND Flag = 1
	LOWER_BOUND Flag = 2
)

type Node struct {
	lowerBound float64
	upperBound float64
	bestMove   byte
	bestBoard  byte
	depth      byte
	flag       Flag
}

func NewNode(state *Game.Game) (*Node, bool) {
	// Rotate and invert board to check if it already exists in cache
	var oldNode *Node = nil
	var exists bool = false
	var cacheExists bool = false
	for i := 0; i < 2; i++ {
		if !cacheExists && !exists {
			for r := 0; r < 4; r++ {
				// Check if the board exists in the cache
				if !cacheExists && !exists {
					if oldNode, exists = TranspositionTable.Get(state.Hash()); exists {
						cacheExists = true
					}
				}
				state.Rotate(2)
			}
		}
		state.Invert()
	}

	return oldNode, cacheExists
}

const inf float64 = 10000

func Search(state *Game.Game, alpha float64, beta float64, depth byte, maxPlayer Game.Player) (float64, byte, byte) {
	// Restore the values from the last node
	n, cached := NewNode(state)
	if cached && n.depth >= depth {
		if n.flag == EXACT {
			return n.lowerBound, n.bestMove, n.bestBoard
		} else if n.flag == LOWER_BOUND {
			alpha = math.Max(alpha, n.lowerBound)
		} else if n.flag == UPPER_BOUND {
			beta = math.Min(beta, n.upperBound)
		}

		if alpha >= beta {
			return n.lowerBound, n.bestMove, n.bestBoard
		}
	}

	var value float64 = 0
	var currentBestMove byte = 0
	var currentBestBoard byte = 0
	var prevBoard = byte(state.Board[Game.PlayerBoardIndex] >> 1)
	if depth == 0 || state.IsTerminal() {
		return state.HeuristicPlayer(maxPlayer), 0, 0

		// This is a max node
	} else if Game.Player(state.Board[Game.PlayerBoardIndex]&0x1) == maxPlayer {
		value = -inf
		a := alpha
		state.GetMoves(func(boardIndex byte, move byte) bool {
			state.MakeMove(boardIndex, move)
			searchValue, _, _ := Search(state, a, beta, depth-1, maxPlayer)
			state.UnMakeMove(move, boardIndex, prevBoard)

			if searchValue >= value {
				value = searchValue
				currentBestMove = move
				currentBestBoard = boardIndex
			}

			a = math.Max(a, value)
			return value >= beta
		})
	} else {
		value = inf
		b := beta
		state.GetMoves(func(boardIndex byte, move byte) bool {
			state.MakeMove(boardIndex, move)
			searchValue, _, _ := Search(state, alpha, b, depth-1, maxPlayer)
			state.UnMakeMove(move, boardIndex, prevBoard)
			if searchValue <= value {
				value = searchValue
				currentBestMove = move
				currentBestBoard = boardIndex
			}
			b = math.Min(b, value)
			return value <= alpha
		})
	}

	if !cached {
		n = &Node{
			lowerBound: -inf,
			upperBound: inf,
			bestBoard:  251,
			bestMove:   251,
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
		n.bestBoard = currentBestBoard
		n.flag = EXACT
	}
	/* Fail high result implies a lower bound */
	if value >= beta {
		n.lowerBound = value
		n.bestMove = currentBestMove
		n.bestBoard = currentBestBoard
		n.flag = LOWER_BOUND
	}
	n.depth = depth
	if !cached {
		TranspositionTable.Set(state.Hash(), n)
	}

	return value, n.bestMove, n.bestBoard
}

type Minimax struct {
	root  *Node
	Depth byte
}

func NewMinimax() *Minimax {
	return &Minimax{root: &Node{}, Depth: 0}
}

func (minimax *Minimax) Search(state *Game.Game) (byte, byte) {
	if minimax.Depth == 0 {
		minimax.setDepth(state)
	}

	TranspositionTable.Reset()
	// Game.HeuristicStorage.Reset()
	_, bestMove, bestBoard := Search(state, math.Inf(-1), math.Inf(1), minimax.Depth, Game.Player(state.Board[Game.PlayerBoardIndex]&0x1))
	fmt.Printf("Stored nodes, %d Depth %d \n", len(TranspositionTable.nodeStore), minimax.Depth)
	return bestMove, bestBoard
}

func (minimax *Minimax) setDepth(state *Game.Game) {
	minimax.Depth = GetDepth(state)
}

func GetDepth(g *Game.Game) byte {
	movesPlayed := g.MovesMade()
	if movesPlayed < 4 {
		return 8
	} else if movesPlayed < 16 {
		return 9
	} else if movesPlayed < 22 {
		return 10
	} else if movesPlayed < 32 {
		return 11
	} else if movesPlayed < 34 {
		return 12
	} else if movesPlayed < 38 {
		return 14
	} else {
		return 15
	}
}
