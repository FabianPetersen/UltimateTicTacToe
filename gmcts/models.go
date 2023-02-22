package gmcts

/*
type gameState struct {
	*Game2.Game
	Game2.GameHash
}

//MCTS contains functionality for the MCTS algorithm
type MCTS struct {
	init  *Game2.Game
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

	nodeScore  map[Game2.Player]float64
	nodeVisits int

	heuristicScore map[Game2.Player]float64
}

//Tree represents a game state tree
type Tree struct {
	current          *node
	gameStates       map[Game2.GameHash]*node
	explorationConst float64
	bestActionPolicy BestActionPolicy
	treePolicy       TreePolicy
}

type BestActionPolicy byte

const (
	MAX_CHILD_SCORE BestActionPolicy = 0
	ROBUST_CHILD    BestActionPolicy = 1
)

type TreePolicy byte

const (
	UCT2      TreePolicy = 0
	SMITSIMAX TreePolicy = 1
)
*/
