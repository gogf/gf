package defines

import (
	"io"
	"net/http"
	"os"
)

// A File provides access to a single file.
// The File interface is the minimum implementation required of the file.
// Directory files should also implement [ReadDirFile].
// A file may implement [io.ReaderAt] or [io.Seeker] as optimizations.
type File interface {
	// Name returns the path of the file.
	Name() string
	Open() (io.ReadCloser, error)
	Content() []byte
	FileInfo() os.FileInfo
	Export(dst string, option ...ExportOption) error
	HttpFile() (http.File, error)
}

// FS is the interface that defines a virtual file system.
type FS interface {
	// Get returns the file with given path.
	Get(path string) File

	// IsEmpty checks and returns whether the FS is empty.
	IsEmpty() bool

	// ScanDir returns the files under the given path,
	// the parameter `path` should be a folder type.
	ScanDir(path string, pattern string, recursive ...bool) []File

	ListAll() []File
}

// PackOption contains the extra options for Pack functions.
type PackOption struct {
	Prefix   string // The file path prefix for each file item in resource manager.
	KeepPath bool   // Keep the passed path when packing, usually for relative path.
}

// ExportOption contains options for Export.
type ExportOption struct {
	RemovePrefix string // Remove the prefix from source file before export.
}
