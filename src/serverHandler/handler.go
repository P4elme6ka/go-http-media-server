package serverHandler

import (
	"github.com/P4elme6ka/go-http-media-server/src/param"
	"github.com/P4elme6ka/go-http-media-server/src/serverErrHandler"
	"github.com/P4elme6ka/go-http-media-server/src/serverLog"
	"github.com/P4elme6ka/go-http-media-server/src/tpl"
	"github.com/P4elme6ka/go-http-media-server/src/user"
	"net/http"
	"regexp"
	"strings"
)

var defaultHandler = http.NotFoundHandler()

var createFileServer func(root string) http.Handler

type handler struct {
	root        string
	emptyRoot   bool
	globalHsts  bool
	globalHttps bool
	httpsPort   string // with prefix ":"
	defaultSort string
	aliasPrefix string

	dirIndexes []string
	aliases    aliases

	globalHeaders [][2]string

	globalUpload bool
	uploadUrls   []string
	uploadDirs   []string

	globalMkdir bool
	mkdirUrls   []string
	mkdirDirs   []string

	globalDelete bool
	deleteUrls   []string
	deleteDirs   []string

	globalArchive bool
	archiveUrls   []string
	archiveDirs   []string

	globalCors bool
	corsUrls   []string
	corsDirs   []string

	globalAuth bool
	authUrls   []string
	authDirs   []string
	users      user.List

	shows     *regexp.Regexp
	showDirs  *regexp.Regexp
	showFiles *regexp.Regexp
	hides     *regexp.Regexp
	hideDirs  *regexp.Regexp
	hideFiles *regexp.Regexp
	theme     tpl.Theme

	fileServer http.Handler

	logger     *serverLog.Logger
	errHandler *serverErrHandler.ErrHandler
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logRequest(r)

	// hsts redirect
	if h.globalHsts && h.hsts(w, r) {
		return
	}

	// https redirect
	if h.globalHttps && h.https(w, r) {
		return
	}

	// asset
	const assetPrefix = "asset="
	if strings.HasPrefix(r.URL.RawQuery, assetPrefix) {
		assetPath := r.URL.RawQuery[len(assetPrefix):]
		h.asset(w, r, assetPath)
		return
	}

	// data
	data := h.getResponseData(r)
	if len(data.errors) > 0 {
		h.logErrors(data.errors...)
	}
	file := data.File
	if file != nil {
		defer file.Close()
	}

	if data.NeedAuth {
		h.needAuth(w, r)
	}
	if !data.AuthSuccess {
		h.authFailed(w)
		return
	}

	h.header(w)

	if data.CanCors {
		h.cors(w)
	}

	if data.IsMutate {
		h.mutate(w, r, data)
		return
	}

	// regular flows
	if len(r.URL.RawQuery) >= 3 {
		switch r.URL.RawQuery[:3] {
		case "tar":
			h.tar(w, r, data)
			return
		case "tgz":
			h.tgz(w, r, data)
			return
		case "zip":
			h.zip(w, r, data)
			return
		}
	}

	item := data.Item
	if data.WantJson {
		h.json(w, r, data)
	} else if file != nil && item != nil && !item.IsDir() {
		h.content(w, r, data)
	} else {
		h.page(w, r, data)
	}
}

func newHandler(
	p *param.Param,
	root string,
	aliasPrefix string,
	allAliases aliases,
	users user.List,
	theme tpl.Theme,
	logger *serverLog.Logger,
	errHandler *serverErrHandler.ErrHandler,
) http.Handler {
	emptyRoot := p.EmptyRoot && aliasPrefix == "/"

	aliases := aliases{}
	for _, alias := range allAliases {
		if alias.isSuccessorOf(aliasPrefix) {
			aliases = append(aliases, alias)
		}
	}

	var fileServer http.Handler
	if !emptyRoot && createFileServer != nil { // for WSL 1 fix
		fileServer = createFileServer(root)
	}

	h := &handler{
		root:        root,
		emptyRoot:   emptyRoot,
		globalHsts:  p.GlobalHsts,
		globalHttps: p.GlobalHttps,
		httpsPort:   p.HttpsPort,
		defaultSort: p.DefaultSort,
		aliasPrefix: aliasPrefix,
		aliases:     aliases,

		dirIndexes: p.DirIndexes,

		globalHeaders: p.GlobalHeaders,

		globalUpload: p.GlobalUpload,
		uploadUrls:   p.UploadUrls,
		uploadDirs:   p.UploadDirs,

		globalMkdir: p.GlobalMkdir,
		mkdirUrls:   p.MkdirUrls,
		mkdirDirs:   p.MkdirDirs,

		globalDelete: p.GlobalDelete,
		deleteUrls:   p.DeleteUrls,
		deleteDirs:   p.DeleteDirs,

		globalArchive: p.GlobalArchive,
		archiveUrls:   p.ArchiveUrls,
		archiveDirs:   p.ArchiveDirs,

		globalCors: p.GlobalCors,
		corsUrls:   p.CorsUrls,
		corsDirs:   p.CorsDirs,

		globalAuth: p.GlobalAuth,
		authUrls:   p.AuthUrls,
		authDirs:   p.AuthDirs,
		users:      users,

		shows:     p.Shows,
		showDirs:  p.ShowDirs,
		showFiles: p.ShowFiles,
		hides:     p.Hides,
		hideDirs:  p.HideDirs,
		hideFiles: p.HideFiles,
		theme:     theme,

		fileServer: fileServer,

		logger:     logger,
		errHandler: errHandler,
	}
	return h
}
