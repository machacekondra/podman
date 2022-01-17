package entities

import (
	"time"

	"github.com/containers/podman/v3/pkg/errorhandling"
)

type ConfigMapCreateReport struct {
	ID string
}

type ConfigMapCreateOptions struct {
	Driver     string
	DriverOpts map[string]string
}

type ConfigMapListRequest struct {
	Filters map[string][]string
}

type ConfigMapListReport struct {
	ID        string
	Name      string
	Driver    string
	CreatedAt string
	UpdatedAt string
}

type ConfigMapRmOptions struct {
	All bool
}

type ConfigMapRmReport struct {
	ID  string
	Err error
}

type ConfigMapInfoReport struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Spec      ConfigMapSpec
}

type ConfigMapInfoReportCompat struct {
	ConfigMapInfoReport
	Version ConfigMapVersion
}

type ConfigMapVersion struct {
	Index int
}

type ConfigMapSpec struct {
	Name   string
	Driver ConfigMapDriverSpec
}

type ConfigMapDriverSpec struct {
	Name    string
	Options map[string]string
}

// swagger:model ConfigMapCreate
type ConfigMapCreateRequest struct {
	// User-defined name of the config map.
	Name string
	// Base64-url-safe-encoded (RFC 4648) data to store as config map.
	Data string
	// Driver represents a driver (default "file")
	Driver ConfigMapDriverSpec
}

// ConfigMap create response
// swagger:response ConfigMapCreateResponse
type SwagConfigMapCreateResponse struct {
	// in:body
	Body struct {
		ConfigMapCreateReport
	}
}

// ConfigMap list response
// swagger:response ConfigMapListResponse
type SwagConfigMapListResponse struct {
	// in:body
	Body []*ConfigMapInfoReport
}

// ConfigMap list response
// swagger:response ConfigMapListCompatResponse
type SwagConfigMapListCompatResponse struct {
	// in:body
	Body []*ConfigMapInfoReportCompat
}

// ConfigMap inspect response
// swagger:response ConfigMapInspectResponse
type SwagConfigMapInspectResponse struct {
	// in:body
	Body ConfigMapInfoReport
}

// ConfigMap inspect compat
// swagger:response ConfigMapInspectCompatResponse
type SwagConfigMapInspectCompatResponse struct {
	// in:body
	Body ConfigMapInfoReportCompat
}

// No such configmap
// swagger:response NoSuchConfigMap
type SwagErrNoSuchConfigMap struct {
	// in:body
	Body struct {
		errorhandling.ErrorModel
	}
}

// ConfigMap in use
// swagger:response ConfigMapInUse
type SwagErrConfigMapInUse struct {
	// in:body
	Body struct {
		errorhandling.ErrorModel
	}
}
