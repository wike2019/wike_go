package Archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"github.com/wike2019/wike_go/src/util/help"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)
type  Archive struct {

}

func New() *Archive {
	return &Archive{}
}
func (this * Archive) TarFile(src,dist string) error {
	fw, err := os.Create(dist)
	defer fw.Close()
	if err != nil{
		return  err
	}

	tw := tar.NewWriter(fw)
	tw.Close()

	fi, err := os.Stat(src)
	if err != nil{
		return  err
	}
	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil{
		return  err
	}
	return tw.WriteHeader(hdr)

	fr, err := os.Open(src)
	if err != nil{
		return  err
	}
	defer fr.Close()
	written, err := io.Copy(tw, fr)
	if err != nil{
		return  err
	}
	log.Printf("共写入了 %d 个字符的数据\n",written)
}
func (this * Archive) UnTarFile(src string) error{
	fr, err := os.Open(src)
	if err != nil{
		return  err
	}
	defer fr.Close()
	tr := tar.NewReader(fr)
	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next(){
		// 处理 err ！= nil 的情况
		if err != nil{
			return  err
		}
		// 获取文件信息
		fi := hdr.FileInfo()

		// 创建一个空文件，用来写入解包后的数据
		fw, err := os.Create(fi.Name())
		if err != nil{
			return  err
		}

		// 将 tr 写入到 fw
		n, err := io.Copy(fw, tr)
		if err != nil{
			return  err
		}
		log.Printf("解包： %s 到 %s ，共处理了 %d 个字符的数据。", src,fi.Name(),n)

		// 设置文件权限，这样可以保证和原始文件权限相同，如果不设置，会根据当前系统的 umask 来设置。
		os.Chmod(fi.Name(),fi.Mode().Perm())

		// 注意，因为是在循环中，所以就没有使用 defer 关闭文件
		// 如果想使用 defer 的话，可以将文件写入的步骤单独封装在一个函数中即可
		fw.Close()
	}

}

func (this * Archive) TarDir(src,dist string) error{
	fw, err := os.Create(dist)
	if err != nil{
		return  err
	}
	defer fw.Close()

	// 将 tar 包使用 gzip 压缩，其实添加压缩功能很简单，
	// 只需要在 fw 和 tw 之前加上一层压缩就行了，和 Linux 的管道的感觉类似
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// 创建 Tar.Writer 结构
	tw := tar.NewWriter(gw)
	// 如果需要启用 gzip 将上面代码注释，换成下面的

	defer tw.Close()

	// 下面就该开始处理数据了，这里的思路就是递归处理目录及目录下的所有文件和目录
	// 这里可以自己写个递归来处理，不过 Golang 提供了 filepath.Walk 函数，可以很方便的做这个事情
	// 直接将这个函数的处理结果返回就行，需要传给它一个源文件或目录，它就可以自己去处理
	// 我们就只需要去实现我们自己的 打包逻辑即可，不需要再去路径相关的事情
	err= filepath.Walk(src, func(fileName string, fi os.FileInfo, err error) error {
		// 因为这个闭包会返回个 error ，所以先要处理一下这个
		if err != nil{
			return  err
		}

		// 这里就不需要我们自己再 os.Stat 了，它已经做好了，我们直接使用 fi 即可
		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil{
			return  err
		}
		// 这里需要处理下 hdr 中的 Name，因为默认文件的名字是不带路径的，
		// 打包之后所有文件就会堆在一起，这样就破坏了原本的目录结果
		// 例如： 将原本 hdr.Name 的 syslog 替换程 log/syslog
		// 这个其实也很简单，回调函数的 fileName 字段给我们返回来的就是完整路径的 log/syslog
		// strings.TrimPrefix 将 fileName 的最左侧的 / 去掉，
		// 熟悉 Linux 的都知道为什么要去掉这个
		hdr.Name = strings.TrimPrefix(fileName, string(filepath.Separator))

		// 写入文件信息
		return tw.WriteHeader(hdr)

		// 判断下文件是否是标准文件，如果不是就不处理了，
		// 如： 目录，这里就只记录了文件信息，不会执行下面的 copy
		if !fi.Mode().IsRegular() {
			return nil
		}

		// 打开文件
		fr, err := os.Open(fileName)
		defer fr.Close()
		if err != nil{
			return  err
		}

		// copy 文件数据到 tw
		n, err := io.Copy(tw, fr)
		if err != nil{
			return  err
		}

		// 记录下过程，这个可以不记录，这个看需要，这样可以看到打包的过程
		log.Printf("成功打包 %s ，共写入了 %d 字节的数据\n", fileName, n)
		return  nil
	})
	if err != nil{
		return  err
	}
}
func (this * Archive) UnTarDir(src,dist string)error{
	// 打开准备解压的 tar 包
	fr, err := os.Open(src)
	if err != nil{
		return  err
	}
	defer fr.Close()

	// 将打开的文件先解压
	gr, err := gzip.NewReader(fr)
	if err != nil{
		return  err
	}
	defer gr.Close()

	// 通过 gr 创建 tar.Reader
	tr := tar.NewReader(gr)

	// 现在已经获得了 tar.Reader 结构了，只需要循环里面的数据写入文件就可以了
	for {
		hdr, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			if err != nil{
				return  err
			}
		case hdr == nil:
			continue
		}

		// 处理下保存路径，将要保存的目录加上 header 中的 Name
		// 这个变量保存的有可能是目录，有可能是文件，所以就叫 FileDir 了……
		dstFileDir := filepath.Join(dist, hdr.Name)

		// 根据 header 的 Typeflag 字段，判断文件的类型
		switch hdr.Typeflag {
		case tar.TypeDir: // 如果是目录时候，创建目录
			// 判断下目录是否存在，不存在就创建
			if b := help.ExistDir(dstFileDir); !b {
				// 使用 MkdirAll 不使用 Mkdir ，就类似 Linux 终端下的 mkdir -p，
				// 可以递归创建每一级目录
				return(os.MkdirAll(dstFileDir, 0775))
			}
		case tar.TypeReg: // 如果是文件就写入到磁盘
			// 创建一个可以读写的文件，权限就使用 header 中记录的权限
			// 因为操作系统的 FileMode 是 int32 类型的，hdr 中的是 int64，所以转换下
			file, err := os.OpenFile(dstFileDir, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil{
				return  err
			}
			n, err := io.Copy(file, tr)
			if err != nil{
				return  err
			}
			// 将解压结果输出显示
			fmt.Printf("成功解压： %s , 共处理了 %d 个字符\n", dstFileDir, n)

			// 不要忘记关闭打开的文件，因为它是在 for 循环中，不能使用 defer
			// 如果想使用 defer 就放在一个单独的函数中
			file.Close()
		}
	}


}

