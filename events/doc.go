// Package events provides an observer abstraction that allows functions
// to be run when certain events are fired from any part of a codebase.
//
// # Usage
//
// Create a new event dispatcher:
//
//	dispatcher := events.NewDispatcher()
//
// Register event IDs and listeners:
//
//	const SomeEventId events.EventId = "Some description"
//
//	listener := func(event events.Event) {
//	    eventData := event.Get().(*SomeType)
//
//	    // Do something with eventData
//	}
//	dispatcher.Register(SomeEventId, listener)
//
// It is possible to register multiple listeners for the same event ID.
// A listener is just a function that receives an events.Event and returns
// nothing. Inside the listener a type assertion is needed due to
// events.Event.Get() returning any. eventData will have anything that
// was passed to the dispatcher when triggering the event.
//
// Then one can dispatch an event using:
//
//	event := events.Make(SomeEventId, &someData)
//	dispatcher.Dispatch(event)
//
// Listeners will be executed one after the other in the registred order
// (or with goroutines if configured that way). someData will be available
// under the events.Event passed to the listener.
//
// # Usage as facade
//
// Although not recomended, a dispatcher can also be available in the package
// as facade. Registering listeners and dispatchers work the same but they are
// available package wide; useful for quick prototypes to avoid passing the
// dispatcher too deeply the call chain:
//
//	package one
//
//	import "github.com/jProgr/events"
//
//	func f() {
//	    events.NewDispatcher(func(config *events.Config) {
//	        // This makes this dispatcher to be stored in
//	        // the package and makes it available by just
//	        // importing events.
//	        config.AsFacade(true)
//	    })
//
//	    // Register directly on the package, without calling
//	    // the dispatcher.
//	    events.Register(SomeEventId, someListener)
//	}
//
//	package two
//
//	import "github.com/jProgr/events"
//
//	func g() {
//	    // Dispatch directly by just importing the package
//	    events.Dispatch(events.Make(SomeEventId, &someData))
//	}
//
// In the example, f() should be run before g() for everything to work. Calling
// Register() or Dispatch() directly on the package without configuring a dispatcher
// to work as facade will raise a panic().
//
// # Async execution
//
// The default execution order of listeners is just one after the other in the goroutine
// where the event was triggered (usually the main one). If the listeners are to do slow
// work (usually network stuff), one can configure the dispatcher to run each listener in
// its own goroutine. There are two modes:
//
//   - Trigger the event and wait for every listener to finish work before continuing.
//   - Trigger the event, launch the listeners and continue work without waiting on the
//     goroutines.
//
// To wait on every goroutine to finish:
//
//	dispatcher := events.NewDispatcher(func(config *events.Config) {
//	    config.ShouldAsync(true)
//	})
//
// To continue immediately after triggering an event:
//
//	waitGroup := new(sync.WaitGroup)
//	dispatcher := events.NewDispatcher(func(config *events.Config) {
//	    config.ShouldAsync(true)
//	    config.ShouldWait(false, waitGroup)
//	})
//
// Failing to provide a sync.WaitGroup instance will result in an error of type
// events.AsyncConfigError. On this mode the caller is responsible for managing the WaitGroup
// to avoid the main goroutine finishing before the listeners are done.
package events
