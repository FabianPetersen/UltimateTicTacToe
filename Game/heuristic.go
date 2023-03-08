package Game

var BoardHeuristicCacheP1 = map[uint32]float64{}
var BoardHeuristicCacheP2 = map[uint32]float64{}

type HeuristicScores struct {
	BoardRating [9]float64
	PosRating   [9]float64

	BoardCornerRating                        float64
	BoardSideRating                          float64
	BoardMiddleRating                        float64
	PosCornerRating                          float64
	PosSideRating                            float64
	PosMiddleRating                          float64
	OverallWinLossRating                     float64
	OverallAlmostDrawWinLossRating           float64
	GlobalStateRating                        float64
	OverallBoardMultiplierRating             float64
	WinMovesMadeLossRating                   float64
	LossMovesMadeAdvantageRating             float64
	TwoInARowAdvantageRating                 float64
	EnemyTwoInARowLossRating                 float64
	EnemyWonBoardLossRating                  float64
	EnemyWonBoardDiscountRating              float64
	WonBoardRating                           float64
	DrawBoardScoreEnemyDiscountRating        float64
	DrawBoardScorePlayerDiscountRating       float64
	LocalBoardWinPlayedMovesDiscountRating   float64
	OverallBoardWinPlayedMovesDiscountRating float64
}

func DefaultHeuristic() *HeuristicScores {
	h := HeuristicScores{
		BoardCornerRating:                        1.5,
		BoardSideRating:                          0.49,
		BoardMiddleRating:                        2.32,
		PosCornerRating:                          2.03,
		PosSideRating:                            0.61,
		PosMiddleRating:                          1.21,
		OverallWinLossRating:                     5000,
		OverallAlmostDrawWinLossRating:           2500,
		GlobalStateRating:                        70.06,
		OverallBoardMultiplierRating:             150,
		WinMovesMadeLossRating:                   10,
		LossMovesMadeAdvantageRating:             12,
		TwoInARowAdvantageRating:                 9,
		EnemyTwoInARowLossRating:                 9,
		EnemyWonBoardLossRating:                  24,
		EnemyWonBoardDiscountRating:              0.25,
		WonBoardRating:                           24,
		DrawBoardScoreEnemyDiscountRating:        0.12,
		DrawBoardScorePlayerDiscountRating:       0.12,
		LocalBoardWinPlayedMovesDiscountRating:   0.36,
		OverallBoardWinPlayedMovesDiscountRating: 1.32,
	}

	h.BoardRating = [9]float64{h.BoardCornerRating, h.BoardSideRating, h.BoardCornerRating, h.BoardSideRating, h.BoardCornerRating, h.BoardSideRating, h.BoardCornerRating, h.BoardSideRating, h.BoardMiddleRating}
	h.PosRating = [9]float64{h.PosCornerRating, h.PosSideRating, h.PosCornerRating, h.PosSideRating, h.PosCornerRating, h.PosSideRating, h.PosCornerRating, h.PosSideRating, h.PosMiddleRating}

	return &h
}

func getOffset(player Player) (int, int) {
	if player == Player1 {
		return 0, 9
	}
	return 9, 0
}

func (g *Game) HeuristicBoard(player Player, board uint32, isOverallBoard bool) float64 {
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
	jointBoard := (board>>18)&0x1FF | playerBoard | enemyBoard

	// Check if won
	if CheckCompleted(playerBoard) {
		// Give a discount on the amount of moves made in the board
		// To incentivise a lower number of total moves
		if !isOverallBoard {
			score += g.HeuristicScores.WonBoardRating - float64(bitCount(playerBoard))*g.HeuristicScores.LocalBoardWinPlayedMovesDiscountRating
		} else {
			score += g.HeuristicScores.WonBoardRating - float64(bitCount(playerBoard))*g.HeuristicScores.OverallBoardWinPlayedMovesDiscountRating
		}

		// The board is a draw
	} else if jointBoard == 0x1FF && !CheckCompleted(enemyBoard) {
		// Give a reward for the amount of wasted enemy moves (or won moves
		if !isOverallBoard {
			score += float64(bitCount(enemyBoard)) * g.HeuristicScores.DrawBoardScoreEnemyDiscountRating
		} else {
			score += float64(bitCount(playerBoard)) * g.HeuristicScores.DrawBoardScorePlayerDiscountRating
		}

	} else {
		// Calculate pos for items
		for i := 0; i < boardLength; i++ {
			if playerBoard&(0x1<<i) > 0 {
				score += g.HeuristicScores.PosRating[i]
			}
		}

		// If the board is not won by enemy
		if !CheckCompleted(enemyBoard) {
			// Check 2 joint items
			if checkCloseWinningSequence(playerBoard, jointBoard) > 0 {
				score += g.HeuristicScores.TwoInARowAdvantageRating
			}

			// Check 2 joint items
			if checkCloseWinningSequence(enemyBoard, jointBoard) > 0 {
				score -= g.HeuristicScores.EnemyTwoInARowLossRating
			}

			// The enemy has won a square
		} else {
			// Give a reward for enemy moves
			score -= g.HeuristicScores.EnemyWonBoardLossRating - float64(bitCount(enemyBoard))*g.HeuristicScores.EnemyWonBoardDiscountRating
		}
	}
	return score
}

