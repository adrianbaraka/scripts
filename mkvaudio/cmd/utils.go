package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/adrianbaraka/goutils/cli"
	"github.com/adrianbaraka/goutils/echo"
)

// TODO structured output

// returns if it has been downmixed, alreadyStereo, beforesize, aftersize (bytes), and error
func handleFile(filename string) (bool, bool, int64, int64, error) {
	oldFileSize, e := fileSize(filename)
	if e != nil {
		return false, false, 0, 0, e
	}

	chanels, er1 := getNumberOfChannels(filename)
	if er1 != nil {
		return false, false, oldFileSize, oldFileSize, er1
	}

	// already stereo
	if chanels <= 2 {
		return false, true, oldFileSize, oldFileSize, nil
	}

	if config.dryRun {
		config.Logger.Echof(echo.Blue, echo.Info, "Would downmix %v to stereo.\n", filename)
		return false, false, oldFileSize, oldFileSize, nil
	}

	// downmix the file
	newfile, cleanDir, er2 := downmix(filename)
	if er2 != nil {
		return false, false, oldFileSize, oldFileSize, er2
	}

	newFileSize, e2 := fileSize(newfile)
	if e2 != nil {
		return true, false, oldFileSize, oldFileSize, e2
	}

	// update the stats
	er3 := updateStats(newfile)
	if er3 != nil {
		return true, false, oldFileSize, newFileSize, er3
	}

	er4 := handleBackup(newfile, filename, cleanDir)
	if er4 != nil {
		return true, false, oldFileSize, newFileSize, er4
	}

	return true, false, oldFileSize, newFileSize, nil
}

// no of audio chanels in the file.
//
//	checks the first audio track
func getNumberOfChannels(filename string) (int, error) {
	err := cli.RequireTools(config.ffprobeExe)
	if err != nil {
		return 0, err
	}

	args := make([]string, 0, 9)
	args = append(args, "-v", "error")                                                                   // suppress to only errors
	args = append(args, "-select_streams", "a:0")                                                        // first audio stream
	args = append(args, "-show_entries", "stream=channels", "-of", "default=noprint_wrappers=1:nokey=1") // get no of channels only
	args = append(args, filename)

	stdout, err, exitCode := config.Runner.RunCmd(echo.Info, config.ffprobeExe, args...)
	if exitCode != 0 {
		return 0, err
	}
	// it is possible there are no audio channels
	if len(stdout) == 0 {
		return 0, fmt.Errorf("The file possibly has no audio files.")
	}
	return strconv.Atoi(stdout[0])
}

// Downmix to stereo the given file
// Returns the newfile location cleandir and error if it occurred
func downmix(filename string) (string, string, error) {
	err := cli.RequireTools(config.ffmpegExe)
	if err != nil {
		return "", "", err
	}

	cleanFile, cleanDir, err := makeCleanDir(filename)

	if err != nil {
		return "", "", err
	}

	args := make([]string, 0, 15)
	if config.Runner.LogLevel < echo.Trace {
		args = append(args, "-loglevel", "error") // suppress to only errors
	}

	if config.Runner.LogLevel >= echo.Info {
		args = append(args, "-stats") // show progress
	}
	args = append(args, "-i", filename)
	args = append(args, "-c:v", "copy", "-c:s", "copy") // copy video and subs
	args = append(args, "-ac", "2", "-c:a", "aac")      // downmix to 2 channels
	args = append(args, cleanFile)                      // output file
	args = append(args, "-y")                           // overwrite

	err, exitCode := config.Runner.RunCmdStreamer(echo.Info, config.ffmpegExe, args...)
	//err, exitCode := r.RunCmd(echo.Info, "ffmpeg", args...)
	if exitCode != 0 {
		return "", "", err
	}

	return cleanFile, cleanDir, nil
}

func updateStats(filename string) error {
	err := cli.RequireTools(config.mkvpropeditExe)
	if err != nil {
		return err
	}
	//mkvpropedit --add-track-statistics-tags output.mkv
	args := make([]string, 0, 2)
	args = append(args, "--add-track-statistics-tags") // suppress to only errors
	args = append(args, filename)

	_, err, exitCode := config.Runner.RunCmd(echo.Debug, config.mkvpropeditExe, args...)
	if exitCode != 0 {
		return err
	}
	return nil
}

// get the file size
func fileSize(filename string) (int64, error) {
	fileinfo, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return fileinfo.Size(), nil
}

// get percent reduction
func reduction(before, after float64) float64 {
	if before == 0 || after == 0 {
		return 0
	}
	return ((before - after) / before) * 100
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
		config.Logger.Echof(echo.Yellow, echo.Warn, "\tCould not remove %v: %v\n", filename, err)
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

// return OS specific executable name
func getExe(executable string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%v.exe", executable)
	}
	return executable
}
