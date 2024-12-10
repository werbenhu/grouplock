# grouplock


`grouplock` is a Go package for key-based locks with automatic cleanup.

`grouplock` provides a mechanism for managing locks on keys with support for automatic cleanup of unused locks. It allows you to acquire and release locks for specific keys and automatically cleans up locks that are no longer in use.

## Features

- Lock management for keys with reference counting.
- Automatic cleanup of unused locks at configurable intervals.

## Installation

To use the `grouplock` package, import it into your Go project:

```go
import "github.com/werbenhu/grouplock"
```

### Usage

#### Creating a New GroupLock

You can create a new GroupLock with an optional cleanup interval (default is 30 minutes):


```go
// Default interval 30 minutes
gl := grouplock.New() 

// Custom interval
gl := grouplock.New(grouplock.WithCleanInterval(time.Hour)) 
```

#### Locking and Unlocking Keys
To acquire and release the lock for a specific key:

```go
gl.Lock("key")
gl.Unlock("key")
```

### Stopping the Cleaner

To stop the cleaner goroutine that periodically cleans up unused locks:
```go
gl.Stop()
```

### Options

`WithCleanInterval(interval time.Duration)`: Set the interval for automatic lock cleanup.

### Internal Working

- Each lock is associated with a reference count. The lock is only removed when there are no active references.
- A background goroutine periodically checks for unused locks and removes them.
