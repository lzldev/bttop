package procinfo

import (
	"fmt"
	"os"
	"strconv"
)

const procinfopath = "/proc/"

type ProcInfoReader struct {
	Entries map[int]*ProcEntry
}

func NewProcInfoReader() (*ProcInfoReader, error) {
	return &ProcInfoReader{
		Entries: make(map[int]*ProcEntry),
	}, nil
}

func (reader *ProcInfoReader) Update() error {
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
		proc, ok := reader.Entries[pid]

		if !ok {
			reader.Entries[pid] = NewProcEntry(pid)
			reader.Entries[pid].Update()
		} else {
			proc.Update()
		}
	}

	return nil
}
