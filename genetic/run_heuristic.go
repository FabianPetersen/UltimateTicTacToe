package main

import (
	"encoding/json"
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"github.com/FabianPetersen/UltimateTicTacToe/mtd"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/tomcraven/goga"
)

const totalParams = 22
const bitsPerParam = 12
const bitsPerParamSmall = 8
const divisionPerParam = 100
const population = 200

var bbitsPerParam = []int{
	bitsPerParamSmall,
	bitsPerParamSmall,
	bitsPerParamSmall,
	bitsPerParamSmall,
	bitsPerParamSmall,
	bitsPerParamSmall,
	bitsPerParam,
	bitsPerParam,
	bitsPerParam,
	bitsPerParam,
	bitsPerParam,
	bitsPerParam,
	bitsPerParam,
	bitsPerParam,
	bitsPerParamSmall,
	bitsPerParam,
	bitsPerParamSmall,
	bitsPerParamSmall,
	bitsPerParamSmall,
	bitsPerParamSmall,
}
var bbitsPerParamOffset = []int{}
var genAlgo = goga.NewGeneticAlgorithm()
var oppHeuristic = Game.DefaultHeuristic()

func getFloat64(bits goga.Bitset) float64 {
	value := uint32(0)
	for i := 0; i < bits.GetSize(); i++ {
		if bits.Get(i) == 1 {
			value |= 1 << i
		}
	}
	return float64(value) / divisionPerParam
}

func getBits(f float64) goga.Bitset {
	value := uint32(f * divisionPerParam)
	b := goga.Bitset{}
	b.Create(bitsPerParam)
	for i := 0; i < bitsPerParam; i++ {
		if value&(1<<i) > 0 {
			b.Set(i, 1)
		}
	}
	return b
}

type utttHeuristicBitsetCreate struct{}

func (bc *utttHeuristicBitsetCreate) Go() goga.Bitset {
	b := goga.Bitset{}
	b.Create(totalParams * bitsPerParam)

	// Convert heuristic object to bitset
	h := Game.DefaultHeuristic()
	params := []float64{
		h.BoardCornerRating,
		h.BoardSideRating,
		h.BoardMiddleRating,
		h.PosCornerRating,
		h.PosSideRating,
		h.PosMiddleRating,
		h.OverallWinLossRating,
		h.OverallAlmostDrawWinLossRating,
		h.GlobalStateRating,
		h.OverallBoardMultiplierRating,
		h.WinMovesMadeLossRating,
		h.LossMovesMadeAdvantageRating,
		h.TwoInARowAdvantageRating,
		h.EnemyTwoInARowLossRating,
		h.EnemyWonBoardLossRating,
		h.EnemyWonBoardDiscountRating,
		h.WonBoardRating,
		h.DrawBoardScoreEnemyDiscountRating,
		h.DrawBoardScorePlayerDiscountRating,
		h.LocalBoardWinPlayedMovesDiscountRating,
		h.OverallBoardWinPlayedMovesDiscountRating,
	}

	for i, param := range params {
		bitset := getBits(param)
		for x := 0; x < bitset.GetSize(); x++ {
			b.Set(bbitsPerParamOffset[i]+x, bitset.Get(x))
		}
	}

	return b
}

func GetHeuristic(bits *goga.Bitset) *Game.HeuristicScores {
	if len(bbitsPerParamOffset) == 0 {
		offset := 0
		for _, i2 := range bbitsPerParam {
			bbitsPerParamOffset = append(bbitsPerParamOffset, offset)
			offset += i2
		}
	}

	h := Game.HeuristicScores{
		BoardCornerRating:                        getFloat64(bits.Slice(bbitsPerParamOffset[0], bbitsPerParam[0])),
		BoardSideRating:                          getFloat64(bits.Slice(bbitsPerParamOffset[1], bbitsPerParam[1])),
		BoardMiddleRating:                        getFloat64(bits.Slice(bbitsPerParamOffset[2], bbitsPerParam[2])),
		PosCornerRating:                          getFloat64(bits.Slice(bbitsPerParamOffset[3], bbitsPerParam[3])),
		PosSideRating:                            getFloat64(bits.Slice(bbitsPerParamOffset[4], bbitsPerParam[4])),
		PosMiddleRating:                          getFloat64(bits.Slice(bbitsPerParamOffset[5], bbitsPerParam[5])),
		OverallWinLossRating:                     getFloat64(bits.Slice(bbitsPerParamOffset[6], bbitsPerParam[6])),
		OverallAlmostDrawWinLossRating:           getFloat64(bits.Slice(bbitsPerParamOffset[7], bbitsPerParam[7])),
		GlobalStateRating:                        getFloat64(bits.Slice(bbitsPerParamOffset[8], bbitsPerParam[8])),
		OverallBoardMultiplierRating:             getFloat64(bits.Slice(bbitsPerParamOffset[9], bbitsPerParam[9])),
		WinMovesMadeLossRating:                   getFloat64(bits.Slice(bbitsPerParamOffset[11], bbitsPerParam[11])),
		LossMovesMadeAdvantageRating:             getFloat64(bits.Slice(bbitsPerParamOffset[12], bbitsPerParam[12])),
		TwoInARowAdvantageRating:                 getFloat64(bits.Slice(bbitsPerParamOffset[13], bbitsPerParam[13])),
		EnemyTwoInARowLossRating:                 getFloat64(bits.Slice(bbitsPerParamOffset[14], bbitsPerParam[14])),
		EnemyWonBoardLossRating:                  getFloat64(bits.Slice(bbitsPerParamOffset[15], bbitsPerParam[15])),
		EnemyWonBoardDiscountRating:              getFloat64(bits.Slice(bbitsPerParamOffset[16], bbitsPerParam[16])),
		WonBoardRating:                           getFloat64(bits.Slice(bbitsPerParamOffset[17], bbitsPerParam[17])),
		DrawBoardScoreEnemyDiscountRating:        getFloat64(bits.Slice(bbitsPerParamOffset[18], bbitsPerParam[18])),
		DrawBoardScorePlayerDiscountRating:       getFloat64(bits.Slice(bbitsPerParamOffset[19], bbitsPerParam[19])),
		LocalBoardWinPlayedMovesDiscountRating:   getFloat64(bits.Slice(bbitsPerParamOffset[20], bbitsPerParam[20])),
		OverallBoardWinPlayedMovesDiscountRating: getFloat64(bits.Slice(bbitsPerParamOffset[21], bbitsPerParam[21])),
	}
	h.BoardRating = [9]float64{h.BoardCornerRating, h.BoardSideRating, h.BoardCornerRating, h.BoardSideRating, h.BoardCornerRating, h.BoardSideRating, h.BoardCornerRating, h.BoardSideRating, h.BoardMiddleRating}
	h.PosRating = [9]float64{h.PosCornerRating, h.PosSideRating, h.PosCornerRating, h.PosSideRating, h.PosCornerRating, h.PosSideRating, h.PosCornerRating, h.PosSideRating, h.PosMiddleRating}
	return &h
}

