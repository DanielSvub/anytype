/*
AnyType Library for Go
Object (dictionary) type
*/

package anytype

/*
Interface for an object.

Extends:
  - field.
*/
type Object interface {
	field

	/*
		Initializes the ego pointer, which allows deriving.

		Parameters:
		  - ptr - ego pointer.
	*/
	Init(ptr Object)

	/*
		Acquires the ego pointer previously set by Init.

		Returns:
		  - ego pointer.
	*/
	Ego() Object

	/*
		Sets a values of the fields.
		If the key already exists, the value is overwritten, if not, new field is created.
		If one key is given multiple times, the value is set to the last one.

		Parameters:
		  - values... - any amount of key-value pairs to set.

		Returns:
		  - updated object.
	*/
	Set(values ...any) Object

	/*
		Deletes the fields with given keys.

		Parameters:
		  - keys... - any amount of keys to delete.

		Returns:
		  - updated object.
	*/
	Unset(keys ...string) Object

	/*
		Deletes all field of the object.

		Returns:
		  - updated object.
	*/
	Clear() Object

	/*
		Acquires the value under the specified key of the object.

		Parameters:
		  - key - key of the field to get.

		Returns:
		  - corresponding value (any type, has to be asserted).
	*/
	Get(key string) any

	/*
		Acquires the nested object under the specified key of the object.
		Causes a panic if the field has another type.

		Parameters:
		  - key - key of the field to get.

		Returns:
		  - corresponding value asserted as object.
	*/
	GetObject(key string) Object

	/*
		Acquires the list under the specified key of the object.
		Causes a panic if the field has another type.

		Parameters:
		  - key - key of the field to get.

		Returns:
		  - corresponding value asserted as list.
	*/
	GetList(key string) List

	/*
		Acquires the string under the specified key of the object.
		Causes a panic if the field has another type.

		Parameters:
		  - key - key of the field to get.

		Returns:
		  - corresponding value asserted as string.
	*/
	GetString(key string) string

	/*
		Acquires the bool under the specified key of the object.
		Causes a panic if the field has another type.

		Parameters:
		  - key - key of the field to get.

		Returns:
		  - corresponding value asserted as bool.
	*/
	GetBool(key string) bool

	/*
		Acquires the int under the specified key of the object.
		Causes a panic if the field has another type.

		Parameters:
		  - key - key of the field to get.

		Returns:
		  - corresponding value asserted as int.
	*/
	GetInt(key string) int

	/*
		Acquires the float under the specified key of the object.
		Causes a panic if the field has another type.

		Parameters:
		  - key - key of the field to get.

		Returns:
		  - corresponding value asserted as float.
	*/
	GetFloat(key string) float64

	/*
		Gives a type of the field under the specified key of the object.

		Parameters:
		  - key - key of the field.

		Returns:
		  - integer constant representing the type (see type enum).
	*/
	TypeOf(key string) Type

	/*
		Gives a JSON representation of the object, including nested objects and lists.

		Returns:
		  - JSON string.
	*/
	String() string

	/*
		Gives a JSON representation of the object in standardized format with the given indentation.

		Parameters:
		  - indent - indentation spaces (0-10).

		Returns:
		  - JSON string.
	*/
	FormatString(indent int) string

	/*
		Converts the object into a Go map of empty interfaces.

		Returns:
		  - map.
	*/
	Dict() map[string]any

	/*
		Convers the object to a list of its keys.

		Returns:
		  - list of keys of the object.
	*/
	Keys() List

	/*
		Convers the object to a list of its values.

		Returns:
		  - list of values of the object.
	*/
	Values() List

	/*
		Creates a deep copy of the object.

		Returns:
		  - copied object.
	*/
	Clone() Object

	/*
		Gives a number of fields of the object.

		Returns:
		  - number of fields.
	*/
	Count() int

	/*
		Checks whether the object is empty.

		Returns:
		  - true if the object is empty, false otherwise.
	*/
	Empty() bool

	/*
		Checks if the content of the object is equal to the content of another object.
		Nested objects and lists are compared recursively (by value).

		Parameters:
		  - another - an object to compare with.

		Returns:
		  - true if the objects are equal, false otherwise.
	*/
	Equals(another Object) bool

	/*
		Creates a new object containing all elements of the old object and another object.
		The old object remains unchanged.
		If both objects contain a key, the value from another object is used.

		Parameters:
		  - another - an object to merge.

		Returns:
		  - new object.
	*/
	Merge(another Object) Object

	/*
		Creates a new object containing the given fields of the existing object.

		Parameters:
		  - keys... - any amount of keys to be in the new object.

		Returns:
		  - created plucked object.
	*/
	Pluck(keys ...string) Object

	/*
		Checks if the object contains a field with a given value.
		Objects and lists are compared by reference.

		Parameters:
		  - value - the value to check.

		Returns:
		  - true if the object contains the value, false otherwise.
	*/
	Contains(value any) bool

	/*
		Gives a key containing a given value.
		If multiple keys contain the value, any of them is returned.

		Parameters:
		  - value - the value to check.

		Returns:
		  - key for the value (empty string if the object does not contain the value).
	*/
	KeyOf(value any) string

	/*
		Checks if a given key exists within the object.

		Parameters:
		  - key - the key to check.

		Returns:
		  - true if the key exists, false otherwise.
	*/
	KeyExists(key string) bool

	/*
		Executes a given function over an every field of the object.
		The function has two parameters: key of the current field and its value.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEach(function func(string, any)) Object

	/*
		Executes a given function over an every field of the object.
		The function has one parameter, value of the current field.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachValue(function func(any)) Object

	/*
		Executes a given function over all objects nested in the object.
		Fields with other types are ignored.
		The function has one parameter, the current object.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachObject(function func(Object)) Object

	/*
		Executes a given function over all lists in the object.
		Fields with other types are ignored.
		The function has one parameter, the current list.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachList(function func(List)) Object

	/*
		Executes a given function over all strings in the object.
		Fields with other types are ignored.
		The function has one parameter, the current string.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachString(function func(string)) Object

	/*
		Executes a given function over all bools in the object.
		Fields with other types are ignored.
		The function has one parameter, the current bool.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachBool(function func(bool)) Object

	/*
		Executes a given function over all ints in the object.
		Fields with other types are ignored.
		The function has one parameter, the current int.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachInt(function func(int)) Object

	/*
		Executes a given function over all floats in the object.
		Fields with other types are ignored.
		The function has one parameter, the current float.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachFloat(function func(float64)) Object

	/*
		Copies the object and modifies each field by a given mapping function.
		The resulting field can have a different type than the original one.
		The function has two parameters: current key and value of the current element. Returns empty interface.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new list.
	*/
	Map(function func(string, any) any) Object

	/*
		Copies the object and modifies each field by a given mapping function.
		The resulting field can have a different type than the original one.
		The function has one parameter, value of the current field, and returns empty interface.
		The old list remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapValues(function func(any) any) Object

	/*
		Selects all nested objects of the object and modifies each of them by a given mapping function.
		Fields with other types are ignored.
		The resulting field can have a different type than the original one.
		The function has one parameter, the current object, and returns empty interface.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapObjects(function func(Object) any) Object

	/*
		Selects all lists of the object and modifies each of them by a given mapping function.
		Fields with other types are ignored.
		The resulting field can have a different type than the original one.
		The function has one parameter, the current list, and returns empty interface.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapLists(function func(List) any) Object

	/*
		Selects all nested strings of the object and modifies each of them by a given mapping function.
		Fields with other types are ignored.
		The resulting field can have a different type than the original one.
		The function has one parameter, the current string, and returns empty interface.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapStrings(function func(string) any) Object

	/*
		Selects all nested ints of the object and modifies each of them by a given mapping function.
		Fields with other types are ignored.
		The resulting field can have a different type than the original one.
		The function has one parameter, the current int, and returns empty interface.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapInts(function func(int) any) Object

	/*
		Selects all nested floats of the object and modifies each of them by a given mapping function.
		Fields with other types are ignored.
		The resulting field can have a different type than the original one.
		The function has one parameter, the current float, and returns empty interface.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapFloats(function func(float64) any) Object

	/*
		Parallelly executes a given function over an every field of the object.
		The function has two parameters: key of the current field and its value.
		The order of the iterations is random.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachAsync(function func(string, any)) Object

	/*
		Copies the object and paralelly modifies each field by a given mapping function.
		The resulting field can have a different type than the original one.
		The function has two parameters: key of the current field and its value.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapAsync(function func(string, any) any) Object

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
		  - updated object.
	*/
	SetTF(tf string, value any) Object
}
