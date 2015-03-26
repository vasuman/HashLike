package pow

import "testing"

func TestCountLeadingZeros(t *testing.T) {
	assertZeros := func(num byte, zeros int) {
		b := make([]byte, 1)
		b[0] = num
		if c := countLeadingZeros(b); c != zeros {
			t.Errorf("zero count %d. expected %d, got %d", num, zeros, c)
		}
	}
	b := []byte{0, 0, 0, 0}
	if countLeadingZeros(b) != 32 {
		t.Error("slice zero count fail")
	}
	b = []byte{0, 0, 48, 0}
	if c := countLeadingZeros(b); c != 18 {
		t.Errorf("slice zero count fail. Got %d\n", c)
	}
	assertZeros(7, 5)
	assertZeros(5, 5)
	assertZeros(1, 7)
	assertZeros(3, 6)
	assertZeros(17, 3)
	assertZeros(127, 1)
	assertZeros(129, 0)
}
