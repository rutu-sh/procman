package process

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rutu-sh/procman/internal/common"
	"github.com/rutu-sh/procman/internal/image"
	"gopkg.in/yaml.v3"
	"encoding/json"
)

func getAllProcDir() string {
	dir := "/var/lib/procman/proc"
	os.MkdirAll(dir, 0755)
	return dir
}

func getProcRootFS(process_id string) string {
	dir := fmt.Sprintf("%v/%v/rootfs", getAllProcDir(), process_id)
	os.MkdirAll(dir, 0755)
	return dir
}

func getProcConfDir() string {
	return "etc/procman"
}

func getProcConfPath(process_id string) string {
	return fmt.Sprintf("%v/%v/process.json", getProcRootFS(process_id), getProcConfDir())
}

/*
Do not confuse process_id with process_pid
*/
func getProcessDir(process_id string) string {
	dir := fmt.Sprintf("%v/%v", getAllProcDir(), process_id)
	os.MkdirAll(dir, 0755)
	return dir
}

func runCmd(env []string, command string, args ...string) *common.ProcStartErr {
	_logger := common.GetLogger()

	_logger.Info().Msgf("executing command: %v with args: %v", command, args)

	cmd := exec.Command(command, args...)
	if len(env) > 0 {
		cmd.Env = env
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		_logger.Error().Msgf("error running command: %v", err)
		return &common.ProcStartErr{
			Code:    500,
			Message: "error running the command",
		}
	}

	_logger.Info().Msgf("command %v executed", command)

	return nil
}

func run(command []string) *common.ProcStartErr {
	err := runCmd([]string{}, command[0], command[1:]...)
	if err != nil {
		return &common.ProcStartErr{Code: 500, Message: fmt.Sprintf("error running command image: %v", err)}
	}
	return nil
}

func parseProcJob(process *Process) (*image.ImageJob, *common.ProcStartErr) {
	_logger := common.GetLogger()
	jobYamlFile := fmt.Sprintf("%v/%v/job.yaml", getProcRootFS(process.Id), getProcConfDir())

	_logger.Info().Msgf("parsing job yaml at: %v", jobYamlFile)
	if _, err := os.Stat(jobYamlFile); err != nil {
		_logger.Error().Msgf("file %v not found", jobYamlFile)
		return nil, &common.ProcStartErr{Code: 500, Message: fmt.Sprintf("file %v not found", jobYamlFile)}
	}

	data, errRead := os.ReadFile(jobYamlFile)
	if errRead != nil {
		_logger.Error().Msgf("error unmarshal: %v", errRead)
		return nil, &common.ProcStartErr{
			Code:    500,
			Message: fmt.Sprintf("error reading the yaml spec %v: %v", jobYamlFile, errRead),
		}
	}

	parsedJob := &image.ImageJob{}
	errUnmarshal := yaml.Unmarshal(data, &parsedJob)
	if errUnmarshal != nil {
		_logger.Error().Msgf("error unmarshal: %v", errUnmarshal)
		return nil, &common.ProcStartErr{
			Code:    500,
			Message: fmt.Sprintf("error reading the yaml spec %v: %v", jobYamlFile, errUnmarshal),
		}
	}

	return parsedJob, nil
}

func getJobDefaultEnvs() ProcessEnv {
	m := make(ProcessEnv)
	m["PATH"] = "/bin:/sbin:/usr/bin:/usr/sbin"
	m["HOME"] = "/home"
	m["TERM"] = "xterm-256color"
	m["LANG"] = "en_US.UTF-8"
	m["LANGUAGE"] = "en_US:en"
	m["LC_ALL"] = "en_US.UTF-8"
	m["PS1"] = "[namepace] > "
	return m
}

func getProcEnv(proc *ProcessCreate) ProcessEnv {
	jobEnv := getJobDefaultEnvs()
	for key, val := range proc.Env {
		jobEnv[key] = val
	}
	return jobEnv
}

// Function to write Process object to a JSON file
func WriteProcessToJson(proc Process, filepath string) error {
	data, err := json.Marshal(&proc)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
