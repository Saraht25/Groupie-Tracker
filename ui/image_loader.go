// Package ui - image_loader.go fournit le chargement asynchrone des images avec placeholder.
// Il télécharge les images des artistes en arrière-plan tandis qu'un placeholder s'affiche.
// C'est crucial pour maintenir une UI réactive même lors du téléchargement d'images volumineuses.
package ui

import (
	"image/color"
	"io"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// LoadImageAsync charge une image de manière asynchrone avec un placeholder
func LoadImageAsync(imageURL string, width, height float32) fyne.CanvasObject {
	placeholder := widget.NewLabel("⏳")
	placeholder.Alignment = fyne.TextAlignCenter

	sizeRect := canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	sizeRect.SetMinSize(fyne.NewSize(width, height))

	imgContainer := container.NewStack(sizeRect, placeholder)

	go func() {
		resp, err := http.Get(imageURL)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return
		}

		tmpFile, err := os.CreateTemp("", "artist-*.jpg")
		if err != nil {
			return
		}
		defer tmpFile.Close()

		_, err = io.Copy(tmpFile, resp.Body)
		if err != nil {
			return
		}

		img := canvas.NewImageFromFile(tmpFile.Name())
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(width, height))

		go func() {
			imgContainer.Objects = []fyne.CanvasObject{sizeRect, img}
			imgContainer.Refresh()
		}()
	}()

	return imgContainer
}

// loadDetailImage télécharge et affiche l'image principale d'un artiste
func loadDetailImage(imageURL string) fyne.CanvasObject {
	if imageURL == "" {
		return widget.NewLabel("Image indisponible")
	}

	resp, err := http.Get(imageURL)
	if err != nil {
		return widget.NewLabel("Erreur de chargement")
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "artist-detail-*.jpg")
	if err != nil {
		return widget.NewLabel("Erreur de cache")
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return widget.NewLabel("Erreur de copie")
	}

	img := canvas.NewImageFromFile(tmpFile.Name())
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(200, 200))
	return img
}
