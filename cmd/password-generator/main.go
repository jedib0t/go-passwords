package main

import (
	"flag"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/jedib0t/go-passwords/password"
)

var (
	flagCharset     = flag.String("charset", string(password.AlphaNumeric), "Charset for Passwords")
	flagCount       = flag.String("count", "10", "Number of Passwords to generate (0=Max-Possible, 1<=x<=100 for non-sequenced)")
	flagHelp        = flag.Bool("help", false, "Display this Help text")
	flagNoAmbiguity = flag.Bool("no-ambiguity", false, "Avoid Ambiguous Characters?")
	flagNumChars    = flag.Int("num-chars", 12, "Number of characters in Password")
	flagPrintIndex  = flag.Bool("index", false, "Print Index Value (1-indexed)")
	flagSeed        = flag.Int64("seed", 0, "Seed value for non-sequenced mode (ignored if zero)")
	flagSequenced   = flag.Bool("sequenced", false, "Generate passwords in a sequence")
	flagStartIdx    = flag.String("start", "0", "Index to start from in the sequence (1-indexed)")

	numZero    = big.NewInt(0)
	numOne     = big.NewInt(1)
	numHundred = big.NewInt(100)
)

func main() {
	flag.Parse()
	if *flagHelp {
		printHelp()
		os.Exit(0)
	}

	charset := password.Charset(*flagCharset).WithoutDuplicates()
	if *flagNoAmbiguity {
		charset = charset.WithoutAmbiguity()
	}
	count, ok := new(big.Int).SetString(*flagCount, 10)
	if !ok {
		fmt.Printf("ERROR: failed to parse value of flag --count: %#v\n", *flagCount)
		os.Exit(1)
	}
	numChars := *flagNumChars
	if numChars > 64 {
		fmt.Printf("ERROR: value of --num-chars way too high: %#v\n", numChars)
		os.Exit(1)
	}
	printIndex := *flagPrintIndex
	seed := *flagSeed
	if seed <= 0 {
		seed = time.Now().UnixNano()
	}
	startIdx, ok := new(big.Int).SetString(*flagStartIdx, 10)
	if !ok {
		fmt.Printf("ERROR: failed to parse value of flag --start: %#v\n", *flagStartIdx)
		os.Exit(1)
	}
	if startIdx.Cmp(numZero) > 0 {
		startIdx = startIdx.Sub(startIdx, numOne) // startIdx is 0-indexed
	}

	if *flagSequenced {
		if count.Cmp(numZero) == 0 {
			count.Set(password.MaximumPossibleWords(charset, numChars))
		}
		generateSequencedPasswords(charset, numChars, count, startIdx, printIndex)
	} else {
		if count.Cmp(numZero) == 0 {
			count.Set(numOne)
		} else if count.Cmp(numHundred) > 0 {
			count.Set(numHundred)
		}
		generateRandomPasswords(charset, numChars, count, printIndex, seed)
	}
}

func generateRandomPasswords(charset password.Charset, numChars int, count *big.Int, printIndex bool, seed int64) {
	generator, err := password.NewGenerator(charset, numChars)
	if err != nil {
		fmt.Printf("ERROR: failed to instantiate generator: %v\n", err)
		os.Exit(1)
	}
	generator.SetSeed(seed)

	for idx := big.NewInt(1); count.Cmp(idx) >= 0; idx = idx.Add(idx, numOne) {
		if printIndex {
			fmt.Printf("%s\t", idx.String())
		}
		fmt.Println(generator.Generate())
	}
}

func generateSequencedPasswords(charset password.Charset, numChars int, count *big.Int, startIdx *big.Int, printIndex bool) {
	sequencer, err := password.NewSequencer(charset, numChars)
	if err != nil {
		fmt.Printf("ERROR: failed to instantiate generator: %v\n", err)
		os.Exit(1)
	}

	_, err = sequencer.GotoN(startIdx)
	if err != nil {
		fmt.Printf("ERROR: failed to go to sequence start @ %v: %v\n", startIdx, err)
		os.Exit(1)
	}

	for idx := big.NewInt(1); count.Cmp(idx) >= 0; idx = idx.Add(idx, numOne) {
		pw := sequencer.Get()

		if printIndex {
			idxPlusStart := new(big.Int).Set(startIdx)
			idxPlusStart.Add(idxPlusStart, idx)
			fmt.Printf("%s\t", idxPlusStart.String())
		}
		fmt.Println(pw)

		if !sequencer.HasNext() {
			break
		}
		sequencer.Next()
	}
}

func printHelp() {
	fmt.Println(`password-generator: Generate random/sequential passwords.

Examples:
  * password-generator
    // generate 10 random passwords with default charset
  * password-generator --sequenced --charset "AB" --count 0 --num-chars 5
    // generate all possible five-character passwords in sequence with letters 'A' & 'B'
  * password-generator --sequenced --charset "AB" --count 10 --num-chars 5 --start 15
    // generate 10 five-character passwords from 15th in the sequence with letters 'A' & 'B'
  * password-generator --sequenced --charset "AB" --count 10 --num-chars 5 --start 15 --index
    // generate 10 five-character passwords from 15th in the sequence with letters 'A' & 'B' with index

Flags:`)
	flag.PrintDefaults()
}
