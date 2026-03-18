#!/usr/bin/env python3
"""
Resume Scanner - PDF Resume Text Extraction and Analysis
Extracts text from PDF resumes and applies filtering criteria
"""

import json
import sys
import argparse
from pathlib import Path
import PyPDF2
import re
from typing import Dict, List, Any

def extract_text_from_pdf(pdf_path: str) -> str:
    """Extract text from PDF file"""
    text = ""
    try:
        with open(pdf_path, 'rb') as file:
            reader = PyPDF2.PdfReader(file)
            for page in reader.pages:
                text += page.extract_text() + "\n"
    except Exception as e:
        print(f"Error reading {pdf_path}: {e}")
        return ""
    return text

def analyze_resume(text: str, criteria: Dict[str, Any]) -> Dict[str, Any]:
    """Analyze resume text against criteria"""
    analysis = {
        "text": text,
        "matches": {},
        "score": 0,
        "details": {}
    }
    
    # Keyword matching
    if "keywords" in criteria:
        found_keywords = []
        for keyword in criteria["keywords"]:
            if re.search(r'\b' + re.escape(keyword.lower()) + r'\b', text.lower()):
                found_keywords.append(keyword)
        analysis["matches"]["keywords"] = found_keywords
        analysis["details"]["keyword_match_rate"] = len(found_keywords) / len(criteria["keywords"]) if criteria["keywords"] else 0
    
    # Education matching
    if "education" in criteria:
        found_education = []
        for edu in criteria["education"]:
            if re.search(r'\b' + re.escape(edu.lower()) + r'\b', text.lower()):
                found_education.append(edu)
        analysis["matches"]["education"] = found_education
        analysis["details"]["education_match_rate"] = len(found_education) / len(criteria["education"]) if criteria["education"] else 0
    
    # Experience years extraction
    if "experience_years" in criteria:
        experience_pattern = r'(\d+)\+?\s*(?:years?|yrs?)\s*(?:of\s*)?(?:experience|exp)'
        matches = re.findall(experience_pattern, text.lower())
        if matches:
            max_experience = max(int(m) for m in matches)
            analysis["matches"]["experience_years"] = max_experience
            analysis["details"]["experience_meets_requirement"] = max_experience >= criteria["experience_years"]
        else:
            analysis["matches"]["experience_years"] = 0
            analysis["details"]["experience_meets_requirement"] = False
    
    # Industry matching
    if "industries" in criteria:
        found_industries = []
        for industry in criteria["industries"]:
            if re.search(r'\b' + re.escape(industry.lower()) + r'\b', text.lower()):
                found_industries.append(industry)
        analysis["matches"]["industries"] = found_industries
        analysis["details"]["industry_match_rate"] = len(found_industries) / len(criteria["industries"]) if criteria["industries"] else 0
    
    # Calculate overall score
    scores = []
    if "keyword_match_rate" in analysis["details"]:
        scores.append(analysis["details"]["keyword_match_rate"])
    if "education_match_rate" in analysis["details"]:
        scores.append(analysis["details"]["education_match_rate"])
    if "experience_meets_requirement" in analysis["details"]:
        scores.append(1.0 if analysis["details"]["experience_meets_requirement"] else 0.0)
    if "industry_match_rate" in analysis["details"]:
        scores.append(analysis["details"]["industry_match_rate"])
    
    analysis["score"] = sum(scores) / len(scores) if scores else 0.0
    
    return analysis

def scan_resumes(input_dir: str, criteria_file: str, output_file: str):
    """Scan all PDF resumes in input directory"""
    # Load criteria
    try:
        with open(criteria_file, 'r') as f:
            criteria = json.load(f)
    except Exception as e:
        print(f"Error loading criteria file: {e}")
        return
    
    # Scan resumes
    results = {
        "criteria": criteria,
        "scanned_at": str(Path().resolve()),
        "resumes": [],
        "summary": {}
    }
    
    input_path = Path(input_dir)
    pdf_files = list(input_path.glob("*.pdf"))
    
    print(f"Found {len(pdf_files)} PDF files to scan...")
    
    for pdf_file in pdf_files:
        print(f"Scanning {pdf_file.name}...")
        
        text = extract_text_from_pdf(str(pdf_file))
        if not text:
            continue
            
        analysis = analyze_resume(text, criteria)
        
        resume_data = {
            "filename": pdf_file.name,
            "path": str(pdf_file),
            "analysis": analysis
        }
        
        results["resumes"].append(resume_data)
    
    # Generate summary
    total_resumes = len(results["resumes"])
    if total_resumes > 0:
        scores = [r["analysis"]["score"] for r in results["resumes"]]
        results["summary"] = {
            "total_resumes": total_resumes,
            "average_score": sum(scores) / len(scores),
            "high_match_count": len([s for s in scores if s >= 0.7]),
            "medium_match_count": len([s for s in scores if 0.4 <= s < 0.7]),
            "low_match_count": len([s for s in scores if s < 0.4])
        }
    
    # Save results
    try:
        with open(output_file, 'w') as f:
            json.dump(results, f, indent=2)
        print(f"Results saved to {output_file}")
        print(f"Summary: {results['summary']}")
    except Exception as e:
        print(f"Error saving results: {e}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Scan and analyze PDF resumes')
    parser.add_argument('--input-dir', required=True, help='Directory containing PDF resumes')
    parser.add_argument('--criteria-file', required=True, help='JSON file with filtering criteria')
    parser.add_argument('--output', required=True, help='Output JSON file for results')
    
    args = parser.parse_args()
    scan_resumes(args.input_dir, args.criteria_file, args.output)