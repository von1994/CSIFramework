package cbs

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	storageHttpClient "github.com/von1994/CSIFramework/pkg/storageHttpCLient"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	GB uint64 = 1 << (10 * 3)
	Endpoint = "cbs"
	// controllerCaps represents the capability of controller service
	controllerCaps = []csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
		csi.ControllerServiceCapability_RPC_EXPAND_VOLUME,
	}
)

type cbsController struct {
	httpClient  *storageHttpClient.StorageHttpClient
}

func newCbsController(storageURL string) (*cbsController, error) {
	return &cbsController{
		httpClient:     storageHttpClient.NewStorageHttpClient(storageURL, Endpoint),
	}, nil
}

// CreateVolume 通过调用存储API，创建后端真实存储
//  @receiver ctrl
//  @param ctx
//  @param req
//  @return *csi.CreateVolumeResponse
//  @return error
func (ctrl *cbsController) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	//if req.Name == "" {
	//	return nil, status.Error(codes.InvalidArgument, "volume name is empty")
	//}
	//
	//volumeIdempotencyName := req.Name
	//volumeCapacity := req.CapacityRange.RequiredBytes
	//
	//if len(req.VolumeCapabilities) <= 0 {
	//	return nil, status.Error(codes.InvalidArgument, "volume has no capabilities")
	//}
	//
	//for _, c := range req.VolumeCapabilities {
	//	if c.GetBlock() != nil {
	//		return nil, status.Error(codes.InvalidArgument, "block volume is not supported")
	//	}
	//	if c.AccessMode.Mode != csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER {
	//		return nil, status.Error(codes.InvalidArgument, "block access mode only support singer node writer")
	//	}
	//}

	return &csi.CreateVolumeResponse{}, nil
}

// DeleteVolume 通过调用存储API，删除后端真实存储
//  @receiver ctrl
//  @param ctx
//  @param req
//  @return *csi.DeleteVolumeResponse
//  @return error
func (ctrl *cbsController) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	//if req.VolumeId == "" {
	//	return nil, status.Error(codes.InvalidArgument, "volume id is empty")
	//}

	return &csi.DeleteVolumeResponse{}, nil
}

/*
kubernetes 1.20 Breaking Changes
Kubelet no longer creates the target_path for NodePublishVolume in accordance with the CSI spec.
Kubelet also no longer checks if staging and target paths are mounts or corrupted.
CSI drivers need to be idempotent and do any necessary mount verification.
*/

// ControllerPublishVolume 将创建完毕的后端存储attach到指定节点
//  @receiver ctrl
//  @param ctx
//  @param req
//  @return *csi.ControllerPublishVolumeResponse
//  @return error
func (ctrl *cbsController) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	//diskId := req.VolumeId
	//instanceId := req.NodeId
	return &csi.ControllerPublishVolumeResponse{}, nil
}

// ControllerUnpublishVolume 将后端存储从指定节点detach
//  @receiver ctrl
//  @param ctx
//  @param req
//  @return *csi.ControllerUnpublishVolumeResponse
//  @return error
func (ctrl *cbsController) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

// ControllerGetCapabilities 返回controller能力属性
//  @receiver ctrl
//  @param ctx
//  @param req
//  @return *csi.ControllerGetCapabilitiesResponse
//  @return error
func (ctrl *cbsController) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	glog.Infof("ControllerGetCapabilities: called with args %+v", *req)
	var caps []*csi.ControllerServiceCapability
	for _, cap := range controllerCaps {
		c := &csi.ControllerServiceCapability{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: cap,
				},
			},
		}
		caps = append(caps, c)
	}
	return &csi.ControllerGetCapabilitiesResponse{}, nil
}

func (ctrl *cbsController) ValidateVolumeCapabilities(context.Context, *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (ctrl *cbsController) ListVolumes(context.Context, *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ControllerExpandVolume 扩容后端存储
//  @receiver ctrl
//  @param ctx
//  @param req
//  @return *csi.ControllerExpandVolumeResponse
//  @return error
func (ctrl *cbsController) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	//glog.Infof("ControllerExpandVolume: ControllerExpandVolumeRequest is %v", *req)
	//
	//volumeID := req.GetVolumeId()
	//if len(volumeID) == 0 {
	//	return nil, status.Error(codes.InvalidArgument, "Volume ID not provided")
	//}
	//
	//capacityRange := req.GetCapacityRange()
	//if capacityRange == nil {
	//	return nil, status.Error(codes.InvalidArgument, "Capacity range not provided")
	//}
	//
	//newCbsSizeGB := util.RoundUpGiB(capacityRange.GetRequiredBytes())
	//
	//diskId := req.VolumeId

	return &csi.ControllerExpandVolumeResponse{}, nil
}

func (ctrl *cbsController) GetCapacity(context.Context, *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (ctrl *cbsController) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (ctrl *cbsController) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (ctrl *cbsController) ListSnapshots(context.Context, *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}