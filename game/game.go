package game

import (
	"errors"
	"fmt"
	"math/rand"
)

var ValidationError = errors.New("can't create Proox game field with provided parameters")

type CellsState struct {
	IsBlackBomb     bool
	IsVisible       bool
	AdjacentCounter int
}

type ProxxField struct {
	Width            int
	Height           int
	BlackHolesNumber int
	Cells            [][]*CellsState
}

func NewProxxField(blackHolesNumber int, width int, height int) (*ProxxField, error) {
	if blackHolesNumber <= 0 || width <= 0 || height <= 0 || height*width <= blackHolesNumber {
		return nil, ValidationError
	}

	cells := make([][]*CellsState, height)
	for h := 0; h < height; h++ {
		row := make([]*CellsState, width)

		for w := 0; w < width; w++ {
			row[w] = &CellsState{
				IsBlackBomb:     false,
				IsVisible:       false,
				AdjacentCounter: 0,
			}
		}

		cells[h] = row
	}

	field := &ProxxField{
		Width:            width,
		Height:           height,
		BlackHolesNumber: blackHolesNumber,
		Cells:            cells,
	}
	field.PlaceBlackBombs()
	field.PopulateAdjacentCounters()

	return field, nil
}

func (field *ProxxField) PlaceBlackBombs() {
	for i := 0; i < field.BlackHolesNumber; {
		w := rand.Intn(field.Width)
		h := rand.Intn(field.Height)
		if field.Cells[h][w].IsBlackBomb {
			continue
		}
		field.Cells[h][w].IsBlackBomb = true
		i++
	}
}

func (field *ProxxField) PopulateAdjacentCounters() {
	for h := 0; h < field.Height; h++ {
		for w := 0; w < field.Width; w++ {
			if field.Cells[h][w].IsBlackBomb {
				if h+1 < field.Height {
					nextRow := h + 1
					if (w - 1) >= 0 {
						field.Cells[nextRow][w-1].AdjacentCounter++
					}
					field.Cells[nextRow][w].AdjacentCounter++
					if (w + 1) < field.Width {
						field.Cells[nextRow][w+1].AdjacentCounter++
					}
				}
				if (h - 1) >= 0 {
					previousRow := h - 1
					if (w - 1) >= 0 {
						field.Cells[previousRow][w-1].AdjacentCounter++
					}
					field.Cells[h-1][w].AdjacentCounter++
					if (w + 1) < field.Width {
						field.Cells[previousRow][w+1].AdjacentCounter++
					}
				}
				if (w - 1) >= 0 {
					field.Cells[h][w-1].AdjacentCounter++
				}
				if (w + 1) < field.Width {
					field.Cells[h][w+1].AdjacentCounter++
				}
			}
		}
	}
}

// ToString is for debug purpose only
func (field *ProxxField) ToString() string {
	result := ""
	for i := 0; i < field.Height; i++ {
		for j := 0; j < field.Width; j++ {
			result += fmt.Sprintf("%t ", field.Cells[i][j].IsBlackBomb)
		}
		result += "\n"
	}
	for i := 0; i < field.Height; i++ {
		for j := 0; j < field.Width; j++ {
			result += fmt.Sprintf("%d ", field.Cells[i][j].AdjacentCounter)
		}
		result += "\n"
	}
	return result
}

// Direction offsets for adjacent cells
var directions = []struct {
	dx, dy int
}{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

func (field *ProxxField) RevealCell(h, w int) {
	queue := [][]int{{h, w}}

	visited := make([][]bool, field.Height)
	for i := range visited {
		visited[i] = make([]bool, field.Width)
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		ch, cw := curr[0], curr[1]
		cell := field.Cells[ch][cw]

		if visited[ch][cw] {
			continue
		}

		visited[ch][cw] = true
		cell.IsVisible = true

		if cell.AdjacentCounter == 0 && !cell.IsBlackBomb {
			for _, d := range directions {
				nh, nw := ch+d.dx, cw+d.dy
				if 0 <= nh && nh < field.Height && 0 <= nw && nw < field.Width && !visited[nh][nw] && !field.Cells[nh][nw].IsBlackBomb {
					queue = append(queue, []int{nh, nw})
				}
			}
		}
	}
}
