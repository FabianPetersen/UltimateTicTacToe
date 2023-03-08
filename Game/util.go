package Game

func CheckCompleted(test uint32) bool {
	if test&0x100 > 0 && ((test&0x111) == 0x111 || (test&0x144) == 0x144 || (test&0x188) == 0x188 || (test&0x122) == 0x122) {
		return true
	}

	return (test&0x7) == 0x7 || (test&0x70) == 0x70 || (test&0xc1) == 0xc1 || (test&0x1c) == 0x1c
}

func checkCloseWinningSequence(player uint32, board uint32) uint32 {
	var i byte = 0
	for ; i < boardLength; i++ {
		if checkCloseWinningSequenceMove(moveOrder[i], player, board) {
			return 1
		}
	}
	return 0
}

func checkCloseWinningSequenceMove(i byte, player uint32, board uint32) bool {
	// Move already occupied
	if board&(0x1<<i) > 0 {
		return false
	}

	player |= 0x1 << i
	return CheckCompleted(player)
}

func bitCount(u uint32) uint32 {
	uCount := uint32(0)
	uCount = u - ((u >> 1) & 033333333333) - ((u >> 2) & 011111111111)
	return ((uCount + (uCount >> 3)) & 030707070707) % 63
}

func rotl(x uint32, by uint32) uint32 {
	x &= 0xff
	return (x<<by | x>>(8-by)) & 0xFF
}
