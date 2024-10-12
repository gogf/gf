package base

type Base struct {
	sBase `gen:"extend"`
}
type sBase struct {
	baseDestory `gen:"extend"`
}

// base Init
func (*sBase) Init() {

}

// base Destory
func (*sBase) Destory() {

}
