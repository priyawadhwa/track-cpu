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
	duration float64
	interval float64
)

func init() {
	flag.IntVar(&pid, "p", -1, "pid value")
	flag.Float64Var(&duration, "d", 60.0, "duration to collect data in seconds")
	flag.Float64Var(&interval, "i", .9, "interval at which to collect cpu data")
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
	for _, d := range data {
		fmt.Println(d)
	}

	fmt.Println()
	fmt.Printf("Average: %f\n", average(data))
	return nil
}

func getCPUData() ([]float64, error) {
	var cpuData []float64
	total := 0.0

	for {
		if total > duration {
			break
		}

		p, err := process.NewProcess(int32(pid))
		if err != nil {
			return nil, err
		}
		cpu, err := p.Percent(900 * time.Millisecond)
	
		if err != nil {
			return nil, err
		}
		cpuData = append(cpuData, cpu)
		total += .9
	}
	return cpuData, nil
}

func average(nums []float64) float64 {
	total := 0.0
	for _, n :=range nums {
		total += n
	}
	return total/float64(len(nums))
}
