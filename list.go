/*
AnyType Library for Go
List (array) type
*/

package anytype

/*
Interface for a list.

Extends:
  - field.
*/
type List interface {
	field

	/*
		Initializes the ego pointer, which allows deriving.

		Parameters:
		  - ptr - ego pointer.
	*/
	Init(ptr List)

	/*
		Acquires the ego pointer previously set by Init.

		Returns:
		  - ego pointer.
	*/
	Ego() List

	/*
		Adds new elements at the end of the list.

		Parameters:
		  - values... - any amount of elements to add.

		Returns:
		  - updated list.
	*/
	Add(val ...any) List

	/*
		Inserts a new element at the specified position in the list.

		Parameters:
		  - index - position where the element should be inserted,
		  - value - element to insert.

		Returns:
		  - updated list.
	*/
	Insert(index int, value any) List

	/*
		Replaces an existing element with a new one.

		Parameters:
		  - index - position of the element which should be replaced,
		  - value - new element.

		Returns:
		  - updated list.
	*/
	Replace(index int, value any) List

	/*
		Deletes the elements at the specified positions in the list.

		Parameters:
		  - indexes... - any amount of positions of the elements to delete.

		Returns:
		  - updated list.
	*/
	Delete(index ...int) List

	/*
		Deletes the last element in the list.

		Returns:
		  - updated list.
	*/
	Pop() List

	/*
		Deletes all elements in the list.

		Returns:
		  - updated list.
	*/
	Clear() List

	/*
		Acquires the element at the specified position in the list.

		Parameters:
		  - index - position of the element to get.

		Returns:
		  - corresponding value (any type, has to be asserted).
	*/
	Get(index int) any

	/*
		Acquires the object at the specified position in the list.
		Causes a panic if the element has another type.

		Parameters:
		  - index - position of the element to get.

		Returns:
		  - corresponding value asserted as object.
	*/
	GetObject(index int) Object

	/*
		Acquires the list at the specified position in the list.
		Causes a panic if the element has another type.

		Parameters:
		  - index - position of the element to get.

		Returns:
		  - corresponding value asserted as list.
	*/
	GetList(index int) List

	/*
		Acquires the string at the specified position in the list.
		Causes a panic if the element has another type.

		Parameters:
		  - index - position of the element to get.

		Returns:
		  - corresponding value asserted as string.
	*/
	GetString(index int) string

	/*
		Acquires the boolean at the specified position in the list.
		Causes a panic if the element has another type.

		Parameters:
		  - index - position of the element to get.

		Returns:
		  - corresponding value asserted as bool.
	*/
	GetBool(index int) bool

	/*
		Acquires the integer at the specified position in the list.
		Causes a panic if the element has another type.

		Parameters:
		  - index - position of the element to get.

		Returns:
		  - corresponding value asserted as int.
	*/
	GetInt(index int) int

	/*
		Acquires the float at the specified position in the list.
		Causes a panic if the element has another type.

		Parameters:
		  - index - position of the element to get.

		Returns:
		  - corresponding value asserted as float64.
	*/
	GetFloat(index int) float64

	/*
		Gives a type of the element at the specified position in the list.

		Parameters:
		  - index - position of the element.

		Returns:
		  - integer constant representing the type (see type enum).
	*/
	TypeOf(index int) Type

	/*
		Gives a JSON representation of the list, including nested lists and objects.

		Returns:
		  - JSON string.
	*/
	String() string

	/*
		Gives a JSON representation of the list in standardized format with the given indentation.

		Parameters:
		  - indent - indentation spaces (0-10).

		Returns:
		  - JSON string.
	*/
	FormatString(indent int) string

	/*
		Converts the list into a Go slice of empty interfaces.

		Returns:
		  - slice.
	*/
	Slice() []any

	/*
		Converts the list of objects into a Go slice.
		The list has to be homogeneous and all elements have to be objects.

		Returns:
		  - slice.
	*/
	ObjectSlice() []Object

	/*
		Converts the list of lists into a Go slice.
		The list has to be homogeneous and all elements have to be lists.

		Returns:
		  - slice.
	*/
	ListSlice() []List

	/*
		Converts the list of strings into a Go slice.
		The list has to be homogeneous and all elements have to be strings.

		Returns:
		  - slice.
	*/
	StringSlice() []string

	/*
		Converts the list of bools into a Go slice.
		The list has to be homogeneous and all elements have to be bools.

		Returns:
		  - slice.
	*/
	BoolSlice() []bool

	/*
		Converts the list of ints into a Go slice.
		The list has to be homogeneous and all elements have to be ints.

		Returns:
		  - slice.
	*/
	IntSlice() []int

	/*
		Converts the list of floats into a Go slice.
		The list has to be homogeneous and all elements have to be floats.

		Returns:
		  - slice.
	*/
	FloatSlice() []float64

	/*
		Creates a deep copy of the list.

		Returns:
		  - copied list.
	*/
	Clone() List

	/*
		Gives a number of elements in the list.

		Returns:
		  - number of elements.
	*/
	Count() int

	/*
		Checks whether the list is empty.

		Returns:
		  - true if the list is empty, false otherwise.
	*/
	Empty() bool

	/*
		Checks if the content of the list is equal to the content of another list.
		Nested objects and lists are compared recursively (by value).

		Parameters:
		  - another - a list to compare with.

		Returns:
		  - true if the lists are equal, false otherwise.
	*/
	Equals(another List) bool

	/*
		Creates a new list containing all elements of the old list and another list.
		The old list remains unchanged.

		Parameters:
		  - another - a list to append.

		Returns:
		  - new list.
	*/
	Concat(another List) List

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
	SubList(start int, end int) List

	/*
		Checks if the list contains a given element.
		Objects and lists are compared by reference.

		Parameters:
		  - elem - the element to check.

		Returns:
		  - true if the list contains the element, false otherwise.
	*/
	Contains(elem any) bool

	/*
		Gives a position of the first occurrence of a given element.

		Parameters:
		  - elem - the element to check.

		Returns:
		  - index of the element (-1 if the list does not contain the element).
	*/
	IndexOf(elem any) int

	/*
		Sorts elements in the list (ascending).
		The list has to be homogeneous, all elements have to be either strings, ints or floats.

		Returns:
		  - updated list.
	*/
	Sort() List

	/*
		Reverses the order of elements in the list.

		Returns:
		  - updated list.
	*/
	Reverse() List

	/*
		Checks if the list is homogeneous and all of its elements are objects.

		Returns:
		  - true if all elements are objects, false otherwise.
	*/
	AllObjects() bool

	/*
		Checks if the list is homogeneous and all of its elements are lists.

		Returns:
		  - true if all elements are lists, false otherwise.
	*/
	AllLists() bool

	/*
		Checks if the list is homogeneous and all of its elements are strings.

		Returns:
		  - true if all elements are strings, false otherwise.
	*/
	AllStrings() bool

	/*
		Checks if the list is homogeneous and all of its elements are bools.

		Returns:
		  - true if all elements are bools, false otherwise.
	*/
	AllBools() bool

	/*
		Checks if the list is homogeneous and all of its elements are ints.

		Returns:
		  - true if all elements are ints, false otherwise.
	*/
	AllInts() bool

	/*
		Checks if the list is homogeneous and all of its elements are floats.

		Returns:
		  - true if all elements are floats, false otherwise.
	*/
	AllFloats() bool

	/*
		Checks if all elements of the list are numeric (ints or floats).

		Returns:
		  - true if all elements are numeric, false otherwise.
	*/
	AllNumeric() bool

	/*
		Executes a given function over an every element of the list.
		The function has two parameters: index of the current element and its value.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEach(function func(int, any)) List

	/*
		Executes a given function over an every element of the list.
		The function has one parameter, value of the current element.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachValue(function func(any)) List

	/*
		Executes a given function over all objects in the list.
		Elements with other types are ignored.
		The function has one parameter, the current object.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachObject(function func(Object)) List

	/*
		Executes a given function over all lists nested in the list.
		Elements with other types are ignored.
		The function has one parameter, the current list.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachList(function func(List)) List

	/*
		Executes a given function over all strings in the list.
		Elements with other types are ignored.
		The function has one parameter, the current string.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachString(function func(string)) List

	/*
		Executes a given function over all bools in the list.
		Elements with other types are ignored.
		The function has one parameter, the current bool.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachBool(function func(bool)) List

	/*
		Executes a given function over all ints in the list.
		Elements with other types are ignored.
		The function has one parameter, the current int.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachInt(function func(int)) List

	/*
		Executes a given function over all floats in the list.
		Elements with other types are ignored.
		The function has one parameter, the current float.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachFloat(function func(float64)) List

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
	Map(function func(int, any) any) List

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
	MapValues(function func(any) any) List

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
	MapObjects(function func(Object) any) List

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
	MapLists(function func(List) any) List

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
	MapStrings(function func(string) any) List

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
	MapInts(function func(int) any) List

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
	MapFloats(function func(float64) any) List

	/*
		Reduces all elements of the list into a single value.
		The function has two parameters: value returned by the previous iteration and value of the current element. Returns empty interface.
		The list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - computed value.
	*/
	Reduce(initial any, function func(any, any) any) any

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
	ReduceStrings(initial string, function func(string, string) string) string

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
	ReduceInts(initial int, function func(int, int) int) int

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
	ReduceFloats(initial float64, function func(float64, float64) float64) float64

	/*
		Creates a new list containing elements of the old one, satisfying a condition.
		The function has one parameter, value of the current element, and returns bool.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - filtered list.
	*/
	Filter(function func(any) bool) List

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
	FilterObjects(function func(Object) bool) List

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
	FilterLists(function func(List) bool) List

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
	FilterStrings(function func(string) bool) List

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
	FilterInts(function func(int) bool) List

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
	FilterFloats(function func(float64) bool) List

	/*
		Computes a sum of all elements in the list.
		The list has to be homogeneous and all its elements have to be ints.

		Returns:
		  - computed sum (int).
	*/
	IntSum() int

	/*
		Computes a sum of all elements in the list.
		The list has to be homogeneous and all its elements have to be numeric.

		Returns:
		  - computed sum (float).
	*/
	Sum() float64

	/*
		Computes a product of all elements in the list.
		The list has to be homogeneous and all its elements have to be ints.

		Returns:
		  - computed product (int).
	*/
	IntProd() int

	/*
		Computes a product of all elements in the list.
		The list has to be homogeneous and all its elements have to be numeric.

		Returns:
		  - computed pruduct (float).
	*/
	Prod() float64

	/*
		Computes an arithmetic average of all elements in the list.
		The list has to be homogeneous and all its elements have to be numeric.

		Returns:
		  - computed average (float).
	*/
	Avg() float64

	/*
		Finds a minimum of the list.
		The list has to be homogeneous and all its elements have to be ints.

		Returns:
		  - found minimum (int).
	*/
	IntMin() int

	/*
		Finds a minimum of the list.
		The list has to be homogeneous and all its elements have to be numeric.

		Returns:
		  - found minimum (float).
	*/
	Min() float64

	/*
		Finds a maximum of the list.
		The list has to be homogeneous and all its elements have to be ints.

		Returns:
		  - found maximum (int).
	*/
	IntMax() int

	/*
		Finds a maximum of the list.
		The list has to be homogeneous and all its elements have to be numeric.

		Returns:
		  - found maximum (float).
	*/
	Max() float64

	/*
		Parallelly executes a given function over an every element of the list.
		The function has two parameters: index of the current element and its value.
		The order of the iterations is random.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged list.
	*/
	ForEachAsync(function func(int, any)) List

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
	MapAsync(function func(int, any) any) List

	/*
		Acquires the element specified by the given tree form.

		Parameters:
		  - tf - tree form string.

		Returns:
		  - corresponding value (any type, has to be asserted).
	*/
	GetTF(tf string) any

	/*
		Sets the element specified by the given tree form.

		Parameters:
		  - tf - tree form string,
		  - value - value to set.

		Returns:
		  - updated list.
	*/
	SetTF(tf string, value any) List
}
