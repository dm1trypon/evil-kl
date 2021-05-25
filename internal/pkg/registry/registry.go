package registry

import (
	"fmt"

	logger "github.com/dm1trypon/easy-logger"
	"golang.org/x/sys/windows/registry"
)

/*
Create <Registry> - init Registry structure
	Returns <*Registry>:
		1. Structure pointer
*/
func (r *Registry) Create() *Registry {
	r = &Registry{
		lc: "REGISTRY",
	}

	return r
}

/*
SetStringValue <Registry> - SetStringValue - sets a string value in the registry
	Returns <error>:
		1. error
	Args:
		1. path <string> - path in the registry
		2. name <string> - name of value
		3. value <string> - value
*/
func (r *Registry) SetStringValue(path, name, value string) error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		logger.ErrorJ(r.lc, fmt.Sprint("Failed opening key by path '", path, "': ", err.Error()))
		return err
	}

	defer key.Close()

	if err = key.SetStringValue(name, value); err != nil {
		logger.ErrorJ(r.lc, fmt.Sprint("Failed setting value of '", name, "': ", err.Error()))
		return err
	}

	return nil
}

/*
SetStringValue <Registry> - SetStringValue - gets a string value in the registry
	Returns <string, error>:
		1. value
		2. error
	Args:
		1. path <string> - path in the registry
		2. name <string> - name of value
*/
func (r *Registry) GetStringValue(path, name string) (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		logger.ErrorJ(r.lc, fmt.Sprint("Failed opening key by path '", path, "': ", err.Error()))
		return "", err
	}

	defer key.Close()

	value, _, err := key.GetStringValue(name)
	if err != nil {
		logger.ErrorJ(r.lc, fmt.Sprint("Failed getting value of '", name, "': ", err.Error()))
		return "", err
	}

	return value, nil
}
