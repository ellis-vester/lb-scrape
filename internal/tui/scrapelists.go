package tui

import (
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ellis-vester/lb-scrape/files"
	lb "github.com/ellis-vester/lb-scrape/letterboxd"
	"github.com/ellis-vester/lb-scrape/scraper"
)

var _ tea.Model = &ScrapeListsModel{}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
var titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EDA4BD")).Render
var textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DEEFB7")).Render
var headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#DEEFB7")).Render

var statusStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#EDA4BD")).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#725AC1")).Render

func NewScrapeListsModel(listsPath string, outputDir string, pollInterval int) *ScrapeListsModel {

	listProgress := progress.New(progress.WithSolidFill("#DEEFB7"))
	listProgress.Width = 80

	filmProgress := progress.New(progress.WithSolidFill("#DEEFB7"))
	filmProgress.Width = 80

	progressSpinner := spinner.New()
	progressSpinner.Spinner = spinner.Pulse

	return &ScrapeListsModel{
		spinner:      progressSpinner,
		listProgress: listProgress,
		filmProgress: filmProgress,
		err:          nil,
		ListsPath:    listsPath,
		OutputDir:    outputDir,
		PollInterval: pollInterval,
	}
}

type ScrapeListsModel struct {
	spinner      spinner.Model
	listProgress progress.Model
	filmProgress progress.Model
	status       string
	err          error

	UnscrapedLists []string
	ScrapedLists   [][]lb.FilmListEntry

	UnscrapedFilms []lb.FilmListEntry
	ScrapedFilms   []lb.Film
	Directors      []lb.Director

	ListsPath    string
	OutputDir    string
	PollInterval int
}

func (m ScrapeListsModel) Init() tea.Cmd {
	return tea.Batch(getLists(m.ListsPath), m.spinner.Tick)
}

type startMsg string

