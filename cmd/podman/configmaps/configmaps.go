package configmaps

import (
	"github.com/containers/podman/v4/cmd/podman/registry"
	"github.com/containers/podman/v4/cmd/podman/validate"
	"github.com/spf13/cobra"
)

var (
	// Command: podman _configmap_
	configmapCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage configmaps",
		Long:  "Manage configmaps",
		RunE:  validate.SubCommandExists,
	}
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Command: configmapCmd,
	})
}
