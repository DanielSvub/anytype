/*
AnyType Library for Go
List (array) implementation
*/

package anytype

import (
	"bytes"
	"encoding/json"
	"math"
	"math/bits"
	"sort"
	"strconv"
	"strings"
	"sync"
)

/*
Slice list, a reference type. Contains a slice of elements.

Implements:
  - field,
  - List.
*/
type list struct {
	val []field
	ptr List
}

/*
List constructor.
Creates a new list.

Parameters:
  - values... - any amount of initial elements.

Returns:
  - pointer to the created list.
*/
func NewList(values ...any) List {
	ego := &list{val: make([]field, 0)}
	ego.ptr = ego
	ego.Add(values...)
	return ego
}

/*
List constructor.
Creates a new list of n repeated values.

Parameters:
  - value - value to repeat,
  - count - number of repetitions.

Returns:
  - pointer to the created list.
*/
func NewListOf(value any, count int) List {
	ego := &list{val: make([]field, count)}
	ego.ptr = ego
	elem := parseVal(value)
	for i := 0; i < ego.Count(); i++ {
		ego.val[i] = elem
	}
	return ego
}

/*
List constructor.
Converts a slice of supported types to a list.

Parameters:
  - slice - original slice.

Returns:
  - created list.
*/
func NewListFrom(slice any) List {
	ego := NewList()
	switch s := slice.(type) {
	case []any:
		for _, item := range s {
			ego.Add(item)
		}
	case []Object:
		for _, item := range s {
			ego.Add(item)
		}
	case []List:
		for _, item := range s {
			ego.Add(item)
		}
	case []string:
		for _, item := range s {
			ego.Add(item)
		}
	case []bool:
		for _, item := range s {
			ego.Add(item)
		}
	case []int:
		for _, item := range s {
			ego.Add(item)
		}
	case []float64:
		for _, item := range s {
			ego.Add(item)
		}
	default:
		panic("Unknown slice type.")
	}
	return ego
}

/*
Asserts that the list is initialized.
*/
func (ego *list) assert() {
	if ego == nil || ego.val == nil {
		panic("List is not initialized.")
	}
}

/*
Defined in the field interface.
Acquires the value of the field, in this case a reference to the whole struct (List is reference type).

Returns:
  - value of the field.
*/
func (ego *list) getVal() any {
	return ego.ptr
}

/*
Defined in the field interface.
Creates a deep copy of the field, in this case a new list with identical elements.
Can be called recursively.

Returns:
  - deep copy of the field.
*/
func (ego *list) copy() any {
	list := NewList()
	for _, value := range ego.val {
		list.Add(value.copy())
	}
	return list
}

/*
Defined in the field interface.
Serializes the field into the JSON format, in this case prints all elements of the list.
Can be called recursively.

Returns:
  - string representing serialized field.
*/
func (ego *list) serialize() string {
	result := "["
	for i, value := range ego.val {
		result += value.serialize()
		if i+1 < len(ego.val) {
			result += ","
		}
	}
	result += "]"
	return result
}

/*
Defined in the field interface.
Checks if the content of the field is equal to the given field.
Can be called recursively.

Returns:
  - true if the fields are equal, false otherwise.
*/
func (ego *list) isEqual(another any) bool {
	list, ok := another.(*list)
	if !ok || ego.Count() != list.Count() {
		return false
	}
	for i := range ego.val {
		if !ego.val[i].isEqual(list.val[i]) {
			return false
		}
	}
	return true
}

func (ego *list) Init(ptr List) {
	ego.ptr = ptr
}

func (ego *list) Ego() List {
	return ego.ptr
}

func (ego *list) Add(values ...any) List {
	ego.assert()
	for _, val := range values {
		ego.val = append(ego.val, parseVal(val))
	}
	return ego.ptr
}

func (ego *list) Insert(index int, value any) List {
	ego.assert()
	if index < 0 || index > ego.Count() {
		panic("Index " + strconv.Itoa(index) + " out of range.")
	}
	if index == ego.Count() {
		return ego.ptr.Add(value)
	}
	ego.val = append(ego.val[:index+1], ego.val[index:]...)
	ego.val[index] = parseVal(value)
	return ego.ptr
}

