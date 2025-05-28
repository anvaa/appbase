package global

import (
	"fmt"
	"image/color"
	"image/png"
	
	"srv/filefunc"
	"srv/srv_conf"
	"app/app_models"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
	"github.com/skip2/go-qrcode"
)

func MakeQRCode(qr_conf app_models.PrintQR) error {
	// Generate QR code for an item
	imgpath := fmt.Sprintf("%s/%d.png", srv_conf.QRImgDir, qr_conf.UUID)
	//fmt.Println("QR code path", imgpath)
	if filefunc.IsExists(imgpath) {
		if err := filefunc.DeleteFile(imgpath); err != nil {
			fmt.Println(err)
		}
	}

	// mnipulate the date and lines for the QR code
	line1 := qr_conf.Date
	line2 := qr_conf.Line1
	line3 := qr_conf.Line2
	line4 := qr_conf.Line3

	// Create the QR code data
	qrtxt := fmt.Sprintf("RAADIG %s\n %s %s %s %s", qr_conf.UUID, line1, line2, line3, line4)
	qrcx, err := qrcode.New(qrtxt, qrcode.Medium)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	qrcx.DisableBorder = true
	qrcx.ForegroundColor = color.Black
	qrcx.BackgroundColor = color.White

	// Write the QR code to a file
	if err := qrcx.WriteFile(256, imgpath); err != nil {
		fmt.Println(err.Error())
	}
	
	return nil
}

func MakeBARCode(bar_conf app_models.PrintBarCode) error {
	// Generate a barcode for an item

	// convert from int to string
	imgpath := fmt.Sprintf("%s/%d.png", srv_conf.BarImgDir, bar_conf.UUID)
	if filefunc.IsExists(imgpath) {
		if err := filefunc.DeleteFile(imgpath); err != nil {
			return err
		}
	}

	// Write the barcode to a file
	file, _ := filefunc.CreateFile(imgpath)
	defer file.Close()

	line1 := bar_conf.Line1
	bartxt := fmt.Sprintf("RAADIG %s %s", bar_conf.UUID, line1)
	
	writer := oned.NewCode128Writer()
	barCode, err := writer.Encode(bartxt, gozxing.BarcodeFormat_CODE_128, 150, 40, nil)
	if err != nil {
		return err
	}

	return png.Encode(file, barCode)
}