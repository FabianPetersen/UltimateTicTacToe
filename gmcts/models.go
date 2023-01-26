package gmcts

import (
	"math/rand"
	"sync"
)

//Player is an id for the player
type Player byte

type gameState struct {
	*Game
	GameHash
}

type GameHash uint64

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
	unvisitedChildren []*node
	childVisits       []float64
	actionCount       int

	nodeScore  map[Player]float64
	nodeVisits int
}

//Tree represents a game state tree
type Tree struct {
	current          *node
	gameStates       map[GameHash]*node
	explorationConst float64
	randSource       *rand.Rand
}
