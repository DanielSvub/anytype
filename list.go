/*
AnyType Library for Go
List (array) type
*/

package anytype

/*
List is an ordered sequence of elements.
*/
type List interface {
	field

	/*
		Init initializes the ego pointer, which allows deriving.

		Parameters:
		  - ptr - ego pointer.
	*/
	Init(ptr List)

	/*
		Ego acquires the ego pointer previously set by Init.

		Returns:
		  - ego pointer.
	*/
	Ego() List

	/*
		Add adds new elements at the end of the list.

		Parameters:
		  - values... - any amount of elements to add.

		Returns:
		  - updated list.
	*/
	Add(val ...any) List

	/*
		Insert inserts a new element at the specified position in the list.

		Parameters:
		  - index - position where the element should be inserted,
		  - value - element to insert.

		Returns:
		  - updated list.
	*/
	Insert(index int, value any) List

	/*
		Replace replaces an existing element with a new one.

		Parameters:
		  - index - position of the element which should be replaced,
		  - value - new element.

		Returns:
		  - updated list.
	*/
	Replace(index int, value any) List

	/*
		Delete deletes the elements at the specified positions in the list.

		Parameters:
		  - indexes... - any amount of positions of the elements to delete.

		Returns:
		  - updated list.
	*/
	Delete(index ...int) List

	/*
		Pop deletes the last element in the list.

		Returns:
		  - updated list.
	*/
	Pop() List

	/*
		Clear deletes all elements in the list.

		Returns:
		  - updated list.
	*/
	Clear() List

	/*
		Get acquires an element at the specified position in the list.

		Parameters:
		  - index - position of the element to get.

		Returns:
		  - corresponding value (any type, has to be asserted).
	*/
	Get(index int) any

	/*
		GetObject acquires an object at the specified position in the list.
		Causes a panic if the element has another type.

		Parameters:
		  - index - position of the element to get.

		Returns:
		  - corresponding value asserted as object.
	*/
	GetObject(index int) Object

	/*
		GetList acquires a list at the specified position in the list.
		Causes a panic if the element has another type.

		Parameters:
		  - index - position of the element to get.

		Returns:
		  - corresponding value asserted as list.
	*/
	GetList(index int) List

	/*
		GetString acquires a string at the specified position in the list.
		Causes a panic if the element has another type.

		Parameters:
		  - index - position of the element to get.

		Returns:
		  - corresponding value asserted as string.
	*/
	GetString(index int) string

	/*
		GetBool acquires a boolean at the specified position in the list.
		Causes a panic if the element has another type.

		Parameters:
		  - index - position of the element to get.

		Returns:
		  - corresponding value asserted as bool.
	*/
	GetBool(index int) bool

	/*
		GetInt acquires an integer at the specified position in the list.
		Causes a panic if the element has another type.

		Parameters:
		  - index - position of the element to get.

		Returns:
		  - corresponding value asserted as int.
	*/
	GetInt(index int) int

	/*
		GetFloat acquires a float at the specified position in the list.
		Causes a panic if the element has another type.

		Parameters:
		  - index - position of the element to get.

		Returns:
		  - corresponding value asserted as float64.
	*/
	GetFloat(index int) float64

	/*
		TypeOf gives a type of the element at the specified position in the list.
		If the index is out of range, 0 (TypeUndefined) is returned.

		Parameters:
		  - index - position of the element.

		Returns:
		  - integer constant representing the type (see type enum).
	*/
	TypeOf(index int) Type

	/*
		String gives a JSON representation of the list, including nested lists and objects.

		Returns:
		  - JSON string.
	*/
	String() string

	/*
		FormatString gives a JSON representation of the list in standardized format with the given indentation.

		Parameters:
		  - indent - indentation spaces (0-10).

		Returns:
		  - JSON string.
	*/
	FormatString(indent int) string

	/*
		Slice converts the list into a Go slice of any.

		Returns:
		  - slice.
	*/
	Slice() []any

	/*
		NativeSlice converts the list into a Go slice of any.
		All nested objects and lists are converted recursively.

		Returns:
		  - slice.
	*/
	NativeSlice() []any

	/*
		ObjectSlice converts the list of objects into a Go slice.
		Elements of other types are ignored.

		Returns:
		  - slice.
	*/
	ObjectSlice() []Object

	/*
		ListSlice converts the list of lists into a Go slice.
		Elements of other types are ignored.

		Returns:
		  - slice.
	*/
	ListSlice() []List

	/*
		StringSlice converts the list of strings into a Go slice.
		Elements of other types are ignored.

		Returns:
		  - slice.
	*/
	StringSlice() []string

	/*
		BoolSlice converts the list of bools into a Go slice.
		Elements of other types are ignored.

		Returns:
		  - slice.
	*/
	BoolSlice() []bool

	/*
		IntSlice converts the list of ints into a Go slice.
		Elements of other types are ignored.

		Returns:
		  - slice.
	*/
	IntSlice() []int

	/*
		FloatSlice converts the list of floats into a Go slice.
		Elements of other types are ignored.

		Returns:
		  - slice.
	*/
	FloatSlice() []float64

	/*
		Clone creates a deep copy of the list.

		Returns:
		  - copied list.
	*/
	Clone() List

	/*
		Count gives a number of elements in the list.

		Returns:
		  - number of elements.
	*/
	Count() int

	/*
		Empty checks whether the list is empty.

		Returns:
		  - true if the list is empty, false otherwise.
	*/
	Empty() bool

	/*
		Equals checks if the content of the list is equal to the content of another list.
		Nested objects and lists are compared recursively (by value).

		Parameters:
		  - another - a list to compare with.

		Returns:
		  - true if the lists are equal, false otherwise.
	*/
	Equals(another List) bool

	/*
		Concat creates a new list containing all elements of the old list and another list.
		The old list remains unchanged.

		Parameters:
		  - another - a list to concat.

		Returns:
		  - new list.
	*/
	Concat(another List) List

	/*
		SubList creates a new list containing the elements from the starting index (including) to the ending index (excluding).
		If the ending index is zero, it is set to the length of the list. If negative, it is counted from the end of the list.
		Starting index has to be non-negative and cannot be higher than the ending index.

		Parameters:
		  - start - starting index,
		  - end - ending index.

		Returns:
		  - created sub list.
	*/
	SubList(start int, end int) List

	/*
		Contains checks if the list contains a given element.
		Objects and lists are compared by reference.

		Parameters:
		  - elem - the element to check.

		Returns:
		  - true if the list contains the element, false otherwise.
	*/
	Contains(elem any) bool

	/*
		IndexOf gives a position of the first occurrence of a given element.

		Parameters:
		  - elem - the element to check.

		Returns:
		  - index of the element (-1 if the list does not contain the element).
	*/
	IndexOf(elem any) int

	/*
		Sort sorts elements in the list (ascending).
		The first element of the list determines the sorting type and has to be either string, int or float.
		Values of other types are ignored.

		Returns:
		  - updated list.
	*/
	Sort() List

	/*
		Reverse reverses the order of elements in the list.

		Returns:
		  - updated list.
	*/
	Reverse() List

	/*
		AllObjects checks if the list is homogeneous and all of its elements are objects.

		Returns:
		  - true if all elements are objects, false otherwise.
	*/
	AllObjects() bool

	/*
		AllLists checks if the list is homogeneous and all of its elements are lists.

		Returns:
		  - true if all elements are lists, false otherwise.
	*/
	AllLists() bool

	/*
		AllStrings checks if the list is homogeneous and all of its elements are strings.

		Returns:
		  - true if all elements are strings, false otherwise.
	*/
	AllStrings() bool

	/*
		AllBools checks if the list is homogeneous and all of its elements are bools.

		Returns:
		  - true if all elements are bools, false otherwise.
	*/
	AllBools() bool

	/*
		AllInts checks if the list is homogeneous and all of its elements are ints.

		Returns:
		  - true if all elements are ints, false otherwise.
	*/
	AllInts() bool

	/*
		AllFloats checks if the list is homogeneous and all of its elements are floats.

		Returns:
		  - true if all elements are floats, false otherwise.
	*/
	AllFloats() bool

	/*
		AllNumeric checks if all elements of the list are numeric (ints or floats).

		Returns:
		  - true if all elements are numeric, false otherwise.
	*/
	AllNumeric() bool

	/*
		ForEach executes a given function over an every element of the list.
		The function has two parameters: index of the current element and its value.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEach(function func(i int, val any)) List

	/*
		ForEachValue executes a given function over an every element of the list.
		The function has one parameter, value of the current element.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachValue(function func(x any)) List

	/*
		ForEachObject executes a given function over all objects in the list.
		Elements with other types are ignored.
		The function has one parameter, the current object.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachObject(function func(x Object)) List

	/*
		ForEachList executes a given function over all lists nested in the list.
		Elements with other types are ignored.
		The function has one parameter, the current list.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachList(function func(x List)) List

	/*
		ForEachString executes a given function over all strings in the list.
		Elements with other types are ignored.
		The function has one parameter, the current string.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachString(function func(x string)) List

	/*
		ForEachBool executes a given function over all bools in the list.
		Elements with other types are ignored.
		The function has one parameter, the current bool.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachBool(function func(x bool)) List

	/*
		ForEachInt executes a given function over all ints in the list.
		Elements with other types are ignored.
		The function has one parameter, the current int.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachInt(function func(x int)) List

	/*
		ForEachFloat executes a given function over all floats in the list.
		Elements with other types are ignored.
		The function has one parameter, the current float.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachFloat(function func(x float64)) List

	/*
		Map copies the list and modifies each element by a given mapping function.
		The resulting element can have a different type than the original one.
		The function has two parameters: current index and value of the current element. Returns any.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new list.
	*/
	Map(function func(i int, val any) any) List

	/*
		MapValues copies the list and modifies each element by a given mapping function.
		The resulting element can have a different type than the original one.
		The function has one parameter, value of the current element, and returns any.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new list.
	*/
	MapValues(function func(x any) any) List

	/*
		MapObjects selects all objects from the list and modifies each of them by a given mapping function.
		Elements of other types are ignored.
		The resulting element can have a different type than the original one.
		The function has one parameter, the current object, and returns any.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new list.
	*/
	MapObjects(function func(x Object) any) List

	/*
		MapLists selects all nested lists from the list and modifies each of them by a given mapping function.
		Elements of other types are ignored.
		The resulting element can have a different type than the original one.
		The function has one parameter, the current list, and returns any.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new list.
	*/
	MapLists(function func(x List) any) List

	/*
		MapStrings selects all strings from the list and modifies each of them by a given mapping function.
		Elements of other types are ignored.
		The resulting element can have a different type than the original one.
		The function has one parameter, the current string, and returns any.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new list.
	*/
	MapStrings(function func(x string) any) List

	/*
		MapBools selects all bools from the list and modifies each of them by a given mapping function.
		Elements of other types are ignored.
		The resulting element can have a different type than the original one.
		The function has one parameter, the current bool, and returns any.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new list.
	*/
	MapBools(function func(x bool) any) List

	/*
		MapInts selects all ints from the list and modifies each of them by a given mapping function.
		Elements of other types are ignored.
		The resulting element can have a different type than the original one.
		The function has one parameter, the current int, and returns any.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new list.
	*/
	MapInts(function func(x int) any) List

	/*
		MapFloats selects all floats from the list and modifies each of them by a given mapping function.
		Elements of other types are ignored.
		The resulting element can have a different type than the original one.
		The function has one parameter, the current float, and returns any.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new list.
	*/
	MapFloats(function func(x float64) any) List

	/*
		Reduce reduces all elements of the list into a single value.
		The function has two parameters: value returned by the previous iteration and value of the current element. Returns any.
		The list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - computed value.
	*/
	Reduce(initial any, function func(acc any, val any) any) any

	/*
		ReduceStrings reduces all strings in the list into a single string.
		Elements of other types are ignored.
		The function has two parameters: string returned by the previous iteration and current string. Returns string.
		The list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - computed value.
	*/
	ReduceStrings(initial string, function func(acc string, val string) string) string

	/*
		ReduceInts reduces all ints in the list into a single int.
		Elements of other types are ignored.
		The function has two parameters: int returned by the previous iteration and current int. Returns int.
		The list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - computed value.
	*/
	ReduceInts(initial int, function func(acc int, val int) int) int

	/*
		ReduceFloats reduces all floats in the list into a single float.
		Elements of other types are ignored.
		The function has two parameters: float returned by the previous iteration and current float. Returns float.
		The list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - computed value.
	*/
	ReduceFloats(initial float64, function func(acc float64, val float64) float64) float64

	/*
		Filter creates a new list containing elements of the old one satisfying a condition.
		The function has one parameter, value of the current element, and returns bool.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - filtered list.
	*/
	Filter(function func(x any) bool) List

	/*
		FilterObjects creates a new list containing objects of the old one satisfying a condition.
		Elements of other types are ignored.
		The function has one parameter, current object, and returns bool.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - filtered list.
	*/
	FilterObjects(function func(x Object) bool) List

	/*
		FilterLists creates a new list containing nested lists of the old one satisfying a condition.
		Elements of other types are ignored.
		The function has one parameter, current list, and returns bool.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - filtered list.
	*/
	FilterLists(function func(x List) bool) List

	/*
		FilterStrings creates a new list containing strings of the old one satisfying a condition.
		Elements of other types are ignored.
		The function has one parameter, current string, and returns bool.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - filtered list.
	*/
	FilterStrings(function func(x string) bool) List

	/*
		FilterInts creates a new list containing ints of the old one satisfying a condition.
		Elements of other types are ignored.
		The function has one parameter, current int, and returns bool.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - filtered list.
	*/
	FilterInts(function func(x int) bool) List

	/*
		FilterFloats creates a new list containing floats of the old one satisfying a condition.
		Elements of other types are ignored.
		The function has one parameter, current float, and returns bool.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - filtered list.
	*/
	FilterFloats(function func(x float64) bool) List

	/*
		IntSum computes a sum of all elements in the list.
		All elements of the list have to be ints.

		Returns:
		  - computed sum (int).
	*/
	IntSum() int

	/*
		Sum computes a sum of all elements in the list.
		All elements of the list have to be numeric.

		Returns:
		  - computed sum (float).
	*/
	Sum() float64

	/*
		IntProd computes a product of all elements in the list.
		All elements of the list have to be ints.

		Returns:
		  - computed product (int).
	*/
	IntProd() int

	/*
		Prod computes a product of all elements in the list.
		All elements of the list have to be numeric.

		Returns:
		  - computed pruduct (float).
	*/
	Prod() float64

	/*
		Avg computes an arithmetic mean of all elements in the list.
		All elements of the list have to be numeric.

		Returns:
		  - computed average value (float).
	*/
	Avg() float64

	/*
		IntMin finds a minimum int of the list.

		Returns:
		  - found minimum (int).
	*/
	IntMin() int

	/*
		Min finds a minimum number of the list.

		Returns:
		  - found minimum (float).
	*/
	Min() float64

	/*
		IntMax finds a maximum int of the list.

		Returns:
		  - found maximum (int).
	*/
	IntMax() int

	/*
		Max finds a maximum number of the list.

		Returns:
		  - found maximum (float).
	*/
	Max() float64

	/*
		ForEachAsync parallelly executes a given function over an every element of the list.
		The function has two parameters: index of the current element and its value.
		The order of the iterations is random.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachAsync(function func(i int, val any)) List

	/*
		MapAsync copies the list and paralelly modifies each element by a given mapping function.
		The resulting element can have a different type than the original one.
		The function has two parameters: index of the current element and its value.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new list.
	*/
	MapAsync(function func(i int, val any) any) List

	/*
		GetTF acquires a value specified by a given tree form.

		Parameters:
		  - tf - tree form string.

		Returns:
		  - corresponding value (any type, has to be asserted).
	*/
	GetTF(tf string) any

	/*
		SetTF sets a value specified by a given tree form.
		If the index exceeds the count, the interspace will be filled with nils.

		Parameters:
		  - tf - tree form string,
		  - value - value to set.

		Returns:
		  - updated list.
	*/
	SetTF(tf string, value any) List

	/*
		Unset deletes the value specified by a given tree form.
		If the TF path does not exist, nothing happens.

		Parameters:
		  - tf - tree form string.

		Returns:
		  - updated object.
	*/
	UnsetTF(tf string) List

	/*
		TypeOfTF gives a type of the element specified by a given tree form.
		If the TF path does not exist, 0 (TypeUndefined) is returned.

		Parameters:
		  - tf - tree form string.

		Returns:
		  - integer constant representing the type (see type enum).
	*/
	TypeOfTF(tf string) Type
}
