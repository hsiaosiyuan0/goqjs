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

const (
	/* all tags with a reference count are negative */
	JS_TAG_FIRST             int = -11 /* first negative tag */
	JS_TAG_BIG_DECIMAL           = -11
	JS_TAG_BIG_INT               = -10
	JS_TAG_BIG_FLOAT             = -9
	JS_TAG_SYMBOL                = -8
	JS_TAG_STRING                = -7
	JS_TAG_MODULE                = -3 /* used internally */
	JS_TAG_FUNCTION_BYTECODE     = -2 /* used internally */
	JS_TAG_OBJECT                = -1

	JS_TAG_INT           = 0
	JS_TAG_BOOL          = 1
	JS_TAG_NULL          = 2
	JS_TAG_UNDEFINED     = 3
	JS_TAG_UNINITIALIZED = 4
	JS_TAG_CATCH_OFFSET  = 5
	JS_TAG_EXCEPTION     = 6
	JS_TAG_FLOAT64       = 7
	/* any larger tag is FLOAT64 if JS_NAN_BOXING */
)

func NaN() float64 {
	return float64(C.js_float64_nan())
}

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

	c_res := C.JS_Eval(c.c, c_code, C.size_t(len(code)), c_filename, C.int(flags))
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

func (c *Context) NewBool(v bool) Value {
	i := C.int(*((*uint8)(unsafe.Pointer(&v))))
	return Value{C.JS_NewBool(c.c, i)}
}

func (c *Context) NewInt32(v int32) Value {
	return Value{C.JS_NewInt32(c.c, C.int32_t(v))}
}

func (c *Context) NewUint32(v uint32) Value {
	return Value{C.JS_NewUint32(c.c, C.uint32_t(v))}
}

func (c *Context) NewInt64(v int64) Value {
	return Value{C.JS_NewInt64(c.c, C.int64_t(v))}
}

func (c *Context) NewFloat64(v float64) Value {
	return Value{C.JS_NewFloat64(c.c, C.double(v))}
}

func (c *Context) NewBigInt64(v int64) Value {
	jv := Value{C.JS_NewBigInt64(c.c, C.int64_t(v))}
	return jv
}

func (c *Context) NewBigUint64(v uint64) Value {
	return Value{C.JS_NewBigUint64(c.c, C.uint64_t(v))}
}

// free `JSRuntime::current_exception` at first then set
// set `v` to it.
func (c *Context) Throw(v Value) Value {
	return Value{C.JS_Throw(c.c, v.c)}
}

func (c *Context) NewError() Value {
	return Value{C.JS_NewError(c.c)}
}

func (c *Context) ThrowOutOfMemory() Value {
	return Value{C.JS_ThrowOutOfMemory(c.c)}
}

type JSAtom uint32

func (a JSAtom) Dup(ctx Context) {
	C.JS_DupAtom(ctx.c, C.uint32_t(a))
}

func (a JSAtom) Free(ctx Context) {
	C.JS_FreeAtom(ctx.c, C.uint32_t(a))
}

func (a JSAtom) ToValue(ctx Context) Value {
	return Value{C.JS_AtomToValue(ctx.c, C.uint32_t(a))}
}

func (a JSAtom) ToString(ctx Context) Value {
	return Value{C.JS_AtomToString(ctx.c, C.uint32_t(a))}
}

func (a JSAtom) String(ctx Context) string {
	c := C.JS_AtomToCString(ctx.c, C.uint32_t(a))
	defer C.JS_FreeCString(ctx.c, c)
	return C.GoString(c)
}

func (c *Context) NewAtom(str string) JSAtom {
	// it's safe to use the underlying data point directly without fearing about the point being
	// freed twice since `JS_NewAtomLen` copies its argument internally
	bytes := *(*[]byte)(unsafe.Pointer(&str))
	return JSAtom(C.JS_NewAtomLen(c.c, (*C.char)(unsafe.Pointer(&bytes[0])), C.size_t(len(str))))
}

