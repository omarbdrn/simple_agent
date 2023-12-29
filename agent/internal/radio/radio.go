package radio

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/omarbdrn/simple_agent/pkg/components/network"
	"github.com/omarbdrn/simple_agent/pkg/components/shares"

	"github.com/omarbdrn/fingerprintx/pkg/scan"
	"github.com/projectdiscovery/ipranger"
)

type RadioConfig struct {
	MaxConcurrentScans int
	ScanInterval       time.Duration
}

func StartRadio() {
	scan.Init()

	config := RadioConfig{
		MaxConcurrentScans: 50,
		ScanInterval:       5 * time.Minute,
	}

	pool := NewWorkerPool(config.MaxConcurrentScans)

	scannedIPs := make(map[string]bool)

	go radioLoop(config, pool, scannedIPs)
}

func radioLoop(config RadioConfig, pool *WorkerPool, scannedIPs map[string]bool) {
	var mu sync.Mutex

	for {
		share, err := shares.GetShare()
		if err != nil {
			log.Printf("Error retrieving share: %v", err)
			time.Sleep(config.ScanInterval)
			continue
		}

		for _, cidr := range share.CIDRs {
			start := time.Now()
			ips, err := ipranger.Ips(cidr)
			if err != nil {
				log.Printf("Error getting IPs from CIDR %s: %v", cidr, err)
				continue
			}

			for _, ip := range ips {
				mu.Lock()
				if scannedIPs[ip] {
					mu.Unlock()
					continue
				}

				scannedIPs[ip] = true
				mu.Unlock()

				pool.Submit(func() {
					isAlive := network.Ping(ip)
					if isAlive {
						network.NetworkScan(ip)
					}
				})
			}

			elapsed := time.Since(start)
			fmt.Printf("Total time taken for the CIDR: %s\n", elapsed)
		}

		time.Sleep(config.ScanInterval)
	}
}
