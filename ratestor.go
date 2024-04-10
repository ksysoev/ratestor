package ratestor

import (
	"container/heap"
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	GC_BATCH_SIZE = 100
)

// ErrRateLimitExceeded is an error that is returned when the rate limit is exceeded.
// This error indicates that the maximum number of requests allowed within a certain time period has been reached.
var ErrRateLimitExceeded = fmt.Errorf("rate limit exceeded")

type RateValue struct {
	Value     uint64
	ExpiresAt time.Time
	Limit     uint64
}

type RateStor struct {
	lock       *sync.Mutex
	rates      map[string]RateValue
	index      expIndex
	gcInterval time.Duration
	stop       context.CancelFunc
	wg         *sync.WaitGroup
}

type indexValue struct {
	Key       string
	ExpiresAt time.Time
}

type expIndex []indexValue // Heap type

func (h expIndex) Len() int           { return len(h) }
func (h expIndex) Less(i, j int) bool { return h[i].ExpiresAt.Before(h[j].ExpiresAt) }
func (h expIndex) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *expIndex) Push(x any)        { *h = append(*h, x.(indexValue)) }
func (h *expIndex) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

type Optition func(*RateStor)

// NewRateStor creates a new instance of RateStor with the provided options.
// It initializes the necessary fields and starts a goroutine for periodic cleaning.
// The cleaning interval is set to 1 second by default.
// The provided options can be used to customize the behavior of the RateStor instance.
func NewRateStor(opts ...Optition) *RateStor {

	ctx, cancel := context.WithCancel(context.Background())

	stor := &RateStor{
		lock:       &sync.Mutex{},
		rates:      make(map[string]RateValue),
		gcInterval: 1 * time.Second,
		stop:       cancel,
		wg:         &sync.WaitGroup{},
	}

	for _, opt := range opts {
		opt(stor)
	}

	stor.wg.Add(1)
	go stor.cleaner(ctx)

	return stor
}

// Allow allows a request with the given key if the rate limit is not exceeded.
// It takes the key, period, and limit as parameters and returns an error if the rate limit is exceeded.
// The key is used to identify the request, the period is the duration for which the rate limit is enforced,
// and the limit is the maximum number of requests allowed within the given period.
// If the rate limit is not exceeded, the function increments the rate value for the given key.
// If the rate limit is exceeded, it returns an ErrRateLimitExceeded error.
func (rs *RateStor) Allow(key string, period time.Duration, limit uint64) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	if rate, ok := rs.rates[key]; ok {
		if !rate.ExpiresAt.Before(time.Now()) {
			if rate.Value < rate.Limit {
				rate.Value++
				rs.rates[key] = rate

				return nil
			}

			return ErrRateLimitExceeded
		}
	}

	ExpiresAt := time.Now().Add(period)
	rs.rates[key] = RateValue{
		Value:     1,
		ExpiresAt: ExpiresAt,
		Limit:     limit,
	}

	heap.Push(&rs.index, indexValue{
		Key:       key,
		ExpiresAt: ExpiresAt,
	})

	return nil
}

// cleaner is a goroutine that periodically cleans up expired rate limit entries from the RateStor.
// It runs in the background and removes expired entries from the index and rates map.
// The cleaning interval is determined by the gcInterval field of the RateStor struct.
// This function should be called as a goroutine using the go keyword.
// It takes a context.Context as a parameter to allow for cancellation.
// The function will exit when the context is canceled.
func (rs *RateStor) cleaner(ctx context.Context) {
	var ticker = time.NewTicker(rs.gcInterval)
	defer ticker.Stop()
	defer rs.wg.Done()

	for {
		select {
		case <-ticker.C:
			isRunning := true

			for isRunning {
				rs.lock.Lock()
				for i := 0; i < GC_BATCH_SIZE; i++ {
					if rs.index.Len() == 0 {
						isRunning = false
						break
					}

					item := heap.Pop(&rs.index).(indexValue)
					isReady := item.ExpiresAt.Before(time.Now())
					if !isReady {
						heap.Push(&rs.index, item)
						isRunning = false
						break
					}

					val, ok := rs.rates[item.Key]
					if ok && val.ExpiresAt == item.ExpiresAt {
						delete(rs.rates, item.Key)
					}
				}
				rs.lock.Unlock()
			}
		case <-ctx.Done():
			return
		}
	}
}

// Close stops the RateStor instance and waits for all goroutines to complete.
func (rs *RateStor) Close() {
	rs.stop()
	rs.wg.Wait()
}

// WithGCInterval sets the garbage collection interval for the RateStor instance.
// The garbage collection interval determines how often the RateStor instance
// will perform garbage collection to remove expired rate limit entries.
// The default garbage collection interval is 1 second.
func WithGCInterval(interval time.Duration) Optition {
	return func(rs *RateStor) {
		rs.gcInterval = interval
	}
}
