package game

//The purpose of this is to solve the

//Function that accepts a board and returns a array of board inputs to reach a solve

func BoardSolve(b *Board) ([]*Point, error) {

	if points, ok := findSolution(b); ok {
		return points, nil
	}

	return []*Point{}, nil
}

type solveBuilder struct {
	pointTapArray []*Point
	visitedPoints map[Point]*Pipe
}

//Function that takes in a Point and returns if there is a solve or not.
func findSolution(b *Board) ([]*Point, bool) {

	///*loop:
	for x := 0; x < len(b.Cells); x++ {
		for y := 0; y < len(b.Cells[x]); y++ {

			solveBuilder, pathFound := findPathAtPoint(b, x, y)

			if pathFound {
				return solveBuilder.pointTapArray, true
			}
		}
	}

	return nil, false

}

func findPathAtPoint(b *Board, x int, y int) (*solveBuilder, bool) {
	pipe := copyPipe(b.Cells[x][y])
	return recursivefindPathAtPoint(map[Point]*Pipe{}, b, nil, pipe)
}

func recursivefindPathAtPoint(originalVistedPoints map[Point]*Pipe, b *Board, parentPipe *Pipe, originalPipe *Pipe) (*solveBuilder, bool) {

	pipe := copyPipe(originalPipe)
	currentPoint := *newPoint(pipe.X, pipe.Y)
	maxNumberOfRotates := getMaxNumberOfRotations(pipe)

	for numberOfRotations := 0; numberOfRotations < maxNumberOfRotates; numberOfRotations++ {

		//When you rotate a pipe to search for a different route
		//You need to reset the visited locations so the new route can re-use those locations as that could lead to a path
		newVisitedPoints := map[Point]*Pipe{
			currentPoint: pipe,
		}

		for k, v := range originalVistedPoints {
			newVisitedPoints[k] = v
		}

		//Check that the pipe is pointing to the parent every rotation.
		//If the pipe does not have a parent assume it is the root of the solve
		if parentPipe == nil || isPipePointingToPipe(pipe, parentPipe) {

			//This method also checks that the pipe is not pointing out of the bounds. (Make method name clearer?)
			if pipes, ok := getPipesThatAreBeingPointedTo(pipe, b); ok {
				pointsArray := make([]*Point, 0)
				pathFound := true
				for _, childPipe := range pipes {

					//check if pipe has been visited
					if solvedPipe, isVisited := newVisitedPoints[Point{childPipe.X, childPipe.Y}]; !isVisited {
						childSolveBuilder, ok := recursivefindPathAtPoint(newVisitedPoints, b, pipe, childPipe)
						if !ok {
							pathFound = false
							break
						} else {
							pointsArray = append(childSolveBuilder.pointTapArray, pointsArray...)
							for k, v := range childSolveBuilder.visitedPoints {
								newVisitedPoints[k] = v
							}
						}
					} else {
						//If the location has been visited before it will be stored in the current 'solve'
						//If the 'solved' copy of the pipe points toward the current pipe it is part of the path.
						if !isPipePointingToPipe(solvedPipe, pipe) {
							pathFound = false
							break
						}
					}
				}

				if pathFound {

					//Prepend point based on the amount of rotations
					for i := 0; i < numberOfRotations; i++ {
						pointsArray = append([]*Point{&Point{X: pipe.X, Y: pipe.Y}}, pointsArray...)
					}

					return &solveBuilder{
						pointTapArray: pointsArray,
						visitedPoints: newVisitedPoints,
					}, true
				}
			}
		}

		pipe.RotateClockWise()
	}

	return nil, false
}

func isPipePointingOutsideOfBoard(pipe *Pipe, b *Board) bool {
	pointsTo := pipe.pointsTo()
	for i := 0; i < len(pointsTo); i++ {
		if !b.containsPoint(&pointsTo[i]) {
			return true
		}
	}
	return false
}

func isPipePointingToPipe(pipe1 *Pipe, pipe2 *Pipe) bool {

	for _, point := range pipe1.pointsTo() {
		if point.X == pipe2.X && point.Y == pipe2.Y {
			return true
		}
	}
	return false
}

func getPipesThatAreBeingPointedTo(pipe *Pipe, b *Board) ([]*Pipe, bool) {

	pointsTo := pipe.pointsTo()

	pipesPointedTo := make([]*Pipe, 0)

	for i := 0; i < len(pointsTo); i++ {
		if !b.containsPoint(&pointsTo[i]) {
			return nil, false
		}

		pipesPointedTo = append(pipesPointedTo, b.Cells[pointsTo[i].X][pointsTo[i].Y])
	}

	return pipesPointedTo, true
}

func getMaxNumberOfRotations(p *Pipe) int {
	switch p.Type {
	case LINE:
		return 2
	default:
		return 4
	}
}

func copyPipe(p *Pipe) *Pipe {
	return &Pipe{
		X:         p.X,
		Y:         p.Y,
		Type:      p.Type,
		Direction: p.Direction,
	}
}
