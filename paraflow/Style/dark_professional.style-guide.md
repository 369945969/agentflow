# AgentFlow Dark Professional Style Guide

**Style Overview**:
A clean, dark professional aesthetic optimized for data-focused clarity, featuring vibrant cyan as the primary brand color against a deep charcoal foundation. Uses solid color surfaces with clear contrast to establish structural boundaries and functional hierarchy, avoiding shadows entirely for maximum information density and technical precision.

## Colors
### Primary Colors
  - **primary-base**: `text-[#00D9FF]` or `bg-[#00D9FF]`
  - **primary-lighter**: `text-[#33E3FF]` or `bg-[#33E3FF]`
  - **primary-darker**: `text-[#00B8DB]` or `bg-[#00B8DB]`

### Background Colors

#### Structural Backgrounds
Choose based on layout type:

**For Vertical Layout** (Top Header + Optional Side Panels):
- **bg-nav-primary**: `bg-[#1A1D23]` - Top header
- **bg-nav-secondary**: `bg-[#1F2228]` - Inner Left sidebar (if present)
- **bg-page**: `bg-[#24272E]` - Page background (bg of Main Content area)

**For Horizontal Layout** (Side Navigation + Optional Top Bar):
- **bg-nav-primary**: `bg-[#1A1D23]` - Left main sidebar
- **bg-nav-secondary**: `bg-[#1F2228]` - Inner Top header (if present)
- **bg-page**: `bg-[#24272E]` - Page background (bg of Main Content area)

#### Container Backgrounds
For main content area. Adjust values when used on navigation backgrounds to ensure sufficient contrast.
- **bg-container-primary**: `bg-[#2A2E36]`
- **bg-container-secondary**: `bg-[#2F333C]`
- **bg-container-inset**: `bg-[#00D9FF]/8`
- **bg-container-inset-strong**: `bg-[#00D9FF]/15`

### Text Colors
- **color-text-primary**: `text-white/95`
- **color-text-secondary**: `text-white/75`
- **color-text-tertiary**: `text-white/55`
- **color-text-quaternary**: `text-white/35`
- **color-text-on-light-primary**: `text-black/85` - Text on light backgrounds and primary-base color surfaces
- **color-text-on-light-secondary**: `text-black/65` - Text on light backgrounds and primary-base color surfaces
- **color-text-link**: `text-[#00D9FF]` - Links, text-only buttons without backgrounds, and clickable text in tables

### Functional Colors
Use **sparingly** to maintain a professional and technical overall style. Used for the surfaces of text-only cards, simple cards, buttons, and tags.
  - **color-success-default**: `#00C896` - status indicators
  - **color-success-light**: `#00C896/20` - tag/label bg
  - **color-error-default**: `#FF4757` - alert banner bg, status indicators
  - **color-error-light**: `#FF4757/20` - tag/label bg
  - **color-warning-default**: `#FFA726` - alert banner bg, status indicators
  - **color-warning-light**: `#FFA726/20` - tag/label bg
  - **color-function-default**: `#5B8DEE`
  - **color-function-light**: `#5B8DEE/20` - tag/label bg

### Accent Colors
  - A secondary palette for data categorization and visual hierarchy. Use strategically for clarity.
  - **accent-electric-blue**: `text-[#4D9FFF]` or `bg-[#4D9FFF]`
  - **accent-soft-teal**: `text-[#4ECDC4]` or `bg-[#4ECDC4]`
  - **accent-slate-gray**: `text-[#8B93A7]` or `bg-[#8B93A7]`

### Data Visualization Charts
For data visualization charts only.
  - Standard data colors: #00D9FF, #4D9FFF, #4ECDC4, #8B93A7, #7C5CDB, #A78BFA
  - Important/Critical data: #FF4757, #FFA726, #00C896
  - Neutral data: #52596B, #6B7280, #9CA3AF

## Typography
- **Font Stack**:
  - **font-family-base**: `-apple-system, BlinkMacSystemFont, "Segoe UI"` — For regular UI copy

- **Font Size & Weight**:
  - **Caption**: `text-sm font-normal`
  - **Body**: `text-base font-normal`
  - **Body Emphasized**: `text-base font-semibold`
  - **Card Title / Subtitle**: `text-lg font-semibold`
  - **Page Title**: `text-2xl font-semibold`
  - **Headline**: `text-4xl font-semibold`