func (g *Game) HeuristicPlayer(player Player) float64 {
	// Don't rerun calculation unless needed
	var score float64 = 0
	var playerOffset, enemyOffset = getOffset(player)

	// Calculate pos rating
	playerBoard := (g.OverallBoard >> playerOffset) & 0x1FF
	enemyBoard := (g.OverallBoard >> enemyOffset) & 0x1FF
	jointBoard := (g.OverallBoard>>18)&0x1FF | playerBoard | enemyBoard

	if CheckCompleted(playerBoard) {
		// Incentivise quicker wins
		return g.HeuristicScores.OverallWinLossRating - float64(g.MovesMade())*g.HeuristicScores.WinMovesMadeLossRating
	} else if CheckCompleted(enemyBoard) {
		// Incentivise slower losses
		return -g.HeuristicScores.OverallWinLossRating + float64(g.MovesMade())*g.HeuristicScores.LossMovesMadeAdvantageRating
	} else if jointBoard == 0x1FF {
		if bitCount(playerBoard) > bitCount(enemyBoard) {
			return g.HeuristicScores.OverallAlmostDrawWinLossRating - float64(g.MovesMade())*g.HeuristicScores.WinMovesMadeLossRating
		}
		return -g.HeuristicScores.OverallAlmostDrawWinLossRating + float64(g.MovesMade())*g.HeuristicScores.LossMovesMadeAdvantageRating
	}

	if byte(g.Board[PlayerBoardIndex]>>1) == GlobalBoard {
		if byte(g.Board[PlayerBoardIndex]&0x1) == byte(player) {
			score += g.HeuristicScores.GlobalStateRating

		} else {
			score -= g.HeuristicScores.GlobalStateRating
		}
	}

	for i := 0; i < boardLength; i++ {
		boardScore := g.HeuristicBoard(player, g.Board[i], false)
		score += boardScore * g.HeuristicScores.BoardRating[i]
	}

	score += g.HeuristicBoard(player, g.OverallBoard, true) * g.HeuristicScores.OverallBoardMultiplierRating
	return score
}

func (g *Game) MovesMade() uint32 {
	// Count how far the game has progressed
	var movesPlayed uint32 = 0
	for i := 0; i < 9; i++ {
		movesPlayed += bitCount(g.Board[i] & 0x3FFFF)
	}
	return movesPlayed
}

func popBoardHelper(b uint32, pos int, player int) uint32 {
	if player == 0 {
		b |= 0x1 << pos
	} else if player == 1 {
		b |= 0x1 << (pos + 9)
	}
	return b
}

func (g *Game) PopulateBoards() {
	var board uint32 = 0
	for i0 := 0; i0 < 3; i0++ {
		for i1 := 0; i1 < 3; i1++ {
			for i2 := 0; i2 < 3; i2++ {
				for i3 := 0; i3 < 3; i3++ {
					for i4 := 0; i4 < 3; i4++ {
						for i5 := 0; i5 < 3; i5++ {
							for i6 := 0; i6 < 3; i6++ {
								for i7 := 0; i7 < 3; i7++ {
									for i8 := 0; i8 < 3; i8++ {
										board = 0
										board = popBoardHelper(board, 0, i0)
										board = popBoardHelper(board, 1, i1)
										board = popBoardHelper(board, 2, i2)
										board = popBoardHelper(board, 3, i3)
										board = popBoardHelper(board, 4, i4)
										board = popBoardHelper(board, 5, i5)
										board = popBoardHelper(board, 6, i6)
										board = popBoardHelper(board, 7, i7)
										board = popBoardHelper(board, 8, i8)

										// Check if board is valid
										isValid := false
										xBoard, oBoard, jointBoard := board&0x1FF, (board>>9)&0x1FF, uint16((board&0x1FF)|((board>>9)&0x1FF))

										// Board can be valid only if either
										xWin, oWin := CheckCompleted(xBoard), CheckCompleted(oBoard)
										if xWin && !oWin {
											isValid = true
										} else if !xWin && oWin {
											isValid = true
										} else if !xWin && !oWin {
											isValid = true
										}

										if isValid {
											BoardHeuristicCacheP1[board] = g.HeuristicBoard(Player1, board, false)
											BoardHeuristicCacheP2[board] = g.HeuristicBoard(Player2, board, false)
										}

										MovesStorage[jointBoard] = []byte{}
										for _, move := range moveOrder {
											if jointBoard&(0x1<<move) == 0 {
												MovesStorage[jointBoard] = append(MovesStorage[jointBoard], move)
											}
										}
										MovesLengthStorage[jointBoard] = byte(len(MovesStorage[jointBoard]))
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
