package pow

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
	"time"
)

func encodeTime(t time.Time) []byte {
	tb := make([]byte, 8)
	binary.BigEndian.PutUint64(tb, t.Unix())
	return tb
}

type Hashcash struct{}

func (Hashcash) Challenge(url string, addr []byte, when time.Time) []byte {
	b := []byte(url)
	b = append(b, encodeTime(when)...)
	return append(b, addr...)
}

func countLeadingZeros(h []byte) int {
	var z int
	for i := 0; i < len(h); i++ {
		if h[i] != 0 {
			z = i * 8
			v := h[i]
			if v < 16 {
				v <<= 4
				z += 4
			}
			if v < 64 {
				v <<= 2
				z += 2
			}
			if v < 128 {
				z += 1
			}
			return z
		}
	}
	// Input is zero
	return len(h) * 8
}

func (Hashcash) Verify(challenge, hash []byte, nonce int) float64 {
	nB := make([]byte, 4)
	binary.BigEndian.PutUint32(nB, uint32(nonce))
	hC := append(hash, challenge...)
	b := sha256.Sum256(append(hC, nB...))
	leadZeros := countLeadingZeros(b[:])
	return math.Pow(2, float64(leadZeros))
}

func (Hashcash) Name() {
	return "Hashcash"
}
