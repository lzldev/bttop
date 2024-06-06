package procinfo

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/lzldev/bttop/internal/sysmanager/procinfo/proctype"
)

func TestProcInfo(t *testing.T) {
	reader, err := NewProcInfoReader()
	if err != nil {
		panic(err)
	}

	reader.Update()

	for k, v := range reader.Entries {
		dir := fmt.Sprintf("/proc/%v/task/", k)
		if v.Type != proctype.Proccess {
			continue
		}

		if _, err := os.Open(dir); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				fmt.Printf("no task on %v skipping\n", k)
				continue
			}
			panic(err)
		}

		entries, _ := os.ReadDir(dir)
		Equal(t, len(entries), v.threads, fmt.Sprintf("%+v", v))

	}
}

func Equal[T comparable](t *testing.T, a T, b T, info string) {
	t.Helper()

	if a != b {
		t.Errorf("Want :%v | Got :%v | %v", a, b, info)
	}
}
