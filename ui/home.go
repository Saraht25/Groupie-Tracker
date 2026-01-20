package ui

import (
	"Groupie-Tracker/api"
	"Groupie-Tracker/models"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// homeState contient l'état global de la page d'accueil
type homeState struct {
	allArtists     []models.Artist
	filtered       []models.Artist
	cards          *fyne.Container
	suggestions    []Suggestion
	list           *widget.List
	listWrap       *container.Scroll
	searchEntry    *widget.Entry
	app            fyne.App
	filterLabel    *widget.Label
	creationMin    *widget.Entry
	creationMax    *widget.Entry
	albumMin       *widget.Entry
	albumMax       *widget.Entry
	memberCountMin *widget.Entry
	memberCountMax *widget.Entry
	mainContainer  *fyne.Container
	window         fyne.Window
	homeView       fyne.CanvasObject
	locationQuery  *widget.Entry
}

// Home affiche l'écran principal avec grille d'artistes, recherche et filtres en français
func Home(app fyne.App, window fyne.Window) fyne.CanvasObject {
	artists, err := api.GetArtists()
	if err != nil {
		return widget.NewLabel("Erreur lors du chargement: " + err.Error())
	}

	state := &homeState{
		allArtists: artists,
		filtered:   artists,
		app:        app,
		window:     window,
	}

	state.cards = container.NewVBox()
	state.renderCards()

	// Barre de recherche
	state.searchEntry = widget.NewEntry()
	state.searchEntry.SetPlaceHolder("Chercher artiste, membre, lieu...")
	state.searchEntry.OnChanged = func(q string) { state.applySearch(q) }

	// Titre
	titleLabel := widget.NewLabel("Groupie Tracker")
	titleLabel.TextStyle.Bold = true
	titleContainer := container.NewCenter(titleLabel)

	// Liste des suggestions
	state.list = widget.NewList(
		func() int { return len(state.suggestions) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, co fyne.CanvasObject) {
			if id < 0 || id >= len(state.suggestions) {
				return
			}
			sug := state.suggestions[id]
			co.(*widget.Label).SetText(fmt.Sprintf("%s [%s]", sug.Label, sug.Type))
		},
	)
	state.list.OnSelected = func(id widget.ListItemID) {
		if id < 0 || id >= len(state.suggestions) {
			return
		}
		choice := state.suggestions[id]
		state.searchEntry.SetText(choice.Label)
		state.applySearch(choice.Label)
		state.list.Unselect(id)
	}
	state.listWrap = container.NewVScroll(state.list)
	state.listWrap.SetMinSize(fyne.NewSize(0, 100))
	state.listWrap.Hide()

	// Bouton de filtres avancés
	filterButton := widget.NewButton("Filtres avances", func() {
		state.showAdvancedFilters()
	})

	state.filterLabel = widget.NewLabel("Tous les artistes (" + strconv.Itoa(len(state.filtered)) + ")")

	// En-tête avec recherche
	searchBox := container.NewVBox(
		titleContainer,
		state.searchEntry,
		state.listWrap,
		container.NewHBox(filterButton, state.filterLabel),
	)
	searchBox = container.NewPadded(searchBox)

	// Contenu principal avec scroll
	mainContent := container.NewVScroll(state.cards)

	// Ajouter une couleur d'accent derrière l'en-tête
	headerBg := canvas.NewRectangle(color.NRGBA{R: 30, G: 60, B: 120, A: 255})
	header := container.NewStack(headerBg, searchBox)

	content := container.NewBorder(
		header, nil, nil, nil,
		mainContent,
	)

	// Stocker la vue principale et le conteneur principal pour pouvoir changer de vue
	state.homeView = content
	state.mainContainer = container.NewStack(content)

	return state.mainContainer
}

// showAdvancedFilters affiche une fenêtre avec des filtres avancés
func (s *homeState) showAdvancedFilters() {
	filterWindow := s.app.NewWindow("Filtres avances")
	filterWindow.Resize(fyne.NewSize(400, 500))

	// Année de création
	creationLabel := widget.NewLabel("Année de creation:")
	s.creationMin = widget.NewEntry()
	s.creationMin.SetPlaceHolder("Min")
	s.creationMax = widget.NewEntry()
	s.creationMax.SetPlaceHolder("Max")
	creationBox := container.NewHBox(
		widget.NewLabel("De"),
		s.creationMin,
		widget.NewLabel("à"),
		s.creationMax,
	)

	// Année du premier album
	albumLabel := widget.NewLabel("Année du premier album:")
	s.albumMin = widget.NewEntry()
	s.albumMin.SetPlaceHolder("Min")
	s.albumMax = widget.NewEntry()
	s.albumMax.SetPlaceHolder("Max")
	albumBox := container.NewHBox(
		widget.NewLabel("De"),
		s.albumMin,
		widget.NewLabel("à"),
		s.albumMax,
	)

	// Nombre de membres
	memberLabel := widget.NewLabel("Nombre de membres:")
	s.memberCountMin = widget.NewEntry()
	s.memberCountMin.SetPlaceHolder("Min")
	s.memberCountMax = widget.NewEntry()
	s.memberCountMax.SetPlaceHolder("Max")
	memberBox := container.NewHBox(
		widget.NewLabel("De"),
		s.memberCountMin,
		widget.NewLabel("à"),
		s.memberCountMax,
	)

	// Filtre par localisation (texte)
	locationLabel := widget.NewLabel("Localisation contient:")
	s.locationQuery = widget.NewEntry()
	s.locationQuery.SetPlaceHolder("Ex: Paris, France")

	// Boutons d'action
	applyButton := widget.NewButton("Appliquer filtres", func() {
		s.applyAdvancedFilters()
		filterWindow.Close()
	})

	resetButton := widget.NewButton("Reinitialiser", func() {
		// Remettre à zéro tous les champs
		s.creationMin.SetText("")
		s.creationMax.SetText("")
		s.albumMin.SetText("")
		s.albumMax.SetText("")
		s.memberCountMin.SetText("")
		s.memberCountMax.SetText("")
		if s.locationQuery != nil {
			s.locationQuery.SetText("")
		}

		// Recalculer l'affichage avec filtres vides
		s.applyAdvancedFilters()

		// Fermer la fenêtre des filtres
		filterWindow.Close()
	})

	buttonBox := container.NewHBox(applyButton, resetButton)

	content := container.NewVBox(
		creationLabel,
		creationBox,
		widget.NewSeparator(),
		albumLabel,
		albumBox,
		widget.NewSeparator(),
		memberLabel,
		memberBox,
		widget.NewSeparator(),
		locationLabel,
		s.locationQuery,
		widget.NewSeparator(),
		buttonBox,
	)

	scroll := container.NewVScroll(content)
	filterWindow.SetContent(scroll)
	filterWindow.Show()
}

// applyAdvancedFilters applique tous les filtres actifs
func (s *homeState) applyAdvancedFilters() {
	s.filtered = s.allArtists

	// Filtre année de création
	if s.creationMin.Text != "" || s.creationMax.Text != "" {
		var min, max int
		if s.creationMin.Text != "" {
			min, _ = strconv.Atoi(s.creationMin.Text)
		} else {
			min = 0
		}
		if s.creationMax.Text != "" {
			max, _ = strconv.Atoi(s.creationMax.Text)
		} else {
			max = 9999
		}

		var temp []models.Artist
		for _, a := range s.filtered {
			if a.CreationDate >= min && a.CreationDate <= max {
				temp = append(temp, a)
			}
		}
		s.filtered = temp
	}

	// Filtre année du premier album
	if s.albumMin.Text != "" || s.albumMax.Text != "" {
		var min, max int
		if s.albumMin.Text != "" {
			min, _ = strconv.Atoi(s.albumMin.Text)
		} else {
			min = 0
		}
		if s.albumMax.Text != "" {
			max, _ = strconv.Atoi(s.albumMax.Text)
		} else {
			max = 9999
		}

		var temp []models.Artist
		for _, a := range s.filtered {
			if year, ok := firstYearFromString(a.FirstAlbum); ok {
				if year >= min && year <= max {
					temp = append(temp, a)
				}
			}
		}
		s.filtered = temp
	}

	// Filtre nombre de membres
	if s.memberCountMin.Text != "" || s.memberCountMax.Text != "" {
		var min, max int
		if s.memberCountMin.Text != "" {
			min, _ = strconv.Atoi(s.memberCountMin.Text)
		} else {
			min = 0
		}
		if s.memberCountMax.Text != "" {
			max, _ = strconv.Atoi(s.memberCountMax.Text)
		} else {
			max = 9999
		}

		var temp []models.Artist
		for _, a := range s.filtered {
			count := len(a.Members)
			if count >= min && count <= max {
				temp = append(temp, a)
			}
		}
		s.filtered = temp
	}

	// Filtre par localisation (texte en sous-chaîne)
	if s.locationQuery != nil {
		locQ := strings.TrimSpace(strings.ToLower(s.locationQuery.Text))
		if locQ != "" {
			var temp []models.Artist
			for _, a := range s.filtered {
				matched := false
				for _, loc := range a.Locations {
					if strings.Contains(strings.ToLower(loc), locQ) {
						matched = true
						break
					}
				}
				if matched {
					temp = append(temp, a)
				}
			}
			s.filtered = temp
		}
	}

	s.renderCards()
	s.updateFilterLabel()
}

// applySearch applique la recherche par texte
func (s *homeState) applySearch(q string) {
	query := strings.TrimSpace(q)
	s.filtered = SearchArtists(query, s.allArtists)
	s.suggestions = BuildSuggestions(query, s.allArtists)
	if len(s.suggestions) == 0 {
		s.listWrap.Hide()
	} else {
		s.listWrap.Show()
	}
	s.list.Refresh()
	s.renderCards()
	s.updateFilterLabel()
}

// updateFilterLabel met à jour le label affichant le nombre d'artistes
func (s *homeState) updateFilterLabel() {
	count := len(s.filtered)
	label := fmt.Sprintf("Artistes affiches (%d)", count)
	if s.locationQuery != nil {
		// Protéger contre nil ou valeurs inattendues
		if strings.TrimSpace(s.locationQuery.Text) != "" {
		label += " • filtre lieu: " + strings.TrimSpace(s.locationQuery.Text)
		}
	}
	if s.filterLabel != nil {
		s.filterLabel.SetText(label)
	}
}

// renderCards crée la grille d'artistes
func (s *homeState) renderCards() {
	const cardsPerRow = 4
	var rows []fyne.CanvasObject

	for i := 0; i < len(s.filtered); i += cardsPerRow {
		var rowCards []fyne.CanvasObject
		end := i + cardsPerRow
		if end > len(s.filtered) {
			end = len(s.filtered)
		}

		for j := i; j < end; j++ {
			card := s.createArtistCard(s.filtered[j])
			rowCards = append(rowCards, card)
		}

		// Centrer chaque ligne de cartes
		row := container.NewCenter(container.NewHBox(rowCards...))
		rows = append(rows, row)
	}

	s.cards.Objects = rows
	s.cards.Refresh()
}

// createArtistCard crée une carte cliquable pour un artiste
func (s *homeState) createArtistCard(artist models.Artist) fyne.CanvasObject {
	// Image artiste
	img := s.loadArtistImage(artist.Image)

	// Nom artiste
	nameLabel := widget.NewLabel(artist.Name)
	nameLabel.TextStyle.Bold = true
	nameLabel.Alignment = fyne.TextAlignCenter

	// Contenu de la carte
	cardContent := container.NewVBox(img, nameLabel)
	paddedBox := container.NewPadded(cardContent)

	// Fond sombre
	rect := canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 255})
	cardWithBg := container.NewStack(rect, paddedBox)

	// Bouton cliquable pour ouvrir les détails
	clickableArea := widget.NewButton("", func() {
		s.showArtistDetail(artist)
	})
	clickableArea.Importance = widget.LowImportance

	// Retourner le bouton qui contient la carte
	return container.NewStack(
		clickableArea,
		cardWithBg,
	)
}

