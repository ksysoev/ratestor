package ratestor

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkSingleKey(b *testing.B) {
	limiter := NewRateStor()
	defer limiter.Close()
	for i := 0; i < b.N; i++ {
		limiter.Allow("key", 1*time.Second, 1)
	}
}

func BenchmarkDifferentKeys(b *testing.B) {
	limiter := NewRateStor()
	defer limiter.Close()

	for i := 0; i < b.N; i++ {
		limiter.Allow(strconv.Itoa(i), 1*time.Second, 1)
	}
}

func BenchmarkSameKeyParallelSingleKey(b *testing.B) {
	limiter := NewRateStor()
	defer limiter.Close()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow("key", 1*time.Second, 1)
		}
	})
}

func BenchmarkSameKeyParallelDifferentKeys(b *testing.B) {
	limiter := NewRateStor()
	defer limiter.Close()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			random_key := strconv.Itoa(rand.Intn(1000))
			limiter.Allow(random_key, 1*time.Second, 1)
		}
	})
}
