package eventer

import (
	"sync"
)

// Event is an empty interface that represents any event type.
type Event interface{}

// An EventListener handles events given to it.
// HandleEvent should be thread safe.
type EventListener interface {
	HandleEvent(event Event)
}

// EventEmitter can add / remove EventListener which will be invoked when an event
// is emitted.
type EventEmitter interface {
	// AddListener adds a EventListener that will be invoked when an event is emitted via EmitEvent.
	// The returned value indicates if the listener has been added or not.
	AddListener(EventListener) bool

	// RemoveListener removes a EventListener and prevents it being invoked in subsequent emitted events.
	// The returned value indicates if the listener has been removed or not.
	RemoveListener(EventListener) bool

	// EmitEvent invokes all the EventListener that are registered with the EventEmitter.
	EmitEvent(Event)
}

type sharedEventEmitter struct {
	// listeners is a slice keeping all added EvenListener.
	// The slice is only assinged, but never altered in place.
	listeners []EventListener
	lock      sync.RWMutex
}

func (s *sharedEventEmitter) AddListener(l EventListener) bool {
	if l == nil {
		// do not register nil listener
		return false
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	// check if listener is already registered
	for _, listener := range s.listeners {
		if listener == l {
			return false
		}
	}

	// by making a copy of the underlying array, it will never be changed after creation
	// so it's ok copy the slice while locked but iterate while not locked
	c := make([]EventListener, 0, len(s.listeners)+1)
	c = append(c, s.listeners...)
	c = append(c, l)

	s.listeners = c

	return true
}

func (s *sharedEventEmitter) RemoveListener(l EventListener) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	for i := range s.listeners {
		if s.listeners[i] == l {
			// create a new slice excluding the listener that need to be removed
			c := make([]EventListener, 0, len(s.listeners)-1)
			c = append(c, s.listeners[:i]...)
			c = append(c, s.listeners[i+1:]...)
			s.listeners = c

			return true
		}
	}

	return false
}

// AsyncEventEmitter is an implementation of EventEmitter
// that emits events in separate go routine.
type AsyncEventEmitter struct {
	sharedEventEmitter
}

// EmitEvent will send the event to all the registered listeners.
func (asy *AsyncEventEmitter) EmitEvent(event Event) {
	asy.lock.RLock()
	listeners := asy.listeners
	asy.lock.RUnlock()

	for _, listener := range listeners {
		go listener.HandleEvent(event)
	}
}

// SyncEventEmitter is an implementation of EventEmitter
// that emits events in the same go routine.
type SyncEventEmitter struct {
	sharedEventEmitter
}

// EmitEvent will send the event to all the registered listeners.
func (sy *SyncEventEmitter) EmitEvent(event Event) {
	sy.lock.RLock()
	listeners := sy.listeners
	sy.lock.RUnlock()

	for _, listener := range listeners {
		listener.HandleEvent(event)
	}
}
