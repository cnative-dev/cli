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
package build

import (
	"github.com/cnative-dev/cli/internal"
	"github.com/spf13/cobra"
)

type Build struct {
	ProjectId  string            `json:"project" view:"项目ID,2"`
	BuildId    string            `json:"id" view:"构建ID,1"`
	Sha        string            `json:"sha" view:"Git SHA,3"`
	Status     string            `json:"status" view:"构建状态,4"`
	Product    bool              `json:"product" view:"产物包,5"`
	Configs    map[string]string `json:"configs" view:"配置,6"`
	Secrets    []string          `json:"secrets" view:"Secrets,7"`
	CreateTime string            `json:"createdAt" view:"创建时间,8"`
}

type ListBuild struct {
	// ProjectId  string `json:"project" view:"项目ID,2"`
	BuildId    string `json:"id" view:"构建ID,1"`
	Sha        string `json:"sha" view:"Git SHA,3"`
	Status     string `json:"status" view:"构建状态,4"`
	Product    bool   `json:"product" view:"产物包,5"`
	CreateTime string `json:"createdAt" view:"创建时间,6"`
}

func NewBuildCommand() *cobra.Command {

	var cmd = &cobra.Command{
		Use:     "build",
		Aliases: []string{"builds"},
		Short:   "管理构建过程和产物",
	}
	internal.InitCommand(cmd)

	cmd.AddCommand(NewBuildListCommand())
	cmd.AddCommand(NewBuildDeleteCommand())
	cmd.AddCommand(NewBuildLogCommand())
	cmd.AddCommand(NewBuildCancelCommand())
	cmd.AddCommand(NewBuildRetryCommand())
	return cmd
}
