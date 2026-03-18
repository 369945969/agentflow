#!/usr/bin/env python3
"""
Dependencies Installer - Install required Python packages
Installs all dependencies needed for the Resume Scanner system
"""

import subprocess
import sys

def install_package(package):
    """Install a Python package using pip"""
    try:
        subprocess.check_call([sys.executable, "-m", "pip", "install", package])
        print(f"✓ {package} installed successfully")
        return True
    except subprocess.CalledProcessError as e:
        print(f"✗ Failed to install {package}: {e}")
        return False

def main():
    """Install all required packages"""
    
    print("=== Resume Scanner Dependencies Installation ===")
    print()
    
    packages = [
        "PyPDF2",          # PDF text extraction
        "matplotlib",      # Chart generation
        "pandas",          # Data manipulation (optional but recommended)
        "Pillow",          # Image processing for charts
        "requests",        # HTTP requests (for future features)
    ]
    
    failed_packages = []
    
    for package in packages:
        print(f"Installing {package}...")
        if not install_package(package):
            failed_packages.append(package)
        print()
    
    if failed_packages:
        print(f"Installation completed with {len(failed_packages)} failures:")
        for package in failed_packages:
            print(f"  - {package}")
        print("\nPlease install these packages manually:")
        print(f"pip install {' '.join(failed_packages)}")
    else:
        print("✓ All packages installed successfully!")
    
    print("\nYou can now use the Resume Scanner system.")
    print("Next steps:")
    print("1. Configure email: python3 scripts/setup_email.py")
    print("2. Prepare your criteria: edit references/criteria_template.json")
    print("3. Scan resumes: python3 scripts/scan_resumes.py --input-dir ./resumes --criteria-file references/criteria_template.json --output results.json")
    print("4. Generate report: python3 scripts/generate_report.py --results results.json --email your-email@example.com")

if __name__ == "__main__":
    main()