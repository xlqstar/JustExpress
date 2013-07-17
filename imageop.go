//图像操作函数集

package just

import (
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"
	"strings"
)

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
