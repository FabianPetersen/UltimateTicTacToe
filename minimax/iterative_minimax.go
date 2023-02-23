package minimax

import (
	"fmt"
	"github.com/FabianPetersen/UltimateTicTacToe/Game"
	"math"
)

/*

def negamax_nr(game, target_depth, scoring, alpha=-INF, beta=+INF):

    ################################################
    #
    #    INITIALIZE AND CHECK ENTRY CONDITIONS
    #
    ################################################

    if not hasattr(game, "ttentry"):
        raise AttributeError('Method "ttentry()" missing from game.')
    if not hasattr(game, "ttrestore"):
        raise AttributeError('Method "ttrestore()" missing from game.')

    if game.is_over():
        score = scoring(game)
        game.ai_move = None
        return score

    if target_depth == 0:
        current_game = game.ttentry()
        move_list = game.possible_moves()
        best_move = None
        best_score = -INF
        for move in move_list:
            game.make_move(move)
            score = scoring(game)
            if score > best_score:
                best_move = copy.copy(move)
                best_score = score
            game.ttrestore(current_game)
        game.ai_move = best_move
        return best_score

    states = StateList(target_depth)

    ################################################
    #
    #    START GRAND LOOP
    #
    ################################################

    depth = -1  # proto-parent
    states[depth].alpha = alpha
    states[depth].beta = beta
    direction = DOWN
    depth = 0

    while True:
        parent = depth - 1
        if direction == DOWN:
            if (depth < target_depth) and not game.is_over():  # down we go...
                states[depth].image = game.ttentry()
                states[depth].move_list = game.possible_moves()
                states[depth].best_move = 0
                states[depth].best_score = -INF
                states[depth].current_move = 0
                states[depth].player = game.current_player
                states[depth].alpha = -states[parent].beta  # inherit alpha from -beta
                states[depth].beta = -states[parent].alpha  # inherit beta from -alpha
                index = states[depth].current_move
                game.make_move(states[depth].move_list[index])
                game.switch_player()
                direction = DOWN
                depth += 1
            else:  # reached a leaf or the game is over; going back up
                leaf_score = -scoring(game)
                if leaf_score > states[parent].best_score:
                    states[parent].best_score = leaf_score
                    states[parent].best_move = states[parent].current_move
                if states[parent].alpha < leaf_score:
                    states[parent].alpha = leaf_score
                direction = UP
                depth = parent
            continue
        elif direction == UP:
            prune_time = states[depth].alpha >= states[depth].beta
            if states[depth].out_of_moves() or prune_time:  # out of moves
                bs = -states[depth].best_score
                if bs > states[parent].best_score:
                    states[parent].best_score = bs
                    states[parent].best_move = states[parent].current_move
                if states[parent].alpha < bs:
                    states[parent].alpha = bs
                if depth <= 0:
                    break  # we are done.
                direction = UP
                depth = parent
                continue
            # else go down the next branch
            game.ttrestore(states[depth].image)
            game.current_player = states[depth].player
            next_move = states[depth].goto_next_move()
            game.make_move(next_move)
            game.switch_player()
            direction = DOWN
            depth += 1

    best_move_index = states[0].best_move
    best_move = states[0].move_list[best_move_index]
    best_value = states[0].best_score
    game.ai_move = best_move
    return best_value
*/

type StateHolder struct {
	Alpha       float64
	Beta        float64
	BestMove    int
	MoveLength  int
	BestValue   float64
	CurrentMove int
	Player      Game.Player
	GameCopy    Game.Game
}

func (state *StateHolder) OutOfMoves() bool {
	return state.CurrentMove >= state.GameCopy.Len()-1
}

func (state *StateHolder) NextMove() (int, error) {
	state.CurrentMove += 1
	return state.GameCopy.GetMove(state.CurrentMove)
}

const UP = 1
const DOWN = 0

func (n *Node) SearchIterative(alpha float64, beta float64, targetDepth byte, maxPlayer Game.Player) (float64, int) {
	if targetDepth < 2 || n.State.IsTerminal() {
		return n.State.HeuristicPlayer(maxPlayer), n.bestMove
	}

	states := make([]StateHolder, targetDepth)
	depth := byte(0)
	game := n.State.Copy()
	states[depth].GameCopy = game
	states[depth].Alpha = alpha
	states[depth].Beta = beta
	direction := DOWN
	depth = 1

	for true {
		parent := depth - 1
		if direction == DOWN {
			if depth < targetDepth-1 && !game.IsTerminal() {
				states[depth].GameCopy = game.Copy()
				states[depth].BestMove = 0
				states[depth].BestValue = math.Inf(-1)
				states[depth].CurrentMove = 0
				states[depth].Player = game.CurrentPlayer
				states[depth].Alpha = -states[parent].Beta // inherit alpha from -beta
				states[depth].Beta = -states[parent].Alpha // inherit beta from -alpha

				index := states[depth].CurrentMove
				game.ApplyActionModify(index)
				direction = DOWN
				depth += 1
			} else { // reached a leaf or the game is over; going back up
				leafValue := -game.HeuristicPlayer(maxPlayer)
				if game.CurrentPlayer != maxPlayer {
					leafValue = -leafValue
				}
				if leafValue > states[parent].BestValue {
					states[parent].BestValue = leafValue
					states[parent].BestMove = states[parent].CurrentMove
				}
				if states[parent].Alpha < leafValue {
					states[parent].Alpha = leafValue
				}
				direction = UP
				depth = parent
			}
			continue
		} else if direction == UP {
			shouldPrune := states[depth].Alpha >= states[depth].Beta
			if states[depth].OutOfMoves() || shouldPrune {
				bestValue := -states[depth].BestValue
				if bestValue > states[parent].BestValue {
					states[parent].BestValue = bestValue
					states[parent].BestMove = states[parent].CurrentMove
				}
				if states[parent].Alpha < bestValue {
					states[parent].Alpha = bestValue
				}
				if depth <= 1 { // We are done TODO: should maybe be 0
					break
				}
				direction = UP
				depth = parent
				continue
			} else { // or go down the next branch
				game = states[depth].GameCopy.Copy()
				move, _ := states[depth].NextMove()
				game.ApplyActionModify(move)
				direction = DOWN
				depth += 1
			}
		}
	}

	return states[1].BestValue, states[1].BestMove
}

func (minimax *Minimax) SearchIterative() int {
	if minimax.Depth == 0 {
		minimax.setDepth()
	}

	TranspositionTable.Reset()
	_, bestMove := minimax.root.SearchIterative(math.Inf(-1), math.Inf(1), minimax.Depth, Game.Player2)
	fmt.Printf("Stored nodes, %d %d, Depth %d \n", len(TranspositionTable.nodeStore), minimax.Depth)
	return bestMove
}
