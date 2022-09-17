package isolate

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
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
	max := 99999

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

const (
	LanguageRuby       string = "rb"
	LanguageGo         string = "go"
	LanguageC          string = "c"
	LanguageCPP        string = "c++"
	LanguageCSharp     string = "c#"
	LanguageJava       string = "java"
	LanguageJavaScript string = "js"
	LanguagePython     string = "py"
)

//{
//id: 50,
//name: "C (GCC 9.2.0)",
//is_archived: false,
//source_file: "main.c",
//compile_cmd: "/usr/local/gcc-9.2.0/bin/gcc %s main.c",
//run_cmd: "./a.out"
//},
//{
//id: 51,
//name: "C# (Mono 6.6.0.161)",
//is_archived: false,
//source_file: "Main.cs",
//compile_cmd: "/usr/local/mono-6.6.0.161/bin/mcs %s Main.cs",
//run_cmd: "/usr/local/mono-6.6.0.161/bin/mono Main.exe"
//},
//{
//id: 52,
//name: "C++ (GCC 7.4.0)",
//is_archived: false,
//source_file: "main.cpp",
//compile_cmd: "/usr/local/gcc-7.4.0/bin/g++ %s main.cpp",
//run_cmd: "LD_LIBRARY_PATH=/usr/local/gcc-7.4.0/lib64 ./a.out"
//},
//{
//id: 60,
//name: "Go (1.13.5)",
//is_archived: false,
//source_file: "main.go",
//compile_cmd: "GOCACHE=/tmp/.cache/go-build /usr/local/go-1.13.5/bin/go build %s main.go",
//run_cmd: "./main"
//},
//{
//id: 62,
//name: "Java (OpenJDK 13.0.1)",
//is_archived: false,
//source_file: "Main.java",
//compile_cmd: "/usr/local/openjdk13/bin/javac %s Main.java",
//run_cmd: "/usr/local/openjdk13/bin/java Main"
//},
//{
//id: 63,
//name: "JavaScript (Node.js 12.14.0)",
//is_archived: false,
//source_file: "script.js",
//run_cmd: "/usr/local/node-12.14.0/bin/node script.js"
//},
//{
//id: 71,
//name: "Python (3.8.1)",
//is_archived: false,
//source_file: "script.py",
//run_cmd: "/usr/local/python-3.8.1/bin/python3 script.py"
//},
//{
//id: 72,
//name: "Ruby (2.7.0)",
//is_archived: false,
//source_file: "script.rb",
//run_cmd: "/usr/local/ruby-2.7.0/bin/ruby script.rb"
//}

type JobType string

const (
	JobCompile JobType = "compile"
	JobRun     JobType = "run"
)

var LanguageNameToJobType = map[string]JobType{
	LanguageRuby:       JobRun,
	LanguageGo:         JobCompile,
	LanguageC:          JobCompile,
	LanguageCPP:        JobCompile,
	LanguageCSharp:     JobCompile,
	LanguageJava:       JobCompile,
	LanguageJavaScript: JobRun,
	LanguagePython:     JobRun,
}

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
		CompileCommand: "/usr/local/go-1.19.1/bin/go build  main.go",
		RunCommand:     "./main",
	},
}

// needs filepath as input /var/local/lib/source/1/
func InitializeFile(path string) (err error) {
	files := []string{STDIN_FILE_NAME,
		STDOUT_FILE_NAME,
		STDERR_FILE_NAME,
		METADATA_FILE_NAME,
		COMPILED_FILE_NAME,
	}
	for _, fileName := range files {
		fileName = filepath.Join(path, filepath.Base(fileName))
		_, err = exec.Command("touch", fileName).Output()
		if err != nil {
			logger.Debug("Failed to Initialize File", fileName)
			return
		}
		logger.Debug("created file at " + fileName)
	}
	return
}

func Cleanup(ctx context.Context, boxId string) (err error) {
	dirBytes, err := exec.Command(isolateCommand, "--cg", "-b", boxId, "--cleanup").Output()
	if err != nil {
		logger.Debug(ctx, "Isolate : Failed Cleanup", err.Error())
		fmt.Println("cleanup failed", err.Error())
		return
	}
	fmt.Println("cleanup done")
	fmt.Println(string(dirBytes))
	return
}

func ReadFromFile(filePath string) (output string, err error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	output = string(content)
	return
}
