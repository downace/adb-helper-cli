package cmd

import (
	"context"
	"fmt"
	"github.com/downace/adb-helper-cli/internal/adb"
	"github.com/downace/adb-helper-cli/internal/mdns"
	"github.com/downace/adb-helper-cli/internal/ui"
	"github.com/libp2p/zeroconf/v2"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
	"strings"
	"time"
)

var timeout uint
var connectFirst bool

var connectCmd = &cobra.Command{
	Use:     "connect",
	Short:   "Search devices and connect to them",
	Run:     connect,
	GroupID: cmdGroupApp,
}

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.Flags().UintVarP(&timeout, "timeout", "t", 5, "Search timeout in seconds. Specify 0 to search indefinitely")
	connectCmd.Flags().BoolVarP(&connectFirst, "use-first", "f", false, "When first device found, stop searching and connect to this device")
}

func connect(_ *cobra.Command, _ []string) {
	var hosts []*mdns.DeviceWithHost

	hosts = discover(time.Second*time.Duration(timeout), connectFirst)

	if len(hosts) == 0 {
		fmt.Println(chalk.Yellow.Color("No devices found"))
		fmt.Println(
			chalk.Magenta.Color("TIP: Ensure that device is paired. You can use"),
			chalk.Cyan.Color("pair"),
			chalk.Magenta.Color("command"),
		)
		return
	}

	var host *mdns.DeviceWithHost

	if connectFirst {
		host = hosts[0]
	} else {
		host = hosts[ui.SelectPrompt("Select device to connect:", hosts)]
	}

	if host == nil {
		return
	}

	connectToHost(host)
}

func discover(timeout time.Duration, stopOnFirst bool) []*mdns.DeviceWithHost {
	fmt.Println(chalk.Blue.Color("Searching devices..."))

	hosts := make([]*mdns.DeviceWithHost, 0)

	showTip := true

	mdns.DiscoverServices(timeout, "_adb-tls-connect._tcp", func(entry *zeroconf.ServiceEntry, stop context.CancelFunc) {
		for _, ip := range entry.AddrIPv4 {
			host := mdns.Host{Addr: ip, Port: entry.Port}
			device := mdns.DeviceWithHost{
				Label:        fmt.Sprintf("%s (%s)", host.String(), entry.ServiceRecord.Instance),
				ServiceEntry: entry,
				Host:         host,
			}
			fmt.Println(chalk.Green.Color(fmt.Sprintf("Device found: %v", device)))
			if showTip {
				fmt.Println(
					chalk.Magenta.Color("TIP: You can add"),
					chalk.Cyan.Color("--use-first"),
					chalk.Magenta.Color("to immediately connect to the first found device"),
				)
				showTip = false
			}
			hosts = append(hosts, &device)
		}

		if stopOnFirst {
			stop()
		}
	})

	return hosts
}

func connectToHost(host *mdns.DeviceWithHost) {
	output, err := adb.ExecAdb("connect", host.Host.String())
	if err != nil {
		return
	}
	// `adb connect` returns exit code 0 irrespective of whether the connection is established or not.
	if strings.Contains(output, "failed to connect") {
		fmt.Println(chalk.Red.Color(output))
		fmt.Println(chalk.Magenta.Color("TIP: Maybe device is not paired? Try using ") + chalk.Blue.Color("pair") + chalk.Magenta.Color(" command"))
	} else {
		fmt.Println(chalk.Green.Color(output))
	}
}
