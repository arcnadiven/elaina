package server

import (
	"context"
	"github.com/arcnadiven/elaina/models"
	"github.com/arcnadiven/elaina/tracelog"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"strings"
)

func (es *ElainaServerImpl) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	// query volume from db
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "volume name is empty")
	}
	if persi_vol, err := es.opera.QueryPersistentVolumeByName(req.GetName()); err != nil && err != gorm.ErrRecordNotFound {
		es.bl.Errorln(err)
		return nil, status.Error(codes.Unavailable, "query persi-vol failed")
	} else if err == nil {
		es.bl.Infof("volume %s already exist", req.Name)
		return &csi.CreateVolumeResponse{
			Volume: &csi.Volume{
				CapacityBytes: persi_vol.Size,
				VolumeId:      persi_vol.VolumeID,
			},
		}, nil
	}

	// if volume name not exist create a new one
	vol_id := strings.ToLower(uuid.New().String())
	tl := tracelog.NewTraceLogger(es.bl)
	tl.WithValue("volume-id", vol_id)

	if req.GetCapacityRange() == nil {
		es.bl.Errorln("get capacity range failed")
		return nil, status.Error(codes.Unavailable, "get capacity range failed")
	}
	persi_vol := &models.CSIPersiVol{
		VolumeID:     vol_id,
		PersiVolName: req.GetName(),
		Size:         req.GetCapacityRange().RequiredBytes,
		State:        models.Persi_Vol_Created,
	}
	if err := es.opera.InsertPersistentVolume(persi_vol); err != nil {
		tl.Errorln(err)
		return nil, status.Error(codes.Unavailable, "insert persi-vol failed")
	}
	tl.Infof("insert persi-vol name: %s success", req.GetName())
	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			CapacityBytes: req.GetCapacityRange().RequiredBytes,
			VolumeId:      vol_id,
		},
	}, nil
}

func (es *ElainaServerImpl) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	// query volume from db
	if req.GetVolumeId() == "" {
		return nil, status.Error(codes.InvalidArgument, "volume id is empty")
	}
	tl := tracelog.NewTraceLogger(es.bl)
	tl.WithValue("volume-id", req.GetVolumeId())
	persi_vol, err := es.opera.QueryPersistentVolume(req.GetVolumeId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			tl.Infof("volume id %s already deleted", req.GetVolumeId())
			return nil, nil
		}
		tl.Errorln(err)
		return nil, status.Error(codes.Unavailable, "query persi-vol failed")
	}

	// query volume state is created
	if persi_vol.State != models.Persi_Vol_Created {
		tl.Errorf("persi-vol state error: %s", persi_vol.State)
		return nil, status.Error(codes.Unavailable, "persi-vol state error")
	}
	// delete volume if volume exist
	if err := es.opera.DeletePersistentVolume(persi_vol.VolumeID); err != nil {
		tl.Errorln("delete persi-vol failed")
		return nil, status.Error(codes.Unavailable, "delete persi-vol failed")
	}
	tl.Infoln("delete persi-vol success")
	return nil, nil
}

func (es *ElainaServerImpl) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	// query volume from db
	if req.GetVolumeId() == "" {
		return nil, status.Error(codes.InvalidArgument, "volume id is empty")
	}
	tl := tracelog.NewTraceLogger(es.bl)
	tl.WithValue("volume-id", req.GetVolumeId())
	if req.GetNodeId() == "" {
		tl.Errorln("cant get the node id")
		return nil, status.Error(codes.InvalidArgument, "node id is empty")
	}
	vol, err := es.opera.QueryPersistentVolume(req.GetVolumeId())
	if err != nil {
		tl.Errorln(err)
		return nil, status.Error(codes.Unavailable, "query persi-vol failed")
	}
	if vol.State != models.Persi_Vol_Created {
		if vol.State == models.Persi_Vol_Attached {
			return nil, nil
		}
		tl.Errorf("persi-vol state error: %s", vol.State)
		return nil, status.Error(codes.Unavailable, "persi-vol state error")
	}
	vol.State = models.Persi_Vol_Attached
	vol.NodeID = req.GetNodeId()
	if err := es.opera.UpdatePersistentVolume(vol); err != nil {
		tl.Errorln(err)
		return nil, status.Error(codes.Unavailable, "update persi-vol failed")
	}
	return nil, nil
}

func (es *ElainaServerImpl) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	// query volume from db
	if req.GetVolumeId() == "" {
		return nil, status.Error(codes.InvalidArgument, "volume id is empty")
	}
	tl := tracelog.NewTraceLogger(es.bl)
	tl.WithValue("volume-id", req.GetVolumeId())
	if req.GetNodeId() == "" {
		tl.Errorln("cant get the node id")
		return nil, status.Error(codes.InvalidArgument, "node id is empty")
	}
	vol, err := es.opera.QueryPersistentVolume(req.GetVolumeId())
	if err != nil {
		tl.Errorln(err)
		return nil, status.Error(codes.Unavailable, "query persi-vol failed")
	}
	if vol.State != models.Persi_Vol_Attached {
		if vol.State == models.Persi_Vol_Created {
			return nil, nil
		}
		tl.Errorf("persi-vol state error: %s", vol.State)
		return nil, status.Error(codes.Unavailable, "persi-vol state error")
	}
	vol.State = models.Persi_Vol_Created
	vol.NodeID = ""
	if err := es.opera.UpdatePersistentVolume(vol); err != nil {
		tl.Errorln(err)
		return nil, status.Error(codes.Unavailable, "update persi-vol failed")
	}
	return nil, nil
}

func (es *ElainaServerImpl) ValidateVolumeCapabilities(context.Context, *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (es *ElainaServerImpl) ControllerGetCapabilities(context.Context, *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	caps := []*csi.ControllerServiceCapability{
		{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
				},
			},
		},
		{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
				},
			},
		},
	}
	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: caps,
	}, nil
}

func (es *ElainaServerImpl) CreateSnapshot(context.Context, *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (es *ElainaServerImpl) GetCapacity(context.Context, *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (es *ElainaServerImpl) ListVolumes(context.Context, *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (es *ElainaServerImpl) DeleteSnapshot(context.Context, *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (es *ElainaServerImpl) ListSnapshots(context.Context, *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (es *ElainaServerImpl) ControllerExpandVolume(context.Context, *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (es *ElainaServerImpl) ControllerGetVolume(context.Context, *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
