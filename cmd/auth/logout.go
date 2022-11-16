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
package auth

import (
	"fmt"

	"github.com/cnative-dev/cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewAuthLogoutCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "logout",
		Short: "登出 CNative",
		Run: func(cmd *cobra.Command, args []string) {
			viper.Set("token", "")
			// // 设置 netrc
			// rcs, _ := internal.ReadNetrc()
			// filtered := []internal.NetrcLine{}
			// for _, rc := range rcs {
			// 	if rc.Machine != "git.cnative.dev" {
			// 		filtered = append(filtered, rc)
			// 	}
			// }
			// internal.WriteNetrc(filtered)
			fmt.Println("已登出")
		},
	}
	internal.InitCommand(cmd)

	return cmd
}
