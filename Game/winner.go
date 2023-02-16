package Game

func (g *Game) OverallWinner() []Player {
	if checkCompleted(g.overallBoard & 0x1FF) {
		return []Player{Player1}
	} else if checkCompleted((g.overallBoard >> 9) & 0x1FF) {
		return []Player{Player2}
	} else if ((g.overallBoard>>18)|(g.overallBoard>>9)|g.overallBoard)&0x1FF == 0x1FF {
		return []Player{Player1, Player2}
	}
	return []Player{}
}

func (g *Game) Winner(boardIndex byte) []Player {
	// Player 1
	if g.Board[boardIndex]&0x80000000 > 0 || checkCompleted(g.Board[boardIndex]&0x1FF) {
		g.Board[boardIndex] |= 0x80000000
		g.overallBoard |= 0x1 << boardIndex
		return []Player{Player1}

		// Player 2
	} else if g.Board[boardIndex]&0x40000000 > 0 || checkCompleted((g.Board[boardIndex]>>9)&0x1FF) {
		g.Board[boardIndex] |= 0x40000000
		g.overallBoard |= 1 << (boardIndex + 9)
		return []Player{Player2}

		// Draw
	} else if g.Board[boardIndex]&0x20000000 > 0 || (g.Board[boardIndex]|(g.Board[boardIndex]>>9))&0x1FF == 0x1FF {
		g.Board[boardIndex] |= 0x20000000
		g.overallBoard |= 1 << (boardIndex + 18)
		return []Player{Player1, Player2}
	}

	return []Player{}
}
