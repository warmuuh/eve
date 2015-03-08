package eve

import (
	"testing"
	 "time"
)

func TestEventDispatch(t *testing.T){
	b := Bus()

	result := make(chan string)
	
	go func () {
		<-b.From("test")
		result <- "finished"
	}()
	
	b.To("test") <- "Test"

	<- result

}

func TestEventDispatchFailing(t *testing.T){
	b := Bus()
	
	go func () {
		time.Sleep(1 * time.Second)
		<-b.From("test")
		t.Error("should not receive event")
	}()

	//will trigger earlier then listener will listen
	b.To("test") <- "Test"

}