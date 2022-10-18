package goqjs

import (
	"reflect"
	"testing"
)

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
