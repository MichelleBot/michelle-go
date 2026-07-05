package serialize

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
)

type StickerMetadata struct {
	PackName   string
	Author     string
	Categories []string
}

func StickerMetadataDefault() StickerMetadata {
	return StickerMetadata{
		PackName:   "WhatsApp Bot",
		Author:     "Siputzx",
		Categories: []string{""},
	}
}

func tmpFileExif(prefix, ext string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s%d_%d%s", prefix, os.Getpid(), rand.Int63(), ext))
}

func buildExifBinary(metadata StickerMetadata) ([]byte, error) {
	cats := metadata.Categories
	if len(cats) == 0 {
		cats = []string{""}
	}

	jsonData := map[string]interface{}{
		"sticker-pack-id":        "https://github.com/siputzx",
		"sticker-pack-name":      metadata.PackName,
		"sticker-pack-publisher": metadata.Author,
		"emojis":                 cats,
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	startingBytes := []byte{0x49, 0x49, 0x2A, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01, 0x00, 0x41, 0x57, 0x07, 0x00}
	endingBytes := []byte{0x16, 0x00, 0x00, 0x00}

	var b bytes.Buffer
	b.Write(startingBytes)

	lenBuffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(lenBuffer, uint32(len(jsonBytes)))
	b.Write(lenBuffer)

	b.Write(endingBytes)
	b.Write(jsonBytes)

	return b.Bytes(), nil
}

func AddExifToWebp(webpData []byte, metadata StickerMetadata) ([]byte, error) {
	if metadata.PackName == "" && metadata.Author == "" {
		return webpData, nil
	}

	tmpIn := tmpFileExif("webp_in_", ".webp")
	tmpOut := tmpFileExif("webp_out_", ".webp")
	tmpExif := tmpFileExif("exif_", ".bin")
	defer os.Remove(tmpIn)
	defer os.Remove(tmpOut)
	defer os.Remove(tmpExif)

	if err := os.WriteFile(tmpIn, webpData, 0600); err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}

	exifData, err := buildExifBinary(metadata)
	if err != nil {
		return nil, fmt.Errorf("build exif: %w", err)
	}

	if err := os.WriteFile(tmpExif, exifData, 0600); err != nil {
		return nil, fmt.Errorf("write exif: %w", err)
	}

	cmd := exec.Command("webpmux", "-set", "exif", tmpExif, tmpIn, "-o", tmpOut)
	if _, err := cmd.CombinedOutput(); err != nil {
		return webpData, nil
	}

	result, err := os.ReadFile(tmpOut)
	if err != nil {
		return webpData, nil
	}
	return result, nil
}

func ToStaticWebpExif(input []byte, ext string, metadata StickerMetadata) ([]byte, error) {
	webpData, err := ToStaticWebp(input, ext)
	if err != nil {
		return nil, err
	}
	return AddExifToWebp(webpData, metadata)
}

func ToAnimatedWebpExif(input []byte, ext string, trim bool, metadata StickerMetadata) ([]byte, error) {
	webpData, err := ToAnimatedWebp(input, ext, trim)
	if err != nil {
		return nil, err
	}
	return AddExifToWebp(webpData, metadata)
}
