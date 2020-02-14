package internal

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
)

func Tag(path string) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = LastVersion(r)
	return err
}

func LastVersion(r *git.Repository) (semver.Version, error) {
	headRef, err := r.Head()
	checkErr(err)

	logOpts := &git.LogOptions{
		From:  headRef.Hash(),
		Order: git.LogOrderCommitterTime,
	}

	iter, err := r.Log(logOpts)
	checkErr(err)

	const maxHistory = 100

	// Loop through commits until we find a semver tag OR we hit our limit

	for i := 0; i < maxHistory; i++ {
		c, err := iter.Next()
		checkErr(err)

		// if c is a semver tag: store the semver, break out of loop

		fmt.Println(c.Hash)
	}

	return semver.Version{}, err
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
