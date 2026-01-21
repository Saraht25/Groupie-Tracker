// Package ui - home.go affiche la page d'accueil avec recherche, suggestions et filtres avancés.
// C'est la première interface que l'utilisateur voit, permettant de chercher et de filtrer les artistes
// avant d'accéder aux détails. C'est le cœur de la navigation et de la découverte.
package ui

import (
	"Groupie-Tracker/api"
	"Groupie-Tracker/models"
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

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

// Home construit la page d'accueil avec recherche, suggestions et filtres
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

	state.searchEntry = widget.NewEntry()
	state.searchEntry.SetPlaceHolder("Chercher artiste, membre, lieu...")
	state.searchEntry.OnChanged = func(q string) { state.applySearch(q) }

	titleLabel := widget.NewLabel("Groupie Tracker")
	titleLabel.TextStyle.Bold = true
	titleContainer := container.NewCenter(titleLabel)

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

	filterButton := widget.NewButton("Filtres avances", func() {
		state.showAdvancedFilters()
	})

	state.filterLabel = widget.NewLabel("Tous les artistes (" + strconv.Itoa(len(state.filtered)) + ")")

	searchBox := container.NewVBox(
		titleContainer,
		state.searchEntry,
		state.listWrap,
		container.NewHBox(filterButton, state.filterLabel),
	)
	searchBox = container.NewPadded(searchBox)

	mainContent := container.NewVScroll(state.cards)

	headerBg := canvas.NewRectangle(color.NRGBA{R: 30, G: 60, B: 120, A: 255})
	header := container.NewStack(headerBg, searchBox)

	content := container.NewBorder(
		header, nil, nil, nil,
		mainContent,
	)

	state.homeView = content
	state.mainContainer = container.NewStack(content)

	return state.mainContainer
}

// showAdvancedFilters ouvre une fenêtre pour affiner la recherche
func (s *homeState) showAdvancedFilters() {
	filterWindow := s.app.NewWindow("Filtres avances")
	filterWindow.Resize(fyne.NewSize(400, 500))

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

	locationLabel := widget.NewLabel("Localisation contient:")
	s.locationQuery = widget.NewEntry()
	s.locationQuery.SetPlaceHolder("Ex: Paris, France")

	applyButton := widget.NewButton("Appliquer filtres", func() {
		s.applyAdvancedFilters()
		filterWindow.Close()
	})

	resetButton := widget.NewButton("Reinitialiser", func() {
		s.creationMin.SetText("")
		s.creationMax.SetText("")
		s.albumMin.SetText("")
		s.albumMax.SetText("")
		s.memberCountMin.SetText("")
		s.memberCountMax.SetText("")
		if s.locationQuery != nil {
			s.locationQuery.SetText("")
		}

		s.applyAdvancedFilters()

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

// applyAdvancedFilters applique tous les filtres saisis et rafraichit la grille
func (s *homeState) applyAdvancedFilters() {
	criteria := FilterCriteria{
		CreationMin:    s.creationMin.Text,
		CreationMax:    s.creationMax.Text,
		AlbumMin:       s.albumMin.Text,
		AlbumMax:       s.albumMax.Text,
		MemberCountMin: s.memberCountMin.Text,
		MemberCountMax: s.memberCountMax.Text,
	}
	if s.locationQuery != nil {
		criteria.LocationQuery = s.locationQuery.Text
	}

	s.filtered = ApplyFilters(s.allArtists, criteria)
	s.renderCards()
	s.updateFilterLabel()
}

// applySearch filtre et suggère à partir du texte tapé
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

// updateFilterLabel met à jour le compteur et le libellé des filtres actifs
func (s *homeState) updateFilterLabel() {
	count := len(s.filtered)
	label := fmt.Sprintf("Artistes affiches (%d)", count)
	if s.locationQuery != nil {
		if strings.TrimSpace(s.locationQuery.Text) != "" {
			label += " • filtre lieu: " + strings.TrimSpace(s.locationQuery.Text)
		}
	}
	if s.filterLabel != nil {
		s.filterLabel.SetText(label)
	}
}

// renderCards reconstruit la grille d'artistes
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

		row := container.NewCenter(container.NewHBox(rowCards...))
		rows = append(rows, row)
	}

	s.cards.Objects = rows
	s.cards.Refresh()
}

// createArtistCard fabrique une carte cliquable pour un artiste
func (s *homeState) createArtistCard(artist models.Artist) fyne.CanvasObject {
	img := s.loadArtistImage(artist.Image)

	nameLabel := widget.NewLabel(artist.Name)
	nameLabel.TextStyle.Bold = true
	nameLabel.Alignment = fyne.TextAlignCenter

	cardContent := container.NewVBox(img, nameLabel)
	paddedBox := container.NewPadded(cardContent)

	rect := canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 255})
	cardWithBg := container.NewStack(rect, paddedBox)

	clickableArea := widget.NewButton("", func() {
		s.showArtistDetail(artist)
	})
	clickableArea.Importance = widget.LowImportance

	return container.NewStack(
		clickableArea,
		cardWithBg,
	)
}

// loadArtistImage charge l'image d'un artiste en arrière-plan avec un placeholder
func (s *homeState) loadArtistImage(imageURL string) fyne.CanvasObject {
	return LoadImageAsync(imageURL, 160, 160)
}

// showArtistDetail remplace la vue courante par les détails de l'artiste
func (s *homeState) showArtistDetail(artist models.Artist) {
	detailView := CreateArtistDetailView(artist, s.app, func() {
		if s.homeView != nil && s.window != nil {
			s.window.SetContent(s.homeView)
		}
	})

	if s.window != nil {
		s.window.SetContent(detailView)
	} else {
		s.mainContainer.Objects = []fyne.CanvasObject{detailView}
		s.mainContainer.Refresh()
	}
}
