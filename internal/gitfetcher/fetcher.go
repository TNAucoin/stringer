package gitfetcher

import (
	"fmt"
	"io"
	"log"
	"net/http"

	gp "github.com/tnaucoin/stringer/parser"
	"github.com/tnaucoin/stringer/types"
)

const githubRawURL = "https://raw.githubusercontent.com"

type Options struct {
	Repo string
	Ref  string
}

func FetchCompositeActionsFromRepo(opts Options) ([]types.CompositeAction, error) {
	if opts.Repo == "" {
		return nil, fmt.Errorf("repo is required")
	}
	if opts.Ref == "" {
		opts.Ref = "main"
	}

	// TODO: placeholder for now replace with something real
	paths := []string{
		"",
	}
	var actions []types.CompositeAction

	for _, path := range paths {
		data, err := fetchFileFromGithub(opts.Repo, opts.Ref, path)
		if err != nil {
			log.Printf("warning: fetch failed for %s: %v", path, err)
			continue
		}
		action, err := gp.ParseCompositeActionFromBytes(data, path)
		if err != nil {
			log.Printf("warning: failed to parse %s: %v", path, err)
			continue
		}
		actions = append(actions, action)

	}
	return actions, nil
}

func fetchFileFromGithub(repo, ref, path string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s/%s", githubRawURL, repo, ref, path)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from github: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github returned %d for %s", resp.StatusCode, url)
	}

	return io.ReadAll(resp.Body)
}
