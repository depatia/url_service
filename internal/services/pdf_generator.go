package services

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/sirupsen/logrus"
)

type PdfGeneratorService struct {
	logger *logrus.Logger
}

func NewPdfGeneratorService(logger *logrus.Logger) *PdfGeneratorService {
	return &PdfGeneratorService{
		logger: logger,
	}
}

// GeneratePdf - генерирует pdf из переданного текста и выдает байты (дабы не сохранять pdf)
func (pg *PdfGeneratorService) GeneratePdf(ctx context.Context, text string) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetAutoPageBreak(true, 20) // Автоматический разрыв страницы с отступом 20mm

	pdf.SetFont("Arial", "", 16)

	_, lineHt := pdf.GetFontSize()

	lines := strings.Split(text, "\n") // разбиваем на строки вручную

	for _, line := range lines { // для каждой строки создаем новую строку в пдф
		pdf.Cell(40, lineHt, line)
		pdf.Ln(10)
	}

	buf := bytes.NewBuffer(nil)
	if err := pdf.Output(buf); err != nil {
		pg.logger.Warnln("Failed to generate pdf report: ", err.Error())

		return nil, fmt.Errorf("generate pdf: %v", err)
	}

	return buf.Bytes(), nil
}
