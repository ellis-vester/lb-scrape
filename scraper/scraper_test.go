package scraper

import (
	"reflect"
	"testing"

	lb "github.com/ellis-vester/lb-scrape/letterboxd"
)

func TestParseFilmList_ReturnsValidFilmListEntry(t *testing.T) {

	got, err := ParseFilmList(`
		<li class="poster-container" data-owner-rating="10"> <div class="film-poster" data-target-link="/film/faust-1926/"></div></li>
		<li class="poster-container" data-owner-rating="8"> <div class="film-poster" data-target-link="/film/parasite/"></div></li>`)
	if err != nil {
		t.Errorf("got %v, want %v", err, nil)
	}

	want := []lb.FilmListEntry{
		{Rating: 10, Link: "/film/faust-1926/"},
		{Rating: 8, Link: "/film/parasite/"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseFilmList_ReturnsNonNilErrorWhenRatingEmpty(t *testing.T) {

	_, err := ParseFilmList(`<li class="poster-container" data-owner-rating=""> <div class="film-poster" data-target-link="/film/faust-1926/"></div></li>`)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParseFilmList_ReturnsNonNilErrorWhenRatingNotPresent(t *testing.T) {

	_, err := ParseFilmList(`<li class="poster-container"> <div class="film-poster" data-target-link="/film/faust-1926/"></div></li>`)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParseFilmList_ReturnsNonNilErrorWhenRatingNotInt(t *testing.T) {

	_, err := ParseFilmList(`<li class="poster-container" data-owner-rating="ten"> <div class="film-poster" data-target-link="/film/faust-1926/"></div></li>`)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParseFilmList_ReturnsNonNilErrorWhenLinkEmpty(t *testing.T) {

	_, err := ParseFilmList(`<li class="poster-container" data-owner-rating="10"> <div class="film-poster" data-target-link=""></div></li>`)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParseFilmList_ReturnsNonNilErrorWhenLinkNotPresent(t *testing.T) {

	_, err := ParseFilmList(`<li class="poster-container" data-owner-rating="10"> <div class="film-poster"></div></li>`)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParseFilm_ReturnsValidFilm(t *testing.T) {

	got, err := ParseFilm(`
	<div class="details">			
		<h1 class="headline-1 filmtitle">
		<span class="name js-widont prettify">Wild at Heart</span>
		</h1>
		<div class="metablock">
			<div class="releaseyear">
				<a href="/films/year/1990/">1990</a>
			</div>
			<p class="credits">
				<span class="introduction">Directed by</span>
				<span class="directorlist">
					<a class="contributor" href="/director/david-lynch/">
					<span class="prettify">David Lynch</span></a>
				</span>
			</p>
		</div>
	</div>`)
	if err != nil {
		t.Errorf("got %v, want %v", err, nil)
	}

	want := lb.Film{
		Title:    "Wild at Heart",
		Director: "David Lynch",
		Year:     1990,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseFilm_ReturnsNonNilErrorWhenTitleEmpty(t *testing.T) {

	_, err := ParseFilm(`
	<div class="details">			
		<h1 class="headline-1 filmtitle">
		<span class="name js-widont prettify"></span>
		</h1>
		<div class="metablock">
			<div class="releaseyear">
				<a href="/films/year/1990/">1990</a>
			</div>
			<p class="credits">
				<span class="introduction">Directed by</span>
				<span class="directorlist">
					<a class="contributor" href="/director/david-lynch/">
					<span class="prettify">David Lynch</span></a>
				</span>
			</p>
		</div>
	</div>`)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParseFilm_ReturnsNonNilErrorWhenTitleNotPresent(t *testing.T) {

	_, err := ParseFilm(`
	<div class="details">			
		<div class="metablock">
			<div class="releaseyear">
				<a href="/films/year/1990/">1990</a>
			</div>
			<p class="credits">
				<span class="introduction">Directed by</span>
				<span class="directorlist">
					<a class="contributor" href="/director/david-lynch/">
					<span class="prettify">David Lynch</span></a>
				</span>
			</p>
		</div>
	</div>`)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParseFilm_ReturnsNonNilErrorWhenYearEmpty(t *testing.T) {

	_, err := ParseFilm(`
	<div class="details">			
		<h1 class="headline-1 filmtitle">
		<span class="name js-widont prettify">Wild at Heart</span>
		</h1>
		<div class="metablock">
			<div class="releaseyear">
				<a href="/films/year/1990/"></a>
			</div>
			<p class="credits">
				<span class="introduction">Directed by</span>
				<span class="directorlist">
					<a class="contributor" href="/director/david-lynch/">
					<span class="prettify">David Lynch</span></a>
				</span>
			</p>
		</div>
	</div>`)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParseFilm_ReturnsNonNilErrorWhenYearNotPresent(t *testing.T) {

	_, err := ParseFilm(`
	<div class="details">			
		<h1 class="headline-1 filmtitle">
		<span class="name js-widont prettify">Wild at Heart</span>
		</h1>
		<div class="metablock">
			<p class="credits">
				<span class="introduction">Directed by</span>
				<span class="directorlist">
					<a class="contributor" href="/director/david-lynch/">
					<span class="prettify">David Lynch</span></a>
				</span>
			</p>
		</div>
	</div>`)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParseFilm_ReturnsNonNilErrorWhenYearNotInt(t *testing.T) {

	_, err := ParseFilm(`
	<div class="details">			
		<h1 class="headline-1 filmtitle">
		<span class="name js-widont prettify">Wild at Heart</span>
		</h1>
		<div class="metablock">
			<div class="releaseyear">
				<a href="/films/year/1990/">nineteen ninety</a>
			</div>
			<p class="credits">
				<span class="introduction">Directed by</span>
				<span class="directorlist">
					<a class="contributor" href="/director/david-lynch/">
					<span class="prettify">David Lynch</span></a>
				</span>
			</p>
		</div>
	</div>`)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParseFilm_ReturnsNonNilErrorWhenDirectorEmpty(t *testing.T) {

	_, err := ParseFilm(`
	<div class="details">			
		<h1 class="headline-1 filmtitle">
		<span class="name js-widont prettify">Wild at Heart</span>
		</h1>
		<div class="metablock">
			<div class="releaseyear">
				<a href="/films/year/1990/">1990</a>
			</div>
			<p class="credits">
				<span class="introduction">Directed by</span>
				<span class="directorlist">
					<a class="contributor" href="/director/david-lynch/">
					<span class="prettify"></span></a>
				</span>
			</p>
		</div>
	</div>`)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParseFilm_ReturnsNonNilErrorWhenDirectorNotPresent(t *testing.T) {

	_, err := ParseFilm(`
	<div class="details">			
		<h1 class="headline-1 filmtitle">
		<span class="name js-widont prettify">Wild at Heart</span>
		</h1>
		<div class="metablock">
			<div class="releaseyear">
				<a href="/films/year/1990/">1990</a>
			</div>
			<p class="credits">
				<span class="introduction">Directed by</span>
				<span class="directorlist">
				</span>
			</p>
		</div>
	</div>`)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestParseUsername(t *testing.T) {
	got := ParseUsername("https://letterboxd.com/username/list/2023-favs/")
	want := "username"
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSumFilmInclusions(t *testing.T) {
	films := [][]lb.FilmListEntry{
		{
			{Rating: 8, Link: "/film/parasite/"},
			{Rating: 10, Link: "/film/faust-1926/"},
		},
		{
			{Rating: 10, Link: "/film/faust-1926/"},
			{Rating: 8, Link: "/film/wild-at-heart/"},
		},
		{
			{Rating: 8, Link: "/film/wild-at-heart/"},
			{Rating: 10, Link: "/film/faust-1926/"},
		},
	}

	got := SumFilmInclusions(films)
	want := []lb.FilmListEntry{
		{Inclusions: 3, Link: "/film/faust-1926/", Rating: 10},
		{Inclusions: 2, Link: "/film/wild-at-heart/", Rating: 8},
		{Inclusions: 1, Link: "/film/parasite/", Rating: 8},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSumDirectorInclusions(t *testing.T) {
	films := []lb.Film{
		{Rating: 8, Link: "/film/totally-fucked-up/", Director: "Gregg Araki"},
		{Rating: 8, Link: "/film/nowhere/", Director: "Gregg Araki"},
		{Rating: 10, Link: "/film/wild-at-heart/", Director: "David Lynch"},
		{Rating: 8, Link: "/film/inland-empire/", Director: "David Lynch"},
		{Rating: 8, Link: "/film/eraser-head/", Director: "David Lynch"},
		{Rating: 8, Link: "/film/dead-ringers/", Director: "David Cronenburg"},
	}

	got := SumDirectorInclusions(films)
	want := []lb.Director{
		{Name: "David Lynch", Inclusions: 3},
		{Name: "Gregg Araki", Inclusions: 2},
		{Name: "David Cronenburg", Inclusions: 1},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
