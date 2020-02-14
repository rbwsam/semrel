package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"gituhub.com/rbwsam/semrel/internal"
)

const (
	app = "semrel"
)

func main() {
	app := &cli.App{
		Name:     app,
		Usage: "Manage git releases with Semantic Versioning",
		Commands: []*cli.Command{tag()},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	/**

	  LAST_VERSION = last valid semver tag OR v0.0.0
	  NEW_VERSION = ""

	  if HEAD.commit_sha != LAST_VERSION.commit_sha {
	  	if LAST_VERSION is not a prerelease version {
	  		if current branch == master {
	  			set NEW_VERSION to LAST_VERSION with incremented minor version
	  		} else {
	  			set NEW_VERSION to LAST_VERSION with incremented patch version
	  		}
	  	} else {
	  		set NEW_VERSION to LAST_VERSION with PRE-RELEASE == ${NUM_REV_PAST_LAST_VERSION_NO_LEADING_ZEROS}.g${HEAD.commit_sha}
	  	}
	  }

	  if remote branch release-$(VERSION_MAJOR).$(VERSION_MINOR) does not exist {
	  	create it from HEAD and push it
	  }
	*/
}

func tag() *cli.Command {
	return &cli.Command{
		Name:        "tag",
		Usage: "Generates a Semantic Version for HEAD and tags it",
		Action: func(c *cli.Context) error {
			path, err := os.Getwd()
			if err != nil {
				return err
			}

			return internal.Tag(path)
		},
	}

}
