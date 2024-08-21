// SyncJobItemStatus.go

package emby

// SyncJobItemStatus
type SyncJobItemStatus string

const (
	Converting      SyncJobItemStatus = "Converting"
	Failed          SyncJobItemStatus = "Failed"
	Queued          SyncJobItemStatus = "Queued"
	ReadyToTransfer SyncJobItemStatus = "ReadyToTransfer"
	Synced          SyncJobItemStatus = "Synced"
	Transferring    SyncJobItemStatus = "Transferring"
)
