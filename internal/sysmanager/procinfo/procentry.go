package procinfo

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/lzldev/bttop/internal/sysmanager/procinfo/proctype"
)

type ProcEntry struct {
	t            proctype.ProcType
	Pid          int
	ParentPid    int
	Name         string
	utime        int
	stime        int
	virtual_mem  int
	resident_mem int
	threads      int
	stat_file    *os.File
}

func NewProcEntry(pid int) *ProcEntry {
	dir := fmt.Sprint(procinfopath, pid, "/stat")
	file, err := os.Open(dir)
	if err != nil {
		panic(fmt.Errorf("to open process stat_file %v directory %v :  %w", pid, dir, err))
	}

	return &ProcEntry{
		Pid:       pid,
		stat_file: file,
	}
}

func (p *ProcEntry) Update() {
	if _, err := p.stat_file.Seek(1, io.SeekStart); err != nil {
		panic(fmt.Errorf("process %v stat seek: %w", p.Pid, err))
	}
	scanner := bufio.NewScanner(p.stat_file)
	scanner.Scan()

	stat := scanner.Text()

	values := strings.Split(stat, " ")

	if p.t == 0 {
		parent_pid, _ := strconv.Atoi(values[3])
		p.ParentPid = parent_pid

		if p.ParentPid == p.Pid {
			p.t = proctype.Proccess
		} else {
			p.t = proctype.Thread
		}
	}

	if p.Name == "" {
		l := len(values[1])
		p.Name = values[1][min(l, 1):l]
	}

	utime, _ := strconv.Atoi(values[15]) // /proc/1/stat | 14th value
	stime, _ := strconv.Atoi(values[16]) // /proc/1/stat | 15th value

	virtual_mem, _ := strconv.Atoi(values[24])
	resident_mem, _ := strconv.Atoi(values[25])

	threads, _ := strconv.Atoi(values[19])

	p.utime = utime
	p.stime = stime
	p.virtual_mem = virtual_mem
	p.resident_mem = resident_mem
	p.threads = threads
}
