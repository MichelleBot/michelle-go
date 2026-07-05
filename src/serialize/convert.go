package serialize

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const (
	StickerSize    = "512:512"
	StickerFPS     = "10"
	StickerMaxSecs = 6
	StickerQuality = "50"
)

var stickerVF = "scale=512:512:force_original_aspect_ratio=decrease,pad=512:512:(ow-iw)/2:(oh-ih)/2:color=0x00000000"

func tmpFile(prefix, ext string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s%d_%d%s", prefix, os.Getpid(), rand.Int63(), ext))
}

func ToStaticWebp(input []byte, ext string) ([]byte, error) {
	tmpIn := tmpFile("stk_in_", ext)
	tmpOut := tmpFile("stk_out_", ".webp")
	defer os.Remove(tmpIn)
	defer os.Remove(tmpOut)

	if err := os.WriteFile(tmpIn, input, 0600); err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}

	cmd := exec.Command("ffmpeg", "-y",
		"-i", tmpIn,
		"-vf", stickerVF,
		"-vframes", "1",
		"-vcodec", "libwebp",
		"-quality", StickerQuality,
		"-compression_level", "6",
		tmpOut,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg static webp: %w\n%s", err, string(out))
	}
	return os.ReadFile(tmpOut)
}

func ToAnimatedWebp(input []byte, ext string, trim bool) ([]byte, error) {
	tmpIn := tmpFile("stka_in_", ext)
	tmpOut := tmpFile("stka_out_", ".webp")
	defer os.Remove(tmpIn)
	defer os.Remove(tmpOut)

	if err := os.WriteFile(tmpIn, input, 0600); err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}

	vf := stickerVF + ",fps=" + StickerFPS
	maxSecs := StickerMaxSecs
	if !trim {
		maxSecs = 6
	}

	cmd := exec.Command("ffmpeg", "-y",
		"-i", tmpIn,
		"-t", strconv.Itoa(maxSecs),
		"-vf", vf,
		"-vcodec", "libwebp_anim",
		"-loop", "0",
		"-quality", StickerQuality,
		"-compression_level", "6",
		"-qmin", "30",
		"-qmax", "60",
		"-preset", "default",
		"-an",
		tmpOut,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg libwebp_anim: %w\n%s", err, string(out))
	}

	data, err := os.ReadFile(tmpOut)
	if err != nil {
		return nil, err
	}

	if len(data) < 12 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WEBP" {
		return nil, fmt.Errorf("invalid webp output")
	}

	durMs := GetAnimatedWebpDurationMs(data)
	maxMs := int64(StickerMaxSecs * 1000)
	if durMs > maxMs {
		trimmed, err := TrimAnimatedWebp(data)
		if err == nil {
			return trimmed, nil
		}
	}

	return data, nil
}

func webpChunk(tag string, data []byte) []byte {
	size := len(data)
	chunk := []byte(tag)
	sz := make([]byte, 4)
	binary.LittleEndian.PutUint32(sz, uint32(size))
	chunk = append(chunk, sz...)
	chunk = append(chunk, data...)
	if size&1 != 0 {
		chunk = append(chunk, 0x00)
	}
	return chunk
}

func putU24(b []byte, n int) {
	b[0] = byte(n)
	b[1] = byte(n >> 8)
	b[2] = byte(n >> 16)
}

func GenerateJPEGThumbnail(input []byte, ext string) ([]byte, error) {
	tmpIn := tmpFile("thumb_in_", ext)
	tmpOut := tmpFile("thumb_out_", ".jpg")
	defer os.Remove(tmpIn)
	defer os.Remove(tmpOut)

	if err := os.WriteFile(tmpIn, input, 0600); err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}
	cmd := exec.Command("ffmpeg", "-y",
		"-ss", "0.1",
		"-i", tmpIn,
		"-vframes", "1",
		"-vf", "scale=72:72:force_original_aspect_ratio=decrease",
		"-f", "image2",
		"-vcodec", "mjpeg",
		"-q:v", "5",
		tmpOut,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		cmd2 := exec.Command("ffmpeg", "-y",
			"-i", tmpIn,
			"-vframes", "1",
			"-vf", "scale=72:72:force_original_aspect_ratio=decrease",
			"-f", "image2",
			"-vcodec", "mjpeg",
			"-q:v", "5",
			tmpOut,
		)
		if out2, err2 := cmd2.CombinedOutput(); err2 != nil {
			return nil, fmt.Errorf("ffmpeg thumbnail: %w\n%s\n%s", err, string(out), string(out2))
		}
	}
	return os.ReadFile(tmpOut)
}

