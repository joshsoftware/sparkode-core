package isolateutil

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/joshsoftware/sparkode-core/config"
	logger "github.com/sirupsen/logrus"
)

const (
	workdir            = "/var/local/lib/isolate/"
	STDIN_FILE_NAME    = "stdin.txt"
	STDOUT_FILE_NAME   = "stdout.txt"
	STDERR_FILE_NAME   = "stderr.txt"
	METADATA_FILE_NAME = "metadata.txt"
)

var (
	isolateCommand string = "isolate"
)

func Run(ctx context.Context, sourceFile string) (string, error) {
	boxId := strconv.Itoa(rand.Intn(200-1) + 1)
	defer func() {
		fmt.Println("running cleanup")
		Cleanup(ctx, boxId)
	}()
	_, err := exec.Command(isolateCommand, "--init", "--cg", "-b", boxId).Output()
	if err != nil {
		logger.Debug(ctx, "Isolate : Failed Run", err.Error())
	}

	// initialize files
	InitializeFile(boxId)

	// copy all files to this dir
	runCfg := config.RunConfig{
		TimeLimit:   5.0,
		WallLimit:   10.2,
		MemoryLimit: 128000,
	}
	tmp := os.TempDir()

	metaFile := filepath.Join(tmp, "meta.txt")
	command := createCMD(runCfg, filepath.Join(workdir, boxId), metaFile, sourceFile)

	err, out, errout := Shellout(command)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	fmt.Println("--- stdout ---")
	fmt.Println(out)
	fmt.Println("--- stderr ---")
	fmt.Println(errout)
	fmt.Println("created box", boxId)
	return "", err
}

func Cleanup(ctx context.Context, boxId string) {
	dirBytes, err := exec.Command(isolateCommand, "--cg", "-b", boxId, "--cleanup").Output()
	if err != nil {
		logger.Debug(ctx, "Isolate : Failed Cleanup", err.Error())
		fmt.Println("cleanup failed", err.Error())
		return
	}
	fmt.Println("cleanup done")
	fmt.Println(string(dirBytes))
}

const ShellToUse = "bash"

func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

// create command to --run
func createCMD(cfg config.RunConfig, workdir, metaFile, sourceFile string) string {

	cmd := fmt.Sprintf(
		"isolate --cg -s -b 1 -M %s/metadata.txt -t %.1f -x 1.0 -w %.1f -k 64000 -p60 --cg-timing --cg-mem=%d -f 1024 -E HOME=/tmp -E PATH=\"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"  -d /etc:noexec --run -- /bin/bash %s < %s/stdin.txt > %s/stdout.txt 2> %s/stderr.txt", workdir, cfg.TimeLimit, cfg.WallLimit, cfg.MemoryLimit, sourceFile, workdir, workdir, workdir)

	return cmd
}

func InitializeFile(boxId string) {
	files := []string{STDIN_FILE_NAME,
		STDOUT_FILE_NAME,
		STDERR_FILE_NAME,
		METADATA_FILE_NAME,
	}
	for _, fileName := range files {
		fileName = filepath.Join(workdir, boxId, filepath.Base(fileName))
		_, err := exec.Command("touch", fileName).Output()
		if err != nil {
			logger.Debug("Failed to InitializeFile", fileName)
		}
		logger.Debug("crated file at " + fileName)
	}

}
