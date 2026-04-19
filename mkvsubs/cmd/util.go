package cmd

import (
	"encoding/json"
	"fmt"
	"mkvsubs/subs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrianbaraka/goutils/echo"
)

// uses package level config struct defined in root.go

// main function that handles all logic
func handleFile(filename string) error {
	// get json from mkvmerge of the given file
	// call extract subinfo

	// switch based on results
	// case 1. no subs in file ok
	// case 2. chosen sub is subrip ok
	// case 3. chosen sub is ssa -> convert to subrip ok
	// case 4. chosen sub is imagebased -> remove or leave alone (flag --delete-image-based) ok
	// case 5. chosen sub is unknown -> leave it alone ok
	// case 6. chosen sub does not exist in said file ok

	// case 7. look for external sub in same dir and merge if subrip (flag --merge-external) ok
	// case 8. path to external sub is given merge if subrip. both 7 and 8 ok
	// case 9. touch only passed sub num leave any other intact ok

	// case 10. dryrun

	jsonRes, e1 := getsubInfo(filename)
	if e1 != nil {
		return e1
	}
	subinfoList, e2 := subs.ExtractSubInfo(jsonRes)
	if e2 != nil {
		return e2
	}

	// if it  is not found it is already handled in subtitle case
	requiredSub, _ := subs.GetRequiredSub(subinfoList, config.subtitleNumber)
	subtitleCase := subs.DetermineCase(subinfoList, config.subtitleNumber)

	// handle 7 and 8 first as they are flag based
	if config.externalSub != "" {
		config.Logger.Echof(echo.Blue, echo.Info, "\tMerging  external sub %v to file %v.\n", config.externalSub, filename)
		return handleExtsub(filename, config.externalSub, subinfoList)
	}

	if config.mergeScan {
		// name of the external sub
		externalSub := subs.FileNameNoExtension(filename) + subs.SubExtension(subs.SubText)
		config.Logger.Echof(echo.Blue, echo.Info, "\tChecking and merging external sub %v to file %v.\n", externalSub, filename)
		return handleExtsub(filename, externalSub, subinfoList)
	}

	switch subtitleCase {
	case subs.CaseNoSubtitles:
		m := fmt.Sprintln("\t", filename+". "+subtitleCase.String())
		return fmt.Errorf("%v", m)
	case subs.CaseAlreadySubRip:
		config.Logger.Echoln(echo.Green, echo.Info, "\t", filename+". "+subtitleCase.String())
		return nil
	case subs.CaseMissingTrack:
		m := fmt.Sprintf("\t%v. %v.\nFile has %v subtitle(s). Requested id %v.", filename, subtitleCase, len(subinfoList), config.subtitleNumber)
		return fmt.Errorf("%v", m)
	case subs.CaseUnknownType:
		m := fmt.Sprintf("\t%v. %v. %v.", filename, subtitleCase, subinfoList)
		config.Logger.Echoln(echo.Yellow, echo.Warn, m)
		return nil
	case subs.CaseConvertSSA:
		return convertSSA(filename, requiredSub, subinfoList)
	case subs.CaseImageBased:
		if !config.delImagesubs {
			return nil
		}
		return delAllSubs(filename)
	default:
		return fmt.Errorf("%v", subtitleCase)
	}

}

func success(m ...any) {
	config.Logger.Echoln(echo.Green, echo.Info, m...)
}

// get json from mkvmerge of the given file if an error occurs fatal
func getsubInfo(filename string) ([]byte, error) {
	args := []string{
		"--identification-format", "json", "--identify", filename,
	}
	config.Logger.Echoln(echo.Blue, echo.Trace, "Running the command: ", config.mkvmergeExe, strings.Join(args, " "))
	res, err, exitCode := config.Runner.RunCmd(echo.Trace, config.mkvmergeExe, args...)

	stringres := strings.Join(res, "\n")
	resBytes := []byte(stringres)

	// get the errors from the json object
	errors, e1 := extractErrorJson(resBytes)
	if e1 != nil {
		return nil, e1
	}
	if err != nil && exitCode != 1 {
		return nil, fmt.Errorf("%v", errors)
	}

	config.Logger.Echoln(echo.Green, echo.Trace, stringres)
	// convert to a list of bytes
	return resBytes, nil
}

