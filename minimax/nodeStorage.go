package minimax

import "github.com/FabianPetersen/UltimateTicTacToe/Game"

type Storage struct {
	nodeStore map[[9]uint32]*Node
}

func (storage *Storage) Count() int {
	return len(storage.nodeStore)
}

func (storage *Storage) Get(hash Game.GameHash) (*Node, bool) {
	node, exists := storage.nodeStore[*hash]
	return node, exists
}

func (storage *Storage) Set(node *Node) {
	var key = [9]uint32{}
	copy(key[:], (*node.State.Hash())[:])
	storage.nodeStore[key] = node
}

func (storage *Storage) Reset() {
	storage.nodeStore = make(map[[9]uint32]*Node, 50000)
}

func NewStorage() Storage {
	return Storage{
		nodeStore: make(map[[9]uint32]*Node, 50000),
	}
}
