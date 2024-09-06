package anytype_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/DanielSvub/anytype"
)

func TestList(t *testing.T) {

	t.Run("basics", func(t *testing.T) {
		l := List(1, 2, 3)
		if l.TypeOf(0) != anytype.TypeInt {
			println(l.TypeOf(0))
			t.Error("element has a wrong type")
		}
		if !l.Insert(1, 4).Equals(List(1, 4, 2, 3)) {
			t.Error("element has not been inserted properly")
		}
		if !l.Replace(1, 5).Equals(List(1, 5, 2, 3)) {
			t.Error("element has not been replaced properly")
		}
		if !l.Delete(1).Equals(List(1, 2, 3)) {
			t.Error("element has not been deleted properly")
		}
		if !l.Add(4).Pop().Equals(List(1, 2, 3)) {
			t.Error("pop is not working properly")
		}
		if !l.Clone().Clear().Equals(List()) {
			t.Error("list has not been cleared properly")
		}
		if l.Count() != 3 {
			t.Error("list should have 3 elements")
		}
		if l.Empty() {
			t.Error("list should not be empty")
		}
		if !List(1, 2, 3, 4).Delete(2, 0, 3).Equals(List(2)) {
			t.Error("deleting multiple items does not work properly")
		}
		if !List().Empty() {
			t.Error("empty list should be empty")
		}
		if List(1, 2).Equals(List(1, 2, 3, 4)) {
			t.Error("lists should not be equal")
		}
		if List(1, 2).Equals(List(0, 1)) {
			t.Error("lists should not be equal")
		}
		if !List(1, 2).Concat(List(3, 4)).Equals(List(1, 2, 3, 4)) {
			t.Error("concatenation does not work properly")
		}
		if !l.Contains(1) {
			t.Error("list should contain element 1")
		}
		if l.Contains(4) {
			t.Error("list should not contain element 4")
		}
		if l.IndexOf(2) != 1 {
			t.Error("element 2 should be at index 1")
		}
		if l.IndexOf(4) != -1 {
			t.Error("IndexOf should return -1 if the element is not present")
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
			t.Error("element should be an object")
		}
		if l.TypeOf(1) != anytype.TypeList {
			t.Error("element should be a list")
		}
		if l.TypeOf(2) != anytype.TypeString {
			t.Error("element should be a string")
		}
		if l.TypeOf(3) != anytype.TypeBool {
			t.Error("element should be a bool")
		}
		if l.TypeOf(4) != anytype.TypeInt {
			t.Error("element should be an int")
		}
		if l.TypeOf(5) != anytype.TypeFloat {
			t.Error("element should be a float")
		}
		if l.TypeOf(6) != anytype.TypeNil {
			t.Error("element should be nil")
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
			t.Error("list should be equal to itself")
		}
	})

	t.Run("constructors", func(t *testing.T) {
		if !ListOf(1, 3).Equals(List(1, 1, 1)) {
			t.Error("listOf does not work properly")
		}
		if !ListFrom(make([]int, 3)).Equals(List(0, 0, 0)) {
			t.Error("listFrom does not work properly")
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
			t.Error("serialization does not work properly")
		}
	})

	t.Run("sorting", func(t *testing.T) {
		if !List(2, 4, 3, 5, 1).Sort().Equals(List(1, 2, 3, 4, 5)) {
			t.Error("ascending int sorting does not work properly")
		}
		if !List(2, 4, 3, 5, 1).Sort().Reverse().Equals(List(5, 4, 3, 2, 1)) {
			t.Error("descending int sorting does not work properly")
		}
		if !List(2.0, 4.0, 3.0, 5.0, 1.0).Sort().Equals(List(1.0, 2.0, 3.0, 4.0, 5.0)) {
			t.Error("ascending float sorting does not work properly")
		}
		if !List(2.0, 4.0, 3.0, 5.0, 1.0).Sort().Reverse().Equals(List(5.0, 4.0, 3.0, 2.0, 1.0)) {
			t.Error("descending float sorting does not work properly")
		}
		if !List("b", "c", "a").Sort().Equals(List("a", "b", "c")) {
			t.Error("ascending string sorting does not work properly")
		}
		if !List("b", "c", "a").Sort().Reverse().Equals(List("c", "b", "a")) {
			t.Error("descending string sorting does not work properly")
		}
	})

	t.Run("numeric", func(t *testing.T) {
		if List(2, 4, 3, 5, 1).IntMax() != 5 {
			t.Error("IntMax does not work")
		}
		if List().IntMax() != 0 {
			t.Error("IntMax does not work")
		}
		if List(2, 4, 3, 5, 1).IntMin() != 1 {
			t.Error("IntMin does not work")
		}
		if List().IntMin() != 0 {
			t.Error("IntMin does not work")
		}
		if List(2.0, 4, 3, 5.0, 1.0).Max() != 5 {
			t.Error("Max does not work")
		}
		if List(2, 4, 3, 5, 1).Max() != 5 {
			t.Error("Max does not work")
		}
		if List().Max() != 0 {
			t.Error("Max does not work")
		}
		if List(2.0, 4, 3, 5.0, 1.0).Min() != 1 {
			t.Error("Min does not work")
		}
		if List(2, 4, 3, 5, 1).Min() != 1 {
			t.Error("Min does not work")
		}
		if List().Min() != 0 {
			t.Error("Min does not work")
		}
		if List(1, 4, 5).IntSum() != 10 {
			t.Error("IntSum does not work")
		}
		if List(1.0, 4, 5.0).Sum() != 10 {
			t.Error("Sum does not work")
		}
		if List(1, 4, 5).IntProd() != 20 {
			t.Error("IntProd does not work")
		}
		if List(1.0, 4, 5.0).Prod() != 20 {
			t.Error("Prod does not work")
		}
		if List(0, 5, 5, 10).Avg() != 5 {
			t.Error("Avg does not work")
		}
	})

	t.Run("sublist", func(t *testing.T) {
		l := List(0, 1, 2, 3, 4)
		if !l.SubList(0, 0).Equals(l) {
			t.Error("SubList(0, 0) should return original list")
		}
		if !l.SubList(2, 4).Equals(List(2, 3)) {
			t.Error("SubList(2, 4) should return two elements")
		}
		if !l.SubList(0, -2).Equals(List(0, 1, 2)) {
			t.Error("SubList(0, -2) should cut last two elements")
		}
	})

	t.Run("asserts", func(t *testing.T) {
		l0 := List("test", 0)
		if l0.AllObjects() || l0.AllLists() || l0.AllStrings() || l0.AllBools() || l0.AllInts() || l0.AllFloats() || l0.AllNumeric() {
			t.Error("list should not be homogeneous")
		}
		if !List(Object()).AllObjects() {
			t.Error("list should be homogeneous")
		}
		if !List(List()).AllLists() {
			t.Error("list should be homogeneous")
		}
		if !List("test").AllStrings() {
			t.Error("list should be homogeneous")
		}
		if !List(true).AllBools() {
			t.Error("list should be homogeneous")
		}
		if !List(1).AllInts() {
			t.Error("list should be homogeneous")
		}
		if !List(3.14).AllFloats() {
			t.Error("list should be homogeneous")
		}
		if !List(1, 3.14).AllNumeric() {
			t.Error("list should be numeric")
		}
	})

	t.Run("slices", func(t *testing.T) {
		if !ListFrom(List(1, "test").Slice()).Equals(List(1, "test")) {
			t.Error("conversion to []any does not work")
		}
		if !ListFrom(List(Object()).ObjectSlice()).Equals(List(Object())) {
			t.Error("conversion to []object does not work")
		}
		if !ListFrom(List(List()).ListSlice()).Equals(List(List())) {
			t.Error("conversion to []list does not work")
		}
		if !ListFrom(List("test").StringSlice()).Equals(List("test")) {
			t.Error("conversion to []int does not work")
		}
		if !ListFrom(List(true).BoolSlice()).Equals(List(true)) {
			t.Error("conversion to []bool does not work")
		}
		if !ListFrom(List(1).IntSlice()).Equals(List(1)) {
			t.Error("conversion to []int does not work")
		}
		if !ListFrom(List(3.14).FloatSlice()).Equals(List(3.14)) {
			t.Error("conversion to []float64 does not work")
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
			t.Error("int element should be 1")
		}
		if !l.GetObject(0).Equals(Object("test", 0)) {
			t.Error("cannot acquire an object")
		}
		if !l.GetList(1).Equals(List(0)) {
			t.Error("cannot acquire a list")
		}
		if l.GetString(2) != "test" {
			t.Error("cannot acquire a string")
		}
		if !l.GetBool(3) {
			t.Error("cannot acquire a bool")
		}
		if l.GetInt(4) != 1 {
			t.Error("cannot acquire an int")
		}
		if l.GetFloat(5) != 3.14 {
			t.Error("cannot acquire a float")
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
			t.Error("standard ForEach does not work properly")
		}
		l2 := List()
		l.ForEachValue(func(value any) { l2.Add(value) })
		if !l2.Equals(l) {
			t.Error("ForEachValue does not work properly")
		}
		l3 := List()
		l.ForEachObject(func(value object) { l3.Add(value) })
		if !l3.Equals(List(Object("test", 0))) {
			t.Error("ForEachObject does not work properly")
		}
		l4 := List()
		l.ForEachList(func(value list) { l4.Add(value) })
		if !l4.Equals(List(List(0))) {
			t.Error("ForEachList does not work properly")
		}
		l5 := List()
		l.ForEachString(func(value string) { l5.Add(value) })
		if !l5.Equals(List("test")) {
			t.Error("ForEachString does not work properly")
		}
		l6 := List()
		l.ForEachBool(func(value bool) { l6.Add(value) })
		if !l6.Equals(List(true)) {
			t.Error("ForEachBool does not work properly")
		}
		l7 := List()
		l.ForEachInt(func(value int) { l7.Add(value) })
		if !l7.Equals(List(1)) {
			t.Error("ForEachInt does not work properly")
		}
		l8 := List()
		l.ForEachFloat(func(value float64) { l8.Add(value) })
		if !l8.Equals(List(3.14)) {
			t.Error("ForEachFloat does not work properly")
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
			t.Error("async ForEach does not work properly")
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
			t.Error("standard map does not work properly")
		}
		if !l.MapValues(func(value any) any { return value }).Equals(l) {
			t.Error("MapValues does not work properly")
		}
		if !l.MapObjects(func(value object) any { return value }).Equals(List(Object("test", 0))) {
			t.Error("MapObjects does not work properly")
		}
		if !l.MapLists(func(value list) any { return value }).Equals(List(List(0))) {
			t.Error("MapLists does not work properly")
		}
		if !l.MapStrings(func(value string) any { return value }).Equals(List("test")) {
			t.Error("MapStrings does not work properly")
		}
		if !l.MapInts(func(value int) any { return value }).Equals(List(1)) {
			t.Error("MapInts does not work properly")
		}
		if !l.MapFloats(func(value float64) any { return value }).Equals(List(3.14)) {
			t.Error("MapFloats does not work properly")
		}
		if !l.MapAsync(func(index int, value any) any { return value }).Equals(l) {
			t.Error("async map does not work properly")
		}
	})

	t.Run("reductions", func(t *testing.T) {
		l := List(1, 4, 5)
		if l.Reduce(0, func(sum any, value any) any {
			return sum.(int) + value.(int)
		}).(int) != 10 {
			t.Error("standard reduce does not work properly")
		}
		if List("a", "b", "c").ReduceStrings("Res: ", func(res string, value string) string {
			return res + value
		}) != "Res: abc" {
			t.Error("ReduceStrings does not work properly")
		}
		if l.ReduceInts(0, func(sum int, value int) int {
			return sum + value
		}) != 10 {
			t.Error("ReduceInts does not work properly")
		}
		if List(1.0, 4.0, 5.0).ReduceFloats(0, func(sum float64, value float64) float64 {
			return sum + value
		}) != 10.0 {
			t.Error("ReduceFloats does not work properly")
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
			t.Error("only 1 element should pass the standard filter")
		}
		if l.FilterObjects(func(value object) bool {
			return true
		}).Count() != 1 {
			t.Error("only 1 element should pass the object filter")
		}
		if l.FilterLists(func(value list) bool {
			return true
		}).Count() != 1 {
			t.Error("only 1 element should pass the list filter")
		}
		if l.FilterStrings(func(value string) bool {
			return true
		}).Count() != 1 {
			t.Error("only 1 element should pass the string filter")
		}
		if l.FilterInts(func(value int) bool {
			return true
		}).Count() != 1 {
			t.Error("only 1 element should pass the int filter")
		}
		if l.FilterFloats(func(value float64) bool {
			return true
		}).Count() != 1 {
			t.Error("only 1 element should pass the float filter")
		}
	})

}

