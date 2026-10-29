gomsgs
===

[![Go](https://github.com/caixw/gomsgs/actions/workflows/go.yml/badge.svg)](https://github.com/caixw/gomsgs/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/caixw/gomsgs/branch/master/graph/badge.svg)](https://codecov.io/gh/caixw/gomsgs)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/caixw/gomsgs)](https://pkg.go.dev/github.com/caixw/gomsgs)
![Go version](https://img.shields.io/github/go-mod/go-version/caixw/gomsgs)
![License](https://img.shields.io/github/license/caixw/gomsgs)

提供了以下几个命令：

- extract 从 Go 源码提取未翻译的内容；
- update 根据从 extract 提取的内容，更新其它语言的翻译内容；

安装
---

```bash
go install github.com/caixw/gomsgs
```

如果是在 `go:generate` 中使用，也可以直接使用以下命令写入 go.mod 中：

```bash
go get -tool github.com/caixw/gomsgs
```

版权
---

本项目采用 [MIT](http://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
