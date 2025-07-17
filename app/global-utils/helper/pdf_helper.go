package helper

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"

	"appsku-golang/app/global-utils/model"

	"html/template"

	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// Create a PDF from an HTML file and a data struct
// WARNING: MUST USE DOCKER IMAGE 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com/alpine:wkhtml
func CreatePDF(fileName string, filePath string, data interface{}) (*bytes.Buffer, *model.ErrorLog) {

	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		errorLog := WriteLog(err, http.StatusInternalServerError, nil)
		return nil, errorLog
	}

	newFilePath := strings.Replace(filePath, ".html", "-temp.html", 1)

	f, err := os.Create(newFilePath)
	if err != nil {
		errorLog := WriteLog(err, http.StatusInternalServerError, nil)
		return nil, errorLog
	}

	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		errorLog := WriteLog(err, http.StatusInternalServerError, nil)
		return nil, errorLog
	}

	html, err := os.ReadFile(newFilePath)
	if err != nil {
		errorLog := WriteLog(err, http.StatusInternalServerError, nil)
		return nil, errorLog
	}

	// Create new PDF generator
	pdfg, err := pdf.NewPDFGenerator()
	if err != nil {
		errorLog := WriteLog(err, http.StatusInternalServerError, nil)
		return nil, errorLog
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(pdf.OrientationPortrait)
	pdfg.Grayscale.Set(false)

	page := pdf.NewPageReader(bytes.NewReader(html))
	page.EnableLocalFileAccess.Set(true)

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		errorLog := WriteLog(err, http.StatusInternalServerError, nil)
		return nil, errorLog
	}

	// Write buffer contents to file on disk
	pdfOutput := fmt.Sprintf("%s.pdf", fileName)
	err = pdfg.WriteFile(pdfOutput)
	if err != nil {
		errorLog := WriteLog(err, http.StatusInternalServerError, nil)
		return nil, errorLog
	}

	// Read the generated PDF file
	pdfBytes, err := os.ReadFile(pdfOutput)
	if err != nil {
		errorLog := WriteLog(err, http.StatusInternalServerError, nil)
		return nil, errorLog
	}

	// Delete temporary HTML file
	err = os.Remove(newFilePath)
	if err != nil {
		errorLog := WriteLog(err, http.StatusInternalServerError, nil)
		return nil, errorLog
	}

	// Delete generated PDF file
	err = os.Remove(pdfOutput)
	if err != nil {
		errorLog := WriteLog(err, http.StatusInternalServerError, nil)
		return nil, errorLog
	}

	// Convert PDF bytes to buffer
	buffer := bytes.NewBuffer(pdfBytes)

	return buffer, nil
}
