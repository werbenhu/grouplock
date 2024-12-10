package grouplock

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	sharedMutex  sync.Mutex
	sleepTime    = time.Millisecond * 10
	numberOfKeys = 1000
)

// sharedLock simulates locking and unlocking a shared mutex.
func sharedLock(wg *sync.WaitGroup) {
	defer wg.Done()

	sharedMutex.Lock()
	time.Sleep(sleepTime)
	sharedMutex.Unlock()
}

// groupLockTest simulates locking and unlocking keys using GroupLock.
func groupLockTest(gl *GroupLock, keys []string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Randomly lock one of the keys
	key := keys[rand.Intn(len(keys))]
	gl.Lock(key)
	// Simulate some work
	time.Sleep(sleepTime)
	gl.Unlock(key)
}

// BenchmarkSharedMutex tests the performance of shared mutex approach.
func BenchmarkSharedMutex(b *testing.B) {
	goroutineCount := 1000
	keys := make([]string, numberOfKeys)

	// Generate numberOfKeys keys
	for i := 0; i < numberOfKeys; i++ {
		keys[i] = "key" + strconv.Itoa(i)
	}

	b.ResetTimer() // Reset the timer to only measure the actual test duration
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		// Start goroutines
		for k := 0; k < goroutineCount; k++ {
			wg.Add(1)
			go sharedLock(&wg)
		}
		wg.Wait()
	}
}

// BenchmarkGroupLock tests the performance of GroupLock approach.
func BenchmarkGroupLock(b *testing.B) {
	goroutineCount := 1000
	keys := make([]string, numberOfKeys)
	gl := New()

	// Generate numberOfKeys keys
	for i := 0; i < numberOfKeys; i++ {
		keys[i] = "key" + strconv.Itoa(i)
	}

	b.ResetTimer() // Reset the timer to only measure the actual test duration
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		// Start goroutines
		for k := 0; k < goroutineCount; k++ {
			wg.Add(1)
			go groupLockTest(gl, keys, &wg)
		}
		wg.Wait()
	}
}
