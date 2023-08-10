package game_test

import (
	"proxx/game"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProxxField(t *testing.T) {
	type test struct {
		width            int
		height           int
		blackHolesNumber int
		expectsErr       bool
	}

	tests := []test{
		{
			width:            0,
			height:           0,
			blackHolesNumber: 0,
			expectsErr:       true,
		},
		{
			width:            1,
			height:           1,
			blackHolesNumber: 1,
			expectsErr:       true,
		},
		{
			width:            3,
			height:           3,
			blackHolesNumber: 8,
			expectsErr:       false,
		},
		{
			width:            -3,
			height:           3,
			blackHolesNumber: 7,
			expectsErr:       true,
		},
		{
			width:            3,
			height:           -3,
			blackHolesNumber: 7,
			expectsErr:       true,
		},
		{
			width:            3,
			height:           3,
			blackHolesNumber: 11,
			expectsErr:       true,
		},
	}

	for _, tc := range tests {
		field, err := game.NewProxxField(tc.blackHolesNumber, tc.width, tc.height)
		if tc.expectsErr {
			assert.NotNil(t, err, "must be error")
			assert.Nil(t, field, "must be nil")
		} else {
			assert.Nil(t, err, "must be no error")
			assert.NotNil(t, field, "must be field")
		}
	}
}

func CountBlackHoles(field *game.ProxxField) int {
	count := 0
	for h := 0; h < field.Height; h++ {
		for w := 0; w < field.Width; w++ {
			if field.Cells[h][w].IsBlackBomb {
				count++
			}
		}
	}

	return count
}

func TestPlaceBlackBombs(t *testing.T) {
	maxHeight := 5
	maxWidth := 5
	for height := 1; height <= maxHeight; height++ {
		for width := 1; width <= maxWidth; width++ {
			maxBlackHolesNumber := (height * width) - 1
			for blackHolesNumber := 1; blackHolesNumber <= maxBlackHolesNumber; blackHolesNumber++ {
				field, err := game.NewProxxField(blackHolesNumber, width, height)
				assert.Nil(t, err, "must be no error")
				assert.NotNil(t, field, "must be field")
				assert.Equal(t, blackHolesNumber, CountBlackHoles(field))
			}
		}
	}
}

func CheckBlackHoleAdjacentCounters(field *game.ProxxField) (*game.ProxxField, bool) {
	// avoid boundchecks at cost of additional memory

	cells := make([][]*game.CellsState, field.Height+2)
	biggerField := &game.ProxxField{
		BlackHolesNumber: 0,
		Width:            field.Width + 2,
		Height:           field.Height + 2,
		Cells:            cells,
	}
	for h := 0; h < field.Height+2; h++ {
		row := make([]*game.CellsState, field.Width+2)

		for w := 0; w < field.Width+2; w++ {
			row[w] = &game.CellsState{
				IsBlackBomb:     false,
				IsVisible:       false,
				AdjacentCounter: 0,
			}
		}
		cells[h] = row
	}

	for h := 0; h < field.Height; h++ {
		for w := 0; w < field.Width; w++ {
			cells[h+1][w+1] = field.Cells[h][w]
		}
	}
	for h := 1; h < field.Height+1; h++ {
		for w := 1; w < field.Width+1; w++ {
			count := 0
			if cells[h-1][w-1].IsBlackBomb {
				count++
			}
			if cells[h-1][w].IsBlackBomb {
				count++
			}
			if cells[h-1][w+1].IsBlackBomb {
				count++
			}
			if cells[h][w-1].IsBlackBomb {
				count++
			}
			if cells[h][w+1].IsBlackBomb {
				count++
			}
			if cells[h+1][w-1].IsBlackBomb {
				count++
			}
			if cells[h+1][w].IsBlackBomb {
				count++
			}
			if cells[h+1][w+1].IsBlackBomb {
				count++
			}
			if count != cells[h][w].AdjacentCounter {
				return biggerField, false
			}
		}
	}
	return biggerField, true
}

func TestPopulateAdjacentCounters(t *testing.T) {
	maxHeight := 10
	maxWidth := 10
	for height := 1; height <= maxHeight; height++ {
		for width := 1; width <= maxWidth; width++ {
			maxBlackHolesNumber := (height * width) - 1
			for blackHolesNumber := 1; blackHolesNumber <= maxBlackHolesNumber; blackHolesNumber++ {
				field, err := game.NewProxxField(blackHolesNumber, width, height)
				assert.Nil(t, err, "must be no error")
				assert.NotNil(t, field, "must be field")
				biggerField, correct := CheckBlackHoleAdjacentCounters(field)
				if !correct {
					assert.Fail(t, "This field is incorrect: %s\n %s", field.ToString(), biggerField.ToString())
					return
				}
			}
		}
	}
}

func TestRevealCell(t *testing.T) {
	gameField, _ := game.NewProxxField(1, 3, 3) // Small 3x3 board with 1 black hole for easy testing

	countVisibleCells := func() int {
		count := 0
		for i := range gameField.Cells {
			for _, cell := range gameField.Cells[i] {
				if cell.IsVisible {
					count++
				}
			}
		}
		return count
	}

	assert.Equal(t, 0, countVisibleCells(), "Expected all cells to be hidden at start")

	// Reveal a cell that's guaranteed not to be a black hole (since there's only 1 black hole)
	gameField.RevealCell(0, 0)

	assert.True(t, gameField.Cells[0][0].IsVisible, "Expected the clicked cell (0, 0) to be visible")

	assert.LessOrEqual(t, 1, countVisibleCells(), "Expected more cells to be revealed due to cascading reveals")

	// Finally, make sure the black hole itself is not revealed
	for i := range gameField.Cells {
		for j, cell := range gameField.Cells[i] {
			if cell.IsBlackBomb && cell.IsVisible {
				assert.Falsef(t, cell.IsVisible, "Black hole at (%d, %d) should not be revealed", i, j)
			}
		}
	}
}
