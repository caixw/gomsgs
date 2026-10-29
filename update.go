// SPDX-FileCopyrightText: 2018-2026 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json/v2"
	"errors"
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
	"github.com/issue9/messages/serialize"
)

const (
	updateTitle = localeutil.StringPhrase("update locale file")
	updateUsage = localeutil.StringPhrase(`update locale file usage

flags:
{{flags}}
`)
	srcUsage  = localeutil.StringPhrase("src locale file")
	destUsage = localeutil.StringPhrase("dest locale files")
)

func initUpdate(opt *cmdopt.CmdOpt, p *localeutil.Printer) {
	opt.New("update", updateTitle.LocaleString(p), updateUsage.LocaleString(p), func(fs *flag.FlagSet) cmdopt.DoFunc {
		src := fs.String("src", "", srcUsage.LocaleString(p))
		dest := fs.String("dest", "", destUsage.LocaleString(p))

		return func(io.Writer) error {
			log := newLogs(p, os.Stdout)

			srcMsg, err := getSrc(*src)
			if err != nil {
				return err
			}

			for path := range strings.SplitSeq(*dest, ",") {
				var dest *messages.File

				stat, err := os.Stat(path)
				if errors.Is(err, os.ErrNotExist) {
					dest = &messages.File{}
				} else if err != nil {
					return err
				} else if stat.IsDir() {
					return localeutil.Error("the dest file %s is dir", path)
				}

				filename := filepath.Base(path)
				ext := filepath.Ext(filename)
				u, err := getUnmarshalByExt(ext)
				if err != nil {
					return err
				}
				if dest == nil { // dest != nil，说明因为不存在文件，已经被初始经默认值。
					if dest, err = serialize.LoadFile(path, u); err != nil {
						return err
					}
				}

				srcMsg.MergeTo(log.WARN().LocaleString, dest, path)

				m, _, err := getMarshalByExt(ext)
				if err != nil {
					return err
				}
				if err = serialize.SaveFile(dest, path, m, os.ModePerm); err != nil {
					return err
				}
			}

			return nil
		}
	})
}

func getSrc(src string) (*messages.File, error) {
	u, err := getUnmarshalByExt(filepath.Ext(src))
	if err != nil {
		return nil, err
	}

	return serialize.LoadFile(src, u)
}

func getUnmarshalByExt(ext string) (serialize.UnmarshalFunc, error) {
	switch strings.ToLower(ext) {
	case "json", ".json":
		return func(data []byte, v any) error { return json.Unmarshal(data, v) }, nil
	case "yaml", "yml", ".yaml", ".yml":
		return yaml.Unmarshal, nil
	case "toml", ".toml":
		return toml.Unmarshal, nil
	default:
		return nil, localeutil.Error("unsupported unmarshal for %s", ext)
	}
}
