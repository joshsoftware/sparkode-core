package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/joshsoftware/sparkode-core/isolate"
	"log"
	"os"
	"os/exec"
)

const (
	workdir                            = "/"
	STDIN_FILE_NAME                    = "stdin.txt"
	STDOUT_FILE_NAME                   = "stdout.txt"
	STDERR_FILE_NAME                   = "stderr.txt"
	METADATA_FILE_NAME                 = "metadata.txt"
	ADDITIONAL_FILES_ARCHIVE_FILE_NAME = "additional_files.zip"
)

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

type LanguageDetails struct {
	ID             int
	Name           string
	SourceFile     string
	CompileCommand string
	RunCommand     string
}

//var SupportedLanguage = map[Language]LanguageDetails{
//	LanguageRuby: {
//		ID:         1,
//		Name:       "Ruby (2.7.0)",
//		SourceFile: "script.rb",
//		RunCommand: "/usr/local/ruby-2.7.0/bin/ruby script.rb",
//	},
//	LanguageGo: {
//		ID:             2,
//		Name:           "Go (1.19.1)",
//		SourceFile:     "main.go",
//		CompileCommand: "GOCACHE=/tmp/.cache/go-build /usr/local/go-1.19.1/bin/go build %s main.go",
//		RunCommand:     "./main",
//	},
//}

var (
	InitIsolate      string = "isolate -b %s --init"
	MoveToIsolateDir string = "cd /var/local/lib/isolate/%s"

	FileSTDIN         string = "/var/local/lib/isolate/1/stdin.txt"
	FileSTDOUT        string = "/var/local/lib/isolate/1/stdout.txt"
	FileSTDERR        string = "/var/local/lib/isolate/1/stderr.txt"
	FileCompileOutput string = "/var/local/lib/isolate/1/compile_output.txt"

	MoveToBox string = "/var/local/lib/isolate/1/box"

	RunIsolateCommandForInterpreted string = "isolate -s -b %s  -t 5.0 -x	 1.0 -w 10.0 -E HOME=/tmp -E PATH=\"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\" -d /etc:noexec --run -- /bin/bash %s < /var/local/lib/isolate/%s/stdin.txt > /var/local/lib/isolate/%s/stdout.txt 2> /var/local/lib/isolate/%s/stderr.txt"

	CleanUpIsolate string = "isolate -b #{%s} --cleanup"

	IsolateCommandToCompileExecutable string = "isolate --cg -s -b %s -M /var/local/lib/isolate/5/metadata.txt --stderr-to-stdout -i /dev/null -t 15.0 -x 0 -w 20.0  -E HOME=/tmp -E PATH=\"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\" -d /etc:noexec --run -- /bin/bash compile > /var/local/lib/isolate/%s/compile_output.txt"
	IsolateCommandToRunExecutable     string = "isolate --cg -s -b %s -M /var/local/lib/isolate/5/metadata.txt -t 5.0 -x 1.0 -w 10.0 -E HOME=/tmp -E PATH=\"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\" -d /etc:noexec --run -- /bin/bash run < /var/local/lib/isolate/%s/stdin.txt > /var/local/lib/isolate/%s/stdout.txt 2> /var/local/lib/isolate/%s/stderr.txt"
)

func execute(code string, input string, langSpecs LanguageDetails, boxNum string) {

	if langSpecs.CompileCommand != "" {
		executeCompiledCode(code, input, langSpecs, boxNum)
	}

	executeInterpretedCode(code, input, langSpecs, boxNum)
}