func GeneratePNGThumbnail(input []byte, ext string) ([]byte, error) {
	tmpIn := tmpFile("pthumb_in_", ext)
	tmpOut := tmpFile("pthumb_out_", ".png")
	defer os.Remove(tmpIn)
	defer os.Remove(tmpOut)

	if err := os.WriteFile(tmpIn, input, 0600); err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}
	cmd := exec.Command("ffmpeg", "-y",
		"-i", tmpIn,
		"-vframes", "1",
		"-vf", "scale=72:72:force_original_aspect_ratio=decrease",
		"-f", "image2",
		"-vcodec", "png",
		tmpOut,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg png thumbnail: %w\n%s", err, string(out))
	}
	return os.ReadFile(tmpOut)
}

type mediaDimensions struct {
	Width  uint32
	Height uint32
}

func GetMediaDimensions(input []byte, ext string) (mediaDimensions, error) {
	tmpIn := tmpFile("dim_in_", ext)
	defer os.Remove(tmpIn)

	if err := os.WriteFile(tmpIn, input, 0600); err != nil {
		return mediaDimensions{}, fmt.Errorf("write input: %w", err)
	}
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_streams",
		"-select_streams", "v:0",
		tmpIn,
	)
	out, err := cmd.Output()
	if err != nil {
		return mediaDimensions{}, fmt.Errorf("ffprobe: %w", err)
	}
	var probe struct {
		Streams []struct {
			Width  uint32 `json:"width"`
			Height uint32 `json:"height"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(out, &probe); err != nil {
		return mediaDimensions{}, fmt.Errorf("parse ffprobe: %w", err)
	}
	if len(probe.Streams) == 0 {
		return mediaDimensions{}, fmt.Errorf("no video stream found")
	}
	return mediaDimensions{Width: probe.Streams[0].Width, Height: probe.Streams[0].Height}, nil
}

func GetVideoDurationSeconds(input []byte, ext string) (uint32, error) {
	tmpIn := tmpFile("dur_in_", ext)
	defer os.Remove(tmpIn)

	if err := os.WriteFile(tmpIn, input, 0600); err != nil {
		return 0, fmt.Errorf("write input: %w", err)
	}
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		tmpIn,
	)
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe: %w", err)
	}
	var probe struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}
	if err := json.Unmarshal(out, &probe); err != nil {
		return 0, fmt.Errorf("parse ffprobe: %w", err)
	}
	f, err := strconv.ParseFloat(probe.Format.Duration, 64)
	if err != nil {
		return 0, nil
	}
	return uint32(f), nil
}

func GetAnimatedWebpDurationMs(data []byte) int64 {
	if len(data) < 12 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WEBP" {
		return 0
	}
	var total int64
	pos := 12
	for pos+8 <= len(data) {
		tag := string(data[pos : pos+4])
		size := int(binary.LittleEndian.Uint32(data[pos+4 : pos+8]))
		if tag == "ANMF" {
			durOff := pos + 8 + 12
			if durOff+3 <= len(data) {
				total += int64(uint32(data[durOff]) |
					uint32(data[durOff+1])<<8 |
					uint32(data[durOff+2])<<16)
			}
		}
		pos += 8 + size + (size & 1)
	}
	return total
}

func TrimAnimatedWebp(data []byte) ([]byte, error) {
	durMs := GetAnimatedWebpDurationMs(data)
	if durMs == 0 || durMs <= int64(StickerMaxSecs*1000) {
		return data, nil
	}
	if len(data) < 12 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WEBP" {
		return nil, fmt.Errorf("invalid WEBP format")
	}

	const maxMs = int64(StickerMaxSecs * 1000)
	var preChunks, anmfChunks [][]byte
	var cumulativeMs int64
	pos := 12

	for pos+8 <= len(data) {
		tag := string(data[pos : pos+4])
		size := int(binary.LittleEndian.Uint32(data[pos+4 : pos+8]))
		aligned := size + (size & 1)
		end := pos + 8 + aligned
		if end > len(data) {
			end = len(data)
		}
		chunk := make([]byte, end-pos)
		copy(chunk, data[pos:end])

		if tag == "ANMF" {
			durOff := 8 + 12
			if durOff+3 <= len(chunk) {
				dur := int64(uint32(chunk[durOff]) |
					uint32(chunk[durOff+1])<<8 |
					uint32(chunk[durOff+2])<<16)
				if cumulativeMs+dur <= maxMs {
					cumulativeMs += dur
					anmfChunks = append(anmfChunks, chunk)
				}
			}
		} else {
			preChunks = append(preChunks, chunk)
		}
		pos += 8 + aligned
	}

	if len(anmfChunks) == 0 {
		return data, nil
	}

	all := append(preChunks, anmfChunks...)
	payloadSize := 4
	for _, c := range all {
		payloadSize += len(c)
	}

	out := make([]byte, 0, 8+payloadSize)
	out = append(out, "RIFF"...)
	sz := make([]byte, 4)
	binary.LittleEndian.PutUint32(sz, uint32(payloadSize))
	out = append(out, sz...)
	out = append(out, "WEBP"...)
	for _, c := range all {
		out = append(out, c...)
	}
	return out, nil
}

func ToOggOpus(input []byte, ext string) ([]byte, error) {
	tmpIn := tmpFile("ogg_in_", ext)
	tmpOut := tmpFile("ogg_out_", ".ogg")
	defer os.Remove(tmpIn)
	defer os.Remove(tmpOut)

	if err := os.WriteFile(tmpIn, input, 0600); err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}
	cmd := exec.Command("ffmpeg", "-y",
		"-i", tmpIn,
		"-c:a", "libopus",
		"-b:a", "128k",
		"-vbr", "on",
		"-compression_level", "10",
		"-frame_duration", "20",
		"-application", "voip",
		tmpOut,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg ogg opus: %w\n%s", err, string(out))
	}
	return os.ReadFile(tmpOut)
}

func ToJPEG(input []byte, ext string) ([]byte, error) {
	tmpIn := tmpFile("jpg_in_", ext)
	tmpOut := tmpFile("jpg_out_", ".jpg")
	defer os.Remove(tmpIn)
	defer os.Remove(tmpOut)

	if err := os.WriteFile(tmpIn, input, 0600); err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}
	cmd := exec.Command("ffmpeg", "-y",
		"-i", tmpIn,
		"-vframes", "1",
		"-q:v", "2",
		"-f", "image2",
		"-vcodec", "mjpeg",
		tmpOut,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg tojpeg: %w\n%s", err, string(out))
	}
	return os.ReadFile(tmpOut)
}

func ReencodeMP4(input []byte) ([]byte, error) {
	tmpIn := tmpFile("re_in_", ".mp4")
	tmpOut := tmpFile("re_out_", ".mp4")
	defer os.Remove(tmpIn)
	defer os.Remove(tmpOut)

	if err := os.WriteFile(tmpIn, input, 0600); err != nil {
		return nil, fmt.Errorf("write input: %w", err)
	}

	cmd := exec.Command("ffmpeg", "-y",
		"-i", tmpIn,
		"-c:v", "copy",
		"-c:a", "copy",
		"-movflags", "+faststart",
		"-avoid_negative_ts", "make_zero",
		tmpOut,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		cmd2 := exec.Command("ffmpeg", "-y",
			"-i", tmpIn,
			"-c:v", "libx264",
			"-c:a", "aac",
			"-movflags", "+faststart",
			"-preset", "fast",
			"-crf", "23",
			tmpOut,
		)
		if out2, err2 := cmd2.CombinedOutput(); err2 != nil {
			return nil, fmt.Errorf("ffmpeg reencode: %w\n%s\n%s", err, string(out), string(out2))
		}
	}
	return os.ReadFile(tmpOut)
}
