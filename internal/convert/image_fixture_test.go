package convert

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	imagedraw "image/draw"
	"image/jpeg"
	"sort"
	"testing"
)

func makeJPEGWithEXIF(t *testing.T, tags map[uint16]string) []byte {
	t.Helper()

	base := image.NewRGBA(image.Rect(0, 0, 64, 64))
	imagedraw.Draw(base, base.Bounds(), &image.Uniform{C: color.White}, image.Point{}, imagedraw.Src)

	var jpegData bytes.Buffer
	if err := jpeg.Encode(&jpegData, base, &jpeg.Options{Quality: 95}); err != nil {
		t.Fatalf("jpeg.Encode returned error: %v", err)
	}

	exifSegment := buildEXIFAPP1(t, tags)
	data := jpegData.Bytes()

	output := append([]byte{}, data[:2]...)
	output = append(output, exifSegment...)
	output = append(output, data[2:]...)
	return output
}

func buildEXIFAPP1(t *testing.T, tags map[uint16]string) []byte {
	t.Helper()

	keys := make([]int, 0, len(tags))
	for tag := range tags {
		keys = append(keys, int(tag))
	}
	sort.Ints(keys)

	var ifd bytes.Buffer
	if err := binary.Write(&ifd, binary.LittleEndian, uint16(len(keys))); err != nil {
		t.Fatalf("binary.Write entry count returned error: %v", err)
	}

	var values bytes.Buffer
	valueOffset := uint32(8 + 2 + len(keys)*12 + 4)
	for _, key := range keys {
		tag := uint16(key)
		value := append([]byte(tags[tag]), 0)

		if err := binary.Write(&ifd, binary.LittleEndian, tag); err != nil {
			t.Fatalf("binary.Write tag returned error: %v", err)
		}
		if err := binary.Write(&ifd, binary.LittleEndian, uint16(exifTypeASCII)); err != nil {
			t.Fatalf("binary.Write type returned error: %v", err)
		}
		if err := binary.Write(&ifd, binary.LittleEndian, uint32(len(value))); err != nil {
			t.Fatalf("binary.Write count returned error: %v", err)
		}

		if len(value) <= 4 {
			padded := [4]byte{}
			copy(padded[:], value)
			if _, err := ifd.Write(padded[:]); err != nil {
				t.Fatalf("ifd.Write inline value returned error: %v", err)
			}
			continue
		}

		if err := binary.Write(&ifd, binary.LittleEndian, valueOffset); err != nil {
			t.Fatalf("binary.Write value offset returned error: %v", err)
		}
		if _, err := values.Write(value); err != nil {
			t.Fatalf("values.Write returned error: %v", err)
		}
		valueOffset += uint32(len(value))
	}

	if err := binary.Write(&ifd, binary.LittleEndian, uint32(0)); err != nil {
		t.Fatalf("binary.Write next ifd returned error: %v", err)
	}

	var tiff bytes.Buffer
	tiff.Write([]byte{'I', 'I', '*', 0})
	if err := binary.Write(&tiff, binary.LittleEndian, uint32(8)); err != nil {
		t.Fatalf("binary.Write root IFD offset returned error: %v", err)
	}
	if _, err := tiff.Write(ifd.Bytes()); err != nil {
		t.Fatalf("tiff.Write ifd returned error: %v", err)
	}
	if _, err := tiff.Write(values.Bytes()); err != nil {
		t.Fatalf("tiff.Write values returned error: %v", err)
	}

	payload := append([]byte("Exif\x00\x00"), tiff.Bytes()...)

	var segment bytes.Buffer
	segment.Write([]byte{0xFF, 0xE1})
	if err := binary.Write(&segment, binary.BigEndian, uint16(len(payload)+2)); err != nil {
		t.Fatalf("binary.Write APP1 length returned error: %v", err)
	}
	if _, err := segment.Write(payload); err != nil {
		t.Fatalf("segment.Write payload returned error: %v", err)
	}

	return segment.Bytes()
}
