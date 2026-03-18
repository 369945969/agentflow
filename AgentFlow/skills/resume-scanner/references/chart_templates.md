# Chart Templates and Examples

This document describes the visualization templates used in resume analysis reports.

## Available Charts

### 1. Score Distribution Bar Chart

- **Purpose**: Shows individual resume match scores
- **X-axis**: Resume files (ordered by score)
- **Y-axis**: Match score (0-1 or 0-100%)
- **Color coding**:
  - Green: High match (≥70%)
  - Orange: Medium match (40-70%)
  - Red: Low match (<40%)

### 2. Match Categories Pie Chart

- **Purpose**: Shows distribution of resume quality
- **Segments**: High, Medium, Low match categories
- **Percentages**: Automatically calculated
- **Colors**: Green, Orange, Red

### 3. Keyword Analysis Bar Chart

- **Purpose**: Shows frequency of matched keywords
- **X-axis**: Keywords found in resumes
- **Y-axis**: Number of resumes containing each keyword
- **Color**: Sky blue

## Customization Options

### Chart Size

```python
plt.figure(figsize=(12, 6))  # Width, height in inches
```

### Colors

```python
colors = ['green', 'orange', 'red']  # Standard color palette
colors = ['skyblue', 'lightcoral']   # Custom colors
```

### Chart Titles and Labels

```python
plt.title('Custom Chart Title')
plt.xlabel('X-axis Label')
plt.ylabel('Y-axis Label')
```

## File Formats

- Charts saved as PNG format
- DPI: 300 (high quality for email embedding)
- Transparent backgrounds

## Additional Chart Ideas

### Experience Years Distribution

```python
# Histogram of experience years
plt.hist(experience_years, bins=5, alpha=0.7)
```

### Education Level Breakdown

```python
# Pie chart of education levels
education_counts = {'Bachelor': 45, 'Master': 25, 'PhD': 10}
plt.pie(education_counts.values(), labels=education_counts.keys())
```

### Industry Distribution

```python
# Bar chart of industry matches
industry_counts = {'Technology': 35, 'Finance': 15, 'Healthcare': 10}
plt.bar(industry_counts.keys(), industry_counts.values())
```

## Integration with Reports

Charts are embedded in HTML emails using CID (Content-ID) references:

- HTML: `<img src="cid:score_distribution" alt="Score Distribution">`
- MIME header: `img.add_header('Content-ID', '<score_distribution>')`

## Performance Notes

- Large datasets (>100 resumes) may need reduced chart complexity
- Consider sampling for very large resume pools
- Chart generation time scales linearly with resume count
