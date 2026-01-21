// Package ui contient le cœur graphique de l'application Fyne.
// app.go orchestre la fenêtre principale, la barre latérale de navigation et la gestion d'état globale.
// C'est le composant central qui connecte tous les autres modules UI et gère le flux d'affichage principal.
package ui

import (
	"Groupie-Tracker/api"
	"Groupie-Tracker/models"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// StartApp lance l'application Fyne principale
func StartApp() {
	a := app.New()
	w := a.NewWindow("Groupie Tracker")

	w.Resize(fyne.NewSize(1200, 800))
	w.CenterOnScreen()

	w.SetContent(CreateMainLayout(a, w))
	w.ShowAndRun()
}

type AppState struct {
	app            fyne.App
	window         fyne.Window
	allArtists     []models.Artist
	mainContent    *fyne.Container
	selectedArtist *models.Artist
}

// CreateMainLayout crée le layout Spotify-like avec sidebar et contenu principal
func CreateMainLayout(app fyne.App, window fyne.Window) fyne.CanvasObject {
	artists, err := api.GetArtists()
	if err != nil {
		return widget.NewLabel("Erreur: " + err.Error())
	}

	state := &AppState{
		app:        app,
		window:     window,
		allArtists: artists,
	}

	sidebar := createSidebar(state)

	state.mainContent = container.NewVBox()
	mainScroll := container.NewVScroll(state.mainContent)

	displayArtistGrid(state, artists)

	mainLayout := container.NewHBox(
		sidebar,
		mainScroll,
	)

	return mainLayout
}

// createSidebar construit la barre latérale avec navigation
func createSidebar(state *AppState) fyne.CanvasObject {
	bg := canvas.NewRectangle(color.NRGBA{R: 20, G: 20, B: 20, A: 255})

	title := widget.NewLabel("🎵 GROUPIE\nTRACKER 🎵")
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	sep1 := widget.NewSeparator()

	allArtistsBtn := widget.NewButton("🎤 Tous les artistes", func() {
		displayArtistGrid(state, state.allArtists)
	})
	allArtistsBtn.Importance = widget.MediumImportance

	searchBtn := widget.NewButton("🔍 Rechercher", func() {
		displaySearchView(state)
	})
	searchBtn.Importance = widget.MediumImportance

	filterBtn := widget.NewButton("⚙️ Filtres", func() {
		displayFilterView(state)
	})
	filterBtn.Importance = widget.MediumImportance

	buttonsBox := container.NewVBox(
		allArtistsBtn,
		searchBtn,
		filterBtn,
	)

	sidebarContent := container.NewVBox(
		title,
		sep1,
		buttonsBox,
	)

	sidebarRect := container.NewStack(bg)
	sizeRect := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	sizeRect.SetMinSize(fyne.NewSize(200, 600))

	return container.NewStack(sizeRect, sidebarRect, container.NewPadded(sidebarContent))
}

// displayArtistGrid affiche la grille principale des artistes
func displayArtistGrid(state *AppState, artists []models.Artist) {
	state.mainContent.Objects = nil

	header := createMainHeader()
	state.mainContent.Add(header)
	state.mainContent.Add(widget.NewSeparator())

	gridTitle := widget.NewLabel("Tous les artistes")
	gridTitle.TextStyle.Bold = true
	gridTitle.Alignment = fyne.TextAlignCenter
	state.mainContent.Add(gridTitle)

	grid := createArtistGrid(state, artists)
	state.mainContent.Add(grid)
	state.mainContent.Refresh()
}

// createMainHeader construit le bandeau visuel avec image vinyle
func createMainHeader() fyne.CanvasObject {
	vinylURL := "https://images.unsplash.com/photo-1603048588665-791ca8aea617?w=1200&h=300&fit=crop"

	headerBg := canvas.NewRectangle(color.NRGBA{R: 20, G: 20, B: 20, A: 255})
	headerBg.SetMinSize(fyne.NewSize(1000, 250))

	vinylContainer := container.NewStack(headerBg)
	go func() {
		resp, err := http.Get(vinylURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			tmpFile, err := os.CreateTemp("", "vinyl-bg-*.jpg")
			if err == nil {
				defer tmpFile.Close()
				io.Copy(tmpFile, resp.Body)

				vinylImg := canvas.NewImageFromFile(tmpFile.Name())
				vinylImg.FillMode = canvas.ImageFillStretch
				vinylImg.SetMinSize(fyne.NewSize(1000, 250))

				vinylContainer.Objects = []fyne.CanvasObject{vinylImg}
				vinylContainer.Refresh()
			}
		}
	}()

	titleText := canvas.NewText("Groupie Tracker", color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	titleText.TextSize = 60
	titleText.TextStyle.Bold = true

	subtitleLabel := widget.NewLabel("Fondé en B1")
	subtitleLabel.TextStyle.Bold = true
	subtitleLabel.Alignment = fyne.TextAlignCenter

	creatorsLabel := widget.NewLabel("par Olivier, Adama et Sarah")
	creatorsLabel.Alignment = fyne.TextAlignCenter

	overlay := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 180})
	overlay.SetMinSize(fyne.NewSize(1000, 250))

	content := container.NewVBox(
		container.NewCenter(titleText),
		subtitleLabel,
		creatorsLabel,
	)

	return container.NewStack(vinylContainer, overlay, container.NewPadded(content))
}

