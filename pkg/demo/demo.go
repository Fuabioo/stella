package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type state struct {
	sync.Mutex
	total   int
	current int
}

func (s *state) GetProgress() (total, current int) {
	s.Lock()
	defer s.Unlock()
	return s.total, s.current
}

func (s *state) IncProgress() {
	s.Lock()
	defer s.Unlock()
	s.current++
	if s.current > s.total {
		s.current = s.total
	}
}

func main() {
	step := time.Second / 8
	port := 10301
	state := &state{
		total:   100,
		current: 0,
	}

	http.HandleFunc("/progress", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Unsupported method", http.StatusBadRequest)
			return
		}

		total, current := state.GetProgress()
		text := "Lorem ipsum"

		w.Header().Set("Content-Type", "application/json")
		w.Write(fmt.Appendf(nil, `{"total": %d, "current": %d, "text": %q, "loading": true}`,
			total,
			current,
			text,
		))
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			time.Sleep(step)
			if state.current >= state.total {
				wg.Done()
				return
			}
			state.IncProgress()
		}
	}()

	go func() {
		addr := fmt.Sprintf(":%d", port)

		log.Printf("Serving at %s", addr)
		log.Printf("\t curl 'http://0.0.0.0%s/progress'", addr)

		err := http.ListenAndServe(addr, nil)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	wg.Wait()

	log.Println("Got to 100%")
	time.Sleep(time.Second)
}
