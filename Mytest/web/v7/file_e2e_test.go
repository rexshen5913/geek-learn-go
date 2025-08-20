package web

import (
	"path/filepath"
	"testing"
)

// func TestFileUpload(t *testing.T) {
// 	tpl, err := template.ParseGlob("testdata/tpls/*.html")
// 	require.NoError(t, err)
// 	engine := &GoTemplateEngine{
// 		T: tpl,
// 	}

// 	h := NewHttpServer(ServerWithTemplateEngine(engine))
// 	h.Get("/upload", func(c *Context) {
// 		err := c.Render("upload.html", nil)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	})

// 	fu := &FileUploader{
// 		// 這裡的 myfile 是 html 表單中的 name 屬性
// 		// <input type="file" name="myfile">
// 		FileField: "myfile",
// 		// 這裡的 DstPathFunc 是目標路徑的計算邏輯
// 		DstPathFunc: func(fileHeader *multipart.FileHeader) string {
// 			return filepath.Join("testdata/uploads", fileHeader.Filename)
// 		},
// 	}

// 	h.Post("/upload", fu.Handle())

// 	h.Start(":8081")
// }

// func TestFileUploaderOption(t *testing.T) {
// 	tpl, err := template.ParseGlob("testdata/tpls/*.html")
// 	require.NoError(t, err)
// 	engine := &GoTemplateEngine{
// 		T: tpl,
// 	}

// 	h := NewHttpServer(ServerWithTemplateEngine(engine))

// 	h.Get("/upload", func(c *Context) {
// 		err := c.Render("upload.html", nil)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	})

// 	fu := NewFileUploader(
// 		FileUploaderOptionWithDstPathFunc(func(fileHeader *multipart.FileHeader) string {
// 			return filepath.Join("testdata/uploads", fileHeader.Filename)
// 		}),
// 	)

// 	h.Post("/upload", fu.Handle())

// 	h.Start(":8081")
// }

// func TestFileDownloader(t *testing.T) {
// 	h := NewHttpServer()

// 	fu := &FileDownloader{
// 		Dir: filepath.Join("testdata", "downloads"),
// 	}

// 	h.Get("/download", fu.Handle())

// 	h.Start(":8081")
// }

func TestStaticResourceHandler_Handle(t *testing.T) {
	h := NewHttpServer()
	fu := NewStaticResourceHandler(filepath.Join("testdata", "static"))
	// localhost:8081/static/test.txt
	h.Get("/static/:file", fu.Handle)

	h.Start(":8081")
}
