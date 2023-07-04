package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: Changing required response property to optional
func TestResponsePropertyBecameRequiredlCheck(t *testing.T) {
	s1, err := open("../data/checker/response_property_became_optional_revision.yaml")
	require.Empty(t, err)
	s2, err := open("../data/checker/response_property_became_optional_base.yaml")
	require.Empty(t, err)
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyBecameRequiredCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "response-property-became-required",
			Text:        "the response property 'data/name' became required for the status '200'",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/api/v1.0/groups",
			Source:      "../data/checker/response_property_became_optional_base.yaml",
			OperationId: "createOneGroup",
		},
	}, errs)
}

// CL: Changing required response write-only property to optional
func TestResponseWriteOnlyPropertyBecameRequiredCheck(t *testing.T) {
	s1, err := open("../data/checker/response_property_became_optional_revision.yaml")
	require.Empty(t, err)
	s2, err := open("../data/checker/response_property_became_optional_base.yaml")
	require.Empty(t, err)
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	s1.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["name"].Value.WriteOnly = true

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyBecameRequiredCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "response-write-only-property-became-required",
			Text:        "the response write-only property 'data/name' became required for the status '200'",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/api/v1.0/groups",
			Source:      "../data/checker/response_property_became_optional_base.yaml",
			OperationId: "createOneGroup",
		},
	}, errs)
}
