package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jedib0t/go-passwords/passphrase"
	"github.com/jedib0t/go-passwords/passphrase/dictionaries"
	"github.com/jedib0t/go-passwords/password"
)

func main() {
	fmt.Println("Passphrases:")
	passphraseGenerator()
	fmt.Println()

	fmt.Println("Passwords:")
	passwordGenerator()
	fmt.Println()

	fmt.Println("Passwords Sequenced:")
	passwordSequencer()
	fmt.Println()

	fmt.Println("Passwords Sequenced & Streamed:")
	passwordSequencerStreaming()
	fmt.Println()
}

func passphraseGenerator() {
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
}

func passwordGenerator() {
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
}

func passwordSequencer() {
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
}

func passwordSequencerStreaming() {
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
}
