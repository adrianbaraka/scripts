package subs_test

import (
	"log"
	"mkvsubs/subs"
	"os"
	"slices"
	"testing"
)

// reads the test json files
func readTestData(name string) []byte {
	t, err := os.ReadFile(name)

	if err != nil {
		log.Fatal(err)
	}
	return t
}

func TestExtractSubInfo(t *testing.T) {
	var t1 []subs.SubInfo
	t1 = append(t1, subs.GetSubinfo(2, "SubRip/SRT", 1))

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		data    []byte
		want    []subs.SubInfo
		wantErr bool
	}{
		//TODO add more test cases
		{name: "Empty Json", wantErr: true},
		{name: "Normal case", data: readTestData("testdata/1.json"), want: t1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := subs.ExtractSubInfo(tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ExtractSubInfo() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ExtractSubInfo() succeeded unexpectedly")
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("ExtractSubInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToUTF8(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		data    []byte
		want    []byte
		wantErr bool
	}{
		{name: "Base case", data: []byte{98}, want: []byte{98}},
		{name: "Windows-1252 Euro", data: []byte{0x80}, want: []byte{0xE2, 0x82, 0xAC}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := subs.ToUTF8(tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ToUTF8() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ToUTF8() succeeded unexpectedly")
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("ToUTF8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategorizeSub(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		codec string
		want  subs.SubtitleType
	}{
		{name: "Base", codec: "SubRip/SRT", want: subs.SubText},
		{name: "Empty", want: subs.SubUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subs.CategorizeSub(tt.codec)
			if got != tt.want {
				t.Errorf("CategorizeSub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubtitleType_Extension(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		codec string
		want  string
	}{
		{"SRT file", "SubRip/SRT", ".srt"},
		{"Rich text", "SubStationAlpha", ".ass"},
		{"Image ", "hdmv_pgs_subtitle", ".sup"},
		{"Empty", "", ".txt"},
		{"Unknown codec", "randomstring", ".txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := subs.CategorizeSub(tt.codec)
			got := subs.SubExtension(s)

			if got != tt.want {
				t.Errorf("Extension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileNameNoExtension(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		filename string
		want     string
	}{
		{"Base", "/home/abc/document.pdf", "/home/abc/document"},
		{"Windows", `C:\Users\abc\document.pdf`, `C:\Users\abc\document`},
		{name: "empty"},
		{"Multi dot", "file.name.tar.gz", "file.name.tar"},
		{"No extension", "/etc/passwd", "/etc/passwd"},
		{"Root Path", "/", "/"},
		{"Hidden file", "/home/user/.config", "/home/user/.config"},
		{"Double hidden file", "/home/.render.yaml", "/home/.render"},
		{"Trailing slash", "/home/", "/home/"},
		{"Trailing slash windows", `C:\Users\`, `C:\Users\`},
		{"Only extension", ".mkv", ".mkv"},
		{"Relative Parent", "../../video.mkv", "../../video"},
		{"Current Dir", "./video.mkv", "video"},
		{"spaces", " file new .mkv", " file new "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subs.FileNameNoExtension(tt.filename)
			if got != tt.want {
				t.Errorf("FileNameNoExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}
