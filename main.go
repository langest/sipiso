package main

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"sort"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Invalid number of arguments, it should be at least 1 and at most two\nThe first one is path to image\nThe second is path to the output image\nIf no output is given it will output to \"<inputfile.ending>.sorted.png\"")
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

	sortedImg := image.NewNRGBA(img.Bounds())
	bounds := img.Bounds()
	xPieces := bounds.Max.X / 2
	startX := bounds.Min.X
	stepX := bounds.Max.X / xPieces
	for x := 0; x < xPieces; x++ {
		currentY := 0
		yVariance := bounds.Max.Y / 10
		nextY := rand.Intn(yVariance)
		for ; nextY < bounds.Max.Y; nextY += 10 + rand.Intn(yVariance) {
			sortVertically(sortedImg, img, startX+x*stepX, startX+(x+1)*stepX, currentY, nextY)
			currentY = nextY
		}
		sortVertically(sortedImg, img, startX+x*stepX, startX+(x+1)*stepX, currentY, bounds.Max.Y)
	}

	output := os.Args[1] + ".sorted.png"
	if len(os.Args) == 3 {
		output = os.Args[2]
	}
	err = saveImage(output, sortedImg)
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
	return ir+ig+ib > jr+jg+jb
}

func (c *column) Swap(i, j int) {
	tmp := c.color[i]
	c.color[i] = c.color[j]
	c.color[j] = tmp
}

func sortVertically(img *image.NRGBA, orig image.Image, xStart, xEnd, yStart, yEnd int) {
	for x := xStart; x < xEnd; x++ {
		column := newColumn(yEnd - yStart) //Calculate the size of the column
		for y := 0; y < yEnd-yStart; y++ {
			column.color[y] = orig.At(x, yStart+y)
		}
		sort.Sort(column)
		for y := 0; y < yEnd-yStart; y++ {
			img.Set(x, yStart+y, column.color[y])
		}
	}
}
