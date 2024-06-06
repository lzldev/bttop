package procinfo

import (
	"fmt"
	"os"
	"strconv"
)

const procinfopath = "/proc/"

type ProcInfoReader struct {
	proc map[int]*ProcEntry
}

func NewProcInfoReader() (*ProcInfoReader, error) {
	return &ProcInfoReader{
		proc: make(map[int]*ProcEntry),
	}, nil
}

func (reader *ProcInfoReader) Read() error {
	dir, err := os.ReadDir(procinfopath)
	if err != nil {
		return fmt.Errorf("procinforeader read : %w", err)
	}

	count := 0
	for _, entry := range dir {
		if entry.Name()[0] > '9' {
			break
		}
		count++

		dirName := entry.Name()
		pid, _ := strconv.Atoi(dirName[:max(len(dirName), 0)])
		proc, ok := reader.proc[pid]

		if !ok {
			reader.proc[pid] = NewProcEntry(pid)
			reader.proc[pid].Update()
		} else {
			proc.Update()
		}

		fmt.Println(entry)
	}
	fmt.Println("End of PIDs | Total:", count)

	return nil
}
