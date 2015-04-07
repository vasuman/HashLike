package pow

// This interface represents a typical challenge-response
// proof-of-work system.
type POW interface {

	// Generates a challenge that the client is required to solve
	Challenge(url string, remoteAddr []byte) []byte

	// Verifies a solution that the client submits and returns a
	// reward value. If the solution is invalid, it should return
	// zero.
	Verify(hash, challenge []byte, nonce int) float64
}
