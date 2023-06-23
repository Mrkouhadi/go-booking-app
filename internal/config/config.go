package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
	"github.com/mrkouhadi/go-booking-app/internal/models"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	MailChan      chan models.MailData
}
