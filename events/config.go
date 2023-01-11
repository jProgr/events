package events

import "sync"

// Config holds config information for events.Dispatcher.
// waitGroup can be nil.
type Config struct {
    isAsync    bool
    isFacade   bool
    shouldWait bool
    waitGroup  *sync.WaitGroup
}

// newDefaultConfig builds a new config struct instance.
func newDefaultConfig() *Config {
    return &Config{
        isAsync:    false,
        isFacade:   false,
        shouldWait: true,
        waitGroup:  nil,
    }
}

// ShouldAsync sets whether a dispatcher should execute listeners one after the other
// or execute them all in goroutines.
func (config *Config) ShouldAsync(shouldAsync bool) *Config {
    config.isAsync = shouldAsync

    return config
}

// AsFacade allows a dispatcher to be available package wide by just importing.
func (config *Config) AsFacade(isFacade bool) *Config {
    config.isFacade = isFacade

    return config
}

// ShouldWait controls whether an event triggering (Dispatcher.Dispatch()) should block
// the execution until every listener is done. On shouldWait = true execution is
// blocked until listeners are done, on shouldWait = false execution continues immediatly
// after triggering an event. In the last case, a *sync.WaitGroup should be provided or
// an error will be returned.
func (config *Config) ShouldWait(shouldWait bool, waitGroup *sync.WaitGroup) error {
    if !shouldWait && waitGroup == nil {
        return newAsyncConfigError("When waiting for goroutines is managed outside the package, a `sync.waitGroup` instance should be provided")
    }

    config.shouldWait = shouldWait
    config.waitGroup = waitGroup

    return nil
}

// AsyncConfigError can be found when there is an error setting up a dispatcher.
type AsyncConfigError struct {
    message string
}

// newAsyncConfigError creates a new AsyncConfigError.
func newAsyncConfigError(text string) error {
    return &AsyncConfigError{text}
}

// Error returns the error message.
func (err *AsyncConfigError) Error() string {
    return err.message
}