func (ego *list) Replace(index int, value any) List {
	ego.assert()
	if index < 0 || index > ego.Count() {
		panic("Index " + strconv.Itoa(index) + " out of range.")
	}
	ego.val[index] = parseVal(value)
	return ego.ptr
}

func (ego *list) Delete(indexes ...int) List {
	ego.assert()
	if len(indexes) > 1 {
		sort.Ints(indexes)
	}
	for i := len(indexes) - 1; i >= 0; i-- {
		index := indexes[i]
		if len(ego.val) <= index || index < 0 {
			panic("Index " + strconv.Itoa(index) + " out of range.")
		}
		ego.val = append(ego.val[:index], ego.val[index+1:]...)
	}
	return ego.ptr
}

func (ego *list) Pop() List {
	return ego.ptr.Delete(ego.Count() - 1)
}

func (ego *list) Clear() List {
	ego.assert()
	ego.val = make([]field, 0)
	return ego.ptr
}

func (ego *list) Get(index int) any {
	ego.assert()
	if len(ego.val) <= index || index < 0 {
		panic("Index " + strconv.Itoa(index) + " out of range.")
	}
	obj := ego.val[index]
	switch obj.(type) {
	case Object, List:
		return obj
	default:
		return obj.getVal()
	}
}

func (ego *list) GetObject(index int) Object {
	o, ok := ego.Get(index).(Object)
	if !ok {
		panic("Item is not an object.")
	}
	return o
}

func (ego *list) GetList(index int) List {
	o, ok := ego.Get(index).(List)
	if !ok {
		panic("Item is not a list.")
	}
	return o
}

func (ego *list) GetString(index int) string {
	o, ok := ego.Get(index).(string)
	if !ok {
		panic("Item is not a string.")
	}
	return o
}

func (ego *list) GetBool(index int) bool {
	o, ok := ego.Get(index).(bool)
	if !ok {
		panic("Item is not a bool.")
	}
	return o
}

func (ego *list) GetInt(index int) int {
	o, ok := ego.Get(index).(int)
	if !ok {
		panic("Item is not an int.")
	}
	return o
}

func (ego *list) GetFloat(index int) float64 {
	o, ok := ego.Get(index).(float64)
	if !ok {
		panic("Item is not a float.")
	}
	return o
}

func (ego *list) TypeOf(index int) Type {
	ego.assert()
	switch ego.val[index].(type) {
	case *atString:
		return TypeString
	case *atInt:
		return TypeInt
	case *atBool:
		return TypeBool
	case *atFloat:
		return TypeFloat
	case Object:
		return TypeObject
	case List:
		return TypeList
	case *atNil:
		return TypeNil
	default:
		panic("Unknown element type.")
	}
}

func (ego *list) String() string {
	ego.assert()
	return ego.ptr.serialize()
}

func (ego *list) FormatString(indent int) string {
	if indent < 0 || indent > 10 {
		panic("Invalid indentation.")
	}
	buffer := new(bytes.Buffer)
	json.Indent(buffer, []byte(ego.String()), "", strings.Repeat(" ", indent))
	return buffer.String()
}

func (ego *list) Slice() []any {
	ego.assert()
	slice := make([]any, 0)
	for _, item := range ego.val {
		slice = append(slice, item.getVal())
	}
	return slice
}

func (ego *list) ObjectSlice() []Object {
	ego.assert()
	if !ego.AllObjects() {
		panic("All elements have to be objects.")
	}
	slice := make([]Object, 0)
	for _, item := range ego.val {
		slice = append(slice, item.(Object))
	}
	return slice
}

func (ego *list) ListSlice() []List {
	ego.assert()
	if !ego.AllLists() {
		panic("All elements have to be lists.")
	}
	slice := make([]List, 0)
	for _, item := range ego.val {
		slice = append(slice, item.(List))
	}
	return slice
}

