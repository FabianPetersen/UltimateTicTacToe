package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
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

const boardCornerRating = 1.4
const boardSideRating = 1
const boardMiddleRating = 1.75
const posCornerRating = 0.2
const posSideRating = 0.17
const posMiddleRating = 0.22

var boardRating = []float64{boardCornerRating, boardSideRating, boardCornerRating, boardSideRating, boardCornerRating, boardSideRating, boardCornerRating, boardSideRating, boardMiddleRating}
var posRating = []float64{posCornerRating, posSideRating, posCornerRating, posSideRating, posCornerRating, posSideRating, posCornerRating, posSideRating, posMiddleRating}
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

func (g *Game) PopulateBoards() {
	for i := 0; i < 19683; i++ {
		c := i
		var board uint32 = 0
		for ii := 0; ii < 9; ii++ {
			switch c % 3 {
			case 0:
			case 1:
				board |= 0x1 << ii
			case 2:
				board |= 0x1 << (ii + 9)
			}
			c /= 3
		}

		isValid := false

		// Check if board is valid
		xBoard, oBoard, _ := board&0x1FF, (board>>9)&0x1FF, (board&0x1FF)|((board>>9)&0x1FF)

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
	}
}

func CheckCompleted(test uint32) bool {
	if bitCount(test) < 2 {
		return false
	}

	if test&0x100 > 0 && ((test&0x111) == 0x111 || (test&0x144) == 0x144 || (test&0x188) == 0x188) || (test&0x122) == 0x122 {
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
	Board           [boardLength + 1]uint32
	OverallBoard    uint32
	HeuristicScores *HeuristicScores
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

// IsTerminal returns true if this game state is a terminal state
func (g *Game) IsTerminal() bool {
	return CheckCompleted(g.OverallBoard&0x1FF) || CheckCompleted((g.OverallBoard>>9)&0x1FF) || ((g.OverallBoard>>18)|(g.OverallBoard>>9)|g.OverallBoard)&0x1FF == 0x1FF
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

type Storage struct {
	nodeStore map[[10]uint32]*Node
}

func (storage *Storage) Count() int {
	return len(storage.nodeStore)
}

func (storage *Storage) Get(hash GameHash) (*Node, bool) {
	node, exists := storage.nodeStore[*hash]
	return node, exists
}

func (storage *Storage) Set(hash GameHash, node *Node) {
	storage.nodeStore[*hash] = node
}

func (storage *Storage) Reset() {
	storage.nodeStore = make(map[[10]uint32]*Node, 150000)
}

func NewStorage() Storage {
	return Storage{}
}

var TranspositionTable = NewStorage()

type Flag byte

const (
	EXACT       Flag = 0
	UPPER_BOUND Flag = 1
	LOWER_BOUND Flag = 2
)

type Node struct {
	lowerBound float64
	upperBound float64
	bestMove   byte
	bestBoard  byte
	depth      byte
	flag       Flag
}

func NewNode(state *Game) (*Node, bool) {
	// Rotate and invert board to check if it already exists in cache
	var oldNode *Node = nil
	var exists bool = false
	var cacheExists bool = false
	for i := 0; i < 2; i++ {
		if !cacheExists && !exists {
			for r := 0; r < 4; r++ {
				// Check if the board exists in the cache
				if !cacheExists && !exists {
					if oldNode, exists = TranspositionTable.Get(state.Hash()); exists {
						cacheExists = true
					}
				}
				state.Rotate(2)
			}
		}
		state.Invert()
	}

	return oldNode, cacheExists
}

const inf float64 = 10000

func Search(state *Game, alpha float64, beta float64, depth byte, maxPlayer Player, start *time.Time, maxDuration *time.Duration) (float64, byte, byte) {
	// Restore the values from the last node
	n, cached := NewNode(state)
	if cached && n.depth >= depth {
		if n.flag == EXACT {
			return n.lowerBound, n.bestMove, n.bestBoard
		} else if n.flag == LOWER_BOUND {
			alpha = math.Max(alpha, n.lowerBound)
		} else if n.flag == UPPER_BOUND {
			beta = math.Min(beta, n.upperBound)
		}

		if alpha >= beta {
			return n.lowerBound, n.bestMove, n.bestBoard
		}
	}

	var value float64 = 0
	var currentBestMove byte = 0
	var currentBestBoard byte = 0
	var prevBoard = byte(state.Board[PlayerBoardIndex] >> 1)
	if depth == 0 || state.IsTerminal() || time.Since(*start) > *maxDuration {
		return state.HeuristicPlayer(maxPlayer), 0, 0

		// This is a max node
	} else if Player(state.Board[PlayerBoardIndex]&0x1) == maxPlayer {
		value = -inf
		a := alpha
		state.GetMoves(func(boardIndex byte, move byte) bool {
			state.MakeMove(boardIndex, move)
			searchValue, _, _ := Search(state, a, beta, depth-1, maxPlayer, start, maxDuration)
			state.UnMakeMove(move, boardIndex, prevBoard)

			if searchValue >= value {
				value = searchValue
				currentBestMove = move
				currentBestBoard = boardIndex
			}

			a = math.Max(a, value)
			return value >= beta
		})
	} else {
		value = inf
		b := beta
		state.GetMoves(func(boardIndex byte, move byte) bool {
			state.MakeMove(boardIndex, move)
			searchValue, _, _ := Search(state, alpha, b, depth-1, maxPlayer, start, maxDuration)
			state.UnMakeMove(move, boardIndex, prevBoard)
			if searchValue <= value {
				value = searchValue
				currentBestMove = move
				currentBestBoard = boardIndex
			}
			b = math.Min(b, value)
			return value <= alpha
		})
	}

	if !cached {
		n = &Node{
			lowerBound: -inf,
			upperBound: inf,
			bestBoard:  251,
			bestMove:   251,
		}
	}

	/* Traditional transposition table storing of bounds */
	/* Fail low result implies an upper bound */
	if value <= alpha {
		n.upperBound = value
		n.flag = UPPER_BOUND
	}
	/* Found an exact minimax value – will not occur if called with zero window */
	if value > alpha && value < beta {
		n.lowerBound = value
		n.upperBound = value
		n.bestMove = currentBestMove
		n.bestBoard = currentBestBoard
		n.flag = EXACT
	}
	/* Fail high result implies a lower bound */
	if value >= beta {
		n.lowerBound = value
		n.bestMove = currentBestMove
		n.bestBoard = currentBestBoard
		n.flag = LOWER_BOUND
	}
	n.depth = depth
	if !cached {
		TranspositionTable.Set(state.Hash(), n)
	}

	return value, n.bestMove, n.bestBoard
}

func mtdF(state *Game, start *time.Time, maxDuration *time.Duration, f float64, d byte, maxPlayer Player) (float64, byte, byte) {
	g := f
	lowerBound, upperBound := -inf, inf
	beta := -inf
	var bestMove byte = 253
	var bestBoard byte = 253
	var nBestMove, nBestBoard = byte(0), byte(0)
	for lowerBound < upperBound && time.Since(*start) < *maxDuration {
		if g == lowerBound {
			beta = g + 1
		} else {
			beta = g
		}

		g, nBestMove, nBestBoard = Search(state, beta-1, beta, d, maxPlayer, start, maxDuration)
		if nBestBoard < 200 && nBestMove < 200 {
			bestMove = nBestMove
			bestBoard = nBestBoard
		}

		if g < beta {
			upperBound = g
		} else {
			lowerBound = g
		}
	}

	return g, bestMove, bestBoard
}

func IterativeDeepeningTime(state *Game, maxDepth byte, maxTime time.Duration) (byte, byte) {
	// Start the guess at the current heuristic
	var maxPlayer = Player(state.Board[PlayerBoardIndex] & 0x1)
	var firstGuess = state.HeuristicPlayer(maxPlayer)

	var bestMove byte = 255
	var bestBoard byte = 255
	var d byte = 0
	// HeuristicStorage.Reset()
	TranspositionTable.Reset()
	start := time.Now()
	for ; time.Since(start) < maxTime && d < maxDepth; d++ {
		firstGuess, bestMove, bestBoard = mtdF(state, &start, &maxTime, firstGuess, d, maxPlayer)
	}
	fmt.Fprintf(os.Stderr, "Stored nodes, %d Depth %d \n", TranspositionTable.Count(), d)
	return bestMove, bestBoard
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
			moveIndex, boardIndex = IterativeDeepeningTime(game, 20, 93*time.Millisecond)
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
