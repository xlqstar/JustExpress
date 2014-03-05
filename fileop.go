//文件操作函数集

package justExpress

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//拷贝文件
func CopyFile(srcPath, dstPath string) {
	src, err := os.Open(srcPath)
	defer src.Close()
	if err != nil {
		log.Println("拷贝文件：读取 " + srcPath + " 文件出错！")
		return
	}

	dst, err := os.Create(dstPath)
	defer src.Close()
	if err != nil {
		log.Println("拷贝文件：创建 " + dstPath + " 文件出错！")
		return
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		log.Println("拷贝文件：从 " + srcPath + " 拷贝至" + dstPath + "文件出错！")
		return
	}
	return
}

//拷贝目录
func CopyDir(srcDirPath string, destDirPath string) {
	srcDirPath = filepath.Clean(srcDirPath)
	destDirPath = filepath.Clean(destDirPath)

	filepath.Walk(srcDirPath,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				log.Println("拷贝目录:" + path + " 文件不存在!")
				return nil
			}
			if f.IsDir() {
				dest_dir := destDirPath + strings.TrimPrefix(path, srcDirPath)
				if _, err := os.Stat(dest_dir); os.IsNotExist(err) {
					err := os.Mkdir(dest_dir, os.ModePerm)
					if err != nil {
						log.Fatal("拷贝目录:" + dest_dir + " 目录创建失败!")
						return nil
					}
				}
			} else {
				dest_file := destDirPath + strings.TrimPrefix(path, srcDirPath)
				CopyFile(path, dest_file)
			}
			return nil
		})
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func Remkdir(dirPath string) {
	os.RemoveAll(dirPath)
	err := os.Mkdir(dirPath, os.ModePerm)
	if err != nil {
		log.Fatal("重建" + dirPath + "目录发生未预料到的错误")
	}
}
