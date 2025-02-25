module github.com/downace/adb-helper-cli

go 1.24.0

retract (
	v0.3.2 // Contains retractions only
	[v0.0.0, v0.3.0] // Contains invalid module name
)

require (
	github.com/evilsocket/islazy v1.11.0
	github.com/imroc/req/v3 v3.49.1
	github.com/libp2p/zeroconf/v2 v2.2.0
	github.com/matoous/go-nanoid/v2 v2.1.0
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/spf13/cobra v1.9.1
	github.com/ttacon/chalk v0.0.0-20160626202418-22c06c80ed31
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/cloudflare/circl v1.6.0 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/pprof v0.0.0-20250208200701-d0013a598941 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/miekg/dns v1.1.63 // indirect
	github.com/onsi/ginkgo/v2 v2.22.2 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.50.0 // indirect
	github.com/refraction-networking/utls v1.6.7 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	go.uber.org/mock v0.5.0 // indirect
	golang.org/x/crypto v0.34.0 // indirect
	golang.org/x/exp v0.0.0-20250218142911-aa4b98e5adaa // indirect
	golang.org/x/mod v0.23.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/tools v0.30.0 // indirect
)
