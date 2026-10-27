package httphandlers

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strconv"
)

const cspNonceSize = 16
const cspNonceBase64Size = 22 // base64.RawURLEncoding.EncodedLen(cspNonceSize)

type ContentReplacementRule struct {
	TargetBytes      []byte
	ReplacementBytes []byte
}

func NewContentReplacementRule(targetText, replacementText string) ContentReplacementRule {
	return ContentReplacementRule{
		TargetBytes:      []byte(targetText),
		ReplacementBytes: []byte(replacementText),
	}
}

func splitCSPNonceContentBytesIntoParts(data []byte, nonceReplacementKeyBytes []byte, contentReplacements []ContentReplacementRule) (staticContentParts [][]byte) {
	staticContentParts = [][]byte{data}
	for _, rule := range contentReplacements {
		var nextStaticContentParts [][]byte
		replacementParts := bytes.Split(rule.ReplacementBytes, nonceReplacementKeyBytes)
		replacementPartCount := len(replacementParts)
		for _, part := range staticContentParts {
			splitParts := bytes.Split(part, rule.TargetBytes)
			splitPartCount := len(splitParts)
			if splitPartCount == 1 {
				nextStaticContentParts = append(nextStaticContentParts, part)
				continue
			}
			if replacementPartCount == 1 {
				nextStaticContentParts = append(nextStaticContentParts, bytes.Join(splitParts, rule.ReplacementBytes))
				continue
			}
			lastSubPartIndex := splitPartCount - 1
			for idx, subPart := range splitParts {
				if idx == lastSubPartIndex {
					aux := append(replacementParts[replacementPartCount-1], subPart...)
					nextStaticContentParts = append(nextStaticContentParts, aux)
					continue
				}
				if idx == 0 {
					aux := append(subPart, replacementParts[0]...)
					nextStaticContentParts = append(nextStaticContentParts, aux)
				} else {
					aux := append(append(replacementParts[replacementPartCount-1], subPart...), replacementParts[0]...)
					nextStaticContentParts = append(nextStaticContentParts, aux)
				}
				if replacementPartCount > 2 {
					nextStaticContentParts = append(nextStaticContentParts, replacementParts[1:replacementPartCount-1]...)
				}
			}
		}
		staticContentParts = nextStaticContentParts
	}
	return
}

type cspNonceBytes struct {
	staticCSPParts     [][]byte
	staticContentParts [][]byte
	contentType        string
	lenText            string
}

func (hnd *cspNonceBytes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(hnd.staticContentParts) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var nonceBytes [cspNonceSize]byte
	var nonceBase64 [cspNonceBase64Size]byte
	rand.Read(nonceBytes[:])
	base64.RawURLEncoding.Encode(nonceBase64[:], nonceBytes[:])
	resultCSPDirectives := bytes.Join(hnd.staticCSPParts, nonceBase64[:])
	resultData := bytes.Join(hnd.staticContentParts, nonceBase64[:])
	hdr := w.Header()
	hdr.Set("Content-Type", hnd.contentType)
	hdr.Set("Content-Length", hnd.lenText)
	hdr.Set("Content-Security-Policy", string(resultCSPDirectives))
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		w.Write(resultData)
	}
}

func NewCSPNonceBytes(data []byte, contentType, cspDirectives, nonceReplacementKey string, contentReplacements []ContentReplacementRule) (hnd http.Handler) {
	nonceReplacementKeyBytes := []byte(nonceReplacementKey)
	cspDirectivesBytes := []byte(cspDirectives)
	staticCSPParts := bytes.Split(cspDirectivesBytes, nonceReplacementKeyBytes)
	staticContentParts := splitCSPNonceContentBytesIntoParts(data, nonceReplacementKeyBytes, contentReplacements)
	var resultContentLen int
	for _, part := range staticContentParts {
		resultContentLen += len(part)
	}
	resultContentLen += ((len(staticContentParts) - 1) * cspNonceBase64Size)
	hnd = &cspNonceBytes{
		staticCSPParts:     staticCSPParts,
		staticContentParts: staticContentParts,
		contentType:        contentType,
		lenText:            strconv.FormatInt(int64(resultContentLen), 10),
	}
	return
}

func NewCSPNonceBytesFromFile(contentFilePath, contentType, cspDirectives, nonceReplacementKey string, contentReplacements []ContentReplacementRule) (hnd http.Handler, err error) {
	data, err := os.ReadFile(contentFilePath)
	if err != nil {
		return
	}
	hnd = NewCSPNonceBytes(data, contentType, cspDirectives, nonceReplacementKey, contentReplacements)
	return
}

// NewCSPNonceBytesFromZipArchive creates a new content handler with CSP Nonce header and content modification from a file within a zip archive.
// The contentFileName specifies the relative path of the file within the zip archive to be used as content. It must not start with a drive letter or leading slash, and only forward slashes are allowed.
func NewCSPNonceBytesFromZipArchive(zipFilePath, contentFileName, contentType, cspDirectives, nonceReplacementKey string, contentReplacements []ContentReplacementRule) (hnd http.Handler, err error) {
	fp, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return
	}
	defer fp.Close()
	var data []byte
	for _, f := range fp.File {
		// Name: It must be a relative path: it must not start with a drive
		// letter (e.g. C:) or leading slash, and only forward slashes
		// are allowed. (https://pkg.go.dev/archive/zip#FileHeader)
		if f.Name == contentFileName {
			var reader io.ReadCloser
			reader, err = f.Open()
			if err != nil {
				return
			}
			defer reader.Close()
			data, err = io.ReadAll(reader)
			if err != nil {
				return
			}
			break
		}
	}
	hnd = NewCSPNonceBytes(data, contentType, cspDirectives, nonceReplacementKey, contentReplacements)
	return
}
