package cmd

import (
	"context"
	"fmt"
	"github.com/downace/adb-helper-cli/internal/adb"
	"github.com/downace/adb-helper-cli/internal/mdns"
	"github.com/downace/adb-helper-cli/internal/ui"
	"github.com/libp2p/zeroconf/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
	"time"
)

var useQrCode bool
var usePairingCode bool

var pairCmd = &cobra.Command{
	Use:     "pair",
	Short:   "Pair device using QR-code or pairing code",
	Run:     pair,
	GroupID: cmdGroupApp,
}

func init() {
	rootCmd.AddCommand(pairCmd)

	pairCmd.Flags().BoolVar(&useQrCode, "qr", false, "use QR-code")
	pairCmd.Flags().BoolVar(&usePairingCode, "code", false, "use pairing code")
	pairCmd.MarkFlagsOneRequired("qr", "code")
	pairCmd.MarkFlagsMutuallyExclusive("qr", "code")
}

func pair(_ *cobra.Command, _ []string) {
	var err error

	if useQrCode {
		err = pairUsingQRCode()
	} else if usePairingCode {
		err = pairUsingPairingCode()
	}

	if err != nil {
		fmt.Println(chalk.Red.Color(err.Error()))
	}
}

func pairUsingQRCode() error {
	name, password, err := genNameAndPassword()

	if err != nil {
		return err
	}

	err = printPairingQrCode(name, password)

	if err != nil {
		return err
	}

	printPairingHelp("Pair device with QR code")

	pairingHost, err := discoverPairingHost()

	if pairingHost == nil {
		return err
	}

	output, err := adb.ExecAdb("pair", pairingHost.String(), password)

	if err == nil {
		fmt.Println(chalk.Green.Color(output))
	}

	return err
}

func pairUsingPairingCode() error {
	printPairingHelp("Pair device with pairing code")

	pairingHost, err := discoverPairingHost()

	if pairingHost == nil {
		return err
	}

	code := ui.StringPrompt("Enter pairing code")

	var output string

	output, err = adb.ExecAdb("pair", pairingHost.String(), code)

	if err == nil {
		fmt.Println(chalk.Green.Color(output))
	}

	return err
}

func printPairingHelp(lastSegment string) {
	fmt.Println(fmt.Sprintf("Go to %s -> %s -> %s",
		chalk.Cyan.Color("Developer options"),
		chalk.Cyan.Color("Wireless debugging"),
		chalk.Cyan.Color(lastSegment),
	))
}

func discoverPairingHost() (*mdns.Host, error) {
	var pairingHost *mdns.Host

	mdns.DiscoverServices(time.Hour, "_adb-tls-pairing._tcp", func(entry *zeroconf.ServiceEntry, stop context.CancelFunc) {
		pairingHost = &mdns.Host{Addr: entry.AddrIPv4[0], Port: entry.Port}
		stop()
	})

	if pairingHost == nil {
		return nil, fmt.Errorf("pairing failed")
	}

	return pairingHost, nil
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

func printPairingQrCode(name string, password string) error {
	content := fmt.Sprintf("WIFI:T:ADB;S:%s;P:%s;;", name, password)

	qr, err := qrcode.New(content, qrcode.Low)

	if err != nil {
		return err
	}
	fmt.Println(qr.ToString(false))

	return nil
}
