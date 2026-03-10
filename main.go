package main

import (
	"fmt"
	"sync"
)

type DataItem struct {
	Value interface{}
}

type FlashDB struct {
	mu    sync.RWMutex
	store map[string]DataItem

	eventChan chan string

	statsMu     sync.Mutex
	totalReads  int
	totalWrites int
}

type Engine interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Delete(key string)
	GetStats() (int, int)
}

func (f *FlashDB) Set(key string, value interface{}) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.store[key] = DataItem{value}
	f.eventChan <- fmt.Sprintf("EVENT SET KEY: %s", key)

	f.statsMu.Lock()
	defer f.statsMu.Unlock()
	f.totalWrites += 1
}

func (f *FlashDB) Get(key string) (interface{}, bool) {
	f.statsMu.Lock()
	defer f.statsMu.Unlock()
	f.totalReads += 1

	f.mu.RLock()
	defer f.mu.RUnlock()
	if i, ok := f.store[key]; ok {
		return i.Value, true
	}
	return nil, false
}

func (f *FlashDB) Delete(key string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.store, key)
	f.eventChan <- fmt.Sprintf("EVENT DELETE KEY: %s", key)
}

func (f *FlashDB) GetStats() (int, int) {
	return f.totalReads, f.totalWrites
}

func main() {
	pipeCh := make(chan string, 100)
	stateCh := make(chan bool)
	var wg sync.WaitGroup
	db := &FlashDB{
		store:     make(map[string]DataItem),
		eventChan: pipeCh,
	}

	go func() {
		for event := range pipeCh {
			fmt.Println("Received event:", event)
		}

		stateCh <- true
	}()

	for range 100 {
		wg.Go(func() {
			db.Set("robert", 12)
		})
	}
	wg.Wait()
	fmt.Println(db.totalWrites)

	for range 1000 {
		wg.Go(func() {
			db.Get("robert")
		})
	}

	wg.Wait()
	fmt.Println(db.totalReads)

	db.Delete("robert")
	fmt.Println(db.GetStats())
	close(pipeCh)
	<-stateCh
	close(stateCh)
	fmt.Println("Exiting main")
}
