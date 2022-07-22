package param

import (
	"github.com/P4elme6ka/go-http-media-server/src/util"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func splitMapping(input string) (k, v string, ok bool) {
	sep, sepLen := utf8.DecodeRuneInString(input)
	if sepLen == 0 {
		return
	}
	entry := input[sepLen:]
	if len(entry) == 0 {
		return
	}

	sepIndex := strings.IndexRune(entry, sep)
	if sepIndex <= 0 || sepIndex+sepLen == len(entry) {
		return
	}

	k = entry[:sepIndex]
	v = entry[sepIndex+sepLen:]
	return k, v, true
}

func normalizePathMaps(inputs []string) map[string]string {
	maps := make(map[string]string, len(inputs))

	for _, input := range inputs {
		urlPath, fsPath, ok := splitMapping(input)
		if !ok {
			continue
		}

		urlPath = util.CleanUrlPath(urlPath)
		fsPath = filepath.Clean(fsPath)
		maps[urlPath] = fsPath
	}

	return maps
}

func normalizePathMapsNoCase(inputs []string) map[string]string {
	maps := make(map[string]string, len(inputs))

	for _, input := range inputs {
		urlPath, fsPath, ok := splitMapping(input)
		if !ok {
			continue
		}

		urlPath = util.CleanUrlPath(urlPath)
		fsPath = filepath.Clean(fsPath)

		for url := range maps {
			if strings.EqualFold(url, urlPath) {
				delete(maps, url)
			}
		}

		maps[urlPath] = fsPath
	}

	return maps
}

func normalizeUrlPaths(inputs []string) []string {
	outputs := make([]string, 0, len(inputs))

	for _, input := range inputs {
		if len(input) == 0 {
			continue
		}
		outputs = append(outputs, util.CleanUrlPath(input))
	}

	return outputs
}

func normalizeFsPaths(inputs []string) []string {
	outputs := make([]string, 0, len(inputs))

	for _, input := range inputs {
		if len(input) == 0 {
			continue
		}

		abs, err := util.NormalizeFsPath(input)
		if err != nil {
			continue
		}

		outputs = append(outputs, abs)
	}

	return outputs
}

func normalizeFilenames(inputs []string) []string {
	outputs := make([]string, 0, len(inputs))

	for _, input := range inputs {
		if len(input) == 0 {
			continue
		}

		if strings.IndexByte(input, '/') >= 0 {
			continue
		}

		if filepath.Separator != '/' && strings.IndexByte(input, filepath.Separator) >= 0 {
			continue
		}

		outputs = append(outputs, input)
	}

	return outputs
}

func validateHstsPort(listensPlain, ListensTLS []string) bool {
	var fromOK, toOK bool

	for _, listen := range listensPlain {
		port := util.ExtractListenPort(listen)
		if len(port) == 0 || port == "80" {
			fromOK = true
			break
		}
	}

	for _, listen := range ListensTLS {
		port := util.ExtractListenPort(listen)
		if len(port) == 0 || port == "443" {
			toOK = true
			break
		}
	}

	return fromOK && toOK
}

func normalizeHttpsPort(httpsPort string, ListensTLS []string) (string, bool) {
	if len(httpsPort) > 0 {
		httpsPort = util.ExtractListenPort(httpsPort)
		if len(httpsPort) == 0 {
			return "", false
		}
	} else if len(ListensTLS) > 0 {
		httpsPort = util.ExtractListenPort(ListensTLS[0])
	}

	lenHttpsPort := len(httpsPort)
	httpsColonPort := ":" + httpsPort
	for _, listen := range ListensTLS {
		if lenHttpsPort == 0 && len(listen) == 0 {
			return "", true
		}

		port := util.ExtractListenPort(listen)
		if lenHttpsPort == 0 && len(port) == 0 {
			return "", true
		}
		if httpsPort == port {
			return httpsColonPort, true
		}

		if httpsPort == "443" && port == "" {
			return "", true
		}
	}

	return "", false
}
