# ADB Helper (CLI)

CLI tool to aid with Android Wireless debugging (searching, pairing and connecting devices) without Android Studio.

Also check out my desktop app with similar features: https://github.com/downace/adb-helper-desktop

### Installation

Download binary for your platform from [Releases](https://github.com/downace/adb-helper-cli/releases) page,
or use `go install github.com/downace/adb-helper-cli`

### Usage

Search devices and connect interactively:

```shell
adb-helper connect
```

<details><summary>Demo</summary>

[![Connect demo GIF](demo_connect.gif)](https://asciinema.org/a/NzVvQTDaPy2X2MrSeE7rR684S)
</details>

Search devices and connect to the first found device:

```shell
adb-helper --use-first
```

Pair to device using QR-code:

```shell
adb-helper pair --qr
```

<details><summary>Demo</summary>

[![Connect demo GIF](demo_pair_qr.gif)](https://asciinema.org/a/wsnxJjCtC9alC3qerjh7ZK5p1)
</details>

Pair to device using pairing code:

```shell
adb-helper pair --code
```

<details><summary>Demo</summary>

[![Connect demo GIF](demo_pair_code.gif)](https://asciinema.org/a/2z9pxvLFdqlm5kISTVibDTQlG)
</details>

If ADB executable is not in your system Path, you can specify it using `--adb`:

```shell
adb-helper --adb "C:\Android\platform-tools\adb.exe" connect
```

Download ADB executable:

```shell
adb-helper download
```

## Roadmap

- [x] Simplify interface. Less interactivity, more arguments
- [x] Add option to download ADB executable
- [ ] Controllable output (`--quiet` and `--verbose` flags), logging
