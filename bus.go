package eve

import (
	"log"
)

type bus struct {
	events   map[string]chan interface{}
	listener map[string][]chan interface{}
}

func Bus() bus {
	return bus{
		events:   make(map[string]chan interface{}),
		listener: make(map[string][]chan interface{}),
	}
}

func (b *bus) dispatcher(evtName string, c chan interface{}) {
	for {
		evt := <-c
		//dispatch:
		if _, ok := b.listener[evtName]; !ok {
			return //no listener registered
		}
		log.Println("dispatching evt: ", evtName)
		for _, l := range b.listener[evtName] {
			//send unblocking:
			select {
			case l <- evt:
				log.Println("Event dispatched")
			default: //did not send message
			}
		}
		log.Println("clearing listeners for ", evtName)
		b.listener[evtName] = make([]chan interface{}, 5)
	}
}

func (b *bus) To(evtName string) chan interface{} {
	if _, ok := b.events[evtName]; !ok {
		log.Println("Registered evt: ", evtName)
		c := make(chan interface{})
		b.events[evtName] = c
		go b.dispatcher(evtName, c)
	}
	return b.events[evtName]
}

func (b *bus) From(evtName string) chan interface{} {
	c := make(chan interface{})
	if _, ok := b.listener[evtName]; !ok {
		b.listener[evtName] = make([]chan interface{}, 5)
	}

	log.Println("Registered listener to evt: ", evtName)
	b.listener[evtName] = append(b.listener[evtName], c)
	return c
}
