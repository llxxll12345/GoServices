package main

import (
	"fmt"
	"sync"
)

var (
	matches     = sync.Map{} // Map from player id to a channel, the picker sends the result through the channel.
	mu          = sync.Mutex{}
	matchResult = make(map[string]string)
	queue       = make(chan string, 10)
	wg          = sync.WaitGroup{}
	done        = make(chan bool, 1)
)

func picker() {
	var p1, p2 string
	for {
		select {
		case val := <-queue:
			if len(p1) == 0 {
				p1 = val
			} else {
				p2 = val

				if ch1, ok := matches.Load(p1); ok {
					*(ch1.(*chan string)) <- p2
				}
				if ch2, ok := matches.Load(p2); ok {
					*(ch2.(*chan string)) <- p1
				}

				p1 = ""
				p2 = ""
			}
		case <-done:
			wg.Done()
			break
		}
	}
}

func requestMatch(id string) {
	ch := make(chan string, 1)
	matches.Store(id, &ch)
	queue <- id
	matched := <-ch
	fmt.Println("id: ", matched)

	mu.Lock()
	matchResult[id] = matched
	mu.Unlock()
	wg.Done()
}

func main() {
	wg.Add(1)
	go picker()
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go requestMatch(fmt.Sprintf("%d", i))
	}
	done <- true
	wg.Wait()

	// check match
	for i := 0; i < 20; i++ {
		player := fmt.Sprintf("%d", i)
		if matchResult[matchResult[player]] != player {
			fmt.Println("Incorrect match: ", matchResult[player], matchResult[matchResult[player]])
		}
	}
}
