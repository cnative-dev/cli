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
package main

import (
	"fmt"
	"os"

	"github.com/cnative-dev/cli/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	// defer profile.Start().Stop()
	cmd.Execute()
	if err := viper.ReadInConfig(); err == nil {
		// 真正 cmd.Execute 了之后的分支，通过 cobra.OnInitialize 完成了 viper 的配置
		viper.WriteConfig()
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cnative" (without extension).
		viper.AddConfigPath(fmt.Sprintf("%s/.cnative/", home))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

		// 所有没有 Execute 的分支，因为不会调用 cobra.OnInitialize
		viper.SafeWriteConfig()
	}
}
