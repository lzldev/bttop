package cpuinfo

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const procStatPath = "/proc/stat"

type CpuInfoReader struct {
	file *os.File
}

func NewCpuInfoReader() (reader *CpuInfoReader, err error) {
	file, err := os.Open(procStatPath)
	if err != nil {
		return nil, err
	}

	return &CpuInfoReader{
		file: file,
	}, nil
}

func (reader *CpuInfoReader) Read() (info *CpuInfo, err error) {
	_, err = reader.file.Seek(1, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("cpu info reader %w", err)
	}

	scanner := bufio.NewScanner(reader.file)

	scanner.Scan()

	cpuLine := strings.Split(scanner.Text(), " ")[2:]

	times, err := parseTimeArray(cpuLine)
	if err != nil {
		return nil, err
	}

	total := 0
	for _, v := range times {
		total += v
	}

	total_work := 0
	for i, v := range times {
		if i == 3 || i == 4 { // Skip IO + Idle time
			continue
		}
		total_work += v
	}

	return &CpuInfo{
		user:       times[0],
		nice:       times[1],
		system:     times[2],
		idle:       times[3],
		iowait:     times[4],
		irq:        times[5],
		softirq:    times[6],
		steal:      times[7],
		guest:      times[8],
		guest_nice: times[9],
		TotalTime:  total,
		WorkTime:   total_work,
	}, nil
}

func parseTimeArray(arr []string) (value []int, err error) {
	times := make([]int, len(arr))

	for i, str := range arr {
		v, err := strconv.Atoi(str)

		if err != nil {
			return nil, fmt.Errorf("error trying to parse arr[%v] %v :%w\narr:%+v", i, str, err, arr)
		}

		times[i] = v
	}

	return times, nil
}
