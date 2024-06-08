package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Tile struct {
	fired  bool
	x, y   int
	mValue int     // manhattan value
	pValue float32 // payout value
}

type Board struct {
	tiles  [][]Tile
	moves  []Tile
	odds   float64
	length int
}

func calcMValue(x, y, xCenter, yCenter int) int {
	return int(math.Abs(float64(x-xCenter)) + math.Abs(float64(y-yCenter)))
}

func calcPValue(mValue, boardLength int) float32 {
	return 0.0
}

func NewBoard(odds float64, length int) *Board {
	if odds < 0 || odds > 1 {
		panic("Odds must be between 0 and 1")
	}
	board := new(Board)
	board.odds = odds
	board.length = length
	board.tiles = make([][]Tile, board.length)
	for i := 0; i < board.length; i++ {
		board.tiles[i] = make([]Tile, board.length)
		for j := 0; j < board.length; j++ {
			mValue := calcMValue(i, j, board.length/2, board.length/2)
			board.tiles[i][j] = Tile{
				false,
				i, j,
				mValue,
				calcPValue(mValue, length),
			}
		}
	}
	return board
}

func updateMValues(board *Board) {
	if len(board.moves) == 0 {
		return
	}
	startingX, startingY := board.moves[0].x, board.moves[0].y
	for i := 0; i < board.length; i++ {
		for j := 0; j < board.length; j++ {
			board.tiles[i][j].mValue = calcMValue(i, j, startingX, startingY)
		}
	}
}

func triggerTile(board *Board, x int, y int, initialTile bool) {
	if x < 0 || x >= board.length || y < 0 || y >= board.length || board.tiles[x][y].fired {
		return
	}
	if initialTile {
		board.tiles[x][y].fired = true
		board.moves = append(board.moves, board.tiles[x][y])
		triggerTile(board, x+1, y, false)
		triggerTile(board, x-1, y, false)
		triggerTile(board, x, y+1, false)
		triggerTile(board, x, y-1, false)
		return
	}

	chance := rand.Float64()
	if chance < board.odds {
		board.tiles[x][y].fired = true
		board.moves = append(board.moves, board.tiles[x][y])
		triggerTile(board, x+1, y, false)
		triggerTile(board, x-1, y, false)
		triggerTile(board, x, y+1, false)
		triggerTile(board, x, y-1, false)
		return
	}
	board.moves = append(board.moves, board.tiles[x][y])
}

func printBoard(board *Board) {
	for i := 0; i < board.length; i++ {
		for j := 0; j < board.length; j++ {
			if board.tiles[i][j].fired {
				if board.tiles[i][j].x == board.moves[0].x && board.tiles[i][j].y == board.moves[0].y {
					fmt.Print("🔥")
				} else {
					fmt.Print(" _ ")
				}
			} else {
				// fmt.Print("0 ")
				if board.tiles[i][j].mValue < 10 {
					fmt.Printf(" %d  ", board.tiles[i][j].mValue)
				} else {
					fmt.Printf("%d  ", board.tiles[i][j].mValue)
				}
			}
		}
		fmt.Println()
	}
}

func printMoves(board *Board) {
	for i := 0; i < len(board.moves); i++ {
		fmt.Printf("(%d, %d), %t \n", board.moves[i].x, board.moves[i].y, board.moves[i].fired)
	}
	fmt.Println()

}

func printFiredCount(board *Board) int {
	count := 0
	for i := 0; i < board.length; i++ {
		for j := 0; j < board.length; j++ {
			if board.tiles[i][j].fired {
				count++
			}
		}
	}
	return count
}

func runGame(board *Board) {
	printBoard(board)
	// get user input for tile to trigger
	x := -1
	y := -1
	for x < 0 || x >= board.length || y < 0 || y >= board.length {
		fmt.Print("Enter x coordinate: ")
		fmt.Scanln(&x)
		fmt.Print("Enter y coordinate: ")
		fmt.Scanln(&y)
	}
	triggerTile(board, x, y, true)
	// updateMValues(board)
	printBoard(board)
}

func polynomialOdds(n, k, a float64) float64 {
	return math.Pow(k, n) / (a + math.Pow(n, k))
}

func exponentialOdds(n, k float64) float64 {
	return 1 - math.Pow(math.E, -k*n)
}

func main() {
	length := 0
	max := 33
	for length%2 == 0 || length < 2 || length > max {
		fmt.Printf("Enter the length of row/cols (odd number), max %d:\n", max)
		fmt.Scanln(&length)
	}
	// odds := polynomialOdds(float64(length), 1.08, 10.0)
	odds := exponentialOdds(float64(length), 0.02)
	fmt.Printf("odds: %0.2f\n", odds)
	board := *NewBoard(odds, length)
	runGame(&board)
	// printMoves(&board)
	fmt.Printf("%d/%d (%0.2f%%) tiles on fire! \n", printFiredCount(&board), length*length, float32(printFiredCount(&board))/float32(length*length)*100)
}
