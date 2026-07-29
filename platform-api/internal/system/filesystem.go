package system

type Filesystems struct {
	Filesystems []Filesystem `json:"filesystems"`
}

type Filesystem struct {
	Type         string  `json:"filesystem_type"`
	MountPoint   string  `json:"mount_point"`
	TotalBytes   uint64  `json:"total_bytes"`
	UsedBytes    uint64  `json:"used_bytes"`
	Available    uint64  `json:"available_bytes"`
	UsagePercent float64 `json:"usage_percent"`
}

func GetFilesystems() (Filesystems, error) {
	provider := NewFilesystemProvider()

	filesystems, err := provider.ListFilesystems()
	if err != nil {
		return Filesystems{}, err
	}

	return Filesystems{Filesystems: filesystems}, nil
}
