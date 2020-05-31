package main

import (
	"fmt"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func main() {
	// extract page 35 from input
	pdfcpuapi.ExtractPagesFile("scenario.pdf", "test", []string{"35"}, nil)
	extractedPage := "test/page_35.pdf"

	// add demo watermark do page
	onTop := true
	wm, _ := pdfcpu.ParseTextWatermarkDetails("Demo", "", onTop)
	extracedPageWithWatermark := "test/page_35_wm.pdf"
	err := pdfcpuapi.AddWatermarksFile(extractedPage, extracedPageWithWatermark, nil, wm, nil)
	if err != nil {
		fmt.Printf("Error while adding watermark: %v", err)
	}
}
