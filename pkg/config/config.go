package config

import "text/template"

// it holds application configurations
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InProduction  bool
}
