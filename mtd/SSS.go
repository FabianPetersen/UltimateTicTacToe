package mtd

import "github.com/FabianPetersen/UltimateTicTacToe/Game"

func (mtd *MTD) SSS() int {
	_ = mtd.mtd(mtd.root, func(game *Game.Game, storage *Storage) float64 {
		return 160 // self.win_score  # essence of SSS algorithm
	}, func(lowerbound float64, upperbound float64, bestValue float64, bound float64) float64 {
		return bestValue
	}, mtd.maxDepth)

	return mtd.aiMove
}
