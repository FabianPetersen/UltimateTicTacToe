package main

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"github.com/FabianPetersen/UltimateTicTacToe/bns"
	"github.com/FabianPetersen/UltimateTicTacToe/mtd"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"image/color"
	"log"
	"math/rand"
	"time"
)

type GameEngine struct {
	game         *Game.Game
	restartCount int
}

var boards = [][2]float64{
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

var boardColors = []color.RGBA{
	colornames.Red,
	colornames.Blue,
	colornames.Yellow,
	colornames.Green,
	colornames.Cyan,
	colornames.Gold,
	colornames.Orange,
	colornames.Magenta,
	colornames.Grey,
}
var activeBoardColor = color.RGBA{G: 255, A: 0x3F}

type BOT_ALGORITHM byte

const (
	MINIMAX                 BOT_ALGORITHM = 0
	MINIMAX_ITERATIVE       BOT_ALGORITHM = 4
	MONTE_CARLO_TREE_SEARCH BOT_ALGORITHM = 1
	MTD_F                   BOT_ALGORITHM = 2
	BNS                     BOT_ALGORITHM = 3
)

var activeBotAlgorithm = MTD_F

const windowSizeW = 320 * 2
const windowSizeH = 320 * 2
const screenSize = 3.0
const offset = 5
const width = (windowSizeW - offset) / screenSize
const height = (windowSizeH - offset*2) / screenSize

var randSource = rand.New(rand.NewSource(time.Now().Unix()))

const HUMAN = true

func (g *GameEngine) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.game = Game.NewGame()
		return nil
	}

	currentPlayer := Game.Player(g.game.Board[Game.PlayerBoardIndex] & 0x1)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if HUMAN && !g.game.IsTerminal() && currentPlayer == Game.Player1 {
			x, y := ebiten.CursorPosition()
			boardIndex, posIndex := g.getBoardPos(float64(x), float64(y))
			if g.game.ValidMove(byte(boardIndex), byte(posIndex)) {
				g.game.MakeMove(byte(boardIndex), byte(posIndex))
				/*
					g.game.UnMakeMove(posIndex, byte(boardIndex))
					for i := 0; i < 4; i++ {
						g.game.Rotate()
					}
					g.game.Invert()
					g.game.Invert()
					g.game.MakeMove(byte(boardIndex), posIndex)
				*/
			}
		}
	} else if HUMAN && !g.game.IsTerminal() && currentPlayer == Game.Player2 {
		board, botmove := g.getBotMove()
		// Check if the board is empty
		if g.game.ValidMove(board, botmove) {
			g.game.MakeMove(board, botmove)
		}
	}

	return nil
}

func (g *GameEngine) getBotMove() (byte, byte) {
	var botMove byte = 254
	var botBoard byte = 254
	switch activeBotAlgorithm {
	case MTD_F:
		botMove, botBoard = mtd.IterativeDeepeningTime(g.game, 15, 100*time.Millisecond)

	case MINIMAX_ITERATIVE:
		//mini := minimax.NewMinimax(g.game)
		//botMove = mini.SearchIterative()

	case BNS:
		botMove = bns.IterativeDeepening(g.game, 10)

		/*
			case MONTE_CARLO_TREE_SEARCH:
				mcts := gmcts.NewMCTS(g.game)
				tree := mcts.SpawnTree(gmcts.ROBUST_CHILD, gmcts.SMITSIMAX)
				tree.SearchTime(timeToSearch)
				mcts.AddTree(tree)
				botMove, _ = mcts.BestAction()
		*/
	}
	return botBoard, botMove
}

func (g *GameEngine) getBoardPos(clickX float64, clickY float64) (boardIndex int, posIndex int) {
	for i, boardPos := range boards {
		x, y := boardPos[0]+1, boardPos[1]+1
		boardIndex = i

		if (x-1)*windowSizeW/screenSize < clickX && clickX <= x*windowSizeW/screenSize && (y-1)*windowSizeH/screenSize < clickY && clickY <= y*windowSizeH/screenSize {
			for iP, pPos := range boards {
				xPos, yPos := pPos[0]+1, pPos[1]+1
				posIndex = iP

				newX := clickX - (x-1)*windowSizeW/screenSize
				newY := clickY - (y-1)*windowSizeH/screenSize
				if (xPos-1)*width/screenSize < newX && newX <= xPos*width/screenSize && (yPos-1)*height/screenSize < newY && newY <= yPos*height/screenSize {
					return
				}
			}
		}
	}
	return 0, 0
}

