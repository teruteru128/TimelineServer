package utils

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/TinyKitten/TimelineServer/logger"
	"go.uber.org/zap"
)

var supportBase64MIME = map[string]string{
	".png":  "image/png",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
}

var ErrFileNotSupported = errors.New("base64 file uris is not supported")

func DetectFileExtension(str string) string {
	fileExt := ""
	for ext, mime := range supportBase64MIME {
		if strings.HasPrefix(str, "data:"+mime) {
			fileExt = ext
		}
	}

	return fileExt
}

func DecodeImage(str string) ([]byte, error) {
	trimed := ""
	for _, mime := range supportBase64MIME {
		if strings.HasPrefix(str, "data:"+mime) {
			trimed = strings.Replace(str, "data:"+mime+";base64,", "", 1)
		}
	}

	if trimed == "" {
		return nil, ErrFileNotSupported
	}

	dat, err := base64.StdEncoding.DecodeString(trimed)
	if err != nil {
		logger := logger.NewLogger()
		logger.Debug("Base64 Error", zap.String("Error", err.Error()))
		return nil, err
	}
	return dat, err
}
