package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/libp2p/zeroconf/v2"
	"github.com/manifoldco/promptui"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/skip2/go-qrcode"
	"github.com/ttacon/chalk"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

const connectServiceName = "_adb-tls-connect._tcp"
const pairServiceName = "_adb-tls-pairing._tcp"
const domain = "local."

type Args = struct {
	Timeout int    `default:"5" help:"Search timeout in seconds"`
	Adb     string `default:"adb" help:"ADB executable path"`
}

type Host struct {
	addr net.IP
	port int
}

func (h *Host) String() string {
	return fmt.Sprintf("%s:%d", h.addr, h.port)
}

type DeviceWithHost struct {
	serviceEntry *zeroconf.ServiceEntry
	label        string
	host         Host
}

func (d DeviceWithHost) String() string {
	return d.label
}

var appArgs Args

func main() {
	parseArgs()

	var hosts []*DeviceWithHost

	for {
		hosts = discover()

		if len(hosts) > 0 || !yesNo("No devices found. Repeat search") {
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

func parseArgs() {
	var args Args
	parser := arg.MustParse(&args)

	if args.Timeout <= 0 {
		parser.Fail("--timeout must be positive")
	}

	appArgs = args
}

func yesNo(question string) bool {
	prompt := promptui.Prompt{
		Label:     question,
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

func execAdb(args ...string) (output string, err error) {
	fmt.Println("Executing", chalk.Blue.Color(appArgs.Adb), "with", chalk.Blue.Color(fmt.Sprint(args)))
	cmd := exec.Command(appArgs.Adb, args...)
	outputBytes, e := cmd.CombinedOutput()

	if e != nil {
		fmt.Println(chalk.Red.Color(e.Error()))
	}
	return string(outputBytes), e
}

func connectToHost(host *DeviceWithHost) {
	fmt.Println("Connecting to", chalk.Blue.Color(host.host.String()))
	output, err := execAdb("connect", host.host.String())
	if err != nil {
		return
	}
	// `adb connect` returns exit code 0 irrespective of whether the connection is established or not.
	if strings.Contains(output, "failed to connect") {
		fmt.Println(chalk.Red.Color(string(output)))
		fmt.Println(chalk.Yellow.Color("Device is probably not paired"))
		promptPair(host)
	} else {
		fmt.Println(chalk.Green.Color(output))
	}
}

func promptPair(host *DeviceWithHost) {
	prompt := promptui.Select{
		Label: "Would you like to pair a device?",
		Items: []string{
			"Yes, pair using QR-code",
			"Yes, pair using pairing code",
			"No",
		},
	}
	i, _, err := prompt.Run()

	if err != nil || i == 2 {
		return
	}

	if i == 0 {
		err = pairUsingQRCode(host)
	} else if i == 1 {
		err = pairUsingPairingCode(host)
	} else {
		err = nil
	}

	if err != nil {
		fmt.Println(chalk.Red.Color(err.Error()))
	}
}

func pairUsingQRCode(host *DeviceWithHost) error {
	name, password, err := genNameAndPassword()

	if err != nil {
		return err
	}

	err = printPairingCode(name, password)

	if err != nil {
		return err
	}

	fmt.Println(
		"Go to",
		chalk.Blue.Color("Developer options"),
		"->",
		chalk.Blue.Color("Wireless debugging"),
		"->",
		chalk.Blue.Color("Pair device with QR code"),
	)

	var pairingHost *Host

	doDiscover(time.Hour, pairServiceName, domain, func(entry *zeroconf.ServiceEntry, cancel context.CancelFunc) {
		if entry.ServiceRecord.Instance == name {
			pairingHost = &Host{entry.AddrIPv4[0], entry.Port}
			cancel()
		}
	})

	if pairingHost == nil {
		return errors.New("pairing failed")
	}

	output, err := execAdb("pair", pairingHost.String(), password)

	if err == nil {
		fmt.Println(chalk.Green.Color(output))
	}

	return err
}

func pairUsingPairingCode(host *DeviceWithHost) error {
	fmt.Println(
		"Go to",
		chalk.Blue.Color("Developer options"),
		"->",
		chalk.Blue.Color("Wireless debugging"),
		"->",
		chalk.Blue.Color("Pair device with pairing code"),
	)
	fmt.Println(chalk.Yellow.Color("Waiting..."))

	var pairingHost *Host

	doDiscover(time.Hour, pairServiceName, domain, func(entry *zeroconf.ServiceEntry, cancel context.CancelFunc) {
		if entry.ServiceRecord.Instance == host.serviceEntry.ServiceRecord.Instance {
			pairingHost = &Host{entry.AddrIPv4[0], entry.Port}
			cancel()
		}
	})

	if pairingHost == nil {
		return errors.New("pairing failed")
	}

	prompt := promptui.Prompt{
		Label: "Enter pairing code",
	}

	code, err := prompt.Run()

	if err != nil {
		return err
	}

	var output string

	output, err = execAdb("pair", pairingHost.String(), code)

	if err == nil {
		fmt.Println(chalk.Green.Color(output))
	}

	return err
}

func discover() []*DeviceWithHost {
	timeout := time.Second * time.Duration(appArgs.Timeout)

	fmt.Println(chalk.Blue, fmt.Sprintf("Searching devices for %d seconds...", int(timeout.Seconds())), chalk.ResetColor)

	devices := doDiscover(timeout, connectServiceName, domain, func(device *zeroconf.ServiceEntry, _ context.CancelFunc) {
		fmt.Println(
			chalk.Green.Color(fmt.Sprintf("Found device: %s %s", device.ServiceRecord.Instance, device.AddrIPv4)),
			"(use Ctrl + C to stop searching)",
		)
	})

	hosts := make([]*DeviceWithHost, 0)

	for _, device := range devices {
		for _, ip := range device.AddrIPv4 {
			host := Host{ip, device.Port}
			hosts = append(hosts, &DeviceWithHost{
				label:        fmt.Sprintf("%s (%s)", host.String(), device.ServiceRecord.Instance),
				serviceEntry: device,
				host:         host,
			})
		}
	}

	return hosts
}

func doDiscover(timeout time.Duration, serviceName string, domain string, onFound func(*zeroconf.ServiceEntry, context.CancelFunc)) []*zeroconf.ServiceEntry {
	entries := make([]*zeroconf.ServiceEntry, 0)

	entriesCh := make(chan *zeroconf.ServiceEntry)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			onFound(entry, cancel)
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

func genNameAndPassword() (name string, password string, err error) {
	name = ""
	password = ""

	uid, err := gonanoid.New()
	if err != nil {
		return
	}
	name = "ADB_WIFI_" + uid
	uid, err = gonanoid.New()
	if err != nil {
		return
	}
	password = uid

	return
}

func printPairingCode(name string, password string) error {
	content := fmt.Sprintf("WIFI:T:ADB;S:%s;P:%s;;", name, password)

	qr, err := qrcode.New(content, qrcode.Low)

	if err != nil {
		return err
	}
	fmt.Println(qr.ToString(false))

	return nil
}