func (g *GameEngine) DrawSingleGameEngine(screen *ebiten.Image, boardX float64, boardY float64, boardIndex byte) {
	rgba := boardColors[boardIndex]

	// Print board
	startX := offset/2 + (width * boardX)
	startY := offset/2 + (height * boardY)

	// Print current board background color
	if byte(g.game.Board[Game.PlayerBoardIndex]>>1) == boardIndex {
		ebitenutil.DrawRect(screen, startX, startY, width, height, activeBoardColor)
	}

	// Draw the lines
	for i := 0.0; i <= 3.0; i++ {
		x := startX + (i * ((width - offset/3) / 3.0))
		y := startY + (i * ((height - offset/3) / 3.0))
		ebitenutil.DrawLine(screen, x, startY, x, startY+height, rgba)
		ebitenutil.DrawLine(screen, startX, y, startX+width, y, rgba)
	}

	i := 0
	for player := 0; player <= 2; player++ {
		for _, boardPos := range boards {
			x, y := boardPos[0], boardPos[1]

			// The item contains a symbol
			if g.game.Board[boardIndex]&(1<<i) != 0 {

				// The symbol is owned by player
				if player == 0 {
					drawX := startX + ((x + 0.5) * (width / 3.0))
					drawy := startY + ((y + 0.5) * (height / 3.0))
					ebitenutil.DrawCircle(screen, drawX, drawy, width/9, rgba)
				} else {
					drawX := startX + ((x + 0.25) * (width / 3.0))
					drawy := startY + ((y + 0.25) * (height / 3.0))
					ebitenutil.DrawRect(screen, drawX, drawy, width/6.0, height/6.0, rgba)
				}
			}
			i++
		}
	}

	if (g.game.OverallBoard>>boardIndex)&0x1 > 0 {
		ebitenutil.DrawCircle(screen, startX+width/2, startY+height/2, 70, rgba)
	} else if (g.game.OverallBoard>>(boardIndex+9))&0x1 > 0 {
		ebitenutil.DrawRect(screen, startX+width/4, startY+height/4, width/2, height/2, rgba)
	} else if (g.game.OverallBoard>>(boardIndex+18))&0x1 > 0 {
		ebitenutil.DebugPrintAt(screen, "draw", int(startX+width/4)-int(float64(len("Draw"))*3.25), int(startY+height/4))
	}

	p1Score := g.game.HeuristicBoard(Game.Player1, g.game.Board[boardIndex], false)
	p2Score := g.game.HeuristicBoard(Game.Player2, g.game.Board[boardIndex], false)
	p1ScoreString := fmt.Sprintf("Player 1: %.2f", p1Score)
	p2ScoreString := fmt.Sprintf("Player 2: %.2f", p2Score)
	ebitenutil.DebugPrintAt(screen, p1ScoreString, int(startX+2), int(startY+5))
	ebitenutil.DebugPrintAt(screen, p2ScoreString, int(startX+2), int(startY+20))
}

func (g *GameEngine) Draw(screen *ebiten.Image) {
	if byte(g.game.Board[Game.PlayerBoardIndex]>>1) == Game.GlobalBoard {
		ebitenutil.DrawRect(screen, 0, 0, windowSizeW, windowSizeH, activeBoardColor)
	}

	for boardIndex, boardPos := range boards {
		g.DrawSingleGameEngine(screen, boardPos[0], boardPos[1], byte(boardIndex))
	}

	winner := g.game.WinningPlayer()
	if winner == Game.Player1 {
		text := fmt.Sprintf("Winner %d", 0)
		ebitenutil.DebugPrintAt(screen, text, windowSizeW/2-int(float64(len(text))*3.25), windowSizeH/2)

	} else if winner == Game.Player2 {
		text := fmt.Sprintf("Winner %d", 1)
		ebitenutil.DebugPrintAt(screen, text, windowSizeW/2-int(float64(len(text))*3.25), windowSizeH/2)

	} else if winner == Game.Draw {
		ebitenutil.DebugPrintAt(screen, "Draw", windowSizeW/2-int(float64(len("Draw"))*3.25), windowSizeH/2)
	}

	totalScoreString := fmt.Sprintf("Total Player 1: %.2f, Total Player 2: %.2f", g.game.HeuristicPlayer(Game.Player1), g.game.HeuristicPlayer(Game.Player2))
	ebitenutil.DebugPrintAt(screen, totalScoreString, windowSizeW/2-int(float64(len(totalScoreString))*3.25), windowSizeH-20)
}

func (g *GameEngine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowSizeW, windowSizeH
}

func main() {
	ebiten.SetWindowSize(windowSizeW, windowSizeH)
	ebiten.SetWindowTitle("Ultimate Tic-Tac-Toe")
	gameEngine := &GameEngine{
		Game.NewGame(),
		5,
	}

	if err := ebiten.RunGame(gameEngine); err != nil {
		log.Fatal(err)
	}
}
