package golang

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	ga "github.com/merico-dev/stream/internal/pkg/plugin/githubactions"
	"github.com/merico-dev/stream/pkg/util/github"
	"github.com/merico-dev/stream/pkg/util/log"
)

func Read(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options

	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	ghOptions := &github.Option{
		Owner:    opts.Owner,
		Repo:     opts.Repo,
		NeedAuth: true,
	}
	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return nil, err
	}

	path, err := ghClient.GetWorkflowPath()
	if err != nil {
		return nil, err
	}
	if path == "" {
		// file not found
		return nil, nil
	}

	log.Debugf("Language is: %s.", ga.GetLanguage(opts.Language))

	return ga.BuildReadState(path), nil
}
