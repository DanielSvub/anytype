/*
AnyType Library for Go
List (array) implementation
*/

package anytype

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/bits"
	"sort"
	"strconv"
	"strings"
	"sync"
)

/*
Slice list, a reference type. Contains a slice.

Implements:
  - field,
  - List.
*/
type list struct {
	val []field
	ptr List
}

/*
NewList creates a new list.

Parameters:
  - values... - any amount of initial elements.

Returns:
  - pointer to the created list.
*/
func NewList(values ...any) List {
	ego := &list{val: []field{}}
	ego.Init(ego)
	ego.Add(values...)
	return ego
}

/*
NewListOf creates a new list of n repeated values.

Parameters:
  - value - value to repeat,
  - count - number of repetitions.

Returns:
  - pointer to the created list.
*/
func NewListOf(value any, count int) List {
	ego := &list{val: make([]field, 0, count)}
	ego.Init(ego)
	elem := parseVal(value)
	for i := 0; i < count; i++ {
		ego.val = append(ego.val, elem)
	}
	return ego
}

/*
NewListFrom converts a slice of supported types to a list.

Parameters:
  - slice - original slice.

Returns:
  - created list.
*/
func NewListFrom(slice any) List {
	var ego List
	init := func(cap int) {
		ego = &list{val: make([]field, 0, cap)}
		ego.Init(ego)
	}
	switch s := slice.(type) {
	case []any:
		init(len(s))
		for _, item := range s {
			ego.Add(item)
		}
	case []Object:
		init(len(s))
		for _, item := range s {
			ego.Add(item)
		}
	case []List:
		init(len(s))
		for _, item := range s {
			ego.Add(item)
		}
	case []string:
		init(len(s))
		for _, item := range s {
			ego.Add(item)
		}
	case []bool:
		init(len(s))
		for _, item := range s {
			ego.Add(item)
		}
	case []int:
		init(len(s))
		for _, item := range s {
			ego.Add(item)
		}
	case []float64:
		init(len(s))
		for _, item := range s {
			ego.Add(item)
		}
	default:
		panic("unsupported slice type")
	}
	return ego
}

/*
Defined in the field interface.
Acquires the value of the field, in this case a reference to the whole struct (List is a reference type).

Returns:
  - value of the field.
*/
func (ego *list) getVal() any {
	return ego.Ego()
}

