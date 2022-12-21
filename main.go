package main

import (
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BoardSpace int

const (
	X     BoardSpace = -1
	Blank BoardSpace = 0
	O     BoardSpace = 1
)

func checkWin(board [][]BoardSpace) BoardSpace {
	for i := 0; i < 3; i++ {
		if board[i][0] == board[i][1] && board[i][1] == board[i][2] && board[i][0] != Blank {
			return board[i][0]
		}
		if board[0][i] == board[1][i] && board[1][i] == board[2][i] && board[0][i] != Blank {
			return board[0][i]
		}
	}

	if board[0][0] == board[1][1] && board[1][1] == board[2][2] && board[0][0] != Blank {
		return board[0][0]
	}
	if board[2][0] == board[1][1] && board[1][1] == board[0][2] && board[2][0] != Blank {
		return board[2][0]
	}

	return Blank
}

func checkFull(board [][]BoardSpace) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == Blank {
				return false
			}
		}
	}
	return true
}

func minimax(board [][]BoardSpace, maximizer bool) int {
	winner := checkWin(board)

	if winner != 0 {
		return int(winner)
	}
	if checkFull(board) {
		return 0
	}

	if maximizer {
		best := math.MinInt
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == Blank {
					board[i][j] = O
					res := minimax(board, false)
					if res > best {
						best = res
					}
					board[i][j] = Blank
				}
			}
		}
		return best
	} else {
		best := math.MaxInt
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == Blank {
					board[i][j] = X
					res := minimax(board, true)
					if res < best {
						best = res
					}
					board[i][j] = Blank
				}
			}
		}
		return best
	}
}

func nextMove(board [][]BoardSpace, player BoardSpace) {
	if player == Blank {
		return
	}

	if checkWin(board) != Blank || checkFull(board) {
		return
	}

	var bestScore int
	var maximizer bool
	if player == X {
		bestScore = math.MaxInt
		maximizer = false
	} else {
		bestScore = math.MinInt
		maximizer = true
	}

	bestRow, bestCol := 0, 0

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == Blank {
				board[i][j] = player

				score := minimax(board, !maximizer)
				if (player == X && score < bestScore) || (player == O && score > bestScore) {
					bestScore = score
					bestRow = i
					bestCol = j
				}

				board[i][j] = Blank
			}
		}
	}

	board[bestRow][bestCol] = player
}

type Response struct {
	Board [][]string `json:"board"`
	Won   string     `json:"won"`
	Full  bool       `json:"full"`
}

type Request struct {
	Board [][]string `json:"board"`
}

func toBoard(boardStrings [][]string) [][]BoardSpace {
	board := [][]BoardSpace{{Blank, Blank, Blank}, {Blank, Blank, Blank}, {Blank, Blank, Blank}}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			var space BoardSpace
			switch strings.ToUpper(boardStrings[i][j]) {
			case "X":
				space = X
			case "O":
				space = O
			default:
				space = Blank
			}

			board[i][j] = space
		}
	}

	return board
}

func toBoardStrings(board [][]BoardSpace) [][]string {
	boardStrings := [][]string{{"_", "_", "_"}, {"_", "_", "_"}, {"_", "_", "_"}}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			var space string
			switch board[i][j] {
			case X:
				space = "X"
			case O:
				space = "O"
			default:
				space = "_"
			}

			boardStrings[i][j] = space
		}
	}

	return boardStrings
}

func checkDims(board [][]string) bool {
	if len(board) != 3 {
		return false
	}

	for i := 0; i < 3; i++ {
		if len(board[i]) != 3 {
			return false
		}
	}

	return true
}

func moveReq(c *gin.Context, player BoardSpace) {
	var req Request

	if err := c.BindJSON(&req); err != nil {
		return
	}

	if !checkDims(req.Board) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	board := toBoard(req.Board)
	nextMove(board, player)
	var won string
	switch checkWin(board) {
	case X:
		won = "X"
	case O:
		won = "O"
	default:
		won = "_"
	}
	full := checkFull(board)

	c.JSON(http.StatusOK, Response{toBoardStrings(board), won, full})
}

func moveX(c *gin.Context) {
	moveReq(c, X)
}

func moveO(c *gin.Context) {
	moveReq(c, O)
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.POST("/move/x", moveX)
	router.POST("/move/o", moveO)

	router.Run()
}
