package anytype_test

import (
	"sync"
	"testing"

	"github.com/DanielSvub/anytype"
)

func TestObject(t *testing.T) {

	t.Run("basics", func(t *testing.T) {
		l := List()
		o := Object(
			"first", 1,
			"second", 2,
			"third", l,
		)
		if !o.Keys().Contains("second") {
			t.Error("key list should contain the key")
		}
		if !o.Values().Contains(2) {
			t.Error("value list should contain the value")
		}
		if !o.Contains(l) {
			t.Error("list should contain the list")
		}
		if !o.Contains(1) {
			t.Error("object should  not contain value 4")
		}
		if o.Contains(4) {
			t.Error("object should  not contain value 4")
		}
		if o.Count() != 3 {
			t.Error("object should have 3 fields")
		}
		if o.KeyOf(2) != "second" {
			t.Error("key for value 2 should be 'second'")
		}
		if o.Equals(Object()) {
			t.Error("object should not be equal to empty object")
		}
		if o.Pluck("first", "second").Count() != 2 {
			t.Error("plucked object should have 2 fields")
		}
		o.Unset("third")
		if !Object("first", 1).Merge(Object("second", 2)).Equals(o) {
			t.Error("merge does not work properly")
		}
		json := o.String()
		if json != `{"first":1,"second":2}` && json != `{"second":2,"first":1}` {
			t.Error("serialization does not work properly")
		}
		o.Clear()
		if !o.Empty() {
			t.Error("object should be empty")
		}
	})

	t.Run("types", func(t *testing.T) {
		o := Object(
			"object", Object("test", 0),
			"list", List(0),
			"string", "test",
			"bool", true,
			"int", 1,
			"float", 3.14,
			"nil", nil,
		)
		if o.TypeOf("object") != anytype.TypeObject {
			t.Error("field should be an object")
		}
		if o.TypeOf("list") != anytype.TypeList {
			t.Error("field should be a list")
		}
		if o.TypeOf("string") != anytype.TypeString {
			t.Error("field should be a string")
		}
		if o.TypeOf("bool") != anytype.TypeBool {
			t.Error("field should be a bool")
		}
		if o.TypeOf("int") != anytype.TypeInt {
			t.Error("field should be an int")
		}
		if o.TypeOf("float") != anytype.TypeFloat {
			t.Error("field should be a float")
		}
		if o.TypeOf("nil") != anytype.TypeNil {
			t.Error("field should be nil")
		}
		if o.TypeOf("undefined") != anytype.TypeUndefined {
			t.Error("field should be undefined")
		}
	})

	t.Run("equality", func(t *testing.T) {
		if Object("first", 1).Equals(Object("second", 2)) {
			t.Error("equality check does not work properly")
		}
	})

	t.Run("cloning", func(t *testing.T) {
		o := Object(
			"object", Object("test", 0),
			"list", List(0),
			"string", "test",
			"bool", true,
			"int", 1,
			"float", 3.14,
			"nil", nil,
		)
		if !o.Equals(o.Clone()) {
			t.Error("object should be equal to itself")
		}
	})

	t.Run("dictionaries", func(t *testing.T) {
		if !ObjectFrom(Object("str", "test", "int", 1).Dict()).Equals(Object("str", "test", "int", 1)) {
			t.Error("conversion to map[string]any does not work")
		}
		if !ObjectFrom(map[string]object{"obj": Object()}).Equals(Object("obj", Object())) {
			t.Error("conversion from map[string]object does not work")
		}
		if !ObjectFrom(map[string]list{"list": List()}).Equals(Object("list", List())) {
			t.Error("conversion from map[string]object does not work")
		}
		if !ObjectFrom(map[string]string{"str": "test"}).Equals(Object("str", "test")) {
			t.Error("conversion from map[string]object does not work")
		}
		if !ObjectFrom(map[string]bool{"bool": true}).Equals(Object("bool", true)) {
			t.Error("conversion from map[string]object does not work")
		}
		if !ObjectFrom(map[string]int{"int": 1}).Equals(Object("int", 1)) {
			t.Error("conversion from map[string]object does not work")
		}
		if !ObjectFrom(map[string]float64{"float": 3.14}).Equals(Object("float", 3.14)) {
			t.Error("conversion from map[string]object does not work")
		}
	})

	t.Run("getters", func(t *testing.T) {
		o := Object(
			"object", Object("test", 0),
			"list", List(0),
			"string", "test",
			"bool", true,
			"int", 1,
			"float", 3.14,
			"nil", nil,
		)
		if o.Get("int").(int) != 1 {
			t.Error("Int field should be 1")
		}
		if !o.GetObject("object").Equals(Object("test", 0)) {
			t.Error("cannot acquire an object")
		}
		if !o.GetList("list").Equals(List(0)) {
			t.Error("cannot acquire a list")
		}
		if o.GetString("string") != "test" {
			t.Error("cannot acquire a string")
		}
		if !o.GetBool("bool") {
			t.Error("cannot acquire a bool")
		}
		if o.GetInt("int") != 1 {
			t.Error("cannot acquire an int")
		}
		if o.GetFloat("float") != 3.14 {
			t.Error("cannot acquire a float")
		}
	})

	t.Run("foreaches", func(t *testing.T) {
		o := Object(
			"object", Object("test", 0),
			"list", List(0),
			"string", "test",
			"bool", true,
			"int", 1,
			"float", 3.14,
			"nil", nil,
		)
		t1 := Object()
		o.ForEach(func(key string, value any) { t1.Set(key, value) })
		if !t1.Equals(o) {
			t.Error("standard ForEach does not work properly")
		}
		t2 := Object()
		o.ForEachValue(func(value any) { t2.Set(o.KeyOf(value), value) })
		if !t2.Equals(o) {
			t.Error("ForEachValue does not work properly")
		}
		t3 := Object()
		o.ForEachObject(func(value object) { t3.Set(o.KeyOf(value), value) })
		if !t3.Equals(Object("object", Object("test", 0))) {
			t.Error("ForEachObject does not work properly")
		}
		t4 := Object()
		o.ForEachList(func(value list) { t4.Set(o.KeyOf(value), value) })
		if !t4.Equals(Object("list", List(0))) {
			t.Error("ForEachList does not work properly")
		}
		t5 := Object()
		o.ForEachString(func(value string) { t5.Set(o.KeyOf(value), value) })
		if !t5.Equals(Object("string", "test")) {
			t.Error("ForEachString does not work properly")
		}
		t6 := Object()
		o.ForEachBool(func(value bool) { t6.Set(o.KeyOf(value), value) })
		if !t6.Equals(Object("bool", true)) {
			t.Error("ForEachBool does not work properly")
		}
		t7 := Object()
		o.ForEachInt(func(value int) { t7.Set(o.KeyOf(value), value) })
		if !t7.Equals(Object("int", 1)) {
			t.Error("ForEachInt does not work properly")
		}
		t8 := Object()
		o.ForEachFloat(func(value float64) { t8.Set(o.KeyOf(value), value) })
		if !t8.Equals(Object("float", 3.14)) {
			t.Error("ForEachFloat does not work properly")
		}
		t9 := Object()
		var mutex sync.Mutex
		o.ForEachAsync(func(key string, value any) {
			mutex.Lock()
			t9.Set(key, value)
			mutex.Unlock()
		})
		if !t1.Equals(o) {
			t.Error("async ForEach does not work properly")
		}
	})

	t.Run("mappings", func(t *testing.T) {
		o := Object(
			"object", Object("test", 0),
			"list", List(0),
			"string", "test",
			"bool", true,
			"int", 1,
			"float", 3.14,
			"nil", nil,
		)
		if !o.Map(func(key string, value any) any { return value }).Equals(o) {
			t.Error("standard map does not work properly")
		}
		if !o.MapValues(func(value any) any { return value }).Equals(o) {
			t.Error("MapValues does not work properly")
		}
		if !o.MapObjects(func(value object) any { return value }).Equals(Object(
			"object", Object("test", 0),
		)) {
			t.Error("MapObjects does not work properly")
		}
		if !o.MapLists(func(value list) any { return value }).Equals(Object(
			"list", List(0),
		)) {
			t.Error("MapLists does not work properly")
		}
		if !o.MapStrings(func(value string) any { return value }).Equals(Object(
			"string", "test",
		)) {
			t.Error("MapStrings does not work properly")
		}
		if !o.MapInts(func(value int) any { return value }).Equals(Object(
			"int", 1,
		)) {
			t.Error("MapInts does not work properly")
		}
		if !o.MapFloats(func(value float64) any { return value }).Equals(Object(
			"float", 3.14,
		)) {
			t.Error("MapFloats does not work properly")
		}
		if !o.MapAsync(func(key string, value any) any { return value }).Equals(o) {
			t.Error("async map does not work properly")
		}
	})

}