- **Line Height**: 1.5

## Border Radius
  - **Small**: 6px — Elements inside cards, small components
  - **Medium**: 8px — Buttons, inputs, tags
  - **Large**: 12px — Cards, panels
  - **Full**: full — Toggles, avatars, status indicators

## Layout & Spacing
  - **Tight**: 8px - For closely related small internal elements, such as icons and text within buttons
  - **Compact**: 12px - For small gaps between small containers, such as a line of tags
  - **Standard**: 16px - For gaps between medium containers like list items
  - **Relaxed**: 24px - For gaps between large containers and sections
  - **Section**: 32px - For major section divisions

## Create Boundaries (contrast of surface color, borders, shadows)
No shadows, primarily relying on surface color contrast with clear differentiation to create distinct structural boundaries and functional hierarchy.

### Borders
  - **Case 1**: For most containers and structural elements - no borders, rely on surface color contrast.
  - **Case 2**: For specific functional elements requiring explicit boundaries:
    - **Default**: 1px solid rgba(255,255,255,0.08). Used for inputs, data tables, separators. `border border-white/8`
    - **Stronger**: 1px solid rgba(255,255,255,0.12). Used for active or focused states. `border border-white/12`
    - **Accent**: 1px solid #00D9FF. Used for selected or highlighted states. `border border-[#00D9FF]`

### Dividers
  - **Case 1**: For major structural divisions - no dividers, rely on surface color contrast.
  - **Case 2**: For data tables, lists, and content sections - `border-t` or `border-b` `border-white/8`.

### Shadows & Effects
  - **Case 1**: No shadows - maintain flat, data-focused clarity throughout the interface.

## Visual Emphasis for Containers
When containers (tags, cards, list items, rows) need visual emphasis to indicate priority, status, or category, use the following techniques:

| Technique | Implementation Notes | Best For | Avoid |
|-----------|---------------------|----------|-------|
| Background Tint | Use progressively lighter/darker shades or adjust opacity | Common approach for moderate emphasis needs | Excessive brightness that reduces readability |
| Border Highlight | Use thin border with accent colors (cyan, blue, teal) | Active/selected states, form validation, critical alerts | - |
| Status Tag/Label | Add colored tag/label inside container | Larger containers, status indicators | - |
| Side Accent Bar | **Left edge only**, for **non-rounded containers**, use 2-3px width with accent colors | Data rows, list items, alert panels | Large cards, rounded containers |

## Assets
### Image
  - For normal `<img>`: object-cover brightness-90 contrast-90
  - For `<img>` with:
    - Slight overlay: object-cover brightness-75 contrast-90
    - Heavy overlay: object-cover brightness-50 contrast-90

### Icon
- Use Lucide icons from Iconify.
- To ensure an aesthetic layout, each icon should be centered in a square container, typically without a background, matching the icon's size.
- Use Tailwind font size to control icon size
- Example:
  ```html
  <div class="flex items-center justify-center bg-transparent w-5 h-5">
  <iconify-icon icon="lucide:flag" class="text-base"></iconify-icon>
  </div>
  ```

### Third-Party Brand Logos:
   - Use Brand Icons from Iconify.
   - Logo Example:
     Monochrome Logo: `<iconify-icon icon="simple-icons:x"></iconify-icon>`
     Colored Logo: `<iconify-icon icon="logos:google-icon"></iconify-icon>`

### User's Own Logo:
- To protect copyright, do **NOT** use real product logos as a logo for a new product, individual user, or other company products.
- **Icon-based**:
  - **Graphic**: Use a simple, relevant icon (e.g., a `workflow` icon for AgentFlow, a `network` icon for orchestration platform).

## Page Layout - Web (*EXTREMELY* important)
### Determine Layout Type
- Choose between Vertical or Horizontal layout based on whether the primary navigation is a full-width top header or a full-height sidebar (left/right).
- User requirements typically indicate the layout preference. If unclear, consider:
  - Marketing/content sites typically use Vertical Layout.
  - Functional/dashboard sites can use either, depending on visual style. Sidebars accommodate more complex navigation than top bars. For complex navigation needs with a preference for minimal chrome (Vertical Layout adds an extra fixed header), choose Horizontal Layout (omits the fixed top header).
