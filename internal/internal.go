package internal

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func Tag(path string) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return errors.WithStack(err)
	}

	v, err := LastVersion(r)
	checkErr(err)

	fmt.Println("Found version:", v.String())

	return err
}

func LastVersion(r *git.Repository) (semver.Version, error) {
	tagTable := map[plumbing.Hash]*object.Tag{}

	tagIter, err := r.TagObjects()
	checkErr(err)

	// Collect all the annotated tags
	_ = tagIter.ForEach(func(tag *object.Tag) error {
		tagTable[tag.Target] = tag
		return nil
	})

	// Get HEAD commit ref
	headRef, err := r.Head()
	checkErr(err)

	logOpts := &git.LogOptions{
		From:  headRef.Hash(),
		Order: git.LogOrderCommitterTime,
	}

	// Get commit log from HEAD
	logIter, err := r.Log(logOpts)
	checkErr(err)
	defer logIter.Close()

	// Iterate through commit history until we find a annotated tag in semver format
	for i := 0; i < 100; i++ {
		commit, err := logIter.Next()
		checkErr(err)

		if tag, ok := tagTable[commit.Hash]; ok {
			v, err := semver.ParseTolerant(tag.Name)
			if err != nil {
				continue
			}

			fmt.Println("Found SemVer tagged commit:", commit.Hash, tag.Name, v.String(), tag.Hash)
			return v, nil
		}
	}

	return semver.Make("0.0.0")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
