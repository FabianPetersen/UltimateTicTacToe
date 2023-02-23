package Game

const boardCornerRating = 1.4
const boardSideRating = 1
const boardMiddleRating = 1.75
const posCornerRating = 0.2
const posSideRating = 0.17
const posMiddleRating = 0.22

var boardRating = []float64{boardCornerRating, boardSideRating, boardCornerRating, boardSideRating, boardCornerRating, boardSideRating, boardCornerRating, boardSideRating, boardMiddleRating}
var posRating = []float64{posCornerRating, posSideRating, posCornerRating, posSideRating, posCornerRating, posSideRating, posCornerRating, posSideRating, posMiddleRating}
var BoardHeuristicCacheP1 = map[uint32]float64{}
var BoardHeuristicCacheP2 = map[uint32]float64{}

// var HeuristicStorage = NewStorage()

func getOffset(player Player) (int, int) {
	if player == Player1 {
		return 0, 9
	}
	return 9, 0
}

func HeuristicBoard(player Player, board uint32, isOverallBoard bool) float64 {
	if !isOverallBoard {
		if player == Player1 {
			if score, ok := BoardHeuristicCacheP1[board&0x3FFFF]; ok {
				return score
			}
		} else if score, ok := BoardHeuristicCacheP2[board&0x3FFFF]; ok {
			return score
		}
	}

	offset, enemyOffset := getOffset(player)
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

			// The enemy has won a square
		} else {
			// Give a reward for enemy moves
			score -= 10 - float64(bitCount(enemyBoard))*0.25
		}
	}
	return score
}

func (g *Game) HeuristicPlayer(player Player) float64 {
	// Don't rerun calculation unless needed
	/*if player == Player2 {
		if score, ok := HeuristicStorage.Get(g.Hash()); ok {
			return score
		}
	}*/

	var score float64 = 0
	var playerOffset, enemyOffset = getOffset(player)
	if checkCompleted((g.overallBoard >> playerOffset) & 0x1FF) {
		// Incentivise quicker wins
		score = 750 - float64(g.MovesMade())*10
		// HeuristicStorage.Set(g.Hash(), score)
		return score
	}

	if checkCompleted((g.overallBoard >> enemyOffset) & 0x1FF) {
		// Incentivise slower losses
		score -= 750 - float64(g.MovesMade())*5
	}

	for i := 0; i < boardLength; i++ {
		boardScore := HeuristicBoard(player, g.Board[i], false)
		score += boardScore * 1.5 * boardRating[i]
	}

	score += HeuristicBoard(player, g.overallBoard, true) * 150
	/*if player == Player2 {
		HeuristicStorage.Set(g.Hash(), score)
	}*/
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

func PopulateBoards() {
	for i := 0; i < 19683; i++ {
		c := i
		var board uint32 = 0
		for ii := 0; ii < 9; ii++ {
			switch c % 3 {
			case 0:
			case 1:
				board |= 0x1 << ii
			case 2:
				board |= 0x1 << (ii + 9)
			}
			c /= 3
		}

		isValid := false

		// Check if board is valid
		xBoard, oBoard, _ := board&0x1FF, (board>>9)&0x1FF, (board&0x1FF)|((board>>9)&0x1FF)
		xCount, oCount := bitCount(xBoard), bitCount(oBoard)

		// Board can be valid only if either
		// xCount and oCount is same or count
		// is one more than oCount
		if xCount == oCount || xCount+1 == oCount || xCount == oCount+1 {
			xWin, oWin := checkCompleted(xBoard), checkCompleted(oBoard)
			if xWin && !oWin {
				isValid = true
			} else if !xWin && oWin {
				isValid = true
			} else if !xWin && !oWin {
				isValid = true
			}
		}

		if isValid {
			BoardHeuristicCacheP1[board] = HeuristicBoard(Player1, board, false)
			BoardHeuristicCacheP2[board] = HeuristicBoard(Player2, board, false)
		}
	}
}
