package internal

import (
	"io"
	"strconv"

	"github.com/blang/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
)

const (
	maxCommitHistory = 100
	initialVersion   = "0.0.0"
)

func CurrentVersion(path string) (semver.Version, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return semver.Version{}, errors.WithStack(err)
	}

	lastVerCommit, lastVer, err := LastVersion(r)
	if err != nil {
		return semver.Version{}, errors.WithStack(err)
	}

	if lastVerCommit == nil {
		return lastVer, nil
	}

	headRef, err := r.Head()
	if err != nil {
		return semver.Version{}, errors.WithStack(err)
	}

	if headRef.Hash() == lastVerCommit.Hash {
		return lastVer, nil
	}

	count, err := revisionCount(r, headRef.Hash(), lastVerCommit.Hash)
	if err != nil {
		return semver.Version{}, errors.WithStack(err)
	}

	// Generate the PreRelease version
	prCount, err := semver.NewPRVersion(strconv.Itoa(count))
	if err != nil {
		return semver.Version{}, errors.WithStack(err)
	}
	prShortHash, err := semver.NewPRVersion(headRef.Hash().String()[:8])
	if err != nil {
		return semver.Version{}, errors.WithStack(err)
	}

	ver := semver.Version{
		Major: lastVer.Major,
		Minor: lastVer.Minor,
		Patch: lastVer.Patch,
		Pre:   []semver.PRVersion{prCount, prShortHash},
	}

	// If the last version was not a preRelease we need to increment the patch version
	if len(lastVer.Pre) == 0 {
		ver.Patch = lastVer.Patch + 1
	}

	return ver, nil
}

func LastVersion(r *git.Repository) (*object.Commit, semver.Version, error) {
	tagTable := map[plumbing.Hash][]*object.Tag{}

	tagIter, err := r.TagObjects()
	if err != nil {
		return nil, semver.Version{}, errors.WithStack(err)
	}

	// Index all annotated tags by their target commit hash
	_ = tagIter.ForEach(func(tag *object.Tag) error {
		tagTable[tag.Target] = append(tagTable[tag.Target], tag)
		return nil
	})

	// Get HEAD commit ref
	headRef, err := r.Head()
	if err != nil {
		return nil, semver.Version{}, errors.WithStack(err)
	}

	logOpts := &git.LogOptions{
		From:  headRef.Hash(),
		Order: git.LogOrderCommitterTime,
	}

	// Get commit log from HEAD
	logIter, err := r.Log(logOpts)
	if err != nil {
		return nil, semver.Version{}, errors.WithStack(err)
	}
	defer logIter.Close()

	// Iterate through commit history until we find commit with an annotated tag in semver format
	// We limit our search here
	for i := 0; i < maxCommitHistory; i++ {
		commit, err := logIter.Next()
		if err != nil {
			if err == io.EOF {
				// We ran out of commits
				break
			}
			return nil, semver.Version{}, errors.WithStack(err)
		}

		// Find tags that point to this commit
		if tags, ok := tagTable[commit.Hash]; ok {
			// Parse each tag as a Semantic Version
			// TODO: Handle case when there are multiple valid SemVer tags pointing to the same commit,
			// 	currently we just pick the first one.
			for _, tag := range tags {
				v, err := semver.ParseTolerant(tag.Name)
				if err == nil {
					return commit, v, nil
				}
			}
		}
	}

	v, err := semver.Make(initialVersion)
	if err != nil {
		return nil, semver.Version{}, errors.WithStack(err)
	}
	return nil, v, nil
}

func revisionCount(repo *git.Repository, newer, older plumbing.Hash) (int, error) {
	logOpts := &git.LogOptions{
		From:  newer,
		Order: git.LogOrderCommitterTime,
	}

	// Get commit log from newer
	logIter, err := repo.Log(logOpts)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer logIter.Close()

	// We limit our search here
	for i := 0; i < maxCommitHistory; i++ {
		commit, err := logIter.Next()
		if err != nil {
			return 0, errors.WithStack(err)
		}
		if commit.Hash == older {
			return i, nil
		}
	}
	return 0, errors.Errorf("unable to find commit %s within the last %d commits", older, maxCommitHistory)
}
