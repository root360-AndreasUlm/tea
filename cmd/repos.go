// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli"
)

// CmdRepos represents to login a gitea server.
var CmdRepos = cli.Command{
	Name:        "repos",
	Usage:       "Operate with repositories",
	Description: `Operate with repositories`,
	Action:      runReposList,
	Subcommands: []cli.Command{
		CmdReposList,
		CmdReposFork,
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "login, l",
			Usage: "Indicate one login, optional when inside a gitea repository",
		},
	},
}

// CmdReposList represents a sub command of issues to list issues
var CmdReposList = cli.Command{
	Name:        "ls",
	Usage:       "List available repositories",
	Description: `List available repositories`,
	Action:      runReposList,
}

func runReposList(ctx *cli.Context) error {
	login := initCommandLoginOnly(ctx)

	rps, err := login.Client().ListMyRepos()

	if err != nil {
		log.Fatal(err)
	}

	if len(rps) == 0 {
		fmt.Println("No repositories found")
		return nil
	}

	fmt.Println("Name | Type/Mode | SSH-URL | Owner")
	for _, rp := range rps {
		var mode = "source"
		if rp.Fork {
			mode = "fork"
		}
		if rp.Mirror {
			mode = "mirror"
		}
		fmt.Printf("%s | %s | %s | %s\n", rp.FullName, mode, rp.SSHURL, rp.Owner.UserName)
	}

	return nil
}

// CmdReposFork represents a sub command of issues to list issues
var CmdReposFork = cli.Command{
	Name:        "fork",
	Usage:       "fork repository",
	Description: `fork repository`,
	Action:      runReposFork,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "login, l",
			Usage: "Indicate one login, optional when inside a gitea repository",
		},
		cli.StringFlag{
			Name:  "repo, r",
			Usage: "Indicate one repo, optional when inside a gitea repository",
		},
		cli.StringFlag{
			Name:  "org",
			Usage: "Organization to fork the repository for (optional, default = logged in user)",
		},
	},
}

func runReposFork(ctx *cli.Context) error {
	login, owner, repo := initCommand(ctx)
	forkOptions := gitea.CreateForkOption{}
	if org := ctx.String("org"); org != "" {
		forkOptions = gitea.CreateForkOption{
			Organization: &org,
		}
	}

	_, err := login.Client().CreateFork(owner, repo, forkOptions)

	if err != nil {
		log.Fatal(err)
	}

	user, _ := login.Client().GetMyUserInfo()

	fmt.Printf("Forked '%s/%s' to '%s/%s'\n", owner, repo, user.UserName, repo)

	return nil
}

func initCommandLoginOnly(ctx *cli.Context) *Login {
	err := loadConfig(yamlConfigPath)
	if err != nil {
		log.Fatal("load config file failed", yamlConfigPath)
	}

	var login *Login
	if loginFlag := getGlobalFlag(ctx, "login"); loginFlag == "" {
		login, err = getActiveLogin()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		login = getLoginByName(loginFlag)
		if login == nil {
			log.Fatal("indicated login name", loginFlag, "does not exist")
		}
	}
	return login
}
