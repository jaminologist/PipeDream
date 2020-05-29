package game

import "github.com/bryjammin/pipedream/lobby/multiplayer/game/model"

func calculateScoreFromBoardReports(boardReports []model.BoardReport) int {

	pipesDestroyed := 0
	for i := 0; i < len(boardReports); i++ {
		pipesDestroyed += len(boardReports[i].DestroyedPipes)
	}

	score := 1250 * pipesDestroyed

	return score
}