func (this * Archive)  Zip(src, dist string)error {
	// 创建准备写入的文件
	fw, err := os.Create(dist)
	defer fw.Close()
	if err != nil{
		return  err
	}

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer zw.Close()


	// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
	filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) error {

		if errBack != nil{
			return  err
		}
		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil{
			return  err
		}

		// 替换文件信息中的文件名
		fh.Name = strings.TrimPrefix(path, string(filepath.Separator))

		// 这步开始没有加，会发现解压的时候说它不是个目录
		if fi.IsDir() {
			fh.Name += "/"
		}

		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil{
			return  err
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(path)
		defer fr.Close()
		if err != nil{
			return  err
		}

		// 将打开的文件 Copy 到 w
		n, err := io.Copy(w, fr)
		if err != nil{
			return  err
		}
		// 输出压缩的内容
		fmt.Printf("成功压缩文件： %s, 共写入了 %d 个字符的数据\n", path, n)
		return  nil
	})
}

func (this * Archive)  UnZip(src, dist string)error {
	// 打开压缩文件，这个 zip 包有个方便的 ReadCloser 类型
	// 这个里面有个方便的 OpenReader 函数，可以比 tar 的时候省去一个打开文件的步骤
	zr, err := zip.OpenReader(src)
	defer zr.Close()
	if err != nil{
		return  err
	}

	// 如果解压后不是放在当前目录就按照保存目录去创建目录
	if dist != "" {
		return(os.MkdirAll(dist, 0755))
	}

	// 遍历 zr ，将文件写入到磁盘
	for _, file := range zr.File {
		path := filepath.Join(dist, file.Name)

		// 如果是目录，就创建目录
		if file.FileInfo().IsDir() {
			return(os.MkdirAll(path, file.Mode()))
			// 因为是目录，跳过当前循环，因为后面都是文件的处理
			continue
		}

		// 获取到 Reader
		fr, err := file.Open()
		if err != nil{
			return  err
		}

		// 创建要写出的文件对应的 Write
		fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil{
			return  err
		}
		n, err := io.Copy(fw, fr)
		if err != nil{
			return  err
		}

		// 将解压的结果输出
		fmt.Printf("成功解压 %s ，共写入了 %d 个字符的数据\n", path, n)

		// 因为是在循环中，无法使用 defer ，直接放在最后
		// 不过这样也有问题，当出现 err 的时候就不会执行这个了，
		// 可以把它单独放在一个函数中，这里是个实验，就这样了
		fw.Close()
		fr.Close()
	}
}
