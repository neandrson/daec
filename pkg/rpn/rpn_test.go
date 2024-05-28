package rpn_test

import (
	"slices"
	"testing"

	"github.com/Vojan-Najov/daec/pkg/rpn"
)

func TestRpn(t *testing.T) {
	input := [...]string{
		"1 + 2",
		"2*3 + 4*5",
	}
	expectedArr := [...][]string{
		[]string{"1", "2", "+"},
		[]string{"2", "3", "*", "4", "5", "*", "+"},
	}
	expectedErrors := [...]error{
		nil,
		nil,
	}
	for i, expected := range expectedArr {
		expectedErr := expectedErrors[i]
		actual, err := rpn.NewRPN(input[i])
		if expectedErr == nil && err != nil ||
			expectedErr != nil && err == nil {
			t.Errorf("expected %v, actual %v", expectedErr, err)
		}
		if !slices.Equal(actual.Token, expected) {
			t.Errorf("expected %v, actual %v", expected, actual)
		}
	}
}
