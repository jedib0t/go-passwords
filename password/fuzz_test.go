package password

import (
	"testing"

	"github.com/jedib0t/go-passwords/charset"
)

// FuzzNewGenerator asserts the constructor/generator contract: any rule
// combination either fails NewGenerator with an error, or yields a generator
// whose Generate never panics, never errors, and always produces a password
// of the requested length using only charset characters.
func FuzzNewGenerator(f *testing.F) {
	f.Add("abcDEF123!@#", 12, 1, 1, 0, 2)
	f.Add(string(charset.AllChars), 12, 3, 1, 1, 1)
	f.Add("!@#$%^&*", 4, 0, 0, 4, 4)
	f.Add("!@#$%^&*", 4, 0, 0, 1, 1)
	f.Add("aB1", 1, 0, 0, 0, 0)
	f.Add("abcdef", 6, 0, 0, 0, 3)
	f.Add("abcdef", 4, 0, 0, 0, 8)
	f.Fuzz(func(t *testing.T, cs string, length, minLower, minUpper, minSym, maxSym int) {
		if len(cs) > 256 || length > 1024 {
			t.Skip()
		}
		g, err := NewGenerator(
			WithCharset(charset.Charset(cs)),
			WithLength(length),
			WithMinLowerCase(minLower),
			WithMinUpperCase(minUpper),
			WithNumSymbols(minSym, maxSym),
		)
		if err != nil {
			// invalid configurations must be rejected, never panic
			return
		}

		charsetMap := make(map[rune]bool)
		for _, r := range cs {
			charsetMap[r] = true
		}
		for i := 0; i < 5; i++ {
			pw, err := g.Generate()
			if err != nil {
				t.Fatalf("Generate failed after successful NewGenerator: %v", err)
			}
			runes := []rune(pw)
			if len(runes) != length {
				t.Fatalf("password %q has %d runes, want %d", pw, len(runes), length)
			}
			for _, r := range runes {
				if !charsetMap[r] {
					t.Fatalf("password %q contains %q which is not in charset %q", pw, r, cs)
				}
			}
		}
	})
}
