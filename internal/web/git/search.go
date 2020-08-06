package gitpls

// MetaGetterFlags are used for top-level search preferences
type MetaGetterFlags struct {
	PerPage int // default: 100
	Page    int // if you want to query a specific page
	/**
	// How to sort the search results. Possible values are:
	//   - for repositories: stars, fork, updated
	//   - for commits: author-date, committer-date
	//   - for code: indexed
	//   - for issues: comments, created, updated
	//   - for users: followers, repositories, joined
	//
	// Default is to sort by best match.
	// ref: https://github.com/google/go-github/blob/master/github/search.go
	*/
	SortBy            string
	Order             string // asc, desc (default: desc)
	TextMatchMetadata bool   // fetch text match metadata with a query

	UseEnterpriseAccount bool
}
