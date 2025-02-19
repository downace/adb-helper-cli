# ADB Helper (CLI)
CLI tool to simplify some ADB operations

Usage:

```shell
adb-helper-cli
```

With custom timeout (in seconds, default is 5)

```shell
adb-helper-cli --timeout 15
```

If ADB executable is not in your system Path, you can specify it using `--adb`:

```shell
adb-helper-cli --adb "C:\Android\platform-tools\adb.exe"
```

## Roadmap

- [ ] Simplify interface. Less interactivity, more arguments
- [ ] Add option to download ADB executable
