package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

func createSimpleGoRoutines() {
	log.Info().Msg("Hello from delegating function.")

	// A Wait Group makes sure that all goroutines launched as a part of it
	// finish before the function exits.
	var wg sync.WaitGroup

	// Add a 1 beause we are launching another single go routine.
	// A Wait Group keeps an internal counter that decrements when `wg.Done`
	// is called. If the counter hits 0, the blocking call `wg.Wait` continues,
	// finishing the function.
	wg.Add(1)

	// Creating a go routine that just prints out something.
	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		log.Info().Msg("Hi from inside a goroutine!")
	}()

	wg.Wait()

	log.Info().Msg("Delegating function is done.")
}

func simpleChannelExample() {
	// Creates a new channel with the Type string
	messages := make(chan string)

	// This Go routine sends a message onto the channel.
	messages <- "ping"

	// This receiver blocks until it finds a message, and then receives the "ping" from the go routine.
	msg := <-messages

	log.Info().Msg(msg)
}

func channelsToPassDataBetweenRoutines() {
	log.Info().Msg("Hello from delegating function.")

	// Creates a new channel with the Type string and a wait group.
	messages := make(chan string)
	var wg sync.WaitGroup

	// We are going to make a creator and consumer.
	// These will need the main method to wait for it to finish.
	wg.Add(1)
	go createMessageEvents(messages, &wg)

	wg.Add(1)
	go processEvents(messages, &wg)

	// Wait in the main method for the producers and consumers to do
	// their thing.
	wg.Wait()

	log.Info().Msg("Delegating function is done.")
}

// Creates a specific number of events onto the channel.
func createMessageEvents(msgChannel chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	// When we are done producing example events, we need to close the channel so the consumer unblocks.
	// Ref: https://golang.org/ref/spec#For_statements (under the "For statements with range clause" #4 bullet)
	defer close(msgChannel)

	for i := 1; i < 10; i++ {
		msgChannel <- "ping number: " + fmt.Sprint(i)
	}
}

// Consumes events from the channel.
func processEvents(msgChannel <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for message := range msgChannel {
		log.Info().Msg(message)
	}
}

func multipleChannelProducers() {
	log.Info().Msg("Hello from delegating function.")

	// Creates a new channel with the Type string and a wait group.
	messages := make(chan string)

	// We are going to make two creators and a single consumer.
	// These will need the main method to wait for it to finish, that is what the wg WaitGroup is for.
	var subroutine sync.WaitGroup
	var prouducer sync.WaitGroup

	subroutine.Add(1)
	prouducer.Add(1)
	go producePingMessages(messages, &subroutine, &prouducer)

	subroutine.Add(1)
	prouducer.Add(1)
	go producePongMessages(messages, &subroutine, &prouducer)

	subroutine.Add(1)
	go consumeMessages(messages, &subroutine)

	// Function to wait for the producuers to be done so we can close the channel.
	go func() {
		prouducer.Wait()
		close(messages)
	}()

	// Wait in the main method for the producers and consumers to do
	// their thing.
	subroutine.Wait()

	log.Info().Msg("Delegating function is done.")
}

// Prouduces Ping messages.
func producePingMessages(msgChannel chan<- string, subroutine *sync.WaitGroup, prouducer *sync.WaitGroup) {
	defer subroutine.Done()
	defer prouducer.Done()

	for i := 1; i < 10; i++ {
		msgChannel <- "ping number: " + fmt.Sprint(i)
		time.Sleep(time.Millisecond)
	}
}

// Prouduces Pong messages.
func producePongMessages(msgChannel chan<- string, subroutine *sync.WaitGroup, prouducer *sync.WaitGroup) {
	defer subroutine.Done()
	defer prouducer.Done()

	for i := 1; i < 10; i++ {
		msgChannel <- "pong number: " + fmt.Sprint(i)
		time.Sleep(time.Millisecond)
	}
}

// Consumes messages.
func consumeMessages(msgChannel <-chan string, subroutine *sync.WaitGroup) {
	defer subroutine.Done()

	for message := range msgChannel {
		log.Info().Msg(message)
	}
}
