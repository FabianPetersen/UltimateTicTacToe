package gmcts

import (
	"sync"
)

//Player is an id for the player
type Player byte

type gameState struct {
	*Game
	GameHash
}

type GameHash interface{}

//MCTS contains functionality for the MCTS algorithm
type MCTS struct {
	init  *Game
	trees []*Tree
	mutex *sync.RWMutex
	seed  int64
}

type node struct {
	state gameState
	tree  *Tree

	children          []*node
	unvisitedChildren int
	childVisits       []float64
	actionCount       int

	nodeScore  map[Player]float64
	nodeVisits int

	heuristicScore map[Player]float64
}

//Tree represents a game state tree
type Tree struct {
	current          *node
	gameStates       map[GameHash]*node
	explorationConst float64
	bestActionPolicy BestActionPolicy
}

type BestActionPolicy byte

const (
	MAX_CHILD_SCORE BestActionPolicy = 0
	ROBUST_CHILD    BestActionPolicy = 1
)
