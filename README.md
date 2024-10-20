# Backend Technical Test at Scalingo - GitHub Repositories API Service

## Overview

This project is a Go-based microservice that interacts with the GitHub API to fetch public repositories and their language details. It is designed to be modular, maintainable, and scalable by separating concerns and using parallel processing to optimize performance. The service allows users to filter repositories based on programming languages and license types.

## Execution

```
docker compose up
```

### Access the API:

The server will be running on http://localhost:5000/repos. You can make a request with filters like language and license:

### Example request:

```
http://localhost:5000/repos?language=Ruby&license=mit
```

Requests parameters are not case sensitive.

### Example output:

```
{
  "repositories": [
    {
      "full_name": "userName/repositoryName",
      "owner": "userName",
      "repository": "repositoryName",
      "license": "mit",
      "languages": {
        "JavaScript": {
          "bytes": 121223
        },
        "Ruby": {
          "bytes": 674
        }
      }
    }
  ]
}
```

## Project Architecture

To achieve maintainability, the project is structured into packages:

### githubclient/client.go:

This file manages GitHub API interactions, including handling authenticated requests and fetching repositories or language data.
It abstracts the HTTP request and response logic into reusable functions, making the service more modular and easier to test.

### githubclient/models.go:

Defines the data structures required to map GitHub API responses to Go objects.
For example, Repository maps to the repository data fetched from the API, and RepositoryOutput provides a refined structure that includes language and owner details for the response.

### service/repository_service.go:

This file contains the business logic for fetching repositories, filtering them based on language and license, and enriching them with language details.
The service is designed to handle API requests concurrently using goroutines to improve performance, especially when fetching language data for multiple repositories.

### main.go:

The main entry point of the application, responsible for initializing the service and setting up the HTTP server.
This file ties everything together by integrating the repository_service.go functions and exposing the /repos API endpoint.

## Design Decisions

### Separation of Concerns:

The code is organized into distinct packages (githubclient and service) to follow the principles of separation of concerns.
This makes the project easier to extend, test, and maintain.

### Parallel Processing:

Fetching language data for repositories is done in parallel using goroutines, improving the overall performance when processing large sets of repositories.

### Error Handling:

Errors are logged but do not crash the application. For example, if a single repository's language data cannot be fetched, the process continues for the remaining repositories.
