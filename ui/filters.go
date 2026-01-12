package ui

import "fyne.io/fyne/v2/widget"

func FilterOptions() *widget.CheckGroup {
	options := []string{"Création", "Premier album", "Nombre de membres", "Lieu concert"}
	return widget.NewCheckGroup(options, nil)
}
