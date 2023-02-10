package mtd

import (
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"time"
)

type Storage struct {
	nodeStore map[Game.GameHash]*MTDNode
}

func (storage *Storage) Get(hash Game.GameHash) (*MTDNode, bool) {
	node, exists := storage.nodeStore[hash]
	return node, exists
}

func (storage *Storage) Set(node *MTDNode) {
	storage.nodeStore[node.State.Hash()] = node
}

func (storage *Storage) Reset() {
	storage.nodeStore = map[Game.GameHash]*MTDNode{}
}

func NewStorage() Storage {
	return Storage{
		nodeStore: map[Game.GameHash]*MTDNode{},
	}
}

const inf float64 = 1000000
const eps = 0.001

type Flag byte

const (
	EXACT       Flag = 0
	UPPER_BOUND Flag = 1
	LOWER_BOUND Flag = 2
)

type MTDNode struct {
	State      *Game.Game
	LowerBound float64
	UpperBound float64
	BestMove   int
	Depth      byte
	Flag       Flag
}

type MTD struct {
	root        *Game.Game
	maxDepth    byte
	aiMove      int
	tt          *Storage
	maxDuration time.Duration
}

func NewMTD(game *Game.Game, maxDepth byte) *MTD {
	storage := NewStorage()
	return &MTD{
		root:        game.Copy(),
		maxDepth:    maxDepth,
		aiMove:      0,
		tt:          &storage,
		maxDuration: time.Second * 10,
	}
}

func (mtd *MTD) mt(game *Game.Game, gamma float64, depth byte, origDepth byte, maxPlayer Game.Player) float64 {
	lowerBound, upperBound := -inf, inf
	bestMove := 0

	// Return the old node if exists
	if lookup, found := mtd.tt.Get(game.Hash()); found && lookup.Depth >= depth {
		lowerBound, upperBound = lookup.LowerBound, lookup.UpperBound
		if lowerBound > gamma {
			if depth == origDepth {
				mtd.aiMove = lookup.BestMove
			}
			return lowerBound
		}

		if upperBound < gamma {
			if depth == origDepth {
				mtd.aiMove = lookup.BestMove
			}
			return upperBound
		}
	}

	bestValue := -inf
	if depth == 0 || game.IsTerminal() {
		score := game.HeuristicPlayer(maxPlayer) * (1 + 0.001*float64(depth))
		lowerBound, upperBound, bestValue = score, score, score

	} else {
		bestMove = 0
		possibleMoves := game.Len()
		for i := 0; bestValue < gamma && i < possibleMoves; i++ {
			ngame := game.Copy()
			ngame.ApplyActionModify(i)

			moveValue := -mtd.mt(ngame, -gamma, depth-1, origDepth, maxPlayer)

			if bestValue < moveValue {
				bestValue = moveValue
				bestMove = i
			}
		}

		if bestValue < gamma {
			upperBound = bestValue
		} else {
			if depth == origDepth {
				mtd.aiMove = bestMove
			}
			lowerBound = bestValue
		}
	}

	if depth > 0 && !game.IsTerminal() {
		mtd.tt.Set(&MTDNode{
			State:      game,
			LowerBound: lowerBound,
			UpperBound: upperBound,
			BestMove:   bestMove,
			Depth:      depth,
		})
	}

	return bestValue
}

type Next func(float64, float64, float64, float64) float64
type First func(*Game.Game, *Storage) float64

func (mtd *MTD) mtd(game *Game.Game, first First, next Next, depth byte) float64 {
	start := time.Now()
	f := first(game, mtd.tt)
	bound, bestValue := f, f
	lowerBound, upperBound := -inf, inf
	for time.Now().Sub(start) < mtd.maxDuration && lowerBound < upperBound {
		bound = next(lowerBound, upperBound, bestValue, bound)
		bestValue = mtd.mt(game, bound-eps, depth, depth, game.CurrentPlayer)
		if bestValue < bound {
			upperBound = bestValue
		} else {
			lowerBound = bestValue
		}
	}

	return bestValue
}
