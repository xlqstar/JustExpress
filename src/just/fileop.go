package just

import (
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path"
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
	filepath.Walk(srcDirPath,
		func(path string, f os.FileInfo, err error) error{
			if f == nil {
				log.Println("拷贝目录:"+path+" 文件不存在!")
				return nil
			}
			if f.IsDir() {
				dest_dir := strings.Replace(path, srcDirPath, destDirPath, -1)
				if _, err := os.Stat(dest_dir); os.IsNotExist(err) {
					err := os.Mkdir(dest_dir, os.ModePerm)
					if err != nil {
						log.Fatal("拷贝目录:"+dest_dir+" 目录创建失败!")
						return nil
					}
				}
			} else {
				dest_file := strings.Replace(path, srcDirPath, destDirPath, -1)
				CopyFile(path, dest_file)
			}
			return nil
		})
}

func Resize(src string, dest string, width uint) {
	extName := strings.ToLower(path.Ext(src))

	if extName == ".gif" {
		CopyFile(src, dest)
	} else {
		file, err := os.Open(src)
		if err != nil {
			log.Fatal("resize图像：读取 " + src + " 图片出错！")
		}

		// decode jpeg into image.Image
		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatal("resize图像：解析 " + src + " 图片出错！")
		}
		file.Close()

		// resize to width 1000 using Lanczos resampling
		// and preserve aspect ratio
		// fmt.Println("==resize图像开始:"+fmt.Sprintf("%d", time.Now().Unix()))
		m := resize.Resize(width, 0, img, resize.NearestNeighbor)
		// fmt.Println("==resize图像结束:"+fmt.Sprintf("%d", time.Now().Unix()))
		out, err := os.Create(dest)
		if err != nil {
			log.Fatal("resize图像：创建 " + dest + " 图片出错！")
		}
		defer out.Close()

		// write new image to file
		if extName == ".png" {
			png.Encode(out, m)
		} else {
			jpeg.Encode(out, m, nil)
		}
	}
}