type utttMaterSimulator struct {
	pop    []goga.Genome
	scores sync.Map
}

func (sms *utttMaterSimulator) OnBeginSimulation() {
	sms.pop = genAlgo.GetPopulation()
	sms.scores.Range(func(key, value any) bool {
		sms.scores.Delete(key)
		return true
	})
}
func (sms *utttMaterSimulator) OnEndSimulation() {}
func (sms *utttMaterSimulator) ExitFunc(g goga.Genome) bool {
	return false
}

func (sms *utttMaterSimulator) Simulate(g goga.Genome) {
	playerH := GetHeuristic(g.GetBits())
	winner1, movesMade1 := sms.sim(playerH, oppHeuristic)
	winner2, movesMade2 := sms.sim(oppHeuristic, playerH)

	var fitness uint32 = 0
	if winner1 == Game.Player1 {
		fitness += 400 - movesMade1
	} else {
		fitness += movesMade1
	}

	if winner2 == Game.Player2 {
		fitness += 400 - movesMade2
	} else {
		fitness += movesMade2
	}

	g.SetFitness(int(fitness))
}

func (sms *utttMaterSimulator) sim(p1 *Game.HeuristicScores, p2 *Game.HeuristicScores) (Game.Player, uint32) {
	playerGame := Game.NewGame()
	playerGame.HeuristicScores = p1

	enemyGame := Game.NewGame()
	enemyGame.HeuristicScores = p2

	for !playerGame.IsTerminal() {
		// Player move
		move, board := mtd.IterativeDeepeningTime(playerGame, 5, time.Millisecond*250)
		playerGame.MakeMove(board, move)
		enemyGame.MakeMove(board, move)

		if playerGame.IsTerminal() {
			break
		}

		// Enemy move
		move, board = mtd.IterativeDeepeningTime(enemyGame, 5, time.Millisecond*250)
		playerGame.MakeMove(board, move)
		enemyGame.MakeMove(board, move)
	}

	return playerGame.WinningPlayer(), playerGame.MovesMade()
}

type utttEliteConsumer struct {
	currentIter int
}

func (ec *utttEliteConsumer) OnElite(g goga.Genome) {
	data, _ := json.Marshal(GetHeuristic(g.GetBits()))
	ec.currentIter++
	fmt.Println(ec.currentIter, "\t", g.GetFitness())

	//if g.GetFitness() > 745 {
	//	oppHeuristic = GetHeuristic(g.GetBits())
	//}

	f, _ := os.OpenFile("elitePop.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	_, _ = f.Write(data)
	_, _ = f.WriteString("\n")
	_ = f.Close()
}

func main() {
	numThreads := 16
	runtime.GOMAXPROCS(numThreads)

	genAlgo.Simulator = &utttMaterSimulator{}
	genAlgo.BitsetCreate = &utttHeuristicBitsetCreate{}
	genAlgo.EliteConsumer = &utttEliteConsumer{}
	genAlgo.Mater = goga.NewMater(
		[]goga.MaterFunctionProbability{
			{P: 1.0, F: goga.UniformCrossover, UseElite: true},
			{P: 1.0, F: goga.TwoPointCrossover, UseElite: true},
			{P: 0.8, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 0.7, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
		},
	)
	genAlgo.Selector = goga.NewSelector(
		[]goga.SelectorFunctionProbability{
			{P: 1.0, F: goga.Roulette},
		},
	)

	genAlgo.Init(population, numThreads)

	startTime := time.Now()
	genAlgo.Simulate()
	fmt.Println(time.Since(startTime))
}