func TestObjectPanics(t *testing.T) {

	catch := func(msg string) {
		if r := recover(); r == nil {
			t.Error(msg)
		}
	}

	t.Run("unsupportedMap", func(t *testing.T) {
		defer catch("creating object from string did not cause panic")
		ObjectFrom("unsupported")
	})

	t.Run("invalidSet", func(t *testing.T) {
		defer catch("invalid setting did not cause panic")
		Object("first", 1, "second")
	})

	t.Run("invalidKey", func(t *testing.T) {
		defer catch("invalid key did not cause panic")
		Object(1, 1)
	})

	t.Run("invalidGet", func(t *testing.T) {
		defer catch("getting non-existing field did not cause panic")
		Object().Get("test")
	})

	t.Run("invalidObjectGet", func(t *testing.T) {
		defer catch("getting non-object field as object did not cause panic")
		Object("test", false).GetObject("test")
	})

	t.Run("invalidListGet", func(t *testing.T) {
		defer catch("getting non-list field as list did not cause panic")
		Object("test", false).GetList("test")
	})

	t.Run("invalidStringGet", func(t *testing.T) {
		defer catch("getting non-string field as string did not cause panic")
		Object("test", false).GetString("test")
	})

	t.Run("invalidBoolGet", func(t *testing.T) {
		defer catch("getting non-bool field as bool did not cause panic")
		Object("test", 0).GetBool("test")
	})

	t.Run("invalidIntGet", func(t *testing.T) {
		defer catch("getting non-int field as int did not cause panic")
		Object("test", false).GetInt("test")
	})

	t.Run("invalidFloatGet", func(t *testing.T) {
		defer catch("getting non-float field as float did not cause panic")
		Object("test", false).GetFloat("test")
	})

	t.Run("invalidIndentation", func(t *testing.T) {
		defer catch("invalid indentation did not cause panic")
		Object().FormatString(-1)
	})

	t.Run("invalidValue", func(t *testing.T) {
		defer catch("non-existing value did not cause panic")
		Object().KeyOf("test")
	})

	t.Run("invalidGetTF", func(t *testing.T) {
		defer catch("getting invalid tree form did not cause panic")
		Object().GetTF("")
	})

	t.Run("invalidSetTF", func(t *testing.T) {
		defer catch("setting invalid tree form did not cause panic")
		Object().SetTF("", 0)
	})

}
