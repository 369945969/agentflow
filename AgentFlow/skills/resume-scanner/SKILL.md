name: resume-scanner
description: Comprehensive resume scanning and analysis system for PDF documents. Use when you need to: (1) Scan and extract text from PDF resumes, (2) Filter resumes based on custom criteria (keywords, education, experience), (3) Generate detailed reports with visualizations, (4) Send analysis reports via email. Supports automated resume screening workflows with customizable matching algorithms and statistical reporting.

# Resume Scanner

## Quick Start

Scan PDF resumes and filter by criteria:

```bash
python3 scripts/scan_resumes.py --input-dir ./resumes --criteria-file criteria.json --output results.json
```

Generate and email report:

```bash
python3 scripts/generate_report.py --results results.json --email recipient@example.com
```

## Core Workflow

1. **Setup criteria**: Define filtering criteria in `references/criteria_template.json`
2. **Scan resumes**: Extract and analyze PDF resumes using `scripts/scan_resumes.py`
3. **Filter matches**: Apply keyword and structured criteria filtering
4. **Generate report**: Create visualization-rich report using `scripts/generate_report.py`
5. **Send email**: Deliver report via email using configured SMTP settings

## Scripts

- **`scan_resumes.py`**: PDF text extraction and resume analysis
- **`generate_report.py`**: Report generation with charts and email delivery
- **`setup_email.py`**: Configure email settings (run once)

## References

- **`criteria_template.json`**: Template for defining filtering criteria
- **`email_config.md`**: Email configuration guide
- **`chart_templates.md`**: Visualization templates and examples

## Email Configuration

Before first use, configure email settings:

```bash
python3 scripts/setup_email.py
```

Or manually edit `references/email_config.json` with your SMTP details.

## Criteria Definition

Create filtering criteria using JSON structure:

```json
{
  "keywords": ["Python", "React", "Machine Learning"],
  "education": ["Bachelor", "Master"],
  "experience_years": 2,
  "industries": ["Technology", "Finance"]
}
```

See `references/criteria_template.json` for complete schema.
