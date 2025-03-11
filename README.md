# AnyType

AnyType is a Go library providing dynamic data structures with JSON support. It contains a number of advanced features with API inspired by Java collections.

It supports following data types compatible with JSON standard:
- nil (null)
- object
- list (array)
- string
- boolean
- integer (number)
- float (number)

Types can be referenced by the `Type` enum (e.g. `TypeNil`, `TypeObject`, ...). If the value does not exist, its type is considered `TypeUndefined`. Attempting to access an undefined value will cause a panic.

AnyType also allows usage of so-called "tree form" for accessing values. It is a string using hash for list elements and dot for object fields. For example `#1.a.b#4` or `.d.c#5#0`.

The library is tested with 100% coverage.

## Objects

Object is an unordered set of key-value pairs, the keys are of type string. The default implementation is based on built-in Go maps. It is possible to make custom implementations by implementing the `Object` interface.

### Constructors

- `NewObject(values ...any) Object` - initial object values are specified as key value pairs. The function panics if an odd number of arguments is given,
```go
emptyObject := anytype.NewObject()
object := anytype.NewObject(
    "number", 1,
    "string", "test",
    "bool", true,
    "null", nil,
)
```

- `NewObjectFrom(dict any) Object` - object can be also created from a given Go map. Any map with string for keys and a compatible type (including any) for values can be used,
```go
object := anytype.NewObjectFrom(map[string]int{
	"first": 1,
	"second": 2,
    "third": 3,
})
```

- `ParseObject(json string) (Object, error)` - loads an object from a JSON string,
```go
object, err := anytype.ParseObject(`{"first":1,"second":2,"third":3}`)
if err != nil {
    // ...
}
```

- `ParseFile(path string) (Object, error)` - loads an object from an UTF-8 encoded JSON file.
```go
object, err := anytype.ParseFile("file.json")
if err != nil {
    // ...
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
    // ...
}
```

### Export
- `String() string` - exports the object into a JSON string,
```go
fmt.Println(object.String())
```

- `FormatString(indent int) string` - exports the object into a well-arranged JSON string with the given indentation, 
```go
fmt.Println(object.FormatString(4))
```

- `Dict() map[string]any` - exports the object into a Go map,
```go
var dict map[string]any
dict = object.Dict()
```

- `Keys() List` - exports all keys of the object into an AnyType list,
```go
var keys anytype.List
keys = object.Keys()
```

