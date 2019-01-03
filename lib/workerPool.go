package lib

import (
	"sync"
)

func WorkerPool(worker func(string), setup func(chan string), workers int) {
	var wg sync.WaitGroup
	channel := make(chan string)
	for w := 1; w <= workers; w++ {
		go func() {
			wg.Add(1)
			for item := range channel {
				worker(item)
			}
			wg.Done()
		}()
	}

	setup(channel)
	close(channel)
	wg.Wait()
}
