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

	// Returns the name of the system
	Name() string
}

const (
	HashcashIdent = "HC"
)

var Systems = map[string]pow{
	HashcashIdent: new(Hashcash),
}
