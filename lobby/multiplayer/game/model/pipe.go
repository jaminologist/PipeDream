package model

import "time"

//Pipe Represents a pipe within the game has a Type, Direction and 'Level' and an X and Y position
type Pipe struct {
	Type      PipeType
	Direction PipeDirection
	Level     PipeLevel
	X         int
	Y         int
}

//PipeType the types of pipe that exist within the game.
type PipeType int

//Collection of all pipe types in the game
const (
	NONE          PipeType = -1
	LINE          PipeType = 0
	LPIPE         PipeType = 2
	END           PipeType = 4
	ENDEXPLOSION2 PipeType = 8
	ENDEXPLOSION3 PipeType = 16
)

//PipeDirection The Direction the pipe is facing
type PipeDirection int

//Collection of pipe directions set using Dir
const (
	UP    PipeDirection = 0
	RIGHT PipeDirection = 90
	DOWN  PipeDirection = 180
	LEFT  PipeDirection = 270
)

var pipeDirections = []PipeDirection{UP, RIGHT, DOWN, LEFT}

//PipeLevel Used to display the the level of the pipes connected to this pipe.
type PipeLevel int

const (
	level0 PipeLevel = 0
	level1 PipeLevel = 1
	level2 PipeLevel = 2
	level3 PipeLevel = 3
)

//PipeMovementAnimation Used to detail information for pipe 'Falling' animation. From StartY to EndY and expected travel time
type PipeMovementAnimation struct {
	X          int
	StartY     int
	EndY       int
	TravelTime time.Duration
}

//DestroyedPipe Used to detail information for Pipe 'Destroyed' animation using X and Y position and pipe Type
type DestroyedPipe struct {
	Type PipeType

	X int
	Y int
}
