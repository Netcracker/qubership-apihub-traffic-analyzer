package view

type RestOperationChange struct {
	Path   string   `json:"path"`
	Method string   `json:"method"`
	Tags   []string `json:"tags,omitempty"`
}

type RestOperationMetadata struct {
	Path   string   `json:"path"`
	Method string   `json:"method"`
	Tags   []string `json:"tags,omitempty"`
}

type RestOperationSingleView struct {
	SingleOperationView
	RestOperationMetadata
}

type RestOperationView struct {
	OperationListView
	RestOperationMetadata
}
type DeprecatedRestOperationView struct {
	DeprecatedOperationView
	RestOperationMetadata
}

type OperationSummary struct {
	Endpoints  int `json:"endpoints"`
	Deprecated int `json:"deprecated"`
	Created    int `json:"created"`
	Deleted    int `json:"deleted"`
}

type RestOperationChangelogView struct {
	OperationChangelogView
	RestOperationChange
}

type RestOperationComparisonChangelogView_deprecated struct {
	OperationComparisonChangelogView_deprecated
	RestOperationChange
}
type RestOperationComparisonChangelogView struct {
	OperationComparisonChangelogView
	RestOperationChange
}

type RestOperationComparisonChangesView struct {
	OperationComparisonChangesView
	RestOperationChange
}

type RestOperations struct {
	Operations []RestOperationView          `json:"operations"`
	Packages   map[string]PackageVersionRef `json:"packages,omitempty"`
}
