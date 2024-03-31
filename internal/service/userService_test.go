package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		Name     string
		Password string
		Want     bool
	}{
		{
			Name:     "No errors check, hash/unhash password",
			Password: "password",
			Want:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			hash := HashPassword(test.Password)
			res := CompareToHash(hash, test.Password)
			fmt.Println()
			assert.Equal(t, test.Want, res)
		})
	}
}
