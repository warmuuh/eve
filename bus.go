package eve

import (
	cmap "github.com/streamrail/concurrent-map"
	"log"
	"sync"
)

type bus struct {
	events        cmap.ConcurrentMap //map[string]chan interface{}
	eventsMutex   sync.Mutex
	listener      cmap.ConcurrentMap //map[string][]chan interface{}
	listenerMutex sync.Mutex
	debug         bool
}

func Bus() bus {
	return bus{
		events:   cmap.New(), //make(map[string]chan interface{}),
		listener: cmap.New(), //make(map[string][]chan interface{}),
		debug:    false,
	}
}

func (b *bus) trace(v ...interface{}) {
	if b.debug {
		log.Println(v)
	}
}

func (b *bus) Close() {
	//for _, v := range b.events {
	//	close(v)
	//}

	//for _, v := range b.listener {
	//	b.trace("Closing ", v)
	//	for _, l := range v {
	//		close(l)
	//	}
	//}

}
func (b *bus) dispatcher(evtName string, c chan interface{}) {
	for {
		evt, more := <-c
		if !more {
			return
		}
		//dispatch:
		obj, ok := b.listener.Get(evtName)
		if !ok {
			b.trace("no listener for event: " + evtName)
			continue //no listener registered
		}
		b.trace("dispatching evt: ", evtName)
		listeners := obj.([]chan interface{})
		log.Println("Registered listeners: ", len(listeners))
		for _, l := range listeners {
			//send unblocking:
			select {
			case l <- evt:
				b.trace("Event dispatched")
			default: //did not send message
			}
		}
		b.trace("clearing listeners for ", evtName)
		//TODO: do i need to close them or what?
		b.listener.Set(evtName, make([]chan interface{}, 0)) //srly 0?
	}
}

func (b *bus) To(evtName string) chan interface{} {

	b.eventsMutex.Lock()
	defer b.eventsMutex.Unlock()

	var c chan interface{}
	obj, ok := b.events.Get(evtName)
	if ok {
		c = obj.(chan interface{})
	} else {
		b.trace("Registered evt: ", evtName)
		c = make(chan interface{})
		b.events.Set(evtName, c)
		go b.dispatcher(evtName, c)
	}

	return c
}

func (b *bus) From(evtName string) chan interface{} {
	c := make(chan interface{})

	b.listenerMutex.Lock()
	defer b.listenerMutex.Unlock()

	var listeners []chan interface{}
	obj, ok := b.listener.Get(evtName)
	if ok {
		listeners = obj.([]chan interface{})
	} else {
		listeners = make([]chan interface{}, 0) //srly 0?
	}

	b.trace("Registered listener to evt: ", evtName)
	b.listener.Set(evtName, append(listeners, c))
	return c
}

func (b *bus) RegisteredListeners(evtName string) int {
	obj, ok := b.listener.Get(evtName)
	if !ok {
		return 0
	}
	listeners := obj.([]chan interface{})
	return len(listeners)
}
