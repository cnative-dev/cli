package flags

import (
	"encoding/json"

	"github.com/cnative-dev/cli/internal"
	"github.com/spf13/cobra"
)

type ProjectName struct {
	Name string `json:"name"`
}

func AddProject(cmd *cobra.Command, p *string) {
	cmd.Flags().StringVarP(p, "project", "p", "", "项目名")
	cmd.MarkFlagRequired("project")
	cmd.RegisterFlagCompletionFunc("project", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		if internal.IsAuthorized() {
			projects := []ProjectName{}
			if resp, err := internal.R().
				Get("/api/projects"); err == nil && resp.StatusCode() == 200 {
				json.Unmarshal(resp.Body(), &projects)
				names := []string{}
				for i := 0; i < len(projects); i++ {
					names = append(names, projects[i].Name)
				}
				return names, cobra.ShellCompDirectiveDefault
			}
		}
		return []string{}, cobra.ShellCompDirectiveError
	})
}
