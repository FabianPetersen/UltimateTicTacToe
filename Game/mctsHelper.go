package Game

func (g *Game) GreedyMove() int {
	length := g.Len()
	var i = 0
	for i = 0; i < length; i++ {
		move, _ := g.GetMove(i)
		if g.CurrentPlayer == Player1 && checkCloseWinningSequenceMove(byte(move), g.Board[i]&0x1FF, (g.Board[g.CurrentBoard]|(g.Board[g.CurrentBoard]>>9))&0x1FF) {
			return i
		} else if g.CurrentPlayer == Player2 && checkCloseWinningSequenceMove(byte(move), (g.Board[i]>>9)&0x1FF, (g.Board[g.CurrentBoard]|(g.Board[g.CurrentBoard]>>9))&0x1FF) {
			return i
		}
	}
	return randSource.Intn(length)
}
