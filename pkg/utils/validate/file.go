package validate

import (
	"AIGO/config"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"
)

var (
	// 允许的文件扩展名
	allowedExtensions = []string{".md", ".txt", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"}
)

// CheckValidateFile 验证文件是否合法
func CheckValidateFile(file *multipart.FileHeader) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}
	// 获取文件的扩展名并将其转成小写
	extension := strings.ToLower(filepath.Ext(file.Filename))
	if !slices.Contains(allowedExtensions, extension) {
		return fmt.Errorf("file extension %s is not allowed", extension)
	}
	// 判断文件大小是否超过限制缓解服务端压力
	if file.Size > config.Cfg.AppCfg.MaxFileSize*(2<<20) { // 10MB
		return fmt.Errorf("file size %d exceeds the limit of %dMB", file.Size, config.Cfg.AppCfg.MaxFileSize)
	}
	return nil
}
