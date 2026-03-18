#!/usr/bin/env python3
"""
Report Generator - Create visualization-rich reports and send via email
Generates comprehensive resume analysis reports with charts and statistics
"""

import json
import sys
import argparse
import smtplib
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from email.mime.base import MIMEBase
from email import encoders
from pathlib import Path
import matplotlib.pyplot as plt
import pandas as pd
from datetime import datetime

def create_visualizations(results: dict, output_dir: str):
    """Create visualization charts for the report"""
    output_path = Path(output_dir)
    output_path.mkdir(exist_ok=True)
    
    # Data for charts
    resumes = results.get("resumes", [])
    if not resumes:
        return {}
    
    scores = [r["analysis"]["score"] for r in resumes]
    filenames = [r["filename"] for r in resumes]
    
    chart_files = {}
    
    # 1. Score Distribution Bar Chart
    plt.figure(figsize=(12, 6))
    colors = ['red' if score < 0.4 else 'orange' if score < 0.7 else 'green' for score in scores]
    plt.bar(range(len(scores)), scores, color=colors, alpha=0.7)
    plt.xlabel('Resumes')
    plt.ylabel('Match Score')
    plt.title('Resume Match Scores')
    plt.xticks(range(len(filenames)), filenames, rotation=45, ha='right')
    plt.tight_layout()
    
    score_chart = output_path / "score_distribution.png"
    plt.savefig(score_chart, dpi=300, bbox_inches='tight')
    plt.close()
    chart_files["score_distribution"] = str(score_chart)
    
    # 2. Match Categories Pie Chart
    high_match = len([s for s in scores if s >= 0.7])
    medium_match = len([s for s in scores if 0.4 <= s < 0.7])
    low_match = len([s for s in scores if s < 0.4])
    
    plt.figure(figsize=(8, 6))
    labels = ['High Match (≥70%)', 'Medium Match (40-70%)', 'Low Match (<40%)']
    sizes = [high_match, medium_match, low_match]
    colors = ['green', 'orange', 'red']
    
    plt.pie(sizes, labels=labels, colors=colors, autopct='%1.1f%%', startangle=90)
    plt.title('Resume Match Categories')
    plt.axis('equal')
    
    pie_chart = output_path / "match_categories.png"
    plt.savefig(pie_chart, dpi=300, bbox_inches='tight')
    plt.close()
    chart_files["match_categories"] = str(pie_chart)
    
    # 3. Keyword Analysis
    keyword_matches = {}
    for resume in resumes:
        for keyword in resume["analysis"]["matches"].get("keywords", []):
            keyword_matches[keyword] = keyword_matches.get(keyword, 0) + 1
    
    if keyword_matches:
        plt.figure(figsize=(10, 6))
        keywords = list(keyword_matches.keys())
        counts = list(keyword_matches.values())
        
        plt.bar(keywords, counts, color='skyblue', alpha=0.7)
        plt.xlabel('Keywords')
        plt.ylabel('Frequency')
        plt.title('Keyword Match Frequency')
        plt.xticks(rotation=45, ha='right')
        plt.tight_layout()
        
        keyword_chart = output_path / "keyword_analysis.png"
        plt.savefig(keyword_chart, dpi=300, bbox_inches='tight')
        plt.close()
        chart_files["keyword_analysis"] = str(keyword_chart)
    
    return chart_files

