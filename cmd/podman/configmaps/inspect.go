package configmaps

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/containers/common/pkg/report"
	"github.com/containers/podman/v4/cmd/podman/common"
	"github.com/containers/podman/v4/cmd/podman/registry"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	inspectCmd = &cobra.Command{
		Use:               "inspect [options] CONFIGMAP [CONFIGMAP...]",
		Short:             "Inspect a configmap",
		Long:              "Display detail information on one or more configmaps",
		RunE:              inspect,
		Example:           "podman config inspect MYCONFIGMAP",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: common.AutocompleteConfigmaps,
	}
)

var format string

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Command: inspectCmd,
		Parent:  configmapCmd,
	})
	flags := inspectCmd.Flags()
	formatFlagName := "format"
	flags.StringVar(&format, formatFlagName, "", "Format volume output using Go template")
	_ = inspectCmd.RegisterFlagCompletionFunc(formatFlagName, common.AutocompleteFormat(entities.ConfigMapInfoReport{}))
}

func inspect(cmd *cobra.Command, args []string) error {
	inspected, errs, _ := registry.ContainerEngine().ConfigMapInspect(context.Background(), args)

	// always print valid list
	if len(inspected) == 0 {
		inspected = []*entities.ConfigMapInfoReport{}
	}

	if cmd.Flags().Changed("format") {
		row := report.NormalizeFormat(format)
		formatted := report.EnforceRange(row)

		tmpl, err := report.NewTemplate("inspect").Parse(formatted)
		if err != nil {
			return err
		}

		w, err := report.NewWriterDefault(os.Stdout)
		if err != nil {
			return err
		}
		defer w.Flush()
		tmpl.Execute(w, inspected)
	} else {
		buf, err := json.MarshalIndent(inspected, "", "    ")
		if err != nil {
			return err
		}
		fmt.Println(string(buf))
	}

	if len(errs) > 0 {
		if len(errs) > 1 {
			for _, err := range errs[1:] {
				fmt.Fprintf(os.Stderr, "error inspecting configmap: %v\n", err)
			}
		}
		return errors.Errorf("error inspecting configmap: %v", errs[0])
	}
	return nil
}
