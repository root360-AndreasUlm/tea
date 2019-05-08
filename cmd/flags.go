// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"

	"github.com/urfave/cli"
)

// create global variables for global Flags to simplify
// access to the options without requiring cli.Context
var (
	loginValue  string
	repoValue   string
	outputValue string
)

// LoginFlag provides flag to specify tea login profile
var LoginFlag = cli.StringFlag{
	Name:        "login, l",
	Usage:       "Indicate one login, optional when inside a gitea repository",
	Destination: &loginValue,
}

// RepoFlag provides flag to specify repository
var RepoFlag = cli.StringFlag{
	Name:        "repo, r",
	Usage:       "Indicate one repository, optional when inside a gitea repository",
	Destination: &repoValue,
}

// OutputFlag provides flag to specify output type
var OutputFlag = cli.StringFlag{
	Name:        "output, o",
	Usage:       "Indicate one repository, optional when inside a gitea repository",
	Destination: &outputValue,
}

// DefaultFlags defines flags that should be available
// for all subcommands and appended to the flags of the
// subcommand to work around issue:
// https://github.com/urfave/cli/issues/585
var DefaultFlags = []cli.Flag{
	LoginFlag,
	OutputFlag,
}

// RepoDefaultFlags defines flags that should be available
// for all subcommands working with dedicated repositories
// to work around issue:
// https://github.com/urfave/cli/issues/585
var RepoDefaultFlags = append([]cli.Flag{
	RepoFlag,
}, DefaultFlags...)

// initCommand returns repository and *Login based on flags
func initCommand() (*Login, string, string) {
	err := loadConfig(yamlConfigPath)
	if err != nil {
		log.Fatal("load config file failed", yamlConfigPath)
	}

	var login *Login
	if loginValue == "" {
		login, err = getActiveLogin()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		login = getLoginByName(loginValue)
		if login == nil {
			log.Fatal("indicated login name ", loginValue, " does not exist")
		}
	}

	repoPath := repoValue
	if repoPath == "" {
		login, repoPath, err = curGitRepoPath()
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	owner, repo := splitRepo(repoPath)
	return login, owner, repo
}