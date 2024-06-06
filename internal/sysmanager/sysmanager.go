package sysmanager

import (
	"fmt"
	"time"

	"github.com/lzldev/bttop/internal/sysmanager/meminfo"
)

type Sender chan SysActorMessage
type Receiver chan SysMessage

type SysMessage struct {
	RamPct float64
}

const startingDuration time.Duration = 5

func StartSysManager(logger chan<- string) (rx Sender, tx Receiver) {
	rx = make(Sender)
	tx = make(Receiver)

	var interval time.Duration = startingDuration
	ticker := time.NewTicker(tickerIntervalMs(interval))

	reader, err := meminfo.NewMemInfoReader()
	if err != nil {
		panic(fmt.Errorf("couldn't start meminfo reader : %w", err))
	}

	go func() {
		defer close(rx)
		defer close(tx)

		for {
			select {
			case <-ticker.C:
				inf, err := reader.Read()
				if err != nil {
					panic(fmt.Errorf("couldn't read", err))
				}

				var ram_pct float64 = float64(inf.Available) / float64(inf.Total)

				tx <- SysMessage{
					RamPct: ram_pct,
				}

				// logger <- fmt.Sprintf("%+v", *inf)
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
