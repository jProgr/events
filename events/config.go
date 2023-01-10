package events

import "sync"

type Config struct {
    isAsync    bool
    isFacade   bool
    shouldWait bool
    waitGroup  *sync.WaitGroup
}

func newDefaultConfig() *Config {
    return &Config{
        isAsync:    false,
        isFacade:   false,
        shouldWait: true,
        waitGroup:  nil,
    }
}

func (config *Config) ShouldAsync(shouldAsync bool) *Config {
    config.isAsync = shouldAsync

    return config
}

func (config *Config) AsFacade(isFacade bool) *Config {
    config.isFacade = isFacade

    return config
}

func (config *Config) ShouldWait(shouldWait bool, waitGroup *sync.WaitGroup) error {
    if !shouldWait && waitGroup == nil {
        return newAsyncConfigError("When waiting for goroutines is managed outside the package, a `sync.waitGroup` instance should be provided")
    }

    config.shouldWait = shouldWait
    config.waitGroup = waitGroup

    return nil
}

type AsyncConfigError struct {
    message string
}

func newAsyncConfigError(text string) error {
    return &AsyncConfigError{text}
}

func (err *AsyncConfigError) Error() string {
    return err.message
}
