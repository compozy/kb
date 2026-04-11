package convert

import (
	"encoding/binary"
	"strings"
)

const (
	exifTypeASCII = 2
	exifTypeShort = 3
	exifTypeLong  = 4
)

var imageTagNames = map[uint16]string{
	0x010E: "imageDescription",
	0x010F: "make",
	0x0110: "model",
	0x0112: "orientation",
	0x0131: "software",
	0x0132: "dateTime",
	0x013B: "artist",
	0x8298: "copyright",
	0x8769: "exifIFDPointer",
	0x9003: "dateTimeOriginal",
}

func acceptsImage(ext string, mimeType string) bool {
	switch normalizeExtension(ext) {
	case ".png", ".jpg", ".jpeg", ".tiff", ".bmp":
		return true
	}

	switch normalizeMIMEType(mimeType) {
	case "image/png", "image/jpeg", "image/tiff", "image/bmp":
		return true
	}

	return false
}

func extractImageMetadata(data []byte) map[string]any {
	tiffData := imageTIFFData(data)
	if len(tiffData) == 0 {
		return nil
	}

	metadata := parseTIFFMetadata(tiffData)
	if len(metadata) == 0 {
		return nil
	}

	return metadata
}

func imageTIFFData(data []byte) []byte {
	if len(data) >= 4 {
		if string(data[:4]) == "II*\x00" || string(data[:4]) == "MM\x00*" {
			return data
		}
	}

	if len(data) < 4 || data[0] != 0xFF || data[1] != 0xD8 {
		return nil
	}

	for offset := 2; offset+4 <= len(data); {
		if data[offset] != 0xFF {
			break
		}

		marker := data[offset+1]
		offset += 2

		switch marker {
		case 0xD8, 0xD9:
			continue
		case 0xDA:
			return nil
		}

		if offset+2 > len(data) {
			return nil
		}

		size := int(binary.BigEndian.Uint16(data[offset : offset+2]))
		if size < 2 || offset+size > len(data) {
			return nil
		}

		segment := data[offset+2 : offset+size]
		if marker == 0xE1 && len(segment) >= 6 && string(segment[:6]) == "Exif\x00\x00" {
			return segment[6:]
		}

		offset += size
	}

	return nil
}

func parseTIFFMetadata(data []byte) map[string]any {
	if len(data) < 8 {
		return nil
	}

	order, ok := tiffByteOrder(data)
	if !ok {
		return nil
	}

	ifdOffset := int(order.Uint32(data[4:8]))
	if ifdOffset <= 0 || ifdOffset >= len(data) {
		return nil
	}

	metadata := map[string]any{}
	visited := map[int]struct{}{}
	parseTIFFIFD(data, order, ifdOffset, metadata, visited)

	return metadata
}

func tiffByteOrder(data []byte) (binary.ByteOrder, bool) {
	switch string(data[:2]) {
	case "II":
		return binary.LittleEndian, true
	case "MM":
		return binary.BigEndian, true
	default:
		return nil, false
	}
}

func parseTIFFIFD(data []byte, order binary.ByteOrder, offset int, metadata map[string]any, visited map[int]struct{}) {
	if offset <= 0 || offset+2 > len(data) {
		return
	}
	if _, ok := visited[offset]; ok {
		return
	}
	visited[offset] = struct{}{}

	entryCount := int(order.Uint16(data[offset : offset+2]))
	entryOffset := offset + 2
	for i := 0; i < entryCount; i++ {
		start := entryOffset + i*12
		if start+12 > len(data) {
			return
		}

		tag := order.Uint16(data[start : start+2])
		valueType := order.Uint16(data[start+2 : start+4])
		count := order.Uint32(data[start+4 : start+8])
		valueField := data[start+8 : start+12]

		if tag == 0x8769 {
			parseTIFFIFD(data, order, int(order.Uint32(valueField)), metadata, visited)
			continue
		}

		name, ok := imageTagNames[tag]
		if !ok || strings.HasSuffix(name, "Pointer") {
			continue
		}

		value, ok := readTIFFValue(data, order, valueType, count, valueField)
		if !ok {
			continue
		}

		metadata[name] = value
	}
}

func readTIFFValue(data []byte, order binary.ByteOrder, valueType uint16, count uint32, valueField []byte) (any, bool) {
	if len(valueField) != 4 || count == 0 {
		return nil, false
	}

	width := tiffTypeWidth(valueType)
	if width == 0 {
		return nil, false
	}

	totalSize := int(count) * width
	valueData := valueField
	if totalSize > len(valueField) {
		offset := int(order.Uint32(valueField))
		if offset < 0 || offset+totalSize > len(data) {
			return nil, false
		}
		valueData = data[offset : offset+totalSize]
	} else {
		valueData = valueField[:totalSize]
	}

	switch valueType {
	case exifTypeASCII:
		text := strings.TrimSpace(strings.TrimRight(string(valueData), "\x00"))
		if text == "" {
			return nil, false
		}
		return text, true
	case exifTypeShort:
		if count != 1 || len(valueData) < 2 {
			return nil, false
		}
		return int(order.Uint16(valueData[:2])), true
	case exifTypeLong:
		if count != 1 || len(valueData) < 4 {
			return nil, false
		}
		return int(order.Uint32(valueData[:4])), true
	default:
		return nil, false
	}
}

func tiffTypeWidth(valueType uint16) int {
	switch valueType {
	case exifTypeASCII:
		return 1
	case exifTypeShort:
		return 2
	case exifTypeLong:
		return 4
	default:
		return 0
	}
}
