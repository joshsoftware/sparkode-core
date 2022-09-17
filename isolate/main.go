package isolate

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joshsoftware/sparkode-core/config"
	"github.com/joshsoftware/sparkode-core/model"
	logger "github.com/sirupsen/logrus"
)

func Run(ctx context.Context, jobType JobType, langSpecs LanguageDetails, codeData model.ExecuteCodeRequest) (stdout string, stderr string, err error) {

	boxID := GenerateRandomID()
	defer func() {
		fmt.Println("running cleanup")
		Cleanup(ctx, boxID)
	}()

	_, err = exec.Command(isolateCommand, "--init", "--cg", "-b", boxID).Output()
	if err != nil {
		logger.Warn(ctx, "Isolate : Failed Run", err.Error())
		err = fmt.Errorf("error occured while initiating isolate : %s", err.Error())
		return
	}

	boxDirPath := fmt.Sprintf("%s%s", workdir, boxID)
	// initialize files
	err = InitializeFile(boxDirPath)
	if err != nil {
		err = fmt.Errorf("error occured while initializing files : %s", err.Error())
		return
	}

	runCfg := config.RunConfig{
		TimeLimit:   5.0,
		WallLimit:   10.2,
		MemoryLimit: 128000,
	}

	switch jobType {
	case JobRun:
		CreateSourceFilesForInterpreted(langSpecs, codeData, boxDirPath)
		if err != nil {
			err = fmt.Errorf("error occured while creating source files : %s", err.Error())
			return
		}
	case JobCompile:
		createSourceFilesForCompile(langSpecs, codeData, boxDirPath)
		if err != nil {
			err = fmt.Errorf("error occured while creating source files : %s", err.Error())
			return
		}

		createCompileCMD(runCfg, filepath.Join(workdir, boxID), boxID)
		if err != nil {
			err = fmt.Errorf("error occured while creating source files : %s", err.Error())
			return
		}
	}

	command := createRunCMD(runCfg, filepath.Join(workdir, boxID), boxID)
	fmt.Println("Final Command : ", command)
	out, errout, err := Shellout(command)
	if err != nil {
		log.Printf("error: %v\n", err)
		err = fmt.Errorf("error occured while creating final command  : %s", err.Error())
		return
	}
	fmt.Println("--- stdout ---")
	fmt.Println(out)
	fmt.Println("--- stderr ---")
	fmt.Println(errout)

	stdoutFilePath := fmt.Sprintf("%s/%s", boxDirPath, STDOUT_FILE_NAME)
	stdout, err = ReadFromFile(stdoutFilePath)
	if err != nil {
		err = fmt.Errorf("error occured while reading from stdout  : %s", err.Error())
		return
	}

	stderrFilePath := fmt.Sprintf("%s/%s", boxDirPath, STDERR_FILE_NAME)
	stderr, err = ReadFromFile(stderrFilePath)
	if err != nil {
		err = fmt.Errorf("error occured while reading from stderr  : %s", err.Error())
		return
	}

	return
}

