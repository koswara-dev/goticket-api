package utils

import (
	"fmt"
	"net/smtp"
	"time"
)

type SMTPConfig struct {
	SMTPHost     string
	SMTPPort     string
	SenderEmail  string
	AuthPassword string
}

// SendOTPEmail mengirimkan email HTML berdesain profesional & responsif untuk verifikasi OTP.
func SendOTPEmail(toEmail, otpCode, otpType string, cfg SMTPConfig) error {
	smtpHost := cfg.SMTPHost
	smtpPort := cfg.SMTPPort
	senderEmail := cfg.SenderEmail
	authPassword := cfg.AuthPassword

	if smtpHost == "" || smtpPort == "" {
		if Log != nil {
			Log.WithFields(map[string]interface{}{"email": toEmail, "type": otpType, "otp": otpCode}).Info("[EMAIL-DEBUG] SMTP belum dikonfigurasi, OTP dicetak ke log")
		} else {
			fmt.Printf("[EMAIL-DEBUG] SMTP tidak dikonfigurasi. OTP untuk %s (%s): %s\n", toEmail, otpType, otpCode)
		}
		return nil
	}

	if senderEmail == "" {
		senderEmail = "noreply@goticket.com"
	}

	var subject, titleText, badgeText, descriptionText string
	var headerGradient string

	if otpType == "register" {
		subject = "🎟️ Kode OTP Verifikasi Akun GoTicket Anda"
		titleText = "Verifikasi Pendaftaran Akun"
		badgeText = "REGISTRASI AKUN BARU"
		descriptionText = "Terima kasih telah mendaftar di <strong>GoTicket</strong>! Gunakan kode OTP di bawah ini untuk memverifikasi dan mengaktifkan akun Anda:"
		headerGradient = "linear-gradient(135deg, #4F46E5 0%, #2563EB 100%)"
	} else {
		subject = "🔒 Kode OTP Reset Password Akun GoTicket"
		titleText = "Permintaan Reset Password"
		badgeText = "KEAMANAN AKUN"
		descriptionText = "Kami menerima permintaan untuk mereset password akun <strong>GoTicket</strong> Anda. Masukkan kode OTP berikut untuk melanjutkan pembuatan password baru:"
		headerGradient = "linear-gradient(135deg, #4338CA 0%, #1D4ED8 100%)"
	}

	currentYear := time.Now().Year()

	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        body { margin: 0; padding: 0; background-color: #F8FAFC; font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; }
        table { border-collapse: collapse; }
        .wrapper { width: 100%%; table-layout: fixed; background-color: #F8FAFC; padding: 40px 0; }
        .main-container { background-color: #FFFFFF; max-width: 560px; margin: 0 auto; border-radius: 16px; overflow: hidden; box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.05), 0 8px 10px -6px rgba(0, 0, 0, 0.01); }
        .header { background: %s; padding: 36px 30px; text-align: center; }
        .brand-logo { font-size: 26px; font-weight: 800; color: #FFFFFF; letter-spacing: 1.5px; text-transform: uppercase; text-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .brand-sub { font-size: 12px; color: #E0E7FF; font-weight: 500; margin-top: 4px; letter-spacing: 0.5px; }
        .content { padding: 36px 32px; }
        .badge { display: inline-block; background-color: #EEF2FF; color: #4F46E5; font-size: 11px; font-weight: 700; padding: 6px 14px; border-radius: 50px; letter-spacing: 1px; margin-bottom: 16px; text-transform: uppercase; }
        .title { font-size: 22px; font-weight: 800; color: #0F172A; margin: 0 0 12px 0; line-height: 1.3; }
        .text { font-size: 14px; color: #475569; line-height: 1.6; margin: 0 0 28px 0; }
        .otp-container { background: #F1F5F9; border: 2px dashed #CBD5E1; border-radius: 12px; padding: 24px 20px; text-align: center; margin-bottom: 28px; }
        .otp-code { font-family: 'Courier New', Courier, monospace; font-size: 38px; font-weight: 800; color: #1E293B; letter-spacing: 10px; margin: 0; padding-left: 10px; }
        .exp-badge { display: inline-block; font-size: 12px; color: #D97706; font-weight: 600; background-color: #FEF3C7; padding: 4px 12px; border-radius: 20px; margin-top: 12px; }
        .warning-box { background-color: #FFFBEB; border-left: 4px solid #F59E0B; border-radius: 0 8px 8px 0; padding: 14px 16px; font-size: 12px; color: #92400E; line-height: 1.5; margin-bottom: 24px; }
        .divider { border-top: 1px solid #E2E8F0; margin: 28px 0; }
        .footer { background-color: #0F172A; padding: 24px 30px; text-align: center; color: #94A3B8; font-size: 12px; line-height: 1.5; }
        .footer a { color: #818CF8; text-decoration: none; }
    </style>
</head>
<body>
    <div class="wrapper">
        <!--[if mso]>
        <table align="center" width="560" style="border-spacing:0;"><tr><td style="padding:0;">
        <![endif]-->
        <div class="main-container">
            <div class="header">
                <div class="brand-logo">🎟️ GOTICKET</div>
                <div class="brand-sub">Sistem Pemesanan Tiket Konser Resmi</div>
            </div>
            
            <div class="content">
                <div style="text-align: center;">
                    <span class="badge">%s</span>
                </div>
                <h1 class="title" style="text-align: center;">%s</h1>
                <p class="text">%s</p>
                
                <div class="otp-container">
                    <div class="otp-code">%s</div>
                    <div><span class="exp-badge">⏱️ KODE BERLAKU SELAMA 10 MENIT</span></div>
                </div>
                
                <div class="warning-box">
                    <strong>🔒 Himbauan Keamanan:</strong><br>
                    Jangan pernah memberitahukan kode OTP ini kepada siapapun, termasuk staf yang mengatasnamakan GoTicket. Tim kami tidak pernah meminta kode rahasia ini.
                </div>
                
                <p class="text" style="font-size: 13px; color: #64748B; margin-bottom: 0;">
                    Jika Anda tidak pernah meminta kode OTP ini, abaikan pesan email ini atau segera hubungi layanan bantuan kami.
                </p>
            </div>
            
            <div class="footer">
                <p style="margin: 0 0 6px 0; font-weight: 600; color: #F8FAFC;">GoTicket Ticketing System &copy; %d</p>
                <p style="margin: 0; color: #64748B;">Email ini dikirim secara otomatis oleh sistem. Harap tidak membalas email ini secara langsung.</p>
            </div>
        </div>
        <!--[if mso]>
        </td></tr></table>
        <![endif]-->
    </div>
</body>
</html>`, subject, headerGradient, badgeText, titleText, descriptionText, otpCode, currentYear)

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("GoTicket Support <%s>", senderEmail)
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	var auth smtp.Auth
	if authPassword != "" {
		auth = smtp.PlainAuth("", senderEmail, authPassword, smtpHost)
	}

	err := smtp.SendMail(addr, auth, senderEmail, []string{toEmail}, []byte(message))
	if err != nil {
		if Log != nil {
			Log.WithFields(map[string]interface{}{"error": err, "email": toEmail, "otp": otpCode}).Warn("[EMAIL-DEBUG] Gagal mengirim email SMTP, OTP fallback ke log")
		} else {
			fmt.Printf("[EMAIL-DEBUG] Gagal kirim SMTP (%v). OTP untuk %s (%s): %s\n", err, toEmail, otpType, otpCode)
		}
		return nil
	}

	if Log != nil {
		Log.WithFields(map[string]interface{}{"email": toEmail, "type": otpType}).Info("Email HTML OTP berhasil dikirim via SMTP")
	}

	return nil
}
