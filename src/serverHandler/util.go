package serverHandler

import (
	"../shimgo"
	"os"
	"path"
	"strings"
)

func needResponseBody(method string) bool {
	return method != shimgo.Net_Http_MethodHead &&
		method != shimgo.Net_Http_MethodOptions &&
		method != shimgo.Net_Http_MethodConnect &&
		method != shimgo.Net_Http_MethodTrace
}

func containsItem(infos []os.FileInfo, name string) bool {
	for i := range infos {
		if infos[i].Name() == name {
			return true
		}
	}
	return false
}

func getCleanFilePath(requestPath string) (filePath string, ok bool) {
	filePath = path.Clean(requestPath)
	ok = filePath == path.Base(filePath)

	return
}

func getCleanDirFilePath(requestPath string) (filePath string, ok bool) {
	filePath = path.Clean(strings.Replace(requestPath, "\\", "/", -1))
	ok = filePath[0] != '/' && filePath != "." && filePath != ".." && !strings.HasPrefix(filePath, "../")

	return
}
