package toolkit

import "crypto/rand"

const randomStringSource = "abcdefghijklmnopqrstuvwxyzABCDEFGHIKHLMNOPQRSTUVWXYZ0123456789_+"

// Tools is the type use to intantiate this module. Any variable
// ofthis type would have access to all methods with the reciver *Tools
type Tools struct{}

// RandomString returns a string of random charactors of length n, using
// randomString Source as the source for the string
func (t *Tools) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randomStringSource)
	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}
	return string(s)
}
