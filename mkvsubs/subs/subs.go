package subs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/asticode/go-astisub"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

type SubtitleType string

const (
	SubText    SubtitleType = "text"    // SRT
	SubRich    SubtitleType = "rich"    // ASS/SSA
	SubImage   SubtitleType = "image"   // PGS, VobSub
	SubUnknown SubtitleType = "unknown" // any other
)

// Extension returns the standard file extension for the category.
func SubExtension(s SubtitleType) string {
	switch s {
	case SubText:
		return ".srt"
	case SubRich:
		return ".ass"
	case SubImage:
		return ".sup" // Standard for PGS; use .sub for VobSub if preferred
	default:
		return ".txt"
	}
}

type SubInfo struct {
	TrackId        int
	Codec          string
	SubtitleNumber int
	Type           SubtitleType
}

func (s SubInfo) String() string {
	return fmt.Sprintf(
		"Track ID: %d | Sub Number: %d | Codec: %-10s | Type: %s",
		s.TrackId, s.SubtitleNumber, s.Codec, s.Type,
	)
}

// constructor for subinfo list
func GetSubinfo(trackId int, codec string, subtitleNumber int) SubInfo {
	return SubInfo{
		TrackId:        trackId,
		Codec:          codec,
		SubtitleNumber: subtitleNumber,
		Type:           CategorizeSub(codec),
	}
}

// categorizes a subtitle codec into its subtitle type
func CategorizeSub(codec string) SubtitleType {
	switch codec {
	case "SubRip/SRT":
		return SubText
	case "SubStationAlpha":
		return SubRich
	case "hdmv_pgs_subtitle", "dvd_subtitle":
		return SubImage
	default:
		return SubUnknown
	}
}

// Read the json file returned from mkvmerge and return a list of the type subInfo
func ExtractSubInfo(data []byte) ([]SubInfo, error) {
	m := MkvmergeRes{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	var res []SubInfo
	subsFound := 1

	for i := 0; i < len(m.Tracks); i++ {
		track := m.Tracks[i]
		if track.Type == "subtitles" {
			sub := SubInfo{}
			sub.Codec = track.Codec
			sub.TrackId = track.ID
			sub.SubtitleNumber = subsFound
			sub.Type = CategorizeSub(track.Codec)

			subsFound++

			res = append(res, sub)
		}
	}
	//fmt.Printf("res: %v\n", res)
	return res, nil
}

// check the encoding of passed file true if it is utf8 false otherwise
func IsUTF8(file string) (bool, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return false, err
	}
	if utf8.Valid(f) {
		return true, nil
	}
	return false, nil

}

// convert the passed bytes to utf8
func ToUTF8(data []byte) ([]byte, error) {
	// get encoding
	enc, name, _ := charset.DetermineEncoding(data, "")

	reader := transform.NewReader(bytes.NewReader(data), enc.NewDecoder())
	result, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert from %s to utf-8: %w", name, err)
	}

	return result, nil
}

// converts the given file to utf8
func ConvertUTF(file string) error {
	data, err := os.ReadFile(file)

	if err != nil {
		return err
	}
	converted, er := ToUTF8(data)
	if er != nil {
		return er
	}

	return os.WriteFile(file, converted, 0644)
}

// TODO use the implementation in my go utils
// Returns the given filename stripped of the last extension eg /home/document.pdf /home/document
func FileNameNoExtension(filename string) string {
	if filename == "" {
		return ""
	}
	// handle trailing slash
	if strings.HasSuffix(filename, string(os.PathSeparator)) {
		return filename
	}

	base := filepath.Base(filename)

	ext := filepath.Ext(base)

	// handle hidden files eg .gitignore
	if base == ext {
		ext = ""
	}

	parent := filepath.Dir(filename)

	//	fmt.Println(base)
	val := strings.TrimSuffix(base, ext)

	return filepath.Join(parent, val)
}

// Convert takes an ssa file subtitle typically .ass and converts it to srt file in same dir
//
//	Returns the filename of the created .srt file and any errors that occcurred
func ConvertSSAtoSRT(filename string, delay time.Duration) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".ssa", ".ass":
		// Valid
	default:
		return "", fmt.Errorf("Unsupported subtitle extension %q: expected .ssa or .ass", ext)
	}

	s1, err := astisub.OpenFile(filename)
	if err != nil {
		return "", err
	}

	// add delay
	s1.Add(delay)

	// create the newfile in the same dir with extension .srt
	newfile := FileNameNoExtension(filename) + ".srt"
	//fmt.Println(newfile)
	f, err := os.Create(newfile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return newfile, s1.WriteToSRT(f)
}
