// Code generated by go generate; DO NOT EDIT.
package configmaps

import (
	"net/url"

	"github.com/containers/podman/v4/pkg/bindings/internal/util"
)

// Changed returns true if named field has been set
func (o *CreateOptions) Changed(fieldName string) bool {
	return util.Changed(o, fieldName)
}

// ToParams formats struct fields to be passed to API service
func (o *CreateOptions) ToParams() (url.Values, error) {
	return util.ToParams(o)
}

// WithName set field Name to given value
func (o *CreateOptions) WithName(value string) *CreateOptions {
	o.Name = &value
	return o
}

// GetName returns value of field Name
func (o *CreateOptions) GetName() string {
	if o.Name == nil {
		var z string
		return z
	}
	return *o.Name
}

// WithDriver set field Driver to given value
func (o *CreateOptions) WithDriver(value string) *CreateOptions {
	o.Driver = &value
	return o
}

// GetDriver returns value of field Driver
func (o *CreateOptions) GetDriver() string {
	if o.Driver == nil {
		var z string
		return z
	}
	return *o.Driver
}

// WithDriverOpts set field DriverOpts to given value
func (o *CreateOptions) WithDriverOpts(value map[string]string) *CreateOptions {
	o.DriverOpts = value
	return o
}

// GetDriverOpts returns value of field DriverOpts
func (o *CreateOptions) GetDriverOpts() map[string]string {
	if o.DriverOpts == nil {
		var z map[string]string
		return z
	}
	return o.DriverOpts
}