func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// create command to --run
func createRunCMD(cfg config.RunConfig, workdir, boxId string) string {

	cmd := fmt.Sprintf(
		"isolate --cg -s -b %s -M %s/metadata.txt -t %.1f -x 1.0 -w %.1f -k 64000 -p60 --cg-timing --cg-mem=%d -f 4009 -E HOME=/tmp -E PATH=\"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"  -d /etc:noexec --run -- /bin/bash %s < %s/stdin.txt > %s/stdout.txt 2> %s/stderr.txt", boxId, workdir, cfg.TimeLimit, cfg.WallLimit, cfg.MemoryLimit, sourceFile, workdir, workdir, workdir)

	return cmd
}

// create command to --compile
func createCompileCMD(cfg config.RunConfig, workdir, boxId string) string {
	cmd := fmt.Sprintf(
		"isolate --cg -s -b %s -M %s/metadata.txt  -i /dev/null -t %.1f -x 1.0 -w %.1f -k 128000 -p60 --cg-timing --cg-mem=%d -f 4009 -E HOME=/tmp -E PATH=\"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"  -d /etc:noexec --run -- /bin/bash compile > %s/compile_output.txt", boxId, workdir, cfg.TimeLimit, cfg.WallLimit, cfg.MemoryLimit, workdir)

	return cmd
}

// First Param: languageSpecs, 2nd argument: boxpath
func CreateSourceFilesForInterpreted(langSpecs LanguageDetails, codeData model.ExecuteCodeRequest, boxPath string) (err error) {
	//Create code script
	fileName := filepath.Join(boxPath, "box", filepath.Base(langSpecs.SourceFile))
	_, err = exec.Command("touch", fileName).Output()
	if err != nil {
		fmt.Println("Failed to Initialize File : ", fileName)
	}
	fmt.Println("Created file : ", fileName)

	code := codeData.Code
	err = os.WriteFile(fileName, []byte(code), 0644)
	if err != nil {
		fmt.Println("Failed to Initialize File :", fileName)
	}
	fmt.Println("Created File : ", fileName)

	//Create run file
	fileName = filepath.Join(boxPath, "box", filepath.Base("run"))
	_, err = exec.Command("touch", fileName).Output()
	if err != nil {
		fmt.Println("Failed to InitializeFile", fileName)
	}
	fmt.Println("Created File : ", fileName)

	runCommand := langSpecs.RunCommand
	err = os.WriteFile(fileName, []byte(runCommand), 0644)
	if err != nil {
		fmt.Println("Failed to InitializeFile", fileName)
	}

	//write input into file
	inputFileName := filepath.Join(boxPath, filepath.Base(STDIN_FILE_NAME))
	err = os.WriteFile(inputFileName, []byte(codeData.Input), 0644)
	if err != nil {
		fmt.Println("Failed to InitializeFile", fileName)
	}
	return
}

func createSourceFilesForCompile(langSpecs LanguageDetails, codeData model.ExecuteCodeRequest, boxPath string) (err error) {
	//Create code script
	fileName := filepath.Join(boxPath, "box", filepath.Base(langSpecs.SourceFile))
	_, err = exec.Command("touch", fileName).Output()
	if err != nil {
		fmt.Println("Failed to Initialize File : ", fileName)
		return
	}
	fmt.Println("Created file : ", fileName)

	code := codeData.Code
	err = os.WriteFile(fileName, []byte(code), 0644)
	if err != nil {
		fmt.Println("Failed to Initialize File :", fileName)
		return
	}
	fmt.Println("Created File : ", fileName)

	// Create compile file
	fileName = filepath.Join(boxPath, "box", filepath.Base("compile"))
	_, err = exec.Command("touch", fileName).Output()
	if err != nil {
		fmt.Println("Failed to InitializeFile", fileName)
	}
	fmt.Println("Created File : ", fileName)

	compileCommand := langSpecs.CompileCommand
	err = os.WriteFile(fileName, []byte(compileCommand), 0644)
	if err != nil {
		fmt.Println("Failed to InitializeFile", fileName)
	}

	//Create run file
	fileName = filepath.Join(boxPath, "box", filepath.Base("run"))
	_, err = exec.Command("touch", fileName).Output()
	if err != nil {
		fmt.Println("Failed to InitializeFile", fileName)
		return
	}
	fmt.Println("Created File : ", fileName)

	runCommand := langSpecs.RunCommand
	err = os.WriteFile(fileName, []byte(runCommand), 0644)
	if err != nil {
		fmt.Println("Failed to InitializeFile", fileName)
		return
	}

	//write input into file
	inputFileName := filepath.Join(boxPath, filepath.Base(STDIN_FILE_NAME))
	err = os.WriteFile(inputFileName, []byte(codeData.Input), 0644)
	if err != nil {
		fmt.Println("Failed to InitializeFile", fileName)
		return
	}
	return
}
