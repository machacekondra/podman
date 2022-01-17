package utils

import (
	"strings"

	"github.com/containers/common/pkg/configmaps"
	"github.com/containers/common/pkg/secrets"
	"github.com/containers/podman/v4/pkg/util"
	"github.com/pkg/errors"
)

func IfPassesSecretsFilter(s secrets.Secret, filters map[string][]string) (bool, error) {
	return ifPassesFilter(s.Name, s.ID, filters)
}

func IfPassesConfigMapsFilter(cm configmaps.ConfigMap, filters map[string][]string) (bool, error) {
	return ifPassesFilter(cm.Name, cm.ID, filters)
}

func ifPassesFilter(name string, id string, filters map[string][]string) (bool, error) {
	result := true
	for key, filterValues := range filters {
		switch strings.ToLower(key) {
		case "name":
			result = util.StringMatchRegexSlice(name, filterValues)
		case "id":
			result = util.StringMatchRegexSlice(id, filterValues)
		default:
			return false, errors.Errorf("invalid filter %q", key)
		}
	}
	return result, nil
}
