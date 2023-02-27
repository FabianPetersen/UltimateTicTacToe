package Game

import (
	"math/rand"
	"time"
)

/*
 last 3 byte (player 1, player 2, draw) = (31, 30, 29)
 Second 9 byte = (9 - 17)
 First 9 byte = (0 - 8)

	// 0 1 2
    // 7 8 3
    // 6 5 4
*/

type Player byte
type GameHash *[10]uint32

const Player1 Player = 0
const Player2 Player = 1
const boardLength = 9
const GlobalBoard byte = 0xF0
const PlayerBoardIndex = 9

var randSource = rand.New(rand.NewSource(time.Now().Unix()))

/*	corner -> middle -> side */
var moveOrder = []byte{0, 4, 2, 6, 8, 1, 3, 5, 7}

type Game struct {
	Board        [boardLength + 1]uint32
	OverallBoard uint32
}

func (g *Game) GetMoves(executeMove func(byte, byte) bool) {
	currentBoard := byte(g.Board[PlayerBoardIndex] >> 1)
	if currentBoard != GlobalBoard {
		board := (g.Board[currentBoard] | (g.Board[currentBoard] >> 9)) & 0x1FF
		for _, move := range moveOrder {
			if board&(0x1<<move) == 0 {
				if executeMove(currentBoard, move) {
					return
				}
			}
		}
	} else {
		for _, i := range moveOrder {
			// Check if the board is open
			if g.OverallBoard&((0x1<<i)|(0x1<<(i+9))|(0x1<<(i+18))) == 0 {
				board := (g.Board[i] | (g.Board[i] >> 9)) & 0x1FF
				for _, move := range moveOrder {
					if board&(0x1<<move) == 0 {
						if executeMove(i, move) {
							return
						}
					}
				}
			}
		}
	}
}

func (g *Game) Hash() GameHash {
	return &g.Board
}

func (g *Game) Compare(c *Game) bool {
	for i := 0; i < boardLength+1; i++ {
		if g.Board[i] != c.Board[i] {
			return false
		}
	}

	return true
}

//IsTerminal returns true if this game state is a terminal state
func (g *Game) IsTerminal() bool {
	return CheckCompleted(g.OverallBoard&0x1FF) || CheckCompleted((g.OverallBoard>>9)&0x1FF) || ((g.OverallBoard>>18)|(g.OverallBoard>>9)|g.OverallBoard)&0x1FF == 0x1FF
}

func NewGame() *Game {
	PopulateBoards()
	return &Game{
		Board: [boardLength + 1]uint32{
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0x11,
		},
		OverallBoard: 0x0,
	}
}

func (g *Game) Copy() Game {
	return Game{
		Board: [boardLength + 1]uint32{
			g.Board[0],
			g.Board[1],
			g.Board[2],
			g.Board[3],
			g.Board[4],
			g.Board[5],
			g.Board[6],
			g.Board[7],
			g.Board[8],
			g.Board[9],
		},
		OverallBoard: g.OverallBoard,
	}
}

func (g *Game) Rotate(rotateBy uint32) {
	g.Board[rotateBy%8], g.Board[(rotateBy+1)%8], g.Board[(rotateBy+2)%8], g.Board[(rotateBy+3)%8], g.Board[(rotateBy+4)%8], g.Board[(rotateBy+5)%8], g.Board[(rotateBy+6)%8], g.Board[(rotateBy+7)%8] = g.Board[(rotateBy+6)%8], g.Board[(rotateBy+7)%8], g.Board[rotateBy%8], g.Board[(rotateBy+1)%8], g.Board[(rotateBy+2)%8], g.Board[(rotateBy+3)%8], g.Board[(rotateBy+4)%8], g.Board[(rotateBy+5)%8]
	for i := 0; i < boardLength; i++ {
		g.Board[i] = g.Board[i]&0xe0020100 | rotl(g.Board[i], rotateBy) | (rotl(g.Board[i]>>9, rotateBy) << 9)
	}
	g.OverallBoard = g.OverallBoard&0xe4020100 | rotl(g.OverallBoard, rotateBy) | (rotl(g.OverallBoard>>9, rotateBy) << 9) | (rotl(g.OverallBoard>>18, rotateBy) << 18)

	// Change current board
	currentBoard := byte(g.Board[PlayerBoardIndex] >> 1)
	if currentBoard != 8 && currentBoard != GlobalBoard {
		g.Board[PlayerBoardIndex] = (((uint32(currentBoard) + rotateBy) % 8) << 1) | (g.Board[PlayerBoardIndex] & 0x1)
	}
}

