package reflect

import (
	stdreflect "reflect"
)

// Like the stdlib reflect.DeepEqual, except ignores order within lists.

// During deepValueEqual, must keep track of checks that are
// in progress.  The comparison algorithm assumes that all
// checks in progress are true when it reencounters them.
// Visited comparisons are stored in a map indexed by visit.
type visit struct {
	a1  uintptr
	a2  uintptr
	typ stdreflect.Type
}

func deepListValueEqualIgnoreOrder(v1, v2 stdreflect.Value, visited map[visit]bool, depth int) bool {
	if v1.Len() != v2.Len() {
		return false
	}
	matched := make(map[int]bool)
	if v1.Kind() == stdreflect.Slice {
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
	}
	for i := 0; i < v1.Len(); i++ {
		foundmatch := false
		for j := 0; j < v2.Len(); j++ {
			if matched[j] {
				continue
			}
			if deepValueEqualIgnoreOrder(v1.Index(i), v2.Index(j), visited, depth+1) {
				foundmatch = true
				matched[j] = true
				break
			}
		}
		if !foundmatch {
			return false
		}
	}
	return len(matched) == v2.Len()
}

// Tests for deep equality using reflected types. The map argument tracks
// comparisons that have already been seen, which allows short circuiting on
// recursive types.
func deepValueEqualIgnoreOrder(v1, v2 stdreflect.Value, visited map[visit]bool, depth int) bool {
	if !v1.IsValid() || !v2.IsValid() {
		return v1.IsValid() == v2.IsValid()
	}
	if v1.Type() != v2.Type() {
		return false
	}

	// if depth > 10 { panic("deepValueEqual") }	// for debugging
	hard := func(k stdreflect.Kind) bool {
		switch k {
		case stdreflect.Array, stdreflect.Map, stdreflect.Slice, stdreflect.Struct:
			return true
		}
		return false
	}

	if v1.CanAddr() && v2.CanAddr() && hard(v1.Kind()) {
		addr1 := v1.UnsafeAddr()
		addr2 := v2.UnsafeAddr()
		if addr1 > addr2 {
			// Canonicalize order to reduce number of entries in visited.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are identical ...
		if addr1 == addr2 {
			return true
		}

		// ... or already seen
		typ := v1.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return true
		}

		// Remember for later.
		visited[v] = true
	}

	switch v1.Kind() {
	case stdreflect.Array:
		return deepListValueEqualIgnoreOrder(v1, v2, visited, depth)
	case stdreflect.Slice:
		return deepListValueEqualIgnoreOrder(v1, v2, visited, depth)
	case stdreflect.Interface:
		if v1.IsNil() || v2.IsNil() {
			return v1.IsNil() == v2.IsNil()
		}
		return deepValueEqualIgnoreOrder(v1.Elem(), v2.Elem(), visited, depth+1)
	case stdreflect.Ptr:
		return deepValueEqualIgnoreOrder(v1.Elem(), v2.Elem(), visited, depth+1)
	case stdreflect.Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {
			if !deepValueEqualIgnoreOrder(v1.Field(i), v2.Field(i), visited, depth+1) {
				return false
			}
		}
		return true
	case stdreflect.Map:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for _, k := range v1.MapKeys() {
			if !deepValueEqualIgnoreOrder(v1.MapIndex(k), v2.MapIndex(k), visited, depth+1) {
				return false
			}
		}
		return true
	case stdreflect.Func:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		// Can't do better than this:
		return false
	default:
		// Normal equality suffices
		// TODO(madadam): stdlib reflect uses an unsafe version of this, but it's not exported.
		// What will break if we don't?  Also not sure what case this catches that isn't caught above
		// by the Interface case.
		//return valueInterface(v1, false) == valueInterface(v2, false)
		return v1.Interface() == v2.Interface()
	}
}

// DeepEqual tests for deep equality. It uses normal == equality where
// possible but will scan elements of arrays, slices, maps, and fields of
// structs. In maps, keys are compared with == but elements use deep
// equality. DeepEqual correctly handles recursive types. Functions are equal
// only if they are both nil.
// An empty slice is not equal to a nil slice.
func DeepEqualIgnoreOrder(a1, a2 interface{}) bool {
	if a1 == nil || a2 == nil {
		return a1 == a2
	}
	v1 := stdreflect.ValueOf(a1)
	v2 := stdreflect.ValueOf(a2)
	if v1.Type() != v2.Type() {
		return false
	}
	return deepValueEqualIgnoreOrder(v1, v2, make(map[visit]bool), 0)
}
