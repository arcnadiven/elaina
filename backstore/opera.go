package backstore

import (
	"errors"
	"fmt"
	"github.com/arcnadiven/elaina/models"
	"gorm.io/gorm"
)

type StoreOperator interface {
	InsertPersistentVolume(vol *models.CSIPersiVol) error
	QueryPersistentVolume(volId string) (*models.CSIPersiVol, error)
	QueryPersistentVolumeByName(volName string) (*models.CSIPersiVol, error)
	UpdatePersistentVolume(vol *models.CSIPersiVol) error
	DeletePersistentVolume(volId string) error
}

type SQLOperator struct {
	*SQLClient
}

func NewStoreOperator(cli *SQLClient) StoreOperator {
	return &SQLOperator{SQLClient: cli}
}

func (op *SQLOperator) InsertPersistentVolume(vol *models.CSIPersiVol) error {
	return op.client.Model(&models.CSIPersiVol{}).Create(vol).Error
}

func (op *SQLOperator) QueryPersistentVolume(volId string) (*models.CSIPersiVol, error) {
	vol := &models.CSIPersiVol{}
	if err := op.client.Raw(fmt.Sprintf("select * from %s where volume_id = ?", vol.TableName()), volId).Scan(vol).Error; err != nil {
		op.log.Errorln(err)
		return nil, err
	}
	if vol.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return vol, nil
}

func (op *SQLOperator) QueryPersistentVolumeByName(volName string) (*models.CSIPersiVol, error) {
	vol := &models.CSIPersiVol{}
	if err := op.client.Raw(fmt.Sprintf("select * from %s where persi_vol = ?", vol.TableName()), volName).Scan(vol).Error; err != nil {
		op.log.Errorln(err)
		return nil, err
	}
	if vol.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return vol, nil
}

func (op *SQLOperator) UpdatePersistentVolume(vol *models.CSIPersiVol) error {
	if vol.VolumeID == "" {
		return errors.New("update failed volume id cant be empty")
	}
	return op.client.Model(&models.CSIPersiVol{}).Where("volume_id = ?", vol.VolumeID).Updates(vol).Error
}

func (op *SQLOperator) DeletePersistentVolume(volId string) error {
	persi_vol, err := op.QueryPersistentVolume(volId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return op.client.Model(&models.CSIPersiVol{}).Delete(persi_vol).Error
}