// extracts the error message from the json returned from the functio getsubinfo
func extractErrorJson(data []byte) (string, error) {
	m := subs.MkvmergeErrorRes{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return "", err
	}
	res := strings.Join(m.Errors, "\n")
	res += strings.Join(m.Warnings, "\n")
	return res, nil
}

func convertSSA(filename string, info subs.SubInfo, subinfolist []subs.SubInfo) error {
	// extract ssa sub to file
	// convert said ssa to srt
	// merge srt to mkv (mindful of flag keep other subs)
	// cleanup of ssa and srt file and old mkv (mindful of backup flag)

	// cleanup
	newASS, e1 := extractSSA(filename, info)
	if e1 != nil {
		return e1
	}
	defer rmFile(newASS)

	// cleanup
	newSrt, err := subs.ConvertSSAtoSRT(newASS, time.Duration(config.delay))
	if err != nil {
		return err
	}
	defer rmFile(newSrt)
	success("\tConverted the SubstationAlpha to SubRip")

	// cleanup
	newMediafile, cleanDir, e := mergeExternalSRT(filename, newSrt, subinfolist)
	if e != nil {
		return e
	}

	return handleBackup(newMediafile, filename, cleanDir)

}

// returns the newsub name
func extractSSA(filename string, info subs.SubInfo) (string, error) {
	// astisub is picky about the extension

	// TODO use in go utils
	newSub := subs.FileNameNoExtension(filename) + subs.SubExtension(subs.SubRich)
	subname := fmt.Sprintf("%v:%v", info.TrackId, newSub)
	//fmt.Println(subname)

	res, err, exitCode := config.Runner.RunCmd(echo.Debug, config.mkvextractExe, filename, "tracks", subname)
	stringres := strings.Join(res, "\n")
	if err != nil && exitCode != 1 {
		return "", fmt.Errorf("%v", stringres)
	}
	config.Logger.Echoln(echo.Green, echo.Debug, "\t", stringres)
	success("\tExtracted the SubstationAlpha subtitle.")
	return newSub, nil
}

// merges the mediafile with the subfile
// returns the new file and new directory to be handled in cleanup
//
//	the newmediafile is in the clean dir
func mergeExternalSRT(mediafile string, subfile string, subinfolist []subs.SubInfo) (string, string, error) {
	// convert the encoding to utf8 if not
	// create clean dir
	// merge

	ok, err := subs.IsUTF8(subfile)
	if err != nil {
		return "", "", err
	}
	if !ok {
		config.Logger.Echoln(echo.Green, echo.Debug, "\t", subfile, "is not in UTF8. converting...")
		// overwrites old file no need for cleanup
		subs.ConvertUTF(subfile)
		success("\tConverted", subfile, "to UTF8")
	}

	newMediafile, cleanDir, e2 := makeCleanDir(mediafile)
	if e2 != nil {
		return "", "", e2
	}

	// disable other subs in the file
	var disableTracksList []string
	for _, track := range subinfolist {
		disableTracksList = append(disableTracksList, "--default-track-flag", fmt.Sprintf("%d:no", track.TrackId))
	}

	//mkvmerge --output "clean/Abbott Elementary S01E01.mkv" -S --default-track-flag 2:no   "Abbott Elementary S01E01 (1).mkv"   --default-track-flag 0:yes --forced-display-flag 0:yes "Abbott Elementary S01E01 (1).srt"
	args := []string{
		"--output", newMediafile,
	}

	if config.keepOthersubs {
		// disable any subs present
		args = append(args, disableTracksList...)

	} else {
		// remove all subs
		args = append(args, "--no-subtitles")
	}

	//main video file
	args = append(args, mediafile)

	// new subtitle with its flags
	args = append(args, "--default-track-flag", "0:yes", "--forced-display-flag", "0:yes", subfile)

	// TODO dry run
	config.Logger.Echoln(echo.Blue, echo.Trace, "Running the command:", config.mkvmergeExe, strings.Join(args, " "))
	// merge the sub
	res, err, exitCode := config.Runner.RunCmd(echo.Debug, config.mkvmergeExe, args...)

	stringres := strings.Join(res, "\n")
	if err != nil && exitCode != 1 {
		return "", "", fmt.Errorf("%v", stringres)
	}
	config.Logger.Echoln(echo.Green, echo.Debug, stringres)
	success("\tMerged the subtitle to new file", newMediafile)

	// returned to cleanup
	return newMediafile, cleanDir, nil
}

