package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

const (
	SizeGb = 1024 * 1024 * 1024
	SizeMb = 1024 * 1024
	SizeKb = 1024
)

var (
	accessCount int64 = 1 * 1000 * 1000 * 100
	indexes     [1 * 1000 * 1000 * 100]int32
	buffer      []int32
)

func main() {
	info, err := linuxproc.ReadCPUInfo("/proc/cpuinfo")
	if err != nil {
		fmt.Printf("Warning: failed to read cpuinfo: %s\n", err)
	} else {
		fmt.Printf("CPU Model: %s\n", info.Processors[0].ModelName)
		fmt.Printf("Total Cores: %d\n", info.NumCore())
	}

	rand.Seed(time.Now().UnixNano())
	size := SizeMb * 42 / 4
	var i int64
	buffer = make([]int32, size)
	for i := 0; i < size; i++ {
		buffer[i] = 0
	}
	for i = 0; i < accessCount; i++ {
		indexes[i] = int32(rand.Intn(size))
	}
	for i = 0; i < 1000; i++ {
		fmt.Printf("%.3f\n", randRead(size))
	}

}

func randRead(bufferSize int) float64 {
	var n int32
	var wg sync.WaitGroup

	var i int64
	for i = 0; i < 8; i++ {
		wg.Add(1)
		go func(index int64) {
			defer wg.Done()
			result := int32(randReadInternal(bufferSize, buffer, index*accessCount/8))
			atomic.AddInt32(&n, result)
		}(i)
	}

	wg.Wait()
	return float64(n) / 8.0
}

func randReadInternal(bufferSize int, buffer []int32, startPos int64) float64 {
	iterations := 20
	durations := make([]time.Duration, 0)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < iterations; i++ {
		var result int32
		start := time.Now()
		for j := startPos; j < accessCount; j++ {
			result |= buffer[indexes[j]]
		}
		duration := time.Since(start)
		durations = append(durations, duration)
		os.WriteFile("/dev/null", []byte(strconv.FormatInt(int64(result), 10)), 0755)
	}

	total := time.Duration(0)
	for i := 0; i < len(durations); i++ {
		total = total + durations[i]
	}

	avg := float64(total.Milliseconds()) / float64(len(durations))
	return avg
}
