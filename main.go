package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"sort"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Invalid number of arguments, it should be 2\nThe first one is path to image\nThe second is path to new image\n")
		return
	}
	file, err := os.Open(os.Args[1])
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		return
	}

	//TODO check format
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
		return
	}

	sortedImg := sortVertically(img)

	err = saveImage(os.Args[2], sortedImg)
	if err != nil {
		log.Fatal(err)
		return
	}

}

func saveImage(fileName string, image *image.NRGBA) (err error) {
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	err = png.Encode(file, image)
	return
}

type column struct {
	color []color.Color
}

func newColumn(height int) *column {
	return &column{
		color: make([]color.Color, height),
	}
}

func (c *column) Len() int {
	return len(c.color)
}

func (c *column) Less(i, j int) bool {
	ir, ig, ib, _ := c.color[i].RGBA()
	jr, jg, jb, _ := c.color[j].RGBA()
	return ir+ig+ib < jr+jg+jb
}

func (c *column) Swap(i, j int) {
	tmp := c.color[i]
	c.color[i] = c.color[j]
	c.color[j] = tmp
}

//TODO will have a bug if image don't span from origo
func sortVertically(img image.Image) *image.NRGBA {
	sortedImage := image.NewNRGBA(img.Bounds())
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		column := newColumn(img.Bounds().Max.Y - img.Bounds().Min.Y)
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			column.color[y] = img.At(x, y)
		}
		sort.Sort(column)
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			sortedImage.Set(x, y, column.color[y])
		}
	}
	return sortedImage
}
