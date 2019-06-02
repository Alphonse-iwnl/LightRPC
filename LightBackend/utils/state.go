package utils

import (
	"github.com/modood/table"
	"log"
	"time"
)

type StateEntry struct {
	ServiceName   string
	ServiceStatus int
	StartTime     int64
	EndTime       int64
	Duration      int64
}

type EntryManager struct {
	EntryList  []*StateEntry
	EntryQueue chan *StateEntry
	Size       int
	log        log.Logger

	// 熔断 流控功能
	// 记录服务的name-rate
	// 熔断功能监控该结构体 rate下降到阈值开始熔断服务
	MethodStatus map[string]float32
}

type StatisticsInfo struct {
	ServiceName string
	StateCode   int
	Average     float32
	Max         float32
	Min         float32
	Rate        float32
	Times       int
}

var Manager EntryManager

func init() {
	Manager = EntryManager{}
	Manager.EntryQueue = make(chan *StateEntry, 100)
	Manager.EntryList = []*StateEntry{}
	Manager.MethodStatus = make(map[string]float32)
	go StartListen()
}

func StartListen() {
	PrintTime := time.Duration(time.Second * 30)
	ticker := time.NewTicker(PrintTime)
	var tmpEntry *StateEntry
	for {
		select {
		case tmpEntry = <-Manager.EntryQueue:
			Manager.EntryList = append(Manager.EntryList, tmpEntry)
			Manager.Size++
		case <-ticker.C:
			statisticsServiceState()
		}
	}

}

func statisticsServiceState() {
	StatisticsList := make(map[string][]*StateEntry)
	var Statistics []StatisticsInfo
	ServiceMaxTime := make(map[string]int64)
	ServiceMinTime := make(map[string]int64)
	ServiceAverageTime := make(map[string]float32)

	// collect service
	for _, item := range Manager.EntryList {
		StatisticsList[item.ServiceName] = append(StatisticsList[item.ServiceName], item)
	}

	for name, entryList := range StatisticsList {
		stateMap := make(map[int]int)
		totalServiceTimes := len(entryList)
		total := int64(0)
		for _, entry := range entryList {
			if max, ok := ServiceMaxTime[entry.ServiceName]; ok {
				if entry.Duration > max {
					ServiceMaxTime[entry.ServiceName] = entry.Duration
				}
			} else {
				ServiceMaxTime[entry.ServiceName] = entry.Duration
			}
			if min, ok := ServiceMinTime[entry.ServiceName]; ok {
				if entry.Duration < min {
					ServiceMinTime[entry.ServiceName] = entry.Duration
				}
			} else {
				ServiceMinTime[entry.ServiceName] = entry.Duration
			}
			total = total + entry.Duration
			stateMap[entry.ServiceStatus]++
		}
		ServiceAverageTime[name] = float32(total/int64(len(entryList))) / 1000000
		for code, times := range stateMap {
			rate := float32(times) / float32(totalServiceTimes)
			var max, min, ava float32
			if code == 200 {
				max = float32(ServiceMaxTime[name]) / float32(1000000)
				min = float32(ServiceMinTime[name]) / float32(1000000)
				ava = float32(ServiceAverageTime[name])
			}
			Statistics = append(Statistics, StatisticsInfo{
				ServiceName: name,
				Times:       times,
				StateCode:   code,
				Rate:        rate,
				Max:         max,
				Min:         min,
				Average:     ava,
			})
			Manager.MethodStatus[name] = rate
		}
	}

	if len(Statistics) != 0 {
		table.Output(Statistics)
		p := table.Table(Statistics)
		LOG.Infof("\nServer statistics in 30s:\n" + p)
	}
	Manager.EntryList = Manager.EntryList[:0]
}

func NewStateEntry(serviceName string) *StateEntry {
	newEntry := &StateEntry{
		ServiceName: serviceName,
		StartTime:   time.Now().UnixNano(),
	}
	return newEntry
}

func EndStateEntry(entry *StateEntry, status int) {
	entry.EndTime = time.Now().UnixNano()
	entry.ServiceStatus = status
	entry.Duration = entry.EndTime - entry.StartTime
	Manager.EntryQueue <- entry
}
