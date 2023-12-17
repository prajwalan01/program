package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

type ProcessRequest struct {
	ToSort [][]int `json:"to_sort"`
}

type ProcessResponse struct {
	SequentialTime []int64 `json:"sequential_time,omitempty"`
	ConcurrentTime []int64 `json:"concurrent_time,omitempty"`
}

func processSequential(arrays [][]int) []int64 {
	sequentialTimes := make([]int64, len(arrays))

	for i, arr := range arrays {
		startTime := time.Now().UnixNano()
		sort.Ints(arr)
		sequentialTimes[i] = time.Now().UnixNano() - startTime
	}

	return sequentialTimes
}

func processConcurrent(arrays [][]int) []int64 {
	concurrentTimes := make([]int64, len(arrays))
	var wg sync.WaitGroup
	wg.Add(len(arrays))

	for i, arr := range arrays {
		go func(i int, arr []int) {
			defer wg.Done()
			startTime := time.Now().UnixNano()
			sort.Ints(arr)
			concurrentTimes[i] = time.Now().UnixNano() - startTime
		}(i, arr)
	}

	wg.Wait()

	return concurrentTimes
}

func processSingleHandler(w http.ResponseWriter, r *http.Request) {
	var requestData ProcessRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sequentialTimes := processSequential(requestData.ToSort)

	response := ProcessResponse{
		SequentialTime: sequentialTimes,
	}

	json.NewEncoder(w).Encode(response)
}

func processConcurrentHandler(w http.ResponseWriter, r *http.Request) {
	var requestData ProcessRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	concurrentTimes := processConcurrent(requestData.ToSort)

	response := ProcessResponse{
		ConcurrentTime: concurrentTimes,
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/process-single", processSingleHandler)
	mux.HandleFunc("/process-concurrent", processConcurrentHandler)

	server := &http.Server{
		Addr:    ":9000",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
