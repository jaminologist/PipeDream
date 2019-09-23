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

	//Add current square to visited points
	currentPoint := *newPoint(x, y)

	pipe := copyPipe(b.Cells[x][y])

	//Number of times you check for a path
	maxNumberOfRotates := getMaxNumberOfRotations(pipe)
	for i := 0; i < maxNumberOfRotates; i++ {

		visitedMap := map[Point]*Pipe{}
		visitedMap[currentPoint] = pipe

		pipes, ok := getPipesThatAreBeingPointedTo(pipe, b)
		if ok {
			pointsArray := make([]*Point, 0)
			pathFound := true
			for _, childPipe := range pipes {

				//check if pipe has been visited
				if solvedPipe, isVisited := visitedMap[Point{childPipe.X, childPipe.Y}]; !isVisited {
					childSolveBuilder, ok := findPathInChild(visitedMap, b, pipe, childPipe)

					if !ok {
						pathFound = false
						break
					} else {
						//prepend
						pointsArray = append(childSolveBuilder.pointTapArray, pointsArray...)
						for k, v := range childSolveBuilder.visitedPoints {
							visitedMap[k] = v
						}
						/*for _, point := range childSolveBuilder.pointTapArray {
							newVisitedMap[*point] = childSolveBuilder.visitedPoints[*point]
						}*/
					}
				} else { //ADD DIRECTION OF VISITED PIPE TO MAP AND USE THAT TO SEE IF THE VISITED POINT IS POINTIUNG TO THE PIPE
					//fmt.Println("isVisited else is visited:", pipe, ":", childPipe.X, ", ", childPipe.Y)
					if !isPipePointingToPipe(solvedPipe, pipe) {
						pathFound = false
						break
					}
				}
			}

			if pathFound {

				//Add number of points based on number of rotations
				for j := 0; j < i; j++ {
					pointsArray = append([]*Point{&Point{X: pipe.X, Y: pipe.Y}}, pointsArray...)
				}

				return &solveBuilder{
					pointTapArray: pointsArray,
					visitedPoints: visitedMap,
				}, true
			}
		}
		pipe.RotateClockWise()
	}
	return nil, false
}

func findPathInChild(originalVistedMap map[Point]*Pipe, b *Board, parentPipe *Pipe, originalPipe *Pipe) (*solveBuilder, bool) {

	pipe := copyPipe(originalPipe)
	currentPoint := *newPoint(pipe.X, pipe.Y)
	maxNumberOfRotates := getMaxNumberOfRotations(pipe)

	for i := 0; i < maxNumberOfRotates; i++ {

		//reset visited points for new pipe rotation
		newVisitedMap := map[Point]*Pipe{}

		for k, v := range originalVistedMap {
			newVisitedMap[k] = v
		}

		newVisitedMap[currentPoint] = pipe

		if isPipePointingToPipe(pipe, parentPipe) && !isPipePointingOutsideOfBoard(pipe, b) {

			if pipes, ok := getPipesThatAreBeingPointedTo(pipe, b); ok {
				pointsArray := make([]*Point, 0)
				pathFound := true
				for _, childPipe := range pipes {

					//check if pipe has been visited
					if solvedPipe, isVisited := newVisitedMap[Point{childPipe.X, childPipe.Y}]; !isVisited {
						childSolveBuilder, ok := findPathInChild(newVisitedMap, b, pipe, childPipe)

						//prepend

						if !ok {
							pathFound = false
							break
						} else {
							pointsArray = append(childSolveBuilder.pointTapArray, pointsArray...)
							for k, v := range childSolveBuilder.visitedPoints {
								newVisitedMap[k] = v
							}
							/*for _, point := range childSolveBuilder.pointTapArray {
								newVisitedMap[*point] = childSolveBuilder.visitedPoints[*point]
							}*/
						}
					} else { //ADD DIRECTION OF VISITED PIPE TO MAP AND USE THAT TO SEE IF THE VISITED POINT IS POINTIUNG TO THE PIPE
						//fmt.Println("isVisited else is visited:", pipe, ":", childPipe.X, ", ", childPipe.Y)
						if !isPipePointingToPipe(solvedPipe, pipe) {
							pathFound = false
							break
						}
					}
				}

				if pathFound {

					//Add number of points based on number of rotations
					for j := 0; j < i; j++ {
						pointsArray = append([]*Point{&Point{X: pipe.X, Y: pipe.Y}}, pointsArray...)
					}

					return &solveBuilder{
						pointTapArray: pointsArray,
						visitedPoints: newVisitedMap,
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
