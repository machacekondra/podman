package configmaps

import (
	"context"
	"io"
	"net/http"

	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/domain/entities"
)

// List returns information about existing configmaps in the form of a slice.
func List(ctx context.Context, options *ListOptions) ([]*entities.ConfigMapInfoReport, error) {
	var (
		secrs []*entities.ConfigMapInfoReport
	)
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	params, err := options.ToParams()
	if err != nil {
		return nil, err
	}
	response, err := conn.DoRequest(ctx, nil, http.MethodGet, "/secrets/json", params, nil)
	if err != nil {
		return secrs, err
	}
	defer response.Body.Close()

	return secrs, response.Process(&secrs)
}

// Inspect returns low-level information about a configmap.
func Inspect(ctx context.Context, nameOrID string, options *InspectOptions) (*entities.ConfigMapInfoReport, error) {
	var (
		inspect *entities.ConfigMapInfoReport
	)
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	response, err := conn.DoRequest(ctx, nil, http.MethodGet, "/secrets/%s/json", nil, nil, nameOrID)
	if err != nil {
		return inspect, err
	}
	defer response.Body.Close()

	return inspect, response.Process(&inspect)
}

// Remove removes a configmap from storage
func Remove(ctx context.Context, nameOrID string) error {
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return err
	}

	response, err := conn.DoRequest(ctx, nil, http.MethodDelete, "/secrets/%s", nil, nil, nameOrID)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return response.Process(nil)
}

// Create creates a configmap given some data
func Create(ctx context.Context, reader io.Reader, options *CreateOptions) (*entities.ConfigMapCreateReport, error) {
	var (
		create *entities.ConfigMapCreateReport
	)
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	params, err := options.ToParams()
	if err != nil {
		return nil, err
	}

	response, err := conn.DoRequest(ctx, reader, http.MethodPost, "/secrets/create", params, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return create, response.Process(&create)
}
