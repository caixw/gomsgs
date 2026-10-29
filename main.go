// SPDX-FileCopyrightText: 2018-2026 caixw
//
// SPDX-License-Identifier: MIT

//go:generate go build ./
//go:generate ./gomsgs extract -l=und -f=yaml ./
//go:generate ./gomsgs update -src=./locales/und.yaml -dest=./locales/zh.yaml

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"github.com/issue9/logs/v7"
	"golang.org/x/text/message"
)

const (
	helpTitle = localeutil.StringPhrase("show help")
	helpUsage = localeutil.StringPhrase("show current help info")
	usageTpl  = localeutil.StringPhrase(`Extract or update localized file

commands：
{{commands}}

flags：
{{flags}}
`)
)

func main() {
	tag, err := localeutil.DetectUserLanguageTag()
	if err != nil {
		panic(err)
	}
	p := message.NewPrinter(tag)

	var opt *cmdopt.CmdOpt
	opt = cmdopt.New(&cmdopt.Options{
		Output:        os.Stdout,
		ErrorHandling: flag.ContinueOnError,
		UsageTemplate: usageTpl.LocaleString(p),
		NotFound:      buildNotFound(p),
		Command: func(fs *flag.FlagSet) cmdopt.DoFunc {
			v := fs.Bool("v", false, localeutil.StringPhrase("show version").LocaleString(p))

			return func(w io.Writer) error {
				if *v {
					_, err := fmt.Fprint(w, localeutil.Phrase("gomsg: %s\nbuild with: %s\n", getAppVersion(), runtime.Version()).LocaleString(p))
					return err
				}

				// 没有任何选项指定，输出帮助信息。
				_, err := io.WriteString(w, opt.Usage())
				return err
			}
		},
	})

	initExtract(opt, p)
	initUpdate(opt, p)
	cmdopt.Help(opt, "help", helpTitle.LocaleString(p), helpUsage.LocaleString(p))

	if err := opt.Exec(os.Args[1:]); err != nil {
		panic(err)
	}
}

func buildNotFound(p *localeutil.Printer) cmdopt.NotFound {
	return func(_ *cmdopt.CmdOpt, s string, _ []string) string {
		return localeutil.Phrase("command %s not found", s).LocaleString(p)
	}
}

func getAppVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	}
	return ""
}

func newLogs(p *localeutil.Printer, out io.Writer) *logs.Logs {
	return logs.New(
		logs.NewTermHandler(out, nil),
		logs.WithLevels(logs.AllLevels()...),
		logs.WithCreated(logs.NanoLayout),
		logs.WithLocale(p),
	)
}
