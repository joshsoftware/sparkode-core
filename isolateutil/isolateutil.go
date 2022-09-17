package isolateutil

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"
	"path/filepath"

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

type LanguageDetails struct {
	ID             int
	Name           string
	SourceFile     string
	CompileCommand string
	RunCommand     string
}

var SupportedLanguage map[int]LanguageDetails = map[int]LanguageDetails{
	1: {
		ID:         1,
		Name:       "Ruby (2.7.0)",
		SourceFile: "script.rb",
		RunCommand: "/usr/local/ruby-2.7.0/bin/ruby script.rb",
	},
	2: {
		ID:             2,
		Name:           "Go (1.19.1)",
		SourceFile:     "main.go",
		CompileCommand: "GOCACHE=/tmp/.cache/go-build /usr/local/go-1.19.1/bin/go build %s main.go",
		RunCommand:     "./main",
	},
}

// 1. Initialise isolate
// isolate --cg -b 1 --init

// cd /var/local/lib/isolate/1

// 2. Initilialize files
// touch stdin.txt
// touch stdout.txt
// touch stderr.txt

// 3. Put code input in stdin
// cd box

// 4. Create script.rb and run inside box
// echo "puts 'hello'" >> script.rb
// echo "/usr/local/ruby-2.7.0/bin/ruby script.rb" >> run

// cd ~

// 5. Run isolate
// isolate --cg -s -b 1 -M /var/local/lib/isolate/1/metadata.txt -t 5.0 -x 1.0 -w 10.0 -k 64000 -p60 --cg-timing --cg-mem=128000 -f 1024 -E HOME=/tmp -E PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" -E LANG -E LANGUAGE -E LC_ALL -E JUDGE0_HOMEPAGE -E JUDGE0_SOURCE_CODE -E JUDGE0_MAINTAINER -E JUDGE0_VERSION -d /etc:noexec --run -- /bin/bash run < /var/local/lib/isolate/1/stdin.txt > /var/local/lib/isolate/1/stdout.txt 2> /var/local/lib/isolate/1/stderr.txt



func Run(ctx context.Context, code string, input string, langSpecs LanguageDetails) (string, error) {
	rand.Seed(time.Now().UnixNano())	
	v := rand.Int()
	boxId := fmt.Sprintf("%v", v)

	defer func() {
		fmt.Println("running cleanup")
		// Cleanup(ctx, boxId)
	}()

	fmt.Println("Step 1")
	// 1. Initialise isolate

	initComand := fmt.Sprintf("isolate --cg -b %s --init",boxId)
	fmt.Println(initComand)
	err, out, errout := Shellout(initComand)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
// rr, err := exec.Command(isolateCommand, "--cg", "-b", boxId, "--init").Output()
// 	if err != nil {
// 		logger.Debug(ctx, "Isolate : Failed Run", err.Error())
// 	}

	fmt.Println("Output : ", string(out))
	fmt.Println("Output : ", err)
	fmt.Println("Output : ", errout)
	fmt.Println("Step 2")

	// 2. Initilialize files
	// initialize files
	InitializeFile(boxId)

	//4. Create script.rb and run inside box
	fmt.Println("Step 3")
	
	fileName := filepath.Join(workdir,"box", filepath.Base(langSpecs.SourceFile))
	_, err = exec.Command("touch", fileName).Output()
	if err != nil {
		logger.Debug("Failed to InitializeFile script.rb", fileName)
	}

	fmt.Println("Step 4")

	fileName = filepath.Join(workdir,"box", filepath.Base("run"))
	_, err = exec.Command("touch", fileName).Output()
	if err != nil {
		logger.Debug("Failed to InitializeFile", fileName)
	}

	fmt.Println("Step 5")

	// copy all files to this dir
	runCfg := config.RunConfig{
		TimeLimit:   5.0,
		WallLimit:   10.2,
		MemoryLimit: 128000,
	}

	fmt.Println("Step 6")

	command := createCMD(runCfg, filepath.Join(workdir, boxId), fileName, boxId)

	fmt.Println("Step 7")

	fmt.Println("Final Command : ", command)
	err, out, errout = Shellout(command)
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
func createCMD(cfg config.RunConfig, workdir, sourceFile, boxId string) string {

	cmd := fmt.Sprintf(
		"isolate --cg -s -b %s -M %s/metadata.txt -t %.1f -x 1.0 -w %.1f -k 64000 -p60 --cg-timing --cg-mem=%d -f 1024 -E HOME=/tmp -E PATH=\"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"  -d /etc:noexec --run -- /bin/bash %s < %s/stdin.txt > %s/stdout.txt 2> %s/stderr.txt",boxId, workdir, cfg.TimeLimit, cfg.WallLimit, cfg.MemoryLimit, sourceFile, workdir, workdir, workdir)

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
