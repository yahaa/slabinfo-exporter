package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExecWithOutput(t *testing.T) {
	tests := []struct {
		Input string
		Want  string
	}{
		{"echo '666'", "666"},
		{"echo '123456'", "123456"},
	}

	for _, test := range tests {
		output, _ := ExecWithOutput(test.Input)
		assert.Equal(t, test.Want, output)
	}
}
