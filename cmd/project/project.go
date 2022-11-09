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
	"github.com/cnative-dev/cli/internal"
	"github.com/spf13/cobra"
)

type Project struct {
	Namespace string `json:"namespace" view:"项目ID,1"`
	Name      string `json:"name" view:"项目名,2"`
	URL       string `json:"url" view:"URL,3"`
	UpdatedAt string `json:"updatedAt" view:"更新时间,4"`
}

func NewProjectCommand() *cobra.Command {

	// projectCmd represents the project command
	var projectCmd = &cobra.Command{
		Use:     "project",
		Aliases: []string{"projects"},
		Short:   "管理项目",
	}
	internal.InitCommand(projectCmd)
	projectCmd.AddCommand(NewProjectCreateCommand())
	projectCmd.AddCommand(NewProjectDescribeCommand())
	projectCmd.AddCommand(NewProjectListCommand())
	projectCmd.AddCommand(NewProjectDeleteCommand())
	projectCmd.AddCommand(NewProjectRenameCommand())
	return projectCmd
}
