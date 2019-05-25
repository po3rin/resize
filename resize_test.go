package resize_test

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/po3rin/resize"
)

var genGoldenFiles = flag.Bool("gen_golden_files", false, "whether to generate the TestXxx golden files.")

func TestResize(t *testing.T) {
	tests := []struct {
		name           string
		imgFilename    string
		xRatio         float64
		yRatio         float64
		goldenFilename string
	}{
		{
			name:           "x2",
			imgFilename:    "testdata/gopher.jpeg",
			xRatio:         2,
			yRatio:         2,
			goldenFilename: "testdata/resize_golden_0.jpg",
		},
		{
			name:           "x0.5",
			imgFilename:    "testdata/gopher.jpeg",
			xRatio:         0.5,
			yRatio:         0.5,
			goldenFilename: "testdata/resize_golden_1.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.imgFilename)
			if err != nil {
				t.Fatalf("failed to open file\nerr: %v", err)
			}
			defer f.Close()
			img, _, err := image.Decode(f)
			if err != nil {
				t.Fatalf("failed to decode file\nerr: %v", err)
			}
			resized := resize.Resize(img, tt.xRatio, tt.yRatio)

			if *genGoldenFiles {
				goldenFile, err := os.Create(tt.goldenFilename)
				if err != nil {
					t.Errorf("failed to create file\nerr: %v", err)
				}
				defer goldenFile.Close()
				// err = jpeg.Encode(goldenFile, resized, &jpeg.Options{Quality: 100})
				err = jpeg.Encode(goldenFile, resized, nil)
				if err != nil {
					t.Errorf("failed to encode file\nerr: %v", err)
				}
				return
			}

			f, err = os.Create("test.jpg")
			if err != nil {
				t.Fatalf("failed to create file\nerr: %v", err)
			}
			defer f.Close()
			err = jpeg.Encode(f, resized, nil)
			if err != nil {
				t.Fatalf("failed to encode file\nerr: %v", err)
			}

			// got
			f, err = os.Open("test.jpg")
			if err != nil {
				t.Fatalf("failed to open file\nerr: %v", err)
			}
			defer f.Close()
			got, _, err := image.Decode(f)
			if err != nil {
				t.Fatalf("failed to decode file\nerr: %v", err)
			}

			// want
			f, err = os.Open(tt.goldenFilename)
			if err != nil {
				t.Fatalf("failed to open file\nerr: %v", err)
			}
			defer f.Close()
			want, _, err := image.Decode(f)
			if err != nil {
				t.Fatalf("failed to decode file\nerr: %v", err)
			}

			// compare RGBA.
			if !reflect.DeepEqual(convertRGBA(got), convertRGBA(want)) {
				t.Errorf("actual image differs from golden image")
			}

			// b := want.Bounds()
			// for y := b.Min.Y; y < b.Max.Y; y++ {
			// 	for x := b.Min.X; x < b.Max.X; x++ {
			// 		if !eq(convertRGBA(got).At(x, y), convertRGBA(want).At(x, y)) {
			// 			t.Fatalf("%v:%v", convertRGBA(got).At(x, y), convertRGBA(want).At(x, y))
			// 		}
			// 	}
			// }

			// remove test file
			cmd := exec.Command("rm", "test.jpg")
			if err := cmd.Run(); err != nil {
				t.Errorf("failed to rm test file\nerr: %v", err)
			}
		})
	}
}

func convertRGBA(raw image.Image) *image.RGBA {
	want, ok := raw.(*image.RGBA)
	if !ok {
		b := raw.Bounds()
		want = image.NewRGBA(b)
		draw.Draw(want, b, raw, b.Min, draw.Src)
	}
	return want
}

func eq(c0, c1 color.Color) bool {
	r0, g0, b0, a0 := c0.RGBA()
	r1, g1, b1, a1 := c1.RGBA()
	return r0 == r1 && g0 == g1 && b0 == b1 && a0 == a1
}
