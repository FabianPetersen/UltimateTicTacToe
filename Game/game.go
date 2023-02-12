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
//Player is an id for the player
type Player byte
type GameHash interface{}

const Player1 Player = 1
const Player2 Player = 2
const boardLength = 9

var hash = make([]byte, (boardLength*4)+2)
var randSource = rand.New(rand.NewSource(time.Now().Unix()))

/*	Middle -> corner -> other */
var moveOrder = []byte{4, 0, 8, 2, 6, 1, 3, 5, 7}

type Game struct {
	CurrentPlayer  Player
	CurrentBoard   byte
	Board          [boardLength]uint32
	overallBoard   uint32
	Turn           byte // Should perhaps be an int in other games
	heuristicStore map[Player]float64
}

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

//Len returns the number of actions to consider.
func (g *Game) Len() int {
	return len(g.getFilteredAvailableMoves())
}

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

func (g *Game) GetMove(i int) (int, error) {
	moves := g.getFilteredAvailableMoves()
	return int(moves[i]), nil
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
	h := 0
	for i := 0; i < boardLength; i++ {
		h = 4 * i
		hash[h+0] = byte(g.Board[i])
		hash[h+1] = byte(g.Board[i] >> 8)
		hash[h+2] = byte(g.Board[i] >> 16)
		hash[h+3] = byte(g.Board[i] >> 24)
	}
	hash[36] = byte(g.CurrentPlayer)
	hash[37] = g.Turn
	return GameHash(XXH3_64bits(hash))
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
		CurrentPlayer:  Player1,
		Board:          [boardLength]uint32{},
		CurrentBoard:   0x4,
		overallBoard:   0x0,
		heuristicStore: map[Player]float64{},
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

func (g *Game) MakeMove(boardIndex byte, pos int) {
	g.heuristicStore = map[Player]float64{}
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
				}
			}
		}
		g.Turn += 1
	}
}

func (g *Game) IsBoardFinished(pos int) bool {
	return !((g.overallBoard>>pos)&0x1 == 0 && (g.overallBoard>>(pos+9))&0x1 == 0 && (g.overallBoard>>(pos+18))&0x1 == 0)
}

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
