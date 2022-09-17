package isolate

import (
	"bytes"
	"context"
	"fmt"
	"github.com/joshsoftware/sparkode-core/config"
	logger "github.com/sirupsen/logrus"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	workdir                  = "/var/local/lib/isolate/"
	sourceFile               = "run"
	STDIN_FILE_NAME          = "stdin.txt"
	STDOUT_FILE_NAME         = "stdout.txt"
	STDERR_FILE_NAME         = "stderr.txt"
	METADATA_FILE_NAME       = "metadata.txt"
	COMPILE_OUTPUT_FILE_NAME = "compile_output.txt"
)

var (
	isolateCommand string = "isolate"
)

//func Run(ctx context.Context) {
//	boxId := "123456789"
//	defer func() {
//		fmt.Println("running cleanup")
//		Cleanup(ctx, boxId)
//	}()
//	_, err := exec.Command(isolateCommand, "--init", "--cg", "-b", boxId).Output()
//	if err != nil {
//		logger.Debug(ctx, "Isolate : Failed Run", err.Error())
//	}
//
//	// initialize files
//	InitializeFile(boxId)
//
//	CreateSourceFile("script.rb", boxId)
//
//	runCfg := config.RunConfig{
//		TimeLimit:   5.0,
//		WallLimit:   10.2,
//		MemoryLimit: 128000,
//	}
//
//	command := createCMD(runCfg, filepath.Join(workdir, boxId), boxId)
//
//	err, out, errout := Shellout(command)
//	if err != nil {
//		log.Printf("error: %v\n", err)
//	}
//	fmt.Println("--- stdout ---")
//	fmt.Println(out)
//	fmt.Println("--- stderr ---")
//	fmt.Println(errout)
//	fmt.Println("created box", boxId)
//	return
//}

