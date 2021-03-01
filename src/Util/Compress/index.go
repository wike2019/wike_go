package Compress

import (
	"bytes"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"io"
	"os"
)

type Compress struct {

}

func New() *Compress{
	return &Compress{}
}
func (this *Compress)  Zlib(data [] byte) [] byte{
	buf := bytes.NewBuffer(nil)
	w := zlib.NewWriter(buf)
	// 写入待压缩内容
	w.Write(data)
	w.Close()
	return buf.Bytes()
}

func (this *Compress)  UnZlib(data [] byte) ([] byte,error){
	buf := bytes.NewBuffer(data)
	buf_res := bytes.NewBuffer(nil)
	r ,err:= zlib.NewReader(buf)
	if err != nil{
		return  nil,err
	}
	defer r.Close()
	io.Copy(buf_res,r)
	return buf_res.Bytes(),nil
}
func (this *Compress)  Lzw(data [] byte) [] byte{
	buf := bytes.NewBuffer(nil)
	w := lzw.NewWriter(buf, lzw.LSB, 8)
	// 写入待压缩内容
	w.Write(data)
	w.Close()
	return buf.Bytes()
}

func (this *Compress)  UnLzw(data [] byte) [] byte{
	buf := bytes.NewBuffer(data)
	buf_res := bytes.NewBuffer(nil)
	r := lzw.NewReader(buf, lzw.LSB, 8)
	defer r.Close()
	io.Copy(buf_res,r)
	return buf_res.Bytes()
}
func (this *Compress) Gzip (src,dist string ) error{
	newFile, err := os.Create(dist)
	if err != nil{
		return  err
	}
	defer newFile.Close()

	file, err := os.Open(src)
	if err != nil{
		return  err
	}
	zw := gzip.NewWriter(newFile)

	filestat, err := file.Stat()
	if err != nil{
		return  err
	}
	zw.Name = filestat.Name()
	zw.ModTime = filestat.ModTime()
	_, err = io.Copy(zw, file)
	if err != nil{
		return  err
	}

	zw.Flush()
	return zw.Close()
}
func (this *Compress) UnGzip (src,dist string ) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	newFile, err := os.Create(dist)
	if err != nil {
		return err
	}
	defer newFile.Close()

	zr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	zr.Name = fileStat.Name()
	zr.ModTime = fileStat.ModTime()
	_, err = io.Copy(newFile, zr)
	if err != nil {
		return err
	}

	return zr.Close()
}