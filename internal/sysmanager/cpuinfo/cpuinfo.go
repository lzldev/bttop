package cpuinfo

type CpuInfo struct {
	user       int // time spent in user mode
	nice       int // time spent in user low priority (nice)
	system     int //time spent in system mode
	idle       int // time spent idle (This value should be USER_HZ times the second entry in the /proc/uptime pseudo-file.)
	iowait     int // Time waiting for I/O to complete.
	irq        int // Time servicing interrupts
	softirq    int // Time servicing softirqs.
	steal      int // Time stolen
	guest      int // Time spent running a virtual CPU for guest operating systems under the control of the Linux kernel.
	guest_nice int // Time spent running a niced guest
	TotalTime  int
	WorkTime   int
}
