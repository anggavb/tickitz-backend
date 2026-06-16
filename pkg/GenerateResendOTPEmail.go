package pkg

import "fmt"

func GenerateResendOTPEmail(otp string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Tickitz New OTP</title>
</head>

<body style="margin:0;padding:0;background:#f4f6f9;font-family:Arial,sans-serif;">

<table width="100%%" cellpadding="0" cellspacing="0">
<tr>
<td align="center" style="padding:40px 20px;">

<table width="600" cellpadding="0" cellspacing="0"
style="
	background:#ffffff;
	border-radius:20px;
	overflow:hidden;
	box-shadow:0 4px 20px rgba(0,0,0,.08);
">

<tr>
<td
style="
	background:linear-gradient(135deg,#f97316,#ea580c);
	padding:35px;
	text-align:center;
	color:white;
">
<table role="presentation" width="100%%">
<tr>
<td align="center">
  <img
    src="https://res.cloudinary.com/dmsxtj60h/image/upload/v1781625330/logo_t2eboa.png"
    alt="Tickitz Logo"
    width="64"
    style="display:block;border:0;"
  >
</td>
</tr>
</table>
	<h1 style="margin:0;">Tickitz</h1>
	<p style="margin-top:10px;">New Activation Code</p>
</td>
</tr>

<tr>
<td style="padding:40px;">

	<h2 style="margin-top:0;color:#111827;">
		Request New OTP
	</h2>

	<p style="color:#6b7280;line-height:1.8;">
		You requested a new activation code because the previous OTP
		has expired or was not received.
	</p>

	<div style="text-align:center;margin:35px 0;">
		<div style="
			display:inline-block;
			padding:18px 40px;
			background:#fff7ed;
			border:2px dashed #f97316;
			border-radius:14px;
			font-size:32px;
			font-weight:bold;
			letter-spacing:8px;
			color:#f97316;
		">
			%s
		</div>
	</div>

	<p style="color:#6b7280;line-height:1.8;">
		Please use this OTP to activate your Tickitz account.
		This code will expire after a certain period.
	</p>

	<div style="
		background:#fef2f2;
		color:#dc2626;
		padding:14px;
		border-radius:10px;
		margin-top:20px;
		font-size:14px;
	">
		For security reasons, never share this code with anyone.
	</div>

</td>
</tr>

<tr>
<td
style="
	background:#f9fafb;
	padding:20px;
	text-align:center;
	font-size:12px;
	color:#9ca3af;
">
	© 2026 Tickitz. All rights reserved.
</td>
</tr>

</table>

</td>
</tr>
</table>

</body>
</html>
`, otp)
}
