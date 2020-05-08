// want package:"1 unsafe.Pointer usages"
package some_unsafe

import "unsafe"

func UnsafeUsages() {
	b := 42
	ptr := unsafe.Pointer(&b)
	_ = *(*int)(ptr)
}