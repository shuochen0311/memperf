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

	fmt.Printf("populating index")
	rand.Seed(time.Now().UnixNano())
	// fmt.Printf("Hashsize(million), Latency(ms)\n")
	// fmt.Printf("%d, %.3f\n", 1, hashTable(1*SizeMb))
	// fmt.Printf("%d, %.3f\n", 2, hashTable(2*SizeMb))
	// fmt.Printf("%d, %.3f\n", 5, hashTable(5*SizeMb))
	// fmt.Printf("%d, %.3f\n", 10, hashTable(10*SizeMb))
	// fmt.Printf("%d, %.3f\n", 20, hashTable(20*SizeMb))
	// fmt.Printf("%d, %.3f\n", 30, hashTable(30*SizeMb))

	fmt.Printf("Buffer(KB), Latency(ms)\n")
	// fmt.Printf("%d, %.3f\n", 1, randWrite(SizeKb*1))
	// fmt.Printf("%d, %.3f\n", 2, randWrite(SizeKb*2))
	// fmt.Printf("%d, %.3f\n", 4, randWrite(SizeKb*4))
	// fmt.Printf("%d, %.3f\n", 8, randWrite(SizeKb*8))
	// fmt.Printf("%d, %.3f\n", 16, randWrite(SizeKb*16))
	// fmt.Printf("%d, %.3f\n", 32, randWrite(SizeKb*32))
	// fmt.Printf("%d, %.3f\n", 64, randWrite(SizeKb*64))
	// fmt.Printf("%d, %.3f\n", 128, randWrite(SizeKb*128))
	// fmt.Printf("%d, %.3f\n", 256, randWrite(SizeKb*256))
	// fmt.Printf("%d, %.3f\n", 2000, randWrite(SizeMb*2))
	// fmt.Printf("%d, %.3f\n", 4000, randWrite(SizeMb*4))
	// fmt.Printf("%d, %.3f\n", 8000, randWrite(SizeMb*8))
	// fmt.Printf("%d, %.3f\n", 16000, randWrite(SizeMb*16))
	// fmt.Printf("%d, %.3f\n", 32000, randWrite(SizeMb*32))
	// fmt.Printf("%d, %.3f\n", 64000, randWrite(SizeMb*64))
	// fmt.Printf("%d, %.3f\n", 128000, randWrite(SizeMb*128))
	// fmt.Printf("%d, %.3f\n", 256000, randWrite(SizeMb*256))
	// fmt.Printf("%d, %.3f\n", 512000, randWrite(SizeMb*512))
	// fmt.Printf("%d, %.3f\n", 768000, randWrite(SizeMb*768))
	// fmt.Printf("%d, %.3f\n", 1000000, randWrite(SizeGb*1))

	size := 0
	// size = SizeKb * 1 / 4
	// buffer = make([]int32, size)
	// for i := 0; i < size; i++ {
	// 	buffer[i] = 0
	// }

	// fmt.Println(len(buffer))

	// var i int64
	// for i = 0; i < accessCount; i++ {
	// 	indexes[i] = int32(rand.Intn(size))
	// }
	// for i = 0; i < 1000; i++ {
	// 	fmt.Printf("%d, %.3f\n", 1, randRead(size))
	// }

	// size = SizeMb * 1 / 4
	// buffer = make([]int32, size)
	// for i := 0; i < size; i++ {
	// 	buffer[i] = 0
	// }

	// fmt.Println(len(buffer))

	// var i int64
	// for i = 0; i < accessCount; i++ {
	// 	indexes[i] = int32(rand.Intn(size))
	// }
	// for i = 0; i < 1000; i++ {
	// 	fmt.Printf("%d, %.3f\n", 1000, randRead(size))
	// }

	// size = SizeMb * 6 / 4
	// buffer = make([]int32, size)
	// for i := 0; i < size; i++ {
	// 	buffer[i] = 0
	// }

	// fmt.Println(len(buffer))

	// var i int64
	// for i = 0; i < accessCount; i++ {
	// 	indexes[i] = int32(rand.Intn(size))
	// }
	// for i = 0; i < 1000; i++ {
	// 	fmt.Printf("%d, %.3f\n", 6000, randRead(size))
	// }

	size = SizeMb * 512 / 4
	buffer = make([]int32, size)
	for i := 0; i < size; i++ {
		buffer[i] = 0
	}

	fmt.Println(len(buffer))

	var i int64
	for i = 0; i < accessCount; i++ {
		indexes[i] = int32(rand.Intn(size))
	}
	for i = 0; i < 1000; i++ {
		fmt.Printf("%d, %.3f\n", 512000, randRead(size))
	}

	// fmt.Println(len(buffer))

	// var i int64
	// for i = 0; i < accessCount; i++ {
	// 	indexes[i] = int32(rand.Intn(size))
	// }
	// for i = 0; i < 1000; i++ {
	// 	fmt.Printf("%d, %.3f\n", 42000, randRead(size))
	// }

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
