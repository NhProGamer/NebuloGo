package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func DoesExistFile(filename string) (bool, error) {
	if _, err := os.Stat(filename); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func IsPathAllowed(baseDir, requestedPath string) bool {
	// Nettoyer le chemin demandé
	cleanedPath := filepath.Clean(requestedPath)

	// Construire le chemin absolu pour le chemin demandé
	absRequestedPath, err := filepath.Abs(cleanedPath)
	if err != nil {
		return false
	}

	// Construire le chemin absolu basé sur le répertoire de base
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return false
	}

	// Vérifier que le chemin demandé commence bien par le répertoire de base
	return strings.HasPrefix(absRequestedPath, absBaseDir)
}
