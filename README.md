# Groupie-Tracker

Groupie-Tracker est un projet scolaire développé en **Go**, dont
l'objectif est de créer une application permettant de visualiser des
informations sur des artistes de musique à partir de l'API **Groupie
Trackers**.

## Objectifs du projet

-   Comprendre et utiliser une API REST
-   Manipuler des données au format JSON
-   Structurer un projet en Go
-   Créer une interface graphique avec le framework Fyne
-   Mettre en relation différentes sources de données

## Description

L'application permet de consulter des informations sur des artistes
musicaux : - Nom de l'artiste ou du groupe - Membres du groupe - Date de
création - Lieux et dates de concerts

Les données sont récupérées depuis l'API publique Groupie Trackers et
affichées dans une interface graphique développée avec Fyne.

## Technologies utilisées

-   Go (Golang)
-   Fyne
-   API Groupie Trackers
-   JSON

## Structure du projet

-   main.go : point d'entrée de l'application
-   api/ : appels à l'API
-   models/ : structures de données
-   ui/ : interface graphique

## Prérequis

-   Go installé
-   Connexion internet
-   Fyne configuré via go.mod

## Installation et exécution

``` bash
git clone https://github.com/Saraht25/Groupie-Tracker.git
cd Groupie-Tracker
go run main.go
```

## Fonctionnalités

-   Affichage des artistes
-   Informations détaillées
-   Lieux et dates de concerts
-   Interface graphique interactive

## Améliorations possibles

-   Barre de recherche
-   Filtres
-   Amélioration UI
-   Système de favoris

## Contexte scolaire

Projet réalisé dans un cadre pédagogique.
