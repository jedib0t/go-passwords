# go-passwords

[![Go Reference](https://pkg.go.dev/badge/github.com/jedib0t/go-passwords/v0.svg)](https://pkg.go.dev/github.com/jedib0t/go-passwords)
[![Build Status](https://github.com/jedib0t/go-passwords/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/jedib0t/go-passwords/actions?query=workflow%3ACI+event%3Apush+branch%3Amain)
[![Coverage Status](https://coveralls.io/repos/github/jedib0t/go-passwords/badge.svg?branch=main)](https://coveralls.io/github/jedib0t/go-passwords?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/jedib0t/go-passwords)](https://goreportcard.com/report/github.com/jedib0t/go-passwords)

Passphrase & Password generation library for GoLang.

## Passphrases

Passphrases are made up of 2 or more words connected by a separator and may have
capitalized words, and numbers. These are easier for humans to remember compared
to passwords.

The `passphrase` package helps generate these and supports the following rules
that be used during generation:
* Capitalize words used in the passphrase (foo -> Foo)
* Use a custom dictionary of words instead of built-in English dictionary
* Use X number of Words
* Insert a random number behind one of the words
* Use a custom separator
* Use words with a specific length-range

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
	for idx := 1; idx <= 10; idx++ {
		fmt.Printf("Passphrase #%3d: %#v\n", idx, g.Generate())
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

Passwords are a random amalgamation of characters.

The `password` package helps generate these and supports the following rules
that be used during generation:
* Use a specific character-set
* Restrict the length of the password
* Use *at least* X lower-case characters
* Use *at least* X upper-case characters
* Use *at least* X and *at most* Y symbols

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
	for idx := 1; idx <= 10; idx++ {
		fmt.Printf("Password #%3d: %#v\n", idx, g.Generate())
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

Enumerator helps generate all possible string combinations of characters given a
list of characters and the expected length of the string.

The `enumerator` package provides optimal interfaces to move through the list:
* Decrement()
* DecrementN(n)
* GoTo(n)
* Increment()
* IncrementN(n)
* etc.

### Example
```golang
	o := enumerator.New(charset.AlphabetsUpper, 8)

	for idx := 1; idx <= 10; idx++ {
		fmt.Printf("Password #%3d: %#v\n", idx, o.String())

		if o.AtEnd() {
			break
		}
		o.Increment()
	}
```
<details>
<summary>Output...</summary>
<pre>
Password #  1: "AAAAAAAA"
Password #  2: "AAAAAAAB"
Password #  3: "AAAAAAAC"
Password #  4: "AAAAAAAD"
Password #  5: "AAAAAAAE"
Password #  6: "AAAAAAAF"
Password #  7: "AAAAAAAG"
Password #  8: "AAAAAAAH"
Password #  9: "AAAAAAAI"
Password # 10: "AAAAAAAJ"
</pre>
</details>

## Benchmarks
```
go test -bench=. -benchmem ./enumerator ./passphrase ./password
goos: linux
goarch: amd64
pkg: github.com/jedib0t/go-passwords/enumerator
cpu: AMD Ryzen 9 9950X3D 16-Core Processor          
BenchmarkEnumerator_Decrement-12                68338802                17.73 ns/op            0 B/op          0 allocs/op
BenchmarkEnumerator_Decrement_Big-12            67899237                17.63 ns/op            0 B/op          0 allocs/op
BenchmarkEnumerator_DecrementN-12               10151359               114.5 ns/op             0 B/op          0 allocs/op
BenchmarkEnumerator_GoTo-12                      8441791               138.2 ns/op            40 B/op          2 allocs/op
BenchmarkEnumerator_Increment-12                69405856                17.29 ns/op            0 B/op          0 allocs/op
BenchmarkEnumerator_Increment_Big-12            70501692                17.36 ns/op            0 B/op          0 allocs/op
BenchmarkEnumerator_IncrementN-12               11806892               109.6 ns/op             0 B/op          0 allocs/op
BenchmarkEnumerator_String-12                   74139985                16.08 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/jedib0t/go-passwords/enumerator      10.132s
goos: linux
goarch: amd64
pkg: github.com/jedib0t/go-passwords/passphrase
cpu: AMD Ryzen 9 9950X3D 16-Core Processor          
BenchmarkGenerator_Generate-12           6747350               159.6 ns/op           144 B/op          5 allocs/op
PASS
ok      github.com/jedib0t/go-passwords/passphrase      1.298s
goos: linux
goarch: amd64
pkg: github.com/jedib0t/go-passwords/password
cpu: AMD Ryzen 9 9950X3D 16-Core Processor          
BenchmarkGenerator_Generate-12           9879747               116.9 ns/op            40 B/op          2 allocs/op
PASS
ok      github.com/jedib0t/go-passwords/password        1.284s
```
