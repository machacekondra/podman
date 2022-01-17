package server

import (
	"net/http"

	"github.com/containers/podman/v4/pkg/api/handlers/compat"
	"github.com/gorilla/mux"
)

func (s *APIServer) registerConfigMapHandlers(r *mux.Router) error {
	// swagger:operation POST /libpod/configmaps/create libpod ConfigMapCreateLibpod
	// ---
	// tags:
	//  - configmaps
	// summary: Create a configmap
	// parameters:
	//   - in: query
	//     name: name
	//     type: string
	//     description: User-defined name of the configmap.
	//     required: true
	//   - in: query
	//     name: driver
	//     type: string
	//     description: ConfigMap driver
	//     default: "file"
	//   - in: body
	//     name: request
	//     description: ConfigMap
	//     schema:
	//       type: string
	// produces:
	// - application/json
	// responses:
	//   '201':
	//     $ref: "#/responses/ConfigMapCreateResponse"
	//   '500':
	//      "$ref": "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/configmaps/create"), s.APIHandler(compat.CreateConfigMap)).Methods(http.MethodPost)
	// swagger:operation GET /libpod/configmaps/json libpod ConfigMapListLibpod
	// ---
	// tags:
	//  - configmaps
	// summary: List configmaps
	// description: Returns a list of configmaps
	// parameters:
	//  - in: query
	//    name: filters
	//    type: string
	//    description: |
	//      JSON encoded value of the filters (a `map[string][]string`) to process on the configmaps list. Currently available filters:
	//        - `name=[name]` Matches configmaps name (accepts regex).
	//        - `id=[id]` Matches for full or partial ID.
	// produces:
	// - application/json
	// parameters:
	// responses:
	//   '200':
	//     "$ref": "#/responses/ConfigMapListResponse"
	//   '500':
	//      "$ref": "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/configmaps/json"), s.APIHandler(compat.ListConfigMaps)).Methods(http.MethodGet)
	// swagger:operation GET /libpod/configmaps/{name}/json libpod ConfigMapInspectLibpod
	// ---
	// tags:
	//  - configmaps
	// summary: Inspect configmap
	// parameters:
	//  - in: path
	//    name: name
	//    type: string
	//    required: true
	//    description: the name or ID of the configmap
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     "$ref": "#/responses/ConfigMapInspectResponse"
	//   '404':
	//     "$ref": "#/responses/NoSuchConfigMap"
	//   '500':
	//     "$ref": "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/configmaps/{name}/json"), s.APIHandler(compat.InspectConfigMap)).Methods(http.MethodGet)
	// swagger:operation DELETE /libpod/configmaps/{name} libpod ConfigMapDeleteLibpod
	// ---
	// tags:
	//  - configmaps
	// summary: Remove configmap
	// parameters:
	//  - in: path
	//    name: name
	//    type: string
	//    required: true
	//    description: the name or ID of the configmap
	//  - in: query
	//    name: all
	//    type: boolean
	//    description: Remove all configmaps
	//    default: false
	// produces:
	// - application/json
	// responses:
	//   '204':
	//     description: no error
	//   '404':
	//     "$ref": "#/responses/NoSuchConfigMap"
	//   '500':
	//     "$ref": "#/responses/InternalError"
	r.Handle(VersionedPath("/libpod/configmaps/{name}"), s.APIHandler(compat.RemoveConfigMap)).Methods(http.MethodDelete)

	/*
	 * Docker compatibility endpoints
	 */
	// swagger:operation GET /configmaps compat ConfigMapList
	// ---
	// tags:
	//  - configmaps (compat)
	// summary: List configmaps
	// description: Returns a list of configmaps
	// parameters:
	//  - in: query
	//    name: filters
	//    type: string
	//    description: |
	//      JSON encoded value of the filters (a `map[string][]string`) to process on the configmaps list. Currently available filters:
	//        - `name=[name]` Matches configmaps name (accepts regex).
	//        - `id=[id]` Matches for full or partial ID.
	// produces:
	// - application/json
	// parameters:
	// responses:
	//   '200':
	//     "$ref": "#/responses/ConfigMapListCompatResponse"
	//   '500':
	//      "$ref": "#/responses/InternalError"
	r.Handle(VersionedPath("/configmaps"), s.APIHandler(compat.ListConfigMaps)).Methods(http.MethodGet)
	r.Handle("/configmaps", s.APIHandler(compat.ListConfigMaps)).Methods(http.MethodGet)
	// swagger:operation POST /configmaps/create compat ConfigMapCreate
	// ---
	// tags:
	//  - configmaps (compat)
	// summary: Create a configmap
	// parameters:
	//  - in: body
	//    name: create
	//    description: |
	//      attributes for creating a configmap
	//    schema:
	//      $ref: "#/definitions/ConfigMapCreate"
	// produces:
	// - application/json
	// responses:
	//   '201':
	//     $ref: "#/responses/ConfigMapCreateResponse"
	//   '409':
	//     "$ref": "#/responses/ConfigMapInUse"
	//   '500':
	//      "$ref": "#/responses/InternalError"
	r.Handle(VersionedPath("/configmaps/create"), s.APIHandler(compat.CreateConfigMap)).Methods(http.MethodPost)
	r.Handle("/configmaps/create", s.APIHandler(compat.CreateConfigMap)).Methods(http.MethodPost)
	// swagger:operation GET /configmaps/{name} compat ConfigMapInspect
	// ---
	// tags:
	//  - configmaps (compat)
	// summary: Inspect configmap
	// parameters:
	//  - in: path
	//    name: name
	//    type: string
	//    required: true
	//    description: the name or ID of the configmap
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     "$ref": "#/responses/ConfigMapInspectCompatResponse"
	//   '404':
	//     "$ref": "#/responses/NoSuchConfigMap"
	//   '500':
	//     "$ref": "#/responses/InternalError"
	r.Handle(VersionedPath("/configmaps/{name}"), s.APIHandler(compat.InspectConfigMap)).Methods(http.MethodGet)
	r.Handle("/configmaps/{name}", s.APIHandler(compat.InspectConfigMap)).Methods(http.MethodGet)
	// swagger:operation DELETE /configmaps/{name} compat ConfigMapDelete
	// ---
	// tags:
	//  - configmaps (compat)
	// summary: Remove configmap
	// parameters:
	//  - in: path
	//    name: name
	//    type: string
	//    required: true
	//    description: the name or ID of the configmap
	// produces:
	// - application/json
	// responses:
	//   '204':
	//     description: no error
	//   '404':
	//     "$ref": "#/responses/NoSuchConfigMap"
	//   '500':
	//     "$ref": "#/responses/InternalError"
	r.Handle(VersionedPath("/configmaps/{name}"), s.APIHandler(compat.RemoveConfigMap)).Methods(http.MethodDelete)
	r.Handle("/configmaps/{name}", s.APIHandler(compat.RemoveConfigMap)).Methods(http.MethodDelete)

	r.Handle(VersionedPath("/configmaps/{name}/update"), s.APIHandler(compat.UpdateConfigMap)).Methods(http.MethodPost)
	r.Handle("/configmaps/{name}/update", s.APIHandler(compat.UpdateConfigMap)).Methods(http.MethodPost)
	return nil
}
