package main

import (
	"errors"
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/gmcts"
)

/*
 last 3 byte (player 1, player 2, draw) = (31, 30, 29)
 Second 9 byte = (9 - 17)
 First 9 byte = (0 - 8)

	// 0 1 2
	// 3 4 5
	// 6 7 8

	// 0-1-2 = 0x7
	// 3-4-5 = 0x38
	// 6-7-8 = 0x1C0

	// 0-3-6 = 0x49
	// 1-4-7 = 0x92
	// 2-5-8 = 0x124

	// 0-4-8 = 0x111
	// 2-4-6 = 0x54
*/
const player1 gmcts.Player = 1
const player2 gmcts.Player = 2
const boardLength = 9

type Game struct {
	gmcts.Game
	turn         bool
	currentBoard byte
	board        [boardLength]uint32
	overallBoard uint32
}

func (g *Game) Len() int {
	return int(boardLength - bitCount((g.board[g.currentBoard]|(g.board[g.currentBoard]>>9))&0x1FF))
}

func (g *Game) GetMove(i int) (int, error) {
	board := (g.board[g.currentBoard] | (g.board[g.currentBoard] >> 9)) & 0x1FF
	count := 0
	for move := 0; move < boardLength; move++ {
		if board&(0x1<<move) == 0 {
			if i == count {
				return move, nil
			}
			count += 1
		}
	}
	return 0, errors.New("could not find move")
}

func (g *Game) ApplyAction(i int) (gmcts.Game, error) {
	newGame := g.Copy()
	move, err := g.GetMove(i)
	newGame.MakeMove(newGame.currentBoard, move)
	return &newGame, err
}

func (g *Game) Hash() interface{} {
	return fmt.Sprintf("%x", g.board)
}

func (g *Game) Player() gmcts.Player {
	if g.turn {
		return player1
	} else {
		return player2
	}
}

func (g *Game) IsTerminal() bool {
	return len(g.OverallWinner()) > 0
}

func (g *Game) Winners() []gmcts.Player {
	return g.OverallWinner()
}

func NewGame() Game {
	return Game{
		turn:         true,
		board:        [boardLength]uint32{},
		currentBoard: 0x4,
		overallBoard: 0x0,
	}
}

func (g *Game) Copy() Game {
	newG := NewGame()
	newG.turn = g.turn
	newG.currentBoard = g.currentBoard
	newG.overallBoard = g.overallBoard
	copy(newG.board[:], g.board[:])
	return newG
}

func (g *Game) MakeMove(boardIndex byte, pos int) {
	// Check if the board is empty
	if boardIndex == g.currentBoard && g.board[boardIndex]&(1<<pos) == 0 && g.board[boardIndex]&(1<<(pos+9)) == 0 {
		if g.turn {
			g.board[boardIndex] |= 1 << pos

		} else {
			g.board[boardIndex] |= 1 << (pos + 9)
		}
		g.turn = !g.turn
		winners := g.Winner(g.currentBoard)

		// Only switch board is a space is available
		if (g.overallBoard>>pos)&0x1 == 0 && (g.overallBoard>>(pos+9))&0x1 == 0 && (g.overallBoard>>(pos+18))&0x1 == 0 {
			g.currentBoard = byte(pos)
			// The current and target board is completed
		} else if len(winners) > 0 {
			for i := byte(0); i < boardLength; i++ {
				if len(g.Winner(i)) == 0 {
					g.currentBoard = i
				}
			}
		}
	}
}

func (g *Game) OverallWinner() []gmcts.Player {
	if checkCompleted(g.overallBoard & 0x1FF) {
		return []gmcts.Player{player1}
	} else if checkCompleted((g.overallBoard >> 9) & 0x1FF) {
		return []gmcts.Player{player2}
	} else if ((g.overallBoard>>18)|(g.overallBoard>>9)|g.overallBoard)&0x1FF == 0x1FF {
		return []gmcts.Player{player1, player2}
	}
	return []gmcts.Player{}
}

func (g *Game) Winner(boardIndex byte) []gmcts.Player {
	// Player 1
	if g.board[boardIndex]&0x80000000 > 0 || checkCompleted(g.board[boardIndex]&0x1FF) {
		g.board[boardIndex] |= 0x80000000
		g.overallBoard |= 0x1 << boardIndex
		return []gmcts.Player{player1}

		// Player 2
	} else if g.board[boardIndex]&0x40000000 > 0 || checkCompleted((g.board[boardIndex]>>9)&0x1FF) {
		g.board[boardIndex] |= 0x40000000
		g.overallBoard |= 1 << (boardIndex + 9)
		return []gmcts.Player{player2}

		// Draw
	} else if g.board[boardIndex]&0x20000000 > 0 || (g.board[boardIndex]|(g.board[boardIndex]>>9))&0x1FF == 0x1FF {
		g.board[boardIndex] |= 0x20000000
		g.overallBoard |= 1 << (boardIndex + 18)
		return []gmcts.Player{player1, player2}
	}

	return []gmcts.Player{}
}

func checkCompleted(test uint32) bool {
	return (test&0x7) == 0x7 || (test&0x38) == 0x38 || (test&0x1C0) == 0x1C0 || (test&0x49) == 0x49 || (test&0x92) == 0x92 || (test&0x124) == 0x124 || (test&0x111) == 0x111 || (test&0x54) == 0x54
}

func bitCount(u uint32) uint32 {
	uCount := uint32(0)
	uCount = u - ((u >> 1) & 033333333333) - ((u >> 2) & 011111111111)
	return ((uCount + (uCount >> 3)) & 030707070707) % 63
}
