# Version 1.0 Release Date: TBD
I would like to add common helper functions/features inspired by the package use in the community. So please, especially before Version 1.0 release, let me know what you would like to see added to the package, but bear in mind the main objective to be a simple wrapper for the API exposed by the GroupMe team.

<br>

# GroupMe API Wrapper
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/densestvoid/groupme?label=version&logo=version&sort=semver)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/densestvoid/groupme)](https://pkg.go.dev/github.com/densestvoid/groupme)
[![codecov](https://codecov.io/gh/densestvoid/groupme/branch/master/graph/badge.svg)](https://codecov.io/gh/densestvoid/groupme)
## Description
The design of this package is meant to be super simple. Wrap the exposed API endpoints [documented](https://dev.groupme.com/docs/v3#v3) by the GroupMe team. While you can achieve the core of this package with cURL, there are some small added features, coupled along with a modern language, that should simplify writing GroupMe [bots](https://dev.groupme.com/bots) and [applications](https://dev.groupme.com/applications).

[*FUTURE*] In addition to the Go package, there is also a CLI application built using this package; all the features are available from the command line.

## Why?
I enjoy programming, I use GroupMe with friends, and I wanted to write a fun add-on application for our group. I happened to start using Go around this time, so it was good practice.

## Example
```golang
package main

import (
	"fmt"

	"github.com/densestvoid/groupme"
)

// This is not a real token. Please find yours by logging
// into the GroupMe development website: https://dev.groupme.com/
const authorizationToken = "0123456789ABCDEF"

// A short program that gets the gets the first 5 groups
// the user is part of, and then the first 10 messages of
// the first group in that list
func main() {
	// Create a new client with your auth token
	client := groupme.NewClient(authorizationToken)

	// Get the groups your user is part of
	groups, err := client.IndexGroups(&groupme.GroupsQuery{
		Page:    0,
		PerPage: 5,
		Omit:    "memberships",
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(groups)

	// Get first 10 messages of the first group
	if len(groups) <= 0 {
		fmt.Println("No groups")
	}

	messages, err := client.IndexMessages(groups[0].ID, &groupme.IndexMessagesQuery{
		Limit: 10,
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(messages)
}
```

## Installation

### Go Package
`go get github.com/densestvoid/groupme`

### [*FUTURE*] CLI

## Contribute
I find the hours I can spend developing personal projects decreasing every year, so I welcome any help I can get. Feel free to tackle any open issues, or if a feature request catches your eye, feel free to reach out to me and we can discuss adding it to the package. However, once version 1.0 is released, I don't foresee much work happening on this project unless the GroupMe API is updated.

## Credits
All credits for the actual platform belong to the GroupMe team; I only used the exposed API they wrote.

## License
GPL-3.0 License © [DensestVoid](https://github.com/densestvoid)
