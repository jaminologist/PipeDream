package model

import "time"

type BlitzGameState struct {
	BoardReports   []BoardReport
	Score          int
	IsOver         bool
	DestroyedPipes []DestroyedPipe
	TimeLimit      *TimeLimit
}

type TimeLimit struct {
	Time time.Duration
}

type VersusPlayerBlitzGamePlayerInformationSentToPlayers struct {
	PlayerID         int
	EnemyInformation *VersusPlayerBlitzGameState
}

type VersusPlayerBlitzGameState struct {
	ID           int
	BoardReports []BoardReport
	Score        int
	IsOver       bool
	IsWinner     bool
}
