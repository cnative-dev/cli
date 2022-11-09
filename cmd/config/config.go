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
package config

import (
	"github.com/cnative-dev/cli/internal"
	"github.com/spf13/cobra"
)

type Config = map[string]string

func NewConfigCommand() *cobra.Command {
	var configCmd = &cobra.Command{
		Use:     "config",
		Aliases: []string{"configs"},
		Short:   "管理配置文件",
		Long:    "管理配置文件，注：所有的写操作都会触发服务重新部署",
	}
	internal.InitCommand(configCmd)
	configCmd.AddCommand(NewConfigListCommand())
	configCmd.AddCommand(NewConfigApplyCommand())
	configCmd.AddCommand(NewConfigDeleteCommand())
	configCmd.AddCommand(NewConfigAddCommand())
	configCmd.AddCommand(NewConfigRemoveCommand())
	return configCmd
}