/*
Defined in the field interface.
Creates a deep copy of the field, in this case a new list with identical elements.
Can be called recursively.

Returns:
  - deep copy of the field.
*/
func (ego *list) copy() any {
	list := &list{val: make([]field, ego.Ego().Count())}
	list.Init(list)
	for i, value := range ego.val {
		list.val[i] = parseVal(value.copy())
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
	var result strings.Builder
	result.WriteRune('[')
	for i, value := range ego.val {
		result.WriteString(value.serialize())
		if i+1 < len(ego.val) {
			result.WriteRune(',')
		}
	}
	result.WriteRune(']')
	return result.String()
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
	if !ok || ego.Ego().Count() != list.Count() {
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
	for _, val := range values {
		ego.val = append(ego.val, parseVal(val))
	}
	return ego.Ego()
}

func (ego *list) Insert(index int, value any) List {
	if index < 0 || index > ego.Ego().Count() {
		panic(fmt.Sprintf("index %d out of range with count %d", index, ego.Ego().Count()))
	}
	if index == ego.Ego().Count() {
		return ego.Ego().Add(value)
	}
	ego.val = append(ego.val[:index+1], ego.val[index:]...)
	ego.val[index] = parseVal(value)
	return ego.Ego()
}

func (ego *list) Replace(index int, value any) List {
	if index < 0 || index >= ego.Ego().Count() {
		panic(fmt.Sprintf("index %d out of range with count %d", index, ego.Ego().Count()))
	}
	ego.val[index] = parseVal(value)
	return ego.Ego()
}

func (ego *list) Delete(indexes ...int) List {
	if len(indexes) > 1 {
		sort.Ints(indexes)
	}
	for i := len(indexes) - 1; i >= 0; i-- {
		index := indexes[i]
		if index < 0 || index >= ego.Ego().Count() {
			panic(fmt.Sprintf("index %d out of range with count %d", index, ego.Ego().Count()))
		}
		ego.val = append(ego.val[:index], ego.val[index+1:]...)
	}
	return ego.Ego()
}

func (ego *list) Pop() List {
	return ego.Ego().Delete(ego.Ego().Count() - 1)
}

func (ego *list) Clear() List {
	ego.val = []field{}
	return ego.Ego()
}

func (ego *list) Get(index int) any {
	if len(ego.val) <= index || index < 0 {
		panic(fmt.Sprintf("index %d out of range with count %d", index, ego.Ego().Count()))
	}
	return ego.val[index].getVal()
}

func (ego *list) GetObject(index int) Object {
	o, ok := ego.Get(index).(Object)
	if !ok {
		panic("item is not an object")
	}
	return o
}

func (ego *list) GetList(index int) List {
	o, ok := ego.Get(index).(List)
	if !ok {
		panic("item is not a list")
	}
	return o
}

func (ego *list) GetString(index int) string {
	o, ok := ego.Get(index).(string)
	if !ok {
		panic("item is not a string")
	}
	return o
}

func (ego *list) GetBool(index int) bool {
	o, ok := ego.Get(index).(bool)
	if !ok {
		panic("item is not a bool")
	}
	return o
}

func (ego *list) GetInt(index int) int {
	o, ok := ego.Get(index).(int)
	if !ok {
		panic("item is not an int")
	}
	return o
}

func (ego *list) GetFloat(index int) float64 {
	o, ok := ego.Get(index).(float64)
	if !ok {
		panic("item is not a float")
	}
	return o
}

func (ego *list) TypeOf(index int) Type {
	if index >= 0 && index < ego.Ego().Count() {
		switch ego.val[index].(type) {
		case Object:
			return TypeObject
		case List:
			return TypeList
		case *atNil:
			return TypeNil
		case *atString:
			return TypeString
		case *atInt:
			return TypeInt
		case *atBool:
			return TypeBool
		case *atFloat:
			return TypeFloat
		}
	}
	return TypeUndefined
}

func (ego *list) String() string {
	return ego.Ego().serialize()
}

func (ego *list) FormatString(indent int) string {
	if indent < 0 || indent > 10 {
		panic("invalid indentation")
	}
	buffer := new(bytes.Buffer)
	json.Indent(buffer, []byte(ego.String()), "", strings.Repeat(" ", indent))
	return buffer.String()
}

func (ego *list) Slice() []any {
	slice := make([]any, 0, ego.Ego().Count())
	for _, item := range ego.val {
		slice = append(slice, item.getVal())
	}
	return slice
}

func (ego *list) NativeSlice() []any {
	return native(ego).([]any)
}

func (ego *list) ObjectSlice() []Object {
	slice := make([]Object, 0, ego.Ego().Count())
	for _, item := range ego.val {
		object, ok := item.(Object)
		if ok {
			slice = append(slice, object)
		}
	}
	return slice
}

func (ego *list) ListSlice() []List {
	slice := make([]List, 0, ego.Ego().Count())
	for _, item := range ego.val {
		list, ok := item.(List)
		if ok {
			slice = append(slice, list)
		}
	}
	return slice
}

func (ego *list) StringSlice() []string {
	slice := make([]string, 0, ego.Ego().Count())
	for _, item := range ego.val {
		str, ok := item.getVal().(string)
		if ok {
			slice = append(slice, str)
		}
	}
	return slice
}

func (ego *list) BoolSlice() []bool {
	slice := make([]bool, 0, ego.Ego().Count())
	for _, item := range ego.val {
		boolean, ok := item.getVal().(bool)
		if ok {
			slice = append(slice, boolean)
		}
	}
	return slice
}

func (ego *list) IntSlice() []int {
	slice := make([]int, 0, ego.Ego().Count())
	for _, item := range ego.val {
		integer, ok := item.getVal().(int)
		if ok {
			slice = append(slice, integer)
		}
	}
	return slice
}

func (ego *list) FloatSlice() []float64 {
	slice := make([]float64, 0, ego.Ego().Count())
	for _, item := range ego.val {
		float, ok := item.getVal().(float64)
		if ok {
			slice = append(slice, float)
		}
	}
	return slice
}

func (ego *list) Clone() List {
	return ego.copy().(*list)
}

func (ego *list) Count() int {
	return len(ego.val)
}

func (ego *list) Empty() bool {
	return ego.Ego().Count() == 0
}

func (ego *list) Equals(another List) bool {
	return ego.isEqual(another)
}

func (ego *list) Concat(another List) List {
	newList := &list{val: append(ego.val, another.getVal().(*list).val...)}
	newList.Init(newList)
	return newList
}

func (ego *list) SubList(start int, end int) List {
	if end > ego.Ego().Count() || end < -ego.Ego().Count() {
		panic(fmt.Sprintf("ending index %d out of range with count %d", end, ego.Ego().Count()))
	}
	if end <= 0 {
		end = ego.Ego().Count() + end
	}
	if start > end {
		panic("starting index is higher than the ending index")
	}
	if start < 0 {
		panic("starting index is lower than zero")
	}
	list := &list{val: make([]field, end-start)}
	list.Init(list)
	copy(list.val, ego.val[start:end])
	return list
}

func (ego *list) Contains(elem any) bool {
	for _, item := range ego.val {
		if item.getVal() == elem {
			return true
		}
	}
	return false
}

func (ego *list) IndexOf(elem any) int {
	for i, item := range ego.val {
		if item.getVal() == elem {
			return i
		}
	}
	return -1
}

func (ego *list) Sort() List {
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
		panic("the first element of the list has to be either string, int or float")
	}
	return ego.Ego()
}

func (ego *list) Reverse() List {
	for i := ego.Ego().Count()/2 - 1; i >= 0; i-- {
		opp := ego.Ego().Count() - 1 - i
		ego.val[i], ego.val[opp] = ego.val[opp], ego.val[i]
	}
	return ego.Ego()
}

func (ego *list) AllObjects() bool {
	for _, item := range ego.val {
		_, ok := item.(Object)
		if !ok {
			return false
		}
	}
	return true
}

func (ego *list) AllLists() bool {
	for _, item := range ego.val {
		_, ok := item.(List)
		if !ok {
			return false
		}
	}
	return true
}

func (ego *list) AllStrings() bool {
	for _, item := range ego.val {
		_, ok := item.(*atString)
		if !ok {
			return false
		}
	}
	return true
}

func (ego *list) AllBools() bool {
	for _, item := range ego.val {
		_, ok := item.(*atBool)
		if !ok {
			return false
		}
	}
	return true
}

func (ego *list) AllInts() bool {
	for _, item := range ego.val {
		_, ok := item.(*atInt)
		if !ok {
			return false
		}
	}
	return true
}

func (ego *list) AllFloats() bool {
	for _, item := range ego.val {
		_, ok := item.(*atFloat)
		if !ok {
			return false
		}
	}
	return true
}

func (ego *list) AllNumeric() bool {
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
	for i, item := range ego.val {
		function(i, item.getVal())
	}
	return ego.Ego()
}

func (ego *list) ForEachValue(function func(any)) List {
	for _, item := range ego.val {
		function(item.getVal())
	}
	return ego.Ego()
}

func (ego *list) ForEachObject(function func(Object)) List {
	for _, item := range ego.val {
		val, ok := item.(Object)
		if ok {
			function(val)
		}
	}
	return ego.Ego()
}

func (ego *list) ForEachList(function func(List)) List {
	for _, item := range ego.val {
		val, ok := item.(List)
		if ok {
			function(val)
		}
	}
	return ego.Ego()
}

func (ego *list) ForEachString(function func(string)) List {
	for _, item := range ego.val {
		val, ok := item.getVal().(string)
		if ok {
			function(val)
		}
	}
	return ego.Ego()
}

func (ego *list) ForEachBool(function func(bool)) List {
	for _, item := range ego.val {
		val, ok := item.getVal().(bool)
		if ok {
			function(val)
		}
	}
	return ego.Ego()
}

func (ego *list) ForEachInt(function func(int)) List {
	for _, item := range ego.val {
		val, ok := item.getVal().(int)
		if ok {
			function(val)
		}
	}
	return ego.Ego()
}

func (ego *list) ForEachFloat(function func(float64)) List {
	for _, item := range ego.val {
		val, ok := item.getVal().(float64)
		if ok {
			function(val)
		}
	}
	return ego.Ego()
}

func (ego *list) Map(function func(int, any) any) List {
	result := NewList()
	for i, item := range ego.val {
		result.Add(function(i, item.getVal()))
	}
	return result
}

func (ego *list) MapValues(function func(any) any) List {
	result := NewList()
	for _, item := range ego.val {
		result.Add(function(item.getVal()))
	}
	return result
}

func (ego *list) MapObjects(function func(Object) any) List {
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
	result := NewList()
	for _, item := range ego.val {
		val, ok := item.getVal().(string)
		if ok {
			result.Add(function(val))
		}
	}
	return result
}

func (ego *list) MapBools(function func(bool) any) List {
	result := NewList()
	for _, item := range ego.val {
		val, ok := item.getVal().(bool)
		if ok {
			result.Add(function(val))
		}
	}
	return result
}

func (ego *list) MapInts(function func(int) any) List {
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
	result := initial
	for _, item := range ego.val {
		result = function(result, item.getVal())
	}
	return result
}

func (ego *list) ReduceStrings(initial string, function func(string, string) string) string {
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
	result := NewList()
	for _, item := range ego.val {
		if function(item.getVal()) {
			result.Add(item.getVal())
		}
	}
	return result
}

func (ego *list) FilterObjects(function func(Object) bool) List {
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
	result := NewList()
	for _, item := range ego.val {
		val, ok := item.getVal().(float64)
		if ok && function(val) {
			result.Add(val)
		}
	}
	return result
}

func (ego *list) IntSum() (result int) {
	for _, item := range ego.val {
		value, ok := item.getVal().(int)
		if ok {
			result += value
		}
	}
	return
}

func (ego *list) Sum() (result float64) {
	for _, item := range ego.val {
		val, ok := item.getVal().(int)
		if ok {
			result += float64(val)
		} else if val, ok := item.getVal().(float64); ok {
			result += val
		}
	}
	return result
}

func (ego *list) IntProd() (result int) {
	result = 1
	for _, item := range ego.val {
		value, ok := item.getVal().(int)
		if ok {
			result *= value
		}
	}
	return
}

func (ego *list) Prod() (result float64) {
	result = 1
	for _, item := range ego.val {
		val, ok := item.getVal().(int)
		if ok {
			result *= float64(val)
		} else if val, ok := item.getVal().(float64); ok {
			result *= val
		}
	}
	return result
}

func (ego *list) Avg() float64 {
	return ego.Ego().Sum() / float64(ego.Ego().Count())
}

func (ego *list) IntMin() int {
	var present bool
	min := ego.Ego().ReduceInts(math.MaxInt, func(min int, item int) int {
		present = true
		if item < min {
			return item
		} else {
			return min
		}
	})
	if present {
		return min
	} else {
		return 0
	}
}

func (ego *list) Min() float64 {
	var present bool
	min := ego.Ego().Reduce(math.MaxFloat64, func(min any, item any) any {
		present = true
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
	if present {
		return min
	} else {
		return 0
	}
}

func (ego *list) IntMax() int {
	var present bool
	max := ego.Ego().ReduceInts(math.MinInt, func(max int, item int) int {
		present = true
		if item > max {
			return item
		} else {
			return max
		}
	})
	if present {
		return max
	} else {
		return 0
	}
}

func (ego *list) Max() float64 {
	var present bool
	max := ego.Ego().Reduce(-math.MaxFloat64, func(max any, item any) any {
		present = true
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
	if present {
		return max
	} else {
		return 0
	}
}

func (ego *list) ForEachAsync(function func(int, any)) List {
	var wg sync.WaitGroup
	step := func(group *sync.WaitGroup, i int, x any) {
		function(i, x)
		group.Done()
	}
	wg.Add(ego.Ego().Count())
	for i, item := range ego.val {
		go step(&wg, i, item.getVal())
	}
	wg.Wait()
	return ego.Ego()
}

func (ego *list) MapAsync(function func(int, any) any) List {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(ego.Ego().Count())
	result := NewListOf(nil, ego.Ego().Count())
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
	if len(tf) < 2 || tf[0] != '#' {
		panic(fmt.Sprintf("'%s' is not a valid tree form for a list", tf))
	}
	tf = tf[1:]
	dot := strings.Index(tf, ".")
	hash := strings.Index(tf, "#")
	if dot > 0 && (hash < 0 || dot < hash) {
		integer, err := strconv.ParseInt(tf[:dot], 0, bits.UintSize)
		if err != nil {
			panic(fmt.Sprintf("'%s' cannot be converted to int", tf[:dot]))
		}
		return ego.Ego().GetObject(int(integer)).GetTF(tf[dot:])
	}
	if hash > 0 && (dot < 0 || hash < dot) {
		integer, err := strconv.ParseInt(tf[:hash], 0, bits.UintSize)
		if err != nil {
			panic(fmt.Sprintf("'%s' cannot be converted to int", tf[:hash]))
		}
		return ego.Ego().GetList(int(integer)).GetTF(tf[hash:])
	}
	integer, err := strconv.ParseInt(tf, 0, bits.UintSize)
	if err != nil {
		panic(fmt.Sprintf("'%s' cannot be converted to int", tf))
	}
	return ego.Ego().Get(int(integer))
}

func (ego *list) SetTF(tf string, value any) List {
	if len(tf) < 2 || tf[0] != '#' {
		panic(fmt.Sprintf("'%s' is not a valid tree form for a list", tf))
	}
	tf = tf[1:]
	dot := strings.Index(tf, ".")
	hash := strings.Index(tf, "#")
	if dot > 0 && (hash < 0 || dot < hash) {
		integer, err := strconv.ParseInt(tf[:dot], 0, bits.UintSize)
		if err != nil {
			panic(fmt.Sprintf("'%s' cannot be converted to int", tf[:dot]))
		}
		var object Object
		index := int(integer)
		if index == ego.Ego().Count() {
			object = NewObject()
			ego.Ego().Add(object)
		} else {
			if ego.Ego().TypeOf(index) == TypeObject {
				object = ego.Ego().GetObject(index)
			} else {
				object = NewObject()
				ego.Ego().Replace(index, object)
			}
		}
		object.SetTF(tf[dot:], value)
		return ego.Ego()
	}
	if hash > 0 && (dot < 0 || hash < dot) {
		integer, err := strconv.ParseInt(tf[:hash], 0, bits.UintSize)
		if err != nil {
			panic(fmt.Sprintf("'%s' cannot be converted to int", tf[:hash]))
		}
		var list List
		index := int(integer)
		if index == ego.Ego().Count() {
			list = NewList()
			ego.Ego().Add(list)
		} else {
			if ego.Ego().TypeOf(index) == TypeList {
				list = ego.Ego().GetList(index)
			} else {
				list = NewList()
				ego.Ego().Replace(index, list)
			}
		}
		list.SetTF(tf[hash:], value)
		return ego.Ego()
	}
	integer, err := strconv.ParseInt(tf, 0, bits.UintSize)
	index := int(integer)
	if err != nil {
		panic(fmt.Sprintf("'%s' cannot be converted to int", tf))
	}
	if index == ego.Ego().Count() {
		return ego.Ego().Add(value)
	}
	return ego.Ego().Replace(index, value)
}

func (ego *list) UnsetTF(tf string) List {
	if len(tf) < 2 || tf[0] != '#' {
		panic(fmt.Sprintf("'%s' is not a valid tree form for a list", tf))
	}
	tf = tf[1:]
	dot := strings.Index(tf, ".")
	hash := strings.Index(tf, "#")
	if dot > 0 && (hash < 0 || dot < hash) {
		integer, err := strconv.ParseInt(tf[:dot], 0, bits.UintSize)
		if err != nil {
			panic(fmt.Sprintf("'%s' cannot be converted to int", tf[:dot]))
		}
		object := ego.GetObject(int(integer))
		object.UnsetTF(tf[dot:])
		return ego.Ego()
	}
	if hash > 0 && (dot < 0 || hash < dot) {
		integer, err := strconv.ParseInt(tf[:hash], 0, bits.UintSize)
		if err != nil {
			panic(fmt.Sprintf("'%s' cannot be converted to int", tf[:hash]))
		}
		list := ego.GetList(int(integer))
		list.UnsetTF(tf[hash:])
		return ego.Ego()
	}
	integer, err := strconv.ParseInt(tf, 0, bits.UintSize)
	if err != nil {
		panic(fmt.Sprintf("'%s' cannot be converted to int", tf))
	}
	return ego.Ego().Delete(int(integer))
}

func (ego *list) TypeOfTF(tf string) Type {
	if len(tf) < 2 || tf[0] != '#' {
		return TypeUndefined
	}
	tf = tf[1:]
	dot := strings.Index(tf, ".")
	hash := strings.Index(tf, "#")
	if dot > 0 && (hash < 0 || dot < hash) {
		integer, err := strconv.ParseInt(tf[:dot], 0, bits.UintSize)
		if err != nil {
			return TypeUndefined
		}
		index := int(integer)
		if ego.Ego().TypeOf(index) != TypeObject {
			return TypeUndefined
		}
		return ego.Ego().GetObject(int(integer)).TypeOfTF(tf[dot:])
	}
	if hash > 0 && (dot < 0 || hash < dot) {
		integer, err := strconv.ParseInt(tf[:hash], 0, bits.UintSize)
		if err != nil {
			return TypeUndefined
		}
		index := int(integer)
		if ego.Ego().TypeOf(index) != TypeList {
			return TypeUndefined
		}
		return ego.Ego().GetList(index).TypeOfTF(tf[hash:])
	}
	integer, err := strconv.ParseInt(tf, 0, bits.UintSize)
	if err != nil {
		return TypeUndefined
	}
	return ego.Ego().TypeOf(int(integer))
}
