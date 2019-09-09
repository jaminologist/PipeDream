package game

import "math/rand"

//The purpose of this is to solve the

//Function that accepts a board and returns a array of board inputs to reach a solve

func BoardSolve(b *Board) ([]*point, error) {

	return []*point{&point{rand.Intn(7), rand.Intn(7)}, &point{rand.Intn(7), rand.Intn(7)}, &point{rand.Intn(7), rand.Intn(7)}}, nil
}

//Function that takes in a point and returns if there is a solve or not.
func findSolution(b *Board) ([]*point, bool) {

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

func findPathAtPoint(b *Board, x int, y int) ([]*point, bool) {

	visitedMap := map[*point]bool{}
	currentPoint := newPoint(x, y)
	visitedMap[currentPoint] = true

	pipe := b.Cells[x][y]

	pointArray := make([]*point, 0)

	//Rotate up to three times
	maxNumberOfRotates := getMaxNumberOfRotations(pipe)
	for i := 0; i < maxNumberOfRotates; i++ {

		_, ok := getPipesThatAreBeingPointedTo(pipe, b)
		if ok {

		}

		if i == maxNumberOfRotates-1 { //No path found leave
			return nil, false
		}

		pointArray = append(pointArray, &point{x, y})
	}

	//_ := b.Cells[x][y]

	//pipesPointsTo :=

	return nil, false
}

func findPathInChild() {

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

func getPipesThatAreBeingPointedTo(pipe *Pipe, b *Board) ([]*Pipe, bool) {

	pointsTo := pipe.pointsTo()

	pipesPointedTo := make([]*Pipe, 0)

	for i := 0; i < len(pointsTo); i++ {
		if !b.containsPoint(&pointsTo[i]) {
			return nil, false
		}

		pipesPointedTo = append(pipesPointedTo, b.Cells[pointsTo[i].x][pointsTo[i].y])
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
