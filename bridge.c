
#ifndef GOQJS_BRIDGE_H
#define GOQJS_BRIDGE_H

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

static js_force_inline JSValue js_undef() { return JS_UNDEFINED; }
static js_force_inline JSValue js_null() { return JS_NULL; }
static js_force_inline JSValue js_false() { return JS_FALSE; }
static js_force_inline JSValue js_true() { return JS_TRUE; }
static js_force_inline JSValue js_exception() { return JS_EXCEPTION; }
static js_force_inline JSValue js_uninitialized() { return JS_UNINITIALIZED; }

#undef js_unlikely
#undef js_force_inline

#endif /* GOQJS_BRIDGE_H */
