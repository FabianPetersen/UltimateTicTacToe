package mtd

import "github.com/FabianPetersen/UltimateTicTacToe/Game"

func first(game *Game.Game, tt *Storage) float64 {
	if lookup, exists := tt.Get(game.Hash()); exists {
		return (lookup.LowerBound + lookup.UpperBound) / 2
	}
	return 0
}

func next(lowerBound float64, upperBound float64, bestValue float64, bound float64) float64 {
	if bestValue < bound {
		return bestValue
	} else {
		return bestValue + 1
	}
}

func (mtd *MTD) MTD_F() int {
	_ = mtd.mtd(mtd.root, first, next, mtd.maxDepth)
	return mtd.aiMove
}
