package anytype_test

import (
	"os"
	"strconv"
	"sync"
	"testing"

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

func TestObject(t *testing.T) {
	t.Run("basics", func(t *testing.T) {
		l := List()
		o := Object(
			"first", 1,
			"second", 2,
			"third", l,
		)
		if !o.Keys().Contains("second") {
			t.Error("Key list should contain the key.")
		}
		if !o.Values().Contains(2) {
			t.Error("Value list should contain the value.")
		}
		if !o.Contains(l) {
			t.Error("Key list should contain the key.")
		}
		if o.Contains(4) {
			t.Error("Object should  not contain value 4.")
		}
		if o.Count() != 3 {
			t.Error("Object should have 3 fields.")
		}
		if o.KeyOf(2) != "second" {
			t.Error("Key for value 2 should be 'second'.")
		}
		if o.Equals(Object()) {
			t.Error("Object should not be equal to empty object.")
		}
		if o.Pluck("first", "second").Count() != 2 {
			t.Error("Plucked object should have 2 fields.")
		}
		o.Unset("third")
		if !Object("first", 1).Merge(Object("second", 2)).Equals(o) {
			t.Error("Merge does not work properly.")
		}
		json := o.String()
		if json != `{"first":1,"second":2}` && json != `{"second":2,"first":1}` {
			t.Error("Serialization does not work properly.")
		}
		o.Clear()
		if !o.Empty() {
			t.Error("Object should be empty.")
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
			t.Error("Field should be an object.")
		}
		if o.TypeOf("list") != anytype.TypeList {
			t.Error("Field should be a list.")
		}
		if o.TypeOf("string") != anytype.TypeString {
			t.Error("Field should be a string.")
		}
		if o.TypeOf("bool") != anytype.TypeBool {
			t.Error("Field should be a bool.")
		}
		if o.TypeOf("int") != anytype.TypeInt {
			t.Error("Field should be an int.")
		}
		if o.TypeOf("float") != anytype.TypeFloat {
			t.Error("Field should be a float.")
		}
		if o.TypeOf("nil") != anytype.TypeNil {
			t.Error("Field should be nil.")
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
			t.Error("Object should be equal to itself.")
		}
	})

	t.Run("dictionaries", func(t *testing.T) {
		if !ObjectFrom(Object("str", "test", "int", 1).Dict()).Equals(Object("str", "test", "int", 1)) {
			t.Error("Conversion to map[string]any does not work.")
		}
		if !ObjectFrom(map[string]object{"obj": Object()}).Equals(Object("obj", Object())) {
			t.Error("Conversion from map[string]object does not work.")
		}
		if !ObjectFrom(map[string]list{"list": List()}).Equals(Object("list", List())) {
			t.Error("Conversion from map[string]object does not work.")
		}
		if !ObjectFrom(map[string]string{"str": "test"}).Equals(Object("str", "test")) {
			t.Error("Conversion from map[string]object does not work.")
		}
		if !ObjectFrom(map[string]bool{"bool": true}).Equals(Object("bool", true)) {
			t.Error("Conversion from map[string]object does not work.")
		}
		if !ObjectFrom(map[string]int{"int": 1}).Equals(Object("int", 1)) {
			t.Error("Conversion from map[string]object does not work.")
		}
		if !ObjectFrom(map[string]float64{"float": 3.14}).Equals(Object("float", 3.14)) {
			t.Error("Conversion from map[string]object does not work.")
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
			t.Error("Int field should be 1.")
		}
		if !o.GetObject("object").Equals(Object("test", 0)) {
			t.Error("Cannot acquire an object.")
		}
		if !o.GetList("list").Equals(List(0)) {
			t.Error("Cannot acquire a list.")
		}
		if o.GetString("string") != "test" {
			t.Error("Cannot acquire a string.")
		}
		if !o.GetBool("bool") {
			t.Error("Cannot acquire a bool.")
		}
		if o.GetInt("int") != 1 {
			t.Error("Cannot acquire an int.")
		}
		if o.GetFloat("float") != 3.14 {
			t.Error("Cannot acquire a float.")
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
			t.Error("Standard ForEach does not work properly.")
		}
		t2 := Object()
		o.ForEachValue(func(value any) { t2.Set(o.KeyOf(value), value) })
		if !t2.Equals(o) {
			t.Error("ForEachValue does not work properly.")
		}
		t3 := Object()
		o.ForEachObject(func(value object) { t3.Set(o.KeyOf(value), value) })
		if !t3.Equals(Object("object", Object("test", 0))) {
			t.Error("ForEachObject does not work properly.")
		}
		t4 := Object()
		o.ForEachList(func(value list) { t4.Set(o.KeyOf(value), value) })
		if !t4.Equals(Object("list", List(0))) {
			t.Error("ForEachList does not work properly.")
		}
		t5 := Object()
		o.ForEachString(func(value string) { t5.Set(o.KeyOf(value), value) })
		if !t5.Equals(Object("string", "test")) {
			t.Error("ForEachString does not work properly.")
		}
		t6 := Object()
		o.ForEachBool(func(value bool) { t6.Set(o.KeyOf(value), value) })
		if !t6.Equals(Object("bool", true)) {
			t.Error("ForEachBool does not work properly.")
		}
		t7 := Object()
		o.ForEachInt(func(value int) { t7.Set(o.KeyOf(value), value) })
		if !t7.Equals(Object("int", 1)) {
			t.Error("ForEachInt does not work properly.")
		}
		t8 := Object()
		o.ForEachFloat(func(value float64) { t8.Set(o.KeyOf(value), value) })
		if !t8.Equals(Object("float", 3.14)) {
			t.Error("ForEachFloat does not work properly.")
		}
		t9 := Object()
		var mutex sync.Mutex
		o.ForEachAsync(func(key string, value any) {
			mutex.Lock()
			t9.Set(key, value)
			mutex.Unlock()
		})
		if !t1.Equals(o) {
			t.Error("Async ForEach does not work properly.")
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
			t.Error("Standard map does not work properly.")
		}
		if !o.MapValues(func(value any) any { return value }).Equals(o) {
			t.Error("MapValues does not work properly.")
		}
		if !o.MapObjects(func(value object) any { return value }).Equals(Object(
			"object", Object("test", 0),
		)) {
			t.Error("MapObjects does not work properly.")
		}
		if !o.MapLists(func(value list) any { return value }).Equals(Object(
			"list", List(0),
		)) {
			t.Error("MapLists does not work properly.")
		}
		if !o.MapStrings(func(value string) any { return value }).Equals(Object(
			"string", "test",
		)) {
			t.Error("MapStrings does not work properly.")
		}
		if !o.MapInts(func(value int) any { return value }).Equals(Object(
			"int", 1,
		)) {
			t.Error("MapInts does not work properly.")
		}
		if !o.MapFloats(func(value float64) any { return value }).Equals(Object(
			"float", 3.14,
		)) {
			t.Error("MapFloats does not work properly.")
		}
		if !o.MapAsync(func(key string, value any) any { return value }).Equals(o) {
			t.Error("Async map does not work properly.")
		}
	})

}

func TestList(t *testing.T) {

	t.Run("basics", func(t *testing.T) {
		l := List(1, 2, 3)
		if l.TypeOf(0) != anytype.TypeInt {
			t.Error("Element has a wrong type.")
		}
		if !l.Insert(1, 4).Equals(List(1, 4, 2, 3)) {
			t.Error("Element has not been inserted properly.")
		}
		if !l.Replace(1, 5).Equals(List(1, 5, 2, 3)) {
			t.Error("Element has not been replaced properly.")
		}
		if !l.Delete(1).Equals(List(1, 2, 3)) {
			t.Error("Element has not been deleted properly.")
		}
		if !l.Add(4).Pop().Equals(List(1, 2, 3)) {
			t.Error("Pop is not working properly.")
		}
		if !l.Clone().Clear().Equals(List()) {
			t.Error("List has not been cleared properly.")
		}
		if l.Count() != 3 {
			t.Error("List should have 3 elements.")
		}
		if l.Empty() {
			t.Error("List should not be empty.")
		}
		if !List().Empty() {
			t.Error("Empty list should be empty.")
		}
		if !List(1, 2).Concat(List(3, 4)).Equals(List(1, 2, 3, 4)) {
			t.Error("Concatenation does not work properly.")
		}
		if !l.Contains(1) {
			t.Error("List should contain element 1.")
		}
		if l.Contains(4) {
			t.Error("List should not contain element 4.")
		}
		if l.IndexOf(2) != 1 {
			t.Error("Element 2 should be at index 1.")
		}
		if l.IndexOf(4) != -1 {
			t.Error("IndexOf should return -1 if the element is not present.")
		}
	})

	t.Run("types", func(t *testing.T) {
		l := List(
			Object("test", 0),
			List(0),
			"test",
			true,
			1,
			3.14,
			nil,
		)
		if l.TypeOf(0) != anytype.TypeObject {
			t.Error("Element should be an object.")
		}
		if l.TypeOf(1) != anytype.TypeList {
			t.Error("Element should be a list.")
		}
		if l.TypeOf(2) != anytype.TypeString {
			t.Error("Element should be a string.")
		}
		if l.TypeOf(3) != anytype.TypeBool {
			t.Error("Element should be a bool.")
		}
		if l.TypeOf(4) != anytype.TypeInt {
			t.Error("Element should be an int.")
		}
		if l.TypeOf(5) != anytype.TypeFloat {
			t.Error("Element should be a float.")
		}
		if l.TypeOf(6) != anytype.TypeNil {
			t.Error("Element should be nil.")
		}
	})

	t.Run("cloning", func(t *testing.T) {
		l := List(
			Object("test", 0),
			List(0),
			"test",
			true,
			1,
			3.14,
			nil,
		)
		if !l.Equals(l.Clone()) {
			t.Error("List should be equal to itself.")
		}
	})

	t.Run("constructors", func(t *testing.T) {
		if !ListOf(1, 3).Equals(List(1, 1, 1)) {
			t.Error("ListOf does not work properly.")
		}
		if !ListFrom(make([]int, 3)).Equals(List(0, 0, 0)) {
			t.Error("ListFrom does not work properly.")
		}
	})

	t.Run("serialization", func(t *testing.T) {
		l := List(
			Object("test", 0),
			List(0),
			"test",
			true,
			1,
			3.14,
			nil,
		)
		json := `[{"test":0},[0],"test",true,1,3.14,null]`
		if l.String() != json {
			t.Error("Serialization does not work properly.")
		}
	})

	t.Run("sorting", func(t *testing.T) {
		if !List(2, 4, 3, 5, 1).Sort().Equals(List(1, 2, 3, 4, 5)) {
			t.Error("Ascending int sorting does not work properly.")
		}
		if !List(2, 4, 3, 5, 1).Sort().Reverse().Equals(List(5, 4, 3, 2, 1)) {
			t.Error("Descending int sorting does not work properly.")
		}
		if !List(2.0, 4.0, 3.0, 5.0, 1.0).Sort().Equals(List(1.0, 2.0, 3.0, 4.0, 5.0)) {
			t.Error("Ascending float sorting does not work properly.")
		}
		if !List(2.0, 4.0, 3.0, 5.0, 1.0).Sort().Reverse().Equals(List(5.0, 4.0, 3.0, 2.0, 1.0)) {
			t.Error("Descending float sorting does not work properly.")
		}
		if !List("b", "c", "a").Sort().Equals(List("a", "b", "c")) {
			t.Error("Ascending string sorting does not work properly.")
		}
		if !List("b", "c", "a").Sort().Reverse().Equals(List("c", "b", "a")) {
			t.Error("Descending string sorting does not work properly.")
		}
	})

	t.Run("numeric", func(t *testing.T) {
		if List(2, 4, 3, 5, 1).IntMax() != 5 {
			t.Error("IntMax does not work.")
		}
		if List(2, 4, 3, 5, 1).IntMin() != 1 {
			t.Error("IntMin does not work.")
		}
		if List(2.0, 4, 3, 5.0, 1.0).Max() != 5.0 {
			t.Error("Max does not work.")
		}
		if List(2.0, 4, 3, 5.0, 1.0).Min() != 1.0 {
			t.Error("Min does not work.")
		}
		if List(1, 4, 5).IntSum() != 10 {
			t.Error("IntSum does not work.")
		}
		if List(1.0, 4, 5.0).Sum() != 10.0 {
			t.Error("Sum does not work.")
		}
		if List(1, 4, 5).IntProd() != 20 {
			t.Error("IntProd does not work.")
		}
		if List(1.0, 4, 5.0).Prod() != 20.0 {
			t.Error("Prod does not work.")
		}
		if List(0, 5, 5, 10).Avg() != 5.0 {
			t.Error("Avg does not work.")
		}
	})

	t.Run("sublist", func(t *testing.T) {
		l := List(0, 1, 2, 3, 4)
		if !l.SubList(0, 0).Equals(l) {
			t.Error("SubList(0, 0) should return original list.")
		}
		if !l.SubList(2, 4).Equals(List(2, 3)) {
			t.Error("SubList(2, 4) should return two elements.")
		}
		if !l.SubList(0, -2).Equals(List(0, 1, 2)) {
			t.Error("SubList(0, -2) should cut last two elements.")
		}
	})

	t.Run("asserts", func(t *testing.T) {
		l0 := List("test", 0)
		if l0.AllObjects() || l0.AllLists() || l0.AllStrings() || l0.AllBools() || l0.AllInts() || l0.AllFloats() || l0.AllNumeric() {
			t.Error("List should not be homogeneous.")
		}
		if !List(Object()).AllObjects() {
			t.Error("List should be homogeneous.")
		}
		if !List(List()).AllLists() {
			t.Error("List should be homogeneous.")
		}
		if !List("test").AllStrings() {
			t.Error("List should be homogeneous.")
		}
		if !List(true).AllBools() {
			t.Error("List should be homogeneous.")
		}
		if !List(1).AllInts() {
			t.Error("List should be homogeneous.")
		}
		if !List(3.14).AllFloats() {
			t.Error("List should be homogeneous.")
		}
		if !List(1, 3.14).AllNumeric() {
			t.Error("List should be numeric.")
		}
	})

	t.Run("slices", func(t *testing.T) {
		if !ListFrom(List(1, "test").Slice()).Equals(List(1, "test")) {
			t.Error("Conversion to []any does not work.")
		}
		if !ListFrom(List(Object()).ObjectSlice()).Equals(List(Object())) {
			t.Error("Conversion to []object does not work.")
		}
		if !ListFrom(List(List()).ListSlice()).Equals(List(List())) {
			t.Error("Conversion to []list does not work.")
		}
		if !ListFrom(List("test").StringSlice()).Equals(List("test")) {
			t.Error("Conversion to []int does not work.")
		}
		if !ListFrom(List(true).BoolSlice()).Equals(List(true)) {
			t.Error("Conversion to []bool does not work.")
		}
		if !ListFrom(List(1).IntSlice()).Equals(List(1)) {
			t.Error("Conversion to []int does not work.")
		}
		if !ListFrom(List(3.14).FloatSlice()).Equals(List(3.14)) {
			t.Error("Conversion to []float64 does not work.")
		}
	})

	t.Run("getters", func(t *testing.T) {
		l := List(
			Object("test", 0),
			List(0),
			"test",
			true,
			1,
			3.14,
			nil,
		)
		if l.Get(4).(int) != 1 {
			t.Error("Int element should be 1.")
		}
		if !l.GetObject(0).Equals(Object("test", 0)) {
			t.Error("Cannot acquire an object.")
		}
		if !l.GetList(1).Equals(List(0)) {
			t.Error("Cannot acquire a list.")
		}
		if l.GetString(2) != "test" {
			t.Error("Cannot acquire a string.")
		}
		if !l.GetBool(3) {
			t.Error("Cannot acquire a bool.")
		}
		if l.GetInt(4) != 1 {
			t.Error("Cannot acquire an int.")
		}
		if l.GetFloat(5) != 3.14 {
			t.Error("Cannot acquire a float.")
		}
	})

	t.Run("foreaches", func(t *testing.T) {
		l := List(
			Object("test", 0),
			List(0),
			"test",
			true,
			1,
			3.14,
			nil,
		)
		l1 := List()
		l.ForEach(func(index int, value any) { l1.Insert(index, value) })
		if !l1.Equals(l) {
			t.Error("Standard ForEach does not work properly.")
		}
		l2 := List()
		l.ForEachValue(func(value any) { l2.Add(value) })
		if !l2.Equals(l) {
			t.Error("ForEachValue does not work properly.")
		}
		l3 := List()
		l.ForEachObject(func(value object) { l3.Add(value) })
		if !l3.Equals(List(Object("test", 0))) {
			t.Error("ForEachObject does not work properly.")
		}
		l4 := List()
		l.ForEachList(func(value list) { l4.Add(value) })
		if !l4.Equals(List(List(0))) {
			t.Error("ForEachList does not work properly.")
		}
		l5 := List()
		l.ForEachString(func(value string) { l5.Add(value) })
		if !l5.Equals(List("test")) {
			t.Error("ForEachString does not work properly.")
		}
		l6 := List()
		l.ForEachBool(func(value bool) { l6.Add(value) })
		if !l6.Equals(List(true)) {
			t.Error("ForEachBool does not work properly.")
		}
		l7 := List()
		l.ForEachInt(func(value int) { l7.Add(value) })
		if !l7.Equals(List(1)) {
			t.Error("ForEachInt does not work properly.")
		}
		l8 := List()
		l.ForEachFloat(func(value float64) { l8.Add(value) })
		if !l8.Equals(List(3.14)) {
			t.Error("ForEachFloat does not work properly.")
		}
		testObj := Object(
			"0", Object("test", 0),
			"1", List(0),
			"2", "test",
			"3", true,
			"4", 1,
			"5", 3.14,
			"6", nil,
		)
		o := Object()
		var mutex sync.Mutex
		l.ForEachAsync(func(index int, value any) {
			mutex.Lock()
			o.Set(strconv.Itoa(index), value)
			mutex.Unlock()
		})
		if !o.Equals(testObj) {
			t.Error("Async ForEach does not work properly.")
		}
	})

	t.Run("mappings", func(t *testing.T) {
		l := List(
			Object("test", 0),
			List(0),
			"test",
			true,
			1,
			3.14,
			nil,
		)
		if !l.Map(func(index int, value any) any { return value }).Equals(l) {
			t.Error("Standard map does not work properly.")
		}
		if !l.MapValues(func(value any) any { return value }).Equals(l) {
			t.Error("MapValues does not work properly.")
		}
		if !l.MapObjects(func(value object) any { return value }).Equals(List(Object("test", 0))) {
			t.Error("MapObjects does not work properly.")
		}
		if !l.MapLists(func(value list) any { return value }).Equals(List(List(0))) {
			t.Error("MapLists does not work properly.")
		}
		if !l.MapStrings(func(value string) any { return value }).Equals(List("test")) {
			t.Error("MapStrings does not work properly.")
		}
		if !l.MapInts(func(value int) any { return value }).Equals(List(1)) {
			t.Error("MapInts does not work properly.")
		}
		if !l.MapFloats(func(value float64) any { return value }).Equals(List(3.14)) {
			t.Error("MapFloats does not work properly.")
		}
		if !l.MapAsync(func(index int, value any) any { return value }).Equals(l) {
			t.Error("Async map does not work properly.")
		}
	})

	t.Run("reductions", func(t *testing.T) {
		l := List(1, 4, 5)
		if l.Reduce(0, func(sum any, value any) any {
			return sum.(int) + value.(int)
		}).(int) != 10 {
			t.Error("Standard reduce does not work properly.")
		}
		if List("a", "b", "c").ReduceStrings("Res: ", func(res string, value string) string {
			return res + value
		}) != "Res: abc" {
			t.Error("ReduceStrings does not work properly.")
		}
		if l.ReduceInts(0, func(sum int, value int) int {
			return sum + value
		}) != 10 {
			t.Error("ReduceInts does not work properly.")
		}
		if List(1.0, 4.0, 5.0).ReduceFloats(0, func(sum float64, value float64) float64 {
			return sum + value
		}) != 10.0 {
			t.Error("ReduceFloats does not work properly.")
		}
	})

	t.Run("filters", func(t *testing.T) {
		l := List(
			Object("test", 0),
			List(0),
			"test",
			true,
			1,
			3.14,
			nil,
		)
		if l.Filter(func(value any) bool {
			_, ok := value.(int)
			return ok
		}).Count() != 1 {
			t.Error("Only 1 element should pass the standard filter.")
		}
		if l.FilterObjects(func(value object) bool {
			return true
		}).Count() != 1 {
			t.Error("Only 1 element should pass the object filter.")
		}
		if l.FilterLists(func(value list) bool {
			return true
		}).Count() != 1 {
			t.Error("Only 1 element should pass the list filter.")
		}
		if l.FilterStrings(func(value string) bool {
			return true
		}).Count() != 1 {
			t.Error("Only 1 element should pass the string filter.")
		}
		if l.FilterInts(func(value int) bool {
			return true
		}).Count() != 1 {
			t.Error("Only 1 element should pass the int filter.")
		}
		if l.FilterFloats(func(value float64) bool {
			return true
		}).Count() != 1 {
			t.Error("Only 1 element should pass the float filter.")
		}
	})

}

func TestParsing(t *testing.T) {
	t.Run("object", func(t *testing.T) {
		o := Object(
			"object", Object(
				"\"tes\u0074\n\"", 0,
				"innerList", List(0, 1, Object(), "2"),
			),
			"list", List(nil, false, 42, 1.6e-8, "\"Žř@./'á?\u0041\n\\\"", Object("\"1\"", "\"1\""), List(0)),
			"string", "test",
			"bool", true,
			"int", 1,
			"float", 3.14,
			"nil", nil,
		)
		parsed, err := anytype.ParseObject(o.String())
		if err != nil {
			t.Error("JSON parser failed.")
		}
		if !parsed.Equals(o) {
			t.Error("JSON parsing or exporting of object does not work properly.")
		}
	})

	t.Run("list", func(t *testing.T) {
		l := List(nil, false, 42, 1.6, "\"Žř@./'á?\u0041\n\\\"", List(0, "1"), Object(
			"bool", true,
			"int", 1,
			"float", 3.14,
			"innerList", List(0, 1, Object(), "2"),
			"nil", nil,
			"string", "test",
		))
		parsed, err := anytype.ParseList(l.String())
		if err != nil {
			t.Error("JSON parser failed.")
		}
		if !parsed.Equals(l) {
			t.Error("JSON parsing or exporting of list does not work properly.")
		}
	})

	t.Run("indent", func(t *testing.T) {
		l := List(1, 2)
		if l.FormatString(2) != "[\n  1,\n  2\n]" {
			t.Error("JSON formatted export does not work properly.")
		}
		o := Object(
			"key", "value",
		)
		if o.FormatString(2) != "{\n  \"key\": \"value\"\n}" {
			t.Error("JSON formatted export does not work properly.")
		}
	})

	t.Run("errors", func(t *testing.T) {
		if _, err := anytype.ParseObject(""); err == nil {
			t.Error("Parser did not return expected error.")
		}
		if _, err := anytype.ParseObject("{\"test\":1"); err == nil {
			t.Error("Parser did not return expected error.")
		}
		if _, err := anytype.ParseObject("{\"test\"1"); err == nil {
			t.Error("Parser did not return expected error.")
		}
		if _, err := anytype.ParseObject("{1:2}"); err == nil {
			t.Error("Parser did not return expected error.")
		}
		if _, err := anytype.ParseObject("{\"test\":[]2}"); err == nil {
			t.Error("Parser did not return expected error.")
		}
		if _, err := anytype.ParseList(""); err == nil {
			t.Error("Parser did not return expected error.")
		}
		if _, err := anytype.ParseList("[1,2"); err == nil {
			t.Error("Parser did not return expected error.")
		}
		if _, err := anytype.ParseList("[test]"); err == nil {
			t.Error("Parser did not return expected error.")
		}
	})

	t.Run("repairs", func(t *testing.T) {
		if _, err := anytype.ParseObject("{\"first\":{}\"second\":2}"); err != nil {
			t.Error("Parser did not repair broken JSON.")
		}
		if _, err := anytype.ParseObject("{\"first\":1 2}"); err != nil {
			t.Error("Parser did not repair broken JSON.")
		}
		if _, err := anytype.ParseList("[nu ll,[]2]"); err != nil {
			t.Error("Parser did not repair broken JSON.")
		}
	})

	t.Run("file", func(t *testing.T) {
		if _, err := anytype.ParseFile("test.json"); err == nil {
			t.Error("Opened a file which should not exist.")
		}
		if err := os.WriteFile("test.json", []byte("{\"first\":[],\"second\":2}"), 0644); err != nil {
			t.Fatal("Unable to create the JSON file.")
		}
		if _, err := anytype.ParseFile("test.json"); err != nil {
			t.Error("Cannot parse JSON from file.")
		}
		if err := os.Remove("test.json"); err != nil {
			t.Fatal("Unable to delete the JSON file.")
		}
	})
}

func TestTF(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		l := List()
		l.SetTF("#0#0", 1)
		l.SetTF("#0#1", 2).SetTF("#0#2.test.inner#0", 1)
		l.SetTF("#0#2.test.inner#1.target#0", 5)
		result := l.GetTF("#0#2.test.inner#1.target#0").(int)
		if result != 5 {
			t.Error("Tree form handling does not work properly.")
		}
		if l.GetTF("#0#2.test") == nil {
			t.Error("Valid TF returns nil.")
		}
	})

	t.Run("invalid", func(t *testing.T) {
		l := List().SetTF("#0#0.test.inner#0", 5)
		if l.GetTF("#5") != nil {
			t.Error("Invalid TF does not return nil.")
		}
		if l.GetTF("#0#0#0") != nil {
			t.Error("Invalid TF does not return nil.")
		}
		if l.GetTF("#0#0.nil") != nil {
			t.Error("Invalid TF does not return nil.")
		}
		if l.GetTF("#0#0.test#0") != nil {
			t.Error("Invalid TF does not return nil.")
		}
		if l.GetTF("#0.test") != nil {
			t.Error("Invalid TF does not return nil.")
		}
		if l.GetTF("#0#0.test.inner.nil") != nil {
			t.Error("Invalid TF does not return nil.")
		}
	})
}
