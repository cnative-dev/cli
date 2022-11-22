/*
Copyright © 2022 Jingsi <lookisliu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package project

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cnative-dev/cli/internal"
	"github.com/spf13/cobra"
)

func NewProjectDeleteCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "删除项目及相关资源和服务",
		Run: func(cmd *cobra.Command, args []string) {
			internal.WithAuthorized(func() {
				if resp, err := internal.R().
					SetPathParam("projectId", args[0]).
					Delete("/api/projects/{projectId}"); err == nil &&
					resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
					fmt.Printf("项目 %s 已删除\n", args[0])
				} else {
					fmt.Fprintln(os.Stderr, resp.Error().(*internal.ErrResp).Details)
					os.Exit(1)
				}
			})
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if internal.IsAuthorized() {
				projects := []Project{}
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
		},
	}
	internal.InitCommand(cmd)
	internal.UsageTemplateWithArgs(cmd, "项目名")

	return cmd
}
