package Game

func checkCompleted(test uint32) bool {
	return (test&0x7) == 0x7 || (test&0x38) == 0x38 || (test&0x1C0) == 0x1C0 || (test&0x49) == 0x49 || (test&0x92) == 0x92 || (test&0x124) == 0x124 || (test&0x111) == 0x111 || (test&0x54) == 0x54
}

func checkCloseWinningSequence(player uint32, board uint32) uint32 {
	var i byte = 0
	var count uint32 = 0
	for ; i < boardLength; i++ {
		if checkCloseWinningSequenceMove(i, player, board) {
			count += 1
		}
	}
	return count
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
