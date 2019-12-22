//          FILE:  gen_prompt.go
//
//         USAGE:  gen_prompt.go
//
//   DESCRIPTION:  $description
//
//       OPTIONS:  ---
//  REQUIREMENTS:  ---
//          BUGS:  ---
//         NOTES:  ---
//        AUTHOR:  John van Zantvoort (jvzantvoort), john@vanzantvoort.org
//       COMPANY:  JDC
//       CREATED:  01-Apr-2019
//
// Copyright (C) 2019 John van Zantvoort
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.
//	"io/ioutil"
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"text/template"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

const (
	defaultMainColor  string = "light_cyan"
	defaultOSColor    string = "green"
	defaultDirColor   string = "yellow"
	color_end         string = "0m"
	color_black       string = "0;30m"
	color_red         string = "0;31m"
	color_green       string = "0;32m"
	color_brown       string = "0;33m"
	color_blue        string = "0;34m"
	color_purple      string = "0;35m"
	color_cyan        string = "0;36m"
	color_light_gray  string = "0;37m"
	color_dark_gray   string = "1;30m"
	color_gray        string = "1:30m"
	color_light_blue  string = "1;34m"
	color_light_cyan  string = "1;36m"
	color_light_green string = "1;32m"
	color_light_purpl string = "1;35m"
	color_light_red   string = "1;31m"
	color_white       string = "1;37m"
	color_yellow      string = "1;33m"
	TemplateText      string = `
PS1="{{.MainColor}}\u@\h{{.EndColor}}/{{.OSColor}}{{.OSName}}{{.EndColor}} \T [{{.DirColor}}\w{{.EndColor}}]
# "
`
)

type ConfigOptions struct {
	MainColor string `toml:"main_color"`
	OSColor   string `toml:"os_color"`
	DirColor  string `toml:"dir_color"`
}

type tomlConfig struct {
	Hosts map[string]ConfigOptions
}

type AppPath struct {
	ConfigFile string
	EnvFile    string
}

type LSBInfo struct {
	Filename string
	Name     string
	Class    string
}

var lsbInfoSets = []LSBInfo{
	LSBInfo{
		Filename: "/etc/centos-release",
		Name:     "CentOS",
		Class:    "redhat",
	},
	LSBInfo{
		Filename: "/etc/fedora-release",
		Name:     "Fedora",
		Class:    "redhat",
	},
	LSBInfo{
		Filename: "/etc/redhat-release",
		Name:     "RedHat",
		Class:    "redhat",
	},
	LSBInfo{
		Filename: "/etc/SuSE-release",
		Name:     "SuSE",
		Class:    "suse",
	},
	LSBInfo{
		Filename: "/etc/mandrake-release",
		Name:     "Mandrake",
		Class:    "mandrake",
	},
	LSBInfo{
		Filename: "/etc/debian_version",
		Name:     "Debian",
		Class:    "debian",
	},
	LSBInfo{
		Filename: "/etc/wrs-release",
		Name:     "WindRiver",
		Class:    "windriver",
	},
	LSBInfo{
		Filename: "/etc/snow-release",
		Name:     "Snow",
		Class:    "snow",
	},
}

type OSInfo struct {
	Name  string
	Class string
}

type TemplateFields struct {
	MainColor string
	OSColor   string
	DirColor  string
	EndColor  string
	OSName    string
	OSClass   string
}

// printc wraps the color definition in escape strings needed.
func printc(color string) string {
	return fmt.Sprintf("\\[\033[%s\\]", color)
}

func GetPath(name string) string {
	var retv string
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	homedir := usr.HomeDir
	hostname := shortHostname()
	switch name {
	case "env":
		retv = path.Join(homedir, ".bash", "prompt.d", hostname+".sh")
	case "config":
		retv = path.Join(homedir, ".userconfig.cfg")
	default:
		retv = homedir
	}
	return retv
}

