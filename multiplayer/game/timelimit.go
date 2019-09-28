package game

import (
	"time"
)

//TimeLimit used to pass down the remaining time to the game client
type TimeLimit struct {
	Time time.Duration
}
