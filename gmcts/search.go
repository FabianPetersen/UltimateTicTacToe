package gmcts

import (
	"fmt"
	"math"
)

const (
	//DefaultExplorationConst is the default exploration constant of UCB1 Formula
	//Sqrt(2) is a frequent choice for this constant as specified by
	//https://en.wikipedia.org/wiki/Monte_Carlo_tree_search
	DefaultExplorationConst = math.Sqrt2
)

func initializeNode(g gameState, tree *Tree) *node {
	return &node{
		state:     g,
		tree:      tree,
		nodeScore: make(map[Player]float64),
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
	} else if len(n.unvisitedChildren) > 0 {
		//Grab the first unvisited child and run a simulation from that point
		selectedChildIndex = n.actionCount - len(n.unvisitedChildren)
		n.children[selectedChildIndex].nodeVisits++
		n.unvisitedChildren = n.unvisitedChildren[1:]

		winners = n.children[selectedChildIndex].simulate()
		scoreToAdd = 1.0 / float64(len(winners))
	} else {
		//Select the child with the max UCT2 score with the current player
		//and get the results to add from its selection
		maxScore := -1.0
		thisPlayer := n.state.Player()
		for i := 0; i < n.actionCount; i++ {
			score := n.UCT2(i, thisPlayer)
			if score > maxScore {
				maxScore = score
				selectedChildIndex = i
			}
		}
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
	n.unvisitedChildren = make([]*node, n.actionCount)
	n.children = n.unvisitedChildren
	n.childVisits = make([]float64, n.actionCount)
	for i := 0; i < n.actionCount; i++ {
		newGame, err := n.state.ApplyAction(i)
		if err != nil {
			panic(fmt.Sprintf("gmcts: Game returned an error when exploring the tree: %s", err))
		}

		newState := gameState{newGame, gameHash{newGame.Hash(), n.state.turn + 1}}

		//If we already have a copy in cache, use that and update
		//this node and its parents
		if cachedNode, made := n.tree.gameStates[newState.gameHash]; made {
			n.unvisitedChildren[i] = cachedNode
		} else {
			newNode := initializeNode(newState, n.tree)
			n.unvisitedChildren[i] = newNode

			//Save node for reuse
			n.tree.gameStates[newState.gameHash] = newNode
		}
	}
}

func (n *node) simulate() []Player {
	game := n.state.Game
	for !game.IsTerminal() {
		var err error

		actions := game.Len()
		if actions <= 0 {
			fmt.Println(!game.IsTerminal())
			fmt.Println(game.Len())
			panic(fmt.Sprintf("gmcts: game returned no actions on a non-terminal state: %#v", game))
		}

		randomIndex := n.tree.randSource.Intn(actions)
		game, err = game.ApplyAction(randomIndex)
		if err != nil {
			panic(fmt.Sprintf("gmcts: game returned an error while searching the tree: %s", err))
		}
	}
	return game.Winners()
}
