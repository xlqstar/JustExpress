package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/xlqstar/Just"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	USAGE = `
Version:   Just 1.0
Author:    Xiao Liqun(xlqstar@gmail.com)

Usage:
	just newsite <site_name>
	just [-site <site_name>] post [<log_type>] <title>
	just [-site <site_name>] delete <title>
	just [-site <site_name>] build
	just [-site <site_name>] rebuild
	just [-site <site_name>] rebuildall
	just [-site <site_name>] switchtheme [<theme_name>]
	just [-site <site_name>] resize
	just [-site <site_name>] preview
	just [-site <site_name>] sitesroot [<sites_root_path>]

PS:
	log_type: album | article
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
	var sitePath string
	if len(args) > 0 {
		if args[0] == "newsite" {
			if len(args) != 2 {
				argumentErr()
			}
			just.NewSite(args[1])
			return
		} else if args[0] == "sitesroot" {
			if len(args) == 1 {
				fmt.Println()
				fmt.Println("当前站点根路径：" + just.SitesRoot(""))
				fmt.Println()
				fmt.Println("您可以通过 just [-site <site_name>] switchsitesroot <sites_root_path> 命令可以更改该路径，也可以直接修改data文件。")
			} else if len(args) == 2 {
				just.SitesRoot(args[1])
			} else {
				argumentErr()
			}
			return
		}
		sitePath = just.GetSitePath(*siteName)
	} else {
		sitePath = just.GetSitePath(*siteName)
		fmt.Print("请输入命令：")
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		args = getArgs(string(data))
	}

	switch args[0] {
	default:
		commandErr()
		/*
			just [-site <site_name>] post <title> <log_type> <categorys> <tags>

			case "post":
				// sitePath := just.GetSitePath(*siteName)
				var title, logType, categorys, tags string
				if len(args) < 4 || len(args) > 5 {
					argumentErr()
				} else {
					title, logType, categorys = args[1], args[2], args[3]
					if len(args) == 5 {
						tags = args[4]
					}
				}
				just.Post(sitePath, title, logType, categorys, tags)
		*/
	case "post": //quick post
		// sitePath := just.GetSitePath(*siteName)
		if len(args) == 2 {
			just.QuickPost(sitePath, "", args[1])
		} else if len(args) == 3 {
			just.QuickPost(sitePath, args[1], args[2])
		} else {
			argumentErr()
		}
	case "delete":
		// sitePath := just.GetSitePath(*siteName)
		if len(args) == 2 {
			just.Delete(sitePath, args[1])
		} else if len(args) == 1 {

			fileList, _ := filepath.Glob(sitePath + "\\*")
			var _fileList []string
			for k := range fileList {
				if filepath.Base(fileList[k]) == "complied" {
					if k-1 < 0 {
						_fileList = fileList[k+1:]
					} else if k+1 > 0 {
						_fileList = fileList[0 : k-1]
					} else {
						_fileList = append(fileList[0:k-1], fileList[k+1:]...)
					}
					break
				}
			}
			if len(_fileList) == 0 {
				fmt.Println("站点目录中没有日志！您可以通过`just post [<log_type>] <title>`命令新建日志！ ")
				os.Exit(1)
			}
			fmt.Println("\n日志列表：\n")
			for k := range _fileList {
				_fileList[k] = filepath.Base(_fileList[k])
				fmt.Println(strconv.Itoa(k) + ". " + _fileList[k])
			}
			var num int
			for true {
				fmt.Println()
				fmt.Print("请输入序号：")
				fmt.Scanf("%d\n", &num)
				if num < len(_fileList) && num >= 0 {
					break
				} else {
					fmt.Println()
					fmt.Println("不存在该日志，请重新输入正确的序号！")
					fmt.Println()
				}
			}

			fmt.Println()
			if fileList[num] != "complied" {
				just.Delete(sitePath, _fileList[num])
			} else {
				fmt.Println("\n\n请按照提示输入！")
			}
		} else {
			argumentErr()
		}
	case "switchtheme":
		// sitePath := just.GetSitePath(*siteName)
		if len(args) == 2 {
			just.SwitchTheme(sitePath, args[1])
		} else if len(args) == 1 {
			fileList, _ := filepath.Glob(".\\themes\\*")
			fmt.Println("\n有如下主题可供选择：")
			for k := range fileList {
				fmt.Println(strconv.Itoa(k) + ". " + filepath.Base(fileList[k]))
			}
			fmt.Print("\n请输入主题名称：")
			var NO int
			for true {
				fmt.Print("请输入序号：")
				fmt.Scanf("%d\n", &NO)
				if NO < len(fileList) && NO >= 0 {
					break
				} else {
					fmt.Println()
					fmt.Println("不存在该主题，请重新输入正确的序号！")
					fmt.Println()
				}
			}

			just.SwitchTheme(sitePath, fileList[NO])
		} else {
			argumentErr()
		}
	case "build":
		// sitePath := just.GetSitePath(*siteName)
		just.Build(sitePath, false)
	case "rebuild": //重新构建(只构建html部分)
		if len(args) != 1 {
			argumentErr()
		}
		// sitePath := just.GetSitePath(*siteName)
		just.Rebuild(sitePath)
	case "rebuildall": //重新构建(彻底重新构建，包括图片及其所有附件)
		if len(args) != 1 {
			argumentErr()
		}
		// sitePath := just.GetSitePath(*siteName)
		just.RebuildAll(sitePath)
	case "resize": //调整所有图片大小
		if len(args) != 1 {
			argumentErr()
		}
		// sitePath := just.GetSitePath(*siteName)
		just.ImgResize(sitePath)
	case "preview":
		// sitePath := just.GetSitePath(*siteName)
		cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", sitePath+"\\complied\\index.html")
		cmd.Run()
	}
}

func getArgs(cmdsStr string) []string {
	cmdsArry := strings.Fields(cmdsStr)
	var args []string
	var tmp string
	var flag bool
	for _, cmd := range cmdsArry {
		if strings.HasPrefix(cmd, "\"") {
			flag = true
			tmp = strings.TrimPrefix(cmd, "\"")
			continue
		}
		if strings.HasSuffix(cmd, "\"") {
			flag = false
			tmp += " " + strings.TrimSuffix(cmd, "\"")
			args = append(args, tmp)
			tmp = ""
			continue
		}

		if flag == true {
			tmp += " " + cmd
		} else {
			args = append(args, cmd)
		}
	}
	return args
}
func argumentErr() {
	fmt.Println()
	fmt.Println("参数数量异常，您可以通过`just -h`命令获取帮助")
	os.Exit(1)
}

func commandErr() {
	fmt.Println()
	fmt.Println("请输入正确命令，您可以通过`just -h`命令获取帮助")
	os.Exit(1)
}