func (ego *list) StringSlice() []string {
	ego.assert()
	if !ego.AllStrings() {
		panic("All elements have to be strings.")
	}
	slice := make([]string, 0)
	for _, item := range ego.val {
		slice = append(slice, item.getVal().(string))
	}
	return slice
}

func (ego *list) BoolSlice() []bool {
	ego.assert()
	if !ego.AllBools() {
		panic("All elements have to be bools.")
	}
	slice := make([]bool, 0)
	for _, item := range ego.val {
		slice = append(slice, item.getVal().(bool))
	}
	return slice
}

func (ego *list) IntSlice() []int {
	ego.assert()
	if !ego.AllInts() {
		panic("All elements have to be ints.")
	}
	slice := make([]int, 0)
	for _, item := range ego.val {
		slice = append(slice, item.getVal().(int))
	}
	return slice
}

func (ego *list) FloatSlice() []float64 {
	ego.assert()
	if !ego.AllFloats() {
		panic("All elements have to be floats.")
	}
	slice := make([]float64, 0)
	for _, item := range ego.val {
		slice = append(slice, item.getVal().(float64))
	}
	return slice
}

func (ego *list) Clone() List {
	ego.assert()
	return ego.ptr.copy().(*list)
}

func (ego *list) Count() int {
	ego.assert()
	return len(ego.val)
}

func (ego *list) Empty() bool {
	return ego.ptr.Count() == 0
}

func (ego *list) Equals(another List) bool {
	ego.assert()
	return ego.ptr.isEqual(another)
}

func (ego *list) Concat(another List) List {
	ego.assert()
	newList := &list{val: append(ego.val, another.getVal().(*list).val...)}
	newList.Init(newList)
	return newList
}

func (ego *list) SubList(start int, end int) List {
	ego.assert()
	if end <= 0 {
		end = ego.Count() + end
	}
	if start > end {
		panic("Starting index higher than ending index.")
	}
	if len(ego.val) < end || start < 0 {
		panic("Index out of range.")
	}
	list := &list{val: make([]field, end-start)}
	list.Init(list)
	copy(list.val, ego.val[start:end])
	return list
}

func (ego *list) Contains(elem any) bool {
	ego.assert()
	for _, item := range ego.val {
		switch item.(type) {
		case Object, List:
			if item == elem {
				return true
			}
		default:
			if item.getVal() == elem {
				return true
			}
		}
	}
	return false
}

func (ego *list) IndexOf(elem any) int {
	ego.assert()
	for i, item := range ego.val {
		switch item.(type) {
		case Object, List:
			if item == elem {
				return i
			}
		default:
			if item.getVal() == elem {
				return i
			}
		}
	}
	return -1
}

func (ego *list) Sort() List {
	ego.assert()
	switch ego.val[0].(type) {
	case *atString:
		slice := ego.StringSlice()
		sort.Strings(slice)
		ego.val = NewListFrom(slice).(*list).val
	case *atInt:
		slice := ego.IntSlice()
		sort.Ints(slice)
		ego.val = NewListFrom(slice).(*list).val
	case *atFloat:
		slice := ego.FloatSlice()
		sort.Float64s(slice)
		ego.val = NewListFrom(slice).(*list).val
	default:
		panic("List has to be homogeneous with all its elements numeric or strings.")
	}
	return ego.ptr
}

func (ego *list) Reverse() List {
	ego.assert()
	for i := ego.Count()/2 - 1; i >= 0; i-- {
		opp := ego.Count() - 1 - i
		ego.val[i], ego.val[opp] = ego.val[opp], ego.val[i]
	}
	return ego.ptr
}

func (ego *list) AllObjects() bool {
	ego.assert()
	for _, item := range ego.val {
		_, ok := item.(Object)
		if !ok {
			return false
		}
	}
	return true
}

func (ego *list) AllLists() bool {
	ego.assert()
	for _, item := range ego.val {
		_, ok := item.(List)
		if !ok {
			return false
		}
	}
	return true
}

func (ego *list) AllStrings() bool {
	ego.assert()
	for _, item := range ego.val {
		_, ok := item.(*atString)
		if !ok {
			return false
		}
	}
	return true
}

