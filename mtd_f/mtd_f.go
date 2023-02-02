package mtd_f

import (
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"github.com/FabianPetersen/UltimateTicTacToe/minimax"
	"math"
	"time"
)

/*
// from https://askeplaat.wordpress.com/534-2/mtdf-algorithm/
function MTDF(root : node_type; f : integer; d : integer) : integer;
      g := f;
      upperbound := +INFINITY;
      lowerbound := -INFINITY;
      repeat
            if g == lowerbound then beta := g + 1 else beta := g;
            g := AlphaBetaWithMemory(root, beta - 1, beta, d);
            if g < beta then upperbound := g else lowerbound := g;
      until lowerbound >= upperbound;
      return g;

function iterative_deepening(root : node_type) : integer;

      firstguess := 0;
      for d = 1 to MAX_SEARCH_DEPTH do
            firstguess := MTDF(root, firstguess, d);
            if times_up() then break;
      return firstguess;
*/

type MTD_F struct {
	root     *minimax.Node
	maxDepth byte
}

func NewMTD_F(state *Game.Game, maxDepth byte) *MTD_F {
	return &MTD_F{
		root: &minimax.Node{
			State: state,
		},
		maxDepth: maxDepth,
	}
}

func (mtd_f *MTD_F) run(f float64, d byte) (float64, int) {
	bestMove := 0
	value := f
	beta := f
	upperBound := math.Inf(1)
	lowerBound := math.Inf(-1)
	for lowerBound < upperBound {
		if value == lowerBound {
			beta = value + 1
		} else {
			beta = value
		}
		value, bestMove = mtd_f.root.Search(beta-1, beta, d, mtd_f.root.State.CurrentPlayer)
		if value < beta {
			upperBound = value
		} else {
			lowerBound = value
		}
	}
	return value, bestMove
}

func (mtd_f *MTD_F) IterativeDeepening(maxDuration time.Duration) int {
	var guess float64 = 0
	var move int = 0
	minimax.TranspositionTable.Reset()
	start := time.Now()
	var d byte = 1
	for ; start.Sub(time.Now()) > maxDuration && d < mtd_f.maxDepth; d++ {
		guess, move = mtd_f.run(guess, d)
	}
	return move
}
