package ggol

import (
	"sync"
	"testing"
)

func shouldInitializeGameWithCorrectSize(t *testing.T) {
	width := 30
	height := 10
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initialTestCell, defaultCellIteratorForTest)
	cellLiveMap := *convertTestCellsMatricToAliveTestCellsMap(g.GetGeneration())

	if len(cellLiveMap) == width && len(cellLiveMap[0]) == height {
		t.Log("Passed")
	} else {
		t.Fatalf("Size should be %v x %v", width, height)
	}
}

func shouldThrowErrorWhenSizeIsInvalid(t *testing.T) {
	width := -1
	height := 3
	size := Size{Width: width, Height: height}
	_, err := NewGame(&size, initialTestCell, defaultCellIteratorForTest)

	if err == nil {
		t.Fatalf("Should get error when giving invalid size.")
	}
	t.Log("Passed")
}

func TestNewGame(t *testing.T) {
	shouldInitializeGameWithCorrectSize(t)
	shouldThrowErrorWhenSizeIsInvalid(t)
}

func shouldThrowErrorWhenCellSeedExceedBoarder(t *testing.T) {
	width := 2
	height := 2
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initialTestCell, defaultCellIteratorForTest)
	c := Coordinate{X: 0, Y: 10}
	err := g.SetCell(&c, TestCell{Alive: true})

	if err == nil {
		t.Fatalf("Should get error when any seed units are outside border.")
	}
	t.Log("Passed")
}

func shouldSetCellCorrectly(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	c := Coordinate{X: 1, Y: 1}
	g, _ := NewGame(&size, initialTestCell, defaultCellIteratorForTest)
	g.SetCell(&c, TestCell{Alive: true})
	cell, _ := g.GetCell(&c)
	newLiveStatus := cell.Alive

	if newLiveStatus {
		t.Log("Passed")
	} else {
		t.Fatalf("Should correctly set cell.")
	}
}

func TestSetCell(t *testing.T) {
	shouldThrowErrorWhenCellSeedExceedBoarder(t)
	shouldSetCellCorrectly(t)
}

func testBlockIteratement(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initialTestCell, defaultCellIteratorForTest)

	// Make a block pattern
	g.SetCell(&Coordinate{X: 0, Y: 0}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 0, Y: 1}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 0}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 1}, TestCell{Alive: true})
	g.Iterate()

	nextAliveCellsMap := *convertTestCellsMatricToAliveTestCellsMap(g.GetGeneration())
	expectedNextAliveCellsMap := aliveTestCellsMap{
		{true, true, false},
		{true, true, false},
		{false, false, false},
	}

	if areAliveTestCellsMapsEqual(nextAliveCellsMap, expectedNextAliveCellsMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Should generate next cellLiveMap of a block, but got %v.", nextAliveCellsMap)
	}
}

func testBlinkerIteratement(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initialTestCell, defaultCellIteratorForTest)

	// Make a blinker pattern
	g.SetCell(&Coordinate{X: 1, Y: 0}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 1}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 2}, TestCell{Alive: true})

	var cellLiveMap aliveTestCellsMap

	expectedNextAliveCellsMapOne := aliveTestCellsMap{
		{false, true, false},
		{false, true, false},
		{false, true, false},
	}
	expectedNextAliveCellsMapTwo := aliveTestCellsMap{
		{false, false, false},
		{true, true, true},
		{false, false, false},
	}

	g.Iterate()
	cellLiveMap = *convertTestCellsMatricToAliveTestCellsMap(g.GetGeneration())
	if !areAliveTestCellsMapsEqual(cellLiveMap, expectedNextAliveCellsMapOne) {
		t.Fatalf("Should generate next cellLiveMap of a blinker, but got %v.", cellLiveMap)
	}

	g.Iterate()
	cellLiveMap = *convertTestCellsMatricToAliveTestCellsMap(g.GetGeneration())
	if !areAliveTestCellsMapsEqual(cellLiveMap, expectedNextAliveCellsMapTwo) {
		t.Fatalf("Should generate 2nd next cellLiveMap of a blinker, but got %v.", cellLiveMap)
	}
}

