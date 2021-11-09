package main

import (
	"github.com/dawidd6/go-appindicator"
	"github.com/gotk3/gotk3/gtk"
	"github.com/shirou/gopsutil/v3/cpu"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	//"github.com/shirou/gopsutil/v3/mem"
	"log"
	"time"
)

var box = []rune("▁▂▃▄▅▆▇██")

func buildCpuIndicator() *appindicator.Indicator {
	item, err := gtk.MenuItemNewWithLabel("Exit")
	if err != nil {
		log.Fatal(err)
	}
	_ = item.Connect("activate", func() {
		gtk.MainQuit()
	})
	//_ = indicator.Object().Connect(appindicator.SignalScrollEvent, func() {
	//	fmt.Println("scroll")
	//})
	menu, err := gtk.MenuNew()
	if err != nil {
		log.Fatal(err)
	}
	menu.Append(item)

	indicator := appindicator.New("cpu-indicator", "", appindicator.CategoryHardware)
	indicator.SetTitle("Bran's CPU Tray")
	exe, _ := os.Executable()
	path := filepath.Dir(exe)
	indicator.SetIcon(filepath.Join(path, "cpu_side.png"))
	indicator.SetAttentionIcon(filepath.Join(path, "cpu_bad.png"))
	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetOrderingIndex(1000000)
	indicator.SetMenu(menu)

	ms := uint(320)
	go func() {
		cpus, _ := cpu.Counts(true)
		guide := strings.Repeat(string(box[len(box)-1]), cpus)
		indicator.SetLabel(guide, guide)
		label := make([]rune, cpus)
		tick := time.NewTicker(time.Millisecond * time.Duration(ms))
		for range tick.C {
			ps, _ := cpu.Percent(0, true)
			sum := 0.0
			for i, p := range ps {
				label[i] = box[int(p*float64(len(box)-1)*0.01)]
				sum += p
			}
			if sum > 80.0*float64(len(ps)) {
				indicator.SetStatus(appindicator.StatusAttention)
			} else {
				indicator.SetStatus(appindicator.StatusActive)
			}
			indicator.SetLabel(string(label), guide)
		}
	}()

	return indicator
}

func main() {
	_ = syscall.Setpriority(syscall.PRIO_PROCESS, 0, -10)
	// GTK init and loop main at the end.

	gtk.Init(nil)
	defer gtk.Main()

	cpuIndicator := buildCpuIndicator()
	cpuIndicator.GetMenu().ShowAll()
}
