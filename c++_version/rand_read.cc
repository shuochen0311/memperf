#include <iostream>
#include <cstdlib>
#include <stdio.h>
#include <thread>
#include <vector>
#include <sys/mman.h>
#include <stdlib.h>
#include <cstring>
#include <chrono>
#include <atomic>
#include <time.h>
#include <fstream>
#include <sched.h>
#include <pthread.h>
#include <sys/time.h>

using namespace std;

long iterations = 20;
long accessCount = 100ULL << 20;
int threadCount = 8;
vector<int32_t> indexes(accessCount, 0);
atomic<int> readyCount;
atomic<int> totalTime;

void* mem;

ofstream myfile;

int cpuPinOffset = 0;

double randReadImpl(long index)
{
    readyCount++;
    while(readyCount != threadCount);
    int time = 0;
    char result = '0';
    for (long i = 0; i < iterations; i++) {
        std::chrono::steady_clock::time_point begin = std::chrono::steady_clock::now();
        for (long n = index; n < accessCount/8; n+=8) {
            result |= ((char*)(mem))[n];
        }
        std::chrono::steady_clock::time_point end = std::chrono::steady_clock::now();
        time += std::chrono::duration_cast<std::chrono::microseconds>(end - begin).count();   
        myfile << result;  
    }

    totalTime += time;
    return time / iterations;
}

double randRead(long size)
{
    readyCount = 0;
    totalTime = 0;

    vector<std::thread> threads;

    for(int i = 0; i < threadCount; i++) {
        threads.emplace_back(randReadImpl, i);
        cpu_set_t cpuset;
        CPU_ZERO(&cpuset);
        CPU_SET(cpuPinOffset + i, &cpuset);
        int rc = pthread_setaffinity_np(threads[i].native_handle(),
                                      sizeof(cpu_set_t), &cpuset);
    }

    for (auto& th : threads) {
        th.join();
    }

    return totalTime / threadCount / iterations;
}


int main(int argc, char const *argv[])
{
    cpuPinOffset = atoi(argv[1]);
    myfile.open("/dev/null");
    long  SizeGb = 1ULL << 30;
    long SizeMb = 1ULL << 20;
    long SizeKb = 1ULL << 10;

    long size = SizeGb * 1;
    printf("size is %ld\n", size);

    // prepare memory
    mem = (void*)malloc(size);
    for (int i = 0; i < size; i++) {
       ((char*)(mem))[i] = 'a';
    }
    memset(mem, 0, size);
    srand(time(NULL));

    for(int i = 0; i < accessCount; i++) {
        indexes[i] = rand() % size;
    }
    
    for (int i = 0; i < 30000; i++) {
        struct timeval tv;
        if(gettimeofday(&tv, NULL) != 0)
                return 0;

        auto ticks = (tv.tv_sec * 1000) + (tv.tv_usec / 1000);
        printf("%ld, %.3f\n", ticks, randRead(size));
    }
    return 0;
}
//g++ rand_read.cc -o test -pthread -O3
// to run ./test