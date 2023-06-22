package utils

// exePath 當前執行檔的絕對路徑
var exePath string

// exeDir 當前執行檔的所在目錄
var exeDir string

// EXEPath 取得當前執行檔的絕對路徑
func EXEPath() string {
	return exePath
}

// EXEDir 取得當前執行檔所在目錄的絕對路徑
func EXEDir() string {
	return exeDir
}
