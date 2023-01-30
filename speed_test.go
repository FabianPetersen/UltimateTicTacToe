package main

import (
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"github.com/FabianPetersen/UltimateTicTacToe/gmcts"
	"testing"
	"time"
)

func TestSpeed(t *testing.T) {
	game := Game.NewGame()
	mcts := gmcts.NewMCTS(game)
	tree := mcts.SpawnTree(gmcts.ROBUST_CHILD, gmcts.SMITSIMAX)
	tree.SearchTime(200 * time.Millisecond)
	mcts.AddTree(tree)
}
