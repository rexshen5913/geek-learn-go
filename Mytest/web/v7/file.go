package web

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
)

type FileUploader struct {
	FileField   string
	DstPathFunc func(*multipart.FileHeader) string
}

func (u *FileUploader) Handle() HandleFunc {
	return func(ctx *Context) {
		// 上傳文件的邏輯
		// 1. 讀取文件內容
		// 2. 計算出目標路徑
		// 3. 保存文件
		// 4. 返回響應
		file, fileHeader, err := ctx.Req.FormFile(u.FileField)
		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("FromFile 上傳失敗" + err.Error())
			return
		}
		defer file.Close()

		// 我怎麼知道目標路徑
		// 這種做法，就是將目標路徑的計算邏輯，交給用戶
		dst := u.DstPathFunc(fileHeader)

		// 自動創建目錄路徑
		// filepath.Dir() 獲取文件路徑的目錄部分
		// os.MkdirAll() 遞歸創建所有必要的目錄，權限設置為 0755
		dir := filepath.Dir(dst)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("MkdirAll 創建目錄失敗" + err.Error())
			return
		}

		// O_WRONLY: 只寫
		// O_CREATE: 如果文件不存在，則創建
		// O_TRUNC: 如果文件存在，則清空
		// 0666: 文件權限
		dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("OpenFile 上傳失敗" + err.Error())
			return
		}
		defer dstFile.Close()

		// 保存文件
		// buffer 是為了避免每次都分配一個新的緩衝區
		// 若是沒有給，那預設會是 32 * 1024 的緩衝區 32KiB
		_, err = io.CopyBuffer(dstFile, file, nil)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("CopyFile 上傳失敗" + err.Error())
			return
		}

		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("上傳成功")
	}
}

type FileUploaderOption func(uploader *FileUploader)

func NewFileUploader(opts ...FileUploaderOption) *FileUploader {
	res := &FileUploader{
		// 可以給他預設值
		FileField: "myfile",
		DstPathFunc: func(fileHeader *multipart.FileHeader) string {
			return filepath.Join("testdata/uploads", fileHeader.Filename)
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func FileUploaderOptionWithDstPathFunc(fn func(*multipart.FileHeader) string) FileUploaderOption {
	return func(uploader *FileUploader) {
		uploader.DstPathFunc = fn
	}
}

// Downloader 下載器
type FileDownloader struct {
	// 目標路徑
	Dir string
}

func (d *FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		// 用的是 xxx?file=xxx 這樣的格式
		req, err := ctx.QueryValue("file")
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("找不到目標文件")
			return
		}
		// 需要做下校驗，避免路徑穿越
		req = filepath.Clean(req)
		dst := filepath.Join(d.Dir, req)
		dst, err = filepath.Abs(dst)
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("非法文件路徑")
			return
		}
		if !strings.Contains(dst, d.Dir) {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("非法文件路徑")
			return
		}
		path := dst
		fn := filepath.Base(path)
		header := ctx.Resp.Header()
		// Content-Disposition 指定 attachment 就是保存到本地，同時也設置了 filename
		// Content-Type 指定文件類型，這裡用的是 application/octet-stream，代表通用的二進制文件
		// Content-Transfer-Encoding 指定編碼方式，這裡用的是 binary，代表二進制編碼，相當於直接傳輸
		header.Set("Content-Disposition", "attachment; filename="+fn) // filename 是文件名
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate") // 過期後一定要回源驗證
		header.Set("Pragma", "public")                 // 舊版瀏覽器可能需要，但現代瀏覽器大多忽略

		http.ServeFile(ctx.Resp, ctx.Req, path)
	}
}

// 靜態資源下載，且須考慮緩存
type StaticResourceHandler struct {
	dir string

	mu sync.RWMutex
	// 緩存 - 改為存儲 fileCache 結構體
	cache *lru.Cache[string, *fileCache]

	//大文件不緩存
	noCacheFileSize int64
	// Content-Type 的映射
	extContentTypeMap map[string]string
}

