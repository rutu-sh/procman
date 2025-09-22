package process

import (
	"fmt"

	"github.com/rutu-sh/procman/internal/common"
)

func StartProcess(proc ProcessCreate) (*Process, *common.ProcStartErr) {
	_logger := common.GetLogger()
	_logger.Info().Msgf("starting process with name: %v", proc.Name)

	process, err := buildProcessContext(proc.Name, proc.Image.Name, proc.Image.Tag)
	if err != nil {
		_logger.Error().Msgf("error starting process: %v", err)
		return nil, &common.ProcStartErr{Code: 500, Message: fmt.Sprintf("error starting process: %v", err)}
	}
	process.Env = getProcEnv(&proc)

	// read process job
	job, errJobParse := parseProcJob(process)
	if errJobParse != nil {
		_logger.Error().Msgf("error parsing job: %v", errJobParse)
		return nil, &common.ProcStartErr{Code: 500, Message: fmt.Sprintf("error starting process: %v", errJobParse)}
	}
	process.Job = *job

	// test code
	process.Network = ProcessNetwork{
		Ports: []PortMapping{
			{HostPort: 8020, ProcPort: 3000},
			{HostPort: 8000, ProcPort: 2000},
			{HostPort: 8080, ProcPort: 4000},
		},
	}

	procConfJson := getProcConfPath(process.Id)
	if errProcWrite := WriteProcessToJson(*process, procConfJson); errProcWrite != nil {
		return nil, &common.ProcStartErr{Code: 500, Message: fmt.Sprintf("error starting process: %v", errProcWrite)}
	}

	return process, nil
}
