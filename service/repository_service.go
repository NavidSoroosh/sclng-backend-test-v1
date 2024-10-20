package service

import (
	"sclng-backend-test-v1/githubclient"
	"strings"
	"sync"
)

type RepositoryService struct {
	Client *githubclient.Client
}

func NewRepositoryService(client *githubclient.Client) *RepositoryService {
	return &RepositoryService{Client: client}
}

// Fetch and filter repositories based on parameters
func (s *RepositoryService) FetchAndProcessRepositories(languageFilter string, licenseFilter string, limit int) ([]githubclient.RepositoryOutput, error) {
	repositories, err := s.Client.FetchRepositories(limit)
	if err != nil {
		return nil, err
	}

	var waitGroup sync.WaitGroup
	var mutex sync.Mutex
	var processedRepos []githubclient.RepositoryOutput

	waitGroup.Add(len(repositories))

	//Process repositories in parallel
	for _, repo := range repositories {
		go func(repo githubclient.Repository) {
			defer waitGroup.Done()

			// Filter by language or license if provided
			license, err := s.Client.FetchLicense(repo.FullName)
			if err != nil {
				return
			}

			if licenseFilter != "" && (license == nil || !strings.EqualFold(license.License.Key, licenseFilter)) {
				return
			}

			languages, err := s.Client.FetchLanguages(repo.FullName)
			if err != nil {
				return
			}

			if languageFilter != "" {
				found := false
				for language := range languages {
					if strings.EqualFold(language, languageFilter) {
						found = true
						break
					}
				}

				if !found {
					return
				}
			}

			// Convert languages to the required format
			languageMap := make(map[string]githubclient.Language)
			for lang, bytes := range languages {
				languageMap[lang] = githubclient.Language{Bytes: bytes}
			}

			mutex.Lock()
			processedRepos = append(processedRepos, githubclient.RepositoryOutput{
				FullName:  repo.FullName,
				Owner:     repo.Owner.Login,
				Name:      repo.Name,
				License:   license.License.Key,
				Languages: languageMap,
			})
			mutex.Unlock()

		}(repo)
	}

	waitGroup.Wait()
	return processedRepos, nil
}
