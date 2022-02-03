package models

import "time"

// if you need define more, see: https://gorm.io/zh_CN/docs/models.html#%E5%AD%97%E6%AE%B5%E6%A0%87%E7%AD%BE
type CSIPersiVol struct {
	ID                 int64         `gorm:"column:id;primaryKey"`
	VolumeID           string        `gorm:"column:volume_id;unique;size:64"`
	PersiVolName       string        `gorm:"column:persi_vol;not null;size:255"`
	RefPersiVolClaim   string        `gorm:"column:ref_persi_vol_claim;size:255"`
	RefPersiVolClaimNS string        `gorm:"column:ref_persi_vol_claim_ns;size:255"`
	GlobalMount        string        `gorm:"column:global_mount"`
	Size               int           `gorm:"column:size;comment:unit: GigaBytes"`
	CreateAt           time.Time     `gorm:"column:create_at;autoCreateTime:milli"`
	UpdateAt           time.Time     `gorm:"column:update_at;autoUpdateTime:milli"`
	State              PersiVolState `gorm:"column:state;size:32"`
}

// define the table name
func (c *CSIPersiVol) TableName() string {
	return "t_csi_persi_vol"
}