// loadArtistImage charge et affiche l'image d'un artiste
func (s *homeState) loadArtistImage(imageURL string) fyne.CanvasObject {
	resp, err := http.Get(imageURL)
	if err != nil {
		return widget.NewLabel("Indisponible")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return widget.NewLabel("Image non trouvee")
	}

	// Fichier temporaire pour l'image
	tmpFile, err := os.CreateTemp("", "artist-*.jpg")
	if err != nil {
		return widget.NewLabel("Erreur image")
	}
	defer tmpFile.Close()

	// Écrire l'image
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return widget.NewLabel("Erreur image")
	}

	tmpPath := tmpFile.Name()

	// Image Fyne
	img := canvas.NewImageFromFile(tmpPath)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(160, 160))

	return img
}

// showArtistDetail affiche les détails de l'artiste dans la même fenêtre
func (s *homeState) showArtistDetail(artist models.Artist) {
	detailView := CreateArtistDetailView(artist, s.app, func() {
		// Revenir à la vue d'accueil
		if s.homeView != nil && s.window != nil {
			s.window.SetContent(s.homeView)
		}
	})

	// Afficher la vue détail en remplaçant entièrement le contenu de la fenêtre
	if s.window != nil {
		s.window.SetContent(detailView)
	} else {
		// Fallback si la fenêtre n'est pas disponible
		s.mainContainer.Objects = []fyne.CanvasObject{detailView}
		s.mainContainer.Refresh()
	}
}
