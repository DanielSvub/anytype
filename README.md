# AnyType

AnyType is a Go library providing dynamic data structures with a JSON support. It contains a number of advanced features with API inspired by Java collections.

It supports following data types compatible with JSON standard:
- string
- integer (number)
- float (number)
- boolean
- object
- list (array)
- nil (null)

## Objects

Object is an unordered set of key-value pairs. The default implementation, `mapObject` structure, is based on built-in Go maps. It is possible to make custom implementations by implementing the `Object` interface.

### Manipulation With Fields
- `Set(values ...any) Object`
- `Unset(keys ...string) Object`
- `Clear() Object`

### Getting Fields
- `Get(key string) any`
- `GetObject(key string) Object`

### Type Check
- `TypeOf(key string) Type`

### Export
- `String() string`
- `FormatString(indent int) string`
- `Dict() map[string]any`
- `Keys() List`
- `Values() List`

### Features Over Whole Object
- `Clone() Object`
- `Count() int`
- `Empty() bool`
- `Equals(another Object) bool`
- `Merge(another Object) Object`
- `Pluck(keys ...string) Object`
- `Contains(elem any) bool`
- `KeyOf(elem any) string`
- `KeyExists(key string) bool`

### ForEaches
- `ForEach(function func(string, any)) Object`
- `ForEachValue(function func(any)) Object`
- `ForEachObject(function func(Object)) Object`

### Mapping
- `Map(function func(string, any) any) Object`
- `MapValues(function func(any) any) Object`
- `MapObjects(function func(Object) any) Object`

### Async
- `ForEachAsync(function func(string, any)) Object`
- `MapAsync(function func(string, any) any) Object`

### Tree Form
- `GetTF(tf string) any`
- `SetTF(tf string, value any) Object`

## Lists

List is an ordered sequence of values. The default implementation, `sliceList` structure, is based on built-in Go slices. It is possible to make custom implementations by implementing the `List` interface.

### Manipulation With Elements
- `Add(val ...any) List`
- `Insert(index int, value any) List`
- `Replace(index int, value any) List`
- `Delete(index ...int) List`
- `Pop() List`
- `Clear() List`

### Getting Elements
- `Get(index int) any`
- `GetObject(index int) Object`
- `GetList(index int) List`
- `GetString(index int) string`
- `GetBool(index int) bool`
- `GetInt(index int) int`
- `GetFloat(index int) float64`

### Type Check
- `TypeOf(index int) Type`

### Export
- `String() string`
- `FormatString(indent int) string`
- `Slice() []any`
- `ObjectSlice() []Object`
- `ListSlice() []List`
- `StringSlice() []string`
- `BoolSlice() []bool`
- `IntSlice() []int`
- `FloatSlice() []float64`

### Features Over Whole List
- `Clone() List`
- `Count() int`
- `Empty() bool`
- `Equals(another List) bool`
- `Concat(another List) List`
- `SubList(start int, end int) List`
- `Contains(elem any) bool`
- `IndexOf(elem any) int`
- `Sort() List`
- `Reverse() List`

### Checks For Homogeneity
- `AllObjects() bool`
- `AllLists() bool`
- `AllStrings() bool`
- `AllBools() bool`
- `AllInts() bool`
- `AllFloats() bool`
- `AllNumeric() bool`

### ForEaches
- `ForEach(function func(int, any)) List`
- `ForEachValue(function func(any)) List`
- `ForEachObject(function func(Object)) List`
- `ForEachList(function func(List)) List`
- `ForEachString(function func(string)) List`
- `ForEachBool(function func(bool)) List`
- `ForEachInt(function func(int)) List`
- `ForEachFloat(function func(float64)) List`

### Mappings
- `Map(function func(int, any) any) List`
- `MapValues(function func(any) any) List`
- `MapObjects(function func(Object) any) List`
- `MapLists(function func(List) any) List`
- `MapStrings(function func(string) any) List`
- `MapInts(function func(int) any) List`
- `MapFloats(function func(float64) any) List`

### Reductions
- `Reduce(initial any, function func(any, any) any) any`
- `ReduceStrings(initial string, function func(string, string) string) string`
- `ReduceInts(initial int, function func(int, int) int) int`
- `ReduceFloats(initial float64, function func(float64, float64) float64) float64`

### Filters
- `Filter(function func(any) bool) List`
- `FilterObjects(function func(Object) bool) List`
- `FilterLists(function func(List) bool) List`
- `FilterStrings(function func(string) bool) List`
- `FilterInts(function func(int) bool) List`
- `FilterFloats(function func(float64) bool) List`

### Numeric Operations
- `IntSum() int`
- `Sum() float64`
- `IntProd() int`
- `Prod() float64`
- `Avg() float64`
- `IntMin() int`
- `Min() float64`
- `IntMax() int`
- `Max() float64`

### Async
- `ForEachAsync(function func(int, any)) List`
- `MapAsync(function func(int, any) any) List`

### Tree Form
- `GetTF(tf string) any`
- `SetTF(tf string, value any) List`

## JSON

AnyType objects can be created from a JSON string using `ParseJson(json string) Object` or `ParseFile(path string) Object` functions.

