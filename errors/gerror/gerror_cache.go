package gerror

import (
	"bytes"
	"strconv"
	"sync"
)

var cacheStackTrace = &stackTrace{
	stacks: make([][]byte, 1),
	pcMap:  make(map[uintptr]int),
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

type stackTrace struct {
	mu     sync.RWMutex
	stacks [][]byte
	// map[pc]stacks.idx
	pcMap map[uintptr]int
}

func (st *stackTrace) addAndGetBytesPtr(pc uintptr, buf []byte) *[]byte {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.stacks = append(st.stacks, buf)
	id := len(st.stacks)
	st.pcMap[pc] = id - 1

	return &st.stacks[id-1]
}

func (st *stackTrace) getBytesPtr(pc uintptr) *[]byte {
	st.mu.RLock()
	defer st.mu.RUnlock()
	id, ok := st.pcMap[pc]
	if ok {
		return &st.stacks[id]
	}
	return nil
}

func init() {
	n := 50
	for i := 0; i < n; i++ {
		// 1).
		// 2).
		// ...
		// 50).
		stackTraceFuncIDHeader = append(stackTraceFuncIDHeader, []byte("\n\t"+strconv.Itoa(i+1)+"). "))
	}
}

func getStackTraceFuncIDHeader(i int) []byte {
	if i < 50 {
		return stackTraceFuncIDHeader[i]
	}
	return []byte("\n\t" + strconv.Itoa(i+1) + "). ")
}

var stackTraceFuncIDHeader = [][]byte{}
