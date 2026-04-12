package subs

type SubCase int

const (
	CaseNoSubtitles SubCase = iota
	CaseAlreadySubRip
	CaseConvertSSA
	CaseImageBased
	CaseUnknownType
	CaseMissingTrack
)

// stringer interface
func (c SubCase) String() string {
	switch c {
	case CaseNoSubtitles:
		return "No subtitles found in file"
	case CaseAlreadySubRip:
		return "Subtitle is already in SubRip (SRT) format"
	case CaseConvertSSA:
		return "Substation Alpha (SSA/ASS) detected."
	case CaseImageBased:
		return "Image-based subtitle (PGS/VobSub) detected,"
	case CaseUnknownType:
		return "Unknown subtitle codec."
	case CaseMissingTrack:
		return "The requested subtitle track index does not exist"
	default:
		return "Undefined subtitle case"
	}
}

// returns the required subinfo from the list. if it is not there false is returned
// TODO write tests
func GetRequiredSub(info []SubInfo, subnum int) (SubInfo, bool) {
	for _, v := range info {
		if v.SubtitleNumber == subnum {
			return v, true
		}
	}
	return SubInfo{}, false
}

func DetermineCase(info []SubInfo, subnum int) SubCase {
	// case 1
	if len(info) == 0 {
		return CaseNoSubtitles
	}

	requiredSub, ok := GetRequiredSub(info, subnum)
	if !ok {
		return CaseMissingTrack
	}
	switch requiredSub.Type {
	case SubText:
		return CaseAlreadySubRip
	case SubRich:
		return CaseConvertSSA
	case SubImage:
		return CaseImageBased
	default:
		return CaseUnknownType
	}
}
