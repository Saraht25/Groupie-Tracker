// Package ui - mapping.go génère les cartes OpenStreetMap pour afficher les lieux de concert.
// Il convertit les coordonnées géographiques en tuiles de cartes et les affiche avec zoom approprié.
// C'est le module essentiel pour la visualisation spatiale des événements musicaux.
package ui

import (
	"fmt"
	"image/color"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// latLonToTile convertit des lat/lon en indices de tuile OSM (Web Mercator)
func latLonToTile(lat, lon float64, z int) (x, y int) {
	n := math.Pow(2, float64(z))

	latRad := lat * math.Pi / 180.0

	x = int((lon + 180.0) / 360.0 * n)
	y = int((1.0 - math.Log(math.Tan(latRad)+(1.0/math.Cos(latRad)))/math.Pi) / 2.0 * n)

	if x < 0 {
		x = 0
	}
	if x >= int(n) {
		x = int(n) - 1
	}
	if y < 0 {
		y = 0
	}
	if y >= int(n) {
		y = int(n) - 1
	}

	return x, y
}

// createMapPlaceholder affiche un visuel de secours quand la carte n'est pas disponible
func createMapPlaceholder(locations []string) fyne.CanvasObject {
	placeholderText := fmt.Sprintf("🗺️ Carte\n\n%d localisations", len(locations))
	label := widget.NewLabel(placeholderText)
	label.Alignment = fyne.TextAlignCenter

	rect := canvas.NewRectangle(color.NRGBA{R: 200, G: 220, B: 240, A: 255})
	rect.SetMinSize(fyne.NewSize(300, 250))

	return container.NewStack(rect, container.NewCenter(label))
}

// createLocationMapForSingle récupère une tuile OSM pour un lieu unique et l'affiche
func createLocationMapForSingle(location string) fyne.CanvasObject {
	placeholder := widget.NewLabel("🗺️ Chargement...")
	placeholder.Alignment = fyne.TextAlignCenter

	sizeRect := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	sizeRect.SetMinSize(fyne.NewSize(300, 200))

	cont := container.NewStack(sizeRect, placeholder)

	go func() {
		lat, lon, ok := geocodeLocation(location)
		if !ok {
			return
		}

		clat, _ := strconv.ParseFloat(lat, 64)
		clon, _ := strconv.ParseFloat(lon, 64)
		z := 12
		x, y := latLonToTile(clat, clon, z)

		mapURL := fmt.Sprintf(
			"https://tile.openstreetmap.org/%d/%d/%d.png",
			z, x, y,
		)

		req, err := http.NewRequest("GET", mapURL, nil)
		if err != nil {
			return
		}
		req.Header.Set("User-Agent", "GroupieTracker/1.0 (+https://github.com/)")

		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			return
		}
		defer resp.Body.Close()

		tmpFile, err := os.CreateTemp("", "map-loc-*.png")
		if err != nil {
			return
		}
		defer tmpFile.Close()

		_, err = io.Copy(tmpFile, resp.Body)
		if err != nil {
			return
		}

		mapImg := canvas.NewImageFromFile(tmpFile.Name())
		mapImg.FillMode = canvas.ImageFillContain
		mapImg.SetMinSize(fyne.NewSize(300, 200))

		go func() {
			cont.Objects = []fyne.CanvasObject{mapImg}
			cont.Refresh()
		}()
	}()

	return cont
}
