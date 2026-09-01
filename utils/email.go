package utils

import (
	"fmt"
	"gotiket-api/config"
	"net/smtp"
)

func SendOTPEmail(toEmail, otpCode, otpType string, cfg config.AppConfig) error {
	smtpHost := cfg.SMTPHost
	smtpPort := cfg.SMTPPort
	senderEmail := cfg.SenderEmail
	authPassword := cfg.AuthPassword

	if smtpHost == "" || smtpPort == "" {
		fmt.Printf("[EMAIL-DEBUG] SMTP tidak dikonfigurasi. OTP untuk %s (%s): %s\n", toEmail, otpType, otpCode)
		return nil
	}

	if senderEmail == "" {
		senderEmail = "noreply@goticket.com"
	}

	var subject, actionName string
	if otpType == "register" {
		subject = "Kode OTP Verifikasi Akun GoTicket"
		actionName = "Verifikasi Pendaftaran Akun"
	} else {
		subject = "Kode OTP Reset Password GoTicket"
		actionName = "Reset Password Akun"
	}

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; background-color: #f4f6f9; margin: 0; padding: 20px; }
        .container { max-width: 500px; margin: 0 auto; background: #ffffff; padding: 30px; border-radius: 8px; box-shadow: 0 4px 10px rgba(0,0,0,0.05); }
        .header { text-align: center; color: #1e293b; font-size: 22px; font-weight: bold; margin-bottom: 20px; }
        .otp-box { text-align: center; background: #f1f5f9; padding: 15px; border-radius: 6px; font-size: 32px; font-weight: bold; letter-spacing: 6px; color: #2563eb; margin: 20px 0; }
        .footer { text-align: center; font-size: 12px; color: #94a3b8; margin-top: 25px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">GoTicket API</div>
        <p>Halo,</p>
        <p>Berikut adalah kode OTP Anda untuk <strong>%s</strong>:</p>
        <div class="otp-box">%s</div>
        <p>Kode OTP ini berlaku selama <strong>10 menit</strong>. Jangan bagikan kode ini kepada siapa pun demi keamanan akun Anda.</p>
        <div class="footer">&copy; GoTicket Ticketing System</div>
    </div>
</body>
</html>`, actionName, otpCode)

	headers := make(map[string]string)
	headers["From"] = senderEmail
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
		fmt.Printf("[EMAIL-DEBUG] Gagal kirim SMTP (%v). OTP untuk %s (%s): %s\n", err, toEmail, otpType, otpCode)
		return nil
	}

	return nil
}
