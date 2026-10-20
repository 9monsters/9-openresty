package stats

import (
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func memoryInfo() (total, free uint64) {
	if runtime.GOOS == "darwin" {
		total = sysctlUint("hw.memsize")
		free = darwinFreeMemory()
		return
	}
	return linuxMemory()
}

func loadAverage() []float64 {
	if runtime.GOOS == "darwin" {
		raw := runOutput("sysctl", "-n", "vm.loadavg")
		raw = strings.Trim(raw, "{ }")
		return parseFloats(strings.Fields(raw))
	}
	fields := strings.Fields(readFileString("/proc/loadavg"))
	if len(fields) > 3 {
		fields = fields[:3]
	}
	return parseFloats(fields)
}

func uptimeSeconds() uint64 {
	if runtime.GOOS == "darwin" {
		raw := runOutput("sysctl", "-n", "kern.boottime")
		match := regexp.MustCompile(`sec\s*=\s*(\d+)`).FindStringSubmatch(raw)
		if match == nil {
			return 0
		}
		boot, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0
		}
		now := time.Now().Unix()
		if now <= boot {
			return 0
		}
		return uint64(now - boot)
	}
	fields := strings.Fields(readFileString("/proc/uptime"))
	if len(fields) == 0 {
		return 0
	}
	seconds, _ := strconv.ParseFloat(fields[0], 64)
	return uint64(seconds)
}

func linuxMemory() (total, free uint64) {
	const kib = 1024
	values := map[string]uint64{}
	for _, line := range strings.Split(readFileString("/proc/meminfo"), "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		fields := strings.Fields(value)
		if len(fields) == 0 {
			continue
		}
		amount, _ := strconv.ParseUint(fields[0], 10, 64)
		values[key] = amount * kib
	}
	return values["MemTotal"], values["MemAvailable"]
}

func darwinFreeMemory() uint64 {
	raw := runOutput("vm_stat")
	pageSize := uint64(4096)
	if match := regexp.MustCompile(`page size of (\d+) bytes`).FindStringSubmatch(raw); match != nil {
		if parsed, _ := strconv.ParseUint(match[1], 10, 64); parsed > 0 {
			pageSize = parsed
		}
	}
	free := uint64(0)
	for _, line := range strings.Split(raw, "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		switch strings.TrimSpace(key) {
		case "Pages free", "Pages inactive", "Pages speculative":
			pages, _ := strconv.ParseUint(strings.Trim(strings.TrimSpace(value), "."), 10, 64)
			free += pages * pageSize
		}
	}
	return free
}

func sysctlUint(key string) uint64 {
	value, err := strconv.ParseUint(strings.TrimSpace(runOutput("sysctl", "-n", key)), 10, 64)
	if err != nil {
		return 0
	}
	return value
}

func runOutput(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func readFileString(path string) string {
	raw, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(raw)
}

func parseFloats(fields []string) []float64 {
	values := make([]float64, 0, len(fields))
	for _, field := range fields {
		parsed, _ := strconv.ParseFloat(field, 64)
		values = append(values, parsed)
	}
	return values
}