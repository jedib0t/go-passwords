# go-passwords

[![Go Reference](https://pkg.go.dev/badge/github.com/jedib0t/go-passwords/v0.svg)](https://pkg.go.dev/github.com/jedib0t/go-passwords)
[![Build Status](https://github.com/jedib0t/go-passwords/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/jedib0t/go-passwords/actions?query=workflow%3ACI+event%3Apush+branch%3Amain)
[![Coverage Status](https://coveralls.io/repos/github/jedib0t/go-passwords/badge.svg?branch=main)](https://coveralls.io/github/jedib0t/go-passwords?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/jedib0t/go-passwords)](https://goreportcard.com/report/github.com/jedib0t/go-passwords)

Passphrase & Password generation library for GoLang.

## Passphrases
```golang
	generator, err := passphrase.NewGenerator(
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
		fmt.Printf("Passphrase #%3d: %#v\n", idx, generator.Generate())
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
```golang
	generator, err := password.NewGenerator(
		password.WithCharset(password.AllChars.WithoutAmbiguity().WithoutDuplicates()),
		password.WithLength(12),
		password.WithMinLowerCase(5),
		password.WithMinUpperCase(2),
		password.WithNumSymbols(1, 1),
	)
	if err != nil {
		panic(err.Error())
	}
	for idx := 1; idx <= 10; idx++ {
		fmt.Printf("Password #%3d: %#v\n", idx, generator.Generate())
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

### Sequential Passwords

```golang
	sequencer, err := password.NewSequencer(
		password.WithCharset(password.AllChars.WithoutAmbiguity()),
		password.WithLength(8),
	)
	if err != nil {
		panic(err.Error())
	}
	for idx := 1; idx <= 10; idx++ {
		fmt.Printf("Password #%3d: %#v\n", idx, sequencer.Get())

		if !sequencer.HasNext() {
			break
		}
		sequencer.Next()
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
Password #  9: "AAAAAAAJ"
Password # 10: "AAAAAAAK"
</pre>
</details>

#### Streamed (for async processing)
```golang
	sequencer, err := password.NewSequencer(
		password.WithCharset(password.Charset("AB")),
		password.WithLength(4),
	)
	if err != nil {
		panic(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	chPasswords := make(chan string, 1)
	go func() {
		err := sequencer.Stream(ctx, chPasswords)
		if err != nil {
			panic(err.Error())
		}
	}()

	idx := 0
	for {
		select {
		case <-ctx.Done():
			panic("timed out")
		case pw, ok := <-chPasswords:
			if !ok {
				return
			}
			idx++
			fmt.Printf("Password #%3d: %#v\n", idx, pw)
		}
	}
```
<details>
<summary>Output...</summary>
<pre>
Password #  1: "AAAA"
Password #  2: "AAAB"
Password #  3: "AABA"
Password #  4: "AABB"
Password #  5: "ABAA"
Password #  6: "ABAB"
Password #  7: "ABBA"
Password #  8: "ABBB"
Password #  9: "BAAA"
Password # 10: "BAAB"
Password # 11: "BABA"
Password # 12: "BABB"
Password # 13: "BBAA"
Password # 14: "BBAB"
Password # 15: "BBBA"
Password # 16: "BBBB"
</pre>
</details>

## Benchmarks
```
goos: linux
goarch: amd64
pkg: github.com/jedib0t/go-passwords/passphrase
cpu: AMD Ryzen 9 5900X 12-Core Processor            
BenchmarkGenerator_Generate-12    	 2862954	       393.5 ns/op	     167 B/op	       8 allocs/op
PASS
ok  	github.com/jedib0t/go-passwords/passphrase	1.567s

goos: linux
goarch: amd64
pkg: github.com/jedib0t/go-passwords/password
cpu: AMD Ryzen 9 5900X 12-Core Processor            
BenchmarkGenerator_Generate-12    	 6413606	       185.3 ns/op	      40 B/op	       2 allocs/op
BenchmarkSequencer_GotoN-12       	 4353010	       272.5 ns/op	      32 B/op	       3 allocs/op
BenchmarkSequencer_Next-12        	13955396	        84.61 ns/op	      16 B/op	       1 allocs/op
BenchmarkSequencer_NextN-12       	 6473270	       183.9 ns/op	      32 B/op	       3 allocs/op
BenchmarkSequencer_Prev-12        	13106161	        87.22 ns/op	      16 B/op	       1 allocs/op
BenchmarkSequencer_PrevN-12       	 3967755	       288.8 ns/op	      32 B/op	       3 allocs/op
PASS
ok  	github.com/jedib0t/go-passwords/password	8.192s
```
