package sysmanager

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/lzldev/bttop/internal/sysmanager/cpuinfo"
	"github.com/lzldev/bttop/internal/sysmanager/meminfo"
	"github.com/lzldev/bttop/internal/sysmanager/procinfo"
	"github.com/lzldev/bttop/internal/sysmanager/procinfo/proctype"
)

type Sender chan SysActorMessage
type Receiver chan SysMessage

type SysMessage struct {
	RamPct   float64
	CpuPct   float64
	ProcRows []table.Row
}

const startingDuration time.Duration = 5

func StartSysManager(logger chan<- string) (rx Sender, tx Receiver) {
	rx = make(Sender)
	tx = make(Receiver)

	var interval time.Duration = startingDuration
	ticker := time.NewTicker(tickerIntervalMs(interval))

	memreader, err := meminfo.NewMemInfoReader()
	if err != nil {
		panic(fmt.Errorf("sysmanager meminfo reader : %w", err))
	}

	cpureader, err := cpuinfo.NewCpuInfoReader()
	if err != nil {
		panic(fmt.Errorf("sysmanager cpuinfo reader : %w", err))
	}

	procreader, err := procinfo.NewProcInfoReader()
	if err != nil {
		panic(fmt.Errorf("sysmanager procinfo reader : %w", err))
	}

	procinf := &procreader.Entries

	var meminf *meminfo.MemInfo = &meminfo.MemInfo{}
	var old_cpuinf *cpuinfo.CpuInfo

	go func() {
		defer close(rx)
		defer close(tx)

		for {
			select {
			case <-ticker.C:
				err := memreader.ReadTo(meminf)
				if err != nil {
					panic(fmt.Errorf("sysmanager meminfo read %w", err))
				}

				cpuinf, err := cpureader.Read()
				if err != nil {
					panic(fmt.Errorf("sysmanager cpuinfo read %w", err))
				}

				var cpu_pct float64 = 0
				if old_cpuinf != nil {
					cpu_pct = cpuinf.CpuUsagePct(old_cpuinf)
				}
				old_cpuinf = cpuinf

				procreader.Update()

				rows := make([]table.Row, 0)
				for _, v := range *procinf {
					if v.Type != proctype.Proccess {
						continue
					}

					rows = append(rows,
						table.Row{strconv.Itoa(v.Pid), v.Name, strconv.Itoa(v.Utime), v.Type.String()},
					)
				}

				tx <- SysMessage{
					RamPct:   meminf.MemUsagePct(),
					CpuPct:   cpu_pct,
					ProcRows: rows,
				}

				logger <- fmt.Sprintf("%#v", len(rows))
				// logger <- fmt.Sprintf("[sys] tick \nfff: %v ffffs:%v free:%v free+swap %v \n %+v", info.Freeram/1024/1024, (info.Freeram+info.Freeswap)/1024/1024, info.Freeram, info.Freeram+info.Freeswap, info)
			case msg := <-rx:
				switch msg {
				case Stop:
					break
				case Nothing:
				case SpeedDown:
					if interval <= 1 {
						continue
					}
					interval--
					ticker.Reset(tickerIntervalMs(interval))
					logger <- fmt.Sprintf("[SYS] New Speed %vms", interval.Nanoseconds()*100)
				case SpeedUp:
					interval++
					ticker.Reset(tickerIntervalMs(interval))
					logger <- fmt.Sprintf("[SYS] New Speed %vms", interval.Nanoseconds()*100)
				}
			}
		}

	}()

	return
}
