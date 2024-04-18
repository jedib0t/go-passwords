# go-passwords

[![Go Reference](https://pkg.go.dev/badge/github.com/jedib0t/go-passwords/v0.svg)](https://pkg.go.dev/github.com/jedib0t/go-passwords)
[![Build Status](https://github.com/jedib0t/go-passwords/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/jedib0t/go-passwords/actions?query=workflow%3ACI+event%3Apush+branch%3Amain)
[![Coverage Status](https://coveralls.io/repos/github/jedib0t/go-passwords/badge.svg?branch=main)](https://coveralls.io/github/jedib0t/go-passwords?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/jedib0t/go-passwords)](https://goreportcard.com/report/github.com/jedib0t/go-passwords)

Password generation library for GoLang.

## Benchmarks
```
$ go test -bench=. -benchmem ./password
goos: linux
goarch: amd64
pkg: github.com/jedib0t/go-passwords/password
cpu: AMD Ryzen 9 5900X 12-Core Processor
BenchmarkGenerator_Generate-12    	 6245260	       188.2 ns/op	      40 B/op	       2 allocs/op
BenchmarkSequencer_GotoN-12       	 4359440	       270.6 ns/op	      32 B/op	       3 allocs/op
BenchmarkSequencer_Next-12        	13632730	        83.67 ns/op	      16 B/op	       1 allocs/op
BenchmarkSequencer_NextN-12       	 6608569	       181.5 ns/op	      32 B/op	       3 allocs/op
BenchmarkSequencer_Prev-12        	13509426	        87.51 ns/op	      16 B/op	       1 allocs/op
BenchmarkSequencer_PrevN-12       	 4266948	       276.8 ns/op	      32 B/op	       3 allocs/op
PASS
ok  	github.com/jedib0t/go-passwords/password	8.178s
```

## Usage

### Random Passwords
```golang
	generator, err := password.NewGenerator(
		password.WithCharset(password.AllChars.WithoutAmbiguity().WithoutDuplicates()),
		password.WithLength(12),
		password.WithMinLowerCase(5),
		password.WithMinpperCase(2),
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
Password #  1: "CmHii4zek_wU"
Password #  2: "m+GicmQEnxkz"
Password #  3: "FruTpkprFNR="
Password #  4: "p@xjqBH3bbfi"
Password #  5: "D(WadeVLTfhm"
Password #  6: "uLpSFv%pMidL"
Password #  7: "bbBQ*gypmhTx"
Password #  8: "abshu4}HNpwE"
Password #  9: "UjGpDsP{4mfi"
Password # 10: "toKue=dvUPzz"
</pre>
</details>

## Sequential Passwords

### In a Loop
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

### Streamed (for async processing)
```golang
	sequencer, err := password.NewSequencer(
		password.WithCharset(password.Charset("AB")),
		password.WithLength(5),
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
Password #  1: "AAAAA"
Password #  2: "AAAAB"
Password #  3: "AAABA"
Password #  4: "AAABB"
Password #  5: "AABAA"
Password #  6: "AABAB"
Password #  7: "AABBA"
Password #  8: "AABBB"
Password #  9: "ABAAA"
Password # 10: "ABAAB"
Password # 11: "ABABA"
Password # 12: "ABABB"
Password # 13: "ABBAA"
Password # 14: "ABBAB"
Password # 15: "ABBBA"
Password # 16: "ABBBB"
Password # 17: "BAAAA"
Password # 18: "BAAAB"
Password # 19: "BAABA"
Password # 20: "BAABB"
Password # 21: "BABAA"
Password # 22: "BABAB"
Password # 23: "BABBA"
Password # 24: "BABBB"
Password # 25: "BBAAA"
Password # 26: "BBAAB"
Password # 27: "BBABA"
Password # 28: "BBABB"
Password # 29: "BBBAA"
Password # 30: "BBBAB"
Password # 31: "BBBBA"
Password # 32: "BBBBB"
</pre>
</details>
