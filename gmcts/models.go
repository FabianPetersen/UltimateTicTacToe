package gmcts

import (
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"math"
)

type Node struct {
	parent   *Node
	children []*Node

	move  byte
	board byte

	nodeScore  [2]float64
	nodeVisits int
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

const (
	//DefaultExplorationConst is the default exploration constant of UCT2 Formula
	//Sqrt(2) is a frequent choice for this constant as specified by
	//https://en.wikipedia.org/wiki/Monte_Carlo_tree_search
	DefaultExplorationConst = math.Sqrt2
)

// UCT2 algorithm is described in this paper
// https://www.csse.uwa.edu.au/cig08/Proceedings/papers/8057.pdf
func (n *Node) UCT2(i int, p *Game.Player) float64 {
	exploit := n.children[i].nodeScore[*p] / float64(n.children[i].nodeVisits)

	explore := math.Log(float64(n.nodeVisits)) / float64(n.children[i].nodeVisits)
	explore = math.Sqrt(explore)

	return exploit + DefaultExplorationConst*explore
}

// smitsimax Node selection algorithm is described in this paper
// https://www.codingame.com/playgrounds/36476/smitsimax
/*
func (n *Node) smitsimax(i int, p Game.Player) float64 {
	exploit := 0.3 * n.children[i].nodeScore[p] / float64(n.children[i].nodeVisits)
	exploit += 0.7 * n.children[i].heuristicScore[p] / float64(n.children[i].nodeVisits)

	explore := math.Log(float64(n.nodeVisits)) / n.childVisits[i]
	explore = math.Sqrt(explore)

	return exploit + DefaultExplorationConst*explore
}
*/

func (node *Node) treePolicy(player *Game.Player) *Node {
	var bestScore float64 = 0
	var bestNode *Node = nil
	for i := 0; i < len(node.children); i++ {
		score := node.UCT2(i, player)
		if score >= bestScore {
			bestScore = score
			bestNode = node.children[i]
		}
	}
	return bestNode
}
