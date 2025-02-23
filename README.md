# ADB Helper (CLI)
CLI tool to simplify some ADB operations

Usage:

```shell
adb-helper connect
```

With custom timeout (in seconds, default is 5)

```shell
adb-helper --timeout 15
```

If ADB executable is not in your system Path, you can specify it using `--adb`:

```shell
adb-helper --adb "C:\Android\platform-tools\adb.exe"
```

## Roadmap

- [ ] Simplify interface. Less interactivity, more arguments
- [ ] Add option to download ADB executable