// createArtistGrid organise les cartes artistes en grille 4 colonnes
func createArtistGrid(state *AppState, artists []models.Artist) fyne.CanvasObject {
	grid := container.NewVBox()

	for i := 0; i < len(artists); i += 4 {
		row := container.NewHBox()
		for j := 0; j < 4 && i+j < len(artists); j++ {
			artist := artists[i+j]
			card := createArtistCard(state, artist)
			row.Add(card)
		}
		grid.Add(row)
	}

	return grid
}

// createArtistCard fabrique la carte individuelle d'un artiste
func createArtistCard(state *AppState, artist models.Artist) fyne.CanvasObject {
	img := state.loadArtistImage(artist.Image)

	nameLabel := widget.NewLabel(artist.Name)
	nameLabel.Alignment = fyne.TextAlignCenter

	btn := widget.NewButton("Détails", func() {
		displayArtistDetail(state, artist)
	})

	card := container.NewVBox(
		img,
		nameLabel,
		btn,
	)

	cardBg := canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 255})
	cardSizeRect := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	cardSizeRect.SetMinSize(fyne.NewSize(280, 350))

	return container.NewStack(cardSizeRect, cardBg, container.NewPadded(card))
}

// loadArtistImage charge l'image d'un artiste en arrière-plan avec placeholder
func (s *AppState) loadArtistImage(imageURL string) fyne.CanvasObject {
	return LoadImageAsync(imageURL, 260, 260)
}

// displayArtistDetail remplace le contenu par la vue détail d'un artiste
func displayArtistDetail(state *AppState, artist models.Artist) {
	state.mainContent.Objects = nil

	detail := CreateArtistDetailView(artist, state.app, func() {
		displayArtistGrid(state, state.allArtists)
	})

	state.mainContent.Add(detail)
	state.mainContent.Refresh()
}

