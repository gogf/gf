package gerror

import (
	"bytes"
	"strconv"
	"sync"
)

var cacheStackTrace = &stackTrace{}

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
	stackBuf []byte
	flag     stackObjectFlag
}

type stackTrace struct {
	pcMap sync.Map
	// pcMap map[uintptr]stackObject
}

func (st *stackTrace) addStackObject(pc uintptr, stackObj stackObject) {
	st.pcMap.LoadOrStore(pc, stackObj)
}

func (st *stackTrace) getStackObject(pc uintptr) stackObject {
	v, ok := st.pcMap.Load(pc)
	if ok {
		return v.(stackObject)
	}
	return stackObject{nil, 0}
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
		idStrs[i] = strconv.Itoa(i) + ". "
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
