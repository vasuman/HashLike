package pow

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
	"time"
)

func encodeTime(t time.Time) []byte {
	tb := make([]byte, 8)
	unixTime := uint64(t.Unix())
	binary.BigEndian.PutUint64(tb, unixTime)
	return tb
}

type Hashcash struct{}

func (*Hashcash) Challenge(params *ChallengeParams) []byte {
	b := []byte(params.Loc.String())
	b = append(b, encodeTime(params.When)...)
	return append(b, params.Addr...)
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

func (*Hashcash) Verify(challenge []byte, nonce int) float64 {
	nB := make([]byte, 4)
	binary.BigEndian.PutUint32(nB, uint32(nonce))
	b := sha256.Sum256(append(challenge, nB...))
	leadZeros := countLeadingZeros(b[:])
	if leadZeros == 0 {
		return 0
	}
	return math.Pow(2, float64(leadZeros))
}

func (*Hashcash) Desc() string {
	return "Hashcash using SHA-256"
}
