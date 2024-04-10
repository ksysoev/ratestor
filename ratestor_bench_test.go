package ratestor

import (
	rand "math/rand/v2"
	"strconv"
	"testing"
	"time"
)

func BenchmarkSingleKey(b *testing.B) {
	limiter := NewRateStor()
	defer limiter.Close()

	for i := 0; i < b.N; i++ {
		_ = limiter.Allow("key", 1*time.Second, 1)
	}
}

func BenchmarkDifferentKeys(b *testing.B) {
	limiter := NewRateStor()
	defer limiter.Close()

	for i := 0; i < b.N; i++ {
		_ = limiter.Allow(strconv.Itoa(i), 1*time.Second, 1)
	}
}

func BenchmarkSameKeyParallelSingleKey(b *testing.B) {
	limiter := NewRateStor()
	defer limiter.Close()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = limiter.Allow("key", 1*time.Second, 1)
		}
	})
}

func BenchmarkSameKeyParallelDifferentKeys(b *testing.B) {
	limiter := NewRateStor()
	defer limiter.Close()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			randomKey := strconv.Itoa(rand.IntN(1000))
			_ = limiter.Allow(randomKey, 1*time.Second, 1)
		}
	})
}
