package gmcts

import (
	"math"
	"unsafe"
)

type Node struct {
	parent        *Node
	children      []*Node
	childrenCount byte
	maxChildren   byte

	move  byte
	board byte

	nodeScore   uint16
	nodeVisits  uint16
	nodeExploit float32
}

type BestActionPolicy byte

const (
	MAX_CHILD_SCORE BestActionPolicy = 0
	ROBUST_CHILD    BestActionPolicy = 1
)

const (
	//DefaultExplorationConst is the default exploration constant of UCT2 Formula
	//Sqrt(2) is a frequent choice for this constant as specified by
	//https://en.wikipedia.org/wiki/Monte_Carlo_tree_search
	DefaultExplorationConst = float32(math.Sqrt2) - 1
)

const magic32 = 0x5F375A86
const th = 1.5

func FastSqrt32(n float32) float32 {
	b := *(*uint64)(unsafe.Pointer(&n))
	b = magic32 - (b >> 1)
	y := *(*float32)(unsafe.Pointer(&b))
	y *= th - (n * 0.5 * y * y)
	return 1 / y
}

func ln(x float32) float32 {
	var bx = *(*uint32)(unsafe.Pointer(&x))
	var ex = bx >> 23
	var t = float32(ex) - 127
	bx = 1065353216 | (bx & 8388607)
	x = *(*float32)(unsafe.Pointer(&bx))
	return -1.49278 + (2.11263+(-0.729104+0.10969*x)*x)*x + 0.6931471806*t
}

// UCT2 algorithm is described in this paper
// https://www.csse.uwa.edu.au/cig08/Proceedings/papers/8057.pdf
func (n *Node) UCT2(i byte) float32 {
	explore := ln(float32(n.nodeVisits)) / float32(n.children[i].nodeVisits) // math.Log(float64(n.nodeVisits)) / float64(n.children[i].nodeVisits)
	explore = float32(math.Sqrt(float64(explore)))                           // FastSqrt32(explore)                                            // float32(math.Sqrt(float64(explore)))                           // 1 / FastInvSqrt64(explore) // math.Sqrt(explore) // 1 / FastInvSqrt64(explore) //

	return n.nodeExploit + DefaultExplorationConst*explore
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

func (node *Node) treePolicy() *Node {
	var bestScore float32 = 0
	var bestNode = node.children[0]
	for i := byte(0); i < node.childrenCount; i++ {
		score := node.UCT2(i)
		if score >= bestScore {
			bestScore = score
			bestNode = node.children[i]
		}
	}
	return bestNode
}