func TestListPanics(t *testing.T) {

	catch := func(msg string) {
		if r := recover(); r == nil {
			t.Error(msg)
		}
	}

	t.Run("unsupportedSlice", func(t *testing.T) {
		defer catch("creating list from string did not cause panic")
		ListFrom("unsupported")
	})

	t.Run("invalidInsert", func(t *testing.T) {
		defer catch("invalid inserting did not cause panic")
		List().Insert(-1, 1)
	})

	t.Run("invalidGet", func(t *testing.T) {
		defer catch("getting non-existing element did not cause panic")
		List().Get(0)
	})

	t.Run("invalidObjectGet", func(t *testing.T) {
		defer catch("getting non-object element as object did not cause panic")
		List(false).GetObject(0)
	})

	t.Run("invalidListGet", func(t *testing.T) {
		defer catch("getting non-list element as list did not cause panic")
		List(false).GetList(0)
	})

	t.Run("invalidStringGet", func(t *testing.T) {
		defer catch("getting non-string element as string did not cause panic")
		List(false).GetString(0)
	})

	t.Run("invalidBoolGet", func(t *testing.T) {
		defer catch("getting non-bool element as bool did not cause panic")
		List(0).GetBool(0)
	})

	t.Run("invalidIntGet", func(t *testing.T) {
		defer catch("getting non-int element as int did not cause panic")
		List(false).GetInt(0)
	})

	t.Run("invalidFloatGet", func(t *testing.T) {
		defer catch("getting non-float element as float did not cause panic")
		List(false).GetFloat(0)
	})

	t.Run("invalidIndentation", func(t *testing.T) {
		defer catch("invalid indentation did not cause panic")
		List().FormatString(-1)
	})

	t.Run("invalidGetTF", func(t *testing.T) {
		defer catch("getting invalid tree form did not cause panic")
		List().GetTF("")
	})

	t.Run("invalidSetTF", func(t *testing.T) {
		defer catch("setting invalid tree form did not cause panic")
		List().SetTF("", 0)
	})

	t.Run("invalidReplace", func(t *testing.T) {
		defer catch("replacing non-existing item did not cause panic")
		List().Replace(0, 0)
	})

	t.Run("invalidDelete", func(t *testing.T) {
		defer catch("deleting non-existing item did not cause panic")
		List().Delete(0)
	})

	t.Run("invalidSubListEnd", func(t *testing.T) {
		defer catch("invalid ending sublist index did not cause panic")
		List().SubList(0, 1)
	})

	t.Run("invalidSubListRange", func(t *testing.T) {
		defer catch("starting sublist index higher than ending did not cause panic")
		List().SubList(1, 0)
	})

	t.Run("invalidSubListStart", func(t *testing.T) {
		defer catch("starting sublist index lower than zero did not cause panic")
		List().SubList(-1, 0)
	})

	t.Run("invalidSort", func(t *testing.T) {
		defer catch("sorting unsupported list did not cause panic")
		List(false).Sort()
	})

	t.Run("invalidGetTF", func(t *testing.T) {
		defer catch("getting invalid tree form did not cause panic")
		List().GetTF("")
	})

	t.Run("invalidGetTFIndex", func(t *testing.T) {
		defer catch("getting tree form with invalid index did not cause panic")
		List().GetTF("#test")
	})

	t.Run("invalidGetTFIndexObject", func(t *testing.T) {
		defer catch("getting tree form with invalid index and nested object did not cause panic")
		List().GetTF("#test.test")
	})

	t.Run("invalidGetTFIndexNested", func(t *testing.T) {
		defer catch("getting tree form with invalid index did not cause panic")
		List().GetTF("#test#0")
	})

	t.Run("invalidSetTF", func(t *testing.T) {
		defer catch("setting invalid tree form did not cause panic")
		List().SetTF("", 0)
	})

	t.Run("invalidSetTFIndex", func(t *testing.T) {
		defer catch("setting tree form with invalid index did not cause panic")
		List().SetTF("#test", 0)
	})

	t.Run("invalidSetTFIndexObject", func(t *testing.T) {
		defer catch("setting tree form with invalid index and nested object did not cause panic")
		List().SetTF("#test.test", 0)
	})

	t.Run("invalidSetTFIndexNested", func(t *testing.T) {
		defer catch("setting tree form with invalid index did not cause panic")
		List().SetTF("#test#0", 0)
	})

}
