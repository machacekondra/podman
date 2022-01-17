package compat

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/containers/podman/v4/libpod"
	"github.com/containers/podman/v4/pkg/api/handlers/utils"
	api "github.com/containers/podman/v4/pkg/api/types"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/containers/podman/v4/pkg/domain/infra/abi"
	"github.com/containers/podman/v4/pkg/util"
	"github.com/pkg/errors"
)

func ListConfigMaps(w http.ResponseWriter, r *http.Request) {
	var (
		runtime = r.Context().Value(api.RuntimeKey).(*libpod.Runtime)
	)
	filtersMap, err := util.PrepareFilters(r)
	if err != nil {
		utils.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError,
			errors.Wrapf(err, "failed to parse parameters for %s", r.URL.String()))
		return
	}
	ic := abi.ContainerEngine{Libpod: runtime}
	listOptions := entities.ConfigMapListRequest{
		Filters: *filtersMap,
	}
	reports, err := ic.ConfigMapList(r.Context(), listOptions)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
	if utils.IsLibpodRequest(r) {
		utils.WriteResponse(w, http.StatusOK, reports)
		return
	}
	// Docker compat expects a version field that increments when the configmap is updated
	// We currently can't update a configmap, so we default the version to 1
	compatReports := make([]entities.ConfigMapInfoReportCompat, 0, len(reports))
	for _, report := range reports {
		compatRep := entities.ConfigMapInfoReportCompat{
			ConfigMapInfoReport: *report,
			Version:             entities.ConfigMapVersion{Index: 1},
		}
		compatReports = append(compatReports, compatRep)
	}
	utils.WriteResponse(w, http.StatusOK, compatReports)
}

func InspectConfigMap(w http.ResponseWriter, r *http.Request) {
	var (
		runtime = r.Context().Value(api.RuntimeKey).(*libpod.Runtime)
	)
	name := utils.GetName(r)
	names := []string{name}
	ic := abi.ContainerEngine{Libpod: runtime}
	reports, errs, err := ic.ConfigMapInspect(r.Context(), names)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
	if len(errs) > 0 {
		utils.ConfigMapNotFound(w, name, errs[0])
		return
	}
	if len(reports) < 1 {
		utils.InternalServerError(w, err)
		return
	}
	if utils.IsLibpodRequest(r) {
		utils.WriteResponse(w, http.StatusOK, reports[0])
		return
	}
	// Docker compat expects a version field that increments when the configmap is updated
	// We currently can't update a configmap, so we default the version to 1
	compatReport := entities.ConfigMapInfoReportCompat{
		ConfigMapInfoReport: *reports[0],
		Version:             entities.ConfigMapVersion{Index: 1},
	}
	utils.WriteResponse(w, http.StatusOK, compatReport)
}

func RemoveConfigMap(w http.ResponseWriter, r *http.Request) {
	var (
		runtime = r.Context().Value(api.RuntimeKey).(*libpod.Runtime)
	)

	opts := entities.ConfigMapRmOptions{}
	name := utils.GetName(r)
	ic := abi.ContainerEngine{Libpod: runtime}
	reports, err := ic.ConfigMapRm(r.Context(), []string{name}, opts)
	if err != nil {
		utils.InternalServerError(w, err)
		return
	}
	if reports[0].Err != nil {
		utils.ConfigMapNotFound(w, name, reports[0].Err)
		return
	}
	utils.WriteResponse(w, http.StatusNoContent, nil)
}

func CreateConfigMap(w http.ResponseWriter, r *http.Request) {
	var (
		runtime = r.Context().Value(api.RuntimeKey).(*libpod.Runtime)
	)
	opts := entities.ConfigMapCreateOptions{}
	createParams := struct {
		*entities.ConfigMapCreateRequest
		Labels map[string]string `schema:"labels"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&createParams); err != nil {
		utils.Error(w, "Something went wrong.", http.StatusInternalServerError, errors.Wrap(err, "Decode()"))
		return
	}
	if len(createParams.Labels) > 0 {
		utils.Error(w, "labels not supported", http.StatusBadRequest,
			errors.Wrapf(errors.New("bad parameter"), "labels not supported"))
		return
	}

	decoded, _ := base64.StdEncoding.DecodeString(createParams.Data)
	reader := bytes.NewReader(decoded)
	opts.Driver = createParams.Driver.Name

	ic := abi.ContainerEngine{Libpod: runtime}
	report, err := ic.ConfigMapCreate(r.Context(), createParams.Name, reader, opts)
	if err != nil {
		if errors.Cause(err).Error() == "configmap name in use" {
			utils.Error(w, "name conflicts with an existing object", http.StatusConflict, err)
			return
		}
		utils.InternalServerError(w, err)
		return
	}
	utils.WriteResponse(w, http.StatusOK, report)
}

func UpdateConfigMap(w http.ResponseWriter, r *http.Request) {
	utils.Error(w, fmt.Sprintf("unsupported endpoint: %v", r.Method), http.StatusNotImplemented, errors.New("update is not supported"))
}
