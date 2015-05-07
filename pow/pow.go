package pow

import (
	"net"
	"net/url"
	"time"
)

// This struct encapsulates the parameters that a proof-of-work system
// uses to generate challenges
type ChallengeParams struct {
	// Location of the page
	Loc *url.URL
	// The remote address of the challenger
	Addr net.IP
	// Timestamp of when the challenge was recieved
	When time.Time
}

// This interface represents a typical challenge-response
// proof-of-work system
type POW interface {

	// Generates a challenge that the client is required to solve
	Challenge(params *ChallengeParams) []byte

	// Verifies a solution that the client submits and returns a
	// reward value. If the solution is invalid, it should return
	// zero
	Verify(challenge []byte, nonce int) float64

	// Returns a string describing the system
	Desc() string
}

// Genrates a map that maps proof-of-work identifiers to strings that
// describe the systems
func descMap() map[string]string {
	ret := make(map[string]string, len(Systems))
	for k, v := range Systems {
		ret[k] = v.Desc()
	}
	return ret
}

var (
	// Maps proof-of-work string identifiers to instances
	Systems = map[string]POW{
		"HC": new(Hashcash),
	}
	Desc = descMap()
)
