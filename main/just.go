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
	USAGE = `
Version:   Just 1.0
Author:    Xiao Liqun(xlqstar@gmail.com)

Usage:
	just newsite <site_name>
	just [-site <site_name>] complie
	just [-site <site_name>] qpost <title>
	just [-site <site_name>] post <title> <log_type> <categorys> <tags>
	just [-site <site_name>] delete <title>
	just [-site <site_name>] switchtheme <theme_name>
	just [-site <site_name>] rebuild
	just [-site <site_name>] rebuildall
	just [-site <site_name>] resize
	just [-site <site_name>] preview
	just [-site <site_name>] switchsitesroot <sites_root_path>
`
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, USAGE)
		// flag.PrintDefaults()
		// fmt.Fprintln(os.Stderr, "  -h   : show help usage")
	}
	// log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	siteName := flag.String("site", "", "站点标识")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		log.Fatal("请输入命令")
	}

	switch args[0] {
	default:
		log.Fatal("请输入正确命令")
	case "complie":
		sitePath := just.GetSitePath(*siteName)
		just.Complie(sitePath, false)
	case "post":
		sitePath := just.GetSitePath(*siteName)
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
		sitePath := just.GetSitePath(*siteName)
		just.Delete(sitePath, args[1])
	case "switchtheme":
		if len(args) != 2 {
			log.Fatal("参数数量异常")
		}
		sitePath := just.GetSitePath(*siteName)
		just.SwitchTheme(sitePath, args[1])
	case "rebuild": //重新构建(只构建html部分)
		if len(args) != 1 {
			log.Fatal("参数数量异常")
		}
		sitePath := just.GetSitePath(*siteName)
		just.Rebuild(sitePath)
	case "rebuildall": //重新构建(彻底重新构建，包括图片及其所有附件)
		if len(args) != 1 {
			log.Fatal("参数数量异常")
		}
		sitePath := just.GetSitePath(*siteName)
		just.RebuildAll(sitePath)
	case "resize": //调整所有图片大小
		if len(args) != 1 {
			log.Fatal("参数数量异常")
		}
		sitePath := just.GetSitePath(*siteName)
		just.ImgResize(sitePath)
	case "newsite":
		if len(args) != 2 {
			log.Fatal("参数数量异常")
		}
		just.NewSite(args[1])
	case "switchsitesroot":
		if len(args) != 2 {
			log.Fatal("参数数量异常")
		}
		just.SitesRoot(args[0])
	case "qpost": //quick post
		sitePath := just.GetSitePath(*siteName)
		if len(args) == 2 {
			just.QuickPost(sitePath, "", args[1])
		} else if len(args) == 3 {
			just.QuickPost(sitePath, args[1], args[2])
		} else {
			log.Fatal("参数数量异常")
		}
	case "preview":
		sitePath := just.GetSitePath(*siteName)
		cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", sitePath+"\\complied\\index.html")
		cmd.Run()
	}
}
