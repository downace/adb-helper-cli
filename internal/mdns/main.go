package mdns

import (
	"context"
	"fmt"
	"github.com/libp2p/zeroconf/v2"
	"github.com/ttacon/chalk"
	"net"
	"time"
)

type Host struct {
	Addr net.IP
	Port int
}

func (h *Host) String() string {
	return fmt.Sprintf("%s:%d", h.Addr, h.Port)
}

type DeviceWithHost struct {
	ServiceEntry *zeroconf.ServiceEntry
	Label        string
	Host         Host
}

func (d DeviceWithHost) String() string {
	return d.Label
}

func DiscoverServices(timeout time.Duration, serviceName string, onDiscover func(entry *zeroconf.ServiceEntry, stop context.CancelFunc)) {
	entriesCh := make(chan *zeroconf.ServiceEntry)

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	if timeout == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	}

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			onDiscover(entry, cancel)
		}
	}(entriesCh)

	defer cancel()
	err := zeroconf.Browse(ctx, serviceName, "local.", entriesCh)
	if err != nil {
		fmt.Println(chalk.Red.Color("Search failed: " + err.Error()))
	}
	<-ctx.Done()
}
