
#ifndef GOQJS_BRIDGE_H
#define GOQJS_BRIDGE_H

#include <math.h>

#include "quickjs-libc.h"
#include "quickjs.h"

#if defined(__GNUC__) || defined(__clang__)
#define js_likely(x) __builtin_expect(!!(x), 1)
#define js_unlikely(x) __builtin_expect(!!(x), 0)
#define js_force_inline inline __attribute__((always_inline))
#define __js_printf_like(f, a) __attribute__((format(printf, f, a)))
#else
#define js_likely(x) (x)
#define js_unlikely(x) (x)
#define js_force_inline inline
#define __js_printf_like(a, b)
#endif

static js_force_inline int32_t js_value_get_tag(JSValue v) {
  return JS_VALUE_GET_TAG(v);
}
static js_force_inline int32_t js_value_get_int(JSValue v) {
  return JS_VALUE_GET_INT(v);
}
static js_force_inline int32_t js_value_get_bool(JSValue v) {
  return JS_VALUE_GET_BOOL(v);
}
static js_force_inline double js_value_get_float64(JSValue v) {
  return JS_VALUE_GET_FLOAT64(v);
}
static js_force_inline void* js_value_get_ptr(JSValue v) {
  return JS_VALUE_GET_PTR(v);
}

static js_force_inline JSValue js_make_val(int64_t tag, int32_t v) {
  return JS_MKVAL(tag, v);
}

static js_force_inline JSValue js_make_ptr(int64_t tag, void* v) {
  return JS_MKPTR(tag, v);
}

static js_force_inline int js_tag_is_float64(int64_t tag) {
  return JS_TAG_IS_FLOAT64(tag);
}

static js_force_inline int js_value_is_both_int(JSValue v1, JSValue v2) {
  return JS_VALUE_IS_BOTH_INT(v1, v2);
}

static js_force_inline int js_value_is_both_float(JSValue v1, JSValue v2) {
  return JS_VALUE_IS_BOTH_FLOAT(v1, v2);
}

static js_force_inline JSObject* js_value_get_obj(JSValue v) {
  return JS_VALUE_GET_OBJ(v);
}

static js_force_inline int js_value_has_ref_count(JSValue v) {
  return JS_VALUE_HAS_REF_COUNT(v);
}


static js_force_inline int js_value_get_ref_count(JSValue v) {
  if (JS_VALUE_HAS_REF_COUNT(v)) {
    JSRefCountHeader *p = (JSRefCountHeader *)JS_VALUE_GET_PTR(v);
    return p->ref_count;
  }
  return -1;
}

static js_force_inline JSValue js_nan() { return JS_NAN; }
static js_force_inline JSValue js_undef() { return JS_UNDEFINED; }
static js_force_inline JSValue js_null() { return JS_NULL; }
static js_force_inline JSValue js_false() { return JS_FALSE; }
static js_force_inline JSValue js_true() { return JS_TRUE; }
static js_force_inline JSValue js_exception() { return JS_EXCEPTION; }
static js_force_inline JSValue js_uninitialized() { return JS_UNINITIALIZED; }

static double js_float64_nan() { return JS_FLOAT64_NAN; }

#undef js_unlikely
#undef js_force_inline

#endif /* GOQJS_BRIDGE_H */
