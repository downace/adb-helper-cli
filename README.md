# ADB Helper (CLI)
CLI tool to simplify some ADB operations

### Usage

Search devices and connect interactively:

```shell
adb-helper connect
```

Search devices and connect to the first found device:

```shell
adb-helper --use-first
```

Pair to device using QR-code:

```shell
adb-helper pair --qr
```

Pair to device using pairing code:

```shell
adb-helper pair --code
```

If ADB executable is not in your system Path, you can specify it using `--adb`:

```shell
adb-helper --adb "C:\Android\platform-tools\adb.exe" connect
```

## Roadmap

- [x] Simplify interface. Less interactivity, more arguments
- [ ] Add option to download ADB executable