func (ego *list) AllBools() bool {
	ego.assert()
	for _, item := range ego.val {
		_, ok := item.(*atBool)
		if !ok {
			return false
		}
	}
	return true
}

func (ego *list) AllInts() bool {
	ego.assert()
	for _, item := range ego.val {
		_, ok := item.(*atInt)
		if !ok {
			return false
		}
	}
	return true
}

func (ego *list) AllFloats() bool {
	ego.assert()
	for _, item := range ego.val {
		_, ok := item.(*atFloat)
		if !ok {
			return false
		}
	}
	return true
}

func (ego *list) AllNumeric() bool {
	ego.assert()
	for _, item := range ego.val {
		_, ok := item.(*atInt)
		if !ok {
			_, ok := item.(*atFloat)
			if !ok {
				return false
			}
		}
	}
	return true
}

func (ego *list) ForEach(function func(int, any)) List {
	ego.assert()
	for i, item := range ego.val {
		function(i, item.getVal())
	}
	return ego.ptr
}

func (ego *list) ForEachValue(function func(any)) List {
	ego.assert()
	for _, item := range ego.val {
		function(item.getVal())
	}
	return ego.ptr
}

func (ego *list) ForEachObject(function func(Object)) List {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.(Object)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

func (ego *list) ForEachList(function func(List)) List {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.(List)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

func (ego *list) ForEachString(function func(string)) List {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(string)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

func (ego *list) ForEachBool(function func(bool)) List {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(bool)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

func (ego *list) ForEachInt(function func(int)) List {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(int)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

func (ego *list) ForEachFloat(function func(float64)) List {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(float64)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

func (ego *list) Map(function func(int, any) any) List {
	ego.assert()
	result := NewList()
	for i, item := range ego.val {
		result.Add(function(i, item.getVal()))
	}
	return result
}

func (ego *list) MapValues(function func(any) any) List {
	ego.assert()
	result := NewList()
	for _, item := range ego.val {
		result.Add(function(item.getVal()))
	}
	return result
}

func (ego *list) MapObjects(function func(Object) any) List {
	ego.assert()
	result := NewList()
	for _, item := range ego.val {
		val, ok := item.(Object)
		if ok {
			result.Add(function(val))
		}
	}
	return result
}

func (ego *list) MapLists(function func(List) any) List {
	ego.assert()
	result := NewList()
	for _, item := range ego.val {
		val, ok := item.(List)
		if ok {
			result.Add(function(val))
		}
	}
	return result
}

func (ego *list) MapStrings(function func(string) any) List {
	ego.assert()
	result := NewList()
	for _, item := range ego.val {
		val, ok := item.getVal().(string)
		if ok {
			result.Add(function(val))
		}
	}
	return result
}

func (ego *list) MapInts(function func(int) any) List {
	ego.assert()
	result := NewList()
	for _, item := range ego.val {
		val, ok := item.getVal().(int)
		if ok {
			result.Add(function(val))
		}
	}
	return result
}

func (ego *list) MapFloats(function func(float64) any) List {
	ego.assert()
	result := NewList()
	for _, item := range ego.val {
		val, ok := item.getVal().(float64)
		if ok {
			result.Add(function(val))
		}
	}
	return result
}

func (ego *list) Reduce(initial any, function func(any, any) any) any {
	ego.assert()
	result := initial
	for _, item := range ego.val {
		result = function(result, item.getVal())
	}
	return result
}

func (ego *list) ReduceStrings(initial string, function func(string, string) string) string {
	ego.assert()
	result := initial
	for _, item := range ego.val {
		val, ok := item.getVal().(string)
		if ok {
			result = function(result, val)
		}
	}
	return result
}

func (ego *list) ReduceInts(initial int, function func(int, int) int) int {
	ego.assert()
	result := initial
	for _, item := range ego.val {
		val, ok := item.getVal().(int)
		if ok {
			result = function(result, val)
		}
	}
	return result
}

func (ego *list) ReduceFloats(initial float64, function func(float64, float64) float64) float64 {
	ego.assert()
	result := initial
	for _, item := range ego.val {
		val, ok := item.getVal().(float64)
		if ok {
			result = function(result, val)
		}
	}
	return result
}

func (ego *list) Filter(function func(any) bool) List {
	ego.assert()
	result := NewList()
	for _, item := range ego.val {
		if function(item.getVal()) {
			result.Add(item.getVal())
		}
	}
	return result
}

func (ego *list) FilterObjects(function func(Object) bool) List {
	ego.assert()
	result := NewList()
	for _, item := range ego.val {
		val, ok := item.(Object)
		if ok && function(val) {
			result.Add(val)
		}
	}
	return result
}

func (ego *list) FilterLists(function func(List) bool) List {
	ego.assert()
	result := NewList()
	for _, item := range ego.val {
		val, ok := item.(List)
		if ok && function(val) {
			result.Add(val)
		}
	}
	return result
}

func (ego *list) FilterStrings(function func(string) bool) List {
	ego.assert()
	result := NewList()
	for _, item := range ego.val {
		val, ok := item.getVal().(string)
		if ok && function(val) {
			result.Add(val)
		}
	}
	return result
}

func (ego *list) FilterInts(function func(int) bool) List {
	ego.assert()
	result := NewList()
	for _, item := range ego.val {
		val, ok := item.getVal().(int)
		if ok && function(val) {
			result.Add(val)
		}
	}
	return result
}

func (ego *list) FilterFloats(function func(float64) bool) List {
	ego.assert()
	result := NewList()
	for _, item := range ego.val {
		val, ok := item.getVal().(float64)
		if ok && function(val) {
			result.Add(val)
		}
	}
	return result
}

func (ego *list) IntSum() int {
	if !ego.AllInts() {
		panic("All elements have to be ints.")
	}
	var result int
	for _, item := range ego.val {
		result += item.getVal().(int)
	}
	return result
}

func (ego *list) Sum() float64 {
	if !ego.AllNumeric() {
		panic("All elements have to be numeric.")
	}
	var result float64
	for _, item := range ego.val {
		val, ok := item.getVal().(int)
		if ok {
			result += float64(val)
		} else {
			result += item.getVal().(float64)
		}

	}
	return result
}

func (ego *list) IntProd() int {
	if !ego.AllInts() {
		panic("All elements have to be ints.")
	}
	result := 1
	for _, item := range ego.val {
		result *= item.getVal().(int)
	}
	return result
}

func (ego *list) Prod() float64 {
	if !ego.AllNumeric() {
		panic("All elements have to be numeric.")
	}
	result := 1.0
	for _, item := range ego.val {
		val, ok := item.getVal().(int)
		if ok {
			result *= float64(val)
		} else {
			result *= item.getVal().(float64)
		}
	}
	return result
}

func (ego *list) Avg() float64 {
	return ego.ptr.Sum() / float64(ego.Count())
}

func (ego *list) IntMin() int {
	if ego.AllInts() {
		return ego.ptr.ReduceInts(math.MaxInt, func(min int, item int) int {
			if item < min {
				return item
			} else {
				return min
			}
		})
	} else {
		panic("All elements have to be ints.")
	}
}

func (ego *list) Min() float64 {
	if ego.AllNumeric() {
		return ego.ptr.Reduce(math.MaxFloat64, func(min any, item any) any {
			val, ok := item.(int)
			if ok {
				if float64(val) < min.(float64) {
					return float64(val)
				} else {
					return min
				}
			} else {
				if item.(float64) < min.(float64) {
					return item
				} else {
					return min
				}
			}
		}).(float64)
	} else {
		panic("All elements have to be numeric.")
	}
}

func (ego *list) IntMax() int {
	if ego.AllInts() {
		return ego.ptr.ReduceInts(math.MinInt, func(max int, item int) int {
			if item > max {
				return item
			} else {
				return max
			}
		})
	} else {
		panic("All elements have to be ints.")
	}
}

func (ego *list) Max() float64 {
	if ego.AllNumeric() {
		return ego.ptr.Reduce(-math.MaxFloat64, func(max any, item any) any {
			val, ok := item.(int)
			if ok {
				if float64(val) > max.(float64) {
					return float64(val)
				} else {
					return max
				}
			} else {
				if item.(float64) > max.(float64) {
					return item
				} else {
					return max
				}
			}
		}).(float64)
	} else {
		panic("All elements have to be numeric.")
	}
}

func (ego *list) ForEachAsync(function func(int, any)) List {
	ego.assert()
	var wg sync.WaitGroup
	step := func(group *sync.WaitGroup, i int, x any) {
		function(i, x)
		group.Done()
	}
	wg.Add(ego.Count())
	for i, item := range ego.val {
		go step(&wg, i, item.getVal())
	}
	wg.Wait()
	return ego.ptr
}

func (ego *list) MapAsync(function func(int, any) any) List {
	ego.assert()
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(ego.Count())
	result := NewListOf(nil, ego.Count())
	step := func(group *sync.WaitGroup, i int, x any) {
		mutex.Lock()
		result.Replace(i, function(i, x))
		mutex.Unlock()
		group.Done()
	}
	for i, item := range ego.val {
		go step(&wg, i, item.getVal())
	}
	wg.Wait()
	return result
}

func (ego *list) GetTF(tf string) any {
	ego.assert()
	if tf[0] != '#' || len(tf) < 2 {
		panic("'" + tf + "' is not a valid tree form.")
	}
	tf = tf[1:]
	dot := strings.Index(tf, ".")
	hash := strings.Index(tf, "#")
	if dot > 0 && (hash < 0 || dot < hash) {
		integer, err := strconv.ParseInt(tf[:dot], 0, bits.UintSize)
		if err != nil {
			panic("'" + tf[:dot] + "' cannot be converted to int.")
		}
		return ego.ptr.GetObject(int(integer)).GetTF(tf[dot:])
	}
	if hash > 0 && (dot < 0 || hash < dot) {
		integer, err := strconv.ParseInt(tf[:hash], 0, bits.UintSize)
		if err != nil {
			panic("'" + tf[:hash] + "' cannot be converted to int.")
		}
		return ego.ptr.GetList(int(integer)).GetTF(tf[hash:])
	}
	integer, err := strconv.ParseInt(tf, 0, bits.UintSize)
	if err != nil {
		panic("'" + tf + "' cannot be converted to int.")
	}
	return ego.ptr.Get(int(integer))
}

func (ego *list) SetTF(tf string, value any) List {
	ego.assert()
	if tf[0] != '#' || len(tf) < 2 {
		panic("'" + tf + "' is not a valid tree form.")
	}
	tf = tf[1:]
	dot := strings.Index(tf, ".")
	hash := strings.Index(tf, "#")
	if dot > 0 && (hash < 0 || dot < hash) {
		integer, err := strconv.ParseInt(tf[:dot], 0, bits.UintSize)
		if err != nil {
			panic("'" + tf[:dot] + "' cannot be converted to int.")
		}
		var object Object
		if int(integer) < ego.Count() {
			object = ego.GetObject(int(integer))
		} else {
			object = NewObject()
			ego.ptr.Insert(int(integer), object)
		}
		object.SetTF(tf[dot:], value)
		return ego.ptr
	}
	if hash > 0 && (dot < 0 || hash < dot) {
		integer, err := strconv.ParseInt(tf[:hash], 0, bits.UintSize)
		if err != nil {
			panic("'" + tf[:hash] + "' cannot be converted to int.")
		}
		var list List
		if int(integer) < ego.Count() {
			list = ego.GetList(int(integer))
		} else {
			list = NewList()
			ego.ptr.Insert(int(integer), list)
		}
		list.SetTF(tf[hash:], value)
		return ego.ptr
	}
	integer, err := strconv.ParseInt(tf, 0, bits.UintSize)
	if err != nil {
		panic("'" + tf + "' cannot be converted to int.")
	}
	return ego.ptr.Insert(int(integer), value)
}
