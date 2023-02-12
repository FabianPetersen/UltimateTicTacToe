package Game

func (g *Game) HeuristicPlayer(player Player) float64 {
	// Don't rerun calculation unless needed
	if score, ok := g.heuristicStore[player]; ok {
		return score
	}

	offset := 0
	if player == Player2 {
		offset = 9
	}

	// Count number of small boards (+5)
	score := bitCount((g.overallBoard>>offset)&0x1FF) * 5

	// Win in center board (+10)
	if ((g.overallBoard >> (offset + 4)) & 0x1) > 0 {
		score += 10
	}

	// Corner boards (+3)
	score += bitCount((g.overallBoard>>offset)&0x145) * 3

	// Center square in any small board (+3)
	/*
		for i := 0; i < boardLength; i++ {
			if ((g.Board[i] >> (offset + 4)) & 0x1) > 0 {
				score += 3
			}
		}
	*/

	// (overall) Sequence of two winning board which can be continued for a winning sequence (+4)
	// overallBoard := ((g.overallBoard >> 18) | (g.overallBoard >> 9) | g.overallBoard) & 0x1FF
	// score += checkCloseWinningSequence((g.overallBoard>>offset)&0x1FF, overallBoard) * 4

	// Sequence of two winning board which can be continued for a winning sequence (+2)
	for i := 0; i < boardLength; i++ {
		board := ((g.Board[i] >> 9) | g.Board[i]) & 0x1FF
		score += checkCloseWinningSequence((g.Board[i]>>offset)&0x1FF, board) * 1
	}

	// Overall win/loss
	if checkCompleted((g.overallBoard >> offset) & 0x1FF) {
		score += 20
	}

	g.heuristicStore[player] = float64(score)
	return float64(score)
}

func (g *Game) Heuristic() map[Player]float64 {
	// var maxScore float64 = 221 // 200 + 9*5 + 10 + 4*3 + 9*4 + 9*2
	// var minScore float64 = -maxScore
	// Min-max normalization (usually called feature scaling)
	// (x - xmin) / (xmax - xmin)

	p1 := g.HeuristicPlayer(Player1)
	p2 := g.HeuristicPlayer(Player2)
	return map[Player]float64{
		Player1: p1 - 0.5*p2,
		Player2: p2 - 0.5*p1,
		//		Player1: (p1 - 0.7*p2 - minScore) / (maxScore - minScore),
		//		Player2: (p2 - 0.7*p1 - minScore) / (maxScore - minScore),
	}
}