func (c *Context) NewAtomUint32(n uint32) JSAtom {
	return JSAtom(C.JS_NewAtomUInt32(c.c, C.uint32_t(n)))
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

func (v Value) IsFunction(ctx Context) bool {
	return C.JS_IsFunction(ctx.c, v.c) == 1
}

func (v Value) IsConstructor(ctx Context) bool {
	return C.JS_IsConstructor(ctx.c, v.c) == 1
}

func (v Value) IsArray(ctx Context) bool {
	return C.JS_IsArray(ctx.c, v.c) == 1
}

func (v Value) GetProp(ctx Context, prop string) Value {
	// it's safe to use the underlying data point directly without fearing about the point being
	// freed twice since `JS_NewAtomLen` copies its argument internally
	bytes := *(*[]byte)(unsafe.Pointer(&prop))
	c_atom := C.JS_NewAtomLen(ctx.c, (*C.char)(unsafe.Pointer(&bytes[0])), C.size_t(len(prop)))
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
	if C.JS_ToInt32(ctx.c, (*C.int32_t)(&i), v.c) == -1 {
		err.c = C.JS_GetException(ctx.c)
	}
	return
}

func (v Value) Uint32(ctx Context) (i uint32, err Value) {
	err.c = C.js_undef()
	if C.JS_ToUint32(ctx.c, (*C.uint32_t)(&i), v.c) == -1 {
		err.c = C.JS_GetException(ctx.c)
	}
	return
}

// use `Int64Ext` instead if v is bigint
func (v Value) Int64(ctx Context) (i int64, err Value) {
	err.c = C.js_undef()
	if C.JS_ToInt64(ctx.c, (*C.int64_t)(&i), v.c) == -1 {
		err.c = C.JS_GetException(ctx.c)
	}
	return
}

func (v Value) Index(ctx Context) (i uint64, err Value) {
	err.c = C.js_undef()
	if C.JS_ToIndex(ctx.c, (*C.uint64_t)(&i), v.c) == -1 {
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

// v => bigint => int64
func (v Value) Bigint64(ctx Context) (i int64, err Value) {
	err.c = C.js_undef()
	if C.JS_ToBigInt64(ctx.c, (*C.int64_t)(&i), v.c) == -1 {
		err.c = C.JS_GetException(ctx.c)
	}
	return
}

// internal logic is:
//
// - bigint => int64, if v is bigint
// - v => int64
func (v Value) Int64Ext(ctx Context) (i int64, err Value) {
	err.c = C.js_undef()
	if C.JS_ToInt64Ext(ctx.c, (*C.int64_t)(&i), v.c) == -1 {
		err.c = C.JS_GetException(ctx.c)
	}
	return
}

func (v Value) Tag() int {
	return int(C.js_value_get_tag(v.c))
}

func (v Value) RefCount() int {
	return int(C.js_value_get_ref_count(v.c))
}

func (v Value) Dup(ctx Context) Value {
	C.JS_DupValue(ctx.c, v.c)
	return v
}

/* flags for object properties */
const (
	JS_PROP_CONFIGURABLE int = 1 << 0
	JS_PROP_WRITABLE         = 1 << 1
	JS_PROP_ENUMERABLE       = 1 << 2
	JS_PROP_C_W_E            = JS_PROP_CONFIGURABLE | JS_PROP_WRITABLE | JS_PROP_ENUMERABLE
	JS_PROP_LENGTH           = 1 << 3 /* used internally in Arrays */
	JS_PROP_TMASK            = 3 << 4 /* mask for NORMAL, GETSET, VARREF, AUTOINIT */
	JS_PROP_NORMAL           = 0 << 4
	JS_PROP_GETSET           = 1 << 4
	JS_PROP_VARREF           = 2 << 4 /* used internally */
	JS_PROP_AUTOINIT         = 3 << 4 /* used internally */
)

/* flags for JS_DefineProperty */
const (
	JS_PROP_HAS_SHIFT        int = 8
	JS_PROP_HAS_CONFIGURABLE     = 1 << 8
	JS_PROP_HAS_WRITABLE         = 1 << 9
	JS_PROP_HAS_ENUMERABLE       = 1 << 10
	JS_PROP_HAS_GET              = 1 << 11
	JS_PROP_HAS_SET              = 1 << 12
	JS_PROP_HAS_VALUE            = 1 << 13

	/* throw an exception if false would be returned
	   (JS_DefineProperty/JS_SetProperty) */
	JS_PROP_THROW = 1 << 14
	/* throw an exception if false would be returned in strict mode
	   (JS_SetProperty) */
	JS_PROP_THROW_STRICT = 1 << 15

	JS_PROP_NO_ADD    = 1 << 16 /* internal use */
	JS_PROP_NO_EXOTIC = 1 << 17 /* internal use */
)

const JS_DEFAULT_STACK_SIZE = 256 * 1024
