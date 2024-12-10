package klocker

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKeyLocker_LockUnlock(t *testing.T) {
	// Initialize KeyLocker
	kl := New()

	// Test Lock and Unlock for a single key
	key := "user1"

	// Lock the key
	kl.Lock(key)

	// Unlock the key
	kl.Unlock(key)

	// Ensure the lock can be used again after unlocking (test with a new lock)
	kl.Lock(key)

	// Unlock again
	kl.Unlock(key)

	// Clean up
	kl.Stop()
}

func TestKeyLocker_LockMultipleUsers(t *testing.T) {
	// Initialize KeyLocker
	kl := New()

	// Define multiple user keys
	keys := []string{"user1", "user2", "user3"}

	// Lock each key in parallel
	var wg sync.WaitGroup
	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			kl.Lock(k)
			time.Sleep(100 * time.Millisecond)
			kl.Unlock(k)
		}(key)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Ensure all locks were released successfully
	for _, key := range keys {
		kl.Lock(key) // Should not block
		kl.Unlock(key)
	}

	// Clean up
	kl.Stop()
}

func TestKeyLocker_AutomaticCleanup(t *testing.T) {
	// Initialize KeyLocker with cleanup interval of 1 second for testing
	kl := New(WithInterval(1 * time.Second))

	// Lock some keys
	keys := []string{"user1", "user2", "user3"}
	for _, key := range keys {
		kl.Lock(key)
	}

	// Sleep for 2 seconds to let the cleaner run
	time.Sleep(2 * time.Second)

	// Verify that the locks are cleaned up after being unlocked
	for _, key := range keys {
		kl.Unlock(key)
	}

	// Sleep to ensure cleanup happens
	time.Sleep(2 * time.Second)

	// Verify that all keys are cleaned up
	kl.locks.Range(func(key, value interface{}) bool {
		t.Errorf("Key %v still exists in the lock map", key)
		return true
	})

	// Clean up
	kl.Stop()
}

func TestKeyLocker_LockCleanupAfterUnlock(t *testing.T) {
	// Initialize KeyLocker with a short cleanup interval
	kl := New(WithInterval(1 * time.Second))

	// Lock a key
	key := "user1"
	kl.Lock(key)

	// Unlock the key
	kl.Unlock(key)

	// Sleep for 2 seconds to allow the cleaner to run
	time.Sleep(2 * time.Second)

	// Verify that the lock is cleaned up after the unlock
	_, loaded := kl.locks.Load(key)
	assert.False(t, loaded, "Lock for key %s should be cleaned up", key)

	// Clean up
	kl.Stop()
}

func TestKeyLocker_StopCleaner(t *testing.T) {
	// Initialize KeyLocker
	kl := New(WithInterval(1 * time.Second))

	// Lock a key
	key := "user1"
	kl.Lock(key)

	// Stop the cleaner
	kl.Stop()

	// Unlock the key after stopping the cleaner
	kl.Unlock(key)

	// Verify that the lock was unlocked
	_, loaded := kl.locks.Load(key)
	assert.True(t, loaded, "Lock for key %s should exist after unlock", key)
}

func TestKeyLocker_MultipleLocksOnSameKey(t *testing.T) {
	// Initialize KeyLocker
	kl := New(WithInterval(1 * time.Second))

	// Lock the same key in two different goroutines
	key := "user1"
	var wg sync.WaitGroup

	wg.Add(2)

	// First goroutine to lock the key
	go func() {
		defer wg.Done()
		kl.Lock(key)
		defer kl.Unlock(key)
	}()

	// Second goroutine to lock the key
	go func() {
		defer wg.Done()
		kl.Lock(key)
		defer kl.Unlock(key)
	}()

	// Wait for both goroutines to finish
	wg.Wait()

	// Sleep for 2 seconds to allow the cleaner to run
	time.Sleep(2 * time.Second)

	// Ensure the lock is released after both unlocks
	// Check that the lock item is cleaned up after both unlocks
	_, loaded := kl.locks.Load(key)
	assert.False(t, loaded, "Lock for key %s should be cleaned up after both unlocks", key)

	// Clean up
	kl.Stop()
}
