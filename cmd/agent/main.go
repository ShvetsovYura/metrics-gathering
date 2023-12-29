package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

type counter int64

// type gauge float64

type Metric struct {
	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapInuse     float64
	HeapObjects   float64
	HeapReleased  float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PauseTotalNs  float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64

	PollCount   int64
	RandomValue float64
}

var PoolCounter counter

func main() {
	poolInterval := 2
	// reportInterval := 10
	var m Metric
	for {
		CollectMetrics(&m)

		time.Sleep(time.Duration(poolInterval) * (time.Second))

		SendMetrics(m)
	}
}

func run(poolInterval int, reportINterval int) {

}

func SendMetrics(m Metric) {
	v := reflect.ValueOf(m)
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Field(i).String()
		link := makeLink(fieldName, v.Interface())
		sendRequest(link)
	}
}

func makeLink(mName string, v any) string {
	var link string
	switch v.(type) {
	case int64:
		link = fmt.Sprintf("http://localhost:8080/update/counter/%s/%d", mName, v)
	case float64:
		link = fmt.Sprintf("http://localhost:8080/update/gauge/%s/%f", mName, v)
	}
	return link
}

func sendRequest(link string) {
	r, err := http.Post(link, "text/plain", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Body.Close()
}

func CollectMetrics(m *Metric) {
	var rtm runtime.MemStats

	runtime.ReadMemStats(&rtm)

	m.HeapSys = float64(rtm.HeapSys)
	m.Alloc = float64(rtm.Alloc)
	m.BuckHashSys = float64(rtm.BuckHashSys)
	m.Frees = float64(rtm.Frees)
	m.GCCPUFraction = rtm.GCCPUFraction
	m.GCSys = float64(rtm.GCSys)
	m.HeapAlloc = float64(rtm.HeapAlloc)
	m.HeapIdle = float64(rtm.HeapIdle)
	m.HeapInuse = float64(rtm.HeapInuse)
	m.HeapObjects = float64(rtm.HeapObjects)
	m.HeapReleased = float64(rtm.HeapReleased)
	m.HeapSys = float64(rtm.HeapSys)
	m.LastGC = float64(rtm.LastGC)
	m.Lookups = float64(rtm.Lookups)
	m.MCacheInuse = float64(rtm.MCacheInuse)
	m.MCacheSys = float64(rtm.MCacheSys)
	m.MSpanInuse = float64(rtm.MSpanInuse)
	m.MSpanSys = float64(rtm.MSpanSys)
	m.Mallocs = float64(rtm.Mallocs)
	m.NextGC = float64(rtm.NextGC)
	m.NumForcedGC = float64(rtm.NumForcedGC)
	m.NumGC = float64(rtm.NumGC)
	m.OtherSys = float64(rtm.OtherSys)
	m.PauseTotalNs = float64(rtm.PauseTotalNs)
	m.StackInuse = float64(rtm.StackInuse)
	m.StackSys = float64(rtm.StackSys)
	m.Sys = float64(rtm.Sys)
	m.TotalAlloc = float64(rtm.TotalAlloc)
	m.RandomValue = rand.Float64()

}
