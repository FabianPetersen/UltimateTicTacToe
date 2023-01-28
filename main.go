package main

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/gmcts"
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
	game *gmcts.Game
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

const windowSizeW = 320 * 2
const windowSizeH = 320 * 2
const screenSize = 3.0
const offset = 5
const width = (windowSizeW - offset) / screenSize
const height = (windowSizeH - offset*2) / screenSize

var randSource = rand.New(rand.NewSource(time.Now().Unix()))

const HUMAN = true

func (g *GameEngine) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.game.IsTerminal() {
			g.game = gmcts.NewGame()
		}

		if !HUMAN {
			for i := 0; i < 6 && !g.game.IsTerminal(); i++ {
				if g.game.Player() == gmcts.Player1 {
					randomIndex := randSource.Intn(g.game.Len())
					actualMove, _ := g.game.GetMove(randomIndex)
					g.game.MakeMove(g.game.CurrentBoard, actualMove)
				} else {
					actualMove := g.getBotMove()
					g.game.MakeMove(g.game.CurrentBoard, actualMove)
				}
			}
		}

		if HUMAN && !g.game.IsTerminal() && g.game.Player() == gmcts.Player1 {
			x, y := ebiten.CursorPosition()
			boardIndex, posIndex := g.getBoardPos(float64(x), float64(y))
			g.game.MakeMove(byte(boardIndex), posIndex)
		}
	} else if HUMAN && !g.game.IsTerminal() && g.game.Player() == gmcts.Player2 {
		botmove := g.getBotMove()
		g.game.MakeMove(g.game.CurrentBoard, botmove)
	}

	return nil
}

func (g *GameEngine) getBotMove() int {
	mcts := gmcts.NewMCTS(g.game)
	tree := mcts.SpawnTree(gmcts.ROBUST_CHILD)
	timeToSearch := 200 * time.Millisecond
	tree.Search(timeToSearch)
	mcts.AddTree(tree)
	botMove, _ := mcts.BestAction()
	actualMove, _ := g.game.GetMove(botMove)
	return actualMove
}

func (g *GameEngine) getBoardPos(clickX float64, clickY float64) (boardIndex int, posIndex int) {
	for y := 1.0; y <= screenSize; y++ {
		for x := 1.0; x <= screenSize; x++ {
			if clickX <= x*windowSizeW/screenSize && clickY <= y*windowSizeH/screenSize {
				for yPos := 1.0; yPos <= screenSize; yPos++ {
					for xPos := 1.0; xPos <= screenSize; xPos++ {
						if clickX-(x-1)*windowSizeW/screenSize <= xPos*width/screenSize && clickY-(y-1)*windowSizeH/screenSize <= yPos*height/screenSize {
							return
						}
						posIndex++
					}
				}
			}
			boardIndex++
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
	if g.game.CurrentBoard == boardIndex {
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
		for y := 0.0; y <= 2.0; y++ {
			for x := 0.0; x <= 2.0; x++ {
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
	}

	winners := g.game.Winner(boardIndex)
	if len(winners) > 0 {
		if len(winners) >= 2 {
			ebitenutil.DebugPrintAt(screen, "draw", int(startX+width/4), int(startY+height/4))
		} else if winners[0] == gmcts.Player1 {
			ebitenutil.DrawCircle(screen, startX+width/2, startY+height/2, 70, rgba)
		} else if winners[0] == gmcts.Player2 {
			ebitenutil.DrawRect(screen, startX+width/4, startY+height/4, width/2, height/2, rgba)
		}
	}
}

func (g *GameEngine) Draw(screen *ebiten.Image) {
	boardIndex := byte(0)
	for y := 0.0; y <= 2.0; y++ {
		for x := 0.0; x <= 2.0; x++ {
			g.DrawSingleGameEngine(screen, x, y, boardIndex)
			boardIndex += 1
		}
	}

	winners := g.game.OverallWinner()
	if len(winners) == 1 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Winner %d", winners[0]), windowSizeW/2, windowSizeH/2)
	} else if len(winners) == 2 {
		ebitenutil.DebugPrintAt(screen, "Draw", windowSizeW/2, windowSizeH/2)
	}
}

func (g *GameEngine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return windowSizeW, windowSizeH
}

func main() {
	ebiten.SetWindowSize(windowSizeW, windowSizeH)
	ebiten.SetWindowTitle("Ultimate Tic-Tac-Toe")
	if err := ebiten.RunGame(&GameEngine{
		gmcts.NewGame(),
	}); err != nil {
		log.Fatal(err)
	}
}
