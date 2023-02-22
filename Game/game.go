package Game

import (
	"errors"
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
type GameHash *[9]uint32

const Player1 Player = 1
const Player2 Player = 2
const boardLength = 9

var randSource = rand.New(rand.NewSource(time.Now().Unix()))

/*	corner -> middle -> side */
var moveOrder = []byte{0, 4, 2, 6, 8, 1, 3, 5, 7}

type Game struct {
	CurrentPlayer Player
	CurrentBoard  byte
	Board         [boardLength]uint32
	overallBoard  uint32

	lastBoard byte
	lastPos   int
}

//Len returns the number of actions to consider.
func (g *Game) Len() int {
	return int(boardLength - bitCount((g.Board[g.CurrentBoard]|(g.Board[g.CurrentBoard]>>9))&0x1FF))
	//return len(g.getFilteredAvailableMoves())
}

func (g *Game) GetMove(i int) (int, error) {
	count := 0
	board := (g.Board[g.CurrentBoard] | (g.Board[g.CurrentBoard] >> 9)) & 0x1FF
	for _, move := range moveOrder {
		if board&(0x1<<move) == 0 {
			if i == count {
				return int(move), nil
			}
			count += 1
		}
	}
	panic(errors.New("data"))
	return 0, errors.New("could not find move")
	/*
		moves := g.getFilteredAvailableMoves()
		return int(moves[i]), nil
	*/
}

//ApplyAction applies the ith action (0-indexed) to the game state,
//and returns a new game state and an error for invalid actions
func (g *Game) ApplyAction(i int) (*Game, error) {
	newGame := g.Copy()
	move, err := g.GetMove(i)
	newGame.MakeMove(newGame.CurrentBoard, move)
	return newGame, err
}

//ApplyAction applies the ith action (0-indexed) to the game state,
// This will modify the game object
func (g *Game) ApplyActionModify(i int) {
	move, _ := g.GetMove(i)
	g.MakeMove(g.CurrentBoard, move)
}

//Hash returns a unique representation of the state.
//Any return value must be comparable.
//This is to separate states that seemingly look the same,
//but actually occur on different turn orders. Without this,
//the directed acyclic graph will become a directed cyclic graph,
//which this MCTS implementation cannot handle properly.
func (g *Game) Hash() GameHash {
	return &g.Board
}

func (g *Game) Compare(c *Game) bool {
	for i := 0; i < boardLength; i++ {
		if g.Board[i] != c.Board[i] {
			return false
		}
	}

	return g.CurrentBoard == c.CurrentBoard && g.CurrentPlayer == c.CurrentPlayer
}

//Player returns the player that can take the next action
func (g *Game) Player() Player {
	if g.CurrentPlayer == Player1 {
		return Player2
	} else {
		return Player1
	}
}

//IsTerminal returns true if this game state is a terminal state
func (g *Game) IsTerminal() bool {
	return len(g.OverallWinner()) > 0
}

//Winners returns a list of players that have won the game if
//IsTerminal() returns true
func (g *Game) Winners() []Player {
	return g.OverallWinner()
}

func NewGame() *Game {
	return &Game{
		CurrentPlayer: Player1,
		Board:         [boardLength]uint32{},
		CurrentBoard:  0x8,
		overallBoard:  0x0,
	}
}

func (g *Game) Copy() *Game {
	newG := NewGame()
	newG.CurrentPlayer = g.CurrentPlayer
	newG.CurrentBoard = g.CurrentBoard
	newG.overallBoard = g.overallBoard
	copy(newG.Board[:], g.Board[:])
	return newG
}

func (g *Game) Rotate() {
	g.Board[2], g.Board[3], g.Board[4], g.Board[5], g.Board[6], g.Board[7], g.Board[0], g.Board[1] = g.Board[0], g.Board[1], g.Board[2], g.Board[3], g.Board[4], g.Board[5], g.Board[6], g.Board[7]
	for i := 0; i < boardLength; i++ {
		g.Board[i] = g.Board[i]&0xe0020100 | rotl2(g.Board[i]) | (rotl2(g.Board[i]>>9) << 9)
	}
	g.overallBoard = rotl2(g.overallBoard) | (rotl2(g.overallBoard>>9) << 9) | (rotl2(g.overallBoard>>18) << 18)

	// Change current board
	if g.CurrentBoard != 8 {
		g.CurrentBoard = (g.CurrentBoard + 2) % 8
	}
}

func (g *Game) Invert() {
	for i := 0; i < boardLength; i++ {
		g.Board[i] = g.Board[i]&0x80000000 | (g.Board[i]&0x40000000>>1)&0x20000000 | (g.Board[i]&0x20000000<<1)&0x40000000 | (g.Board[i]&0x1FF)<<9 | (g.Board[i]>>9)&0x1FF
	}
	g.overallBoard = g.overallBoard&0x7FC0000 | (g.overallBoard&0x1FF)<<9 | (g.overallBoard>>9)&0x1FF

	// Change player
	if g.CurrentPlayer == Player1 {
		g.CurrentPlayer = Player2
	} else {
		g.CurrentPlayer = Player1
	}
}

func (g *Game) UnMakeMove() {
	// Unset move
	if g.CurrentPlayer == Player2 {
		g.Board[g.lastBoard] &^= 1<<g.lastPos | 0x80000000 | 0x40000000 | 0x20000000
		g.CurrentPlayer = Player1

	} else {
		g.Board[g.lastBoard] &^= 1<<(g.lastPos+9) | 0x80000000 | 0x40000000 | 0x20000000
		g.CurrentPlayer = Player2
	}

	// Reset win
	g.overallBoard &^= 0x1<<g.lastBoard | 0x1<<(g.lastBoard+9) | 0x1<<(g.lastBoard+18)
	g.CurrentBoard = g.lastBoard
}

func (g *Game) MakeMove(boardIndex byte, pos int) {
	g.lastBoard = g.CurrentBoard
	g.lastPos = pos

	// Check if the board is empty
	if boardIndex == g.CurrentBoard && g.Board[boardIndex]&(1<<pos) == 0 && g.Board[boardIndex]&(1<<(pos+9)) == 0 {
		if g.CurrentPlayer == Player1 {
			g.Board[boardIndex] |= 1 << pos
			g.CurrentPlayer = Player2

		} else {
			g.Board[boardIndex] |= 1 << (pos + 9)
			g.CurrentPlayer = Player1
		}
		winners := g.Winner(g.CurrentBoard)

		// Only switch board is a space is available
		if !g.IsBoardFinished(pos) {
			g.CurrentBoard = byte(pos)
			// The current and target board is completed
		} else if len(winners) > 0 {
			for i := byte(0); i < boardLength; i++ {
				if len(g.Winner(i)) == 0 {
					g.CurrentBoard = i
					break
				}
			}
		}
	}
}

func (g *Game) IsBoardFinished(pos int) bool {
	return !((g.overallBoard>>pos)&0x1 == 0 && (g.overallBoard>>(pos+9))&0x1 == 0 && (g.overallBoard>>(pos+18))&0x1 == 0)
}
