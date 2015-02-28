package eve

import (
	"log"
	"sync"
)

type bus struct {
	events       map[string]chan interface{}
	eventsLock   sync.Mutex
	listener     map[string][]chan interface{}
	listenerLock sync.Mutex
	debug        bool
}

func Bus() bus {
	return bus{
		events:   make(map[string]chan interface{}),
		listener: make(map[string][]chan interface{}),
		debug:    false,
	}
}

func (b *bus) trace(v ...interface{}) {
	if b.debug {
		log.Println(v)
	}
}

func (b *bus) Close() {
	for _, v := range b.events {
		close(v)
	}

	for _, v := range b.listener {
		//b.trace("Closing ", v)
		for _, l := range v {
			close(l)
		}
	}
}

func (b *bus) dispatcher(evtName string, c chan interface{}) {

	//TODO: mutex has to be added somewhere here because we iterate over listeners and change them

	for {
		evt, more := <-c
		if !more {
			return
		}
		//dispatch:
		if _, ok := b.listener[evtName]; !ok {
			b.trace("no listener for event: " + evtName)
			continue //no listener registered
		}
		b.trace("dispatching evt: ", evtName)
		for _, l := range b.listener[evtName] {
			//send unblocking:
			select {
			case l <- evt:
				//b.trace("Event dispatched")
			default: //did not send message
			}
		}
		b.trace("clearing listeners for ", evtName)
		//TODO: do i need to close them or what?
		b.listener[evtName] = make([]chan interface{}, 0) //srly 0?
	}
}

//returns a channel to send an event to
func (b *bus) To(evtName string) chan interface{} {
	b.eventsLock.Lock()
	defer b.eventsLock.Unlock()

	if _, ok := b.events[evtName]; !ok {
		b.trace("Registered evt: ", evtName)
		c := make(chan interface{})
		b.events[evtName] = c
		go b.dispatcher(evtName, c)
	}
	return b.events[evtName]
}

//returns a channel to receive events from
func (b *bus) From(evtName string) chan interface{} {
	b.eventsLock.Lock()
	defer b.eventsLock.Unlock()

	c := make(chan interface{})
	if _, ok := b.listener[evtName]; !ok {
		b.listener[evtName] = make([]chan interface{}, 0) //srly 0?
	}

	//b.trace("Registered listener to evt: ", evtName)
	b.listener[evtName] = append(b.listener[evtName], c)
	return c
}

func (b *bus) RegisteredListeners(evtName string) int {
	return len(b.listener[evtName])
}
