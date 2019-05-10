package stats

import (
	"time"

	"github.com/shirou/gopsutil/process"

	"github.com/aakashRajur/star-wars/pkg/util"
)

func QueryStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	now := time.Now()
	stats[`now`] = now.UTC()
	compiled := make([]map[string]interface{}, 0)

	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	for _, proc := range processes {
		info := make(map[string]interface{})

		info[`pid`] = proc.Pid

		ppid, err := proc.Ppid()
		if err != nil {
			info[`ppid`] = `N/A`
		} else {
			info[`ppid`] = ppid
		}

		cmd, err := proc.Cmdline()
		if err != nil {
			info[`cmd`] = `N/A`
		} else {
			info[`cmd`] = cmd
		}

		elapsed, err := proc.CreateTime()
		if err != nil {
			info[`created`] = `N/A`
			info[`uptime`] = `N/A`
		} else {
			created := time.Unix(0, elapsed*int64(time.Millisecond))
			uptime := now.Sub(created)
			info[`created`] = created.UTC()
			info[`uptime`] = util.DurationToString(uptime)
		}

		memory, err := proc.MemoryPercent()
		if err != nil {
			info[`memory_percent`] = `N/A`
		} else {
			info[`memory_percent`] = memory
		}

		threads, err := proc.NumThreads()
		if err != nil {
			info[`num_threads`] = `N/A`
		} else {
			info[`num_threads`] = threads
		}

		times, err := proc.Times()
		if err != nil {
			info[`cpu`] = `N/A`
		} else {
			info[`cpu`] = times
		}

		connections, err := proc.Connections()
		if err != nil {
			info[`num_connections`] = `N/A`
		} else {
			info[`num_connections`] = len(connections)
		}

		compiled = append(compiled, info)
	}

	stats[`processes`] = compiled

	return stats, nil
}
