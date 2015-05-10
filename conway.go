package main

import (
	"flag"
	"io/ioutil"
	"math/rand"
	"time"
)

const DeadChar byte = '.'
const LiveChar byte = 'x'

const UninitializedBoardSize int = -1
const DefaultBoardSize int = 40

/*
	Main application:
*/

func main() {

	/*
		Read in command line parameters:
	*/

	ptrBoardSize := flag.Int("size", UninitializedBoardSize, "Size of the NxN grid of cells.")
	ptrInFile := flag.String("input", "", "Text file to initialize the grid.")
	flag.Parse()

	/*
		Declare variables, initialize, and error-check:
	*/

	rand.Seed(time.Now().UTC().UnixNano())

	// TODO: update readme for double buffering details.

	var BoardSize int = *ptrBoardSize
	var inFile string = *ptrInFile

	var UseInFile bool = len(*ptrInFile) > 0
	var UseRand bool = !UseInFile

	var showGrid0 bool = false
	var exit = false

	// Current implementation: double buffer the grid. 
	// Storing the changes in a list may be more efficient.
	var grid0 [][]bool
	var grid1 [][]bool

	if UseInFile && BoardSize != UninitializedBoardSize {
		panic("Error! Cannot initialize from file and set board size.")
	} else if BoardSize < 0 && BoardSize != UninitializedBoardSize {
		panic("Error! Cannot set board size as a negative number.")
	}

	if UseInFile {
		BoardSize, grid0 = readConfigurationFile(inFile)
	} else if UseRand {
		if BoardSize == UninitializedBoardSize {
			BoardSize = DefaultBoardSize
		}
		grid0 = initializeRandom(BoardSize)
	}

	grid1 = createGrid(BoardSize)

	// Initialize UI:
	var ui Ui = Ui{}
	ui.Init(BoardSize)

	/*
		Main application logic:
	*/

	// Loop and update the board, checking UI thread for exit condition in ui.Update()
	for exit == false {
		// Local aliases for grid0, grid1:
		var currentGrid [][]bool
		var nextGrid [][]bool

		// Alternate which buffer is being shown to screen:
		showGrid0 = !showGrid0
		if showGrid0 {
			currentGrid = grid0
			nextGrid = grid1
		} else {
			currentGrid = grid1
			nextGrid = grid0
		}

		// Count neighbours for each cell, update next state (n^2 complexity):
		for row := 0; row < BoardSize; row++ {
			for col := 0; col < BoardSize; col++ {
				var neighbours int = countNeighbours(currentGrid, row, col, BoardSize)
				if currentGrid[row][col] {
					// Live if neighbours equals 2 or 3; otherwise, die.
					nextGrid[row][col] = neighbours == 2 || neighbours == 3
				} else {
					// Live if neighbours equals 3; otherwise, stay dead.
					nextGrid[row][col] = neighbours == 3
				}
			}
		}

		ui.PrintGrid(currentGrid, BoardSize)
		exit = ui.Update()
	}

	ui.Destroy()
}

/*
	Initialization functions:
*/

func readConfigurationFile(inFile string) (int, [][]bool) {
	const CarriageReturnChar byte = '\r'
	const LineFeedChar byte = '\n'
	const MessageGridSizeError = "Error! Input file is not a correctly formatted NxN grid."

	var dimension int
	var grid [][]bool

	data, err := ioutil.ReadFile(inFile)
	if err != nil {
		panic(err)
	}

	// For a (correctly) formatted as an NxN grid, find N:
	for _, char := range data {
		if char == DeadChar || char == LiveChar {
			dimension++
		} else if char == CarriageReturnChar || char == LineFeedChar {
			break
		}
	}

	grid = createGrid(dimension)

	// Read the input file into the grid. Panic if incorrectly formatted.
	row, col := 0, 0
	for _, char := range data {
		if char == LiveChar {
			grid[col][row] = true
			col++
		} else if char == DeadChar {
			grid[col][row] = false
			col++
		} else if char == CarriageReturnChar || char == LineFeedChar {
			if col == dimension {
				row++
				col = 0
			} else if col == 0 {
				// Do nothing.
			} else {
				panic(MessageGridSizeError)
			}
		}
	}

	return dimension, grid
}

func initializeRandom(dimension int) [][]bool {
	var grid [][]bool = createGrid(dimension)
	for i := 0; i < dimension; i++ {
		for j := 0; j < dimension; j++ {
			grid[i][j] = rand.Intn(1+1) == 1 // Random 0 or 1
		}
	}

	return grid
}

func createGrid(dimension int) [][]bool {
	var grid [][]bool = make([][]bool, dimension)

	for i := 0; i < dimension; i++ {
		grid[i] = make([]bool, dimension)
	}

	return grid
}

/*
	Helper functions:
*/

func countNeighbours(currentGrid [][]bool, row int, col int, BoardSize int) int {
	var sum int = 0

	for rowIdx := -1; rowIdx <= 1; rowIdx++ {
		for colIdx := -1; colIdx <= 1; colIdx++ {
			var nbrRow = row + rowIdx
			var nbrCol = col + colIdx

			if rowIdx == 0 && colIdx == 0 {
				continue
			}

			// Implements a toroidal grid. (ie. grid wraps around)
			if nbrRow < 0 {
				nbrRow = BoardSize - 1
			} else if nbrRow >= BoardSize {
				nbrRow = 0
			}

			if nbrCol < 0 {
				nbrCol = BoardSize - 1
			} else if nbrCol >= BoardSize {
				nbrCol = 0
			}

			if currentGrid[nbrRow][nbrCol] {
				sum++
			}
		}
	}

	return sum
}