// displaySearchView affiche la page de recherche textuelle
func displaySearchView(state *AppState) {
	state.mainContent.Objects = nil

	titleText := canvas.NewText("🔍 Rechercher", color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	titleText.TextSize = 40
	titleText.TextStyle.Bold = true
	titleText.Alignment = fyne.TextAlignCenter

	titleBg := canvas.NewRectangle(color.NRGBA{R: 20, G: 20, B: 20, A: 255})
	titleBg.SetMinSize(fyne.NewSize(800, 80))
	titleHeader := container.NewStack(titleBg, container.NewCenter(titleText))

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Entrez le nom d'un artiste...")

	resultsContainer := container.NewVBox()
	resultsScroll := container.NewVScroll(resultsContainer)
	resultsScroll.SetMinSize(fyne.NewSize(900, 600))

	searchBtn := widget.NewButton("Rechercher", func() {
		resultsContainer.Objects = nil
		query := searchEntry.Text
		if query == "" {
			resultsContainer.Add(widget.NewLabel("Veuillez entrer un terme de recherche"))
			resultsContainer.Refresh()
			return
		}

		// Filtrer les artistes
		var results []models.Artist
		query = strings.ToLower(query)
		for _, artist := range state.allArtists {
			if strings.Contains(strings.ToLower(artist.Name), query) {
				results = append(results, artist)
			}
		}

		if len(results) == 0 {
			resultsContainer.Add(widget.NewLabel("Aucun artiste trouvé"))
			resultsContainer.Refresh()
			return
		}

		// Afficher les résultats en grille
		grid := createArtistGrid(state, results)
		resultsContainer.Add(grid)
		resultsContainer.Refresh()
	})

	backBtn := widget.NewButton("← Retour", func() {
		displayArtistGrid(state, state.allArtists)
	})

	searchBox := container.NewVBox(
		searchEntry,
		container.NewHBox(backBtn, searchBtn),
	)

	content := container.NewBorder(
		container.NewVBox(titleHeader, widget.NewSeparator(), searchBox, widget.NewSeparator()),
		nil, nil, nil,
		resultsScroll,
	)

	state.mainContent.Objects = []fyne.CanvasObject{content}
	state.mainContent.Refresh()
}

// displayFilterView affiche la vue des filtres
func displayFilterView(state *AppState) {
	state.mainContent.Objects = nil

	// Titre
	title := widget.NewLabel("⚙️ Filtres")
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	// Filtre par année
	yearLabel := widget.NewLabel("Année de création:")
	yearLabel.TextStyle.Bold = true
	yearMinEntry := widget.NewEntry()
	yearMinEntry.SetPlaceHolder("Min")
	yearMaxEntry := widget.NewEntry()
	yearMaxEntry.SetPlaceHolder("Max")
	yearBox := container.NewHBox(yearMinEntry, widget.NewLabel("à"), yearMaxEntry)

	// Filtre par album
	albumLabel := widget.NewLabel("Premier album:")
	albumLabel.TextStyle.Bold = true
	albumMinEntry := widget.NewEntry()
	albumMinEntry.SetPlaceHolder("Min")
	albumMaxEntry := widget.NewEntry()
	albumMaxEntry.SetPlaceHolder("Max")
	albumBox := container.NewHBox(albumMinEntry, widget.NewLabel("à"), albumMaxEntry)

	// Filtre par nombre de membres
	memberLabel := widget.NewLabel("Nombre de membres:")
	memberLabel.TextStyle.Bold = true
	memberMinEntry := widget.NewEntry()
	memberMinEntry.SetPlaceHolder("Min")
	memberMaxEntry := widget.NewEntry()
	memberMaxEntry.SetPlaceHolder("Max")
	memberBox := container.NewHBox(memberMinEntry, widget.NewLabel("à"), memberMaxEntry)

	// Filtre par localisation
	locationLabel := widget.NewLabel("Localisation:")
	locationLabel.TextStyle.Bold = true
	locationEntry := widget.NewEntry()
	locationEntry.SetPlaceHolder("Entrez une localisation...")

	// Bouton de recherche
	filterBtn := widget.NewButton("Appliquer les filtres", func() {
		// Implémenter la logique de filtrage
		var results []models.Artist

		for _, artist := range state.allArtists {
			// Filtrer par année
			if yearMinEntry.Text != "" || yearMaxEntry.Text != "" {
				minYear, maxYear := 0, 9999
				if yearMinEntry.Text != "" {
					fmt.Sscanf(yearMinEntry.Text, "%d", &minYear)
				}
				if yearMaxEntry.Text != "" {
					fmt.Sscanf(yearMaxEntry.Text, "%d", &maxYear)
				}
				if artist.CreationDate < minYear || artist.CreationDate > maxYear {
					continue
				}
			}

			// Filtrer par album
			if albumMinEntry.Text != "" || albumMaxEntry.Text != "" {
				minYear, maxYear := 0, 9999
				if albumMinEntry.Text != "" {
					fmt.Sscanf(albumMinEntry.Text, "%d", &minYear)
				}
				if albumMaxEntry.Text != "" {
					fmt.Sscanf(albumMaxEntry.Text, "%d", &maxYear)
				}
				albumYear := 0
				fmt.Sscanf(artist.FirstAlbum, "%d", &albumYear)
				if albumYear < minYear || albumYear > maxYear {
					continue
				}
			}

			// Filtrer par nombre de membres
			if memberMinEntry.Text != "" || memberMaxEntry.Text != "" {
				minMembers, maxMembers := 0, 9999
				if memberMinEntry.Text != "" {
					fmt.Sscanf(memberMinEntry.Text, "%d", &minMembers)
				}
				if memberMaxEntry.Text != "" {
					fmt.Sscanf(memberMaxEntry.Text, "%d", &maxMembers)
				}
				memberCount := len(artist.Members)
				if memberCount < minMembers || memberCount > maxMembers {
					continue
				}
			}

			// Filtrer par localisation
			if locationEntry.Text != "" {
				locationFound := false
				locQuery := strings.ToLower(locationEntry.Text)
				for _, loc := range artist.Locations {
					if strings.Contains(strings.ToLower(loc), locQuery) {
						locationFound = true
						break
					}
				}
				if !locationFound {
					continue
				}
			}

			results = append(results, artist)
		}

		// Afficher les résultats
		state.mainContent.Objects = nil
		resultsTitle := widget.NewLabel(fmt.Sprintf("Résultats: %d artiste(s)", len(results)))
		resultsTitle.TextStyle.Bold = true
		resultsTitle.Alignment = fyne.TextAlignCenter

		state.mainContent.Add(resultsTitle)
		if len(results) > 0 {
			grid := createArtistGrid(state, results)
			state.mainContent.Add(grid)
		} else {
			state.mainContent.Add(widget.NewLabel("Aucun artiste ne correspond aux filtres"))
		}
		state.mainContent.Refresh()
	})

	// Bouton Retour
	backBtn := widget.NewButton("← Retour", func() {
		displayArtistGrid(state, state.allArtists)
	})

	// Layout
	filtersScroll := container.NewVScroll(
		container.NewVBox(
			yearLabel,
			yearBox,
			widget.NewSeparator(),
			albumLabel,
			albumBox,
			widget.NewSeparator(),
			memberLabel,
			memberBox,
			widget.NewSeparator(),
			locationLabel,
			locationEntry,
		),
	)

	// Rendre filtersScroll plus grand
	filtersScroll.SetMinSize(fyne.NewSize(800, 600))

	content := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator()),
		container.NewHBox(backBtn, filterBtn),
		nil, nil,
		filtersScroll,
	)

	state.mainContent.Objects = []fyne.CanvasObject{content}
	state.mainContent.Refresh()
}