func (g *Game) Invert() {
	for i := 0; i < boardLength; i++ {
		g.Board[i] = (g.Board[i]&0x1FF)<<9 | (g.Board[i]>>9)&0x1FF
	}
	g.OverallBoard = g.OverallBoard&0xFFFC0000 | (g.OverallBoard&0x1FF)<<9 | (g.OverallBoard>>9)&0x1FF

	// Change player
	g.Board[PlayerBoardIndex] = (g.Board[PlayerBoardIndex] & 0x1FE) | ((g.Board[PlayerBoardIndex] & 0x1) ^ 0x1)
}

func (g *Game) UnMakeMove(lastPos byte, lastBoard byte, prevBoard byte) {
	// Unset move
	if Player(g.Board[PlayerBoardIndex]&0x1) == Player2 {
		g.Board[lastBoard] &^= 1 << lastPos

	} else {
		g.Board[lastBoard] &^= 1 << (lastPos + 9)
	}

	// Reset win
	g.OverallBoard &^= 0x1<<lastBoard | 0x1<<(lastBoard+9) | 0x1<<(lastBoard+18)
	g.Board[PlayerBoardIndex] = uint32(prevBoard)<<1 | ((g.Board[PlayerBoardIndex] & 0x1) ^ 0x1)
}

func (g *Game) MakeMove(boardIndex byte, pos byte) {
	if Player(g.Board[PlayerBoardIndex]&0x1) == Player1 {
		g.Board[boardIndex] |= 1 << pos

		// Player 1 - WIN
		if CheckCompleted(g.Board[boardIndex] & 0x1FF) {
			g.OverallBoard |= 1 << boardIndex
		}

	} else {
		g.Board[boardIndex] |= 1 << (pos + 9)

		// Player 2 - WIN
		if CheckCompleted((g.Board[boardIndex] >> 9) & 0x1FF) {
			g.OverallBoard |= 1 << (boardIndex + 9)
		}
	}

	// Draw
	if (g.Board[boardIndex]|(g.Board[boardIndex]>>9))&0x1FF == 0x1FF {
		g.OverallBoard |= 1 << (boardIndex + 18)
	}

	// Only switch board is a space is available
	if g.IsBoardFinished(pos) {
		g.Board[PlayerBoardIndex] = uint32(GlobalBoard)<<1 | ((g.Board[PlayerBoardIndex] & 0x1) ^ 0x1)
	} else {
		g.Board[PlayerBoardIndex] = uint32(pos)<<1 | ((g.Board[PlayerBoardIndex] & 0x1) ^ 0x1)
	}
}

func (g *Game) ValidMove(boardIndex byte, pos byte) bool {
	currentBoard := byte(g.Board[PlayerBoardIndex] >> 1)
	return !g.IsBoardFinished(boardIndex) && (currentBoard == GlobalBoard || boardIndex == currentBoard) && g.Board[boardIndex]&(1<<pos) == 0 && g.Board[boardIndex]&(1<<(pos+9)) == 0
}

func (g *Game) IsBoardFinished(pos byte) bool {
	return !((g.OverallBoard>>pos)&0x1 == 0 && (g.OverallBoard>>(pos+9))&0x1 == 0 && (g.OverallBoard>>(pos+18))&0x1 == 0)
}
