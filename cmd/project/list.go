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
	"reflect"

	"github.com/cnative-dev/cli/internal"
	"github.com/spf13/cobra"
)

func NewProjectListCommand() *cobra.Command {

	// projectListCmd represents the list command
	var projectListCmd = &cobra.Command{
		Use:   "list",
		Short: "项目列表",
		Run: func(cmd *cobra.Command, args []string) {
			internal.WithAuthorized(func() {
				projects := &[]Project{}
				if resp, err := internal.R().
					Get("/api/projects"); err == nil &&
					resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
					json.Unmarshal(resp.Body(), projects)
					internal.PrettyStructArray(projects, reflect.TypeOf(Project{}))
				} else {
					fmt.Fprintln(os.Stderr, resp.Error().(*internal.ErrResp).Details)
					return
				}
			})
		},
	}

	internal.InitCommand(projectListCmd)
	return projectListCmd
}
