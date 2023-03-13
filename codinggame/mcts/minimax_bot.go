package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
	"unsafe"
)

var boards = [][2]byte{
	{0, 0},
	{1, 0},
	{2, 0},

	{2, 1},
	{2, 2},

	{1, 2},
	{0, 2},
	{0, 1},
	{1, 1},
}

func translateOwnMove(board byte, move byte) (byte, byte) {
	start := boards[board]
	moveOffset := boards[move]

	return (start[0] * 3) + moveOffset[0], (start[1] * 3) + moveOffset[1]
}

func translateOppMove(x byte, y byte) (byte, byte) {
	for boardIndex := 0; boardIndex < len(boards); boardIndex++ {
		board := boards[boardIndex]
		if x/3 == board[0] && y/3 == board[1] {
			for moveIndex := 0; moveIndex < len(boards); moveIndex++ {
				move := boards[moveIndex]
				if x%3 == move[0] && y%3 == move[1] {
					return byte(boardIndex), byte(moveIndex)
				}
			}
		}
	}
	return 0, 0
}

var BoardHeuristicCacheP1 = map[uint32]float64{}
var BoardHeuristicCacheP2 = map[uint32]float64{}

type HeuristicScores struct {
	BoardRating [9]float64
	PosRating   [9]float64

	BoardCornerRating                        float64
	BoardSideRating                          float64
	BoardMiddleRating                        float64
	PosCornerRating                          float64
	PosSideRating                            float64
	PosMiddleRating                          float64
	OverallWinLossRating                     float64
	OverallAlmostDrawWinLossRating           float64
	GlobalStateRating                        float64
	OverallBoardMultiplierRating             float64
	WinMovesMadeLossRating                   float64
	LossMovesMadeAdvantageRating             float64
	TwoInARowAdvantageRating                 float64
	EnemyTwoInARowLossRating                 float64
	EnemyWonBoardLossRating                  float64
	EnemyWonBoardDiscountRating              float64
	WonBoardRating                           float64
	DrawBoardScoreEnemyDiscountRating        float64
	DrawBoardScorePlayerDiscountRating       float64
	LocalBoardWinPlayedMovesDiscountRating   float64
	OverallBoardWinPlayedMovesDiscountRating float64
}

func DefaultHeuristic() *HeuristicScores {
	h := HeuristicScores{
		BoardCornerRating:                        1.5,
		BoardSideRating:                          0.49,
		BoardMiddleRating:                        2.32,
		PosCornerRating:                          2.03,
		PosSideRating:                            0.61,
		PosMiddleRating:                          1.21,
		OverallWinLossRating:                     5000,
		OverallAlmostDrawWinLossRating:           2500,
		GlobalStateRating:                        70.06,
		OverallBoardMultiplierRating:             150,
		WinMovesMadeLossRating:                   10,
		LossMovesMadeAdvantageRating:             12,
		TwoInARowAdvantageRating:                 9,
		EnemyTwoInARowLossRating:                 9,
		EnemyWonBoardLossRating:                  24,
		EnemyWonBoardDiscountRating:              0.25,
		WonBoardRating:                           24,
		DrawBoardScoreEnemyDiscountRating:        0.12,
		DrawBoardScorePlayerDiscountRating:       0.12,
		LocalBoardWinPlayedMovesDiscountRating:   0.36,
		OverallBoardWinPlayedMovesDiscountRating: 1.32,
	}

	h.BoardRating = [9]float64{h.BoardCornerRating, h.BoardSideRating, h.BoardCornerRating, h.BoardSideRating, h.BoardCornerRating, h.BoardSideRating, h.BoardCornerRating, h.BoardSideRating, h.BoardMiddleRating}
	h.PosRating = [9]float64{h.PosCornerRating, h.PosSideRating, h.PosCornerRating, h.PosSideRating, h.PosCornerRating, h.PosSideRating, h.PosCornerRating, h.PosSideRating, h.PosMiddleRating}

	return &h
}

