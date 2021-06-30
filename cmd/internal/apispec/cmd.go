package apispec

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/akitasoftware/akita-cli/apispec"
	"github.com/akitasoftware/akita-cli/cmd/internal/akiflag"
	"github.com/akitasoftware/akita-cli/cmd/internal/cmderr"
	"github.com/akitasoftware/akita-cli/cmd/internal/pluginloader"
	"github.com/akitasoftware/akita-cli/location"
	"github.com/akitasoftware/akita-cli/util"
	"github.com/akitasoftware/akita-libs/gitlab"
	"github.com/akitasoftware/akita-libs/tags"
)

var Cmd = &cobra.Command{
	Use:          "apispec",
	Short:        "Convert traces into an OpenAPI3 specification.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		traces, err := toLocations(tracesFlag)
		if err != nil {
			return err
		}

		traceTags, err := tags.FromPairs(tracesByTagFlag)
		if err != nil {
			return err
		}

		if len(traces) == 0 && len(traceTags) == 0 {
			return errors.New("Must specify at least one input via \"traces\" or \"trace-tag\"")
		}

		if len(traceTags) > 0 {
			var serviceName string
			if serviceFlag != "" {
				serviceName = serviceFlag
			} else if outFlag.AkitaURI != nil {
				serviceName = outFlag.AkitaURI.ServiceName
			} else {
				return errors.New("Must specify \"service\" or \"out\" to use \"trace-tag\"")
			}
			destURI, err := util.GetTraceURIByTags(akiflag.Domain, akiflag.GetClientID(), serviceName, traceTags, "trace-tag")
			if err != nil {
				return err
			}
			if destURI.ObjectName == "" {
				return cmderr.AkitaErr{Err: errors.New("No traces matching specified tag, cannot create spec")}
			}
			traces = append(traces, location.Location{AkitaURI: &destURI})
		}

		tags, err := tags.FromPairs(tagsFlag)
		if err != nil {
			return err
		}

		var gitlabMR *gitlab.MRInfo
		if gitlabProjectFlag != "" {
			gitlabMR = &gitlab.MRInfo{
				Project: gitlabProjectFlag,
				IID:     gitlabMRFlag,
				Branch:  gitlabBranchFlag,
				Commit:  gitlabCommitFlag,
			}
		}

		plugins, err := pluginloader.Load(pluginsFlag)
		if err != nil {
			return errors.Wrap(err, "failed to load plugins")
		}

		args := apispec.Args{
			ClientID:       akiflag.GetClientID(),
			Domain:         akiflag.Domain,
			Traces:         traces,
			Out:            outFlag,
			Service:        serviceFlag,
			Format:         formatFlag,
			Tags:           tags,
			PathParams:     pathParamsFlag,
			PathExclusions: pathExclusionsFlag,

			GitHubBranch: githubBranchFlag,
			GitHubCommit: githubCommitFlag,
			GitHubRepo:   githubRepoFlag,
			GitHubPR:     githubPRFlag,

			GitLabMR: gitlabMR,

			GetSpecEnableRelatedFields: getSpecEnableRelatedFieldsFlag,
			IncludeTrackers:            includeTrackersFlag,

			Plugins: plugins,
		}
		if err := apispec.Run(args); err != nil {
			return cmderr.AkitaErr{Err: err}
		}
		return nil
	},
}
