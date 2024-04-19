package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-passwords/passphrase"
	"github.com/jedib0t/go-passwords/passphrase/dictionaries"
)

var (
	flagCapitalized    = flag.Bool("capitalize", false, "Capitalize all words?")
	flagDictionary     = flag.String("dictionary", "", "Path to dictionary file (default: built-in English words)")
	flagHelp           = flag.Bool("help", false, "Display this Help text")
	flagNumPassphrases = flag.Int("num-passphrases", 10, "Number of passphrases to generate")
	flagNumWords       = flag.Int("num-words", 3, "Number of words in Passphrase")
	flagPrintIndex     = flag.Bool("index", false, "Print Index Value (1-indexed)")
	flagSeed           = flag.Uint64("seed", 0, "Seed value for non-sequenced mode (ignored if zero)")
	flagWithNumber     = flag.Bool("with-number", false, "Inject random number suffix?")
)

func main() {
	flag.Parse()
	if *flagHelp {
		printHelp()
		os.Exit(0)
	}

	dictionary := dictionaries.English()
	if *flagDictionary != "" {
		bytes, err := os.ReadFile(*flagDictionary)
		if err != nil {
			panic(err.Error())
		}
		dictionary = strings.Split(string(bytes), "\n")
	}
	numWords := *flagNumWords
	if numWords < 2 {
		fmt.Printf("ERROR: value of --num-words way too low: %#v\n", numWords)
		os.Exit(1)
	} else if numWords > 16 {
		fmt.Printf("ERROR: value of --num-words way too high: %#v\n", numWords)
		os.Exit(1)
	}
	printIndex := *flagPrintIndex
	seed := *flagSeed
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}

	generator, err := passphrase.NewGenerator(
		passphrase.WithCapitalizedWords(*flagCapitalized),
		passphrase.WithDictionary(dictionary),
		passphrase.WithNumWords(numWords),
		passphrase.WithNumber(*flagWithNumber),
	)
	if err != nil {
		fmt.Printf("ERROR: failed to instantiate generator: %v\n", err)
		os.Exit(1)
	}
	generator.SetSeed(seed)

	for idx := 0; idx < *flagNumPassphrases; idx++ {
		if printIndex {
			fmt.Printf("%d\t", idx+1)
		}
		fmt.Println(generator.Generate())
	}
}

func printHelp() {
	fmt.Println(`passphrase-generator: Generate random passphrases.

Examples:
  * passphrase-generator
    // generate 10 passphrases
  * passphrase-generator --capitalize --num-words=4 --with-number
    // generate 10 passphrases with 4 words capitalized and with a number in there

Flags:`)
	flag.PrintDefaults()
}