func testGliderIteratement(t *testing.T) {
	width := 5
	height := 5
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initialTestCell, defaultCellIteratorForTest)

	// Make a glider pattern
	g.SetCell(&Coordinate{X: 1, Y: 1}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 2, Y: 2}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 3, Y: 2}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 3}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 2, Y: 3}, TestCell{Alive: true})

	var cellLiveMap aliveTestCellsMap

	expectedAliveCellsMapOne := aliveTestCellsMap{
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, true, false, true, false},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}
	expectedAliveCellsMapTwo := aliveTestCellsMap{
		{false, false, false, false, false},
		{false, false, true, false, false},
		{false, false, false, true, true},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}
	expectedAliveCellsMapThree := aliveTestCellsMap{
		{false, false, false, false, false},
		{false, false, false, true, false},
		{false, false, false, false, true},
		{false, false, true, true, true},
		{false, false, false, false, false},
	}
	expectedAliveCellsMapFour := aliveTestCellsMap{
		{false, false, false, false, false},
		{false, false, false, false, false},
		{false, false, true, false, true},
		{false, false, false, true, true},
		{false, false, false, true, false},
	}

	g.Iterate()
	cellLiveMap = *convertTestCellsMatricToAliveTestCellsMap(g.GetGeneration())
	if !areAliveTestCellsMapsEqual(cellLiveMap, expectedAliveCellsMapOne) {
		t.Fatalf("Should generate next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Iterate()
	cellLiveMap = *convertTestCellsMatricToAliveTestCellsMap(g.GetGeneration())
	if !areAliveTestCellsMapsEqual(cellLiveMap, expectedAliveCellsMapTwo) {
		t.Fatalf("Should generate 2nd next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Iterate()
	cellLiveMap = *convertTestCellsMatricToAliveTestCellsMap(g.GetGeneration())
	if !areAliveTestCellsMapsEqual(cellLiveMap, expectedAliveCellsMapThree) {
		t.Fatalf("Should generate 3rd next next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	g.Iterate()
	cellLiveMap = *convertTestCellsMatricToAliveTestCellsMap(g.GetGeneration())
	if !areAliveTestCellsMapsEqual(cellLiveMap, expectedAliveCellsMapFour) {
		t.Fatalf("Should generate 4th next next cellLiveMap of a glider, but got %v.", cellLiveMap)
	}

	t.Log("Passed")
}

func testIteratementWithConcurrency(t *testing.T) {
	width := 200
	height := 200
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initialTestCell, defaultCellIteratorForTest)

	// Make a glider pattern
	g.SetCell(&Coordinate{X: 0, Y: 0}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 1}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 2, Y: 1}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 2, Y: 1}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 0, Y: 2}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 2}, TestCell{Alive: true})

	wg := sync.WaitGroup{}

	step := 100

	wg.Add(step)
	for i := 0; i < step; i++ {
		// Let the glider fly to digonal cell in four steps.
		go func() {
			g.Iterate()
			g.Iterate()
			g.Iterate()
			g.Iterate()
			wg.Done()
		}()
	}
	wg.Wait()

	cellOne, _ := g.GetCell(&Coordinate{X: 0 + step, Y: 0 + step})
	cellTwo, _ := g.GetCell(&Coordinate{X: 0 + step, Y: 2 + step})
	cellThree, _ := g.GetCell(&Coordinate{X: 1 + step, Y: 1 + step})
	cellFour, _ := g.GetCell(&Coordinate{X: 1 + step, Y: 2 + step})
	cellFive, _ := g.GetCell(&Coordinate{X: 2 + step, Y: 1 + step})

	if !cellOne.Alive || !cellTwo.Alive || !cellThree.Alive || !cellFour.Alive || !cellFive.Alive {
		t.Fatalf("Should still be a glider pattern.")
	}

	t.Log("Passed")
}

