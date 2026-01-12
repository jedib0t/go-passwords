# go-passwords

[![Go Reference](https://pkg.go.dev/badge/github.com/jedib0t/go-passwords/v0.svg)](https://pkg.go.dev/github.com/jedib0t/go-passwords)
[![Build Status](https://github.com/jedib0t/go-passwords/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/jedib0t/go-passwords/actions?query=workflow%3ACI+event%3Apush+branch%3Amain)
[![Coverage Status](https://coveralls.io/repos/github/jedib0t/go-passwords/badge.svg?branch=main)](https://coveralls.io/github/jedib0t/go-passwords?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/jedib0t/go-passwords)](https://goreportcard.com/report/github.com/jedib0t/go-passwords)

A high-performance Go library for generating secure passphrases and passwords with extensive customization options.

## Passphrases

Passphrases combine 2+ words with separators, optionally capitalized and with numbers. They're easier to remember than passwords while maintaining security.

**Features:**
- Capitalize words (e.g., `foo` → `Foo`)
- Custom dictionaries or built-in English dictionary
- Configurable word count (2-32 words)
- Optional random number insertion
- Custom separators
- Word length filtering
- **Zero-allocation** via `GenerateTo([]byte)`

### Example
```golang
	g, err := passphrase.NewGenerator(
		passphrase.WithCapitalizedWords(true),
		passphrase.WithDictionary(dictionaries.English()),
		passphrase.WithNumWords(3),
		passphrase.WithNumber(true),
		passphrase.WithSeparator("-"),
		passphrase.WithWordLength(4, 6),
	)
	if err != nil {
		panic(err.Error())
	}
	for i := 1; i <= 10; i++ {
		passphrase, err := g.Generate()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Passphrase #%3d: %#v\n", i, passphrase)
	}
```
<details>
<summary>Output...</summary>
<pre>
Passphrase #  1: "Peage6-Blousy-Whaup"
Passphrase #  2: "Crape0-Natter-Pecs"
Passphrase #  3: "Facers-Razzed-Jupes6"
Passphrase #  4: "Jingko1-Shell-Stupor"
Passphrase #  5: "Nailer-Turgid-Sancta4"
Passphrase #  6: "Rodeo5-Cysts-Pinons"
Passphrase #  7: "Mind-Regina-Swinks9"
Passphrase #  8: "Babas5-Lupous-Xylems"
Passphrase #  9: "Ocreae-Fusel0-Jujube"
Passphrase # 10: "Mirks6-Woofer-Lase"
</pre>
</details>

## Passwords

Generate cryptographically secure random passwords with fine-grained character requirements.

**Features:**
- Custom character sets with ambiguity/duplicate filtering
- Configurable length
- Minimum lower-case character requirements
- Minimum upper-case character requirements
- Symbol count range (min/max)
- **Zero-allocation** via `GenerateTo([]byte)`

### Example
```golang
	g, err := password.NewGenerator(
		password.WithCharset(charset.AllChars.WithoutAmbiguity().WithoutDuplicates()),
		password.WithLength(12),
		password.WithMinLowerCase(5),
		password.WithMinUpperCase(2),
		password.WithNumSymbols(1, 1),
	)
	if err != nil {
		panic(err.Error())
	}
	for i := 1; i <= 10; i++ {
		password, err := g.Generate()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Password #%3d: %#v\n", i, password)
	}
```
<details>
<summary>Output...</summary>
<pre>
Password #  1: "jQwRvL#oye7q"
Password #  2: "T2WRwSbwghc^"
Password #  3: "S@DxkUwkunhy"
Password #  4: "NJ4wxhSygLm&"
Password #  5: "phHfuqw*uAPq"
Password #  6: "$3XDCoLXdeqq"
Password #  7: "enzB*ENGhsQm"
Password #  8: "ioCfs&cLJgyd"
Password #  9: "obwEEEthM$MC"
Password # 10: "kmQVb&fPqexj"
</pre>
</details>

## Enumerator

Systematically enumerate all possible string combinations from a character set and length. Useful for brute-force testing, password cracking research, or exhaustive search scenarios.

**Features:**
- Efficient iteration through all combinations
- Jump to specific positions
- Increment/decrement by N steps
- Optional rollover mode
- Zero-allocation operations (after initial setup)

### Example
```golang
	o := enumerator.New(charset.AlphabetsUpper, 8)

	for i := 1; i <= 10; i++ {
		fmt.Printf("Combination #%3d: %#v\n", i, o.String())
		if o.AtEnd() {
			break
		}
		o.Increment()
	}
```
<details>
<summary>Output...</summary>
<pre>
Combination #  1: "AAAAAAAA"
Combination #  2: "AAAAAAAB"
Combination #  3: "AAAAAAAC"
Combination #  4: "AAAAAAAD"
Combination #  5: "AAAAAAAE"
Combination #  6: "AAAAAAAF"
Combination #  7: "AAAAAAAG"
Combination #  8: "AAAAAAAH"
Combination #  9: "AAAAAAAI"
Combination # 10: "AAAAAAAJ"
</pre>
</details>

## Performance

Benchmarked on AMD Ryzen 9 9950X3D:

| Package | Operation | Time | Allocations |
|---------|-----------|------|-------------|
| **Enumerator** | Increment/Decrement | ~17 ns/op | 0 B/op, 0 allocs/op |
| **Enumerator** | IncrementN/DecrementN | ~104 ns/op | 0 B/op, 0 allocs/op |
| **Enumerator** | String | ~16 ns/op | 0 B/op, 0 allocs/op |
| **Passphrase** | Generate | ~98 ns/op | 24 B/op, 1 allocs/op |
| **Passphrase** | GenerateTo | ~87 ns/op | 0 B/op, 0 allocs/op |
| **Password** | Generate | ~322 ns/op | 64 B/op, 2 allocs/op |
| **Password** | GenerateTo | ~292 ns/op | 0 B/op, 0 allocs/op |
| **RNG** | IntN | ~12 ns/op | 0 B/op, 0 allocs/op |
| **RNG** | Shuffle (Small) | ~108 ns/op | 0 B/op, 0 allocs/op |
| **RNG** | Shuffle (Medium) | ~1190 ns/op | 0 B/op, 0 allocs/op |
| **RNG** | Shuffle (Large) | ~16022 ns/op | 0 B/op, 0 allocs/op |

Run benchmarks: `make bench`
