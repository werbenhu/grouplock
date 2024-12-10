package klocker

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

// groupLockTest simulates locking and unlocking keys using KLocker.
func kLockerTest(kl *KLocker, keys []string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Randomly lock one of the keys
	key := keys[rand.Intn(len(keys))]
	kl.Lock(key)
	// Simulate some work
	time.Sleep(sleepTime)
	kl.Unlock(key)
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

// BenchmarkKeyLocker tests the performance of KLocker approach.
func BenchmarkKeyLocker(b *testing.B) {
	goroutineCount := 1000
	keys := make([]string, numberOfKeys)
	kl := New()

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
			go kLockerTest(kl, keys, &wg)
		}
		wg.Wait()
	}
}
