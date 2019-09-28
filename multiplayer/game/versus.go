package game

import (
	"time"

	"bryjamin.com/multiplayer/player"
	"bryjamin.com/multiplayer/send"
)

type VersusPlayerBlitzGame struct {
	playerGameInformation map[*player.Player](*VersusPlayerBlitzGamePlayerInformation)
	timeLimit             time.Duration
	isOver                bool

	sendMessageToPlayerCh      chan *player.PlayerMessage
	receiveMessageFromPlayerCh chan *player.PlayerMessage

	playerInputChannel   chan *player.PlayerBoardInput
	gameOverInputChannel chan bool
}

type VersusPlayerBlitzGamePlayerInformation struct {
	ID       int
	Board    *Board
	Score    int
	IsWinner bool
}

type VersusPlayerBlitzGamePlayerInformationSentToPlayers struct {
	PlayerID         int
	EnemyInformation *VersusPlayerBlitzGameState
}

type VersusPlayerBlitzGameState struct {
	ID int

	Board        *Board
	BoardReports []BoardReport
	Score        int
	IsOver       bool
	IsWinner     bool
}

func NewVersusPlayerBlitzGame(timeLimit time.Duration, players []*player.Player, sendMessageToPlayerCh chan *player.PlayerMessage, receiveMessageFromPlayerCh chan *player.PlayerMessage) *VersusPlayerBlitzGame {

	playerGameInformation := make(map[*player.Player](*VersusPlayerBlitzGamePlayerInformation))

	i := 0
	for _, player := range players {
		newBoard := NewBoard(7, 7)
		newBoard.UpdateBoardPipeConnections() //Note: Need to add a way to generate a board where there are no connections straight away.
		playerGameInformation[player] = &VersusPlayerBlitzGamePlayerInformation{
			i,
			&newBoard,
			0,
			false,
		}
		i++
	}

	return &VersusPlayerBlitzGame{
		playerGameInformation:      playerGameInformation,
		timeLimit:                  timeLimit,
		sendMessageToPlayerCh:      sendMessageToPlayerCh,
		receiveMessageFromPlayerCh: receiveMessageFromPlayerCh,
		playerInputChannel:         make(chan *player.PlayerBoardInput),
		gameOverInputChannel:       make(chan bool),
	}
}

func (vpbg *VersusPlayerBlitzGame) Run() {

	go func() {

		for player, info := range vpbg.playerGameInformation {
			send.SendMessageToPlayer(&BlitzGameState{
				Board: info.Board,
				Score: info.Score,
			}, player, vpbg.sendMessageToPlayerCh)

			opponent := vpbg.getOpponent(player)

			send.SendMessageToPlayer(&VersusPlayerBlitzGamePlayerInformationSentToPlayers{
				EnemyInformation: &VersusPlayerBlitzGameState{
					Board: info.Board,
					Score: info.Score,
				},
			}, opponent, vpbg.sendMessageToPlayerCh)

		}

		for !vpbg.isOver {
			vpbg.timeLimit = vpbg.timeLimit - serverTick
			for player := range vpbg.playerGameInformation {
				go send.SendMessageToPlayer(&TimeLimit{
					Time: vpbg.timeLimit,
				}, player, vpbg.sendMessageToPlayerCh)
			}

			vpbg.isOver = vpbg.timeLimit <= 0
			time.Sleep(serverTick)
		}

		vpbg.gameOverInputChannel <- vpbg.isOver
	}()

OuterLoop:
	for {
		select {
		case isOver := <-vpbg.gameOverInputChannel:
			if isOver {

				var winner *player.Player
				winnerScore := -1
				for player, info := range vpbg.playerGameInformation {
					if info.Score > winnerScore {
						winner = player
						winnerScore = info.Score
					}
				}

				vpbg.playerGameInformation[winner].IsWinner = true

				for player, info := range vpbg.playerGameInformation {
					send.SendMessageToPlayer(&VersusPlayerBlitzGameState{
						Board:    info.Board,
						IsOver:   vpbg.isOver,
						Score:    info.Score,
						IsWinner: info.IsWinner,
					}, player, vpbg.sendMessageToPlayerCh)
				}
				break OuterLoop
			}
		case playerBoardInput := <-vpbg.playerInputChannel:

			player := playerBoardInput.Player
			info := vpbg.playerGameInformation[player]
			info.Board.RotatePipeClockwise(playerBoardInput.X, playerBoardInput.Y)

			boardReports := info.Board.UpdateBoardPipeConnections()

			info.Score += calculateScoreFromBoardReports(boardReports)

			gameState := BlitzGameState{
				BoardReports: boardReports,
				Score:        info.Score,
				IsOver:       vpbg.isOver,
			}

			send.SendMessageToPlayer(gameState, player, vpbg.sendMessageToPlayerCh)

			opponent := vpbg.getOpponent(player)
			gameStateSentToOpponent := &VersusPlayerBlitzGamePlayerInformationSentToPlayers{
				EnemyInformation: &VersusPlayerBlitzGameState{
					BoardReports: boardReports,
					Score:        info.Score,
				},
			}

			send.SendMessageToPlayer(gameStateSentToOpponent, opponent, vpbg.sendMessageToPlayerCh)
		}
	}

}

func (vpbg *VersusPlayerBlitzGame) getOpponent(p *player.Player) *player.Player {
	for opponent := range vpbg.playerGameInformation {
		if p != opponent {
			return opponent
		}
	}
	return nil
}

func (vpbg *VersusPlayerBlitzGame) SendPlayerBoardInputToGame(pbi *player.PlayerBoardInput) {
	vpbg.playerInputChannel <- pbi
}
