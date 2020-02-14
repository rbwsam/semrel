package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "boom",
		Usage: "make an explosive entrance",
		Action: func(c *cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
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