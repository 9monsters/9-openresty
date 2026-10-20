package stats

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"openresty-admin/internal/nginx"
)

type Collector struct {
	host      *nginx.Host
	interval  time.Duration
	sysRing   *Ring[SystemPoint]
	nginxRing *Ring[NginxPoint]
}

func NewCollector(host *nginx.Host, interval time.Duration) *Collector {
	return &Collector{
		host:      host,
		interval:  interval,
		sysRing:   SystemHistory,
		nginxRing: NginxHistory,
	}
}

func (c *Collector) Run(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	c.sample()
	slog.Info("stats collector started", "interval", c.interval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.sample()
		}
	}
}

func (c *Collector) sample() {
	total, free := memoryInfo()
	load := loadAverage()
	uptime := uptimeSeconds()

	c.sysRing.Push(SystemPoint{
		Time:          time.Now(),
		CPUCount:      runtime.NumCPU(),
		MemoryUsed:    total - free,
		MemoryTotal:   total,
		LoadAvg1:      load[0],
		LoadAvg5:      load[1],
		LoadAvg15:     load[2],
		UptimeSeconds: uptime,
	})

	ngxStats := c.host.Stats()
	if ngxStats != nil {
		c.nginxRing.Push(NginxPoint{
			Time:              time.Now(),
			ActiveConnections: ngxStats.ActiveConnections,
			Accepts:           ngxStats.Accepts,
			Handled:           ngxStats.Handled,
			Requests:          ngxStats.Requests,
			Reading:           ngxStats.Reading,
			Writing:           ngxStats.Writing,
			Waiting:           ngxStats.Waiting,
		})
	}
}
