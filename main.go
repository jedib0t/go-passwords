package main

import (
	"fmt"

	"github.com/jedib0t/go-passwords/charset"
	"github.com/jedib0t/go-passwords/enumerator"
	"github.com/jedib0t/go-passwords/passphrase"
	"github.com/jedib0t/go-passwords/passphrase/dictionaries"
	"github.com/jedib0t/go-passwords/password"
)

func main() {
	fmt.Println("Passphrases:")
	demoPassphraseGenerator()
	fmt.Println()

	fmt.Println("Passwords:")
	demoPasswordGenerator()
	fmt.Println()

	fmt.Println("Enumerator:")
	demoEnumerator()
	fmt.Println()
}

func demoPassphraseGenerator() {
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
		phrase, err := g.Generate()
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Passphrase #%3d: %#v\n", idx, phrase)
	}
}

func demoPasswordGenerator() {
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
		pw, err := g.Generate()
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Password #%3d: %#v\n", idx, pw)
	}
}

func demoEnumerator() {
	o := enumerator.New(charset.AlphabetsUpper, 8)

	for idx := 1; idx <= 10; idx++ {
		fmt.Printf("Password #%3d: %#v\n", idx, o.String())

		if o.AtEnd() {
			break
		}
		o.Increment()
	}
}
