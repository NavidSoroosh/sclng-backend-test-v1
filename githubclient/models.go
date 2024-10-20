package githubclient

// Repository represents the raw data for each repository retrieved from GitHub.
type Repository struct {
	FullName string `json:"full_name"`
	Owner    struct {
		Login string `json:"login"`
	} `json:"owner"`
	Name string `json:"name"`
}

// License represents the raw data for each repository license retrieved from GitHub.
type License struct {
	License *struct {
		Key string `json:"key"`
	} `json:"license"`
}

// Language represents the language details for a repository.
type Language struct {
	Bytes int `json:"bytes"`
}

// RepositoryOutput represents the final structure of a repository with its languages and other details.
type RepositoryOutput struct {
	FullName  string              `json:"full_name"`
	Owner     string              `json:"owner"`
	Name      string              `json:"repository"`
	License   string              `json:"license"`
	Languages map[string]Language `json:"languages"`
}
