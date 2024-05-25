package letterboxd

// Film contains the data about a single film on Letterboxd.
type Film struct {
	Link       string
	Rating     int8
	UserName   string
	Inclusions int
	Director   string
	Year       int
	Title      string
}