// wrapper over mergeextrnalsrt
func handleExtsub(mediafile string, subfile string, subinfolist []subs.SubInfo) error {
	newfile, newdir, e := mergeExternalSRT(mediafile, subfile, subinfolist)
	if e != nil {
		return e
	}

	// if not backup delete the external sub
	if !config.backup {
		rmFile(subfile)
	}

	return handleBackup(newfile, mediafile, newdir)

}

// makes the clean dir in the same dir as the mediafile
// returns the newmediafile cleandir and error
func makeCleanDir(mediafile string) (string, string, error) {
	// make the clean dir
	parent := filepath.Dir(mediafile)
	cleanDir := filepath.Join(parent, "clean")
	e := os.MkdirAll(cleanDir, 0755)
	if e != nil {
		return "", "", e
	}

	newMediafile := filepath.Join(cleanDir, filepath.Base(mediafile))
	return newMediafile, cleanDir, nil
}

// function to remove a file
func rmFile(filename string) {
	config.Logger.Echoln(echo.Blue, echo.Debug, "\tRemoving file: ", filename)
	if err := os.Remove(filename); err != nil {
		config.Logger.Echof(echo.Yellow, echo.Warn, "Could not remove %v: %v\n", filename, err)
	}
}

func moveAndOverwrite(src, dst string) error {

	// Ensure the source exists
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("source file does not exist: %s", src)
	}
	config.Logger.Echoln(echo.Blue, echo.Debug, "\tRemoving file: ", dst)
	err := os.Remove(dst)
	if err != nil {
		return err
	}

	return os.Rename(src, dst)
}

// handles moving of the new file
func handleBackup(newMediafile string, oldMediaFile string, cleanDir string) error {

	// if backup leave as is else remove the old file and move new file
	if config.backup {
		config.Logger.Echof(echo.Green, echo.Info, "\tThe new file is at %v. Old file kept as backup at %v\n", newMediafile, oldMediaFile)
		return nil
	}

	e2 := moveAndOverwrite(newMediafile, oldMediaFile)
	if e2 != nil {
		return e2
	}
	defer rmFile(cleanDir)
	config.Logger.Echof(echo.Green, echo.Info, "\tReplaced the old file. New media file at %v.\n", oldMediaFile)
	return nil
}

// deletes all subtitles in file
func delAllSubs(filename string) error {
	newMediafile, cleanDir, e2 := makeCleanDir(filename)
	if e2 != nil {
		return e2
	}
	args := []string{
		"--output", newMediafile,
	}
	// remove all subs
	args = append(args, "--no-subtitles")
	//main video file
	args = append(args, filename)

	config.Logger.Echoln(echo.Blue, echo.Trace, "Running the command:", config.mkvmergeExe, strings.Join(args, " "))

	// remove all subs
	res, err, exitCode := config.Runner.RunCmd(echo.Debug, config.mkvmergeExe, args...)

	stringres := strings.Join(res, "\n")
	if err != nil && exitCode != 1 {
		return fmt.Errorf("%v", stringres)
	}
	config.Logger.Echoln(echo.Green, echo.Debug, stringres)
	success("\tMultiplexed  to new file", newMediafile)

	return handleBackup(newMediafile, filename, cleanDir)
}