// 可以優化的部分，將這些都放進緩存中
type fileCache struct {
	fileName    string
	fileSize    int64
	contentType string
	data        []byte
}

// 若是想要讓用戶自行設定參數，可以這樣寫
type StaticResourceHandlerOption func(handler *StaticResourceHandler)

func NewStaticResourceHandler(dir string, opts ...StaticResourceHandlerOption) *StaticResourceHandler {
	// 這個 cache 最多能放 50 個不同的 fileCache。
	// 创建一个容量为50的LRU缓存
	cache, _ := lru.New[string, *fileCache](50)
	res := &StaticResourceHandler{
		dir:             dir,
		cache:           cache,
		noCacheFileSize: 1024 * 1024 * 10, // 10MB
		extContentTypeMap: map[string]string{
			".js":   "text/javascript",
			".css":  "text/css",
			".html": "text/html",
			".png":  "image/png",
			".jpg":  "image/jpeg",
		},
	}

	for _, opt := range opts {
		opt(res)
	}

	return res

}

func StaticWithMaxFileSize(size int64) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.noCacheFileSize = size
	}
}

func StaticWithCache(cache *lru.Cache[string, *fileCache]) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.cache = cache
	}
}

func StaticWithExtContentTypeMap(extContentTypeMap map[string]string) StaticResourceHandlerOption {
	return func(handler *StaticResourceHandler) {
		handler.extContentTypeMap = extContentTypeMap
	}
}

// 若是用這種寫法，我們會使用 Build-Option 模式，比較有彈性
func (s *StaticResourceHandler) Handle(ctx *Context) {
	// 1. 拿到目標文件名
	// 2. 定為目標文件，並且讀出來
	// 3. 返回給前端
	// 先用無緩存的版本
	file, err := ctx.PathValue("file")
	if err != nil {
		ctx.RespStatusCode = http.StatusBadRequest
		ctx.RespData = []byte("請求路徑不對，找不到目標文件")
		return
	}
	dst := filepath.Join(s.dir, file)

	// 這邊需要加入緩存
	// 1. 先讀取緩存
	// 2. 如果緩存有，則直接返回
	// 3. 如果緩存沒有，則讀取文件，並且將文件放入緩存
	// 4. 返回文件
	s.mu.RLock()
	cached, ok := s.cache.Get(file)
	if ok {
		header := ctx.Resp.Header()
		header.Set("Content-Type", cached.contentType)
		header.Set("Content-Size", strconv.FormatInt(cached.fileSize, 10))
		header.Set("Content-Length", strconv.Itoa(len(cached.data)))
		header.Set("X-Cache", "HIT") // 標記來自緩存
		ctx.RespData = cached.data
		ctx.RespStatusCode = http.StatusOK
		s.mu.RUnlock()
		return
	}
	s.mu.RUnlock()

	data, err := os.ReadFile(dst)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte("服務器錯誤")
		return
	}

	// 計算文件元信息
	fileSize := int64(len(data))
	contentType := s.extContentTypeMap[filepath.Ext(file)]

	// 設置響應頭
	header := ctx.Resp.Header()
	header.Set("Content-Type", contentType)
	header.Set("Content-Size", strconv.FormatInt(fileSize, 10))
	header.Set("Content-Length", strconv.Itoa(len(data)))
	header.Set("X-Cache", "MISS") // 標記來自文件系統

	// 判斷是否需要緩存（基於文件大小）
	if fileSize < s.noCacheFileSize {
		// 小文件加入緩存
		cacheEntry := &fileCache{
			fileName:    file,
			fileSize:    fileSize,
			contentType: contentType,
			data:        data,
		}
		s.mu.Lock()
		s.cache.Add(file, cacheEntry)
		s.mu.Unlock()
	}

	ctx.RespData = data
	ctx.RespStatusCode = http.StatusOK

}
