package pkg

import "fmt"

func GenerateForgotPasswordEmail(resetURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Reset Password</title>
</head>

<body style="
	margin:0;
	padding:0;
	background:#f4f6f9;
	font-family:Arial,sans-serif;
">

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
      background:#f3f4f6;
      padding:40px 35px;
      text-align:center;
    "
  >

    <img
      src="https://res.cloudinary.com/dmsxtj60h/image/upload/v1781625330/logo_t2eboa.png"
      alt="Tickitz Logo"
      width="140"
      style="
        display:block;
        margin:0 auto 20px auto;
        border:0;
      "
    >

    <p
      style="
        margin:0;
        color:#6b7280;
        font-size:16px;
        font-weight:500;
      "
    >
      Password Recovery
    </p>

  </td>
</tr>
	<tr>
		<td style="padding:40px;">

			<h2 style="margin-top:0;color:#111827;">
				Reset Your Password
			</h2>

			<p style="
				color:#6b7280;
				line-height:1.8;
			">
				We received a request to reset the password
				for your Tickitz account.
			</p>

			<p style="
				color:#6b7280;
				line-height:1.8;
			">
				Click the button below to create a new password.
			</p>

			<div style="
				text-align:center;
				margin:40px 0;
			">
				<a
				href="%s"
				style="
					background:#f97316;
					color:white;
					text-decoration:none;
					padding:14px 28px;
					border-radius:10px;
					display:inline-block;
					font-weight:600;
					font-size:16px;
				">
					Reset Password
				</a>
			</div>

			<div style="
				background:#fff7ed;
				padding:16px;
				border-radius:12px;
				color:#9a3412;
				font-size:14px;
				line-height:1.7;
			">
				This reset link will expire in <strong>30 minutes</strong>.
			</div>

			<p style="
				margin-top:25px;
				color:#6b7280;
				line-height:1.8;
			">
				If you didn't request a password reset,
				you can safely ignore this email.
			</p>

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
`, resetURL)
}
