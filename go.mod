module github.com/networkmachinery/networkmachinery-operators

require (
	github.com/coreos/go-semver v0.3.0
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/elazarl/goproxy v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/elazarl/goproxy/ext v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190212212710-3befbb6ad0cc // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	gopkg.in/resty.v1 v1.12.0
	k8s.io/api v0.0.0-20191016110408-35e52d86657a // kubernetes-1.16.2
	k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8 // kubernetes-1.16.2
	k8s.io/cli-runtime v0.0.0-20191016114015-74ad18325ed5 //kubernetes-1.16.2
	k8s.io/client-go v0.0.0-20191016111102-bec269661e48 // v12.0.0 => supporting 1.16
	k8s.io/klog v0.4.0
	sigs.k8s.io/controller-runtime v0.3.0
)

replace (
	github.com/coreos/go-semver => github.com/coreos/go-semver v0.2.0
	k8s.io/api => k8s.io/api v0.0.0-20191016110408-35e52d86657a // kubernetes-1.16.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8 // kubernetes-1.16.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20191016114015-74ad18325ed5 // kubernetes-1.16.2
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190731143132-de47f833b8db
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20191004115455-8e001e5d1894 // kubernetes-1.16.2
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.3.1-0.20191025102020-402e4d3278d3 // commit hash 402e4d3278d3049e01bcef306f81a863ad0d9ee3 (1.16.2)
)

go 1.13
