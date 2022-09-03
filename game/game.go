package game

import (
	"errors"
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

func NewProxxField(blackHolesNumber int, width int, height int) (*ProxxField, error) {
	if blackHolesNumber <= 0 || width <= 0 || height <= 0 || height*width <= blackHolesNumber {
		return nil, ValidationError
	}

	cells := make([][]*CellsState, height)
	for i := 0; i < height; i++ {
		row := make([]*CellsState, width)

		for j := 0; j < width; j++ {
			row[j] = &CellsState{
				IsBlackBomb:     false,
				IsVisible:       false,
				AdjacentCounter: 0,
			}
		}

		cells[i] = row
	}

	field := &ProxxField{
		Width:            width,
		Height:           height,
		BlackHolesNumber: blackHolesNumber,
		Cells:            cells,
	}
	field.PlaceBlackBombs()

	return field, nil
}
