package Game

func (g *Game) getAllAvailableMoves() []byte {
	var moves []byte
	board := (g.Board[g.CurrentBoard] | (g.Board[g.CurrentBoard] >> 9)) & 0x1FF
	for _, move := range moveOrder {
		if board&(0x1<<move) == 0 {
			moves = append(moves, move)
		}
	}
	return moves
}

func (g *Game) getFilteredAvailableMoves() []byte {
	moves := g.getAllAvailableMoves()
	/*
		filtered := g.filterSuddenDeathMoves(moves)
		filterSafeGreedyMove := g.filterGreedyMove(filtered)
		filterUnsafeGreedyMove := g.filterGreedyMove(moves)

		if len(filterSafeGreedyMove) > 0 {
			return filterSafeGreedyMove
		} else if len(filterUnsafeGreedyMove) > 0 {
			return filterUnsafeGreedyMove
		} else if len(filtered) > 0 {
			return filtered
		}
	*/
	return moves
}

func (g *Game) filterGreedyMove(moves []byte) []byte {
	offset, _ := g.getOffset(g.CurrentPlayer)
	closeWinning := []byte{}
	for _, move := range moves {
		if checkCloseWinningSequenceMove(move, (g.Board[g.CurrentBoard]>>offset)&0x1FF, (g.Board[g.CurrentBoard]|(g.Board[g.CurrentBoard]>>9))&0x1FF) {
			closeWinning = append(closeWinning, move)
		}
	}
	return closeWinning
}

func (g *Game) filterSuddenDeathMoves(moves []byte) []byte {
	_, enemyOffset := g.getOffset(g.CurrentPlayer)

	// Check if enemy can win next board if we choose a move
	var acceptableMoves []byte
	for _, move := range moves {
		// check if the enemy has any
		board := ((g.Board[move] >> 9) | g.Board[move]) & 0x1FF
		board |= 0x1 << move
		if g.IsBoardFinished(int(move)) || checkCloseWinningSequence((g.Board[move]>>enemyOffset)&0x1FF, board) == 0 {
			acceptableMoves = append(acceptableMoves, move)
		}
	}
	return acceptableMoves
}
