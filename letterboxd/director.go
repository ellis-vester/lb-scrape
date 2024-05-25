// Package letterboxd contains structs representing the data available on Letterboxd.
package letterboxd

// Director represents a film director and the number of
// times they appear in a list or lists on Letterboxd.
type Director struct {
	Name       string
	Inclusions int
}
