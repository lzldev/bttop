package meminfo

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const memInfoPath = "/proc/meminfo"

type MemInfoReader struct {
	file *os.File
}

func NewMemInfoReader() (*MemInfoReader, error) {
	file, err := os.Open(memInfoPath)

	if err != nil {
		return nil, err
	}

	return &MemInfoReader{
		file: file,
	}, nil
}

func (reader *MemInfoReader) Read() (info *MemInfo, err error) {
	if _, err = reader.file.Seek(1, io.SeekStart); err != nil {
		return nil, fmt.Errorf("meminfo read: %w", err)
	}

	scanner := bufio.NewScanner(reader.file)

	scanner.Scan()
	total, err := parseMemInfoVal(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("meminfo read: %w", err)
	}

	scanner.Scan()
	free, err := parseMemInfoVal(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("meminfo read: %w", err)
	}

	scanner.Scan()
	available, err := parseMemInfoVal(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("meminfo read: %w", err)
	}

	return &MemInfo{
		Total:     total,
		Free:      free,
		Available: available,
	}, nil
}

func (reader *MemInfoReader) ReadTo(to *MemInfo) (err error) {
	_, err = reader.file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("meminfo readTo: %w", err)
	}
	scanner := bufio.NewScanner(reader.file)
	scanner.Scan()

	total, err := parseMemInfoVal(scanner.Text())
	if err != nil {
		return fmt.Errorf("meminfo readTo: %w", err)
	}

	scanner.Scan()
	free, err := parseMemInfoVal(scanner.Text())
	if err != nil {
		return fmt.Errorf("meminfo readTo: %w", err)
	}

	to.Free = free
	to.Total = total

	return
}
