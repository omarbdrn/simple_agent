package network

import (
	"fmt"
	"time"

	"github.com/omarbdrn/simple_agent/internal/extools/naabu"
	"github.com/omarbdrn/simple_agent/internal/extools/naabu/banners"
	"github.com/omarbdrn/simple_agent/pkg/components/reporter"
	"github.com/omarbdrn/simple_agent/pkg/models"
	probing "github.com/prometheus-community/pro-bing"
)

func Ping(ipAddress string) bool {
	pinger, err := probing.NewPinger(ipAddress)
	if err != nil {
		panic(err)
	}
	pinger.Timeout = 2 * time.Second
	err = pinger.Run()
	if err != nil {
		return false
	}

	stats := pinger.Statistics()
	if stats.PacketLoss > 50 {
		return false
	}
	return true
}

func NetworkScan(ipAddress string) {
	naabuResult := naabu.Naabu(ipAddress)
	if naabuResult.IP != "" {
		hostService := models.HostServices{
			IP:    naabuResult.IP,
			Ports: []models.Service{},
		}

		for _, port := range naabuResult.Ports {
			hostPrefix := fmt.Sprintf("%s:%d", ipAddress, port)
			results := banners.BannerGrab(hostPrefix)

			service := models.Service{
				Port: port,
			}

			if len(results) > 0 {
				serviceBanner := results[0]
				service.Protocol = serviceBanner.Protocol
				service.Version = serviceBanner.Version
				service.Metadata = serviceBanner.Metadata()
			}

			hostService.Ports = append(hostService.Ports, service)
		}

		reporter.Report(hostService)
	}
}
