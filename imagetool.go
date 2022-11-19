// Author: Muhamad Surya Iksanudin<surya.iksanudin@gmail.com>
package imagetool

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const THRESHOLD = 0.1

type (
	blurryDetection struct {
		IsBlur bool
		Score  float64
	}

	blackWhiteDetection struct {
		IsBlackWhite bool
		Score        float64
	}
)

func NewBlackWhiteDetection() blackWhiteDetection {
	return blackWhiteDetection{IsBlackWhite: false, Score: 100.00}
}

func NewBlurryDetection() blurryDetection {
	return blurryDetection{IsBlur: false, Score: -1000.00}
}

func IsBlackWhite(imagePath string) (blackWhiteDetection, error) {
	return NewBlackWhiteDetection().Detect(imagePath)
}

func IsBlur(imagePath string) (blurryDetection, error) {
	return NewBlurryDetection().Detect(imagePath)
}

func (b blurryDetection) Detect(ImagePath string) (blurryDetection, error) {
	if !b.commandAvailable() {
		return b, errors.New("imageMagick is not available")
	}

	if !b.validate(ImagePath) {
		return b, errors.New("please remove space in your image")
	}

	output, err := b.run(ImagePath)
	if err != nil {
		return b, err
	}

	result := strings.Split(output, "\n")
	b.Score, _ = strconv.ParseFloat(strings.Trim(strings.Split(result[len(result)-2], "(")[1], ")"), 64)
	if THRESHOLD > b.Score {
		b.IsBlur = true
	}

	return b, nil
}

func (b blurryDetection) commandAvailable() bool {
	cmd := exec.Command("identify", "-version")
	_, err := cmd.Output()

	return err == nil
}

func (b blurryDetection) validate(ImagePath string) bool {
	return len(strings.Split(ImagePath, " ")) == 1
}

func (b blurryDetection) run(ImagePath string) (string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("identify -verbose %s | grep deviation", ImagePath))
	stdout, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(stdout), err
}

func (b blackWhiteDetection) Detect(ImagePath string) (blackWhiteDetection, error) {
	if !b.commandAvailable() {
		return b, errors.New("imageMagick is not available")
	}

	if !b.validate(ImagePath) {
		return b, errors.New("please remove space in your image")
	}

	output, err := b.run(ImagePath)
	if err != nil {
		return b, err
	}

	result := strings.Split(output, "\n")

	b.Score, _ = strconv.ParseFloat(strings.Trim(strings.Split(result[len(result)-2], " ")[1], ""), 64)
	if THRESHOLD > b.Score {
		b.IsBlackWhite = true
	}

	return b, nil
}

func (b blackWhiteDetection) commandAvailable() bool {
	cmd := exec.Command("convert", "-version")
	_, err := cmd.Output()

	return err == nil
}

func (b blackWhiteDetection) validate(ImagePath string) bool {
	return len(strings.Split(ImagePath, " ")) == 1
}

func (b blackWhiteDetection) run(ImagePath string) (string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("convert %s -colorspace HCL -channel g -separate +channel -format \"%%M: %%[fx:mean]\n\" info: ", ImagePath))
	stdout, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(stdout), err
}
