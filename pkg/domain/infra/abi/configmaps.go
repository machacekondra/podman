package abi

import (
	"context"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/containers/podman/v4/pkg/domain/utils"
	"github.com/pkg/errors"
)

func (ic *ContainerEngine) ConfigMapCreate(ctx context.Context, name string, reader io.Reader, options entities.ConfigMapCreateOptions) (*entities.ConfigMapCreateReport, error) {
	data, _ := ioutil.ReadAll(reader)
	configmapssPath := ic.Libpod.GetConfigMapsStorageDir()
	manager, err := ic.Libpod.ConfigMapsManager()
	if err != nil {
		return nil, err
	}

	// set defaults from config for the case they are not set by an upper layer
	// (-> i.e. tests that talk directly to the api)
	cfg, err := ic.Libpod.GetConfigNoCopy()
	if err != nil {
		return nil, err
	}
	if options.Driver == "" {
		options.Driver = cfg.Secrets.Driver
	}
	if len(options.DriverOpts) == 0 {
		options.DriverOpts = cfg.Secrets.Opts
	}
	if options.DriverOpts == nil {
		options.DriverOpts = make(map[string]string)
	}

	if options.Driver == "file" {
		if _, ok := options.DriverOpts["path"]; !ok {
			options.DriverOpts["path"] = filepath.Join(configmapssPath, "filedriver")
		}
	}

	cmID, err := manager.Store(name, data, options.Driver, options.DriverOpts)
	if err != nil {
		return nil, err
	}
	return &entities.ConfigMapCreateReport{
		ID: cmID,
	}, nil
}

func (ic *ContainerEngine) ConfigMapInspect(ctx context.Context, nameOrIDs []string) ([]*entities.ConfigMapInfoReport, []error, error) {
	manager, err := ic.Libpod.ConfigMapsManager()
	if err != nil {
		return nil, nil, err
	}
	errs := make([]error, 0, len(nameOrIDs))
	reports := make([]*entities.ConfigMapInfoReport, 0, len(nameOrIDs))
	for _, nameOrID := range nameOrIDs {
		secret, err := manager.Lookup(nameOrID)
		if err != nil {
			if errors.Cause(err).Error() == "no such configmap" {
				errs = append(errs, err)
				continue
			} else {
				return nil, nil, errors.Wrapf(err, "error inspecting configmap %s", nameOrID)
			}
		}
		report := &entities.ConfigMapInfoReport{
			ID:        secret.ID,
			CreatedAt: secret.CreatedAt,
			UpdatedAt: secret.CreatedAt,
			Spec: entities.ConfigMapSpec{
				Name: secret.Name,
				Driver: entities.ConfigMapDriverSpec{
					Name:    secret.Driver,
					Options: secret.DriverOptions,
				},
			},
		}
		reports = append(reports, report)
	}

	return reports, errs, nil
}

func (ic *ContainerEngine) ConfigMapList(ctx context.Context, opts entities.ConfigMapListRequest) ([]*entities.ConfigMapInfoReport, error) {
	manager, err := ic.Libpod.ConfigMapsManager()
	if err != nil {
		return nil, err
	}
	configMapList, err := manager.List()
	if err != nil {
		return nil, err
	}
	report := make([]*entities.ConfigMapInfoReport, 0, len(configMapList))
	for _, cm := range configMapList {
		result, err := utils.IfPassesConfigMapsFilter(cm, opts.Filters)
		if err != nil {
			return nil, err
		}
		if result {
			reportItem := entities.ConfigMapInfoReport{
				ID:        cm.ID,
				CreatedAt: cm.CreatedAt,
				UpdatedAt: cm.CreatedAt,
				Spec: entities.ConfigMapSpec{
					Name: cm.Name,
					Driver: entities.ConfigMapDriverSpec{
						Name:    cm.Driver,
						Options: cm.DriverOptions,
					},
				},
			}
			report = append(report, &reportItem)
		}
	}
	return report, nil
}

func (ic *ContainerEngine) ConfigMapRm(ctx context.Context, nameOrIDs []string, options entities.ConfigMapRmOptions) ([]*entities.ConfigMapRmReport, error) {
	var (
		err      error
		toRemove []string
		reports  = []*entities.ConfigMapRmReport{}
	)
	manager, err := ic.Libpod.ConfigMapsManager()
	if err != nil {
		return nil, err
	}
	toRemove = nameOrIDs
	if options.All {
		allSecrs, err := manager.List()
		if err != nil {
			return nil, err
		}
		for _, secr := range allSecrs {
			toRemove = append(toRemove, secr.ID)
		}
	}
	for _, nameOrID := range toRemove {
		deletedID, err := manager.Delete(nameOrID)
		if err == nil || errors.Cause(err).Error() == "no such configmap" {
			reports = append(reports, &entities.ConfigMapRmReport{
				Err: err,
				ID:  deletedID,
			})
			continue
		} else {
			return nil, err
		}
	}

	return reports, nil
}
