package ui

import (
	"Groupie-Tracker/models"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CreateArtistDetailView crée une vue de détails pour un artiste
func CreateArtistDetailView(artist models.Artist, app fyne.App, onBack func()) fyne.CanvasObject {

	// Image artiste
	img := loadDetailImage(artist.Image)

	// Infos de base
	nameLabel := widget.NewLabel(artist.Name)
	nameLabel.TextStyle.Bold = true
	nameLabel.Alignment = fyne.TextAlignCenter

	creationLabel := widget.NewLabel(
		fmt.Sprintf("Année de création: %d", artist.CreationDate),
	)
	albumLabel := widget.NewLabel(
		fmt.Sprintf("Premier album: %s", artist.FirstAlbum),
	)

	// Membres
	membersLabel := widget.NewLabel("Membres:")
	membersLabel.TextStyle.Bold = true
	membersBox := container.NewVBox()
	for _, member := range artist.Members {
		membersBox.Add(widget.NewLabel("  • " + member))
	}

	// Lieux de concert avec carte
	locationsLabel := widget.NewLabel("Lieux de concert:")
	locationsLabel.TextStyle.Bold = true
	locationsBox := container.NewVBox()
	for _, loc := range artist.Locations {
		locationsBox.Add(widget.NewLabel("  • " + loc))
	}

	// Carte des localisations
	mapImage := createLocationMap(artist.Locations)

	// Bouton retour
	backButton := widget.NewButtonWithIcon("Retour", theme.NavigateBackIcon(), func() {
		onBack()
	})
	backButton.Importance = widget.HighImportance

	// Lecteur Spotify (recherche l'artiste sur Spotify)
	spotifyLabel := widget.NewLabel("Écouter sur Spotify:")
	spotifyLabel.TextStyle.Bold = true
	spotifyButton := widget.NewButton("Ouvrir Spotify", func() {
		// Créer un lien de recherche Spotify encodé
		searchQuery := url.QueryEscape(artist.Name)
		spotifyURL, _ := url.Parse(fmt.Sprintf("https://open.spotify.com/search/%s", searchQuery))
		app.OpenURL(spotifyURL)
	})

	// Bouton pour ouvrir la carte interactive OpenStreetMap
	var openMapButton fyne.CanvasObject
	if len(artist.Locations) > 0 {
		if lat, lon, ok := geocodeLocation(artist.Locations[0]); ok {
			osmURL := fmt.Sprintf("https://www.openstreetmap.org/?mlat=%s&mlon=%s#map=5/%s/%s", lat, lon, lat, lon)
			u, _ := url.Parse(osmURL)
			openMapButton = widget.NewButton("Voir sur OpenStreetMap", func() { app.OpenURL(u) })
		} else {
			openMapButton = widget.NewLabel("Localisation introuvable")
		}
	} else {
		openMapButton = widget.NewLabel("Aucune localisation")
	}

	// Colonne gauche: Image et carte (plus petite)
	leftCol := container.NewVBox(
		container.NewCenter(img),
		widget.NewSeparator(),
		container.NewCenter(widget.NewLabel("Carte des localisations")),
		container.NewCenter(mapImage),
		container.NewCenter(openMapButton),
	)

	// Colonne droite: Informations
	rightCol := container.NewVBox(
		nameLabel,
		creationLabel,
		albumLabel,
		widget.NewSeparator(),
		membersLabel,
		membersBox,
		widget.NewSeparator(),
		locationsLabel,
		locationsBox,
		widget.NewSeparator(),
		spotifyLabel,
		spotifyButton,
	)

	// Layout en deux colonnes avec scroll
	contentBody := container.NewHBox(
		container.NewPadded(leftCol),
		container.NewPadded(rightCol),
	)

	scroll := container.NewVScroll(contentBody)

	// En-tête avec bouton retour
	// En-tête coloré avec bouton retour
	headerBg := canvas.NewRectangle(color.NRGBA{R: 30, G: 60, B: 120, A: 255})
	header := container.NewStack(headerBg, container.NewHBox(backButton))

	// Vue complète
	view := container.NewBorder(
		container.NewPadded(header), nil, nil, nil,
		scroll,
	)

	return view
}

func loadDetailImage(imageURL string) fyne.CanvasObject {
	if imageURL == "" {
		return widget.NewLabel("Image indisponible")
	}

	// Télécharger l'image
	resp, err := http.Get(imageURL)
	if err != nil {
		return widget.NewLabel("Erreur de chargement")
	}
	defer resp.Body.Close()

	// Créer un fichier temporaire
	tmpFile, err := os.CreateTemp("", "artist-detail-*.jpg")
	if err != nil {
		return widget.NewLabel("Erreur de cache")
	}
	defer tmpFile.Close()

	// Copier l'image
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return widget.NewLabel("Erreur de copie")
	}

	// Charger l'image
	img := canvas.NewImageFromFile(tmpFile.Name())
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(200, 200))
	return img
}

