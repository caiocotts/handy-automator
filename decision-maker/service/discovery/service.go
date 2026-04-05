package discovery

import (
	"context"
	"decisionMaker/persistence"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
)

const serviceType = "_handy-automator._tcp"
const domain = "local."

type Service struct {
	mu               sync.RWMutex
	registry         map[string]net.IP // hostname -> current IP
	deviceRepository persistence.DeviceRepository
}

func NewService(r persistence.DeviceRepository) *Service {
	return &Service{
		registry:         make(map[string]net.IP),
		deviceRepository: r,
	}
}

func (s *Service) Start(ctx context.Context) {
	processed := make(chan *zeroconf.ServiceEntry)
	go s.processEntries(ctx, processed)

	for {
		browseCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		entries := make(chan *zeroconf.ServiceEntry)

		resolver, err := zeroconf.NewResolver(nil)
		if err != nil {
			log.Fatal("discovery: failed to create mDNS resolver: ", err)
		}

		go func() {
			for entry := range entries {
				processed <- entry
			}
		}()

		go func() {
			if err := resolver.Browse(browseCtx, serviceType, domain, entries); err != nil {
				log.Println("discovery: browse error: ", err)
			}
		}()

		select {
		case <-ctx.Done():
			cancel()
			return
		case <-browseCtx.Done():
			cancel()
		}
	}
}

func (s *Service) processEntries(ctx context.Context, entries chan *zeroconf.ServiceEntry) {
	for {
		select {
		case <-ctx.Done():
			return
		case entry, ok := <-entries:
			if !ok {
				return
			}
			if len(entry.AddrIPv4) == 0 {
				continue
			}
			hostname := strings.TrimSuffix(entry.HostName, ".")
			ip := entry.AddrIPv4[0]

			s.mu.Lock()
			s.registry[hostname] = ip
			s.mu.Unlock()

			result, err := s.deviceRepository.Upsert(ctx, hostname, ip)
			if err != nil {
				log.Printf("discovery: failed to upsert device %s: %v", hostname, err)
			} else if result.IsNew {
				log.Printf("discovery: registered new device %s at %s", hostname, ip)
			} else if !result.PreviousIP.Equal(ip) {
				log.Printf("discovery: device %s IP changed: %s -> %s", hostname, result.PreviousIP, ip)
			}
		}
	}
}

// Resolve returns the current IP for a given mDNS hostname.
func (s *Service) Resolve(hostname string) (net.IP, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ip, ok := s.registry[hostname]
	return ip, ok
}
