package gerror

import (
	"bytes"
	"strconv"
	"sync"
)

var cacheStackTrace = &stackTrace{
	stackBufs: make([][]byte, 0),
	pcMap:     make(map[uintptr]int),
}

var bytesBufferPool = &sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func getBytesBuffer() *bytes.Buffer {
	return bytesBufferPool.Get().(*bytes.Buffer)
}

func putBytesBuffer(b *bytes.Buffer) {
	b.Reset()
	bytesBufferPool.Put(b)
}

type stackObjectFlag int32

const (
	// consts.StackFilterKeyForGoFrame = "github.com/gogf/gf/"
	stackObjectFlagStackFilterKeyForGoFrame stackObjectFlag = 1
	// stackFilterKeyLocal = "/errors/gerror/gerror"
	stackObjectFlagStackFilterKeyLocal stackObjectFlag = 2
)

type stackObject struct {
	// stackBufs[bufId]
	stackBufId int32
	flag       stackObjectFlag
}

type stackTrace struct {
	mu           sync.RWMutex
	stackBufs    [][]byte
	stackObjects []stackObject
	// map[pc]stackBufs.idx
	pcMap map[uintptr]int
	//pcMap sync.Map
}

func (st *stackTrace) addAndGetBytesPtr(pc uintptr, buf []byte, flag stackObjectFlag) *[]byte {
	st.mu.Lock()
	defer st.mu.Unlock()

	v, ok := st.pcMap[pc]
	if ok {
		return &st.stackBufs[v]
	}

	st.stackBufs = append(st.stackBufs, buf)
	id := len(st.stackBufs)
	st.stackObjects = append(st.stackObjects, stackObject{
		stackBufId: int32(id - 1),
		flag:       stackObjectFlag(flag),
	})
	st.pcMap[pc] = id - 1
	return &st.stackBufs[id-1]
}

func (st *stackTrace) getBytesPtr(pc uintptr) (*[]byte, stackObjectFlag) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	idx, ok := st.pcMap[pc]
	if ok {
		flag := st.stackObjects[idx].flag
		return &st.stackBufs[idx], flag
	}
	return nil, -1
}

func init() {
	n := 50
	lineIdPrefix = make([][]byte, n+1)
	for i := 0; i <= n; i++ {
		// 1).
		// 2).
		// ...
		// 50).
		lineIdPrefix[i] = []byte("\n\t" + strconv.Itoa(i+1) + "). ")
	}

	idStrs = make([]string, n+1)
	for i := 0; i <= n; i++ {
		idStrs[i] = strconv.Itoa(i)
	}
}

func getLineIdPrefix(i int) []byte {
	if i < len(lineIdPrefix) {
		return lineIdPrefix[i]
	}
	return []byte("\n\t" + strconv.Itoa(i+1) + "). ")
}

var (
	lineIdPrefix [][]byte
	idStrs       []string
)

func getIdStr(id int) string {
	if id < len(idStrs) {
		return idStrs[id]
	}
	return strconv.Itoa(id)
}
