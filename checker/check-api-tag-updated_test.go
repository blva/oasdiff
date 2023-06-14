package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: Adding a new tag is detecteds
func TestTagAdded(t *testing.T) {
	s1, _ := open("../data/checker/tag_added_base.yaml")
	s2, err := open("../data/checker/tag_added_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Tags = []string{"newTag"}

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APITagUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "api-tag-added",
			Text:        "api tag 'newTag' added",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/api/v1.0/groups",
			Source:      "../data/checker/tag_added_base.yaml",
			OperationId: "createOneGroup",
		}}, errs)
}

// CL: Removing an existing tag is detected
func TestTagRemoved(t *testing.T) {
	s1, _ := open("../data/checker/tag_removed_base.yaml")
	s2, err := open("../data/checker/tag_removed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Tags = []string{}

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APITagUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "api-tag-removed",
			Text:        "api tag 'Test' removed",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/api/v1.0/groups",
			Source:      "../data/checker/tag_removed_base.yaml",
			OperationId: "createOneGroup",
		}}, errs)
}

// CL: Updating an existing tag is detected
func TestTagUpdated(t *testing.T) {
	s1, _ := open("../data/checker/tag_removed_base.yaml")
	s2, err := open("../data/checker/tag_removed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Tags = []string{"newTag"}

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APITagUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 2)
	for cl := range errs {
		require.Equal(t, checker.INFO, errs[cl].Level)
		if errs[cl].Id == "api-tag-removed" {
			require.Equal(t, checker.BackwardCompatibilityError{
				Id:          "api-tag-removed",
				Text:        "api tag 'Test' removed",
				Comment:     "",
				Level:       checker.INFO,
				Operation:   "POST",
				Path:        "/api/v1.0/groups",
				Source:      "../data/checker/tag_removed_base.yaml",
				OperationId: "createOneGroup",
			}, errs[cl])
		}

		if errs[cl].Id == "api-tag-added" {
			require.Equal(t, checker.BackwardCompatibilityError{
				Id:          "api-tag-added",
				Text:        "api tag 'newTag' added",
				Comment:     "",
				Level:       checker.INFO,
				Operation:   "POST",
				Path:        "/api/v1.0/groups",
				Source:      "../data/checker/tag_removed_base.yaml",
				OperationId: "createOneGroup",
			}, errs[cl])
		}
	}
}