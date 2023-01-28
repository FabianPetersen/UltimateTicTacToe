package gmcts

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	//DefaultExplorationConst is the default exploration constant of UCB1 Formula
	//Sqrt(2) is a frequent choice for this constant as specified by
	//https://en.wikipedia.org/wiki/Monte_Carlo_tree_search
	DefaultExplorationConst = math.Sqrt2
	useUCT2                 = false
	greedyRate              = 0
	greedyBot               = false
	botPlayer               = Player1
)

var randSource = rand.New(rand.NewSource(time.Now().Unix()))

func initializeNode(g gameState, tree *Tree) *node {
	return &node{
		state:          g,
		tree:           tree,
		nodeScore:      make(map[Player]float64),
		heuristicScore: g.Heuristic(),
	}
}

//UCT2 algorithm is described in this paper
//https://www.csse.uwa.edu.au/cig08/Proceedings/papers/8057.pdf
func (n *node) UCT2(i int, p Player) float64 {
	exploit := n.children[i].nodeScore[p] / float64(n.children[i].nodeVisits)

	explore := math.Log(float64(n.nodeVisits)) / n.childVisits[i]
	explore = math.Sqrt(explore)

	return exploit + n.tree.explorationConst*explore
}

//smitsimax node selection algorithm is described in this paper
//https://www.codingame.com/playgrounds/36476/smitsimax
func (n *node) smitsimax(i int, p Player) float64 {
	exploit := n.heuristicScore[p] / float64(n.children[i].nodeVisits)

	explore := math.Log(float64(n.nodeVisits)) / n.childVisits[i]
	explore = math.Sqrt(explore)

	return exploit + n.tree.explorationConst*explore
}

func (n *node) treePolicy() int {
	maxScore := -1.0
	thisPlayer := n.state.Player()
	selectedChildIndex := 0
	for i := 0; i < n.actionCount; i++ {
		var score float64
		if useUCT2 {
			score = n.UCT2(i, thisPlayer)
		} else {
			score = n.smitsimax(i, thisPlayer)
		}
		if score > maxScore {
			maxScore = score
			selectedChildIndex = i
		}
	}

	return selectedChildIndex
}

func (n *node) runSimulation() ([]Player, float64) {
	var selectedChildIndex int
	var winners []Player
	var scoreToAdd float64
	var terminalState bool

	//If we have actions, then there's no need to expand.
	if n.actionCount == 0 {
		//If we don't have any actions, then either the state
		//is terminal, or we haven't expanded the node yet.
		terminalState = n.state.IsTerminal()
		if !terminalState {
			n.expand()
		}
	}

	if terminalState {
		//Get the result of the game
		winners = n.simulate()
		scoreToAdd = 1.0 / float64(len(winners))

	} else if n.unvisitedChildren > 0 {
		//Grab the first unvisited child and run a simulation from that point
		selectedChildIndex = n.actionCount - n.unvisitedChildren
		n.children[selectedChildIndex].nodeVisits++
		n.unvisitedChildren -= 1

		winners = n.children[selectedChildIndex].simulate()
		scoreToAdd = 1.0 / float64(len(winners))

	} else {
		//Select the child with the max UCT2 score with the current player
		//and get the results to add from its selection
		selectedChildIndex = n.treePolicy()
		winners, scoreToAdd = n.children[selectedChildIndex].runSimulation()
	}

	//Update this node along with each parent in this path recursively
	n.nodeVisits++
	if n.actionCount != 0 {
		n.childVisits[selectedChildIndex]++
	}

	for _, p := range winners {
		n.nodeScore[p] += scoreToAdd
	}
	return winners, scoreToAdd
}

func (n *node) expand() {
	n.actionCount = n.state.Len()
	n.unvisitedChildren = n.actionCount
	n.children = make([]*node, n.actionCount)
	n.childVisits = make([]float64, n.actionCount)
	for i := 0; i < n.actionCount; i++ {
		newGame, err := n.state.ApplyAction(i)
		if err != nil {
			panic(fmt.Sprintf("gmcts: Game returned an error when exploring the tree: %s", err))
		}

		newState := gameState{newGame, newGame.Hash()}

		//If we already have a copy in cache, use that and update
		//this node and its parents
		if cachedNode, made := n.tree.gameStates[newState.GameHash]; made {
			n.children[i] = cachedNode
		} else {
			newNode := initializeNode(newState, n.tree)
			n.children[i] = newNode

			//Save node for reuse
			n.tree.gameStates[newState.GameHash] = newNode
		}
	}
}

func (n *node) simulate() []Player {
	game := n.state.Game.Copy()
	i := 0
	move := 0
	for !game.IsTerminal() {
		// Greedy for first x moves?
		if (greedyBot && game.currentPlayer != botPlayer) || (!greedyBot && i < greedyRate) {
			move = game.GreedyMove()
		} else {
			move = randSource.Intn(game.Len())
		}
		game.ApplyActionModify(move)
		i++
	}
	return game.Winners()
}
