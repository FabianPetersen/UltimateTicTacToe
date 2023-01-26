package gmcts

import (
	"errors"
	"math/rand"
	"sync"
)

var (
	//ErrNoTrees notifies the callee that the MCTS wrapper has recieved to trees to analyze
	ErrNoTrees = errors.New("gmcts: mcts wrapper has collected to trees to analyze")

	//ErrTerminal notifies the callee that the given state is terminal
	ErrTerminal = errors.New("gmcts: given game state is a terminal state, therefore, it cannot return an action")

	//ErrNoActions notifies the callee that the given state has <= 0 actions
	ErrNoActions = errors.New("gmcts: given game state is not terminal, yet the state has <= 0 actions to search through")
)

//NewMCTS returns a new MCTS wrapper
func NewMCTS(initial *Game) *MCTS {
	return &MCTS{
		init:  initial,
		trees: make([]*Tree, 0),
		mutex: new(sync.RWMutex),
	}
}

//SpawnTree creates a new search tree. The tree returned uses Sqrt(2) as the
//exploration constant.
func (m *MCTS) SpawnTree() *Tree {
	return m.SpawnCustomTree(DefaultExplorationConst)
}

//SetSeed sets the seed of the next tree to be spawned.
//This value is initially set to 0, and increments on each
//spawned tree.
func (m *MCTS) SetSeed(seed int64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.seed = seed
}

//SpawnCustomTree creates a new search tree with a given exploration constant.
func (m *MCTS) SpawnCustomTree(explorationConst float64) *Tree {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	t := &Tree{
		gameStates:       make(map[GameHash]*node),
		explorationConst: explorationConst,
		randSource:       rand.New(rand.NewSource(m.seed)),
	}
	t.current = initializeNode(gameState{m.init, m.init.Hash()}, t)

	m.seed++
	return t
}

//AddTree adds a searched tree to its list of trees to consider
//when deciding upon an action to take.
func (m *MCTS) AddTree(t *Tree) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.trees = append(m.trees, t)
}

//BestAction takes all of the searched trees and returns
//the index of the best action based on the highest win
//percentage of each action.
//
//BestAction returns ErrNoTrees if it has received no trees
//to search through, ErrNoActions if the current state
//it's considering has no legal actions, or ErrTerminal
//if the current state it's considering is terminal.
func (m *MCTS) BestAction() (int, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	//Error checking
	if len(m.trees) == 0 {
		return -1, ErrNoTrees
	} else if m.init.IsTerminal() {
		return -1, ErrTerminal
	} else if m.init.Len() <= 0 {
		return -1, ErrNoActions
	}

	//Democracy Section: each tree votes for an action
	actionScore := make([]int, m.init.Len())
	for _, t := range m.trees {
		actionScore[t.bestAction()]++
	}

	//Democracy Section: the action with the most votes wins
	var bestAction int
	var mostVotes int
	for a, s := range actionScore {
		if s > mostVotes {
			bestAction = a
			mostVotes = s
		}
	}
	return bestAction, nil
}
