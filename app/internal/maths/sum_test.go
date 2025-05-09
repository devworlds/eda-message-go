package maths

import "testing"

func TestSum(t *testing.T) {
	cases := []struct {
		a, b, result int
	}{
		{1, 1, 2},
		{2, 2, 4},
		{10, -2, 8},
	}
	for _, c := range cases {
		resultado := Sum(c.a, c.b)
		if resultado != c.result {
			t.Errorf("Sum(%d, %d) == %d, result %d", c.a, c.b, resultado, c.result)
		}
	}
}
