package cbs

import (
	"context"
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/von1994/CSIFramework/driver/metrics"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

const (

	// DriverName 为CSI Driver的名称
	DriverName      = "io.von1994.cbs"

	// DriverVersion 为CSI Driver版本号
	DriverVersion  = "1.0.0"

	// TopologyZoneKey 为拓扑相关Label的Key
	TopologyZoneKey = "topology." + DriverName + "/zone"
)

// Driver 为CSI Driver结构体定义
type Driver struct {
	volumeAttachLimit int64
	// kube client
	client kubernetes.Interface
}

// NewDriver 为构建CSI Driver的方法
//  @param volumeAttachLimit
//  @param client
//  @return *Driver
//  @return error
func NewDriver(volumeAttachLimit int64, client kubernetes.Interface) (*Driver, error) {
	driver := Driver{
		volumeAttachLimit: volumeAttachLimit,
		client:            client,
	}

	return &driver, nil
}

// Run 为Driver运行函数
//  @receiver drv
//  @param endpoint
//  @param storageURL
//  @param enableMetricsServer
//  @param metricPort
//  @return error
func (drv *Driver) Run(endpoint *url.URL, storageURL string, enableMetricsServer bool, metricPort int64) error {
	controller, err := newCbsController(storageURL)
	if err != nil {
		return err
	}

	identity, err := newCbsIdentity()
	if err != nil {
		return err
	}

	node, err := newCbsNode(drv.volumeAttachLimit)
	if err != nil {
		return err
	}

	logGRPC := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		glog.Infof("GRPC call: %s, request: %+v", info.FullMethod, req)
		resp, err := handler(ctx, req)
		if err != nil {
			glog.Errorf("GRPC error: %v", err)
		} else {
			glog.Infof("GRPC error: %v, response: %+v", err, resp)
		}
		return resp, err
	}

	if enableMetricsServer {
		// expose driver metrics
		metrics.RegisterMetrics()
		http.Handle("/metrics", promhttp.Handler())
		address := fmt.Sprintf(":%d", metricPort)
		glog.Infof("Starting metrics server at %s\n", address)
		go wait.Forever(func() {
			err := http.ListenAndServe(address, nil)
			if err != nil {
				glog.Errorf("Failed to listen on %s: %v", address, err)
			}
		}, 5*time.Second)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(logGRPC),
	}

	srv := grpc.NewServer(opts...)

	csi.RegisterControllerServer(srv, controller)
	csi.RegisterIdentityServer(srv, identity)
	csi.RegisterNodeServer(srv, node)

	if endpoint.Scheme == "unix" {
		sockPath := path.Join(endpoint.Host, endpoint.Path)
		if _, err := os.Stat(sockPath); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			if err := os.Remove(sockPath); err != nil {
				return err
			}
		}
	}

	listener, err := net.Listen(endpoint.Scheme, path.Join(endpoint.Host, endpoint.Path))
	if err != nil {
		return err
	}

	return srv.Serve(listener)
}
