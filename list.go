/*
AnyType Library for Go
List (array) type
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
Interface for a list.

Extends:
  - field.
*/
type List interface {
	field

	Init(ptr List)

	// Manipulation with elements
	Add(val ...any) List
	Insert(index int, value any) List
	Replace(index int, value any) List
	Delete(index ...int) List
	Pop() List
	Clear() List

	// Getting elements
	Get(index int) any
	GetObject(index int) Object
	GetList(index int) List
	GetString(index int) string
	GetBool(index int) bool
	GetInt(index int) int
	GetFloat(index int) float64

	// TypeOf check
	TypeOf(index int) Type

	// Export
	String() string
	FormatString(indent int) string
	Slice() []any
	ObjectSlice() []Object
	ListSlice() []List
	StringSlice() []string
	BoolSlice() []bool
	IntSlice() []int
	FloatSlice() []float64

	// Features over whole list
	Clone() List
	Count() int
	Empty() bool
	Equals(another List) bool
	Concat(another List) List
	SubList(start int, end int) List
	Contains(elem any) bool
	IndexOf(elem any) int
	Sort() List
	Reverse() List

	// Checks for homogeneity
	AllObjects() bool
	AllLists() bool
	AllStrings() bool
	AllBools() bool
	AllInts() bool
	AllFloats() bool
	AllNumeric() bool

	// ForEaches
	ForEach(function func(int, any)) List
	ForEachValue(function func(any)) List
	ForEachObject(function func(Object)) List
	ForEachList(function func(List)) List
	ForEachString(function func(string)) List
	ForEachBool(function func(bool)) List
	ForEachInt(function func(int)) List
	ForEachFloat(function func(float64)) List

	// Mappings
	Map(function func(int, any) any) List
	MapValues(function func(any) any) List
	MapObjects(function func(Object) any) List
	MapLists(function func(List) any) List
	MapStrings(function func(string) any) List
	MapInts(function func(int) any) List
	MapFloats(function func(float64) any) List

	// Reductions
	Reduce(initial any, function func(any, any) any) any
	ReduceStrings(initial string, function func(string, string) string) string
	ReduceInts(initial int, function func(int, int) int) int
	ReduceFloats(initial float64, function func(float64, float64) float64) float64

	// Filters
	Filter(function func(any) bool) List
	FilterObjects(function func(Object) bool) List
	FilterLists(function func(List) bool) List
	FilterStrings(function func(string) bool) List
	FilterInts(function func(int) bool) List
	FilterFloats(function func(float64) bool) List

	// Numeric operations
	IntSum() int
	Sum() float64
	IntProd() int
	Prod() float64
	Avg() float64
	IntMin() int
	Min() float64
	IntMax() int
	Max() float64

	// Async
	ForEachAsync(function func(int, any)) List
	MapAsync(function func(int, any) any) List

	// Tree form
	GetTF(tf string) any
	SetTF(tf string, value any) List
}

/*
Slice list, a reference type. Contains a slice of elements.

Implements:
  - Fielder,
  - Lister.
*/
type SliceList struct {
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
	ego := &SliceList{val: make([]field, 0)}
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
	ego := &SliceList{val: make([]field, count)}
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
	list := NewList()
	switch s := slice.(type) {
	case []any:
		for _, item := range s {
			list.Add(item)
		}
	case []Object:
		for _, item := range s {
			list.Add(item)
		}
	case []List:
		for _, item := range s {
			list.Add(item)
		}
	case []string:
		for _, item := range s {
			list.Add(item)
		}
	case []bool:
		for _, item := range s {
			list.Add(item)
		}
	case []int:
		for _, item := range s {
			list.Add(item)
		}
	case []float64:
		for _, item := range s {
			list.Add(item)
		}
	default:
		panic("Unknown slice type.")
	}
	return list
}

/*
Asserts that the list is initialized.
*/
func (ego *SliceList) assert() {
	if ego == nil || ego.val == nil {
		panic("List is not initialized.")
	}
}

