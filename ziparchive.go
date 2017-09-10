package httphandlers

import (
	"archive/zip"
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type contentRecord struct {
	zipFile     *zip.File
	modifyTime  time.Time
	eTag        string
	contentType string
	l           sync.Mutex
}

// Build ETag string from given time.
func makeETagFromTime(t time.Time) (eTag string) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(t.Unix()))
	eTag = "\"" + hex.EncodeToString(buf) + "\""
	return eTag
}

// Get content type from file extension.
// MIME type application/octet-stream will be applied if failed to guess from file extension.
func contentTypeFromFileName(fileName string) (mimeType string) {
	if ctype := mime.TypeByExtension(filepath.Ext(fileName)); "" != ctype {
		return ctype
	}
	return "application/octet-stream"
}

func newContentRecordFromZipFile(zipFile *zip.File) (r *contentRecord) {
	modifyTime := zipFile.ModTime().Truncate(time.Second)
	eTag := makeETagFromTime(modifyTime)
	contentType := contentTypeFromFileName(zipFile.Name)
	return &contentRecord{
		zipFile:     zipFile,
		modifyTime:  modifyTime,
		eTag:        eTag,
		contentType: contentType,
	}
}

func (r *contentRecord) openZipFileReader() (reader io.ReadCloser, contentSize int64, err error) {
	r.l.Lock()
	defer r.l.Unlock()
	reader, err = r.zipFile.Open()
	return reader, int64(r.zipFile.UncompressedSize64), err
}

// ZipArchiveContentServer is a http.Handler to serve content from Zip archive file.
//
// Limitation: Range request is not supported.
//
// Here is an example of usage:
//
//	func main() {
//		var err error
//		handler, err := utilhttphandlers.NewZipArchiveContentServer("/path/to/file.zip", "web/", "index.html")
//		if nil != err {
//			log.Fatalf("failed on setting up zip content serving handler: %v", err)
//			return
//		}
//		defer handler.Close()
//		err = http.ListenAndServe(":8080", handler)
//		log.Fatalf("result of http.ListenAndServe(): %v", err)
//	}
//
// In above example, if http://localhost:8080/ is requested the content will
// be served from `web/index.html` of zip archive.
type ZipArchiveContentServer struct {
	fp                 *zip.ReadCloser
	contentMap         map[string]*contentRecord
	contentMapLock     sync.RWMutex
	defaultContentPath string
}

// NewZipArchiveContentServer creates a new instance of ZipArchiveContentServer.
//
// The fileName is the path of zip archive file to be serve. The pathPrefix is
// the path of content folder inside zip archive. Empty string can be given if
// content should be serve from root of zip archive. The defaultContentPath is
// the path to default content with pathPrefix stripped (eg. index.html).
func NewZipArchiveContentServer(fileName, pathPrefix, defaultContentPath string) (h *ZipArchiveContentServer, err error) {
	fp, err := zip.OpenReader(fileName)
	if nil != err {
		return nil, err
	}
	contentMap := make(map[string]*contentRecord)
	pathPrefix = strings.Trim(pathPrefix, "/")
	doPrefixCheck := false
	if "" != pathPrefix {
		pathPrefix = pathPrefix + "/"
		doPrefixCheck = true
	}
	defaultContentPath = strings.Trim(defaultContentPath, "/")
	for _, f := range fp.File {
		mode := f.Mode()
		if 0 != (mode & os.ModeType) {
			// not regular file
			continue
		}
		name := f.Name
		// Name: It must be a relative path: it must not start with a drive
		// letter (e.g. C:) or leading slash, and only forward slashes
		// are allowed. (https://godoc.org/archive/zip#FileHeader)
		if doPrefixCheck {
			if p := strings.TrimPrefix(name, pathPrefix); len(p) < len(name) {
				name = p
			} else {
				log.Printf("NewZipArchiveContentServer: skip-content: %v (p=%v; pathPrefix=%v)", name, p, pathPrefix)
				continue
			}
		}
		contentMap[name] = newContentRecordFromZipFile(f)
	}
	h = &ZipArchiveContentServer{
		fp:                 fp,
		contentMap:         contentMap,
		defaultContentPath: defaultContentPath,
	}
	return h, nil
}

func (h *ZipArchiveContentServer) lookupContent(name string) (c *contentRecord) {
	h.contentMapLock.RLock()
	defer h.contentMapLock.RUnlock()
	return h.contentMap[name]
}

// ServeHTTP fulfill the request with content in zip archive.
func (h *ZipArchiveContentServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")
	if "" == path {
		path = h.defaultContentPath
	}
	c := h.lookupContent(path)
	if nil == c {
		http.NotFound(w, r)
		return
	}
	if ConditionalGet(w, r, c.eTag, c.modifyTime, time.Second*2) {
		return
	}
	zipReader, contentSize, err := c.openZipFileReader()
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer zipReader.Close()
	attachModificationTagHeader(w, c.eTag, c.modifyTime)
	w.Header().Set("Content-Type", c.contentType)
	w.Header().Set("Content-Length", strconv.FormatInt(contentSize, 10))
	w.WriteHeader(http.StatusOK)
	if written, err := io.CopyN(w, zipReader, contentSize); nil != err {
		log.Printf("ZipArchiveContentServer.ServeHTTP: failed on sending zip content to remote: %v (written=%v)", err, written)
	}
}

// Close the zip archive file.
func (h *ZipArchiveContentServer) Close() (err error) {
	h.contentMapLock.Lock()
	defer h.contentMapLock.Unlock()
	h.contentMap = make(map[string]*contentRecord)
	return h.fp.Close()
}
