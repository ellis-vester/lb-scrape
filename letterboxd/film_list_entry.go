package letterboxd

// FilmListEntry represents the data for a film in Letterboxd list.
type FilmListEntry struct {
	Rating     int8
	Link       string
	UserName   string
	Inclusions int
}
