# Groupie Tracker - Fyne App

Une application Go utilisant la framework Fyne pour explorer les artistes de musique et leurs informations de concert.

## Corrections et améliorations apportées

### 1. **Gestion des imports**
- ✅ Changé de `import "models"` à `import "Groupie-Tracker/models"`
- ✅ Ajouté Fyne à `go.mod`: `fyne.io/fyne/v2 v2.4.0`
- ✅ Correction des imports Fyne (app, widget, container)

### 2. **Correction du modèle Artist**
- ✅ Ajouté le champ `Locations []string` au modèle Artist
- ✅ Changé `a.ID` (mauvais) à `a.Id` (correct) partout dans le code

### 3. **Structuration du code**
- ✅ `main.go` maintenant propre et simple, appelle `ui.StartApp()`
- ✅ `ui/app.go` initialise correctement la fenêtre Fyne
- ✅ `ui/home.go` affiche les artistes depuis l'API
- ✅ Correction des imports dans tous les fichiers UI et services

### 4. **Logique de recherche et filtrage**
- ✅ Utilisation cohérente de `artist.MatchesSearch()` dans services/search.go
- ✅ Filters appliqués correctement sur les artistes
- ✅ Suggestions de recherche basées sur les propriétés de l'artiste

### 5. **Compilabilité**
- ✅ Le projet compile sans erreurs: `go build -o groupie-tracker.exe`
- ✅ Tous les chemins d'import sont maintenant cohérents

## Installation et exécution

```bash
# Télécharger les dépendances
go mod tidy

# Compiler
go build -o groupie-tracker.exe

# Exécuter
./groupie-tracker.exe
```

## Structure du projet

```
Groupie-Tracker/
├── main.go               # Point d'entrée
├── api/
│   └── api.go           # Appels API
├── models/
│   ├── artist.go        # Modèle Artist avec méthodes
│   ├── date.go
│   ├── location.go
│   └── relation.go
├── ui/
│   ├── app.go           # Initialisation Fyne
│   ├── home.go          # Page d'accueil
│   ├── artist.go        # Cartes d'artistes
│   ├── search.go        # Logique de recherche
│   ├── filters.go       # Application de filtres
│   ├── map.go
│   └── markers.go
└── services/
    ├── search.go        # Service de recherche
    └── filter.go        # Service de filtrage
```

## Améliorations futures possibles

- Intégration d'une véritable carte pour les concerts
- Meilleure gestion des erreurs API
- Cache des données API
- Support complet du filtrage par date/location