package naabu

import (
	"sync"

	"github.com/omarbdrn/simple_agent/internal/extools/naabu/port"
	naabuResult "github.com/omarbdrn/simple_agent/internal/extools/naabu/result"
	naabuRunner "github.com/omarbdrn/simple_agent/internal/extools/naabu/runner"
	"github.com/projectdiscovery/goflags"
)

type NaabuServiceInfo struct {
	Host  string `json:"host"`
	IP    string `json:"ip"`
	Ports []int  `json:"ports"`
}

type naabuResultCallback struct {
	sync.Mutex
	results      NaabuServiceInfo
	stdinResults []NaabuServiceInfo
}

func NaabuBatch(hostsfile string) (results []NaabuServiceInfo) {
	var callback naabuResultCallback

	naabuOptions := naabuRunner.Options{
		HostsFile:        hostsfile,
		Silent:           true,
		Verbose:          false,
		ServiceDiscovery: true,
		Debug:            false,
		Retries:          3,
		Timeout:          2000,
		ScanType:         "c",
		OnResult:         callback.onResultCallBackSTDIN,
	}

	naabu, err := naabuRunner.NewRunner(&naabuOptions)

	if err != nil {
		// fmt.Println(err)
	}

	defer naabu.Close()

	naabu.RunEnumeration()

	callback.Lock()
	callback.Unlock()

	return callback.stdinResults

}

func Naabu(host string) (results NaabuServiceInfo) {
	var callback naabuResultCallback

	naabuOptions := naabuRunner.Options{
		Host:             goflags.StringSlice{host},
		Silent:           true,
		Verbose:          false,
		Debug:            false,
		ServiceDiscovery: true,
		Retries:          3,
		Timeout:          2000,
		ScanType:         "c",
		OnResult:         callback.onResultCallBack,
	}

	naabu, err := naabuRunner.NewRunner(&naabuOptions)

	if err != nil {
		// fmt.Println(err)
	}

	defer naabu.Close()

	naabu.RunEnumeration()

	callback.Lock()
	callback.Unlock()

	return callback.results
}

func (cb *naabuResultCallback) onResultCallBackSTDIN(hr *naabuResult.HostResult) {
	serviceInfo := NaabuServiceInfo{
		Host:  hr.Host,
		IP:    hr.IP,
		Ports: extractPorts(hr.Ports),
	}

	cb.Lock()
	cb.stdinResults = append(cb.stdinResults, serviceInfo)
	cb.Unlock()
}

func (cb *naabuResultCallback) onResultCallBack(hr *naabuResult.HostResult) {
	serviceInfo := NaabuServiceInfo{
		Host:  hr.Host,
		IP:    hr.IP,
		Ports: extractPorts(hr.Ports),
	}

	cb.Lock()
	cb.results = serviceInfo
	cb.Unlock()
}

func extractPorts(ports []*port.Port) []int {
	portNumbers := make([]int, len(ports))
	for i, port := range ports {
		portNumbers[i] = port.Port
	}
	return portNumbers
}
