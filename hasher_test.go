package freelru

//lint:file-ignore U1000 unused fields are necessary to access the hasher
//lint:file-ignore SA4000 hash code comparisons use identical expressions
// From: https://github.com/puzpuzpuz/xsync/blob/main/util_hash_test.go

import (
	"fmt"
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

	if hashString("foo") != hashString("foo") {
		t.Error("hashString is not deterministic")
	}

	if hashUser(User{Name: "Ivan", City: "Sofia"}) != hashUser(User{Name: "Ivan", City: "Sofia"}) {
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
