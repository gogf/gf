package conv

import (
	"reflect"

	"github.com/gogf/gf/v2/util/gconv"
)

// Convert
func Convert[T any](v interface{}) (result T) {
	t := reflect.TypeOf(result)
	return gconv.Convert(v, t.String()).(T)
}