func getOffset(player Player) (int, int) {
	if player == Player1 {
		return 0, 9
	}
	return 9, 0
}

func (g *Game) HeuristicBoard(player Player, board uint32, isOverallBoard bool) float64 {
	if !isOverallBoard {
		if player == Player1 {
			if score, ok := BoardHeuristicCacheP1[board&0x3FFFF]; ok {
				return score
			}
		} else if score, ok := BoardHeuristicCacheP2[board&0x3FFFF]; ok {
			return score
		}
	}

	offset, enemyOffset := getOffset(player)
	var score float64 = 0

	// Calculate pos rating
	playerBoard := (board >> offset) & 0x1FF
	enemyBoard := (board >> enemyOffset) & 0x1FF
	jointBoard := (board>>18)&0x1FF | playerBoard | enemyBoard

	// Check if won
	if CheckCompleted(playerBoard) {
		// Give a discount on the amount of moves made in the board
		// To incentivise a lower number of total moves
		if !isOverallBoard {
			score += g.HeuristicScores.WonBoardRating - float64(bitCount(playerBoard))*g.HeuristicScores.LocalBoardWinPlayedMovesDiscountRating
		} else {
			score += g.HeuristicScores.WonBoardRating - float64(bitCount(playerBoard))*g.HeuristicScores.OverallBoardWinPlayedMovesDiscountRating
		}

		// The board is a draw
	} else if jointBoard == 0x1FF && !CheckCompleted(enemyBoard) {
		// Give a reward for the amount of wasted enemy moves (or won moves
		if !isOverallBoard {
			score += float64(bitCount(enemyBoard)) * g.HeuristicScores.DrawBoardScoreEnemyDiscountRating
		} else {
			score += float64(bitCount(playerBoard)) * g.HeuristicScores.DrawBoardScorePlayerDiscountRating
		}

	} else {
		// Calculate pos for items
		for i := 0; i < boardLength; i++ {
			if playerBoard&(0x1<<i) > 0 {
				score += g.HeuristicScores.PosRating[i]
			}
		}

		// If the board is not won by enemy
		if !CheckCompleted(enemyBoard) {
			// Check 2 joint items
			if checkCloseWinningSequence(playerBoard, jointBoard) > 0 {
				score += g.HeuristicScores.TwoInARowAdvantageRating
			}

			// Check 2 joint items
			if checkCloseWinningSequence(enemyBoard, jointBoard) > 0 {
				score -= g.HeuristicScores.EnemyTwoInARowLossRating
			}

			// The enemy has won a square
		} else {
			// Give a reward for enemy moves
			score -= g.HeuristicScores.EnemyWonBoardLossRating - float64(bitCount(enemyBoard))*g.HeuristicScores.EnemyWonBoardDiscountRating
		}
	}
	return score
}

func (g *Game) HeuristicPlayer(player Player) float64 {
	// Don't rerun calculation unless needed
	var score float64 = 0
	var playerOffset, enemyOffset = getOffset(player)

	// Calculate pos rating
	playerBoard := (g.OverallBoard >> playerOffset) & 0x1FF
	enemyBoard := (g.OverallBoard >> enemyOffset) & 0x1FF
	jointBoard := (g.OverallBoard>>18)&0x1FF | playerBoard | enemyBoard

	if CheckCompleted(playerBoard) {
		// Incentivise quicker wins
		return g.HeuristicScores.OverallWinLossRating - float64(g.MovesMade())*g.HeuristicScores.WinMovesMadeLossRating
	} else if CheckCompleted(enemyBoard) {
		// Incentivise slower losses
		return -g.HeuristicScores.OverallWinLossRating + float64(g.MovesMade())*g.HeuristicScores.LossMovesMadeAdvantageRating
	} else if jointBoard == 0x1FF {
		if bitCount(playerBoard) > bitCount(enemyBoard) {
			return g.HeuristicScores.OverallAlmostDrawWinLossRating - float64(g.MovesMade())*g.HeuristicScores.WinMovesMadeLossRating
		}
		return -g.HeuristicScores.OverallAlmostDrawWinLossRating + float64(g.MovesMade())*g.HeuristicScores.LossMovesMadeAdvantageRating
	}

	if byte(g.Board[PlayerBoardIndex]>>1) == GlobalBoard {
		if byte(g.Board[PlayerBoardIndex]&0x1) == byte(player) {
			score += g.HeuristicScores.GlobalStateRating

		} else {
			score -= g.HeuristicScores.GlobalStateRating
		}
	}

	for i := 0; i < boardLength; i++ {
		boardScore := g.HeuristicBoard(player, g.Board[i], false)
		score += boardScore * g.HeuristicScores.BoardRating[i]
	}

	score += g.HeuristicBoard(player, g.OverallBoard, true) * g.HeuristicScores.OverallBoardMultiplierRating
	return score
}

