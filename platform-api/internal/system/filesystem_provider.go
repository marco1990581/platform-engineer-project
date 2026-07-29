package system

type FilesystemProvider interface {
	ListFilesystems() ([]Filesystem, error)
}

func NewFilesystemProvider() FilesystemProvider {
	return LinuxFilesystemProvider{}
}
