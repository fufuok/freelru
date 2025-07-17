package freelru

import (
	"reflect"
	"unsafe"
)

// MakeHasher creates a fast hash function for the given comparable type.
// The only limitation is that the type should not contain interfaces inside
// based on runtime.typehash.
// Thanks to @puzpuzpuz for the idea in xsync.
func MakeHasher[T comparable]() func(T) uint32 {
	var zero T
	if containsInterface(reflect.TypeOf(&zero).Elem()) {
		panic("freelru: interface types require custom hash functions")
	}

	seed := makeSeed()
	if reflect.TypeOf(&zero).Elem().Kind() == reflect.Interface {
		return func(value T) uint32 {
			iValue := any(value)
			i := (*iface)(unsafe.Pointer(&iValue))
			return runtimeTypehash32(i.typ, i.word, seed)
		}
	}
	var iZero any = zero
	i := (*iface)(unsafe.Pointer(&iZero))
	return func(value T) uint32 {
		return runtimeTypehash32(i.typ, unsafe.Pointer(&value), seed)
	}
}

func containsInterface(t reflect.Type) bool {
	visited := map[reflect.Type]bool{}
	var check func(reflect.Type) bool
	check = func(t reflect.Type) bool {
		if t == nil {
			return false
		}
		// Prevent recursion loops
		if visited[t] {
			return false
		}
		visited[t] = true

		switch t.Kind() {
		case reflect.Interface:
			return true
		case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Chan:
			return check(t.Elem())
		case reflect.Map:
			return check(t.Key()) || check(t.Elem())
		case reflect.Struct:
			for i := 0; i < t.NumField(); i++ {
				if check(t.Field(i).Type) {
					return true
				}
			}
			return false
		case reflect.Func:
			for i := 0; i < t.NumIn(); i++ {
				if check(t.In(i)) {
					return true
				}
			}
			for i := 0; i < t.NumOut(); i++ {
				if check(t.Out(i)) {
					return true
				}
			}
			return false
		default:
			return false
		}
	}
	return check(t)
}

// makeSeed creates a random seed.
func makeSeed() uint32 {
	var s1 uint32
	for {
		s1 = runtimeFastrand()
		// We use seed 0 to indicate an uninitialized seed/hash,
		// so keep trying until we get a non-zero seed.
		if s1 != 0 {
			break
		}
	}
	return s1
}

// how interface is represented in memory
type iface struct {
	typ  uintptr
	word unsafe.Pointer
}

// same as runtimeTypehash, but always returns a uint64
// see: maphash.rthash function for details
func runtimeTypehash32(t uintptr, p unsafe.Pointer, seed uint32) uint32 {
	return uint32(runtimeTypehash(t, p, uintptr(seed)))
}

//go:noescape
//go:linkname runtimeTypehash runtime.typehash
func runtimeTypehash(t uintptr, p unsafe.Pointer, h uintptr) uintptr

//go:noescape
//go:linkname runtimeFastrand runtime.fastrand
func runtimeFastrand() uint32
