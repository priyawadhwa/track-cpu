package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/shirou/gopsutil/process"
)

var (
	pid      int
	duration int
	interval int
)

func init() {
	flag.IntVar(&pid, "p", -1, "pid value")
	flag.IntVar(&duration, "d", 60, "duration to collect data (seconds)")
	flag.IntVar(&interval, "i", 900, "interval at which to collect cpu data (ms)")
}

func main() {
	flag.Parse()
	if err := execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func validate() error {
	if pid == -1 {
		return errors.New("must provide pid with -p")
	}
	return nil
}

func execute() error {
	if err := validate(); err != nil {
		return err
	}
	data, err := getCPUData()
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("Average: %f\n", average(data))
	fmt.Printf("Max: %f\n", max(data))
	return nil
}

func getCPUData() ([]float64, error) {
	var cpuData []float64
	total := 0.0

	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, err
	}

	for {
		if total > float64(duration) {
			break
		}
		cpu, err := p.Percent(time.Duration(interval) * time.Millisecond)
		if err != nil {
			return nil, err
		}
		fmt.Println(cpu)
		cpuData = append(cpuData, cpu)
		total = total + float64(interval)/1000.0
	}
	return cpuData, nil
}

func average(nums []float64) float64 {
	total := 0.0
	for _, n := range nums {
		total += n
	}
	return total / float64(len(nums))
}

func max(nums []float64) float64 {
	max := 0.0
	for _, n := range nums {
		if n > max {
			max = n
		}
	}
	return max
}
