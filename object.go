/*
AnyType Library for Go
Object (dictionary) type
*/

package anytype

/*
Object is an unordered set of key-value pairs.
*/
type Object interface {
	field

	/*
		Init initializes the ego pointer, which allows deriving.

		Parameters:
		  - ptr - ego pointer.
	*/
	Init(ptr Object)

	/*
		Ego acquires the ego pointer previously set by Init.

		Returns:
		  - ego pointer.
	*/
	Ego() Object

	/*
		Set sets values of the fields.
		If the key already exists, the value is overwritten, if not, new field is created.
		If one key is given multiple times, the value is set to the last provided value.

		Parameters:
		  - values... - any amount of key-value pairs to set.

		Returns:
		  - updated object.
	*/
	Set(values ...any) Object

	/*
		Unset deletes the fields with the given keys. If the key does not exist, nothing happens.

		Parameters:
		  - keys... - any amount of keys to delete.

		Returns:
		  - updated object.
	*/
	Unset(keys ...string) Object

	/*
		Clear deletes all fields in the object.

		Returns:
		  - updated object.
	*/
	Clear() Object

	/*
		Get acquires a value under the specified key of the object.

		Parameters:
		  - key - key of the field to get.

		Returns:
		  - corresponding value (any type, has to be asserted).
	*/
	Get(key string) any

	/*
		GetObject acquires a nested object under the specified key of the object.
		Causes a panic if the field has another type.

		Parameters:
		  - key - key of the field to get.

		Returns:
		  - corresponding value asserted as object.
	*/
	GetObject(key string) Object

	/*
		GetList acquires a list under the specified key of the object.
		Causes a panic if the field has another type.

		Parameters:
		  - key - key of the field to get.

		Returns:
		  - corresponding value asserted as list.
	*/
	GetList(key string) List

	/*
		GetString acquires a string under the specified key of the object.
		Causes a panic if the field has another type.

		Parameters:
		  - key - key of the field to get.

		Returns:
		  - corresponding value asserted as string.
	*/
	GetString(key string) string

	/*
		GetBool acquires a bool under the specified key of the object.
		Causes a panic if the field has another type.

		Parameters:
		  - key - key of the field to get.

		Returns:
		  - corresponding value asserted as bool.
	*/
	GetBool(key string) bool

	/*
		GetInt acquires an int under the specified key of the object.
		Causes a panic if the field has another type.

		Parameters:
		  - key - key of the field to get.

		Returns:
		  - corresponding value asserted as int.
	*/
	GetInt(key string) int

	/*
		GetFloat acquires a float under the specified key of the object.
		Causes a panic if the field has another type.

		Parameters:
		  - key - key of the field to get.

		Returns:
		  - corresponding value asserted as float.
	*/
	GetFloat(key string) float64

	/*
		TypeOf gives a type of the field under the specified key of the object.
		If the key does not exist, 0 (TypeUndefined) is returned.

		Parameters:
		  - key - key of the field.

		Returns:
		  - integer constant representing the type (see type enum).
	*/
	TypeOf(key string) Type

	/*
		String gives a JSON representation of the object, including nested objects and lists.

		Returns:
		  - JSON string.
	*/
	String() string

	/*
		FormatString gives a JSON representation of the object in standardized format with the given indentation.

		Parameters:
		  - indent - indentation spaces (0-10).

		Returns:
		  - JSON string.
	*/
	FormatString(indent int) string

	/*
		Dict converts the object into a Go map of empty interfaces.

		Returns:
		  - map.
	*/
	Dict() map[string]any

	/*
		Keys convers the object to a list of its keys.

		Returns:
		  - list of keys of the object.
	*/
	Keys() List

	/*
		Values convers the object to a list of its values.

		Returns:
		  - list of values of the object.
	*/
	Values() List

	/*
		Clone creates a deep copy of the object.

		Returns:
		  - copied object.
	*/
	Clone() Object

	/*
		Count gives a number of fields of the object.

		Returns:
		  - number of fields.
	*/
	Count() int

	/*
		Empty checks whether the object is empty.

		Returns:
		  - true if the object is empty, false otherwise.
	*/
	Empty() bool

	/*
		Equals checks if the content of the object is equal to the content of another object.
		Nested objects and lists are compared recursively (by value).

		Parameters:
		  - another - an object to compare with.

		Returns:
		  - true if the objects are equal, false otherwise.
	*/
	Equals(another Object) bool

	/*
		Merge creates a new object containing all elements of the old object and another object.
		The old object remains unchanged.
		If both objects contain the same key, the value from another object is used.

		Parameters:
		  - another - an object to merge.

		Returns:
		  - new object.
	*/
	Merge(another Object) Object

	/*
		Pluck creates a new object containing the given fields of the existing object.
		The old object remains unchanged.

		Parameters:
		  - keys... - any amount of keys to be in the new object.

		Returns:
		  - new object.
	*/
	Pluck(keys ...string) Object

	/*
		Contains checks if the object contains a field with a given value.
		Objects and lists are compared by reference.

		Parameters:
		  - value - the value to check.

		Returns:
		  - true if the object contains the value, false otherwise.
	*/
	Contains(value any) bool

	/*
		KeyOf gives a key containing a given value.
		If multiple keys contain the value, any of them is returned.

		Parameters:
		  - value - the value to check.

		Returns:
		  - key for the value (empty string if the object does not contain the value).
	*/
	KeyOf(value any) string

	/*
		KeyExists checks if a given key exists within the object.

		Parameters:
		  - key - the key to check.

		Returns:
		  - true if the key exists, false otherwise.
	*/
	KeyExists(key string) bool

	/*
		ForEach executes a given function over an every field of the object.
		The function has two parameters: key of the current field and its value.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEach(function func(key string, val any)) Object

	/*
		ForEachValue executes a given function over an every field of the object.
		The function has one parameter, value of the current field.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachValue(function func(x any)) Object

	/*
		ForEachObject executes a given function over all objects nested in the object.
		Fields of other types are ignored.
		The function has one parameter, the current object.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachObject(function func(x Object)) Object

	/*
		ForEachList executes a given function over all lists in the object.
		Fields of other types are ignored.
		The function has one parameter, the current list.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachList(function func(x List)) Object

	/*
		ForEachString executes a given function over all strings in the object.
		Fields of other types are ignored.
		The function has one parameter, the current string.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachString(function func(x string)) Object

	/*
		ForEachBool executes a given function over all bools in the object.
		Fields of other types are ignored.
		The function has one parameter, the current bool.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachBool(function func(x bool)) Object

	/*
		ForEachInt executes a given function over all ints in the object.
		Fields of other types are ignored.
		The function has one parameter, the current int.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachInt(function func(x int)) Object

	/*
		ForEachFloat executes a given function over all floats in the object.
		Fields of other types are ignored.
		The function has one parameter, the current float.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachFloat(function func(x float64)) Object

	/*
		Map copies the object and modifies each field by a given mapping function.
		The resulting field can have a different type than the original one.
		The function has two parameters: current key and value of the current element. Returns any.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new list.
	*/
	Map(function func(key string, val any) any) Object

	/*
		MapValues copies the object and modifies each field by a given mapping function.
		The resulting field can have a different type than the original one.
		The function has one parameter, value of the current field, and returns any.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapValues(function func(x any) any) Object

	/*
		MapObjects selects all nested objects in the object and modifies each of them by a given mapping function.
		Fields with other types are ignored.
		The resulting field can have a different type than the original one.
		The function has one parameter, the current object, and returns any.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapObjects(function func(x Object) any) Object

	/*
		MapLists selects all lists in the object and modifies each of them by a given mapping function.
		Fields with other types are ignored.
		The resulting field can have a different type than the original one.
		The function has one parameter, the current list, and returns any.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapLists(function func(x List) any) Object

	/*
		MapStrings selects all strings in the object and modifies each of them by a given mapping function.
		Fields with other types are ignored.
		The resulting field can have a different type than the original one.
		The function has one parameter, the current string, and returns empty interface.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapStrings(function func(x string) any) Object

	/*
		MapBools selects all bools in the object and modifies each of them by a given mapping function.
		Fields with other types are ignored.
		The resulting field can have a different type than the original one.
		The function has one parameter, the current bool, and returns empty interface.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapBools(function func(x bool) any) Object

	/*
		MapInts selects all ints in the object and modifies each of them by a given mapping function.
		Fields with other types are ignored.
		The resulting field can have a different type than the original one.
		The function has one parameter, the current int, and returns empty interface.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapInts(function func(x int) any) Object

	/*
		MapFloats selects all floats in the object and modifies each of them by a given mapping function.
		Fields with other types are ignored.
		The resulting field can have a different type than the original one.
		The function has one parameter, the current float, and returns empty interface.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapFloats(function func(x float64) any) Object

	/*
		ForEachAsync parallelly executes a given function over an every field of the object.
		The function has two parameters: key of the current field and its value.
		The order of the iterations is random.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - unchanged object.
	*/
	ForEachAsync(function func(key string, val any)) Object

	/*
		MapAsync copies the object and paralelly modifies each field by a given mapping function.
		The resulting field can have a different type than the original one.
		The function has two parameters: key of the current field and its value.
		The old object remains unchanged.

		Parameters:
		  - function - anonymous function to be executed.

		Returns:
		  - new object.
	*/
	MapAsync(function func(key string, val any) any) Object

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

		Parameters:
		  - tf - tree form string,
		  - value - value to set.

		Returns:
		  - updated object.
	*/
	SetTF(tf string, value any) Object

	/*
		TypeOfTF gives a type of the field specified by a given tree form.
		If the TF path does not exist, 0 (TypeUndefined) is returned.

		Parameters:
		  - tf - tree form string.

		Returns:
		  - integer constant representing the type (see type enum).
	*/
	TypeOfTF(tf string) Type
}
