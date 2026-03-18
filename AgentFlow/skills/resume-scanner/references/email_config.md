# Email Configuration Guide

This guide explains how to configure email settings for the Resume Scanner system.

## Quick Setup

Run the interactive setup script:

```bash
python3 scripts/setup_email.py
```

## Manual Configuration

Edit `references/email_config.json` manually:

```json
{
  "smtp_server": "smtp.gmail.com",
  "smtp_port": 587,
  "sender_email": "your-email@gmail.com",
  "sender_password": "your-app-password"
}
```

## Common SMTP Providers

### Gmail

- Server: `smtp.gmail.com`
- Port: `587` (TLS) or `465` (SSL)
- Password: Use App Password (not regular password)

### Outlook/Hotmail

- Server: `smtp-mail.outlook.com`
- Port: `587` (TLS)

### Yahoo Mail

- Server: `smtp.mail.yahoo.com`
- Port: `587` (TLS)

## Gmail App Password Setup

1. Go to: https://myaccount.google.com/apppasswords
2. Sign in to your Google Account
3. Select "Mail" for the app
4. Select "Other (Custom name)" and enter "Resume Scanner"
5. Click "Generate"
6. Copy the 16-character password
7. Use this password in the email configuration

## Security Notes

- Never commit email credentials to version control
- Use app passwords instead of main account passwords
- Consider using environment variables for production deployments
- Enable two-factor authentication on your email account

## Troubleshooting

### Authentication Failed

- Check if you're using an app password (for Gmail)
- Verify email address and password are correct
- Ensure two-factor authentication is enabled

### Connection Timeout

- Verify SMTP server and port settings
- Check firewall settings
- Ensure internet connectivity

### Email Not Received

- Check spam/junk folder
- Verify recipient email address
- Check email sending limits for your provider