- Vertical Layout Diagram:
┌──────────────────────────────────────────────────────┐
│  Header (Primary Nav)                                │
├──────────┬──────────────────────────────┬────────────┤
│Left      │ Sub-header (Tertiary Nav)    │ Right      │
│Sidebar   │ (optional)                   │ Sidebar    │
│(Secondary├──────────────────────────────┤ (Utility   │
│Nav)      │ Main Content                 │ Panel)     │
│(optional)│                              │ (optional) │
│          │                              │            │
└──────────┴──────────────────────────────┴────────────┘
- Horizontal Layout Diagram:
┌──────────┬──────────────────────────────┬───────────┐
│          │ Header (Secondary Nav)       │           │
│ Left     │ (optional)                   │ Right     │
│ Sidebar  ├──────────────────────────────┤ Sidebar   │
│ (Primary │ Main Content                 │ (Utility  │
│ Nav)     │                              │ Panel)    │
│          │                              │ (optional)│
│          │                              │           │
└──────────┴──────────────────────────────┴───────────┘
### Detailed Layout Code
**Vertical Layout**
```html
<!-- Body: Adjust width (w-[1440px]) based on target screen size -->
<body class="w-[1440px] min-h-[700px] font-[-apple-system,BlinkMacSystemFont,'Segoe UI'] leading-[1.5]">

  <!-- Header (Primary Nav): Fixed height -->
  <header class="w-full">
    <!-- Header content -->
  </header>

  <!-- Content Container: Must include 'flex' class -->
  <div class="w-full flex min-h-[700px]">
    <!-- Left Sidebar (Secondary Nav) (Optional): Remove if not needed. If Left Sidebar exists, use its ml to control left page margin -->
    <aside class="flex-shrink-0 min-w-fit">

    </aside>

    <!-- Main Content Area:
     Use Main Content Area's horizontal padding (px) to control distance from main content to sidebars or page edges.
     For pages without sidebars (like Marketing Pages, simple content pages such as help centers, privacy policies) use larger values (px-30 to px-80), for pages with sidebars (Functional/Dashboard Pages, complex content pages with multi-level navigation like knowledge base articles) use moderate values (px-8 to px-16) -->
    <main class="flex-1 overflow-x-hidden flex flex-col">
    <!--  Main Content -->

    </main>

    <!-- Right Sidebar (Utility Panel) (Optional): Remove if not needed. If Right Sidebar exists, use its mr to control right page margin -->
    <aside class="flex-shrink-0 min-w-fit">
    </aside>

  </div>
</body>
```

**Horizontal Layout**

```html
<!-- Body: Adjust width (w-[1440px]) based on target screen size. Must include 'flex' class -->
<body class="w-[1440px] min-h-[700px] flex font-[-apple-system,BlinkMacSystemFont,'Segoe UI'] leading-[1.5]">

<!-- Left Sidebar (Primary Nav): Use its ml to control left page margin -->
  <aside class="flex-shrink-0 min-w-fit">
  </aside>

  <!-- Content Container-->
  <div class="flex-1 overflow-x-hidden flex flex-col min-h-[700px]">

    <!-- Header (Secondary Nav) (Optional): Remove if not needed. If Header exists, use its mx to control distance to left/right sidebars or page margins -->
    <header class="w-full">
    </header>

    <!-- Main Content Area: Use Main Content Area's pl to control distance from main content to left sidebar. Use pr to control distance to right sidebar/right page edge -->
    <main class="w-full">
    </main>


  </div>

  <!-- Right Sidebar (Utility Panel) (Optional): Remove if not needed. If Right Sidebar exists, use its mr to control right page margin -->
  <aside class="flex-shrink-0 min-w-fit">
  </aside>

</body>
```

## Tailwind Component Examples (Key attributes)
**Important Note**: Use utility classes directly. Do NOT create custom CSS classes or add styles in <style> tags for the following components

### Basic

