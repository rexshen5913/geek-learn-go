package web

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileUploader struct {
	FileField string
	// 比如说 DST 是一个目录
	Dst string
	// DstPathFunc 用于计算目标路径
	DstPathFunc func(fh *multipart.FileHeader) string
}

func (f *FileUploader) Handle() HandleFunc {
	return func(ctx *Context) {
		file, header, err := ctx.Req.FormFile(f.FileField)
		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("上传失败")
			return
		}
		dst, err := os.OpenFile(filepath.Join(f.DstPathFunc(header), header.Filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("上传失败")
			return
		}
		io.CopyBuffer(dst, file, nil)
	}
}
