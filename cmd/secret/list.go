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
package secret

import (
	"encoding/json"

	"github.com/cnative-dev/cli/internal"
	"github.com/cnative-dev/cli/internal/flags"
	"github.com/spf13/cobra"
)

func NewSecretListCommand() *cobra.Command {
	var project string
	// secretListCmd represents the list command
	var secretListCmd = &cobra.Command{
		Use:   "list",
		Short: "Secret 列表",
		Run: func(cmd *cobra.Command, args []string) {
			internal.WithAuthorized(func() {
				secret := &MaskSecret{}
				if resp, err := internal.R().
					SetPathParam("project", project).
					Get("/api/projects/{project}/secrets"); err == nil &&
					resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
					json.Unmarshal(resp.Body(), secret)
					internal.PrettyArray(*secret)
				} else {
					internal.HandleError(resp, err)
					return
				}
			})
		},
	}
	internal.InitCommand(secretListCmd)
	flags.AddProject(secretListCmd, &project)
	return secretListCmd
}
