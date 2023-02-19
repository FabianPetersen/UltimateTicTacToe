package Game

const boardCornerRating = 1.4
const boardSideRating = 1
const boardMiddleRating = 1.75
const posCornerRating = 0.2
const posSideRating = 0.17
const posMiddleRating = 0.22

var boardRating = []float64{boardCornerRating, boardSideRating, boardCornerRating, boardSideRating, boardCornerRating, boardSideRating, boardCornerRating, boardSideRating, boardMiddleRating}
var posRating = []float64{posCornerRating, posSideRating, posCornerRating, posSideRating, posCornerRating, posSideRating, posCornerRating, posSideRating, posMiddleRating}
var boardHeuristicCacheP1 = map[uint32]float64{}
var boardHeuristicCacheP2 = map[uint32]float64{}

func (g *Game) getOffset(player Player) (int, int) {
	if player == Player1 {
		return 0, 9
	}
	return 9, 0
}

func (g *Game) HeuristicBoard(player Player, board uint32, isOverallBoard bool) float64 {
	if player == Player1 {
		if score, ok := boardHeuristicCacheP1[board]; ok {
			return score
		}
	} else if score, ok := boardHeuristicCacheP2[board]; ok {
		return score
	}

	offset, enemyOffset := g.getOffset(player)
	var score float64 = 0

	// Calculate pos rating
	playerBoard := (board >> offset) & 0x1FF
	enemyBoard := (board >> enemyOffset) & 0x1FF
	jointBoard := (board>>18)&0x1FF | (board>>offset)&0x1FF | (board>>enemyOffset)&0x1FF

	// Check if won
	if checkCompleted(playerBoard) {
		// Give a discount on the amount of moves made in the board
		// To incentivise a lower number of total moves
		if !isOverallBoard {
			score += 24 - float64(bitCount(playerBoard))*0.25
		} else {
			score += 24
		}

		// The board is a draw
	} else if jointBoard == 0x1FF && !checkCompleted(enemyBoard) {
		// Give a reward for the amount of wasted enemy moves (or won moves
		if !isOverallBoard {
			score += float64(bitCount(enemyBoard)) * 0.12
		} else {
			score += float64(bitCount(playerBoard)) * 0.12
		}

	} else {
		// Calculate pos for items
		for i := 0; i < boardLength; i++ {
			if playerBoard&(0x1<<i) > 0 {
				score += posRating[i]
			}
		}

		// If the board is not won by enemy
		if !checkCompleted(enemyBoard) {
			// Check 2 joint items
			if checkCloseWinningSequence(playerBoard, jointBoard) > 0 {
				score += 9
			}

			// Check 2 joint items
			if checkCloseWinningSequence(enemyBoard, jointBoard) > 0 {
				score -= 5
			}

			// Block move score
			//if checkBlockSequence(playerBoard, enemyBoard) > 0 {
			//	score += 6
			//}
			// The enemy has won a square
		} else {
			// Give a reward for enemy moves
			score -= 10 - float64(bitCount(enemyBoard))*0.25
		}
	}

	if player == Player1 {
		boardHeuristicCacheP1[board] = score
	} else {
		boardHeuristicCacheP2[board] = score
	}
	return score
}

func (g *Game) HeuristicPlayer(player Player) float64 {
	// Don't rerun calculation unless needed
	if score, ok := g.heuristicStore[player]; ok {
		return score
	}

	var score float64 = 0
	var playerOffset, enemyOffset = g.getOffset(player)

	if checkCompleted((g.overallBoard >> playerOffset) & 0x1FF) {
		// Incentivise quicker wins
		score = 5000 - float64(g.MovesMade())*10
		g.heuristicStore[player] = score
		return score

	}

	if checkCompleted((g.overallBoard >> enemyOffset) & 0x1FF) {
		// Incentivise slower losses
		score -= 5000 - float64(g.MovesMade())*5
	}

	for i := 0; i < boardLength; i++ {
		boardScore := g.HeuristicBoard(player, g.Board[i], false)
		score += boardScore * 1.5 * boardRating[i]
		/*if i == int(g.CurrentBoard) {
			score += boardScore * boardRating[i]
		}*/
	}

	score += g.HeuristicBoard(player, g.overallBoard, true) * 150

	g.heuristicStore[player] = score
	return score
}

func (g *Game) Heuristic() map[Player]float64 {
	// var maxScore float64 = 221 // 200 + 9*5 + 10 + 4*3 + 9*4 + 9*2
	// var minScore float64 = -maxScore
	// Min-max normalization (usually called feature scaling)
	// (x - xmin) / (xmax - xmin)

	p1 := g.HeuristicPlayer(Player1)
	p2 := g.HeuristicPlayer(Player2)
	return map[Player]float64{
		Player1: p1,
		Player2: p2,
		//		Player1: (p1 - 0.7*p2 - minScore) / (maxScore - minScore),
		//		Player2: (p2 - 0.7*p1 - minScore) / (maxScore - minScore),
	}
}

func (g *Game) MovesMade() uint32 {
	// Count how far the game has progressed
	var movesPlayed uint32 = 0
	for i := 0; i < 9; i++ {
		movesPlayed += bitCount(g.Board[i] & 0x3FFFF)
	}
	return movesPlayed
}