- `Values() List` - exports all values of the object into an AnyType list.
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
    // ...
}
```

- `Empty() bool` - checks whether the object is empty,
```go
if object.Empty() {
    // ...
}
```

- `Equals(another Object) bool` - checks whether all fields of the object are equal to the fields of another object,
```go
if object.Equals(another) {
    // ...
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

- `Contains(elem any) bool` - checks whether the object contains a certain value,
```go
if object.Contains(1) {
    // ...
}
```

- `KeyOf(elem any) string` - returns any key containing the given value. It panics if the object does not contain the value,
```go
first := object.KeyOf(1)
```

- `KeyExists(key string) bool` - checks whether a key exists within the object.
```go
if object.KeyExists("first") {
    // ...
}
```

### ForEaches
- `ForEach(function func(string, any)) Object` - executes a given function over an every field of the object,
```go
object.ForEach(func(key string, value any) {
    // ...
})
```

- `ForEachValue(function func(any)) Object` - ForEach without the key variable within the anonymous function,
```go
object.ForEachValue(func(value any) {
    // ...
})
```

- type-specific ForEaches - anonymous function is only executed over values of the corresponding type.
```go
object.ForEachObject(func(object anytype.Object) {
    // ...
})
object.ForEachList(func(list anytype.List) {
    // ...
})
object.ForEachString(func(str string) {
    // ...
})
object.ForEachBool(func(object bool) {
    // ...
})
object.ForEachInt(func(integer int) {
    // ...
})
object.ForEachFloat(func(float float64) {
    // ...
})
```

### Mappings
- `Map(function func(string, any) any) Object` - returns a new object with fields modified by a given function,
```go
mapped := object.Map(func(key string, value any) any {
    // ...
	return newValue
})
```

- `MapValues(function func(any) any) Object` - Map without the key variable within the anonymous function,
```go
mapped := object.MapValues(func(value any) any {
    // ...
	return newValue
})
```

- type-specific Maps - selects only fields of the corresponding type.
```go
objects := object.MapObjects(func(object anytype.Object) any {
    // ...
	return newValue
})
lists := object.MapLists(func(list anytype.List) any {
    // ...
	return newValue
})
strs := object.MapStrings(func(str string) any {
    // ...
	return newValue
})
booleans := object.MapBools(func(boolean bool) any {
    // ...
	return newValue
})
integers := object.MapInts(func(integer int) any {
    // ...
	return newValue
})
floats := object.MapFloats(func(float float64) any {
    // ...
	return newValue
})
```

### Asynchronous
- `ForEachAsync(function func(string, any)) Object` - performs the ForEach parallelly,
```go
object.ForEachAsync(func(key string, value any) {
    // ...
})
```

- `MapAsync(function func(string, any) any) Object` - performs the Map parallelly.
```go
mapped := object.MapAsync(func(key string, value any) any {
    // ...
	return newValue
})
```

### Tree Form
- `GetTF(tf string) any` - returns a value specified by the given tree form string,
```go
value := object.GetTF(".first#2")
```

- `SetTF(tf string, value any) Object` - sets a value on the path specified by the given tree form string,
```go
object.SetTF(".first#2", 2)
```

- `UnsetTF(tf string) Object` - unsets a value on the path specified by the given tree form string,
```go
object.UnsetTF(".first#2")
```

- `TypeOfTF(tf string) Type` - returns a type of the field specified by the given tree form string.
```go
if object.TypeOfTF(".first#2") == anytype.TypeInt {
    // ...
}
```

## Lists

List is an ordered sequence of elements. The default implementation is based on built-in Go slices. It is possible to make custom implementations by implementing the `List` interface.

### Constructors

- `NewList(values ...any) List` - initial list elements could be given as variadic arguments,
```go
emptyList := anytype.NewList()
list := anytype.NewList(1, "test", true, nil)
```

- `NewListOf(value any, count int) List` - creates a list of n repeated values,
```go
list := anytype.NewListOf(nil, 10)
```

- `NewListFrom(slice any) List` - creates a list from a given Go slice. Any slice of a compatible type (including any) can be used,
```go
list := anytype.NewListFrom([]int{1, 2, 3})
```

- `ParseList(json string) (List, error)` - loads a list from a JSON string,
```go
list, err := anytype.ParseList(`[1, 2, 3]`)
if err != nil {
    // ...
}
```

### Manipulation With Elements
- `Add(val ...any) List` - adds any amount of new elements to the list,
```go
list.Add(1, 2, 3)
```

- `Insert(index int, value any) List` - inserts a new element to a specific position in the list,
```go
list.Insert(1, 1.5)
```

- `Replace(index int, value any) List` - replaces an existing element,
```go
list.Replace(1, "2")
```

- `Delete(index ...int) List` - removes specified elements,
```go
list.Delete(1, 2)
```

- `Pop() List` - removes the last element from the list,
```go
list.Pop()
```

- `Clear() List` - removes all elements from the list.
```go
list.Clear()
```

### Getting Elements
- Universal getter (requires type assertion),
```go
object := list.Get(0).(anytype.Object)
nested := list.Get(1).(anytype.List)
str := list.Get(2).(string)
boolean := list.Get(3).(bool)
integer := list.Get(4).(int)
float := list.Get(5).(float64)
```
- type-specific getters.
```go
object := list.GetObject(0)
nested := list.GetList(1)
str := list.GetString(2)
boolean := list.GetBool(3)
integer := list.GetInt(4)
float := list.GetFloat(5)
```

### Type Check
- `TypeOf(index int) Type`.
```go
if list.TypeOf(0) == anytype.TypeInt {
    // ...
}
```

### Export
- `String() string` - exports the list into a JSON string,
```go
fmt.Println(list.String())
```

- `FormatString(indent int) string` - exports the list into a well-arranged JSON string with the given indentation, 
```go
fmt.Println(list.FormatString(4))
```

- `Slice() []any` - exports the list into a Go slice,
```go
var slice []any
slice = list.Slice()
```

- export into type-specific slices - values of other types are ignored,
```go
var objects []anytype.Object
objects = list.ObjectSlice()
var lists []anytype.List
lists = list.ListSlice()
var strs []string
strs = list.StringSlice()
var bools []bool
bools = list.BoolSlice()
var ints []int
ints = list.IntSlice()
var floats []float64
floats = list.FloatSlice()
```

### Features Over Whole List
- `Clone() List` - performs a deep copy of the list,
```go
copy := list.Clone()
```

- `Count() int` - returns a number of elements in the list,
```go
for i := 0; i < list.Count(); i++ {
    // ...
}
```

- `Empty() bool` - checks whether the list is empty,
```go
if list.Empty() {
    // ...
}
```

- `Equals(another List) bool` - checks whether all elements of the list are equal to the elements of another list,
```go
if list.Equals(another) {
    // ...
}
```

- `Concat(another List) List` - concates two lists together,
```go
concated := list.Concat(another)
```

- `SubList(start int, end int) List` - cuts a part of the list,
```go
subList := list.SubList(1, 3)
```

- `Contains(elem any) bool` - checks whether the list contains a certain value,
```go
if list.Contains("value") {
    // ...
}
```

- `IndexOf(elem any) int` - returns a position of the first occurrence of the given value,
```go
elem := list.IndexOf("value")
```

- `Sort() List` - sorts the elements in the list. The sorting type is determined by the first element which has to be either string, int or float,
```go
list.Sort()
```

- `Reverse() List` - reverses the list,
```go
list.Reverse()
```

### Checks For Homogeneity
- `AllNumeric() bool` - checks if all elements are numbers (ints or floats),
```go
if list.AllNumeric() {
    // ...
}
```

- type-specific asserts.
```go
if list.AllObjects() {
    // ...
} else if list.AllLists() {
    // ...
} else if list.AllStrings() {
    // ...
} else if list.AllBools() {
    // ...
} else if list.AllInts() {
    // ...
} else if list.AllFloats() {
    // ...
} else {
    // ...
}
```

### ForEaches
- `ForEach(function func(int, any)) List` - executes a given function over an every element of the list,
```go
list.ForEach(func(index int, value any) {
    // ...
})
```

- `ForEachValue(function func(any)) List` - ForEach without the index variable within the anonymous function,
```go
list.ForEachValue(func(value any) {
    // ...
})
```

- type-specific ForEaches - anonymous function is only executed over values of the corresponding type.
```go
list.ForEachObject(func(object anytype.Object) {
    // ...
})
list.ForEachList(func(list anytype.List) {
    // ...
})
list.ForEachString(func(str string) {
    // ...
})
list.ForEachBool(func(object bool) {
    // ...
})
list.ForEachInt(func(integer int) {
    // ...
})
list.ForEachFloat(func(float float64) {
    // ...
})
```

### Mappings
- `Map(function func(int, any) any) List` - returns a new list with elements modified by a given function,
```go
mapped := list.Map(func(index int, value any) any {
    // ...
	return newValue
})
```

- `MapValues(function func(any) any) List` - Map without the index variable within the anonymous function,
```go
mapped := list.MapValues(func(value any) any {
    // ...
	return newValue
})
```

- type-specific Maps - selects only elements of the corresponding type.
```go
objects := list.MapObjects(func(object anytype.Object) any {
    // ...
	return newValue
})
lists := list.MapLists(func(list anytype.List) any {
    // ...
	return newValue
})
strs := list.MapStrings(func(str string) any {
    // ...
	return newValue
})
booleans := list.MapBools(func(boolean bool) any {
    // ...
	return newValue
})
integers := list.MapInts(func(integer int) any {
    // ...
	return newValue
})
floats := list.MapFloats(func(float float64) any {
    // ...
	return newValue
})
```

### Reductions
- `Reduce(initial any, function func(any, any) any) any` - reduces all elements in the list into a single value,
```go
result := list.Reduce(0, func(sum, value any) any {
	return sum.(int) + value.(int)
})
```

- type-specific Reductions - selects only elements of the corresponding type. Return value has to be of the same type.
```go
result := list.ReduceStrings("", func(concated, value string) string {
	return concated + value
})
result := list.ReduceInts(0, func(sum, value int) int {
	return sum + value
})
result := list.ReduceFloats(0, func(sum, value float64) float64 {
	return sum + value
})
```

### Filters
- `Filter(function func(any) bool) List` - filters elements in the list based on a condition,
```go
filtered := list.Filter(func(value any) bool {
    // ...
	return condition
})
```

- type-specific Filters - filters only elements of the corresponding type.
```go
objects := list.FilterObjects(func(value anytype.Object) bool {
    // ...
	return condition
})
lists := list.FilterLists(func(value anytype.List) bool {
    // ...
	return condition
})
strs := list.FilterStrings(func(value string) bool {
    // ...
	return condition
})
integers := list.FilterInts(func(value int) bool {
    // ...
	return condition
})
floats := list.FilterFloats(func(value float64) bool {
    // ...
	return condition
})
```

### Numeric Operations
- `IntSum() int` - computes a sum of all ints in the list (0 if no ints are present),
```go
sum := list.IntSum()
```

- `Sum() float64` - compatible with both numeric types, returns float,
```go
sum := list.Sum()
```
- `IntProd() int` - computes a product of all ints in the list (1 if no ints are present),
```go
product := list.IntProd()
```

- `Prod() float64` - compatible with both numeric types, returns float,
```go
product := list.Prod()
```

- `Avg() float64` - computes an arithmetic mean of all numbers in the list (0 if no ints are present),
```go
average := list.Avg()
```

- `IntMin() int` - returns a minimum integer value in the list (0 if no ints are present),
```go
minimum := list.IntMin()
```

- `Min() float64` - compatible with both numeric types, returns float,
```go
minimum := list.Min()
```

- `IntMax() int` - returns a maximum integer value in the list (0 if no ints are present),
```go
maximum := list.IntMax()
```

- `Max() float64` - compatible with both numeric types, returns float,
```go
maximum := list.Max()
```

### Asynchronous
- `ForEachAsync(function func(int, any)) List` - performs the ForEach parallelly,
```go
list.ForEachAsync(func(index int, value any) {
    // ...
})
```

- `MapAsync(function func(int, any) any) List` - performs the Map parallelly.
```go
mapped := list.MapAsync(func(index int, value any) any {
    // ...
	return newValue
})
```

### Tree Form
- `GetTF(tf string) any` - returns a value specified by the given tree form string,
```go
value := list.GetTF("#2.first")
```

- `SetTF(tf string, value any) List` - sets a value on the path specified by the given tree form string,
```go
list.SetTF("#2.first", 2)
```

- `UnsetTF(tf string) List` - unsets a value on the path specified by the given tree form string,
```go
list.UnsetTF("#2.first")
```

- `TypeOfTF(tf string) Type` - returns a type of the element specified by the given tree form string.
```go
if list.TypeOfTF("#2.first") == anytype.TypeInt {
    // ...
}
```

## Derived Structures
AnyType supports inheritance and method overriding by defining custom structures with embedded object or list. As Go uses the embedded pointer as a receiver instead of the embedding structure, the pointer to the derived structure (so-called "ego pointer") has to be stored using the method `Init(ptr Object)`/`Init(ptr List)`. When overriding a method, the ego pointer can be obtained with `Ego() Object`/`Ego() List`.

```go
// Embeds an object
type Animal struct {
	anytype.Object
	name string
}

func NewAnimal(name string, age int) *Animal {
	ego := &Animal{
		Object: anytype.NewObject(
			"age", age,
		),
		name: name,
	}
	ego.Init(ego) // Ego pointer initialization
	return ego
}

// Method overriding
func (ego *Animal) Clear() anytype.Object {
	fmt.Fprintln(os.Stderr, "fields of an animal cannot be cleared")
	return ego.Ego() // Using stored pointer to return Animal instead of the embedded object
}

func (ego *Animal) Breathe() {
	fmt.Println("breathing")
}

// Inherits from animal, adds another field
type Dog struct {
	*Animal
	breed string
}

func NewDog(name string, age int, breed string) *Dog {
	ego := &Dog{
		Animal: NewAnimal(name, age),
		breed:  breed,
	}
	ego.Init(ego) // Ego pointer initialization
	return ego
}

// Method overriding
func (ego *Dog) Unset(keys ...string) anytype.Object {
	fmt.Fprintln(os.Stderr, "fields of a dog cannot be unset")
	return ego.Ego() // Using stored pointer to return Dog instead of object
}

func (ego *Dog) Bark() {
	fmt.Println("woof")
}

func main() {

	dog := NewDog("Rex", 2, "German Shepherd")

	// Methods from both Dog and Animal can be used
	dog.Breathe()
	dog.Bark()

    // Methods of the object can be used, too
	dog.Set("color", "black")

	// Printing the object inside
	fmt.Println(dog.String())

	// Both methods have been overridden
	dog.Unset("age")
	dog.Clear()

}
```