/*
Defined in the Fielder interface.
Acquires the value of the field, in this case a reference to the whole struct (List is reference type).

Returns:
  - value of the field.
*/
func (ego *SliceList) getVal() any {
	return ego.ptr
}

/*
Defined in the Fielder interface.
Creates a deep copy of the field, in this case a new list with identical elements.
Can be called recursively.

Returns:
  - deep copy of the field.
*/
func (ego *SliceList) copy() any {
	list := NewList()
	for _, value := range ego.val {
		list.Add(value.copy())
	}
	return list
}

/*
Defined in the Fielder interface.
Serializes the field into the JSON format, in this case prints all elements of the list.
Can be called recursively.

Returns:
  - string representing serialized field.
*/
func (ego *SliceList) serialize() string {
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
Defined in the Fielder interface.
Checks if the content of the field is equal to the given field.
Can be called recursively.

Returns:
  - true if the fields are equal, false otherwise.
*/
func (ego *SliceList) isEqual(another any) bool {
	list, ok := another.(List)
	if !ok || ego.Count() != list.Count() {
		return false
	}
	for i := range ego.val {
		if !ego.val[i].isEqual(list.getVal().(*SliceList).val[i]) {
			return false
		}
	}
	return true
}

/*
Initializes the ego pointer, which allows deriving.

Parameters:
  - ptr - ego pointer.
*/
func (ego *SliceList) Init(ptr List) {
	ego.ptr = ptr
}

/*
Adds new elements at the end of the list.

Parameters:
  - values... - any amount of elements to add.

Returns:
  - updated list.
*/
func (ego *SliceList) Add(values ...any) List {
	ego.assert()
	for _, val := range values {
		ego.val = append(ego.val, parseVal(val))
	}
	return ego.ptr
}

/*
Inserts a new element at the specified position in the list.

Parameters:
  - index - position where the element should be inserted,
  - value - element to insert.

Returns:
  - updated list.
*/
func (ego *SliceList) Insert(index int, value any) List {
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

/*
Replaces an existing element with a new one.

Parameters:
  - index - position of the element which should be replaced,
  - value - new element.

Returns:
  - updated list.
*/
func (ego *SliceList) Replace(index int, value any) List {
	ego.assert()
	if index < 0 || index > ego.Count() {
		panic("Index " + strconv.Itoa(index) + " out of range.")
	}
	ego.val[index] = parseVal(value)
	return ego.ptr
}

/*
Deletes the elements at the specified positions in the list.

Parameters:
  - indexes... - any amount of positions of the elements to delete.

Returns:
  - updated list.
*/
func (ego *SliceList) Delete(indexes ...int) List {
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

/*
Deletes the last element in the list.

Returns:
  - updated list.
*/
func (ego *SliceList) Pop() List {
	return ego.ptr.Delete(ego.Count() - 1)
}

/*
Deletes all elements in the list.

Returns:
  - updated list.
*/
func (ego *SliceList) Clear() List {
	ego.assert()
	ego.val = make([]field, 0)
	return ego.ptr
}

/*
Acquires the element at the specified position in the list.

Parameters:
  - index - position of the element to get.

Returns:
  - corresponding value (any type, has to be asserted).
*/
func (ego *SliceList) Get(index int) any {
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

/*
Acquires the object at the specified position in the list.
Causes a panic if the element has another type.

Parameters:
  - index - position of the element to get.

Returns:
  - corresponding value asserted as object.
*/
func (ego *SliceList) GetObject(index int) Object {
	o, ok := ego.Get(index).(Object)
	if !ok {
		panic("Item is not an object.")
	}
	return o
}

/*
Acquires the list at the specified position in the list.
Causes a panic if the element has another type.

Parameters:
  - index - position of the element to get.

Returns:
  - corresponding value asserted as list.
*/
func (ego *SliceList) GetList(index int) List {
	o, ok := ego.Get(index).(List)
	if !ok {
		panic("Item is not a list.")
	}
	return o
}

/*
Acquires the string at the specified position in the list.
Causes a panic if the element has another type.

Parameters:
  - index - position of the element to get.

Returns:
  - corresponding value asserted as string.
*/
func (ego *SliceList) GetString(index int) string {
	o, ok := ego.Get(index).(string)
	if !ok {
		panic("Item is not a string.")
	}
	return o
}

/*
Acquires the boolean at the specified position in the list.
Causes a panic if the element has another type.

Parameters:
  - index - position of the element to get.

Returns:
  - corresponding value asserted as bool.
*/
func (ego *SliceList) GetBool(index int) bool {
	o, ok := ego.Get(index).(bool)
	if !ok {
		panic("Item is not a bool.")
	}
	return o
}

/*
Acquires the integer at the specified position in the list.
Causes a panic if the element has another type.

Parameters:
  - index - position of the element to get.

Returns:
  - corresponding value asserted as int.
*/
func (ego *SliceList) GetInt(index int) int {
	o, ok := ego.Get(index).(int)
	if !ok {
		panic("Item is not an int.")
	}
	return o
}

/*
Acquires the float at the specified position in the list.
Causes a panic if the element has another type.

Parameters:
  - index - position of the element to get.

Returns:
  - corresponding value asserted as float64.
*/
func (ego *SliceList) GetFloat(index int) float64 {
	o, ok := ego.Get(index).(float64)
	if !ok {
		panic("Item is not a float.")
	}
	return o
}

/*
Gives a type of the element at the specified position in the list.

Parameters:
  - index - position of the element.

Returns:
  - integer constant representing the type (see type enum).
*/
func (ego *SliceList) TypeOf(index int) Type {
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

/*
Gives a JSON representation of the list, including nested lists and objects.

Returns:
  - JSON string.
*/
func (ego *SliceList) String() string {
	ego.assert()
	return ego.ptr.serialize()
}

/*
Gives a JSON representation of the list in standardized format with the given indentation.

Parameters:
  - indent - indentation spaces (0-10).

Returns:
  - JSON string.
*/
func (ego *SliceList) FormatString(indent int) string {
	if indent < 0 || indent > 10 {
		panic("Invalid indentation.")
	}
	buffer := new(bytes.Buffer)
	json.Indent(buffer, []byte(ego.String()), "", strings.Repeat(" ", indent))
	return buffer.String()
}

/*
Converts the list into a Go slice of empty interfaces.

Returns:
  - slice.
*/
func (ego *SliceList) Slice() []any {
	ego.assert()
	slice := make([]any, 0)
	for _, item := range ego.val {
		slice = append(slice, item.getVal())
	}
	return slice
}

/*
Converts the list of objects into a Go slice.
The list has to be homogeneous and all elements have to be objects.

Returns:
  - slice.
*/
func (ego *SliceList) ObjectSlice() []Object {
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

/*
Converts the list of lists into a Go slice.
The list has to be homogeneous and all elements have to be lists.

Returns:
  - slice.
*/
func (ego *SliceList) ListSlice() []List {
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

/*
Converts the list of strings into a Go slice.
The list has to be homogeneous and all elements have to be strings.

Returns:
  - slice.
*/
func (ego *SliceList) StringSlice() []string {
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

/*
Converts the list of bools into a Go slice.
The list has to be homogeneous and all elements have to be bools.

Returns:
  - slice.
*/
func (ego *SliceList) BoolSlice() []bool {
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

/*
Converts the list of ints into a Go slice.
The list has to be homogeneous and all elements have to be ints.

Returns:
  - slice.
*/
func (ego *SliceList) IntSlice() []int {
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

/*
Converts the list of floats into a Go slice.
The list has to be homogeneous and all elements have to be floats.

Returns:
  - slice.
*/
func (ego *SliceList) FloatSlice() []float64 {
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

/*
Creates a deep copy of the list.

Returns:
  - copied list.
*/
func (ego *SliceList) Clone() List {
	ego.assert()
	return ego.ptr.copy().(*SliceList)
}

/*
Gives a number of elements in the list.

Returns:
  - number of elements.
*/
func (ego *SliceList) Count() int {
	ego.assert()
	return len(ego.val)
}

/*
Checks whether the list is empty.

Returns:
  - true if the list is empty, false otherwise.
*/
func (ego *SliceList) Empty() bool {
	return ego.ptr.Count() == 0
}

/*
Checks if the content of the list is equal to the content of another list.
Nested objects and lists are compared recursively (by value).

Parameters:
  - another - a list to compare with.

Returns:
  - true if the lists are equal, false otherwise.
*/
func (ego *SliceList) Equals(another List) bool {
	ego.assert()
	return ego.ptr.isEqual(another)
}

/*
Creates a new list containing all elements of the old list and another list.
The old list remains unchanged.

Parameters:
  - another - a list to append.

Returns:
  - new list.
*/
func (ego *SliceList) Concat(another List) List {
	ego.assert()
	newList := &SliceList{val: append(ego.val, another.getVal().(*SliceList).val...)}
	newList.Init(newList)
	return newList
}

/*
Creates a new list containing the elements from the starting index (including) to the ending index (excluding).
If the ending index is zero, it is set to the length of the list. If negative, it is counted from the end of the list.
Starting index has to be non-negative and cannot be higher than the ending index.

Parameters:
  - start - starting index,
  - end - ending index.

Returns:
  - created sub list.
*/
func (ego *SliceList) SubList(start int, end int) List {
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
	list := &SliceList{val: make([]field, end-start)}
	list.Init(list)
	copy(list.val, ego.val[start:end])
	return list
}

/*
Checks if the list contains a given element.
Objects and lists are compared by reference.

Parameters:
  - elem - the element to check.

Returns:
  - true if the list contains the element, false otherwise.
*/
func (ego *SliceList) Contains(elem any) bool {
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

/*
Gives a position of the first occurrence of a given element.

Parameters:
  - elem - the element to check.

Returns:
  - index of the element (-1 if the list does not contain the element).
*/
func (ego *SliceList) IndexOf(elem any) int {
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

/*
Sorts elements in the list (ascending).
The list has to be homogeneous, all elements have to be either strings, ints or floats.

Returns:
  - updated list.
*/
func (ego *SliceList) Sort() List {
	ego.assert()
	switch ego.val[0].(type) {
	case *atString:
		slice := ego.StringSlice()
		sort.Strings(slice)
		ego.val = NewListFrom(slice).(*SliceList).val
	case *atInt:
		slice := ego.IntSlice()
		sort.Ints(slice)
		ego.val = NewListFrom(slice).(*SliceList).val
	case *atFloat:
		slice := ego.FloatSlice()
		sort.Float64s(slice)
		ego.val = NewListFrom(slice).(*SliceList).val
	default:
		panic("List has to be homogeneous with all its elements numeric or strings.")
	}
	return ego.ptr
}

/*
Reverses the order of elements in the list.

Returns:
  - updated list.
*/
func (ego *SliceList) Reverse() List {
	ego.assert()
	for i := ego.Count()/2 - 1; i >= 0; i-- {
		opp := ego.Count() - 1 - i
		ego.val[i], ego.val[opp] = ego.val[opp], ego.val[i]
	}
	return ego.ptr
}

/*
Checks if the list is homogeneous and all of its elements are objects.

Returns:
  - true if all elements are objects, false otherwise.
*/
func (ego *SliceList) AllObjects() bool {
	ego.assert()
	for _, item := range ego.val {
		_, ok := item.(Object)
		if !ok {
			return false
		}
	}
	return true
}

/*
Checks if the list is homogeneous and all of its elements are lists.

Returns:
  - true if all elements are lists, false otherwise.
*/
func (ego *SliceList) AllLists() bool {
	ego.assert()
	for _, item := range ego.val {
		_, ok := item.(List)
		if !ok {
			return false
		}
	}
	return true
}

/*
Checks if the list is homogeneous and all of its elements are strings.

Returns:
  - true if all elements are strings, false otherwise.
*/
func (ego *SliceList) AllStrings() bool {
	ego.assert()
	for _, item := range ego.val {
		_, ok := item.(*atString)
		if !ok {
			return false
		}
	}
	return true
}

/*
Checks if the list is homogeneous and all of its elements are bools.

Returns:
  - true if all elements are bools, false otherwise.
*/
func (ego *SliceList) AllBools() bool {
	ego.assert()
	for _, item := range ego.val {
		_, ok := item.(*atBool)
		if !ok {
			return false
		}
	}
	return true
}

/*
Checks if the list is homogeneous and all of its elements are ints.

Returns:
  - true if all elements are ints, false otherwise.
*/
func (ego *SliceList) AllInts() bool {
	ego.assert()
	for _, item := range ego.val {
		_, ok := item.(*atInt)
		if !ok {
			return false
		}
	}
	return true
}

/*
Checks if the list is homogeneous and all of its elements are floats.

Returns:
  - true if all elements are floats, false otherwise.
*/
func (ego *SliceList) AllFloats() bool {
	ego.assert()
	for _, item := range ego.val {
		_, ok := item.(*atFloat)
		if !ok {
			return false
		}
	}
	return true
}

/*
Checks if all elements of the list are numeric (ints or floats).

Returns:
  - true if all elements are numeric, false otherwise.
*/
func (ego *SliceList) AllNumeric() bool {
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

/*
Executes a given function over an every element of the list.
The function has two parameters: index of the current element and its value.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged list.
*/
func (ego *SliceList) ForEach(function func(int, any)) List {
	ego.assert()
	for i, item := range ego.val {
		function(i, item.getVal())
	}
	return ego.ptr
}

/*
Executes a given function over an every element of the list.
The function has one parameter, value of the current element.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged list.
*/
func (ego *SliceList) ForEachValue(function func(any)) List {
	ego.assert()
	for _, item := range ego.val {
		function(item.getVal())
	}
	return ego.ptr
}

/*
Executes a given function over all objects in the list.
Elements with other types are ignored.
The function has one parameter, the current object.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged list.
*/
func (ego *SliceList) ForEachObject(function func(Object)) List {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.(Object)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

/*
Executes a given function over all lists nested in the list.
Elements with other types are ignored.
The function has one parameter, the current list.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged list.
*/
func (ego *SliceList) ForEachList(function func(List)) List {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.(List)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

/*
Executes a given function over all strings in the list.
Elements with other types are ignored.
The function has one parameter, the current string.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged list.
*/
func (ego *SliceList) ForEachString(function func(string)) List {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(string)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

/*
Executes a given function over all bools in the list.
Elements with other types are ignored.
The function has one parameter, the current bool.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged list.
*/
func (ego *SliceList) ForEachBool(function func(bool)) List {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(bool)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

/*
Executes a given function over all ints in the list.
Elements with other types are ignored.
The function has one parameter, the current int.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged list.
*/
func (ego *SliceList) ForEachInt(function func(int)) List {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(int)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

/*
Executes a given function over all floats in the list.
Elements with other types are ignored.
The function has one parameter, the current float.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged list.
*/
func (ego *SliceList) ForEachFloat(function func(float64)) List {
	ego.assert()
	for _, item := range ego.val {
		val, ok := item.getVal().(float64)
		if ok {
			function(val)
		}
	}
	return ego.ptr
}

/*
Copies the list and modifies each element by a given mapping function.
The resulting element can have a different type than the original one.
The function has two parameters: current index and value of the current element. Returns empty interface.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new list.
*/
func (ego *SliceList) Map(function func(int, any) any) List {
	ego.assert()
	result := NewList()
	for i, item := range ego.val {
		result.Add(function(i, item.getVal()))
	}
	return result
}

/*
Copies the list and modifies each element by a given mapping function.
The resulting element can have a different type than the original one.
The function has one parameter, value of the current element, and returns empty interface.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new list.
*/
func (ego *SliceList) MapValues(function func(any) any) List {
	ego.assert()
	result := NewList()
	for _, item := range ego.val {
		result.Add(function(item.getVal()))
	}
	return result
}

/*
Selects all objects from the list and modifies each of them by a given mapping function.
Elements with other types are ignored.
The resulting element can have a different type than the original one.
The function has one parameter, the current object, and returns empty interface.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new list.
*/
func (ego *SliceList) MapObjects(function func(Object) any) List {
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

/*
Selects all nested lists from the list and modifies each of them by a given mapping function.
Elements with other types are ignored.
The resulting element can have a different type than the original one.
The function has one parameter, the current list, and returns empty interface.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new list.
*/
func (ego *SliceList) MapLists(function func(List) any) List {
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

/*
Selects all strings from the list and modifies each of them by a given mapping function.
Elements with other types are ignored.
The resulting element can have a different type than the original one.
The function has one parameter, the current string, and returns empty interface.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new list.
*/
func (ego *SliceList) MapStrings(function func(string) any) List {
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

/*
Selects all ints from the list and modifies each of them by a given mapping function.
Elements with other types are ignored.
The resulting element can have a different type than the original one.
The function has one parameter, the current int, and returns empty interface.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new list.
*/
func (ego *SliceList) MapInts(function func(int) any) List {
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

/*
Selects all floats from the list and modifies each of them by a given mapping function.
Elements with other types are ignored.
The resulting element can have a different type than the original one.
The function has one parameter, the current float, and returns empty interface.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new list.
*/
func (ego *SliceList) MapFloats(function func(float64) any) List {
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

/*
Reduces all elements of the list into a single value.
The function has two parameters: value returned by the previous iteration and value of the current element. Returns empty interface.
The list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - computed value.
*/
func (ego *SliceList) Reduce(initial any, function func(any, any) any) any {
	ego.assert()
	result := initial
	for _, item := range ego.val {
		result = function(result, item.getVal())
	}
	return result
}

/*
Reduces all strings in the list into a single string.
Elements with other types are ignored.
The function has two parameters: string returned by the previous iteration and current string. Returns string.
The list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - computed value.
*/
func (ego *SliceList) ReduceStrings(initial string, function func(string, string) string) string {
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

/*
Reduces all ints in the list into a single int.
Elements with other types are ignored.
The function has two parameters: int returned by the previous iteration and current int. Returns int.
The list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - computed value.
*/
func (ego *SliceList) ReduceInts(initial int, function func(int, int) int) int {
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

/*
Reduces all floats in the list into a single float.
Elements with other types are ignored.
The function has two parameters: float returned by the previous iteration and current float. Returns float.
The list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - computed value.
*/
func (ego *SliceList) ReduceFloats(initial float64, function func(float64, float64) float64) float64 {
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

/*
Creates a new list containing elements of the old one, satisfying a condition.
The function has one parameter, value of the current element, and returns bool.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - filtered list.
*/
func (ego *SliceList) Filter(function func(any) bool) List {
	ego.assert()
	result := NewList()
	for _, item := range ego.val {
		if function(item.getVal()) {
			result.Add(item.getVal())
		}
	}
	return result
}

/*
Creates a new list containing objects of the old one, satisfying a condition.
Elements with other types are ignored.
The function has one parameter, current object, and returns bool.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - filtered list.
*/
func (ego *SliceList) FilterObjects(function func(Object) bool) List {
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

/*
Creates a new list containing nested lists of the old one, satisfying a condition.
Elements with other types are ignored.
The function has one parameter, current list, and returns bool.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - filtered list.
*/
func (ego *SliceList) FilterLists(function func(List) bool) List {
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

/*
Creates a new list containing strings of the old one, satisfying a condition.
Elements with other types are ignored.
The function has one parameter, current string, and returns bool.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - filtered list.
*/
func (ego *SliceList) FilterStrings(function func(string) bool) List {
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

/*
Creates a new list containing ints of the old one, satisfying a condition.
Elements with other types are ignored.
The function has one parameter, current int, and returns bool.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - filtered list.
*/
func (ego *SliceList) FilterInts(function func(int) bool) List {
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

/*
Creates a new list containing floats of the old one, satisfying a condition.
Elements with other types are ignored.
The function has one parameter, current float, and returns bool.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - filtered list.
*/
func (ego *SliceList) FilterFloats(function func(float64) bool) List {
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

/*
Computes a sum of all elements in the list.
The list has to be homogeneous and all its elements have to be ints.

Returns:
  - computed sum (int).
*/
func (ego *SliceList) IntSum() int {
	if !ego.AllInts() {
		panic("All elements have to be ints.")
	}
	var result int
	for _, item := range ego.val {
		result += item.getVal().(int)
	}
	return result
}

/*
Computes a sum of all elements in the list.
The list has to be homogeneous and all its elements have to be numeric.

Returns:
  - computed sum (float).
*/
func (ego *SliceList) Sum() float64 {
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

/*
Computes a product of all elements in the list.
The list has to be homogeneous and all its elements have to be ints.

Returns:
  - computed product (int).
*/
func (ego *SliceList) IntProd() int {
	if !ego.AllInts() {
		panic("All elements have to be ints.")
	}
	result := 1
	for _, item := range ego.val {
		result *= item.getVal().(int)
	}
	return result
}

/*
Computes a product of all elements in the list.
The list has to be homogeneous and all its elements have to be numeric.

Returns:
  - computed pruduct (float).
*/
func (ego *SliceList) Prod() float64 {
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

/*
Computes an arithmetic average of all elements in the list.
The list has to be homogeneous and all its elements have to be numeric.

Returns:
  - computed average (float).
*/
func (ego *SliceList) Avg() float64 {
	return ego.ptr.Sum() / float64(ego.Count())
}

/*
Finds a minimum of the list.
The list has to be homogeneous and all its elements have to be ints.

Returns:
  - found minimum (int).
*/
func (ego *SliceList) IntMin() int {
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

/*
Finds a minimum of the list.
The list has to be homogeneous and all its elements have to be numeric.

Returns:
  - found minimum (float).
*/
func (ego *SliceList) Min() float64 {
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

/*
Finds a maximum of the list.
The list has to be homogeneous and all its elements have to be ints.

Returns:
  - found maximum (int).
*/
func (ego *SliceList) IntMax() int {
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

/*
Finds a maximum of the list.
The list has to be homogeneous and all its elements have to be numeric.

Returns:
  - found maximum (float).
*/
func (ego *SliceList) Max() float64 {
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

/*
Parallelly executes a given function over an every element of the list.
The function has two parameters: index of the current element and its value.
The order of the iterations is random.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - unchanged list.
*/
func (ego *SliceList) ForEachAsync(function func(int, any)) List {
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

/*
Copies the list and paralelly modifies each element by a given mapping function.
The resulting element can have a different type than the original one.
The function has two parameters: index of the current element and its value.
The old list remains unchanged.

Parameters:
  - function - anonymous function to be executed.

Returns:
  - new list.
*/
func (ego *SliceList) MapAsync(function func(int, any) any) List {
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

/*
Acquires the element specified by the given tree form.

Parameters:
  - tf - tree form string.

Returns:
  - corresponding value (any type, has to be asserted).
*/
func (ego *SliceList) GetTF(tf string) any {
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

/*
Sets the element specified by the given tree form.

Parameters:
  - tf - tree form string,
  - value - value to set.

Returns:
  - updated list.
*/
func (ego *SliceList) SetTF(tf string, value any) List {
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
