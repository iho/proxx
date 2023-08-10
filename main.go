package main

import "errors"

type FieldState struct {
	IsBlackBomb     bool
	IsVisible       bool
	AdjacentCounter int
}

type ProxxField struct {
	Width            int
	Height           int
	BlackHolesNumber int
	Fields           [][]*FieldState
}

var ValidationError = errors.New("can't create Proox game field with provided parameters")

func NewProxxField(blackHolesNumber int, width int, height int) (*ProxxField, error) {
	if blackHolesNumber <= 0 || width <= 0 || height <= 0 {
		return nil, ValidationError
	}

	fields := make([][]*FieldState, height)
	for i := 0; i < height; i++ {
		row := make([]*FieldState, width)

		for j := 0; j < width; j++ {
			row[j] = &FieldState{
				IsBlackBomb:     false,
				IsVisible:       false,
				AdjacentCounter: 0,
			}
		}

		fields[i] = row
	}

	return &ProxxField{
		Width:  width,
		Height: height,
	}, nil
}

func main() {

}
