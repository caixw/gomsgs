// SPDX-FileCopyrightText: 2018-2026 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json/v2"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"
	"github.com/issue9/cmdopt"
	"github.com/issue9/localeutil"
	"github.com/issue9/messages"
	"github.com/issue9/messages/extract"
	"github.com/issue9/messages/serialize"
	"golang.org/x/text/language"
)

const (
	extractTitle = localeutil.StringPhrase("extract locale")
	extractUsage = localeutil.StringPhrase(`extract usage

flags：
{{flags}}
`)
	format    = localeutil.StringPhrase("file format")
	out       = localeutil.StringPhrase("out dir")
	lang      = localeutil.StringPhrase("language")
	recursive = localeutil.StringPhrase("recursive dir")
	funcs     = localeutil.StringPhrase("locale func")
	tag       = localeutil.StringPhrase("locale struct tag")
	skipMod   = localeutil.StringPhrase("skip sub module")
	info      = localeutil.StringPhrase("show info log")
)

const presetFuncs = `github.com/issue9/localeutil.Phrase,github.com/issue9/localeutil.Error,github.com/issue9/localeutil.StringPhrase`

func initExtract(opt *cmdopt.CmdOpt, p *localeutil.Printer) {
	opt.New("extract", extractTitle.LocaleString(p), extractUsage.LocaleString(p), func(fs *flag.FlagSet) cmdopt.DoFunc {
		f := fs.String("f", "yaml", format.LocaleString(p))
		o := fs.String("o", "./locales", out.LocaleString(p))
		l := fs.String("l", "und", lang.LocaleString(p))
		r := fs.Bool("r", true, recursive.LocaleString(p))
		i := fs.Bool("i", false, info.LocaleString(p))
		t := fs.String("tag", "", tag.LocaleString(p))
		funcs := fs.String("func", presetFuncs, funcs.LocaleString(p))
		skip := fs.Bool("m", true, skipMod.LocaleString(p))

		log := newLogs(p, os.Stdout)

		return func(io.Writer) error {
			u, ext, err := getMarshalByExt(*f)
			if err != nil {
				return err
			}

			// l
			lt, err := language.Parse(*l)
			if err != nil {
				return err
			}

			// out
			if err := os.MkdirAll(*o, os.ModePerm); err != nil {
				return err
			}
			out := filepath.Join(*o, *l+ext)

			if fs.NArg() == 0 {
				return localeutil.Error("no src dir")
			}

			file := &messages.File{
				Languages: []language.Tag{lt},
			}
			ctx := context.Background()
			for _, dir := range fs.Args() {
				opt := &extract.Options{
					Language:      lt,
					Root:          dir,
					Recursive:     *r,
					SkipSubModule: *skip,
					WarnLog:       log.WARN().LocaleString,
					InfoLog:       func(localeutil.Stringer) {}, // 默认不输出提示信息
					Funcs:         strings.Split(*funcs, ","),
					Tag:           *t,
				}
				if *i {
					opt.InfoLog = log.INFO().LocaleString
				}

				lang, err := extract.Extract(ctx, opt)
				if err != nil {
					return err
				}
				file.Join(lang)
			}

			return serialize.SaveFile(file, out, u, os.ModePerm)
		}
	})
}

func getMarshalByExt(ext string) (serialize.MarshalFunc, string, error) {
	switch strings.ToLower(ext) {
	case "json", ".json":
		return func(v any) ([]byte, error) {
			return json.Marshal(v)
		}, ".json", nil
	case "yaml", "yml", ".yaml", ".yml":
		return func(v any) ([]byte, error) {
			return yaml.MarshalWithOptions(v, yaml.UseLiteralStyleIfMultiline(true))
		}, ".yaml", nil
	case "toml", ".toml":
		return toml.Marshal, ".yaml", nil
	default:
		return nil, "", localeutil.Error("unsupported marshal for %s", ext)
	}
}