func (g *Game) MovesMade() uint32 {
	// Count how far the game has progressed
	var movesPlayed uint32 = 0
	for i := 0; i < 9; i++ {
		movesPlayed += bitCount(g.Board[i] & 0x3FFFF)
	}
	return movesPlayed
}

func popBoardHelper(b uint32, pos int, player int) uint32 {
	if player == 0 {
		b |= 0x1 << pos
	} else if player == 1 {
		b |= 0x1 << (pos + 9)
	}
	return b
}

func (g *Game) PopulateBoards() {
	var board uint32 = 0
	for i0 := 0; i0 < 3; i0++ {
		for i1 := 0; i1 < 3; i1++ {
			for i2 := 0; i2 < 3; i2++ {
				for i3 := 0; i3 < 3; i3++ {
					for i4 := 0; i4 < 3; i4++ {
						for i5 := 0; i5 < 3; i5++ {
							for i6 := 0; i6 < 3; i6++ {
								for i7 := 0; i7 < 3; i7++ {
									for i8 := 0; i8 < 3; i8++ {
										board = 0
										board = popBoardHelper(board, 0, i0)
										board = popBoardHelper(board, 1, i1)
										board = popBoardHelper(board, 2, i2)
										board = popBoardHelper(board, 3, i3)
										board = popBoardHelper(board, 4, i4)
										board = popBoardHelper(board, 5, i5)
										board = popBoardHelper(board, 6, i6)
										board = popBoardHelper(board, 7, i7)
										board = popBoardHelper(board, 8, i8)

										// Check if board is valid
										isValid := false
										xBoard, oBoard, jointBoard := board&0x1FF, (board>>9)&0x1FF, uint16((board&0x1FF)|((board>>9)&0x1FF))

										// Board can be valid only if either
										xWin, oWin := CheckCompleted(xBoard), CheckCompleted(oBoard)
										if xWin && !oWin {
											isValid = true
										} else if !xWin && oWin {
											isValid = true
										} else if !xWin && !oWin {
											isValid = true
										}

										if isValid {
											BoardHeuristicCacheP1[board] = g.HeuristicBoard(Player1, board, false)
											BoardHeuristicCacheP2[board] = g.HeuristicBoard(Player2, board, false)
										}

										MovesStorage[jointBoard] = []byte{}
										randSource.Shuffle(len(moveOrder), func(i, j int) {
											moveOrder[i], moveOrder[j] = moveOrder[j], moveOrder[i]
										})
										for _, move := range moveOrder {
											if jointBoard&(0x1<<move) == 0 {
												MovesStorage[jointBoard] = append(MovesStorage[jointBoard], move)
											}
										}
										MovesLengthStorage[jointBoard] = byte(len(MovesStorage[jointBoard]))
										BoardCompletedStorage[jointBoard] = CheckCompletedHelper(uint32(jointBoard))
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

var BoardCompletedStorage = [512]bool{}

func CheckCompleted(test uint32) bool {
	return BoardCompletedStorage[test]
}

func CheckCompletedHelper(test uint32) bool {
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
const Draw Player = 2
const boardLength = 9
const GlobalBoard byte = 0xF0
const PlayerBoardIndex = 9

/*	corner -> middle -> side */
var moveOrder = []byte{0, 4, 2, 6, 8, 1, 3, 5, 7}

var randSource = rand.New(rand.NewSource(time.Now().Unix()))

var x = uint64(1) /*  time.Now().Unix() initial seed must be nonzero, don't use a static variable for the state if multithreaded */
func xorshift64star(n byte) byte {
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	return byte((x * 2685821657736338717) % uint64(n))
}

type Game struct {
	Board           [boardLength + 1]uint32
	OverallBoard    uint32
	HeuristicScores *HeuristicScores
}

func (g *Game) GetMoves(executeMove func(byte, byte) bool) {
	currentBoard := byte(g.Board[PlayerBoardIndex] >> 1)
	if currentBoard < 9 {
		for _, move := range MovesStorage[((g.Board[currentBoard] | (g.Board[currentBoard] >> 9)) & 0x1FF)] {
			if executeMove(currentBoard, move) {
				return
			}
		}
	} else {
		jointOverallBoard = (g.OverallBoard>>9 | g.OverallBoard | g.OverallBoard>>18) & 0x1FF
		for _, i := range MovesStorage[jointOverallBoard] {
			for _, move := range MovesStorage[((g.Board[i] | (g.Board[i] >> 9)) & 0x1FF)] {
				if executeMove(i, move) {
					return
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

// IsTerminal returns true if this game state is a terminal state
func (g *Game) IsTerminal() bool {
	return BoardCompletedStorage[g.OverallBoard&0x1FF] || BoardCompletedStorage[(g.OverallBoard>>9)&0x1FF] || ((g.OverallBoard>>18)|(g.OverallBoard>>9)|g.OverallBoard)&0x1FF == 0x1FF
}

func (g *Game) WinningPlayer() Player {
	if CheckCompleted(g.OverallBoard & 0x1FF) {
		return Player1

	} else if CheckCompleted((g.OverallBoard >> 9) & 0x1FF) {
		return Player2

	} else if ((g.OverallBoard>>18)|(g.OverallBoard>>9)|g.OverallBoard)&0x1FF == 0x1FF {
		p1 := bitCount(g.OverallBoard & 0x1FF)
		p2 := bitCount((g.OverallBoard >> 9) & 0x1FF)
		if p1 > p2 {
			return Player1
		} else if p2 > p1 {
			return Player2
		}
	}
	return Draw
}

func NewGame() *Game {
	g := &Game{
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
		OverallBoard:    0x0,
		HeuristicScores: DefaultHeuristic(),
	}
	g.PopulateBoards()
	return g
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
	if currentBoard != 8 && currentBoard < 9 {
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
	p := byte(g.Board[PlayerBoardIndex] & 0x1)
	g.Board[boardIndex] |= 1 << (pos + 9*p)
	if BoardCompletedStorage[(g.Board[boardIndex]>>(9*p))&0x1FF] {
		g.OverallBoard |= 1 << (boardIndex + (9 * p))
	}

	// Draw
	if (g.Board[boardIndex]|(g.Board[boardIndex]>>9))&0x1FF == 0x1FF {
		g.OverallBoard |= 1 << (boardIndex + 18)
	}

	// Only switch board is a space is available
	g.Board[PlayerBoardIndex] = uint32(pos)<<1 | ((g.Board[PlayerBoardIndex] & 0x1) ^ 0x1)

	// If board is finished
	if g.OverallBoard&(0x1<<pos|0x1<<(pos+9)|0x1<<(pos+18)) != 0 {
		g.Board[PlayerBoardIndex] |= 0x100
	}
}

func (g *Game) ValidMove(boardIndex byte, pos byte) bool {
	currentBoard := byte(g.Board[PlayerBoardIndex] >> 1)
	return !g.IsBoardFinished(boardIndex) && (currentBoard == GlobalBoard || boardIndex == currentBoard) && g.Board[boardIndex]&(1<<pos) == 0 && g.Board[boardIndex]&(1<<(pos+9)) == 0
}

func (g *Game) IsBoardFinished(pos byte) bool {
	return g.OverallBoard&(0x1<<pos|0x1<<(pos+9)|0x1<<(pos+18)) != 0
}

var MovesStorage = [512][]byte{}
var MovesLengthStorage = [512]byte{}

var boardIndex byte
var moveIndex byte
var jointOverallBoard uint32
var moves byte
var currentMoves byte
var board uint32
var relativeMoveIndex byte

func (g *Game) Len() byte {
	boardIndex = byte(g.Board[PlayerBoardIndex] >> 1)
	if boardIndex < 9 {
		return MovesLengthStorage[(g.Board[boardIndex]|(g.Board[boardIndex]>>9))&0x1FF]
	}

	moves = 0
	jointOverallBoard = (g.OverallBoard>>9 | g.OverallBoard | g.OverallBoard>>18) & 0x1FF
	for _, i := range MovesStorage[jointOverallBoard] {
		// Check if the board is open
		moves += MovesLengthStorage[(g.Board[i]|(g.Board[i]>>9))&0x1FF]
	}

	return moves
}

func (g *Game) MakeMoveRandUntilTerminal() {
	//for !g.IsTerminal() {
	jointOverallBoard = (g.OverallBoard>>9 | g.OverallBoard | g.OverallBoard>>18) & 0x1FF
	for !(BoardCompletedStorage[g.OverallBoard&0x1FF] || BoardCompletedStorage[(g.OverallBoard>>9)&0x1FF] || jointOverallBoard == 0x1FF) {
		boardIndex = byte(g.Board[PlayerBoardIndex] >> 1)
		moveIndex = byte(randSource.Intn(int(g.Len())))
		// moveIndex = xorshift64star(g.Len())

		if boardIndex < 9 {
			g.MakeMove(boardIndex, MovesStorage[(g.Board[boardIndex]|(g.Board[boardIndex]>>9))&0x1FF][moveIndex])
			jointOverallBoard = (g.OverallBoard>>9 | g.OverallBoard | g.OverallBoard>>18) & 0x1FF
			continue
		}

		moves = 0
		for _, i := range MovesStorage[jointOverallBoard] {
			board = (g.Board[i] | (g.Board[i] >> 9)) & 0x1FF
			currentMoves = MovesLengthStorage[board]

			if moves+currentMoves > moveIndex {
				relativeMoveIndex = moveIndex - moves
				g.MakeMove(i, MovesStorage[board][relativeMoveIndex])
				jointOverallBoard = (g.OverallBoard>>9 | g.OverallBoard | g.OverallBoard>>18) & 0x1FF
				continue
			}
			moves += currentMoves
		}
	}
}

type Node struct {
	parent        *Node
	children      []*Node
	childrenCount byte
	maxChildren   byte

	move  byte
	board byte

	nodeScore   uint16
	nodeVisits  uint16
	nodeExploit float32
}

type BestActionPolicy byte

const (
	MAX_CHILD_SCORE BestActionPolicy = 0
	ROBUST_CHILD    BestActionPolicy = 1
)

const (
	//DefaultExplorationConst is the default exploration constant of UCT2 Formula
	//Sqrt(2) is a frequent choice for this constant as specified by
	//https://en.wikipedia.org/wiki/Monte_Carlo_tree_search
	DefaultExplorationConst = float32(math.Sqrt2) - 1
)

const magic32 = 0x5F375A86
const th = 1.5

func FastSqrt32(n float32) float32 {
	b := *(*uint64)(unsafe.Pointer(&n))
	b = magic32 - (b >> 1)
	y := *(*float32)(unsafe.Pointer(&b))
	y *= th - (n * 0.5 * y * y)
	return 1 / y
}

func ln(x float32) float32 {
	var bx = *(*uint32)(unsafe.Pointer(&x))
	var ex = bx >> 23
	var t = float32(ex) - 127
	bx = 1065353216 | (bx & 8388607)
	x = *(*float32)(unsafe.Pointer(&bx))
	return -1.49278 + (2.11263+(-0.729104+0.10969*x)*x)*x + 0.6931471806*t
}

// UCT2 algorithm is described in this paper
// https://www.csse.uwa.edu.au/cig08/Proceedings/papers/8057.pdf
func (n *Node) UCT2(i byte) float32 {
	explore := ln(float32(n.nodeVisits)) / float32(n.children[i].nodeVisits) // math.Log(float64(n.nodeVisits)) / float64(n.children[i].nodeVisits)
	explore = float32(math.Sqrt(float64(explore)))                           // FastSqrt32(explore)                                            // float32(math.Sqrt(float64(explore)))                           // 1 / FastInvSqrt64(explore) // math.Sqrt(explore) // 1 / FastInvSqrt64(explore) //

	return n.nodeExploit + DefaultExplorationConst*explore
}

// smitsimax Node selection algorithm is described in this paper
// https://www.codingame.com/playgrounds/36476/smitsimax
/*
func (n *Node) smitsimax(i int, p Player) float64 {
	exploit := 0.3 * n.children[i].nodeScore[p] / float64(n.children[i].nodeVisits)
	exploit += 0.7 * n.children[i].heuristicScore[p] / float64(n.children[i].nodeVisits)

	explore := math.Log(float64(n.nodeVisits)) / n.childVisits[i]
	explore = math.Sqrt(explore)

	return exploit + DefaultExplorationConst*explore
}
*/

func (node *Node) treePolicy() *Node {
	var bestScore float32 = 0
	var bestNode = node.children[0]
	for i := byte(0); i < node.childrenCount; i++ {
		score := node.UCT2(i)
		if score >= bestScore {
			bestScore = score
			bestNode = node.children[i]
		}
	}
	return bestNode
}

const bestActionPolicy = ROBUST_CHILD

var nodePool = [700000]Node{}
var NodePoolIndex = 1

// MCTS contains functionality for the MCTS algorithm
type MCTS struct {
	game     *Game
	gameCopy *Game
	root     *Node
}

var ggCopy Game

// NewMCTS returns a new MCTS wrapper
func NewMCTS(initial *Game) *MCTS {
	NodePoolIndex = 1

	nodePool[0].parent = nil
	nodePool[0].nodeVisits = 1
	nodePool[0].nodeScore = 0
	nodePool[0].childrenCount = 0
	return &MCTS{
		game:     initial,
		gameCopy: &ggCopy,
		root:     &nodePool[0],
	}
}

var player Player
var winningPlayer Player
var node *Node = nil
var availableMoves byte = 0

func (m *MCTS) search() {
	// Selection
	node = m.root
	m.gameCopy.OverallBoard = m.game.OverallBoard
	player = Player(m.gameCopy.Board[PlayerBoardIndex] & 0x1)
	for i := 0; i < 10; i++ {
		m.gameCopy.Board[i] = m.game.Board[i]
	}

	for node.childrenCount > 0 {
		// Check children (tree policy)
		node = node.treePolicy()
		m.gameCopy.MakeMove(node.board, node.move)
	}

	// Expansion
	if !m.gameCopy.IsTerminal() {
		// Fill out the slice to make room for new items
		availableMoves = m.gameCopy.Len()
		if node.maxChildren < availableMoves {
			node.children = append(node.children, make([]*Node, availableMoves-node.maxChildren)...)
			node.maxChildren = availableMoves
		}

		// Iterate over all children
		node.childrenCount = 0
		m.gameCopy.GetMoves(func(board byte, move byte) bool {
			NodePoolIndex++
			node.children[node.childrenCount] = &nodePool[NodePoolIndex]
			node.children[node.childrenCount].parent = node
			node.children[node.childrenCount].move = move
			node.children[node.childrenCount].board = board
			node.children[node.childrenCount].nodeVisits = 1
			node.children[node.childrenCount].nodeScore = 0
			node.children[node.childrenCount].childrenCount = 0
			node.childrenCount++
			return false
		})

		node = node.treePolicy()
		m.gameCopy.MakeMove(node.board, node.move)
	}

	// Simulation
	m.gameCopy.MakeMoveRandUntilTerminal()

	// Backpropagation
	player = m.gameCopy.WinningPlayer()
	for node.parent != nil {
		if player == winningPlayer {
			node.nodeScore += 2
		} else if winningPlayer == 2 {
			node.nodeScore += 1
		}
		node.nodeVisits += 1
		node.nodeExploit = float32(node.nodeScore>>1) / float32(node.nodeVisits)
		node = node.parent
	}
}

func (t *MCTS) BestAction() (byte, byte) {
	var bestAction byte
	var bestBoard byte
	//Select the child with the highest winrate
	if bestActionPolicy == MAX_CHILD_SCORE {
		var bestWinRate float32 = 0
		player = Player(t.game.Board[PlayerBoardIndex] & 0x1)
		for i := byte(0); i < t.game.Len(); i++ {
			winRate := float32(t.root.children[i].nodeScore>>1) / float32(t.root.children[i].nodeVisits)
			if winRate >= bestWinRate {
				bestAction = t.root.children[i].move
				bestBoard = t.root.children[i].board
				bestWinRate = winRate
			}
		}
	} else if bestActionPolicy == ROBUST_CHILD {
		var mostVisists uint16 = 1
		for i := byte(0); i < t.game.Len(); i++ {
			if t.root.children[i].nodeVisits >= mostVisists {
				bestAction = t.root.children[i].move
				bestBoard = t.root.children[i].board
				mostVisists = t.root.children[i].nodeVisits
			}
		}
	}

	return bestAction, bestBoard
}

// SearchTime searches the tree for a specified time
func (t *MCTS) SearchTime(duration time.Duration) {
	var i int
	end := time.Now().Add(duration)
	for {
		if i&0x3F == 0 { // Check in every 128th iteration
			if time.Now().After(end) {
				fmt.Fprintf(os.Stderr, "Rounds %d\n", i)
				break
			}
		}
		t.search()
		i++
	}
}

// SearchRounds searches the tree for a specified number of rounds
//
// SearchRounds will panic if the Game's ApplyAction
// method returns an error or if any game state's Hash()
// method returns a noncomparable value.
func (t *MCTS) SearchRounds(rounds int) {
	for i := 0; i < rounds; i++ {
		t.search()
	}
}

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	game := NewGame()
	first := true

	for {
		var opponentRow, opponentCol int
		fmt.Scan(&opponentRow, &opponentCol)
		if opponentRow >= 0 {
			boardIndex, moveIndex := translateOppMove(byte(opponentRow), byte(opponentCol))
			fmt.Fprintf(os.Stderr, "Opp move BoardIndex %d MoveIndex %d \n", moveIndex, boardIndex)
			game.MakeMove(boardIndex, moveIndex)
			first = false
		}

		var validActionCount int
		fmt.Scan(&validActionCount)

		for i := 0; i < validActionCount; i++ {
			var row, col int
			fmt.Scan(&row, &col)
		}

		moveIndex, boardIndex := byte(8), byte(8)
		if !first {
			start := time.Now()
			mcts := NewMCTS(game)
			mcts.SearchTime(99 * time.Millisecond)
			moveIndex, boardIndex = mcts.BestAction()
			fmt.Fprintf(os.Stderr, "Our move BoardIndex %d MoveIndex %d \n", moveIndex, boardIndex)
			fmt.Fprintln(os.Stderr, time.Since(start))
		}
		first = false
		game.MakeMove(boardIndex, moveIndex)

		// fmt.Fprintln(os.Stderr, "Debug messages...")
		x, y := translateOwnMove(boardIndex, moveIndex)
		fmt.Printf("%d %d\n", x, y) // Write action to stdout
	}
}
