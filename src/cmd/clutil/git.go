package main

import (
	"os"

	"gopkg.in/src-d/go-git.v4"
)

func gitPull(p string) error {
	r, err := git.PlainOpen(p)

	if err != nil {
		return err
	}

	w, err := r.Worktree()

	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}

func gitClone(u, p string) error {
	_, err := git.PlainClone(p, false, &git.CloneOptions{
		URL:      u,
		Progress: os.Stdout,
	})

	return err
}
