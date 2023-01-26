package main

import (
	"github.com/FabianPetersen/UltimateTicTacToe/gmcts"
	"testing"
)

func TestSpeed(t *testing.T) {
	game := gmcts.NewGame()
	mcts := gmcts.NewMCTS(game)
	tree := mcts.SpawnTree()
	tree.SearchRounds(1500000)
	mcts.AddTree(tree)
}
