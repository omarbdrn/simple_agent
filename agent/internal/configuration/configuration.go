package configuration

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/mem"
)

type Config struct {
	AvailableMemory  int
	MaxParallelScans int
}

func GetSystemCapabilities() Config {
	numCPU := runtime.NumCPU()

	memory, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Error retrieving memory information:", err)
		return Config{}
	}

	// This will be for testing purposes
	maxParallelScans := numCPU * int(memory.Total/(1024*1024*1024))

	return Config{
		AvailableMemory:  int(memory.Total / (1024 * 1024 * 1024)),
		MaxParallelScans: maxParallelScans,
	}
}
