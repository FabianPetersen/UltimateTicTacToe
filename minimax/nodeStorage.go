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
	storage.nodeStore[*node.State.Hash()] = node
}

func (storage *Storage) Reset() {
	storage.nodeStore = make(map[[9]uint32]*Node, 150000)
}

func NewStorage() Storage {
	return Storage{
		nodeStore: make(map[[9]uint32]*Node, 150000),
	}
}