def generate_html_report(results: dict, chart_files: dict) -> str:
    """Generate HTML report with embedded charts"""
    
    summary = results.get("summary", {})
    resumes = results.get("resumes", [])
    
    html = f"""
    <!DOCTYPE html>
    <html>
    <head>
        <title>Resume Analysis Report</title>
        <style>
            body {{ font-family: Arial, sans-serif; margin: 20px; }}
            .header {{ background-color: #f0f0f0; padding: 20px; border-radius: 5px; }}
            .section {{ margin: 20px 0; }}
            .chart {{ text-align: center; margin: 20px 0; }}
            .chart img {{ max-width: 100%; height: auto; border: 1px solid #ddd; }}
            .summary-table {{ width: 100%; border-collapse: collapse; margin: 20px 0; }}
            .summary-table th, .summary-table td {{ border: 1px solid #ddd; padding: 8px; text-align: left; }}
            .summary-table th {{ background-color: #f2f2f2; }}
            .resume-item {{ border: 1px solid #ddd; margin: 10px 0; padding: 15px; border-radius: 5px; }}
            .high-match {{ border-left: 5px solid green; }}
            .medium-match {{ border-left: 5px solid orange; }}
            .low-match {{ border-left: 5px solid red; }}
        </style>
    </head>
    <body>
        <div class="header">
            <h1>Resume Analysis Report</h1>
            <p>Generated on: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}</p>
            <p>Total Resumes Scanned: {summary.get('total_resumes', 0)}</p>
            <p>Average Match Score: {summary.get('average_score', 0):.2f}</p>
        </div>
        
        <div class="section">
            <h2>Executive Summary</h2>
            <table class="summary-table">
                <tr>
                    <th>Metric</th>
                    <th>Count</th>
                    <th>Percentage</th>
                </tr>
                <tr>
                    <td>High Match (≥70%)</td>
                    <td>{summary.get('high_match_count', 0)}</td>
                    <td>{(summary.get('high_match_count', 0) / summary.get('total_resumes', 1) * 100):.1f}%</td>
                </tr>
                <tr>
                    <td>Medium Match (40-70%)</td>
                    <td>{summary.get('medium_match_count', 0)}</td>
                    <td>{(summary.get('medium_match_count', 0) / summary.get('total_resumes', 1) * 100):.1f}%</td>
                </tr>
                <tr>
                    <td>Low Match (<40%)</td>
                    <td>{summary.get('low_match_count', 0)}</td>
                    <td>{(summary.get('low_match_count', 0) / summary.get('total_resumes', 1) * 100):.1f}%</td>
                </tr>
            </table>
        </div>
        
        <div class="section">
            <h2>Visualizations</h2>
            <div class="chart">
                <h3>Match Score Distribution</h3>
                {f'<img src="cid:score_distribution" alt="Score Distribution">' if 'score_distribution' in chart_files else '<p>No chart available</p>'}
            </div>
            <div class="chart">
                <h3>Match Categories</h3>
                {f'<img src="cid:match_categories" alt="Match Categories">' if 'match_categories' in chart_files else '<p>No chart available</p>'}
            </div>
            <div class="chart">
                <h3>Keyword Analysis</h3>
                {f'<img src="cid:keyword_analysis" alt="Keyword Analysis">' if 'keyword_analysis' in chart_files else '<p>No chart available</p>'}
            </div>
        </div>
        
        <div class="section">
            <h2>Detailed Resume Analysis</h2>
    """
    
    for resume in resumes:
        score = resume["analysis"]["score"]
        match_class = "high-match" if score >= 0.7 else "medium-match" if score >= 0.4 else "low-match"
        
        html += f"""
            <div class="resume-item {match_class}">
                <h3>{resume['filename']}</h3>
                <p><strong>Match Score:</strong> {score:.2f} ({score*100:.1f}%)</p>
                <p><strong>Matched Keywords:</strong> {', '.join(resume['analysis']['matches'].get('keywords', []))}</p>
                <p><strong>Education:</strong> {', '.join(resume['analysis']['matches'].get('education', []))}</p>
                <p><strong>Experience Years:</strong> {resume['analysis']['matches'].get('experience_years', 0)}</p>
                <p><strong>Industries:</strong> {', '.join(resume['analysis']['matches'].get('industries', []))}</p>
            </div>
        """
    
    html += """
        </div>
    </body>
    </html>
    """
    
    return html

def send_email_report(html_content: str, chart_files: dict, recipient_email: str, config_file: str = "references/email_config.json"):
    """Send HTML report via email with embedded images"""
    
    # Load email configuration
    try:
        with open(config_file, 'r') as f:
            config = json.load(f)
    except Exception as e:
        print(f"Error loading email config: {e}")
        return False
    
    # Create message
    msg = MIMEMultipart('related')
    msg['Subject'] = f"Resume Analysis Report - {datetime.now().strftime('%Y-%m-%d')}"
    msg['From'] = config['sender_email']
    msg['To'] = recipient_email
    
    # Attach HTML content
    html_part = MIMEText(html_content, 'html')
    msg.attach(html_part)
    
    # Attach images
    for chart_name, chart_path in chart_files.items():
        try:
            with open(chart_path, 'rb') as f:
                img_data = f.read()
            
            img = MIMEBase('image', 'png')
            img.set_payload(img_data)
            encoders.encode_base64(img)
            img.add_header('Content-ID', f'<{chart_name}>')
            img.add_header('Content-Disposition', f'inline; filename="{chart_name}.png"')
            msg.attach(img)
        except Exception as e:
            print(f"Error attaching {chart_name}: {e}")
    
    # Send email
    try:
        with smtplib.SMTP(config['smtp_server'], config['smtp_port']) as server:
            server.starttls()
            server.login(config['sender_email'], config['sender_password'])
            server.send_message(msg)
        print(f"Report sent successfully to {recipient_email}")
        return True
    except Exception as e:
        print(f"Error sending email: {e}")
        return False

def generate_report(results_file: str, recipient_email: str, output_dir: str = "report_output"):
    """Generate and email report from scan results"""
    
    # Load scan results
    try:
        with open(results_file, 'r') as f:
            results = json.load(f)
    except Exception as e:
        print(f"Error loading results file: {e}")
        return
    
    print("Generating visualizations...")
    chart_files = create_visualizations(results, output_dir)
    
    print("Generating HTML report...")
    html_content = generate_html_report(results, chart_files)
    
    print("Sending email report...")
    success = send_email_report(html_content, chart_files, recipient_email)
    
    if success:
        print("Report generation and email delivery completed successfully!")
    else:
        print("Failed to send email report. Check email configuration.")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Generate and email resume analysis report')
    parser.add_argument('--results', required=True, help='JSON file with scan results')
    parser.add_argument('--email', required=True, help='Recipient email address')
    parser.add_argument('--output-dir', default='report_output', help='Directory for output files')
    
    args = parser.parse_args()
    generate_report(args.results, args.email, args.output_dir)