package server

import (
	"context"
	"github.com/arcnadiven/elaina/models"
	"github.com/arcnadiven/elaina/tracelog"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	nodeId            = "minikube.deepin.linux"
	maxVolumesPerNode = 100
)

func (es *ElainaServerImpl) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	if req.GetVolumeId() == "" {
		return nil, status.Error(codes.InvalidArgument, "volume id is empty")
	}
	tl := tracelog.NewTraceLogger(es.bl)
	tl.WithValue("volume-id", req.GetVolumeId())
	if req.GetStagingTargetPath() == "" {
		tl.Errorln("cant get the target path")
		return nil, status.Error(codes.InvalidArgument, "target path is empty")
	}
	vol, err := es.opera.QueryPersistentVolume(req.GetVolumeId())
	if err != nil {
		tl.Errorln(err)
		return nil, status.Error(codes.Unavailable, "query persi-vol failed")
	}
	if vol.State != models.Persi_Vol_Attached {
		if vol.State == models.Persi_Vol_Mounted {
			return nil, nil
		}
		tl.Errorf("persi-vol state error: %s", vol.State)
		return nil, status.Error(codes.Unavailable, "persi-vol state error")
	}
	vol.State = models.Persi_Vol_Mounted
	vol.GlobalMount = req.GetStagingTargetPath()
	if err := es.opera.UpdatePersistentVolume(vol); err != nil {
		tl.Errorln(err)
		return nil, status.Error(codes.Unavailable, "update persi-vol failed")
	}
	return nil, nil
}

func (es *ElainaServerImpl) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	if req.GetVolumeId() == "" {
		return nil, status.Error(codes.InvalidArgument, "volume id is empty")
	}
	tl := tracelog.NewTraceLogger(es.bl)
	tl.WithValue("volume-id", req.GetVolumeId())
	if req.GetStagingTargetPath() == "" {
		tl.Errorln("cant get the target path")
		return nil, status.Error(codes.InvalidArgument, "target path is empty")
	}
	vol, err := es.opera.QueryPersistentVolume(req.GetVolumeId())
	if err != nil {
		tl.Errorln("query persi-vol failed")
		return nil, status.Error(codes.Unavailable, "query persi-vol failed")
	}
	if vol.State != models.Persi_Vol_Mounted {
		if vol.State == models.Persi_Vol_Attached {
			return nil, nil
		}
		tl.Errorf("persi-vol state error: %s", vol.State)
		return nil, status.Error(codes.Unavailable, "persi-vol state error")
	}
	vol.State = models.Persi_Vol_Attached
	vol.GlobalMount = ""
	if err := es.opera.UpdatePersistentVolume(vol); err != nil {
		tl.Errorln("update persi-vol failed")
		return nil, status.Error(codes.Unavailable, "update persi-vol failed")
	}
	return nil, status.Error(codes.Unimplemented, "")
}

func (es *ElainaServerImpl) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	if req.GetVolumeId() == "" {
		return nil, status.Error(codes.InvalidArgument, "volume id is empty")
	}
	tl := tracelog.NewTraceLogger(es.bl)
	tl.WithValue("volume-id", req.GetVolumeId())
	if req.GetStagingTargetPath() == "" {
		tl.Errorln("cant get the staging target path")
		return nil, status.Error(codes.InvalidArgument, "staging target path is empty")
	}
	if req.GetTargetPath() == "" {
		tl.Errorln("cant get the target path")
		return nil, status.Error(codes.InvalidArgument, "target path is empty")
	}

	vol, err := es.opera.QueryPersistentVolume(req.GetVolumeId())
	if err != nil {
		tl.Errorln("query persi-vol failed")
		return nil, status.Error(codes.Unavailable, "query persi-vol failed")
	}
	if vol.State != models.Persi_Vol_Mounted {
		if vol.State == models.Persi_Vol_Published {
			return nil, nil
		}
		tl.Errorf("persi-vol state error: %s", vol.State)
		return nil, status.Error(codes.Unavailable, "persi-vol state error")
	}
	vol.State = models.Persi_Vol_Published
	vol.SubMount = req.GetTargetPath()
	if err := es.opera.UpdatePersistentVolume(vol); err != nil {
		tl.Errorln("update persi-vol failed")
		return nil, status.Error(codes.Unavailable, "update persi-vol failed")
	}
	return nil, nil
}

func (es *ElainaServerImpl) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	if req.GetVolumeId() == "" {
		return nil, status.Error(codes.InvalidArgument, "volume id is empty")
	}
	tl := tracelog.NewTraceLogger(es.bl)
	tl.WithValue("volume-id", req.GetVolumeId())
	if req.GetTargetPath() == "" {
		tl.Errorln("cant get the target path")
		return nil, status.Error(codes.InvalidArgument, "target path is empty")
	}
	vol, err := es.opera.QueryPersistentVolume(req.GetVolumeId())
	if err != nil {
		tl.Errorln("query persi-vol failed")
		return nil, status.Error(codes.Unavailable, "query persi-vol failed")
	}
	if vol.State != models.Persi_Vol_Published {
		if vol.State == models.Persi_Vol_Mounted {
			return nil, nil
		}
		tl.Errorf("persi-vol state error: %s", vol.State)
		return nil, status.Error(codes.Unavailable, "persi-vol state error")
	}
	vol.State = models.Persi_Vol_Mounted
	vol.SubMount = ""
	if err := es.opera.UpdatePersistentVolume(vol); err != nil {
		tl.Errorln("update persi-vol failed")
		return nil, status.Error(codes.Unavailable, "update persi-vol failed")
	}
	return nil, nil
}

func (es *ElainaServerImpl) NodeGetVolumeStats(context.Context, *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (es *ElainaServerImpl) NodeExpandVolume(context.Context, *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (es *ElainaServerImpl) NodeGetCapabilities(context.Context, *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	caps := []*csi.NodeServiceCapability{
		{
			Type: &csi.NodeServiceCapability_Rpc{
				Rpc: &csi.NodeServiceCapability_RPC{
					Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
				},
			},
		},
	}
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: caps,
	}, nil
}

func (es *ElainaServerImpl) NodeGetInfo(context.Context, *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{
		NodeId:            nodeId,
		MaxVolumesPerNode: maxVolumesPerNode,
	}, nil
}
