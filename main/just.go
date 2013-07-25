package main

import (
	"flag"
	"fmt"
	"github.com/xlqstar/Just"
	"log"
	"os"
	"os/exec"
)

const (
	VER    = "1.0"
	AUTHOR = "肖立群(xlqstar@gmail.com)"
	USAGE  = "Usage:\n\tjust [-site sitename] [command]"
)

var (
	siteName = flag.String("site", "", "站点标识")
	args     []string
)

func init() {
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", "JustExpress")
		flag.PrintDefaults()
		// fmt.Fprintln(os.Stderr, "  -h   : show help usage")
	}
	// log.SetFlags(log.Lshortfile | log.LstdFlags)
	fmt.Println("Version\t:" + VER)
	fmt.Println("Author\t:" + AUTHOR)
	args = flag.Args()
}

func main() {
	if len(args) == 0 {
		log.Fatal("请输入命令")
	}
	sitePath := just.GetSitePath(*siteName)
	switch args[0] {
	default:
		log.Fatal("请输入正确命令")
	case "complie":
		just.Complie(sitePath, false)
	case "post":
		var title, logType, categorys, tags string
		if len(args) < 4 || len(args) > 5 {
			log.Fatal("参数数量异常")
		} else {
			title, logType, categorys = args[1], args[2], args[3]
			if len(args) == 5 {
				tags = args[4]
			}
		}
		just.Post(sitePath, title, logType, categorys, tags)
	case "delete":
		if len(args) != 2 {
			log.Fatal("参数数量异常")
		}
		just.Delete(sitePath, args[1])
	case "switchtheme":
		if len(args) != 2 {
			log.Fatal("参数数量异常")
		}
		just.SwitchTheme(sitePath, args[1])
	case "rebuild": //重新构建(只构建html部分)
		if len(args) != 1 {
			log.Fatal("参数数量异常")
		}
		just.Rebuild(sitePath)
	case "rebuildall": //重新构建(彻底重新构建，包括图片及其所有附件)
		if len(args) != 1 {
			log.Fatal("参数数量异常")
		}
		just.RebuildAll(sitePath)
	case "resize": //调整所有图片大小
		if len(args) != 1 {
			log.Fatal("参数数量异常")
		}
		just.ImgResize(sitePath)
	case "newsite":
		if len(args) != 2 {
			log.Fatal("参数数量异常")
		}
		just.NewSite(*siteName)
	case "switchsitesroot":
		if len(args) != 2 {
			log.Fatal("参数数量异常")
		}
		just.SitesRoot(args[1])
	case "qpost": //quick post
		if len(args) == 2 {
			just.QuickPost(sitePath, "", args[1])
		} else if len(args) == 3 {
			just.QuickPost(sitePath, args[1], args[2])
		} else {
			log.Fatal("参数数量异常")
		}
	case "preview":
		cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", sitePath+"\\complied\\index.html")
		cmd.Run()
	}
}
