package meminfo

import (
	"strconv"
	"strings"
)

type MemInfo struct {
	Total     int
	Free      int
	Available int
}

// MemTotal: 16048468 kB -> 16048468
func parseMemInfoVal(meminfoline string) (int, error) {
	split := strings.Split(meminfoline, " ")
	return strconv.Atoi(split[len(split)-2])
}
