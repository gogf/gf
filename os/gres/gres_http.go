// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres

//var (
//	defaultFileTimestamp = time.Now()
//)
//
//// HttpFile implements http.File interface for a no-directory file with content.
//type HttpFile struct {
//	*bytes.Reader
//	io.Closer
//	File
//}
//
//func NewHttpFile(name string, content []byte, timestamp time.Time) *HttpFile {
//	if timestamp.IsZero() {
//		timestamp = defaultFileTimestamp
//	}
//	return &HttpFile{
//		bytes.NewReader(content),
//		ioutil.NopCloser(nil),
//		File{name, false, int64(len(content)), timestamp},
//	}
//}
//
//func (f *HttpFile) Readdir(count int) ([]os.FileInfo, error) {
//	return nil, errors.New("not a directory")
//}
//
//func (f *HttpFile) Size() int64 {
//	return f.S
//}
//
//func (f *HttpFile) Stat() (os.FileInfo, error) {
//	return f.FileInfo(), nil
//}
//
//// HttpDirectory implements http.File interface for a directory
//type HttpDirectory struct {
//	HttpFile
//	ChildrenRead int
//	Children     []os.FileInfo
//}
//
//func NewHttpDirectory(name string, children []string, fs *HttpFileSystem) *HttpDirectory {
//	infos := make([]os.FileInfo, 0, len(children))
//	for _, child := range children {
//		_, err := fs.List(filepath.Join(name, child))
//		infos = append(infos, &File{child, err == nil, 0, time.Time{}})
//	}
//	return &HttpDirectory{
//		HttpFile{
//			bytes.NewReader(nil),
//			ioutil.NopCloser(nil),
//			File{name, true, 0, time.Time{}},
//		},
//		0,
//		infos,
//	}
//}
//
//func (f *HttpDirectory) Readdir(count int) ([]os.FileInfo, error) {
//	if count <= 0 {
//		return f.Children, nil
//	}
//	if f.ChildrenRead+count > len(f.Children) {
//		count = len(f.Children) - f.ChildrenRead
//	}
//	rv := f.Children[f.ChildrenRead : f.ChildrenRead+count]
//	f.ChildrenRead += count
//	return rv, nil
//}
//
//func (f *HttpDirectory) Stat() (os.FileInfo, error) {
//	return f.FileInfo(), nil
//}
//
//// HttpFileSystem implements http.FileSystem, allowing embedded files to be served from net/http package.
//type HttpFileSystem struct {
//	Data   func(path string) ([]byte, error)      // Data should return content of file in path if exists
//	List   func(path string) ([]string, error)    // List should return list of files in the path
//	Info   func(path string) (os.FileInfo, error) // Info should return the info of file in path if exists
//	Prefix string                                 // Prefix would be prepended to http requests
//}
//
//func (fs *HttpFileSystem) Open(name string) (http.File, error) {
//	name = path.Join(fs.Prefix, name)
//	if len(name) > 0 && name[0] == '/' {
//		name = name[1:]
//	}
//	if b, err := fs.Data(name); err == nil {
//		timestamp := defaultFileTimestamp
//		if fs.Info != nil {
//			if info, err := fs.Info(name); err == nil {
//				timestamp = info.ModTime()
//			}
//		}
//		return NewHttpFile(name, b, timestamp), nil
//	}
//	if children, err := fs.List(name); err == nil {
//		return NewHttpDirectory(name, children, fs), nil
//	} else {
//		// If the error is not found, return an error that will
//		// result in a 404 error. Otherwise the server returns
//		// a 500 error for files not found.
//		if strings.Contains(err.Error(), "not found") {
//			return nil, os.ErrNotExist
//		}
//		return nil, err
//	}
//}
