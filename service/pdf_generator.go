package service

import (
	"bytes"
	"fmt"
	"gotiket-api/model"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func GenerateBookingPDF(bookings []model.Booking, startDate, endDate string) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 15, 10)
	pdf.AddPage()

	// Header Document
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "GOTICKET API - LAPORAN PEMESANAN TIKET KONSER", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 7, fmt.Sprintf("Periode: %s s.d %s", startDate, endDate), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Table Header
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(12, 8, "No", "1", 0, "C", true, 0, "")
	pdf.CellFormat(45, 8, "Kode Booking", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Tanggal", "1", 0, "C", true, 0, "")
	pdf.CellFormat(55, 8, "Nama Customer", "1", 0, "C", true, 0, "")
	pdf.CellFormat(75, 8, "Detail Tiket (Kategori x Jml)", "1", 0, "C", true, 0, "")
	pdf.CellFormat(50, 8, "Total Harga (IDR)", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	var grandTotal float64

	if len(bookings) == 0 {
		pdf.CellFormat(277, 10, "Tidak ada data transaksi pemesanan pada periode ini.", "1", 1, "C", false, 0, "")
	} else {
		for i, b := range bookings {
			grandTotal += b.TotalAmount

			var itemsStr string
			for idx, d := range b.Details {
				if idx > 0 {
					itemsStr += ", "
				}
				categoryName := d.TicketCategory.Name
				if categoryName == "" {
					categoryName = fmt.Sprintf("Cat #%d", d.TicketCategoryID)
				}
				itemsStr += fmt.Sprintf("%s x%d", categoryName, d.Quantity)
			}

			customerName := b.Customer.Name
			if customerName == "" {
				customerName = fmt.Sprintf("Cust #%d", b.CustomerID)
			}

			dateStr := b.BookingDate.Format("2006-01-02 15:04")

			pdf.CellFormat(12, 8, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
			pdf.CellFormat(45, 8, b.BookingCode, "1", 0, "C", false, 0, "")
			pdf.CellFormat(40, 8, dateStr, "1", 0, "C", false, 0, "")
			pdf.CellFormat(55, 8, customerName, "1", 0, "L", false, 0, "")
			pdf.CellFormat(75, 8, itemsStr, "1", 0, "L", false, 0, "")
			pdf.CellFormat(50, 8, fmt.Sprintf("Rp %s", formatRupiah(b.TotalAmount)), "1", 1, "R", false, 0, "")
		}
	}

	// Grand Total Row
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(227, 9, "TOTAL PENDAPATAN", "1", 0, "R", true, 0, "")
	pdf.CellFormat(50, 9, fmt.Sprintf("Rp %s", formatRupiah(grandTotal)), "1", 1, "R", true, 0, "")

	// Footer / Generated Timestamp
	pdf.Ln(8)
	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(0, 5, fmt.Sprintf("Dicetak otomatis oleh Sistem GoTicket API pada %s", time.Now().Format("2006-01-02 15:04:05 MST")), "", 1, "R", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func formatRupiah(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
