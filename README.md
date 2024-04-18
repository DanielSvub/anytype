# AnyType

AnyType is a Go library providing dynamic data structures with JSON support. It contains a number of advanced features with API inspired by Java collections.

It supports following data types compatible with JSON standard:
- string
- integer (number)
- float (number)
- boolean
- object
- list (array)
- nil (null)

Types can be referenced by the `Type` enum (e.g. `anytype.TypeObject`, `anytype.TypeInt`, ...).

## Objects

Object is an unordered set of key-value pairs. The default implementation, `MapObject` structure, is based on built-in Go maps. It is possible to make custom implementations by implementing the `Object` interface.

### Constructors

- `NewObject(vals ...any) Object` - Initial object values are specified as key value pairs. The function panics if an odd number of arguments is given,
```go
emptyObject := anytype.NewObject()
object := anytype.NewObject(
    "first", 1,
    "second", 2, 
)
```

- `NewObjectFrom(dict any) Object` - object can be also created from a given Go map. Any map with string as a key and a compatible type as a value can be used,
```go
object := anytype.NewObjectFrom(map[string]int{
	"first":  1,
	"second": 2,
})
```

- `ParseObject(json string) (Object, error)` - loads an object from a JSON string,
```go
object, err := anytype.ParseObject(`{"first":1,"second":2}`)
if err != nil {
    ...
}
```

- `ParseFile(path string) (Object, error)` - loads an object from an UTF-8 encoded JSON file.
```go
object, err := anytype.ParseFile("file.json")
if err != nil {
    ...
}
```

### Manipulation With Fields
- `Set(values ...any) Object` - multiple new values can be set as key-value pairs, analogically to the constructor,
```go
object.Set(
    "first", 1,
    "second", 2, 
)
```

- `Unset(keys ...string) Object` - removes the given keys from the object,
```go
object.Unset("first", "second")
```

- `Clear() Object` - removes all keys from the object.
```go
object.Clear()
```

### Getting Fields
- Universal getter (requires type assertion),
```go
nested := object.Get("nested").(anytype.Object)
list := object.Get("list").(anytype.List)
str := object.Get("str").(string)
boolean := object.Get("boolean").(bool)
integer := object.Get("integer").(int)
float := object.Get("float").(float64)
```

- type-specific getters.
```go
nested := object.GetObject("nested")
list := object.GetList("list")
str := object.GetString("str")
boolean := object.GetBool("boolean")
integer := object.GetInt("integer")
float := object.GetFloat("float")
```

### Type Check
- `TypeOf(key string) Type`.
```go
if object.TypeOf("integer") == anytype.TypeInt {
    ...
}
```

### Export
- `String() string` - exports the object to a JSON string,
```go
fmt.Println(object.String())
```

- `FormatString(indent int) string` - exports the object to a well-arranged JSON string with the given indentation, 
```go
fmt.Println(object.FormatString(4))
```

- `Dict() map[string]any`
```go
var dict map[string]any
dict = object.Dict()
```

- `Keys() List` - exports all keys of the object to an AnyType list,
```go
var keys anytype.List
keys = object.Keys()
```

- `Values() List` - exports all values of the object to an AnyType list.
```go
var values anytype.List
values = object.Values()
```

### Features Over Whole Object
- `Clone() Object` - performs a deep copy of the object,
```go
copy := object.Clone()
```

- `Count() int` - returns a number of fileds of the object,
```go
for i := 0; i < object.Count(); i++ {
    ...
}
```

- `Empty() bool` - checks whether the object is empty (has 0 fields),
```go
if object.Empty() {
    ...
}
```

- `Equals(another Object) bool` - checks whether all fields of the object are equal to the fields in another object,
```go
if object.Equals(another) {
    ...
}
```

- `Merge(another Object) Object` - merges two objects together,
```go
merged := object.Merge(another)
```

- `Pluck(keys ...string) Object` - creates a new object containing only the selected keys from existing object,
```go
plucked := object.Pluck("first", "second")
```

- `Contains(elem any) bool` - checks whether the objects contains a value,
```go
if object.Contains(1) {
    ...
}
```

- `KeyOf(elem any) string` - returns any key containing the given value,
```go
first := object.KeyOf(1)
```

- `KeyExists(key string) bool` - checks whether a key exists in the object.
```go
if object.KeyExists("first") {
    ...
}
```

### ForEaches
- `ForEach(function func(string, any)) Object` - executes a given function over an every field of the object,
```go
object.ForEach(func(key string, value any) {
    ...
})
```

- `ForEachValue(function func(any)) Object` - ForEach without the key variable within the anonymous function,
```go
object.ForEachValue(func(value any) {
    ...
})
```

- type-specific ForEaches - anonymous function is only executed over values with the corresponding type.
```go
object.ForEachObject(func(object anytype.Object) {
    ...
})
object.ForEachList(func(list anytype.List) {
    ...
})
object.ForEachString(func(str string) {
    ...
})
object.ForEachBool(func(object bool) {
    ...
})
object.ForEachInt(func(integer int) {
    ...
})
object.ForEachFloat(func(float float64) {
    ...
})
```

### Mapping
- `Map(function func(string, any) any) Object` - returns a new object with fields modified by a given function,
```go
mapped := object.Map(func(key string, value any) any {
    ...
	return newValue
})
```

- `MapValues(function func(any) any) Object` - Map without the key variable within the anonymous function,
```go
mapped := object.MapValues(func(value any) any {
    ...
	return newValue
})
```

- type-specific Maps - selects only fields with the corresponding type.
```go
objects := object.MapObjects(func(object anytype.Object) any {
    ...
	return newValue
})
lists := object.MapLists(func(object anytype.List) any {
    ...
	return newValue
})
strs := object.MapStrings(func(object string) any {
    ...
	return newValue
})
integers := object.MapInts(func(object int) any {
    ...
	return newValue
})
floats := object.MapFloats(func(object float64) any {
    ...
	return newValue
})
```

### Async
- `ForEachAsync(function func(string, any)) Object` - performs the ForEach paralelly,
```go
object.ForEachAsync(func(key string, value any) {
    ...
})
```

- `MapAsync(function func(string, any) any) Object` - performs the Map paralelly.
```go
mapped := object.MapAsync(func(key string, value any) any {
    ...
	return newValue
})
```

### Tree Form
- `GetTF(tf string) any` - returns a value specified by the given tree form string,
```go
value := object.GetTF(".first#2")
```

- `SetTF(tf string, value any) Object` - sets a value on the path specified by the given tree form string.
```go
object.SetTF(".first#2", 2)
```

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