func (m ScrapeListsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case listsReadFromDiskMsg:
		m.status = "Read lists from disk"
		if msg.Err != nil {
			m.status = m.err.Error()
			return m, tea.Quit
		}

		m.UnscrapedLists = msg.Lists

		if len(m.UnscrapedLists) != len(m.ScrapedLists) {
			m.status = "Scraping list " + m.UnscrapedLists[len(m.ScrapedLists)]

			cmd = scrapeFilmList(m.UnscrapedLists[len(m.ScrapedLists)], m.PollInterval)
			cmds = append(cmds, cmd)
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	case listScrapedResponseMsg:
		if msg.Err != nil {
			m.status = msg.Err.Error()
			return m, tea.Quit
		}

		m.ScrapedLists = append(m.ScrapedLists, msg.Films)

		if len(m.ScrapedLists) != len(m.UnscrapedLists) {
			m.status = "Scraping list " + m.UnscrapedLists[len(m.ScrapedLists)]

			cmd = scrapeFilmList(m.UnscrapedLists[len(m.ScrapedLists)], m.PollInterval)
			cmds = append(cmds, cmd)
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

		if len(m.ScrapedLists) == len(m.UnscrapedLists) {
			// Dedup, then start scraping films
			m.UnscrapedFilms = scraper.SumFilmInclusions(m.ScrapedLists)
			m.status = "Scraping film " + m.UnscrapedFilms[len(m.ScrapedFilms)].Link

			cmd = scrapeFilm(m.UnscrapedFilms[len(m.ScrapedFilms)], m.PollInterval)
			cmds = append(cmds, cmd)
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	case filmScrapedResponseMsg:
		if msg.Err != nil {
			m.status = msg.Err.Error()
			return m, tea.Quit
		}

		m.ScrapedFilms = append(m.ScrapedFilms, lb.Film{
			Title:      msg.Film.Title,
			Director:   msg.Film.Director,
			Year:       msg.Film.Year,
			Rating:     m.UnscrapedFilms[len(m.ScrapedFilms)].Rating,
			Inclusions: m.UnscrapedFilms[len(m.ScrapedFilms)].Inclusions,
			Link:       m.UnscrapedFilms[len(m.ScrapedFilms)].Link,
		})

		if len(m.UnscrapedFilms) != len(m.ScrapedFilms) {
			m.status = "Scraping film " + m.UnscrapedFilms[len(m.ScrapedFilms)].Link

			cmd = scrapeFilm(m.UnscrapedFilms[len(m.ScrapedFilms)], m.PollInterval)
			cmds = append(cmds, cmd)
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		} else {

			m.Directors = scraper.SumDirectorInclusions(m.ScrapedFilms)
			m.status = "Writing to disk..."
			files.WriteFilmsToCsv(m.ScrapedFilms, m.OutputDir+"/films.csv")
			files.WriteDirectorsToCsv(m.Directors, m.OutputDir+"/directors.csv")
			m.status = "Done!"
			return m, tea.Quit
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, tea.Batch(cmds...)
}

func (m ScrapeListsModel) View() string {
	progressPad := strings.Repeat(" ", 2)
	detailsPad := strings.Repeat(" ", 99)

	listDenominator := len(m.UnscrapedLists)
	if listDenominator == 0 {

		listDenominator = 1
	}

	filmDenominator := len(m.UnscrapedFilms)
	if filmDenominator == 0 {
		filmDenominator = 1
	}

	scrapedFilm := lb.Film{
		Title:      "",
		Director:   "",
		Year:       0,
		Rating:     0,
		Inclusions: 0,
		Link:       "",
	}

	filmDisplay := statusStyle(detailsPad + "\n\n\n\n\n\n")

	title := `
 ___       ________                 ________  ________  ________  ________  ________  _______      
|\  \     |\   __  \               |\   ____\|\   ____\|\   __  \|\   __  \|\   __  \|\  ___ \     
\ \  \    \ \  \|\ /_  ____________\ \  \___|\ \  \___|\ \  \|\  \ \  \|\  \ \  \|\  \ \   __/|    
 \ \  \    \ \   __  \|\____________\ \_____  \ \  \    \ \   _  _\ \   __  \ \   ____\ \  \_|/__  
  \ \  \____\ \  \|\  \|____________|\|____|\  \ \  \____\ \  \\  \\ \  \ \  \ \  \___|\ \  \_|\ \ 
   \ \_______\ \_______\               ____\_\  \ \_______\ \__\\ _\\ \__\ \__\ \__\    \ \_______\
    \|_______|\|_______|              |\_________\|_______|\|__|\|__|\|__|\|__|\|__|     \|_______|
                                      \|_________|
  `

	if len(m.ScrapedFilms) != 0 {
		scrapedFilm = m.ScrapedFilms[len(m.ScrapedFilms)-1]
		filmDisplay = statusStyle("\n" + titleStyle("  Title:      ") + textStyle(scrapedFilm.Title) + "\n" +
			titleStyle("  Director:   ") + textStyle(scrapedFilm.Director) + "\n" +
			titleStyle("  Year:       ") + textStyle(strconv.FormatInt(int64(scrapedFilm.Year), 10)) + "\n" +
			titleStyle("  Inclusions: ") + textStyle(strconv.FormatInt(int64(scrapedFilm.Inclusions), 10)) + "\n" +
			titleStyle("  Link:       ") + textStyle("https://letterboxd.com"+scrapedFilm.Link) + "\n" + detailsPad)
	}

	return headerStyle(title) + "\n" +
		statusStyle("\n"+progressPad+m.spinner.View()+" "+m.status+"\n\n"+
			progressPad+titleStyle("Lists: ")+m.listProgress.ViewAs(float64(len(m.ScrapedLists))/float64(listDenominator))+progressPad+"        "+"\n\n"+
			progressPad+titleStyle("Films: ")+m.filmProgress.ViewAs(float64(len(m.ScrapedFilms))/float64(filmDenominator))+progressPad+"        "+"\n\n") +
		"\n" + filmDisplay + "\n" +
		progressPad + helpStyle("Press q or ctrl+c to quit") + "\n"
}

// Messages
type listsReadFromDiskMsg struct {
	Lists []string
	Err   error
}

type listScrapedResponseMsg struct {
	Films []lb.FilmListEntry
	Err   error
}

type filmScrapedResponseMsg struct {
	Film lb.Film
	Err  error
}

type filmScrapedMsg lb.Film

// Commands
func getLists(path string) tea.Cmd {
	return func() tea.Msg {
		urls, err := files.GetFilmListUrls(path)
		return listsReadFromDiskMsg{
			Lists: urls,
			Err:   err,
		}
	}
}

func scrapeFilmList(url string, interval int) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Duration(interval) * time.Second)

		html, err := scraper.ScrapeFilmListHtml(url)
		if err != nil {
			return listScrapedResponseMsg{
				Films: nil,
				Err:   err,
			}
		}

		films, err := scraper.ParseFilmList(html)
		if err != nil {
			return listScrapedResponseMsg{
				Films: nil,
				Err:   err,
			}
		}

		return listScrapedResponseMsg{
			Films: films,
			Err:   nil,
		}
	}
}

func scrapeFilm(film lb.FilmListEntry, interval int) tea.Cmd {
	return func() tea.Msg {

		time.Sleep(time.Duration(interval) * time.Second)

		html, err := scraper.ScrapeFilmHtml("https://letterboxd.com" + film.Link)
		if err != nil {
			return filmScrapedResponseMsg{
				Film: lb.Film{},
				Err:  err,
			}
		}

		film, err := scraper.ParseFilm(html)
		if err != nil {
			return filmScrapedResponseMsg{
				Film: lb.Film{},
				Err:  err,
			}
		}

		return filmScrapedResponseMsg{
			Film: film,
			Err:  nil,
		}
	}
}
