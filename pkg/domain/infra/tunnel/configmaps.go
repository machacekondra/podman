package tunnel

import (
	"context"
	"io"

	"github.com/containers/podman/v4/pkg/bindings/configmaps"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/containers/podman/v4/pkg/errorhandling"
	"github.com/pkg/errors"
)

func (ic *ContainerEngine) ConfigMapCreate(ctx context.Context, name string, reader io.Reader, options entities.ConfigMapCreateOptions) (*entities.ConfigMapCreateReport, error) {
	opts := new(configmaps.CreateOptions).
		WithDriver(options.Driver).
		WithDriverOpts(options.DriverOpts).
		WithName(name)
	created, err := configmaps.Create(ic.ClientCtx, reader, opts)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (ic *ContainerEngine) ConfigMapInspect(ctx context.Context, nameOrIDs []string) ([]*entities.ConfigMapInfoReport, []error, error) {
	allInspect := make([]*entities.ConfigMapInfoReport, 0, len(nameOrIDs))
	errs := make([]error, 0, len(nameOrIDs))
	for _, name := range nameOrIDs {
		inspected, err := configmaps.Inspect(ic.ClientCtx, name, nil)
		if err != nil {
			errModel, ok := err.(*errorhandling.ErrorModel)
			if !ok {
				return nil, nil, err
			}
			if errModel.ResponseCode == 404 {
				errs = append(errs, errors.Errorf("no such configmap %q", name))
				continue
			}
			return nil, nil, err
		}
		allInspect = append(allInspect, inspected)
	}
	return allInspect, errs, nil
}

func (ic *ContainerEngine) ConfigMapList(ctx context.Context, opts entities.ConfigMapListRequest) ([]*entities.ConfigMapInfoReport, error) {
	options := new(configmaps.ListOptions).WithFilters(opts.Filters)
	secrs, _ := configmaps.List(ic.ClientCtx, options)
	return secrs, nil
}

func (ic *ContainerEngine) ConfigMapRm(ctx context.Context, nameOrIDs []string, options entities.ConfigMapRmOptions) ([]*entities.ConfigMapRmReport, error) {
	allRm := make([]*entities.ConfigMapRmReport, 0, len(nameOrIDs))
	if options.All {
		allSecrets, err := configmaps.List(ic.ClientCtx, nil)
		if err != nil {
			return nil, err
		}
		for _, secret := range allSecrets {
			allRm = append(allRm, &entities.ConfigMapRmReport{
				Err: configmaps.Remove(ic.ClientCtx, secret.ID),
				ID:  secret.ID,
			})
		}
		return allRm, nil
	}
	for _, name := range nameOrIDs {
		secret, err := configmaps.Inspect(ic.ClientCtx, name, nil)
		if err != nil {
			errModel, ok := err.(*errorhandling.ErrorModel)
			if !ok {
				return nil, err
			}
			if errModel.ResponseCode == 404 {
				allRm = append(allRm, &entities.ConfigMapRmReport{
					Err: errors.Errorf("no configmap with name or id %q: no such configmap ", name),
					ID:  "",
				})
				continue
			}
		}
		allRm = append(allRm, &entities.ConfigMapRmReport{
			Err: configmaps.Remove(ic.ClientCtx, name),
			ID:  secret.ID,
		})
	}
	return allRm, nil
}
