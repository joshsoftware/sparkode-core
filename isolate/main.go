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

func Run(ctx context.Context, code string, input string, langSpecs LanguageDetails, codeData model.ExecuteCodeRequest) {

	boxID := GenerateRandomID()
	defer func() {
		fmt.Println("running cleanup")
		Cleanup(ctx, boxID)
	}()

	_, err := exec.Command(isolateCommand, "--init", "--cg", "-b", boxID).Output()
	if err != nil {
		logger.Debug(ctx, "Isolate : Failed Run", err.Error())
	}

	boxDirPath := fmt.Sprintf("%s/%s", workdir, boxID)
	// initialize files
	InitializeFile(boxDirPath)

	CreateSourceFilesForInterpreted(langSpecs, codeData, boxDirPath)

	runCfg := config.RunConfig{
		TimeLimit:   5.0,
		WallLimit:   10.2,
		MemoryLimit: 128000,
	}

	command := createCMD(runCfg, filepath.Join(workdir, boxID), boxID)

	err, out, errout := Shellout(command)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	fmt.Println("--- stdout ---")
	fmt.Println(out)
	fmt.Println("--- stderr ---")
	fmt.Println(errout)
	fmt.Println("created box", boxID)
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
func createCMD(cfg config.RunConfig, workdir, boxId string) string {

	cmd := fmt.Sprintf(
		"isolate --cg -s -b %s -M %s/metadata.txt -t %.1f -x 1.0 -w %.1f -k 64000 -p60 --cg-timing --cg-mem=%d -f 1024 -E HOME=/tmp -E PATH=\"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"  -d /etc:noexec --run -- /bin/bash %s < %s/stdin.txt > %s/stdout.txt 2> %s/stderr.txt", boxId, workdir, cfg.TimeLimit, cfg.WallLimit, cfg.MemoryLimit, sourceFile, workdir, workdir, workdir)

	return cmd
}

//First Param: languageSpecs, 2nd argument: boxpath
func CreateSourceFilesForInterpreted(langSpecs LanguageDetails, codeData model.ExecuteCodeRequest, boxPath string) {

	//Create code script
	fileName := filepath.Join(boxPath, "box", filepath.Base(langSpecs.SourceFile))
	_, err := exec.Command("touch", fileName).Output()
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
}
