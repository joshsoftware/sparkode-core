package isolate

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"path/filepath"
	"time"

	logger "github.com/sirupsen/logrus"
)

var (
	isolateCommand string = "isolate"
	ShellToUse            = "bash"
)

func GenerateRandomID() (boxID string) {
	min := 1
	max := 9999999999

	rand.Seed(time.Now().UnixNano())
	v := rand.Intn(max-min+1) + min
	boxID = fmt.Sprintf("%v", v)
	return
}

const (
	workdir            = "/var/local/lib/isolate/"
	sourceFile         = "run"
	STDIN_FILE_NAME    = "stdin.txt"
	STDOUT_FILE_NAME   = "stdout.txt"
	STDERR_FILE_NAME   = "stderr.txt"
	METADATA_FILE_NAME = "metadata.txt"
	COMPILED_FILE_NAME = "compile_output.txt"
)

var SupportedLanguage = map[int]string{
	1: "Ruby",
	2: "Go",
}

type LanguageDetails struct {
	ID             int
	Name           string
	SourceFile     string
	CompileCommand string
	RunCommand     string
}

var SupportedLanguageSpecs = map[int]LanguageDetails{
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
		RunCommand:     "./sourceCode",
	},
}

//needs filepath as input /var/local/lib/source/1/
func InitializeFile(path string) {
	files := []string{STDIN_FILE_NAME,
		STDOUT_FILE_NAME,
		STDERR_FILE_NAME,
		METADATA_FILE_NAME,
		COMPILED_FILE_NAME,
	}
	for _, fileName := range files {
		fileName = filepath.Join(path, filepath.Base(fileName))
		_, err := exec.Command("touch", fileName).Output()
		if err != nil {
			logger.Debug("Failed to Initialize File", fileName)
		}
		logger.Debug("created file at " + fileName)
	}

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
