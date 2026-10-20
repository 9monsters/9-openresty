package handler

import (
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"openresty-admin/internal/nginx"
)

// --- sysinfo helpers ---

func memInfo() (total, free uint64) {
	if runtime.GOOS == "darwin" {
		total = sysctlUint("hw.memsize")
		free = darwinFree()
		return
	}
	return linuxMem()
}

func cpuModel() string {
	if runtime.GOOS == "darwin" { return runCmd("sysctl", "-n", "machdep.cpu.brand_string") }
	for _, line := range strings.Split(readFile("/proc/cpuinfo"), "\n") {
		if strings.HasPrefix(line, "model name") { if _, v, ok := strings.Cut(line, ":"); ok { return strings.TrimSpace(v) } }
	}
	return "Unknown"
}

func loadAvg() []float64 {
	if runtime.GOOS == "darwin" { return parseFloats(strings.Fields(strings.Trim(runCmd("sysctl", "-n", "vm.loadavg"), "{ }"))) }
	f := strings.Fields(readFile("/proc/loadavg"))
	if len(f) > 3 { f = f[:3] }
	return parseFloats(f)
}

func uptimeSecs() uint64 {
	if runtime.GOOS == "darwin" {
		raw := runCmd("sysctl", "-n", "kern.boottime")
		m := regexp.MustCompile(`sec\s*=\s*(\d+)`).FindStringSubmatch(raw)
		if m == nil { return 0 }
		boot, _ := strconv.ParseInt(m[1], 10, 64)
		now := time.Now().Unix()
		if now <= boot { return 0 }
		return uint64(now - boot)
	}
	f := strings.Fields(readFile("/proc/uptime"))
	if len(f) == 0 { return 0 }
	s, _ := strconv.ParseFloat(f[0], 64)
	return uint64(s)
}

func linuxMem() (total, free uint64) {
	for _, line := range strings.Split(readFile("/proc/meminfo"), "\n") {
		k, v, ok := strings.Cut(line, ":")
		if !ok { continue }
		f := strings.Fields(v)
		if len(f) == 0 { continue }
		a, _ := strconv.ParseUint(f[0], 10, 64)
		switch k {
		case "MemTotal": total = a * 1024
		case "MemAvailable": free = a * 1024
		}
	}
	return
}

func darwinFree() uint64 {
	raw := runCmd("vm_stat")
	ps := uint64(4096)
	if m := regexp.MustCompile(`page size of (\d+) bytes`).FindStringSubmatch(raw); m != nil { ps, _ = strconv.ParseUint(m[1], 10, 64) }
	var free uint64
	for _, line := range strings.Split(raw, "\n") {
		k, v, ok := strings.Cut(line, ":")
		if !ok { continue }
		switch strings.TrimSpace(k) {
		case "Pages free", "Pages inactive", "Pages speculative":
			p, _ := strconv.ParseUint(strings.Trim(strings.TrimSpace(v), "."), 10, 64)
			free += p * ps
		}
	}
	return free
}

func sysctlUint(key string) uint64 { v, _ := strconv.ParseUint(runCmd("sysctl", "-n", key), 10, 64); return v }
func runCmd(name string, args ...string) string { out, _ := exec.Command(name, args...).Output(); return string(out) }
func readFile(path string) string { raw, _ := os.ReadFile(path); return string(raw) }
func parseFloats(fields []string) []float64 { var out []float64; for _, f := range fields { v, _ := strconv.ParseFloat(f, 64); out = append(out, v) }; return out }
func pct(part, total uint64) int { if total == 0 { return 0 }; return int(part * 100 / total) }
func hostname() string { s, _ := os.Hostname(); return s }


func systemInfo(host *nginx.Host) gin.HandlerFunc {
	return func(c *gin.Context) {
		total, free := memInfo()
		c.JSON(http.StatusOK, gin.H{
			"cpu": gin.H{"count": runtime.NumCPU(), "model": cpuModel()},
			"memory": gin.H{"total": total, "free": free, "used": total - free, "usagePercent": pct(total-free, total)},
			"uptime": uptimeSecs(), "loadAvg": loadAvg(),
			"platform": runtime.GOOS, "arch": runtime.GOARCH, "hostname": hostname(), "go": runtime.Version(), "nodejs": "",
		})
	}
}

func selfCheck(host *nginx.Host) gin.HandlerFunc {
	return func(c *gin.Context) {
		checks := []gin.H{}
		if _, err := exec.LookPath(host.Paths().NginxBin); err == nil {
			checks = append(checks, gin.H{"name": "Nginx Binary", "status": "ok", "detail": host.Paths().NginxBin})
		} else {
			checks = append(checks, gin.H{"name": "Nginx Binary", "status": "warning", "detail": "Not found in PATH"})
		}
		if ok, o := host.Test(); ok { checks = append(checks, gin.H{"name": "Config Syntax", "status": "ok", "detail": o}) } else { checks = append(checks, gin.H{"name": "Config Syntax", "status": "error", "detail": o}) }
		c.JSON(http.StatusOK, gin.H{"checks": checks})
	}
}