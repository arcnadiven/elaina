package models

type PersiVolState string

const (
	Persi_Vol_Created   PersiVolState = "Created"
	Persi_Vol_Attached  PersiVolState = "Attached"
	Persi_Vol_Mounted   PersiVolState = "Mounted"
	Persi_Vol_Published PersiVolState = "Published"
)