// createLocationMap crée une image de carte avec les localisations
func createLocationMap(locations []string) fyne.CanvasObject {
	if len(locations) == 0 {
		return widget.NewLabel("Aucune localisation disponible")
	}

	// Géocoder plusieurs localisations (limite pour éviter une URL trop longue)
	maxMarkers := 10
	type coord struct{ lat, lon string }
	var coords []coord
	for i, loc := range locations {
		if i >= maxMarkers {
			break
		}
		lat, lon, ok := geocodeLocation(loc)
		if ok {
			coords = append(coords, coord{lat: lat, lon: lon})
		}
	}
	if len(coords) == 0 {
		return createMapPlaceholder(locations)
	}

	// Centrer sur la première coordonnée
	centerLat := coords[0].lat
	centerLon := coords[0].lon

	// Construire la chaîne des marqueurs: lon,lat,color
	var markers []string
	for _, c := range coords {
		markers = append(markers, fmt.Sprintf("%s,%s,red", c.lon, c.lat))
	}
	markersParam := strings.Join(markers, "|")

	// Créer l'URL de la carte statique OSM
	mapURL := fmt.Sprintf(
		"https://staticmap.openstreetmap.de/staticmap.php?center=%s,%s&zoom=5&size=600x400&maptype=mapnik&markers=%s",
		centerLat, centerLon, markersParam,
	)

	// Télécharger l'image de carte
	resp, err := http.Get(mapURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return createMapPlaceholder(locations)
	}
	defer resp.Body.Close()

	// Créer un fichier temporaire
	tmpFile, err := os.CreateTemp("", "map-*.png")
	if err != nil {
		return createMapPlaceholder(locations)
	}
	defer tmpFile.Close()

	// Copier l'image
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return createMapPlaceholder(locations)
	}

	// Charger l'image
	mapImg := canvas.NewImageFromFile(tmpFile.Name())
	mapImg.FillMode = canvas.ImageFillContain
	mapImg.SetMinSize(fyne.NewSize(600, 400))
	return mapImg
}

// createMapPlaceholder crée un placeholder pour la carte
func createMapPlaceholder(locations []string) fyne.CanvasObject {
	placeholderText := fmt.Sprintf("🗺️ Carte\n\n%d localisations", len(locations))
	label := widget.NewLabel(placeholderText)
	label.Alignment = fyne.TextAlignCenter

	// Créer un rectangle pour simuler une carte
	rect := canvas.NewRectangle(color.NRGBA{R: 200, G: 220, B: 240, A: 255})
	rect.SetMinSize(fyne.NewSize(300, 250))

	return container.NewStack(rect, container.NewCenter(label))
}

// geocodeLocation utilise l'API Nominatim pour transformer une adresse en coordonnées
func geocodeLocation(query string) (lat string, lon string, ok bool) {
	endpoint := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", url.QueryEscape(query))
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", "", false
	}
	// Nominatim requiert un User-Agent
	req.Header.Set("User-Agent", "GroupieTracker/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", false
	}

	var results []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return "", "", false
	}
	if len(results) == 0 {
		return "", "", false
	}
	return results[0].Lat, results[0].Lon, true
}