- **Button**: (Note: Use flex and items-center for the container)
  - Example 1 (Primary button):
    - button: flex items-center gap-2 bg-[#00D9FF] text-black/85 px-4 py-2 rounded-lg hover:bg-[#33E3FF] transition
      - icon (if applicable)
      - span(button copy): whitespace-nowrap font-semibold
  - Example 2 (Secondary button):
    - button: flex items-center gap-2 bg-[#2A2E36] text-white/95 px-4 py-2 rounded-lg hover:bg-[#2F333C] transition border border-white/8
      - icon (if applicable)
      - span(button copy): whitespace-nowrap
  - Example 3 (text button):
    - button: flex items-center gap-2 text-[#00D9FF] hover:text-[#33E3FF] transition
      - icon (if applicable)
      - span(button copy): whitespace-nowrap
  - Example 4 (icon button):
    - button: flex items-center justify-center w-9 h-9 bg-[#2A2E36] rounded-lg hover:bg-[#2F333C] transition
      - icon

- **Tag Group (Filter Tags)** (Note: `overflow-x-auto` and `whitespace-nowrap` are required)
  - container(scrollable): flex gap-2 overflow-x-auto [&::-webkit-scrollbar]:hidden
    - label (Tag item 1):
      - input: type="radio" name="tag1" class="sr-only peer" checked
      - div: bg-[#2A2E36] text-white/75 px-3 py-1.5 rounded-lg peer-checked:bg-[#00D9FF] peer-checked:text-black/85 hover:bg-[#2F333C] transition whitespace-nowrap text-sm

### Data Entry
- **Progress bars/Slider**: h-2 bg-[#2A2E36] rounded-full
  - progress-fill: bg-[#00D9FF] h-full rounded-full
- **Checkbox**
  - label: flex items-center gap-2
    - input: type="checkbox" class="sr-only peer"
    - div: w-5 h-5 bg-[#2A2E36] rounded-md flex items-center justify-center peer-checked:bg-[#00D9FF] text-transparent peer-checked:text-black/85 border border-white/8 peer-checked:border-[#00D9FF] transition
      - svg(Checkmark): stroke="currentColor" stroke-width="3"
    - span(text): text-white/95
- **Radio button**
  - label: flex items-center gap-2
    - input: type="radio" name="option1" class="sr-only peer"
    - div: w-5 h-5 bg-[#2A2E36] rounded-full flex items-center justify-center peer-checked:bg-[#00D9FF] text-transparent peer-checked:text-black/85 border border-white/8 peer-checked:border-[#00D9FF] transition
      - svg(dot indicator): fill="currentColor" width="8" height="8"
    - span(text): text-white/95
- **Switch/Toggle**
  - label: flex items-center gap-2
    - div: relative
      - input: type="checkbox" class="sr-only peer"
      - div(Toggle track): w-11 h-6 bg-[#2A2E36] peer-checked:bg-[#00D9FF] rounded-full transition
      - div(Toggle thumb): absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full peer-checked:translate-x-5 transition shadow-sm
    - span(text): text-white/95

- **Input Field**
  - container: flex flex-col gap-1
    - label: text-sm text-white/75
    - input: bg-[#2A2E36] border border-white/8 rounded-lg px-3 py-2 text-white/95 placeholder:text-white/35 focus:border-[#00D9FF] focus:outline-none transition

- **Select/Dropdown**
  - Select container: flex items-center gap-2 bg-[#2A2E36] border border-white/8 rounded-lg px-3 py-2 hover:bg-[#2F333C] transition
    - text: text-white/95
    - Dropdown icon(square container): flex items-center justify-center bg-transparent w-4 h-4
      - icon: text-white/55

### Container
- **Navigation Menu - horizontal**
    - Navigation with sections/grouping:
        - Nav Container: flex items-center justify-between w-full px-6 py-4
        - Left Section: flex items-center gap-8
          - Menu Item: flex items-center gap-2 text-white/75 hover:text-[#00D9FF] transition
        - Right Section: flex items-center gap-4
          - Menu Item: flex items-center gap-2 text-white/75 hover:text-[#00D9FF] transition
          - Notification (if applicable): relative flex items-center justify-center w-9 h-9 bg-[#2A2E36] rounded-lg hover:bg-[#2F333C] transition
            - notification-icon: w-5 h-5 text-white/75
            - badge (if has unread): absolute -top-1 -right-1 w-5 h-5 bg-[#FF4757] rounded-full flex items-center justify-center
              - badge-count: text-xs text-white font-semibold
          - Avatar(if applicable): flex items-center gap-2
            - avatar-image: w-8 h-8 rounded-full border border-white/8
            - dropdown-icon (if applicable): w-4 h-4 text-white/55

- **Navigation Menu - vertical sidebar**
    - Nav Container: flex flex-col gap-1 p-3
      - Menu Item: flex items-center gap-3 px-3 py-2.5 rounded-lg text-white/75 hover:bg-[#2A2E36] hover:text-[#00D9FF] transition
        - icon: w-5 h-5
        - span: text-sm
      - Menu Item (active): flex items-center gap-3 px-3 py-2.5 rounded-lg bg-[#00D9FF]/15 text-[#00D9FF]
        - icon: w-5 h-5
        - span: text-sm font-semibold

- **Card**
    - Example 1 (Data card - metrics/stats):
        - Card: bg-[#2A2E36] rounded-xl flex flex-col p-5 gap-4 border border-white/8
        - Header: flex items-center justify-between
          - title: text-sm text-white/75
          - icon: w-4 h-4 text-white/55
        - Content: flex flex-col gap-2
          - metric-value: text-3xl font-semibold text-white/95
          - metric-change: text-xs flex items-center gap-1
            - icon: w-3 h-3
            - text
    - Example 2 (Agent card - horizontal layout):
        - Card: bg-[#2A2E36] rounded-xl flex gap-4 p-4 border border-white/8 hover:border-[#00D9FF]/30 transition
        - Icon area: flex items-center justify-center w-12 h-12 bg-[#00D9FF]/10 rounded-lg
          - icon: w-6 h-6 text-[#00D9FF]
        - Text area: flex-1 flex flex-col gap-1
          - card-title: text-base font-semibold text-white/95
          - card-subtitle: text-sm text-white/55
    - Example 3 (Status card with accent bar):
        - Card: bg-[#2A2E36] rounded-xl flex gap-4 p-4 border-l-2 border-l-[#00D9FF] border border-white/8
        - Content: flex flex-col gap-3
          - card-title: text-base font-semibold text-white/95
          - card-content: text-sm text-white/75

- **Table**
    - Table container: bg-[#2A2E36] rounded-xl border border-white/8 overflow-hidden
      - table: w-full
        - thead: bg-[#1F2228]
          - tr
            - th: px-4 py-3 text-left text-xs font-semibold text-white/75 uppercase tracking-wide
        - tbody
          - tr: border-t border-white/8 hover:bg-[#2F333C] transition
            - td: px-4 py-3 text-sm text-white/95

- **Alert/Banner**
    - Info banner: flex items-start gap-3 bg-[#5B8DEE]/15 border border-[#5B8DEE]/30 rounded-lg p-4
      - icon: w-5 h-5 text-[#5B8DEE] flex-shrink-0
      - content: flex flex-col gap-1
        - title: text-sm font-semibold text-white/95
        - message: text-sm text-white/75
    - Warning banner: flex items-start gap-3 bg-[#FFA726]/15 border border-[#FFA726]/30 rounded-lg p-4
      - icon: w-5 h-5 text-[#FFA726] flex-shrink-0
      - content: flex flex-col gap-1
        - title: text-sm font-semibold text-white/95
        - message: text-sm text-white/75
    - Error banner: flex items-start gap-3 bg-[#FF4757]/15 border border-[#FF4757]/30 rounded-lg p-4
      - icon: w-5 h-5 text-[#FF4757] flex-shrink-0
      - content: flex flex-col gap-1
        - title: text-sm font-semibold text-white/95
        - message: text-sm text-white/75

## Additional Notes
- **High Information Density**: Optimize for displaying complex data and monitoring information without visual clutter
- **Technical Precision**: Maintain clarity and readability for detailed technical content and system metrics
- **Consistent Contrast**: Ensure all text and interactive elements meet accessibility standards against dark backgrounds
- **Data Hierarchy**: Use color temperature and brightness strategically to establish clear visual hierarchy in data-heavy interfaces
- **Status Communication**: Leverage functional colors consistently for system status, alerts, and monitoring indicators

<colors_extraction>
#00D9FF
#33E3FF
#00B8DB
#1A1D23
#1F2228
#24272E
#2A2E36
#2F333C
#00D9FF14
#00D9FF26
#FFFFFFF2
#FFFFFFBF
#FFFFFF8C
#FFFFFF59
#000000D9
#000000A6
#00C896
#00C89633
#FF4757
#FF475733
#FFA726
#FFA72633
#5B8DEE
#5B8DEE33
#4D9FFF
#4ECDC4
#8B93A7
#7C5CDB
#A78BFA
#52596B
#6B7280
#9CA3AF
#FFFFFF14
#FFFFFF1F
#00D9FF4D
</colors_extraction>
