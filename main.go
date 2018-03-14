package main

import (
	"flag"
	"fmt"

	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract"
)

func main() {

	var sourceImagePath string
	var config string
	var useWhitelist bool
	var sharpen float64

	var cropTop int
	var cropBottom int
	var cropLeft int
	var cropRight int

	var cropAll int

	flag.StringVar(&sourceImagePath, "image", "", "path to image")
	flag.StringVar(&config, "config", "", "path to config")
	flag.BoolVar(&useWhitelist, "use-whitelist", false, "true|false to enable|disable usage of the hard-coded whitelist")
	flag.Float64Var(&sharpen, "sharpen", 0, "sigma value (float) for sharpening")

	flag.IntVar(&cropTop, "crop-top", 0, "")
	flag.IntVar(&cropBottom, "crop-bottom", 0, "")
	flag.IntVar(&cropLeft, "crop-left", 0, "")
	flag.IntVar(&cropRight, "crop-right", 0, "")

	flag.IntVar(&cropAll, "crop-all", 0, "")

	flag.Parse()

	client := gosseract.NewClient()
	defer client.Close()

	preparedImagePath := "./out/prepared.tiff"

	image, err := imaging.Open(sourceImagePath)
	if err != nil {
		panic(err)
	}

	if sharpen > 0 {
		image = imaging.Sharpen(image, sharpen)
	}

	{
		width := image.Bounds().Max.X
		height := image.Bounds().Max.Y

		if cropAll > 0 {
			width = width - (cropAll * 2)
			height = height - (cropAll * 2)

			image = imaging.CropCenter(image, width, height)
		} else {

			/*
				rectangle := image.Rectangle

				if cropTop {
					imaging.Crop
				}*/

		}
	}

	err = imaging.Save(image, preparedImagePath)
	if err != nil {
		panic(err)
	}

	client.SetImage(preparedImagePath)

	if config != "" {
		client.SetConfigFile(config)
	}
	client.SetPageSegMode(gosseract.PSM_SPARSE_TEXT_OSD)

	if useWhitelist {
		client.SetWhitelist("abcdefghijklmnopqrstuvwxyz	ABCDEFGHIJKLMNOPQRSTUVWXYZ 1234567890 _ - ")
	}

	text, _ := client.Text()
	fmt.Println(text)
}
