module github.com/snapp-incubator/consul-gslb-driver

go 1.16

require (
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/consul/api v1.12.0
	github.com/kubernetes-csi/csi-lib-utils v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.29.0 // indirect
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.11.0
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.28.0
	k8s.io/klog/v2 v2.9.0
)
