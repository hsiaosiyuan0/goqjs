package goqjs

import (
	"math"
	"reflect"
	"testing"
)

func TestNaN(t *testing.T) {
	assertEqual(t, true, math.IsNaN(NaN()), "")
}

func TestGetTag(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	v := ctx.NewBigInt64(1)
	defer v.Free(ctx)

	assertEqual(t, true, v.Tag() == JS_TAG_BIG_INT, "")
}

func TestNewBool(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	yes := ctx.NewBool(true)
	defer yes.Free(ctx)

	i, _ := yes.Int32(ctx)
	assertEqual(t, true, i == 1, "")

	no := ctx.NewBool(1 == 0)
	defer no.Free(ctx)

	i, _ = no.Int32(ctx)
	assertEqual(t, true, i == 0, "")
}

func TestNewInt32(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	yes := ctx.NewInt32(1)
	defer yes.Free(ctx)

	i, _ := yes.Int32(ctx)
	assertEqual(t, true, i == 1, "")
}

func TestNewUint32(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	yes := ctx.NewUint32(1)
	defer yes.Free(ctx)

	i, _ := yes.Int32(ctx)
	assertEqual(t, true, i == 1, "")
}

func TestNewInt64(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	yes := ctx.NewInt64(1)
	defer yes.Free(ctx)

	i, _ := yes.Int32(ctx)
	assertEqual(t, true, i == 1, "")
}

func TestFloat64(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	yes := ctx.NewFloat64(0.3)
	defer yes.Free(ctx)

	i, _ := yes.Float64(ctx)
	assertEqual(t, true, i == 0.3, "")
}

func TestThrow(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	yes := ctx.NewInt64(1)
	defer yes.Free(ctx)

	ctx.Throw(yes)

	e := ctx.Exception()
	assertEqual(t, true, e.c == yes.c, "")
}

func TestNewError(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	obj := ctx.NewError()
	defer obj.Free(ctx)

	assertEqual(t, true, obj.IsError(ctx), "")
}

func TestDup(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	obj := ctx.NewError()
	defer obj.Free(ctx)

	obj.Dup(ctx)
	defer obj.Free(ctx)

	assertEqual(t, 2, obj.RefCount(), "")
}

func TestAtom(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	a1 := ctx.NewAtom("prop")
	defer a1.Free(ctx)

	a2 := ctx.NewAtom("prop")
	defer a2.Free(ctx)

	assertEqual(t, true, a1 == a2, "")

	assertEqual(t, "prop", a1.String(ctx), "")
}

func TestObject(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	res, err := ctx.Eval(`({a: 1})`, "", EVAL_TYPE_GLOBAL)
	defer res.Free(ctx)

	a := res.GetProp(ctx, "a")
	i, err := a.Int32(ctx)
	assertEqual(t, true, err.IsUndefined(), "")
	assertEqual(t, 1, int(i), "")
}

func TestException(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	_, err := ctx.Eval(`test`, "", EVAL_TYPE_GLOBAL)
	defer err.Free(ctx)

	assertEqual(t, true, err.IsError(ctx), "")
	assertEqual(t, "ReferenceError", err.GetProp(ctx, "name").String(ctx), "")
}

func TestConvertException(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	res, _ := ctx.Eval(`Symbol.for("test")`, "", EVAL_TYPE_GLOBAL)
	defer res.Free(ctx)

	_, err := res.Uint32(ctx)
	defer err.Free(ctx)

	assertEqual(t, true, err.IsError(ctx), "")
	n := err.GetProp(ctx, "name").String(ctx)
	assertEqual(t, "TypeError", n, "")
}

func TestConvert(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	res, _ := ctx.Eval(`[-1, "1", 0.1, 3]`, "", EVAL_TYPE_GLOBAL)
	defer res.Free(ctx)

	i, _ := res.GetPropByIdx(ctx, 0).Uint32(ctx)
	assertEqual(t, 4294967295, int(i), "")

	i1, _ := res.GetPropByIdx(ctx, 1).Index(ctx)
	assertEqual(t, 1, int(i1), "")

	i2, _ := res.GetPropByIdx(ctx, 2).Float64(ctx)
	assertEqual(t, 0.1, float64(i2), "")

	i3 := res.GetPropByIdx(ctx, 3).Bool(ctx)
	assertEqual(t, 1, i3, "")
}

func TestBigintToInt64Fail(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	v := ctx.NewBigInt64(1)
	defer v.Free(ctx)

	_, err := v.Int64(ctx)
	defer err.Free(ctx)

	msg := err.GetProp(ctx, "message").String(ctx)
	assertEqual(t, "cannot convert bigint to number", msg, "")
}

func TestBigintToInt64Ext(t *testing.T) {
	rt := NewRuntime()
	defer rt.Free()

	ctx := rt.NewCtx()
	defer ctx.Free()

	v := ctx.NewBigInt64(1)
	defer v.Free(ctx)

	i, _ := v.Int64Ext(ctx)
	assertEqual(t, int64(1), i, "")
}

func isNilPtr(v interface{}) bool {
	if v == nil {
		return true
	}
	vv := reflect.ValueOf(v)
	return vv.Kind() == reflect.Ptr && vv.IsNil()
}

func assertEqual(t *testing.T, except, actual interface{}, msg string) {
	if except == nil && isNilPtr(actual) {
		return
	}
	if !reflect.DeepEqual(except, actual) {
		t.Fatalf("%s Except: \n%v\nActual: \n%v", msg, except, actual)
	}
}
