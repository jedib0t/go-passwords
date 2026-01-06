package rng

import (
	"crypto/rand"
	"sync"
)

const (
	// bufferSize is the size of the random bytes buffer.
	// Larger buffers reduce syscalls but use more memory.
	bufferSize = 1024
)

var (
	// randBuffer holds buffered random bytes
	randBuffer [bufferSize]byte
	// randPos tracks the current position in the buffer
	randPos int
	// randMutex protects access to the buffer
	randMutex sync.Mutex
)

// readBytesBuffered reads the requested number of bytes from the buffered crypto/rand.
// It automatically refills the buffer when needed.
func readBytesBuffered(b []byte) error {
	randMutex.Lock()
	defer randMutex.Unlock()

	needed := len(b)
	available := bufferSize - randPos

	if available < needed {
		// Not enough bytes in buffer, refill it
		if _, err := rand.Read(randBuffer[:]); err != nil {
			return err
		}
		randPos = 0
		available = bufferSize
	}

	// Copy bytes from buffer
	copy(b, randBuffer[randPos:randPos+needed])
	randPos += needed

	return nil
}
