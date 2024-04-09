package mime

import (
	"bytes"
	"io"

	"boilerplate/internal/abstraction"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

var allowedMimeTypes = map[string]struct{}{
	"image/jpeg":               {},
	"image/png":                {},
	"image/jpg":                {},
	"application/pdf":          {},
	"application/vnd.ms-excel": {}, // .xls
}

func isXlsFileType(kind *types.Type, buf []byte) {
	xlsMagicNumber := []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}
	if *kind == filetype.Unknown && len(buf) >= len(xlsMagicNumber) && bytes.Equal(buf[:len(xlsMagicNumber)], xlsMagicNumber) {
		*kind = types.Type{
			MIME:      types.NewMIME("application/vnd.ms-excel"),
			Extension: "xls",
		}
	}
}

func getMime(buf []byte) (*types.Type, error) {
	kind, err := filetype.Match(buf)
	if err != nil {
		return nil, err
	}
	if kind == filetype.Unknown {
		isXlsFileType(&kind, buf)
	}
	return &kind, err
}

func AllowedMimeTypeImages(ctx *abstraction.Context, file io.Reader) (*bool, *string, error) {
	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	kind, err := getMime(buf)
	if err != nil {
		return nil, nil, err
	}

	ctx.Request().Write(bytes.NewBuffer(buf))
	_, ok := allowedMimeTypes[kind.MIME.Value]

	return &ok, &kind.MIME.Value, nil
}
