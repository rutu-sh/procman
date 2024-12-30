package process

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/rutu-sh/procman/internal/common"
	"github.com/rutu-sh/procman/internal/image"
)

func run(command []string) *common.ProcStartErr {
	err := runCmd([]string{}, command[0], command[1:]...)
	if err != nil {
		return &common.ProcStartErr{Code: 500, Message: fmt.Sprintf("error running command image: %v", err)}
	}
	return nil
}

func BuildProcessContext(name string, image_id string, image_name string, image_tag string) (*Process, *common.ProcStartErr) {
	_logger := common.GetLogger()

	_logger.Info().Msgf("starting process with params (%v, %v, %v, %v)", name, image_id, image_name, image_tag)

	img, err := image.GetImage(image_id, image_name, image_tag)
	if img == nil || err != nil {
		_logger.Error().Msgf("error reading image: %v", err)
		return nil, &common.ProcStartErr{Code: 500, Message: fmt.Sprintf("error reading image: %v", err)}
	}

	uid := strings.Split(uuid.New().String(), "-")[0]
	procDir := getProcessDir(uid)

	proc := &Process{
		Id:         uid,
		Image:      *img,
		ContextDir: procDir,
	}

	run([]string{"cp", fmt.Sprintf("%v/img.tar.gz", img.ImgPath), proc.ContextDir})

	wd, _ := os.Getwd()
	if errchdir := os.Chdir(proc.ContextDir); errchdir != nil {
		return nil, &common.ProcStartErr{Code: 500, Message: fmt.Sprintf("error changing dir: %v", err)}
	}
	run([]string{"tar", "-xf", "img.tar.gz"})
	run([]string{"rm", "img.tar.gz"})
	if errchdir := os.Chdir(wd); errchdir != nil {
		return nil, &common.ProcStartErr{Code: 500, Message: fmt.Sprintf("error changing dir: %v", err)}
	}

	return proc, nil
}
