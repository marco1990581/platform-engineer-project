package system

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

const mountInfoPath = "/proc/self/mountinfo"

type LinuxFilesystemProvider struct{}

func (p LinuxFilesystemProvider) ListFilesystems() ([]Filesystem, error) {
	data, err := os.ReadFile(mountInfoPath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	filesystems := make([]Filesystem, 0, len(lines))

	for _, line := range lines {
		mountPoint, filesystemType, err := parseMountInfo(line)
		if err != nil {
			return nil, err
		}

		filesystem, err := filesystemStats(mountPoint, filesystemType)
		if err != nil {
			return nil, err
		}

		filesystems = append(filesystems, filesystem)
	}

	return filesystems, nil
}

func parseMountInfo(line string) (string, string, error) {
	// mountinfo separates optional mount fields from filesystem data with " - ".
	parts := strings.SplitN(line, " - ", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid mountinfo entry: missing separator")
	}

	mountFields := strings.Fields(parts[0])
	filesystemFields := strings.Fields(parts[1])
	if len(mountFields) < 5 || len(filesystemFields) < 1 {
		return "", "", fmt.Errorf("invalid mountinfo entry: insufficient fields")
	}

	return unescapeMountPath(mountFields[4]), filesystemFields[0], nil
}

func filesystemStats(mountPoint, filesystemType string) (Filesystem, error) {
	var stats syscall.Statfs_t
	if err := syscall.Statfs(mountPoint, &stats); err != nil {
		return Filesystem{}, err
	}

	if stats.Bsize <= 0 {
		return Filesystem{}, fmt.Errorf("invalid block size for mount point %q", mountPoint)
	}

	totalBlocks := uint64(stats.Blocks)
	freeBlocks := uint64(stats.Bfree)
	if freeBlocks > totalBlocks {
		return Filesystem{}, fmt.Errorf("invalid free block count for mount point %q", mountPoint)
	}

	blockSize := uint64(stats.Bsize)
	totalBytes := totalBlocks * blockSize
	usedBytes := (totalBlocks - freeBlocks) * blockSize
	// Bavail excludes blocks reserved by the filesystem for privileged users.
	availableBytes := uint64(stats.Bavail) * blockSize

	usagePercent := 0.0
	if totalBytes > 0 {
		usagePercent = float64(usedBytes) / float64(totalBytes) * 100
	}

	return Filesystem{
		Type:         filesystemType,
		MountPoint:   mountPoint,
		TotalBytes:   totalBytes,
		UsedBytes:    usedBytes,
		Available:    availableBytes,
		UsagePercent: usagePercent,
	}, nil
}

func unescapeMountPath(path string) string {
	replacer := strings.NewReplacer(
		`\040`, " ",
		`\011`, "\t",
		`\012`, "\n",
		`\134`, `\`,
	)

	return replacer.Replace(path)
}
