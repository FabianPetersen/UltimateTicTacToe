package Game

func checkCompleted(test uint32) bool {
	return (test&0x111) == 0x111 || (test&0x144) == 0x144 || (test&0x7) == 0x7 || (test&0x188) == 0x188 || (test&0x70) == 0x70 || (test&0xc1) == 0xc1 || (test&0x122) == 0x122 || (test&0x1c) == 0x1c
}

func checkBlockSequence(player uint32, enemy uint32) uint32 {
	var i byte = 0
	for ; i < boardLength; i++ {
		if checkBlockSequenceMove(i, player, enemy) {
			return 1
		}
	}

	return 0
}

func checkCloseWinningSequence(player uint32, board uint32) uint32 {
	var i byte = 0
	for ; i < boardLength; i++ {
		if checkCloseWinningSequenceMove(i, player, board) {
			return 1
		}
	}
	return 0
}

func checkBlockSequenceMove(i byte, player uint32, enemy uint32) bool {
	// IF enemy makes the move, it is a win
	enemy |= 0x1 << i
	if checkCompleted(enemy) {
		// Move already occupied
		if player&(0x1<<i) > 0 {
			return true
		}
	}

	return false
}

func checkCloseWinningSequenceMove(i byte, player uint32, board uint32) bool {
	// Move already occupied
	if board&(0x1<<i) > 0 {
		return false
	}

	player |= 0x1 << i
	return checkCompleted(player)
}

func bitCount(u uint32) uint32 {
	uCount := uint32(0)
	uCount = u - ((u >> 1) & 033333333333) - ((u >> 2) & 011111111111)
	return ((uCount + (uCount >> 3)) & 030707070707) % 63
}
