#!/usr/bin/env python3
"""
Email Setup - Configure SMTP settings for report delivery
Interactive setup script for email configuration
"""

import json
import getpass
from pathlib import Path

def setup_email_config():
    """Interactive email configuration setup"""
    
    print("=== Email Configuration Setup ===")
    print("This will configure SMTP settings for sending resume reports.")
    print()
    
    config = {}
    
    # SMTP Server
    config['smtp_server'] = input("SMTP Server (e.g., smtp.gmail.com): ").strip()
    
    # SMTP Port
    while True:
        try:
            port = input("SMTP Port (e.g., 587 for TLS, 465 for SSL): ").strip()
            config['smtp_port'] = int(port)
            break
        except ValueError:
            print("Please enter a valid port number.")
    
    # Sender Email
    config['sender_email'] = input("Sender Email Address: ").strip()
    
    # Sender Password
    print("\nNote: For Gmail, you may need to use an App Password instead of your regular password.")
    print("Visit: https://myaccount.google.com/apppasswords")
    config['sender_password'] = getpass.getpass("Sender Password/App Password: ")
    
    # Test configuration option
    test_email = input("Test email address (optional): ").strip()
    
    # Save configuration
    config_path = Path("references/email_config.json")
    config_path.parent.mkdir(exist_ok=True)
    
    try:
        with open(config_path, 'w') as f:
            json.dump(config, f, indent=2)
        
        print(f"\nConfiguration saved to {config_path}")
        
        if test_email:
            print(f"Would you like to test the configuration by sending a test email to {test_email}?")
            test_choice = input("Send test email? (y/n): ").strip().lower()
            if test_choice == 'y':
                test_email_config(config, test_email)
    
    except Exception as e:
        print(f"Error saving configuration: {e}")

def test_email_config(config, test_email):
    """Test email configuration"""
    import smtplib
    from email.mime.text import MIMEText
    from email.mime.multipart import MIMEMultipart
    
    try:
        msg = MIMEMultipart()
        msg['Subject'] = "Resume Scanner - Test Email"
        msg['From'] = config['sender_email']
        msg['To'] = test_email
        
        body = """
        This is a test email from the Resume Scanner system.
        
        If you receive this email, your SMTP configuration is working correctly.
        
        Best regards,
        Resume Scanner
        """
        
        msg.attach(MIMEText(body, 'plain'))
        
        print("Sending test email...")
        with smtplib.SMTP(config['smtp_server'], config['smtp_port']) as server:
            server.starttls()
            server.login(config['sender_email'], config['sender_password'])
            server.send_message(msg)
        
        print("✓ Test email sent successfully!")
        
    except Exception as e:
        print(f"✗ Test failed: {e}")
        print("Please check your SMTP settings and try again.")

if __name__ == "__main__":
    setup_email_config()