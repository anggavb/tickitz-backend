package pkg

import "fmt"

func GenerateTicketEmail(
	customerName string,
	movie string,
	date string,
	time string,
	seats string,
	total string,
	orderID string,
) string {

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Tickitz Ticket</title>
</head>

<body style="margin:0;padding:0;background:#f4f6f9;font-family:Arial,sans-serif;">

<table width="100%%" cellpadding="0" cellspacing="0">
<tr>
<td align="center" style="padding:40px 20px;">

<table width="650" cellpadding="0" cellspacing="0"
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
			<h1 style="margin:0;font-size:32px;">
				🎬 Tickitz
			</h1>

			<p style="margin-top:10px;opacity:.9;">
				Your Movie Ticket Platform
			</p>
		</td>
	</tr>

	<tr>
		<td style="padding:40px;">
			<h2 style="margin-top:0;color:#111827;">
				Payment Successful
			</h2>

			<p style="color:#6b7280;line-height:1.8;">
				Hello <strong>%s</strong>,
				<br><br>
				Thank you for your purchase.
				Your movie ticket has been successfully issued.
			</p>

			<table
			width="100%%"
			cellpadding="0"
			cellspacing="0"
			style="
				margin-top:25px;
				border:1px solid #e5e7eb;
				border-radius:12px;
				overflow:hidden;
			">

				<tr>
					<td style="padding:14px;color:#6b7280;width:35%%;">
						Movie
					</td>

					<td style="padding:14px;font-weight:600;">
						%s
					</td>
				</tr>

				<tr style="background:#fafafa;">
					<td style="padding:14px;color:#6b7280;">
						Date
					</td>

					<td style="padding:14px;">
						%s
					</td>
				</tr>

				<tr>
					<td style="padding:14px;color:#6b7280;">
						Time
					</td>

					<td style="padding:14px;">
						%s
					</td>
				</tr>

				<tr style="background:#fafafa;">
					<td style="padding:14px;color:#6b7280;">
						Seats
					</td>

					<td style="padding:14px;">
						%s
					</td>
				</tr>

				<tr>
					<td style="padding:14px;color:#6b7280;">
						Order ID
					</td>

					<td style="padding:14px;font-weight:600;">
						%s
					</td>
				</tr>

				<tr style="background:#fff7ed;">
					<td style="padding:14px;color:#6b7280;">
						Total Payment
					</td>

					<td
					style="
						padding:14px;
						font-size:18px;
						font-weight:bold;
						color:#f97316;
					">
						%s
					</td>
				</tr>

			</table>

			<div style="text-align:center;margin-top:35px;">
				<a
				href="https://tickitz.viketin.my.id/orders"
				style="
					background:#f97316;
					color:white;
					text-decoration:none;
					padding:14px 24px;
					border-radius:10px;
					display:inline-block;
					font-weight:600;
				">
					View My Ticket
				</a>
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
`,
		customerName,
		movie,
		date,
		time,
		seats,
		orderID,
		total,
	)
}
