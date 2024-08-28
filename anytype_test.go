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
		l.SetTF("#0#0", 0)
		l.SetTF("#0#0.test", 1)
		l.SetTF("#0#1", 2).SetTF("#0#2.test.inner#0", 1)
		l.SetTF("#0#2.test.inner#1.target#0", 5)
		result, ok := l.GetTF("#0#2.test.inner#1.target#0").(int)
		if !ok || result != 5 {
			t.Error("tree form handling does not work properly")
		}
		if l.GetTF("#0#2.test") == nil {
			t.Error("valid TF returns nil")
		}
	})

	t.Run("invalid", func(t *testing.T) {
		l := List().SetTF("#0#0.test.inner#0", 5)
		if l.GetTF("#5") != nil {
			t.Error("invalid TF does not return nil")
		}
		if l.GetTF("#0#0#0") != nil {
			t.Error("invalid TF does not return nil")
		}
		if l.GetTF("#0#0.nil") != nil {
			t.Error("invalid TF does not return nil")
		}
		if l.GetTF("#0#0.test#0") != nil {
			t.Error("invalid TF does not return nil")
		}
		if l.GetTF("#0.test") != nil {
			t.Error("invalid TF does not return nil")
		}
		if l.GetTF("#0#0.test.inner.nil") != nil {
			t.Error("invalid TF does not return nil")
		}
	})

}
