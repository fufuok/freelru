package freelru

//lint:file-ignore U1000 unused fields are necessary to access the hasher
//lint:file-ignore SA4000 hash code comparisons use identical expressions
// From: https://github.com/puzpuzpuz/xsync/blob/main/util_hash_test.go

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMakeHashFunc(t *testing.T) {
	type User struct {
		Name string
		City string
	}

	hashString := MakeHasher[string]()
	hashUser := MakeHasher[User]()

	// Not that much to test TBH.
	// check that hash is not always the same
	for i := 0; ; i++ {
		if hashString("foo") != hashString("bar") {
			break
		}
		if i >= 100 {
			t.Error("hashString is always the same")
			break
		}
	}

	a := hashString("foo")
	b := hashString("foo")
	if a != b {
		t.Error("hashString is not deterministic")
	}

	ua := hashUser(User{Name: "Ivan", City: "Sofia"})
	ub := hashUser(User{Name: "Ivan", City: "Sofia"})
	if ua != ub {
		t.Error("hashUser is not deterministic")
	}
}

func BenchmarkMakeHashFunc(b *testing.B) {
	type Point struct {
		X, Y, Z int
	}

	type User struct {
		ID        int
		FirstName string
		LastName  string
		IsActive  bool
		City      string
	}

	type PadInside struct {
		A int
		B byte
		C int
	}

	type PadTrailing struct {
		A int
		B byte
	}

	doBenchmarkMakeHashFunc(b, int64(116))
	doBenchmarkMakeHashFunc(b, int32(116))
	doBenchmarkMakeHashFunc(b, 3.14)
	doBenchmarkMakeHashFunc(b, "test key test key test key test key test key test key test key test key test key ")
	doBenchmarkMakeHashFunc(b, Point{1, 2, 3})
	doBenchmarkMakeHashFunc(b, User{ID: 1, FirstName: "Ivan", LastName: "Ivanov", IsActive: true, City: "Sofia"})
	doBenchmarkMakeHashFunc(b, PadInside{})
	doBenchmarkMakeHashFunc(b, PadTrailing{})
	doBenchmarkMakeHashFunc(b, [1024]byte{})
	doBenchmarkMakeHashFunc(b, [128]Point{})
	doBenchmarkMakeHashFunc(b, [128]User{})
	doBenchmarkMakeHashFunc(b, [128]PadInside{})
	doBenchmarkMakeHashFunc(b, [128]PadTrailing{})
}

func doBenchmarkMakeHashFunc[T comparable](b *testing.B, val T) {
	hash := MakeHasher[T]()
	b.Run(fmt.Sprintf("%T", val), func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = hash(val)
		}
	})
}

type SInt struct {
	A int
	B int64
}

type SWithIface struct {
	X int
	Y interface{}
}

type SWithError struct {
	E error
}

type SNestedIface struct {
	S SWithIface
	Z string
}

type SPtr struct {
	P *interface{}
}

type SSlice []interface{}

type SArray [2]interface{}

type SMap map[string]interface{}

type SChan chan interface{}

type SFuncParam func(interface{}) error

type SFuncNoIface func(int) int

type SFuncRetIface func() interface{}

type SEmbedded struct {
	SInt
}

type SNilPtr struct {
	P *SInt
}

type SMapNonIface map[string]int

type SDeepNested struct {
	A struct {
		B struct {
			C interface{}
		}
	}
}

type SFuncInStruct struct {
	F func(interface{}, int) error
}

type SFuncInStructNoIface struct {
	F func(int) int
}

func TestContainsInterface(t *testing.T) {
	tests := []struct {
		name string
		typ  reflect.Type
		want bool
	}{
		{"plain int struct", reflect.TypeOf(SInt{}), false},
		{"struct with interface", reflect.TypeOf(SWithIface{}), true},
		{"struct with error", reflect.TypeOf(SWithError{}), true},
		{"nested struct with interface", reflect.TypeOf(SNestedIface{}), true},
		{"pointer to interface", reflect.TypeOf(SPtr{}), true},
		{"slice of interface", reflect.TypeOf(SSlice{}), true},
		{"array of interface", reflect.TypeOf(SArray{}), true},
		{"map with interface value", reflect.TypeOf(SMap{}), true},
		{"chan of interface", reflect.TypeOf(SChan(nil)), true},
		{"func param interface", reflect.TypeOf(SFuncParam(nil)), true},
		{"func param no interface", reflect.TypeOf(SFuncNoIface(nil)), false},
		{"func ret interface", reflect.TypeOf(SFuncRetIface(nil)), true},
		{"embedded struct (no interface)", reflect.TypeOf(SEmbedded{}), false},
		{"pointer to struct (no interface)", reflect.TypeOf(&SInt{}), false},
		{"nil pointer field (no interface)", reflect.TypeOf(SNilPtr{}), false},
		{"map without interface", reflect.TypeOf(SMapNonIface{}), false},
		{"plain interface", reflect.TypeOf((*interface{})(nil)).Elem(), true},
		{"plain error", reflect.TypeOf((*error)(nil)).Elem(), true},
		{"slice of ints", reflect.TypeOf([]int{}), false},
		{"deep nested interface", reflect.TypeOf(SDeepNested{}), true},
		{"func in struct (with interface)", reflect.TypeOf(SFuncInStruct{}), true},
		{"func in struct (no interface)", reflect.TypeOf(SFuncInStructNoIface{}), false},
		{"pointer to slice of interface", reflect.TypeOf(&[]interface{}{}), true},
		{"pointer to slice of int", reflect.TypeOf(&[]int{}), false},
		{"map key is interface", reflect.TypeOf(map[interface{}]int{}), true},
		{"map key no interface", reflect.TypeOf(map[int]string{}), false},
		{"array of struct with interface", reflect.TypeOf([2]SWithIface{}), true},
		{"struct with func field returns interface", reflect.TypeOf(struct{ F func() interface{} }{}), true},
		{"struct with func field no interface", reflect.TypeOf(struct{ F func(int) int }{}), false},
		{"direct interface value (int)", reflect.TypeOf(interface{}(42)), false},
		{"direct interface value (string)", reflect.TypeOf(interface{}("hello")), false},
		{"direct error value", reflect.TypeOf(error(nil)), false},
		{"plain interface", reflect.TypeOf((*interface{})(nil)).Elem(), true},
		{"plain error", reflect.TypeOf((*error)(nil)).Elem(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsInterface(tt.typ)
			if got != tt.want {
				t.Errorf("containsInterface(%v) = %v, want %v", tt.typ, got, tt.want)
			}
		})
	}
}

func TestMakeHasher_chanType(t *testing.T) {
	type C chan int
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MakeHasher for chan int should not panic")
		}
	}()
	h := MakeHasher[C]()
	ch := make(chan int)
	h(ch)
}

func TestMakeHasher_chanWithIfaceType(t *testing.T) {
	type C chan interface{}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MakeHasher for chan interface{} should panic")
		}
	}()
	_ = MakeHasher[C]()
}

func TestMakeHasher_IfaceType(t *testing.T) {
	type S struct {
		i any
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MakeHasher for interface{} should panic")
		}
	}()
	_ = MakeHasher[S]()
}

func TestMakeHasher_interface(t *testing.T) {
	type i any
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MakeHasher for interface{} should panic")
		}
	}()
	_ = MakeHasher[i]()
}
