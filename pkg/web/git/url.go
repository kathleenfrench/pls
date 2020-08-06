package git

import (
	"strings"
)

// ExtractOrganizationAndRepoNameFromRepoURL parses the organization/user and repository name from the repo URL - the usecase here is when results are returned from search and information about the repository/organization aren't available absent parsing this value and/or making an additional API call
// example RepositoryURL: "https://api.github.com/repos/counterThreat/chess_app",
func ExtractOrganizationAndRepoNameFromRepoURL(url string) (organization string, repo string) {
	splitAfter := strings.SplitAfter(url, "repos/")
	orgSplit := strings.Split(splitAfter[1], "/")
	return orgSplit[0], orgSplit[1]
}
