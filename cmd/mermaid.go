/*
Copyright © 2020 Víctor Pérez @MetalBlueberry

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
package cmd

import (
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/MetalBlueberry/mdt/pkg/mdt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

// mermaidCmd represents the mermaid command
var mermaidCmd = &cobra.Command{
	Use:   "mermaid",
	Short: "mermaid related commands",
}

var mermaidWrap = &cobra.Command{
	Use:   "wrap",
	Short: "wrap mermaid fenced code blocks with an image tag",
	Run: func(cmd *cobra.Command, args []string) {

		files := []string{}

		for _, pattern := range args {
			match, err := filepath.Glob(pattern)
			if err != nil {
				log.Fatalf("Invalid pattern, %s", pattern)
			}
			files = append(files, match...)
		}

		logrus.WithField("files", files).Info("found files")

		for _, file := range files {
			source, err := ioutil.ReadFile(file)
			if err != nil {
				panic(err)
			}

			root := goldmark.DefaultParser().Parse(text.NewReader(source))
			fences, err := mdt.ParseFences(source, root)
			if err != nil {
				panic(err)
			}

			wraps, err := mdt.NewMermaidInk().WrapAll(fences)
			if err != nil {
				panic(err)
			}

			output := mdt.ApplyWraps(source, wraps)
			ioutil.WriteFile(file, output, fs.ModePerm)
			// fmt.Print(string(output))
		}
	},
}

var mermaidUpdate = &cobra.Command{
	Use:   "update",
	Short: "updated wrap mermaid fenced code blocks generated with wrap method",
	Run: func(cmd *cobra.Command, args []string) {

		files := []string{}

		for _, pattern := range args {
			match, err := filepath.Glob(pattern)
			if err != nil {
				log.Fatalf("Invalid pattern, %s", pattern)
			}
			files = append(files, match...)
		}

		logrus.WithField("files", files).Info("found files")

		for _, file := range files {
			source, err := ioutil.ReadFile(file)
			if err != nil {
				panic(err)
			}

			root := goldmark.DefaultParser().Parse(text.NewReader(source))
			wraps, err := mdt.ParseWrappedFences(source, root)
			if err != nil {
				panic(err)
			}

			err = mdt.NewMermaidInk().UpdateAll(wraps)
			if err != nil {
				panic(err)
			}

			output := mdt.ApplyWraps(source, wraps)
			ioutil.WriteFile(file, output, fs.ModePerm)
			// fmt.Print(string(output))
		}
	},
}

func init() {
	rootCmd.AddCommand(mermaidCmd)

	mermaidCmd.AddCommand(mermaidWrap)
	mermaidCmd.AddCommand(mermaidUpdate)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mermaidCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mermaidCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
