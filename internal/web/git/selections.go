package gitpls

const (
	// openInBrowser is a selection choice
	openInBrowser = "Open in Browser"

	// cloneRepo is a selection choice
	cloneRepo = "Clone Repo"

	// exitSelections is a selection choice
	exitSelections = "Exit"

	// getOrganizationRepos is a selection choice
	getOrganizationRepos = "Get Organization Repositories"

	// openDiff is a selection choice
	openDiff = "Open Diff"

	// readBodyText is a selection choice
	readBodyText = "See Description"

	// returnToMenu is a selection choice
	returnToMenu = "Return to Menu"

	// ------------ editSelection is a selection choice
	editSelection = "Edit"
	// ---- subchoices
	editTitle          = "Title"
	editBody           = "Body"
	editState          = "State"
	editReadyForReview = "Mark as Ready for Review" // only applicable to 'draft' PRs
	// ------ state changes
	stateOpen   = "open"
	stateClosed = "closed"

	// ------------ mergeSelection is a selection choice
	mergeSelection = "Merge"
	// ------ merge options
	mergeSquash   = "squash"
	mergeRebase   = "rebase"
	mergeStraight = "merge"
)
