package resize_test

import (
	"flag"
	"image"
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
		imgFilename    string
		xRatio         float64
		yRatio         float64
		goldenFilename string
	}{
		{
			imgFilename:    "testdata/gopher.jpeg",
			xRatio:         2,
			yRatio:         2,
			goldenFilename: "testdata/resize_golden_0.jpg",
		},
		{
			imgFilename:    "testdata/gopher.jpeg",
			xRatio:         0.5,
			yRatio:         0.5,
			goldenFilename: "testdata/resize_golden_1.jpg",
		},
	}

	for _, tt := range tests {
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
			err = jpeg.Encode(goldenFile, resized, nil)
			if err != nil {
				t.Errorf("failed to encode file\nerr: %v", err)
			}
			continue
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
			continue
		}

		// remove test file
		cmd := exec.Command("rm", "test.jpg")
		if err := cmd.Run(); err != nil {
			t.Errorf("failed to rm test file\nerr: %v", err)
		}
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
