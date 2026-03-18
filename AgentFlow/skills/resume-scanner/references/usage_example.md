# Resume Scanner Usage Example

This example demonstrates how to use the Resume Scanner system to analyze PDF resumes and generate reports.

## Step 1: Install Dependencies

```bash
python3 scripts/install_dependencies.py
```

## Step 2: Configure Email

```bash
python3 scripts/setup_email.py
```

## Step 3: Define Your Criteria

Edit `references/criteria_template.json`:

```json
{
  "keywords": ["Python", "React", "SQL"],
  "education": ["Bachelor", "Master"],
  "experience_years": 3,
  "industries": ["Technology", "Finance"]
}
```

## Step 4: Prepare Resume Directory

Create a directory with PDF resumes:

```
resumes/
├── john_doe.pdf
├── jane_smith.pdf
├── robert_chen.pdf
└── sarah_jones.pdf
```

## Step 5: Scan Resumes

```bash
python3 scripts/scan_resumes.py \
  --input-dir ./resumes \
  --criteria-file references/criteria_template.json \
  --output scan_results.json
```

Expected output:

```
Found 4 PDF files to scan...
Scanning john_doe.pdf...
Scanning jane_smith.pdf...
Scanning robert_chen.pdf...
Scanning sarah_jones.pdf...
Results saved to scan_results.json
Summary: {'total_resumes': 4, 'average_score': 0.65, 'high_match_count': 1, 'medium_match_count': 2, 'low_match_count': 1}
```

## Step 6: Generate and Email Report

```bash
python3 scripts/generate_report.py \
  --results scan_results.json \
  --email hr@company.com
```

## Complete Workflow Script

Create a workflow script `run_full_scan.sh`:

```bash
#!/bin/bash

# Set your parameters
RESUME_DIR="./resumes"
CRITERIA_FILE="references/criteria_template.json"
RESULTS_FILE="scan_results.json"
EMAIL_RECIPIENT="hr@company.com"

echo "Starting Resume Scanner workflow..."

# Step 1: Scan resumes
echo "Step 1: Scanning resumes..."
python3 scripts/scan_resumes.py \
  --input-dir "$RESUME_DIR" \
  --criteria-file "$CRITERIA_FILE" \
  --output "$RESULTS_FILE"

if [ $? -eq 0 ]; then
    echo "✓ Resume scanning completed"

    # Step 2: Generate and send report
    echo "Step 2: Generating report..."
    python3 scripts/generate_report.py \
      --results "$RESULTS_FILE" \
      --email "$EMAIL_RECIPIENT"

    if [ $? -eq 0 ]; then
        echo "✓ Report generation and email completed successfully!"
    else
        echo "✗ Report generation failed"
        exit 1
    fi
else
    echo "✗ Resume scanning failed"
    exit 1
fi

echo "Workflow completed!"
```

## Custom Criteria Examples

### Technical Hiring

```json
{
  "keywords": ["JavaScript", "React", "Node.js", "AWS", "Docker"],
  "education": ["Computer Science", "Engineering"],
  "experience_years": 2,
  "industries": ["Technology", "Software"]
}
```

### Senior Management

```json
{
  "keywords": ["Leadership", "Strategy", "Management", "MBA"],
  "education": ["Master", "MBA", "PhD"],
  "experience_years": 8,
  "industries": ["Technology", "Finance", "Consulting"]
}
```

### Data Science

```json
{
  "keywords": ["Python", "Machine Learning", "Statistics", "SQL", "R"],
  "education": ["Statistics", "Mathematics", "Computer Science"],
  "experience_years": 3,
  "industries": ["Technology", "Research", "Analytics"]
}
```

## Troubleshooting

### PDF Reading Issues

- Ensure PDFs are text-based (not scanned images)
- Check that PDF files are not password-protected
- Verify PDF files are not corrupted

### Email Issues

- Check SMTP configuration in `references/email_config.json`
- Verify firewall allows SMTP traffic
- Ensure app passwords are used for Gmail

### Performance Issues

- For large batches (>100 resumes), consider processing in smaller groups
- Monitor memory usage with large PDF files
- Use SSD storage for faster PDF reading
