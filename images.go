package procman

import (
	"fmt"
	"path/filepath"

	"github.com/rutu-sh/procman/internal/common"
	"github.com/rutu-sh/procman/internal/image"
)

func BuildImage(
	name string,
	tag string,
	context_dir string) (*Image, *ImageError) {

	_logger := common.GetLogger()

	imgfind, imgfind_err := image.GetImage("", name, tag)
	if imgfind_err != nil {
		return nil, nil
	}
	if imgfind != nil {
		img := &Image{
			Id:      imgfind.Id,
			Name:    imgfind.Name,
			ImgPath: imgfind.ImgPath,
			Tag:     imgfind.Tag,
			Created: imgfind.Created,
		}
		return img, &ImageError{Message: "image already exists"}
	}

	_logger.Info().Msgf("building image %v:%v using context dir %v", name, tag, context_dir)

	abs_context_dir, err := filepath.Abs(context_dir)
	if err != nil {
		_logger.Error().Msgf("error getting abs path: %v", err)
		return nil, &ImageError{
			Message: fmt.Sprintf("error getting abs path: %v", err),
		}
	}

	res, errbuild := image.BuildImage(name, tag, abs_context_dir)
	if errbuild != nil {
		_logger.Error().Msgf("error building image: %v", errbuild)
		return nil, &ImageError{
			Message: fmt.Sprintf("error building: %v", errbuild),
		}
	}
	if res == nil {
		return nil, nil
	}

	img := &Image{
		Id:      res.Id,
		Name:    res.Name,
		ImgPath: res.ImgPath,
		Tag:     res.Tag,
		Created: res.Created,
	}

	return img, nil
}

func ListImages() (*[]*Image, *ImageListError) {
	_logger := common.GetLogger()

	res, err := image.ListImages()
	if err != nil {
		_logger.Error().Msgf("error listing images: %v", err)
		return &[]*Image{}, nil
	}
	if res == nil {
		return &[]*Image{}, nil
	}

	images := []*Image{}
	for _, img := range *res {
		transformed := &Image{
			Id:      img.Id,
			Name:    img.Name,
			ImgPath: img.ImgPath,
			Tag:     img.Tag,
			Created: img.Created,
		}
		images = append(images, transformed)
	}

	return &images, nil
}

func GetImage(image_id string, name string, tag string) (*Image, *ImageError) {
	_logger := common.GetLogger()

	res, err := image.GetImage(image_id, name, tag)
	if err != nil {
		_logger.Error().Msgf("error getting image (%v, %v, %v): %v", image_id, name, tag, err)
		return nil, &ImageError{Message: err.Message}
	}
	if res == nil {
		return nil, &ImageError{Message: "not found"}
	}

	img := &Image{
		Id:      res.Id,
		Name:    res.Name,
		ImgPath: res.ImgPath,
		Tag:     res.Tag,
		Created: res.Created,
	}
	return img, nil
}

func DelImage(image_id string, name string, tag string) *ImageError {
	_logger := common.GetLogger()
	_logger.Info().Msgf("deletimg image (%v, %v, %v)", image_id, name, tag)

	res, err := image.GetImage(image_id, name, tag)
	if err != nil {
		_logger.Error().Msgf("error getting image (%v, %v, %v): %v", image_id, name, tag, err)
		return &ImageError{Message: err.Message}
	}
	if res == nil {
		return &ImageError{Message: "not found"}
	}

	delErr := image.DelImage(res.Id)
	if delErr != nil {
		return &ImageError{Message: delErr.Message}
	}
	return nil
}
