package mutex

// MutexDriver is an interface
type MutexDriver interface {
	// Read and Write latest marker
	NewMutex(name string)
	Lock() error
	Unlock() (bool, error)
}