// buildConfig contruct the text from the template definition and arguments.
func (t TemplateFields) buildConfig() string {
	tmpl, err := template.New("prompt").Parse(TemplateText)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, t)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func shortHostname() string {
	fqdn, _ := os.Hostname()
	parts := strings.Split(fqdn, ".")
	return strings.ToLower(parts[0])
}

func uname() string {
	output, err := exec.Command("uname", "-s").CombinedOutput()
	if err == nil {
		return strings.TrimSuffix(string(output), "\n")
	}
	return "unknown"
}

func colornameToColorvalue(name string) string {
	retv := printc(color_white)
	if name == "black" {
		retv = printc(color_black)
	} else if name == "blue" {
		retv = printc(color_blue)
	} else if name == "brown" {
		retv = printc(color_brown)
	} else if name == "cyan" {
		retv = printc(color_cyan)
	} else if name == "dark_gray" {
		retv = printc(color_dark_gray)
	} else if name == "gray" {
		retv = printc(color_gray)
	} else if name == "green" {
		retv = printc(color_green)
	} else if name == "light_blue" {
		retv = printc(color_light_blue)
	} else if name == "light_cyan" {
		retv = printc(color_light_cyan)
	} else if name == "light_gray" {
		retv = printc(color_light_gray)
	} else if name == "light_green" {
		retv = printc(color_light_green)
	} else if name == "light_purpl" {
		retv = printc(color_light_purpl)
	} else if name == "light_red" {
		retv = printc(color_light_red)
	} else if name == "purple" {
		retv = printc(color_purple)
	} else if name == "red" {
		retv = printc(color_red)
	} else if name == "white" {
		retv = printc(color_white)
	} else if name == "yellow" {
		retv = printc(color_yellow)
	}
	return retv
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	if len(os.Args[1:]) == 0 {
		fmt.Printf("source %s\n", GetPath("env"))
		return
	}

	viper.SetConfigName("prompt")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}
	viper.Set("Verbose", true)
	lala := viper.AllSettings()
	b, err := yaml.Marshal(lala)
	check(err)
	fmt.Print(string(b))
	//
	// command argument handling
	//
	// defaultMainColor  string = "light_cyan"
	// defaultOSColor    string = "green"
	// defaultDirColor   string = "yellow"
	flags := flag.NewFlagSet("prompt", flag.ExitOnError)
	fl_MainColor := flags.String("main", defaultMainColor, "Main color")
	fl_OSColor := flags.String("os", defaultOSColor, "Os color")
	fl_DirColor := flags.String("dir", defaultDirColor, "Dir color")
	// fl_write := flags.Bool("write", false, "Write the prompt file")
	flags.Parse(os.Args[1:])

	main_c := colornameToColorvalue(*fl_MainColor)
	os_c := colornameToColorvalue(*fl_OSColor)
	dir_c := colornameToColorvalue(*fl_DirColor)
	//
	// command argument handling, end
	//

	//
	// Obtain name and class
	//
	var os_name string
	var os_class string
	kernel_name := uname()
	if kernel_name == "Linux" {
		for _, i := range lsbInfoSets {
			if _, err := os.Stat(i.Filename); os.IsNotExist(err) {
				continue
			}
			if len(os_name) != 0 {
				continue
			}
			os_name = i.Name
			os_class = i.Class
		}

	} else if kernel_name == "Darwin" {
		os_name = "Darwin"
		os_class = "mac"

	} else if kernel_name == "SunOS" {
		os_name = "SunOS"
		os_class = "solaris"
	}
	//
	// Obtain name and class, end
	//

	// setup the template defaults
	var promptColor TemplateFields
	promptColor.MainColor = main_c
	promptColor.OSColor = os_c
	promptColor.DirColor = dir_c
	promptColor.EndColor = printc(color_end)
	promptColor.OSName = os_name
	promptColor.OSClass = os_class

	fmt.Print(promptColor.buildConfig())
	fmt.Print(GetPath("config"))

}

// vim: noexpandtab filetype=go