func executeCompiledCode(code string, input string, langSpecs LanguageDetails, boxNum string) {

	//step 1: init isolate

	initIsolate := fmt.Sprintf(InitIsolate, boxNum)
	out, err := executeCommand(initIsolate)
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	//step 2: cd /var/local/lib/isolate/5

	moveToIsolateDir := fmt.Sprintf(MoveToIsolateDir, boxNum)
	out, err = executeCommand(moveToIsolateDir)
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	_, err = os.Create(FileSTDIN)
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	_, err = os.Create(FileSTDOUT)
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	_, err = os.Create(FileSTDERR)
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	_, err = os.Create(FileCompileOutput)
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	out, err = executeCommand("cd box")
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	mainFile, err := os.Create("main.go")
	if err != nil {
		log.Fatal(err)
	}

	defer mainFile.Close()

	_, err = mainFile.WriteString(code)
	if err != nil {
		log.Fatal(err)
	}

	compileFile, err := os.Create("compile")
	if err != nil {
		log.Fatal(err)
	}

	defer compileFile.Close()

	_, err = compileFile.WriteString(langSpecs.CompileCommand)
	if err != nil {
		log.Fatal(err)
	}

	isolateCommandToCompileExecutable := fmt.Sprintf(IsolateCommandToCompileExecutable, boxNum)
	out, err = executeCommand(isolateCommandToCompileExecutable)
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	runFile, err := os.Create("run")
	if err != nil {
		log.Fatal(err)
	}

	defer runFile.Close()

	_, err = runFile.WriteString("./main")
	if err != nil {
		log.Fatal(err)
	}

	isolateCommandToRunExecutable := fmt.Sprintf(IsolateCommandToRunExecutable, boxNum)
	out, err = executeCommand(isolateCommandToRunExecutable)
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	fmt.Println("Out : ", string(out))
}

func executeInterpretedCode(code string, input string, langSpecs LanguageDetails, boxNum string) {
	//step 1: init isolate

	initIsolate := fmt.Sprintf(InitIsolate, boxNum)
	out, err := executeCommand(initIsolate)
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	//step 2: cd /var/local/lib/isolate/5

	fmt.Println("One")
	out, err = executeCommand("pwd")
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	moveToIsolateDir := fmt.Sprintf(MoveToIsolateDir, boxNum)
	out, err = executeCommand(moveToIsolateDir)
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	fmt.Println("two")

	out, err = executeCommand("pwd")
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	_, err = os.Create(FileSTDIN)
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	_, err = os.Create(FileSTDOUT)
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	_, err = os.Create(FileSTDERR)
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	out, err = executeCommand("pwd")
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	out, err = executeCommand("cd box")
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	filepath := "/var/local/lib/isolate/1" + langSpecs.SourceFile
	mainFile, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}

	defer mainFile.Close()

	_, err = mainFile.WriteString(code)
	if err != nil {
		log.Fatal(err)
	}

	filepathrun := "/var/local/lib/isolate/1/" + "run"
	runnFile, err := os.Create(filepathrun)
	if err != nil {
		log.Fatal(err)
	}

	defer runnFile.Close()

	_, err = runnFile.WriteString(langSpecs.RunCommand)
	if err != nil {
		log.Fatal(err)
	}

	rrunIsolateCommandForInterpreted := fmt.Sprintf(RunIsolateCommandForInterpreted, boxNum, filepathrun, boxNum, boxNum, boxNum)
	out, err = executeCommand(rrunIsolateCommandForInterpreted)
	if err != nil {
		fmt.Printf("\n%s", err)
	}

	fmt.Println("Out : ", string(out))
}

func executeCommand(command string) (output string, err error) {
	fmt.Println("Running the command : ", command)
	err, out, errout := Shellout(command)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	fmt.Println("--- stdout ---")
	fmt.Println(out)
	fmt.Println("--- stderr ---")
	fmt.Println(errout)
	return
}

func main() {
	ctx := context.Background()
	isolate.Run(ctx)
	//code := "puts 'hello'"
	//boxNum := "1"
	//input := ""
	//
	//langSpecs := SupportedLanguage[1]
	//if runtime.GOOS == "windows" {
	//	fmt.Println("Can't Execute this on a windows machine")
	//} else {
	//	execute(code, input, langSpecs, boxNum)
	//}
}
