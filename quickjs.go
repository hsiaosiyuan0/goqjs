package goqjs

/*
#cgo CFLAGS: -I./quickjs/build/include
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/quickjs/build/linux -lquickjs -lm -ldl
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/quickjs/build/darwin -lquickjs -lm -ldl

#include <stdlib.h>
#include <string.h>
#include "quickjs.h"
#include "quickjs-libc.h"
#include "bridge.c"

*/
import "C"
import (
	"unsafe"
)

type Runtime struct {
	c *C.JSRuntime
}

func NewRuntime() Runtime {
	return Runtime{C.JS_NewRuntime()}
}

func (r Runtime) NewCtx() Context {
	return Context{C.JS_NewContext(r.c)}
}

func (r Runtime) Free() {
	C.JS_FreeRuntime(r.c)
}

type Context struct {
	c *C.JSContext
}

const (
	EVAL_TYPE_GLOBAL   int = 0
	EVAL_TYPE_MODULE       = 1
	EVAL_TYPE_DIRECT       = 2
	EVAL_TYPE_INDIRECT     = 3

	EVAL_FLAG_STRICT            = 1 << 3
	EVAL_FLAG_STRIP             = 1 << 4
	EVAL_FLAG_COMPILE_ONLY      = 1 << 5
	EVAL_FLAG_BACKTRACE_BARRIER = 1 << 6
)

func (c Context) Eval(code string, filename string, flags int) (result Value, err Value) {
	c_code := C.CString(code)
	defer C.free(unsafe.Pointer(c_code))

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	c_res := C.JS_Eval(c.c, c_code, C.ulong(len(code)), c_filename, C.int(flags))
	if C.JS_IsException(c_res) == 1 {
		result.c = C.js_undef()
		err.c = C.JS_GetException(c.c)
		return
	}

	result.c = c_res
	err.c = C.js_undef()
	return
}

func (c Context) Free() {
	C.JS_FreeContext(c.c)
}

func (c *Context) Exception() Value {
	return Value{C.JS_GetException(c.c)}
}

type Value struct {
	c C.JSValue
}

func (v Value) Free(ctx Context) {
	C.JS_FreeValue(ctx.c, v.c)
}

func (v Value) IsNumber() bool {
	return C.JS_IsNumber(v.c) == 1
}

func (v Value) IsBigInt(ctx Context) bool {
	return C.JS_IsBigInt(ctx.c, v.c) == 1
}

func (v Value) IsBigFloat() bool {
	return C.JS_IsBigFloat(v.c) == 1
}

func (v Value) IsBigDecimal() bool {
	return C.JS_IsBigDecimal(v.c) == 1
}

func (v Value) IsBool() bool {
	return C.JS_IsBool(v.c) == 1
}

func (v Value) IsNull() bool {
	return C.JS_IsNull(v.c) == 1
}

func (v Value) IsUndefined() bool {
	return C.JS_IsUndefined(v.c) == 1
}

func (v Value) IsException() bool {
	return C.JS_IsException(v.c) == 1
}

func (v Value) IsUninitialized() bool {
	return C.JS_IsUninitialized(v.c) == 1
}

func (v Value) IsString() bool {
	return C.JS_IsString(v.c) == 1
}

func (v Value) IsSymbol() bool {
	return C.JS_IsSymbol(v.c) == 1
}

func (v Value) IsObject() bool {
	return C.JS_IsObject(v.c) == 1
}

func (v Value) IsError(ctx Context) bool {
	return C.JS_IsError(ctx.c, v.c) == 1
}

func (v Value) GetProp(ctx Context, prop string) Value {
	bytes := *(*[]byte)(unsafe.Pointer(&prop))
	c_atom := C.JS_NewAtomLen(ctx.c, (*C.char)(unsafe.Pointer(&bytes[0])), C.ulong(len(prop)))
	defer C.JS_FreeAtom(ctx.c, c_atom)
	return Value{C.JS_GetProperty(ctx.c, v.c, c_atom)}
}

func (v Value) GetPropByIdx(ctx Context, idx int) Value {
	return Value{C.JS_GetPropertyUint32(ctx.c, v.c, C.uint(idx))}
}

func (v Value) String(ctx Context) string {
	c := C.JS_ToCString(ctx.c, v.c)
	defer C.JS_FreeCString(ctx.c, c)
	return C.GoString(c)
}

func (v Value) Bool(ctx Context) int {
	return int(C.JS_ToBool(ctx.c, v.c))
}

func (v Value) Int32(ctx Context) (i int32, err Value) {
	err.c = C.js_undef()
	if C.JS_ToInt32(ctx.c, (*C.int)(&i), v.c) == -1 {
		err.c = C.JS_GetException(ctx.c)
	}
	return
}

func (v Value) Uint32(ctx Context) (i uint32, err Value) {
	err.c = C.js_undef()
	if C.JS_ToUint32(ctx.c, (*C.uint)(&i), v.c) == -1 {
		err.c = C.JS_GetException(ctx.c)
	}
	return
}

func (v Value) Int64(ctx Context) (i int64, err Value) {
	err.c = C.js_undef()
	if C.JS_ToInt64(ctx.c, (*C.longlong)(&i), v.c) == -1 {
		err.c = C.JS_GetException(ctx.c)
	}
	return
}

func (v Value) Index(ctx Context) (i uint64, err Value) {
	err.c = C.js_undef()
	if C.JS_ToIndex(ctx.c, (*C.ulonglong)(&i), v.c) == -1 {
		err.c = C.JS_GetException(ctx.c)
	}
	return
}

func (v Value) Float64(ctx Context) (i float64, err Value) {
	err.c = C.js_undef()
	if C.JS_ToFloat64(ctx.c, (*C.double)(&i), v.c) == -1 {
		err.c = C.JS_GetException(ctx.c)
	}
	return
}
