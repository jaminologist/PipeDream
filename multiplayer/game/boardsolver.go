package game

//The purpose of this is to solve the

//Function that accepts a board and returns a array of board inputs to reach a solve

func BoardSolve(b *Board) ([]*Point, error) {
	return []*Point{&Point{0, 1}, &Point{5, 1}}, nil
}

//Function that takes in a Point and returns if there is a solve or not.
func findSolution(b *Board) ([]*Point, bool) {

	/*loop:
	for x := 0; x < len(b.Cells); x++ {
		for y := 0; y < len(b.Cells[x]); y++ {

			pointArray, pathFound := findPathAtPoint(b, x, y)

			if true {

				break loop
			}
		}
	}*/

	return nil, false

}

func findPathAtPoint(b *Board, x int, y int) ([]*Point, bool) {

	//Add current square to visited points
	visitedMap := map[*Point]bool{}
	currentPoint := newPoint(x, y)
	visitedMap[currentPoint] = true

	pipe := b.Cells[x][y]
	pointArray := make([]*Point, 0)

	//Number of times you check for a path
	maxNumberOfRotates := getMaxNumberOfRotations(pipe)
	for i := 0; i < maxNumberOfRotates; i++ {

		pipes, ok := getPipesThatAreBeingPointedTo(pipe, b)
		if ok {
			for _, childPipe := range pipes {
				_ = childPipe
			}
		}

		if i == maxNumberOfRotates-1 { //No path found leave
			return nil, false
		}

		pointArray = append(pointArray, &Point{x, y})
	}

	//_ := b.Cells[x][y]

	//pipesPointsTo :=

	return nil, false
}

func findPathInChild(vistedMap map[*Point]bool, b *Board, parentPipe *Pipe, childPipe *Pipe) {

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
