package anytype_test

import (
	"testing"
	"time"

	"github.com/DanielSvub/anytype"
)

type object = anytype.Object
type list = anytype.List

var (
	Object     = anytype.NewObject
	ObjectFrom = anytype.NewObjectFrom
	List       = anytype.NewList
	ListOf     = anytype.NewListOf
	ListFrom   = anytype.NewListFrom
)

func TestValues(t *testing.T) {

	t.Run("equality", func(t *testing.T) {
		if List(1).Equals(List(3.14)) {
			t.Error("equality between atomic types does not work properly")
		}
		if List(3.14).Equals(List(nil)) {
			t.Error("equality between atomic types does not work properly")
		}
		if List("test").Equals(List(false)) {
			t.Error("equality between atomic types does not work properly")
		}
		if List(false).Equals(List(1)) {
			t.Error("equality between atomic types does not work properly")
		}
	})

	t.Run("objectParsing", func(t *testing.T) {
		if !List(map[string]any{"first": 1}).Equals(List(Object("first", 1))) {
			t.Error("value parsing does not work properly")
		}
		if !List(map[string]object{"first": Object()}).Equals(List(Object("first", Object()))) {
			t.Error("value parsing does not work properly")
		}
		if !List(map[string]list{"first": List()}).Equals(List(Object("first", List()))) {
			t.Error("value parsing does not work properly")
		}
		if !List(map[string]string{"first": "test"}).Equals(List(Object("first", "test"))) {
			t.Error("value parsing does not work properly")
		}
		if !List(map[string]bool{"first": false}).Equals(List(Object("first", false))) {
			t.Error("value parsing does not work properly")
		}
		if !List(map[string]int{"first": 1}).Equals(List(Object("first", 1))) {
			t.Error("value parsing does not work properly")
		}
		if !List(map[string]float64{"first": 3.14}).Equals(List(Object("first", 3.14))) {
			t.Error("value parsing does not work properly")
		}
	})

	t.Run("listParsing", func(t *testing.T) {
		if !List([]any{1}).Equals(List(List(1))) {
			t.Error("value parsing does not work properly")
		}
		if !List([]object{Object()}).Equals(List(List(Object()))) {
			t.Error("value parsing does not work properly")
		}
		if !List([]list{List()}).Equals(List(List(List()))) {
			t.Error("value parsing does not work properly")
		}
		if !List([]string{"test"}).Equals(List(List("test"))) {
			t.Error("value parsing does not work properly")
		}
		if !List([]bool{false}).Equals(List(List(false))) {
			t.Error("value parsing does not work properly")
		}
		if !List([]int{1}).Equals(List(List(1))) {
			t.Error("value parsing does not work properly")
		}
		if !List([]float64{3.14}).Equals(List(List(3.14))) {
			t.Error("value parsing does not work properly")
		}
	})

	t.Run("numbersParsing", func(t *testing.T) {
		if !List(int(1)).Equals(List(1)) {
			t.Error("value parsing does not work properly")
		}
		if !List(int64(1)).Equals(List(1)) {
			t.Error("value parsing does not work properly")
		}
		if !List(int32(1)).Equals(List(1)) {
			t.Error("value parsing does not work properly")
		}
		if !List(int16(1)).Equals(List(1)) {
			t.Error("value parsing does not work properly")
		}
		if !List(int8(1)).Equals(List(1)) {
			t.Error("value parsing does not work properly")
		}
		if !List(uint(1)).Equals(List(1)) {
			t.Error("value parsing does not work properly")
		}
		if !List(uint64(1)).Equals(List(1)) {
			t.Error("value parsing does not work properly")
		}
		if !List(uint32(1)).Equals(List(1)) {
			t.Error("value parsing does not work properly")
		}
		if !List(uint16(1)).Equals(List(1)) {
			t.Error("value parsing does not work properly")
		}
		if !List(uint8(1)).Equals(List(1)) {
			t.Error("value parsing does not work properly")
		}
		if !List(float64(3.14)).Equals(List(3.14)) {
			t.Error("value parsing does not work properly")
		}
		if !List(float32(1.0)).Equals(List(1.0)) {
			t.Error("value parsing does not work properly")
		}
	})

}