func TestIterate(t *testing.T) {
	testBlockIteratement(t)
	testBlinkerIteratement(t)
	testGliderIteratement(t)
	testIteratementWithConcurrency(t)
}

func testGetSizeCaseOne(t *testing.T) {
	width := 3
	height := 6
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initialTestCell, defaultCellIteratorForTest)

	if g.GetSize().Width == 3 && g.GetSize().Height == 6 {
		t.Log("Passed")
	} else {
		t.Fatalf("Size is not correct.")
	}
}

func TestGetSize(t *testing.T) {
	testGetSizeCaseOne(t)
}

func testGetCellCaseOne(t *testing.T) {
	width := 2
	height := 2
	size := Size{Width: width, Height: height}
	coord := Coordinate{X: 1, Y: 0}
	g, _ := NewGame(&size, initialTestCell, defaultCellIteratorForTest)
	g.SetCell(&coord, TestCell{Alive: true})
	cell, _ := g.GetCell(&coord)

	if cell.Alive == true {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not get correct cell at the coordinate.")
	}
}

func testGetCellCaseTwo(t *testing.T) {
	width := 2
	height := 2
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initialTestCell, defaultCellIteratorForTest)
	coord := Coordinate{X: 1, Y: 4}
	_, err := g.GetCell(&coord)

	if err == nil {
		t.Fatalf("Should get error when given coordinate is out of border.")
	} else {
		t.Log("Passed")
	}
}

func TestGetCell(t *testing.T) {
	testGetCellCaseOne(t)
	testGetCellCaseTwo(t)
}

func testResetCaseOne(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initialTestCell, defaultCellIteratorForTest)

	// Make a glider pattern
	g.SetCell(&Coordinate{X: 1, Y: 0}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 1}, TestCell{Alive: true})
	g.SetCell(&Coordinate{X: 1, Y: 2}, TestCell{Alive: true})

	g.Reset()
	cellLiveMap := convertTestCellsMatricToAliveTestCellsMap(g.GetGeneration())

	expectedBinaryBoard := aliveTestCellsMap{
		{false, false, false},
		{false, false, false},
		{false, false, false},
	}

	if areAliveTestCellsMapsEqual(*cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not reset cellLiveMap correctly.")
	}
}

func TestReset(t *testing.T) {
	testResetCaseOne(t)
}

func testSetCellIteratorCaseOne(t *testing.T) {
	width := 3
	height := 3
	size := Size{Width: width, Height: height}
	customCellIterator := func(coord *Coordinate, cell TestCell, getAdjacentCell GetAdjacentCell[TestCell]) *TestCell {
		nextCell := TestCell{}

		// Bring back all dead cells to alive in next iteration.
		if !nextCell.Alive {
			nextCell.Alive = true
			return &nextCell
		} else {
			nextCell.Alive = false
			return &nextCell
		}
	}
	g, _ := NewGame(&size, initialTestCell, customCellIterator)
	g.Iterate()
	cellLiveMap := convertTestCellsMatricToAliveTestCellsMap(g.GetGeneration())

	expectedBinaryBoard := aliveTestCellsMap{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}

	if areAliveTestCellsMapsEqual(*cellLiveMap, expectedBinaryBoard) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not set custom 'shouldCellDie' logic correcly.")
	}
}

func TestSetCellIterator(t *testing.T) {
	testSetCellIteratorCaseOne(t)
}

func testGetGenerationCaseOne(t *testing.T) {
	width := 2
	height := 2
	size := Size{Width: width, Height: height}
	g, _ := NewGame(&size, initialTestCell, defaultCellIteratorForTest)
	generation := g.GetGeneration()
	aliveCellsMap := convertTestCellsMatricToAliveTestCellsMap(generation)

	expectedCellsMap := [][]bool{{false, false}, {false, false}}

	if areAliveTestCellsMapsEqual(*aliveCellsMap, expectedCellsMap) {
		t.Log("Passed")
	} else {
		t.Fatalf("Did not get correct generation.")
	}
}

func TestGetGeneration(t *testing.T) {
	testGetGenerationCaseOne(t)
}
