package events

import (
    "sync"
    "testing"
)

func TestTogglesAsync(test *testing.T) {
    config := newDefaultConfig()

    config.ShouldAsync(true)
    if !config.isAsync {
        test.Fatal("`config.isAsync` should be `true`")
    }

    config.ShouldAsync(false)
    if config.isAsync {
        test.Fatal("`config.isAsync` should be `false`")
    }
}

func TestTogglesFacade(test *testing.T) {
    config := newDefaultConfig()

    config.AsFacade(true)
    if !config.isFacade {
        test.Fatal("`config.isFacade` should be `true`")
    }

    config.AsFacade(false)
    if config.isFacade {
        test.Fatal("`config.isFacade` should be `false`")
    }
}

func TestTogglesWaiting(test *testing.T) {
    config := newDefaultConfig()

    err := config.ShouldWait(true, nil)
    if err != nil {
        test.Fatal("`err` should be `nil` on `ShouldWait(true, nil)`")
    }
    if !config.shouldWait || config.waitGroup != nil {
        test.Fatal("`config.shouldWait` should be `true` and `config.waitGroup` should be `nil`")
    }

    err = config.ShouldWait(true, new(sync.WaitGroup))
    if err != nil {
        test.Fatal("`err` should be `nil` on `ShouldWait(true, new(sync.WaitGroup))`")
    }
    if !config.shouldWait || config.waitGroup == nil {
        test.Fatal("`config.shouldWait` should be `true` and `config.waitGroup` shouldn't be `nil`")
    }

    err = config.ShouldWait(false, nil)
    if _, ok := err.(*AsyncConfigError); !ok {
        test.Fatal("`err` should be a `AsyncConfigError`")
    }

    err = config.ShouldWait(false, new(sync.WaitGroup))
    if err != nil {
        test.Fatal("`err` should be `nil` on `ShouldWait(false, new(sync.WaitGroup))`")
    }
    if config.shouldWait || config.waitGroup == nil {
        test.Fatal("`config.shouldWait` should be `false` and `config.waitGroup` shouldn't be `nil`")
    }
}
