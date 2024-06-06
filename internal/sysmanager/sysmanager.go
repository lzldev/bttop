package sysmanager

import (
	"fmt"
	"time"

	"github.com/lzldev/bttop/internal/sysmanager/cpuinfo"
	"github.com/lzldev/bttop/internal/sysmanager/meminfo"
)

type Sender chan SysActorMessage
type Receiver chan SysMessage

type SysMessage struct {
	RamPct float64
	CpuPct float64
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

	var old_cpuinf *cpuinfo.CpuInfo

	go func() {
		defer close(rx)
		defer close(tx)

		for {
			select {
			case <-ticker.C:
				meminf, err := memreader.Read()
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

				tx <- SysMessage{
					RamPct: meminf.MemUsagePct(),
					CpuPct: cpu_pct,
				}

				logger <- fmt.Sprintf("%+v", *cpuinf)
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
