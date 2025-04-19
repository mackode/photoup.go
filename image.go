package main

import (
	"path/filepath"
	"strings"
	"github.com/disintegration/imaging"
)

func thumbName(fileName string) string {
	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)
	return name + "_s" + ext
}

func scaleJPG(inputPath string) error {
	outputPath := thumbName(inputPath)
	img, err := imaging.Open(inputPath, imaging.AutoOrientation(true))
	if err != nil {
		return err
	}

	width := img.Bounds().Dx() / 4
	height := img.Bounds().Dy() / 4
	thumbnail := imaging.Thumbnail(img, width, height, imaging.Lanczos)
	err = imaging.Save(thumbnail, outputPath)
	return err
}