func TestValuePanics(t *testing.T) {

	catch := func(msg string) {
		if r := recover(); r == nil {
			t.Error(msg)
		}
	}

	t.Run("incompatible", func(t *testing.T) {
		defer catch("parsing an incompatible type does not cause panic")
		List(time.Now())
	})

}

func TestTF(t *testing.T) {

	t.Run("valid", func(t *testing.T) {
		l := List()
		l.SetTF("#0", 0)
		l.SetTF("#0#0", 0)
		l.SetTF("#0#0.test", 1)
		l.SetTF("#0#0.test", 2)
		l.SetTF("#0#1", 2).SetTF("#0#2.test.inner#0", 1)
		l.SetTF("#0#1", 3)
		l.SetTF("#0#2.test.inner#1.target#0", 5)
		result, ok := l.GetTF("#0#2.test.inner#1.target#0").(int)
		if !ok || result != 5 {
			t.Error("tree form handling does not work properly")
		}
		if l.GetTF("#0#2.test") == nil {
			t.Error("valid TF returns nil")
		}
	})

	t.Run("types", func(t *testing.T) {
		l := List(
			Object("test", 0, "list", List("test")),
			List(0),
		)
		if l.TypeOfTF("#1#0") != anytype.TypeInt {
			t.Error("tree form type check does not work properly")
		}
		if l.TypeOfTF("#0.test") != anytype.TypeInt {
			t.Error("tree form type check does not work properly")
		}
		if l.TypeOfTF("#0.list#0") != anytype.TypeString {
			t.Error("invalid TF is not undefined")
		}
	})

	t.Run("unset", func(t *testing.T) {
		l := List(
			Object(
				"test", 0,
				"list", List("test"),
				"object", Object("test", 0),
			),
			0,
			1,
			List(0),
		)
		if l.UnsetTF("#1").GetInt(1) != 1 {
			t.Error("unsetting by TF does not work properly")
		}
		if l.UnsetTF("#0.test").GetObject(0).TypeOf("test") != anytype.TypeUndefined {
			t.Error("unsetting by TF does not work properly")
		}
		if !l.UnsetTF("#2#0").GetList(2).Empty() {
			t.Error("unsetting by TF does not work properly")
		}
		if !l.UnsetTF("#0.list#0").GetObject(0).GetList("list").Empty() {
			t.Error("unsetting by TF does not work properly")
		}
		if !l.UnsetTF("#0.object.test").GetObject(0).GetObject("object").Empty() {
			t.Error("unsetting by TF does not work properly")
		}
	})

	t.Run("invalid", func(t *testing.T) {
		l := List().SetTF("#0#0.test.inner#0", 5)
		if Object().TypeOfTF("invalid") != anytype.TypeUndefined {
			t.Error("invalid TF is not undefined")
		}
		if l.TypeOfTF("invalid") != anytype.TypeUndefined {
			t.Error("invalid TF is not undefined")
		}
		if l.TypeOfTF("invalid") != anytype.TypeUndefined {
			t.Error("invalid TF is not undefined")
		}
		if l.TypeOfTF("#5") != anytype.TypeUndefined {
			t.Error("invalid TF is not undefined")
		}
		if l.TypeOfTF("#0#0#0") != anytype.TypeUndefined {
			t.Error("invalid TF is not undefined")
		}
		if l.TypeOfTF("#0#0.nil") != anytype.TypeUndefined {
			t.Error("invalid TF is not undefined")
		}
		if l.TypeOfTF("#0#0.test#0") != anytype.TypeUndefined {
			t.Error("invalid TF is not undefined")
		}
		if l.TypeOfTF("#0.test") != anytype.TypeUndefined {
			t.Error("invalid TF is not undefined")
		}
		if l.TypeOfTF("#0#0.test.inner.nil") != anytype.TypeUndefined {
			t.Error("invalid TF is not undefined")
		}
		if l.TypeOfTF("#invalid") != anytype.TypeUndefined {
			t.Error("invalid TF is not undefined")
		}
		if l.TypeOfTF("#invalid.test") != anytype.TypeUndefined {
			t.Error("invalid TF is not undefined")
		}
		if l.TypeOfTF("#invalid#0") != anytype.TypeUndefined {
			t.Error("invalid TF is not undefined")
		}
	})

}
