package subs_test

import (
	"mkvsubs/subs"
	"testing"
)

func TestDetermineCase(t *testing.T) {

	t1 := []subs.SubInfo{}
	t1 = append(t1, subs.GetSubinfo(2, "SubRip/SRT", 1))
	t1 = append(t1, subs.GetSubinfo(3, "SubStationAlpha", 2))
	t1 = append(t1, subs.GetSubinfo(4, "dvd_subtitle", 3))

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		info   []subs.SubInfo
		subnum int
		want   subs.SubCase
	}{
		// todo add json file tests for actual testing
		{"no subs", []subs.SubInfo{}, 1, subs.CaseNoSubtitles},
		{"missing track", t1, 10, subs.CaseMissingTrack},
		{"Already subrip", t1, 1, subs.CaseAlreadySubRip},
		{"SSA needs conversion", t1, 2, subs.CaseConvertSSA},
		{"Image based sub", t1, 3, subs.CaseImageBased},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := subs.DetermineCase(tt.info, tt.subnum)
			if got != tt.want {
				t.Errorf("DetermineCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
