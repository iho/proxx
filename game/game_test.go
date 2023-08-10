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
	for i := 0; i < field.Height; i++ {
		for j := 0; j < field.Width; j++ {
			if field.Cells[i][j].IsBlackBomb {
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
