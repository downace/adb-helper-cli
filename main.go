package main

import (
	"context"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/libp2p/zeroconf/v2"
	"github.com/manifoldco/promptui"
	"github.com/ttacon/chalk"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

const serviceName = "_adb-tls-connect._tcp"

const domain = "local."

type Args = struct {
	Timeout int `default:"5" help:"Search timeout in seconds"`
}

type DeviceWithHost struct {
	serviceEntry *zeroconf.ServiceEntry
	label        string
	host         string
}

func (d DeviceWithHost) String() string {
	return d.label
}

func main() {
	args := parseArgs()

	var hosts []*DeviceWithHost

	for {
		hosts = discover(time.Second * time.Duration(args.Timeout))

		if len(hosts) > 0 || !askRepeatDiscover() {
			break
		}
	}

	if len(hosts) == 0 {
		return
	}

	host := selectHost(hosts)

	if host == nil {
		return
	}

	connectToHost(host)
}

func parseArgs() Args {
	var args Args
	parser := arg.MustParse(&args)

	if args.Timeout <= 0 {
		parser.Fail("--timeout must be positive")
	}

	return args
}

func askRepeatDiscover() bool {
	prompt := promptui.Prompt{
		Label:     "No devices found. Repeat search",
		IsConfirm: true,
	}
	_, err := prompt.Run()

	return err == nil
}

func selectHost(hosts []*DeviceWithHost) *DeviceWithHost {
	prompt := promptui.Select{
		Label: "Select device to connect",
		Items: hosts,
	}
	i, _, err := prompt.Run()

	if err != nil {
		return nil
	}

	return hosts[i]
}

func connectToHost(host *DeviceWithHost) {
	cmd := exec.Command("adb", "connect", host.host)
	fmt.Println("Connecting via", chalk.Blue.Color(cmd.Path), "with", chalk.Blue.Color(fmt.Sprint(cmd.Args)))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(chalk.Red.Color(string(output)))
		return
	}
	fmt.Println(chalk.Green.Color(string(output)))
}

func discover(timeout time.Duration) []*DeviceWithHost {
	fmt.Println(chalk.Blue, fmt.Sprintf("Searching devices for %d seconds...", int(timeout.Seconds())), chalk.ResetColor)

	devices := doDiscover(timeout, func(device *zeroconf.ServiceEntry) {
		fmt.Println(
			chalk.Green.Color(fmt.Sprintf("Found device: %s %s", device.ServiceRecord.Instance, device.AddrIPv4)),
			"(use Ctrl + C to stop searching)",
		)
	})

	hosts := make([]*DeviceWithHost, 0)

	for _, device := range devices {
		for _, ip := range device.AddrIPv4 {
			host := fmt.Sprintf("%s:%d", ip, device.Port)
			hosts = append(hosts, &DeviceWithHost{
				label:        fmt.Sprintf("%s (%s)", host, device.ServiceRecord.Instance),
				serviceEntry: device,
				host:         host,
			})
		}
	}

	return hosts
}

func doDiscover(timeout time.Duration, onFound func(*zeroconf.ServiceEntry)) []*zeroconf.ServiceEntry {
	entries := make([]*zeroconf.ServiceEntry, 0)

	entriesCh := make(chan *zeroconf.ServiceEntry)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			onFound(entry)
			entries = append(entries, entry)
		}
	}(entriesCh)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		cancel()
	}()

	defer cancel()
	err := zeroconf.Browse(ctx, serviceName, domain, entriesCh)
	if err != nil {
		fmt.Println(chalk.Red.Color("Search failed: " + err.Error()))
	}
	<-ctx.Done()

	return entries
}
