package configmaps

import (
	"context"
	"errors"
	"fmt"

	"github.com/containers/podman/v4/cmd/podman/common"
	"github.com/containers/podman/v4/cmd/podman/registry"
	"github.com/containers/podman/v4/cmd/podman/utils"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/spf13/cobra"
)

var (
	rmCmd = &cobra.Command{
		Use:               "rm [options] CONFIGMAP [CONFIGMAP...]",
		Short:             "Remove one or more configmaps",
		RunE:              rm,
		ValidArgsFunction: common.AutocompleteConfigmaps,
		Example:           "podman config rm mycm1 mycm2",
	}
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Command: rmCmd,
		Parent:  configmapCmd,
	})
	flags := rmCmd.Flags()
	flags.BoolVarP(&rmOptions.All, "all", "a", false, "Remove all configmaps")
}

var (
	rmOptions = entities.ConfigMapRmOptions{}
)

func rm(cmd *cobra.Command, args []string) error {
	var (
		errs utils.OutputErrors
	)
	if (len(args) > 0 && rmOptions.All) || (len(args) < 1 && !rmOptions.All) {
		return errors.New("`podman config rm` requires one argument, or the --all flag")
	}
	responses, err := registry.ContainerEngine().ConfigMapRm(context.Background(), args, rmOptions)
	if err != nil {
		return err
	}
	for _, r := range responses {
		if r.Err == nil {
			fmt.Println(r.ID)
		} else {
			errs = append(errs, r.Err)
		}
	}
	return errs.PrintErrors()
}