func Run(ctx context.Context) {
	boxId := "123456789"
	defer func() {
		fmt.Println("running cleanup")
		//Cleanup(ctx, boxId)
	}()
	_, err := exec.Command(isolateCommand, "--init", "--cg", "-b", boxId).Output()
	if err != nil {
		logger.Debug(ctx, "Isolate : Failed Run", err.Error())
	}

	// initialize files
	InitializeFile(boxId)

	CreateSourceFile("source.go", boxId)

	runCfg := config.RunConfig{
		TimeLimit:   5.0,
		WallLimit:   10.2,
		MemoryLimit: 128000,
	}

	command := CreateCompileCMD(runCfg, filepath.Join(workdir, boxId), boxId)
	fmt.Println("Command filename: ", command)

	err, out, errout := Shellout(command)
	if err != nil {
		fmt.Println("faild to run compile command", err.Error())
		return
	}
	fmt.Println("Successfully executed compile comamnd")
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	fmt.Println("--- stdout ---")
	fmt.Println(out)
	fmt.Println("--- stderr ---")
	fmt.Println(errout)
	fmt.Println("created box", boxId)

	runCfg = config.RunConfig{
		TimeLimit:   5.0,
		WallLimit:   10.2,
		MemoryLimit: 128000,
	}

	command = createRunCMD(runCfg, filepath.Join(workdir, boxId), boxId)
	fmt.Println("final command filename: ", command)
	if err != nil {
		fmt.Println("faild to run  command", err.Error())
	}
	fmt.Println("Successfully executed run comamnd")
	err, out, errout = Shellout(command)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	fmt.Println("--- stdout ---")
	fmt.Println(out)
	fmt.Println("--- stderr ---")
	fmt.Println(errout)
	fmt.Println("created box", boxId)
	return
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
func createRunCMD(cfg config.RunConfig, workdir, boxId string) string {

	cmd := fmt.Sprintf(
		"isolate --cg -s -b %s -M %s/metadata.txt -t %.1f -x 1.0 -w %.1f -k 64000 -p60 --cg-timing --cg-mem=%d -f 1024 -E HOME=/tmp -E PATH=\"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"  -d /etc:noexec --run -- /bin/bash %s < %s/stdin.txt > %s/stdout.txt 2> %s/stderr.txt", boxId, workdir, cfg.TimeLimit, cfg.WallLimit, cfg.MemoryLimit, sourceFile, workdir, workdir, workdir)

	return cmd
}

// create command to --compile
func CreateCompileCMD(cfg config.RunConfig, workdir, boxId string) string {
	cmd := fmt.Sprintf(
		"isolate --cg -s -b %s -M %s/metadata.txt  -i /dev/null -t %.1f -x 1.0 -w %.1f -k 128000 -p60 --cg-timing --cg-mem=%d -f 4009 -E HOME=/tmp -E PATH=\"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"  -d /etc:noexec --run -- /bin/bash compile > %s/compile_output.txt", boxId, workdir, cfg.TimeLimit, cfg.WallLimit, cfg.MemoryLimit, workdir)

	return cmd
}

func InitializeFile(boxId string) {
	files := []string{STDIN_FILE_NAME,
		STDOUT_FILE_NAME,
		STDERR_FILE_NAME,
		METADATA_FILE_NAME,
		COMPILE_OUTPUT_FILE_NAME,
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

//
//func CreateSourceFile(fileName, boxId string) {
//	fileName = filepath.Join(workdir, boxId, "box", filepath.Base(fileName))
//	_, err := exec.Command("touch", fileName).Output()
//	if err != nil {
//		fmt.Println("Failed to InitializeFile", fileName)
//	}
//	fmt.Println("crated file at ", fileName)
//	d1 := []byte("puts 'devdoot'")
//	err = os.WriteFile(fileName, d1, 0644)
//	if err != nil {
//		fmt.Println("Failed to InitializeFile", fileName)
//	}
//	fmt.Println("crated file at ", fileName)
//
//	//////////////
//	fileName = filepath.Join(workdir, boxId, "box", filepath.Base("run"))
//	_, err = exec.Command("touch", fileName).Output()
//	if err != nil {
//		fmt.Println("Failed to InitializeFile", fileName)
//	}
//	fmt.Println("crated file at ", fileName)
//	d1 = []byte("/usr/local/ruby-2.7.0/bin/ruby script.rb")
//	err = os.WriteFile(fileName, d1, 0644)
//	if err != nil {
//		fmt.Println("Failed to InitializeFile", fileName)
//	}
//	fmt.Println("crated file at ", fileName)
//}
//
//

//
//isolate --cg -b 5 --init
//cd /var/local/lib/isolate/5
//touch stdin.txt
//touch stdout.txt
//touch stderr.txt
//touch compile_output.txt
//cd box
//
//echo "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Print(\"hello\");\n\n}" >> main.go
//echo "GOCACHE=/tmp/.cache/go-build /usr/local/go-1.19.1/bin/go build main.go" >> compile
//
//isolate --cg -s -b 5 -M /var/local/lib/isolate/5/metadata.txt --stderr-to-stdout -i /dev/null -t 15.0 -x 0 -w 20.0 -k 128000 -p120 --cg-timing --cg-mem=512000 -f 4096 -E HOME=/tmp -E PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" -E LANG -E LANGUAGE -E LC_ALL -E JUDGE0_HOMEPAGE -E JUDGE0_SOURCE_CODE -E JUDGE0_MAINTAINER -E JUDGE0_VERSION -d /etc:noexec --run -- /bin/bash compile > /var/local/lib/isolate/5/compile_output.txt
//
//echo "./main" >> run
//
//isolate --cg -s -b 5 -M /var/local/lib/isolate/5/metadata.txt -t 5.0 -x 1.0 -w 10.0 -k 64000 -p60 --cg-timing --cg-mem=128000 -f 1024 -E HOME=/tmp -E PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" -E LANG -E LANGUAGE -E LC_ALL -E JUDGE0_HOMEPAGE -E JUDGE0_SOURCE_CODE -E JUDGE0_MAINTAINER -E JUDGE0_VERSION -d /etc:noexec --run -- /bin/bash run < /var/local/lib/isolate/5/stdin.txt > /var/local/lib/isolate/5/stdout.txt 2> /var/local/lib/isolate/5/stderr.txt

func CreateSourceFile(fileName, boxId string) {
	fileName = filepath.Join(workdir, boxId, "box", filepath.Base(fileName))
	fmt.Println("First filename: ", fileName)
	_, err := exec.Command("touch", fileName).Output()
	if err != nil {
		fmt.Println("Failed to InitializeFile", fileName)
	}
	fmt.Println("crated file at ", fileName)
	d1 := []byte("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, 世界\")\n}\n")
	err = os.WriteFile(fileName, d1, 0644)
	if err != nil {
		fmt.Println("Failed to InitializeFile", fileName)
	}
	fmt.Println("crated file at ", fileName)

	//////////////
	fileName = filepath.Join(workdir, boxId, "box", filepath.Base("compile"))
	fmt.Println("second filename: ", fileName)

	_, err = exec.Command("touch", fileName).Output()
	if err != nil {
		fmt.Println("Failed to InitializeFile", fileName)
	}
	fmt.Println("crated file at ", fileName)
	d1 = []byte("/usr/local/go-1.19.1/bin/go build  source.go")
	err = os.WriteFile(fileName, d1, 0644)
	if err != nil {
		fmt.Println("Failed to InitializeFile", fileName)
	}
	fmt.Println("crated file at ", fileName)

	//////////////
	fileName = filepath.Join(workdir, boxId, "box", filepath.Base("run"))
	fmt.Println("third filename: ", fileName)

	_, err = exec.Command("touch", fileName).Output()
	if err != nil {
		fmt.Println("Failed to InitializeFile", fileName)
	}
	fmt.Println("crated file at ", fileName)
	d1 = []byte("./source")
	err = os.WriteFile(fileName, d1, 0644)
	if err != nil {
		fmt.Println("Failed to InitializeFile", fileName)
	}
	fmt.Println("crated file at ", fileName)
}
