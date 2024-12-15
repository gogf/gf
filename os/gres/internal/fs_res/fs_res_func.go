package fs_res

import (
	"archive/zip"
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/encoding/gcompress"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres/internal/defines"
	"github.com/gogf/gf/v2/text/gstr"
)

const (
	packedGoSourceTemplate = `
package %s

import "github.com/gogf/gf/v2/os/gres"

func init() {
	if err := gres.Add("%s"); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
`
)

// PackWithOption packs the path specified by `srcPaths` into bytes.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
func PackWithOption(srcPaths string, option defines.PackOption) ([]byte, error) {
	var buffer = bytes.NewBuffer(nil)
	err := zipPathWriter(srcPaths, buffer, option)
	if err != nil {
		return nil, err
	}
	// Gzip the data bytes to reduce the size.
	return gcompress.Gzip(buffer.Bytes(), 9)
}

// PackToFileWithOption packs the path specified by `srcPaths` to target file `dstPath`.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
func PackToFileWithOption(srcPaths, dstPath string, option defines.PackOption) error {
	data, err := PackWithOption(srcPaths, option)
	if err != nil {
		return err
	}
	return gfile.PutBytes(dstPath, data)
}

// PackToGoFileWithOption packs the path specified by `srcPaths` to target go file `goFilePath`
// with given package name `pkgName`.
//
// Note that parameter `srcPaths` supports multiple paths join with ','.
func PackToGoFileWithOption(srcPath, goFilePath, pkgName string, option defines.PackOption) error {
	data, err := PackWithOption(srcPath, option)
	if err != nil {
		return err
	}
	return gfile.PutContents(
		goFilePath,
		fmt.Sprintf(gstr.TrimLeft(packedGoSourceTemplate), pkgName, gbase64.EncodeToString(data)),
	)
}

// Unpack unpacks the content specified by `path` to []*File.
func Unpack(path string) ([]defines.File, error) {
	realPath, err := gfile.Search(path)
	if err != nil {
		return nil, err
	}
	return UnpackContent(gfile.GetContents(realPath))
}

// UnpackContent unpacks the content to []File.
func UnpackContent(content string) ([]defines.File, error) {
	var (
		err  error
		data []byte
	)
	if isHexStr(content) {
		// It here keeps compatible with old version packing string using hex string.
		// TODO remove this support in the future.
		data, err = gcompress.UnGzip(hexStrToBytes(content))
		if err != nil {
			return nil, err
		}
	} else if isBase64(content) {
		// New version packing string using base64.
		b, err := gbase64.DecodeString(content)
		if err != nil {
			return nil, err
		}
		data, err = gcompress.UnGzip(b)
		if err != nil {
			return nil, err
		}
	} else {
		data, err = gcompress.UnGzip([]byte(content))
		if err != nil {
			return nil, err
		}
	}
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		err = gerror.Wrapf(err, `create zip reader failed`)
		return nil, err
	}
	var (
		fs    = NewFS()
		array = make([]defines.File, len(reader.File))
	)
	for i, file := range reader.File {
		array[i] = &FileImp{
			file: file,
			fs:   fs,
		}
	}
	return array, nil
}

// isBase64 checks and returns whether given content `s` is base64 string.
// It returns true if `s` is base64 string, or false if not.
func isBase64(s string) bool {
	var r bool
	for i := 0; i < len(s); i++ {
		r = (s[i] >= '0' && s[i] <= '9') ||
			(s[i] >= 'a' && s[i] <= 'z') ||
			(s[i] >= 'A' && s[i] <= 'Z') ||
			(s[i] == '+' || s[i] == '-') ||
			(s[i] == '_' || s[i] == '/') || s[i] == '='
		if !r {
			return false
		}
	}
	return true
}

// isHexStr checks and returns whether given content `s` is hex string.
// It returns true if `s` is hex string, or false if not.
func isHexStr(s string) bool {
	var r bool
	for i := 0; i < len(s); i++ {
		r = (s[i] >= '0' && s[i] <= '9') ||
			(s[i] >= 'a' && s[i] <= 'f') ||
			(s[i] >= 'A' && s[i] <= 'F')
		if !r {
			return false
		}
	}
	return true
}

// hexStrToBytes converts hex string content to []byte.
func hexStrToBytes(s string) []byte {
	src := []byte(s)
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, _ = hex.Decode(dst, src)
	return dst
}
