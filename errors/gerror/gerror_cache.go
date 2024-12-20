package gerror

import (
	"strconv"
	"sync"
)

var cacheStackTrace = &stackTrace{
	stacks: make([][]byte, 1),
	pcMap:  make(map[uintptr]int),
}

type stackTrace struct {
	mu     sync.RWMutex
	stacks [][]byte
	// map[pc]stacks.idx
	pcMap map[uintptr]int
}

func (st *stackTrace) add(pc uintptr, buf []byte) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.stacks = append(st.stacks, buf)
	id := len(st.stacks)
	st.pcMap[pc] = id - 1
}

func (st *stackTrace) get(pc uintptr) []byte {
	st.mu.RLock()
	defer st.mu.RUnlock()
	id, ok := st.pcMap[pc]
	if ok {
		return st.stacks[id]
	}
	return nil
}

func init() {
	n := 50
	for i := 0; i < n; i++ {
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
