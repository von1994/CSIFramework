package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/url"

	"github.com/golang/glog"

	"github.com/von1994/CSIFramework/driver/cbs"
)

var (
	endpoint            = flag.String("endpoint", fmt.Sprintf("unix:///var/lib/kubelet/plugins/%s/csi.sock", cbs.DriverName), "CSI endpoint")
	storageURL          = flag.String("storage_url", "xxx.xxx.xxx", "storage api domain")
	volumeAttachLimit   = flag.Int64("volume_attach_limit", -1, "Value for the maximum number of volumes attachable for all nodes. If the flag is not specified then the value is default 20.")
	metricsServerEnable = flag.Bool("enable_metrics_server", true, "enable metrics server, set `false` to close it.")
	metricsPort         = flag.Int64("metric_port", 9099, "metric port")
	master              = flag.String("master", "", "Master URL to build a client config from. Either this or kubeconfig needs to be set if the provisioner is being run out of cluster.")
	kubeconfig          = flag.String("kubeconfig", "", "Absolute path to the kubeconfig file. Either this or master needs to be set if the provisioner is being run out of cluster.")
)

func main() {
	flag.Parse()
	defer glog.Flush()

	var config *rest.Config
	var err error
	if *master != "" || *kubeconfig != "" {
		glog.Infof("Either master or kubeconfig specified. building kube config from that..")
		config, err = clientcmd.BuildConfigFromFlags(*master, *kubeconfig)
	} else {
		glog.Infof("Building kube configs for running in cluster...")
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		glog.Fatalf("Failed to create config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Failed to create client: %v", err)
	}

	u, err := url.Parse(*endpoint)
	if err != nil {
		glog.Fatalf("parse endpoint err: %s", err.Error())
	}

	if u.Scheme != "unix" {
		glog.Fatal("only unix socket is supported currently")
	}

	drv, err := cbs.NewDriver(*volumeAttachLimit, clientset)
	if err != nil {
		glog.Fatal(err)
	}

	if err := drv.Run(u, *storageURL, *metricsServerEnable, *metricsPort); err != nil {
		glog.Fatal(err)
	}

	return
}
