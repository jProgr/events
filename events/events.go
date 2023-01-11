package events

import "sync"

// facade works as global variable when using a dispatcher as facade.
var facade *Dispatcher

// Dispatcher stores the map of event IDs and their listeners, config, registers new
// listeners and dispatchs events. Main interactor of the package.
type Dispatcher struct {
    listeners map[EventId][]Listener
    config    *Config
    waitGroup *sync.WaitGroup
}

// NewDispatcher setups and creates a new Dispatcher. A default dispatcher:
//
//  Config{
//      isAsync:    false,
//      isFacade:   false,
//      shouldWait: true,
//      waitGroup:  nil,
//  }
//
// When no sync.WaitGroup is provided on the config, the dispatcher creates its own.
// To configure a dispatcher provide one or more functions the receive the config struct.
// The following example configures a dispatcher to work as facade, async and does not
// wait on the listeners to finish work before continuing:
//
//  waitGroup := new(sync.WaitGroup)
//
//  events.NewDispatcher(func (config *events.Config) {
//      config.AsFacade(true)
//      config.ShouldAsync(true)
//      config.ShouldWait(false, waitGroup)
//  })
//
// When creating a default dispatcher, no configuration, it is safe to ignore the error;
// it will be nil.
func NewDispatcher(configurers ...func(*Config)) (*Dispatcher, error) {
    config, err := getConfig(configurers)
    if err != nil {
        return &Dispatcher{}, err
    }

    dispatcher := &Dispatcher{
        listeners: make(map[EventId][]Listener),
        config:    config,
        waitGroup: config.waitGroup,
    }

    if config.isFacade {
        facade = dispatcher
    }

    return dispatcher, nil
}

func getConfig(configurers []func(*Config)) (*Config, error) {
    config := newDefaultConfig()
    for _, configurer := range configurers {
        configurer(config)
    }

    if !config.shouldWait && config.waitGroup == nil {
        return config, newAsyncConfigError("When waiting for goroutines is managed outside the package, a `sync.waitGroup` instance should be provided")
    }

    if config.waitGroup == nil {
        config.waitGroup = new(sync.WaitGroup)
    }

    return config, nil
}

// Register adds to the internal map of event IDs and listeners the arguments provided.
func (dispatcher *Dispatcher) Register(id EventId, listener Listener) *Dispatcher {
    if listeners, ok := dispatcher.listeners[id]; ok {
        dispatcher.listeners[id] = append(listeners, listener)
        return dispatcher
    }

    dispatcher.listeners[id] = []Listener{listener}

    return dispatcher
}

// Dispatch calls all the listeners registered under the event ID of the argument
// provided and passes it to them. If the event ID is not registered, nothing is done.
func (dispatcher *Dispatcher) Dispatch(event Event) *Dispatcher {
    listeners, ok := dispatcher.listeners[event.id]
    if !ok || len(listeners) == 0 {
        return dispatcher
    }

    for _, listener := range listeners {
        if dispatcher.config.isAsync {
            dispatchAsync(dispatcher, event, listener)
            continue
        }

        listener(event)
    }

    if dispatcher.config.shouldWait {
        dispatcher.waitGroup.Wait()
    }

    return dispatcher
}

// dispatchAsync executes listener under a goroutine and updates dispatcher.waitGroup
// accordingly.
func dispatchAsync(dispatcher *Dispatcher, event Event, listener Listener) {
    dispatcher.waitGroup.Add(1)

    go func() {
        defer dispatcher.waitGroup.Done()
        listener(event)
    }()
}

// Register works the same as Dispatcher.Register() but panics if no facade is configured.
func Register(id EventId, listener Listener) *Dispatcher {
    if facade == nil {
        panic("No facade registered")
    }

    return facade.Register(id, listener)
}

// Dispatch works the same as Dispatcher.Dispatch() but panics if no facade is configured.
func Dispatch(event Event) *Dispatcher {
    if facade == nil {
        panic("No facade registered")
    }

    return facade.Dispatch(event)
}
