package app

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

// Logger global
var instance *logrus.Logger
var once sync.Once

// CustomFormatter est un formateur personnalisé pour logrus
type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("02-01-2006 15:04")
    level := entry.Level.String()
    log := fmt.Sprintf("%s [%s] [%s] %s %s\n", timestamp, level, entry.Data["method"], entry.Data["path"], entry.Message)
    return []byte(log), nil
}

// GetLogger retourne l'instance unique du logger
func GetLogger() *logrus.Logger {
    once.Do(func() {
        // Initialiser le logger
        instance = logrus.New()
        instance.SetFormatter(&CustomFormatter{})

        // Ouvrir un fichier pour y écrire les logs
        logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err != nil {
            logrus.Fatalf("Erreur lors de l'ouverture du fichier de log : %v", err)
        }

        // Définir les sorties : console + fichier
        instance.SetOutput(io.MultiWriter(os.Stdout, logFile))
    })
    return instance
}