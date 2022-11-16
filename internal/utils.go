package internal

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func WithAuthorized(code func()) {
	if IsAuthorized() {
		code()
	} else {
		fmt.Println("未登录，请先使用 auth 命令登录")
	}
}

func NewSpinner() *spinner.Spinner {
	spin := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigc
		spin.Stop()
		os.Exit(0)
	}()
	return spin
}

func IsAuthorized() bool {
	if viper.IsSet("token") {
		tokenString := viper.GetString("token")
		if len(tokenString) > 0 {
			if token, _ := jwt.Parse(tokenString, nil); token != nil && token.Claims.Valid() == nil {
				return true
			} else {
				viper.Set("token", "")
				return false
			}
		}
	}
	return false
}

func InitCommand(cmd *cobra.Command) {
	cmd.Flags().BoolP("help", "h", false, "Hide HELP Flags")
	cmd.Flags().Lookup("help").Hidden = true
}

func ReadProperties(filename string) (map[string]string, error) {
	config := map[string]string{}
	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		// check if the line has = sign
		// and process the line. Ignore the rest.
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				// assign the config map
				config[key] = value
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

func UsageTemplateWithArgs(cmd *cobra.Command, argNames ...string) {
	var argNameString = ""
	if len(argNames) > 0 {
		cmd.Args = cobra.MinimumNArgs(len(argNames))
	}
	for _, s := range argNames {
		argNameString = argNameString + "<" + s + "> "
	}

	cmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}} ` + argNameString + `{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
}

func Endpoints() (map[string]string, error) {
	endpoints := make(map[string]string)
	_, err := R().SetResult(&endpoints).Get("/api/endpoints")
	return endpoints, err
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func RandSeq(length int) string {
	rand.Seed(time.Now().Unix())
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
