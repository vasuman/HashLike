package pow

import "time"

// This interface represents a typical challenge-response
// proof-of-work system.
type POW interface {

	// Generates a challenge that the client is required to solve
	Challenge(url string, addr []byte, when time.Time) []byte

	// Verifies a solution that the client submits and returns a
	// reward value. If the solution is invalid, it should return
	// zero.
	Verify(challenge []byte, nonce int) float64

	// Returns a string describing the system
	Desc() string
}

var Systems = map[string]POW{
	"HC": new(Hashcash),
}
