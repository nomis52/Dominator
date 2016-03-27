/*
	Package flagutil provides utility types for the standard flag package.

	Package flagutil provides a collection of types to enhance the standard
	flag package, such as slices.
*/
package flagutil

// A StringList is a slice of strings that satisfies the standard library
// flag.Value interface.
type StringList []string
