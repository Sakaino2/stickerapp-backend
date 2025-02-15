package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/johnfercher/maroto/v2/pkg/core"

	"github.com/johnfercher/maroto/v2"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/linestyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/skip2/go-qrcode"
	"jorge.com/inventari_app_backend/models"
)

func GetStickersPDF(context *gin.Context) {
	ginFile, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := ginFile.Open()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer file.Close()

	fmt.Println("Stickers Prosessing Start")
	var counter = 0

	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err != nil {
		log.Println(err)
	}

	var stickers []models.Sticker

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			continue
		}

		sticker := models.Sticker{
			CodigoProducto:  row[0],
			CodigoUbicacion: row[1],
			Area:            row[2],
			Sede:            row[3],
		}

		stickers = append(stickers, sticker)

		counter++
	}

	document := generateStickersPDF(stickers)

	fmt.Println("Stickers Processing Finished")
	fmt.Println("Stickers Processed:", counter)

	/* filePath := "pdf/stickers.pdf"
	err = document.Save(filePath)
	if err != nil {
		log.Println(err)
		return
	} */

	/*pdfFile, err := os.Open(filePath)

	if err != nil {
		log.Println(err)
		return
	}

	 defer pdfFile.Close()

	_, err = initializers.CreateNewFile("stickers.pdf", pdfFile)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error,
		})
		return
	} */

	/* pdfFile, err := os.Open(filePath)

	if err != nil {
		log.Println(err)
		return
	}

	fileBytes, err := io.ReadAll(pdfFile)
	if err != nil {
		log.Println(err)
		return
	} */

	context.Header("Content-Type", "application/pdf")
	context.Header("Content-Disposition", "attachment; filename=stickers.pdf")

	context.Data(http.StatusOK, "application/pdf", document.GetBytes())

	/* pdfFile.Close()

	err = os.Remove(filePath)
	if err != nil {
		log.Println(err)
		return
	} */

	/* context.JSON(http.StatusOK, gin.H{
		"Success": "Uploaded: " + res.Name,
	}) */
}

/* func createZipArchive(pdfFilename string, zipFilename string) error {
	zipFile, err := os.Create(zipFilename)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zw := zip.NewWriter(zipFile)
	defer zw.Close()

	file, err := os.Open(pdfFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header := &zip.FileHeader{
		Name:     filepath.Base(pdfFilename),
		Method:   zip.Deflate,
		Modified: info.ModTime(),
	}

	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	return nil
} */

func generateStickersPDF(stickers []models.Sticker) core.Document {
	var chunks [][]models.Sticker

	for i := 0; i < len(stickers); i += 30 {
		end := min(i+30, len(stickers))
		chunks = append(chunks, stickers[i:end])
	}

	var pdfs []core.Document

	for _, chunk := range chunks {
		pdf, err := generatePDF(chunk)
		if err != nil {
			log.Println(err)
			continue
		}
		document, err := pdf.Generate()
		if err != nil {
			log.Println(err)
			continue
		}
		pdfs = append(pdfs, document)
	}

	combinedPdf := pdfs[0]

	for i := 1; i < len(pdfs); i++ {
		err := combinedPdf.Merge(pdfs[i].GetBytes())
		if err != nil {
			log.Println(err)
			continue
		}
	}

	return combinedPdf
}

func generatePDF(stickers []models.Sticker) (core.Maroto, error) {
	gridSum := 25
	configuration := config.NewBuilder().
		WithOrientation(orientation.Horizontal).
		WithMaxGridSize(gridSum).
		WithBottomMargin(0).
		Build()
	maroto := maroto.New(configuration)

	var chunks [][]models.Sticker

	for i := 0; i < len(stickers); i += 5 {
		end := min(i+5, len(stickers))
		chunks = append(chunks, stickers[i:end])
	}

	for _, chunk := range chunks {
		var components []core.Col
		var png []byte
		png, err := qrcode.Encode(rangeIn(1000, 9999), qrcode.Medium, 150)
		image := image.NewFromBytes(png, extension.Png, props.Rect{
			Percent: 40,
			Top:     15,
			Left:    38,
		})
		//err := qrcode.WriteColorFile("123456789", qrcode.Medium, 256, color.White, color.Black, "pdf/qr.png")

		if err != nil {
			log.Println(err)
			continue
		}

		/* image := image.NewFromFile("pdf/qr.png", props.Rect{
			Percent: 100,
			Top:     15,
			Left:    38,
		}) */

		for _, sticker := range chunk {
			component := col.New(5).WithStyle(defaultCellStyle()).Add(
				/* code.NewQr(sticker.CodigoUbicacion, props.Rect{
					Percent: 100,
					Top:     20,
					Left:    10,
					Center:  false,
				}), */
				text.New(sticker.CodigoProducto, props.Text{
					Top:   2,
					Size:  12,
					Left:  3,
					Style: fontstyle.Bold,
					Align: align.Left,
					Color: getDarkGrayColor(),
				}),
				text.New(sticker.CodigoUbicacion, props.Text{
					Top:   7,
					Size:  12,
					Left:  3,
					Style: fontstyle.Bold,
					Align: align.Left,
					Color: getDarkGrayColor(),
				}),
				text.New(sticker.Area, props.Text{
					Top:   13,
					Size:  10,
					Left:  3,
					Style: fontstyle.Normal,
					Align: align.Left,
					Color: getDarkGrayColor(),
				}),
				text.New(sticker.Sede, props.Text{
					Top:   18,
					Size:  10,
					Left:  3,
					Style: fontstyle.Normal,
					Align: align.Left,
					Color: getDarkGrayColor(),
				}),
				image,
			)
			components = append(components, component)
		}

		maroto.AddRow(31, components...)
	}

	return maroto, nil
}

func defaultCellStyle() *props.Cell {
	colStyle := &props.Cell{
		BorderType:      border.Full,
		BorderColor:     getDarkGrayColor(),
		LineStyle:       linestyle.Solid,
		BorderThickness: 0.25,
	}

	return colStyle
}

func getDarkGrayColor() *props.Color {
	return &props.Color{
		Red:   55,
		Green: 55,
		Blue:  55,
	}
}

func rangeIn(low, hi int) string {
	return strconv.Itoa(low + rand.Intn(hi-low))
}
