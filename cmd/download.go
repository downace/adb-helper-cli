package cmd

import (
	"downace/adb-helper-cli/internal/adb"
	"fmt"
	"github.com/evilsocket/islazy/zip"
	"github.com/imroc/req/v3"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
	"net/url"
	"os"
	"path/filepath"
)

var platformToolsUrl string

var downloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "Download Android Platform Tools which contains ADB binary",
	Run:     download,
	GroupID: cmdGroupApp,
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&platformToolsUrl, "url", "u", adb.PlatformToolsDownloadUrl, "Platform Tools URL")
}

func download(_ *cobra.Command, _ []string) {
	dlUrl, err := url.Parse(platformToolsUrl)
	if err != nil {
		fmt.Println("Invalid download URL:", err)
		return
	}

	cwd := getCwd()
	downloadsDir := cwd + "/downloads"
	zipPath := downloadsDir + "/" + "platform_tools.zip"

	fmt.Println(
		"Downloading Android Platform Tools from", chalk.Blue.Color(dlUrl.String()),
		"into",
		chalk.Blue.Color(zipPath),
	)
	err = downloadFile(zipPath, dlUrl)

	if err != nil {
		fmt.Println(chalk.Red.Color("Download failed: " + err.Error()))
		return
	}
	fmt.Println(
		"Unpacking",
		chalk.Blue.Color(zipPath),
		"into",
		chalk.Blue.Color(downloadsDir),
	)
	err = unzipFile(zipPath, downloadsDir)

	if err != nil {
		fmt.Println(chalk.Red.Color("Unpacking failed: " + err.Error()))
		return
	}

	err = os.Remove(zipPath)

	if err != nil {
		fmt.Println(
			chalk.Yellow.Color("Unable to remove"),
			chalk.Blue.Color(zipPath)+
				chalk.Yellow.Color(". You can remove it manually"),
		)
	}

	fmt.Println(
		chalk.Green.Color("Done. Now add"),
		chalk.Cyan.Color(downloadsDir+"/platform-tools"),
		chalk.Green.Color("to your Path variable,\nor use"),
		chalk.Cyan.Color("--adb"),
		chalk.Green.Color("flag:"),
		chalk.Cyan.Color("adb-helper --adb "+downloadsDir+"/platform-tools/"+adb.BinaryName+" ..."),
	)
}

func downloadFile(path string, url *url.URL) error {
	client := req.C()

	_, err := client.R().SetOutputFile(path).OnAfterResponse(func(client *req.Client, resp *req.Response) error {
		if resp.Err == nil && !resp.IsSuccessState() {
			resp.Err = fmt.Errorf("bad status: %s", resp.Status)
		}
		return nil
	}).Get(url.String())

	if err != nil {
		return err
	}

	return nil
}

func unzipFile(zipPath string, outPath string) (err error) {
	_, err = zip.Unzip(zipPath, outPath)

	return err
}

func getCwd() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
