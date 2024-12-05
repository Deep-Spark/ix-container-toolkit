module gitee.com/deep-spark/ix-container-runtime

go 1.22

require github.com/opencontainers/runtime-spec v1.2.0

require (
	gitee.com/deep-spark/go-ixml v0.0.1
	github.com/pelletier/go-toml v1.9.5
	github.com/sirupsen/logrus v1.9.3
	github.com/urfave/cli/v2 v2.27.4
	golang.org/x/sys v0.24.0
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
	sigs.k8s.io/yaml v1.4.0
	tags.cncf.io/container-device-interface v0.8.0
	tags.cncf.io/container-device-interface/specs-go v0.8.0
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.4 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/opencontainers/runtime-tools v0.9.1-0.20221107090550-2e043c6bd626 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace gitee.com/deep-spark/go-ixml v0.0.1 => ./../go-ixml
