package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
	"github.com/tegveer-singh123/bookings/internal/models"
)

// it holds application configurations
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	Mailchan      chan models.MailData
}
