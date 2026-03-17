# AgentFlow Dark Style Guide

**Style Overview**:
A bold dark minimalist **dark theme** with terminal-inspired precision, featuring electric lime green as primary brand color against pure black backgrounds, complemented by deep orange accents and razor-sharp cyan borders for maximum technical contrast and futuristic command center aesthetics.
Avoid gradients, rounded corners, shadows, and any colors not defined in this style.

## Colors
### Primary Colors
  - **primary-base**: `text-[#C6FF00]` or `bg-[#C6FF00]` - Electric lime green
  - **primary-lighter**: `text-[#DEFF66]` or `bg-[#DEFF66]`
  - **primary-darker**: `text-[#9FCC00]` or `bg-[#9FCC00]`

### Background Colors

#### Structural Backgrounds

Choose based on layout type:

**For Vertical Layout** (Top Header + Optional Side Panels):
- **bg-nav-primary**: `bg-black border-b border-[#00E5FF]` - Top header with cyan border
- **bg-nav-secondary**: `bg-black border-r border-[#00E5FF]` - Inner Left sidebar (if present) with cyan border
- **bg-page**: `bg-black` - Page background (bg of Main Content area)

**For Horizontal Layout** (Side Navigation + Optional Top Bar):
- **bg-nav-primary**: `bg-black border-r border-[#00E5FF]` - Left main sidebar with cyan border
- **bg-nav-secondary**: `bg-black border-b border-[#00E5FF]` - Inner Top header (if present) with cyan border
- **bg-page**: `bg-black` - Page background (bg of Main Content area)

#### Container Backgrounds
For main content area. Match background for seamless integration, using borders to define boundaries.
- **bg-container-primary**: `bg-black`
- **bg-container-secondary**: `bg-[#0A0A0A]` - Subtle lift for nested containers
- **bg-container-inset**: `bg-[#121212]`
- **bg-container-inset-strong**: `bg-[#1A1A1A]`

### Text Colors
- **color-text-primary**: `text-white/95`
- **color-text-secondary**: `text-white/70`
- **color-text-tertiary**: `text-white/50`
- **color-text-quaternary**: `text-white/30`
- **color-text-on-light-primary**: `text-black/90` - Text on primary-base (lime green) surfaces
- **color-text-on-light-secondary**: `text-black/70` - Text on primary-base (lime green) surfaces
- **color-text-link**: `text-[#C6FF00]` - Links, text-only buttons without backgrounds, and clickable text in tables

### Functional Colors
Use **sparingly** to maintain technical minimalist aesthetic. Used for status indicators, alerts, and critical system feedback.
  - **color-success-default**: #4CAF50 - success state bg
  - **color-success-light**: #66BB6A - tag/label bg
  - **color-error-default**: #F44336 - alert banner bg, error state
  - **color-error-light**: #EF5350 - tag/label bg
  - **color-warning-default**: #FFC107 - tag/label bg, alert banner bg
  - **color-warning-light**: #FFD54F - tag/label bg
  - **color-function-default**: #2196F3
  - **color-function-light**: #42A5F5 - tag/label bg

### Accent Colors
  - Secondary accent for categorization and emphasis. **Use sparingly** to maintain high-contrast technical aesthetic.
  - **accent-orange-deep**: `text-[#FF6D00]` or `bg-[#FF6D00]` - Deep orange for critical actions
  - **accent-orange-bright**: `text-[#FF9100]` or `bg-[#FF9100]` - Bright orange for highlights
  - **accent-cyan-border**: `#00E5FF` - Technical boundary color (used in borders only)

### Data Visualization Charts
For data visualization charts only.
  - Standard data colors: #C6FF00, #9FCC00, #7DA300, #5B7A00, #3A5200, #1A2900
  - Important data can use small amounts of: #FF6D00, #FF9100, #00E5FF, #00B8CC

## Typography
- **Font Stack**:
  - **font-family-base**: `-apple-system, BlinkMacSystemFont, "Segoe UI"` — For regular UI copy

- **Font Size & Weight**:
  - **Caption**: `text-base font-normal`
  - **Body**: `text-lg font-normal`
  - **Body Emphasized**: `text-lg font-semibold`
  - **Card Title / Subtitle**: `text-xl font-semibold`
  - **Page Title**: `text-2xl font-semibold`
  - **Headline**: `text-4xl font-semibold`

- **Line Height**: 1.5

## Border Radius
  - **Small**: 0px — Sharp technical precision
  - **Medium**: 0px
  - **Large**: 0px — Cards maintain razor-sharp edges
  - **Full**: full — Only for avatars, status indicators

## Layout & Spacing
  - **Tight**: 12px - For closely related small internal elements, such as icons and text within buttons
  - **Compact**: 16px - For small gaps between small containers, such as a line of tags
  - **Standard**: 24px - For gaps between medium containers like list items
  - **Relaxed**: 32px - For gaps between large containers and sections
  - **Section**: 40px - For major section divisions

## Create Boundaries (contrast of surface color, borders, shadows)
Terminal-inspired precision using high-contrast borders on matching backgrounds. No shadows, no surface color variation—boundaries are defined by razor-sharp cyan borders only.

### Borders
  - **Default**: `border border-[#00E5FF]` - 1px solid cyan for primary container boundaries
  - **Stronger**: `border-2 border-[#00E5FF]` - 2px solid cyan for emphasized/active states
  - **Subtle**: `border border-[#00E5FF]/50` - 1px semi-transparent cyan for secondary elements
  - **Primary Accent**: `border border-[#C6FF00]` - 1px solid lime green for special emphasis

### Dividers
  - **Default**: `border-t border-[#00E5FF]` or `border-b border-[#00E5FF]` - Horizontal dividers
  - **Vertical**: `border-l border-[#00E5FF]` or `border-r border-[#00E5FF]` - Vertical dividers

### Shadows & Effects
  - **No shadows** - Maintain pure flat technical aesthetic with border-based depth hierarchy

## Visual Emphasis for Containers
When containers (tags, cards, list items, rows) need visual emphasis to indicate priority, status, or category, use the following techniques:

| Technique | Implementation Notes | Best For | Avoid |
|-----------|---------------------|----------|-------|
| Background Tint | Use subtle gray variations (#0A0A0A, #121212, #1A1A1A) on black base | Gentle hierarchy within technical aesthetic | Bright colors on large areas |
| Border Highlight | Use cyan (#00E5FF) or lime green (#C6FF00) borders with varying thickness (1px/2px) | Active/selected states, critical elements | - |
| Status Tag/Label | Add colored tag/label inside container using functional or accent colors | Larger containers requiring status indication | - |
| Side Accent Bar | **Left edge only**, 3px solid colored bar (lime green/orange/cyan) for **non-rounded containers** | Terminal-style status indicators, active navigation items | Rounded containers |

## Assets
### Image

- For normal `<img>`: object-cover brightness-90 contrast-90
- For `<img>` with:
  - Slight overlay: object-cover brightness-85 contrast-90
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
  - **Graphic**: Use a simple, relevant icon (e.g., a `calendar` icon for a scheduling app, a `heart` icon for a dating app).

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
  - Example 1 (Primary action button with background):
    - button: flex items-center bg-[#C6FF00] text-black/90 border border-[#C6FF00] hover:bg-[#DEFF66] transition
      - span(button copy): whitespace-nowrap font-semibold
  - Example 2 (Secondary action button with border):
    - button: flex items-center bg-black text-[#C6FF00] border border-[#C6FF00] hover:bg-[#C6FF00]/10 transition
      - span(button copy): whitespace-nowrap font-semibold
  - Example 3 (text button):
    - button: flex items-center text-[#C6FF00] hover:text-[#DEFF66] transition
      - span(button copy): whitespace-nowrap
  - Example 4 (icon button):
    - button: flex items-center justify-center w-10 h-10 bg-black border border-[#00E5FF] text-white/95 hover:bg-[#00E5FF]/10 transition
      - icon

- **Tag Group (Filter Tags)** (Note: `overflow-x-auto` and `whitespace-nowrap` are required)
  - container(scrollable): flex overflow-x-auto gap-3 [&::-webkit-scrollbar]:hidden
    - label (Tag item 1):
      - input: type="radio" name="tag1" class="sr-only peer" checked
      - div: bg-black border border-[#00E5FF] text-white/70 peer-checked:bg-[#C6FF00] peer-checked:text-black/90 peer-checked:border-[#C6FF00] hover:bg-[#00E5FF]/10 transition whitespace-nowrap px-4 py-2

### Data Entry
- **Progress bars/Slider**: h-3 bg-[#121212] border border-[#00E5FF]
  - Progress fill: bg-[#C6FF00] h-full
- **Checkbox**
  - label: flex items-center gap-3
    - input: type="checkbox" class="sr-only peer"
    - div: w-5 h-5 bg-black border border-[#00E5FF] flex items-center justify-center peer-checked:bg-[#C6FF00] peer-checked:border-[#C6FF00] text-transparent peer-checked:text-black/90
      - svg(Checkmark): stroke="currentColor" stroke-width="3"
    - span(text): text-white/95
- **Radio button**
  - label: flex items-center gap-3
    - input: type="radio" name="option1" class="sr-only peer"
    - div: w-5 h-5 bg-black border border-[#00E5FF] rounded-full flex items-center justify-center peer-checked:bg-[#C6FF00] peer-checked:border-[#C6FF00] text-transparent peer-checked:text-black/90
      - svg(dot indicator): fill="currentColor" width="8" height="8"
    - span(text): text-white/95
- **Switch/Toggle**
  - label: flex items-center gap-3
    - div: relative
      - input: type="checkbox" class="sr-only peer"
      - div(Toggle track): w-14 h-7 bg-black border border-[#00E5FF] peer-checked:bg-[#C6FF00] peer-checked:border-[#C6FF00] transition
      - div(Toggle thumb): absolute top-0.5 left-0.5 w-6 h-6 bg-[#00E5FF] peer-checked:bg-black peer-checked:translate-x-7 transition
    - span(text): text-white/95

- **Select/Dropdown**
  - Select container: flex items-center gap-2 bg-black border border-[#00E5FF] px-4 py-2
    - text: text-white/95
    - Dropdown icon(square container): flex items-center justify-center bg-transparent w-5 h-5
      - icon: text-[#00E5FF]

### Container
- **Navigation Menu - horizontal**
    - Navigation with sections/grouping:
        - Nav Container: flex items-center justify-between w-full gap-8
        - Left Section: flex items-center gap-10
          - Menu Item: flex items-center gap-3 text-white/70 hover:text-[#C6FF00] transition border-b-2 border-transparent hover:border-[#C6FF00]
        - Right Section: flex items-center gap-4
          - Menu Item: flex items-center gap-3 text-white/70 hover:text-[#C6FF00] transition
          - Notification (if applicable): relative flex items-center justify-center w-12 h-12 border border-[#00E5FF] hover:bg-[#00E5FF]/10 transition
            - notification-icon: w-6 h-6 text-white/95
            - badge (if has unread): absolute -top-1 -right-1 w-6 h-6 bg-[#FF6D00] border-2 border-black flex items-center justify-center
              - badge-count: text-xs font-semibold text-white
          - Avatar(if applicable): flex items-center gap-2
            - avatar-image: w-10 h-10 rounded-full border-2 border-[#00E5FF]
            - dropdown-icon (if applicable): w-6 h-6 text-white/70

- **Card**
    - Example 1 (Vertical card with image and text):
        - Card: bg-black border border-[#00E5FF] flex flex-col p-4 gap-4
        - Image: w-full h-48 border border-[#00E5FF]/50
        - Text area: flex flex-col gap-3
          - card-title: text-xl font-semibold text-white/95
          - card-subtitle: text-base font-normal text-white/70
    - Example 2 (Horizontal card with image and text):
        - Card: bg-black border border-[#00E5FF] flex gap-4 p-4
        - Image: h-full w-48 border border-[#00E5FF]/50
        - Text area: flex flex-col gap-4 flex-1
          - card-title: text-xl font-semibold text-white/95
          - card-subtitle: text-base font-normal text-white/70
    - Example 3 (Image-focused card: minimal container, emphasis on content):
        - Card: flex flex-col gap-4
        - Image: w-full h-64 border border-[#00E5FF]
        - Text area: flex flex-col gap-3
          - card-title: text-xl font-semibold text-white/95
          - card-subtitle: text-base font-normal text-white/70
    - Example 4 (text-only cards, terminal-style information cards):
        - Card: bg-black border-l-4 border-l-[#C6FF00] border border-[#00E5FF] flex flex-col p-6 gap-4

## Additional Notes
- **Terminal Command Aesthetic**: Code blocks and technical elements should use `bg-[#0A0A0A]` with `border border-[#00E5FF]` for authentic command-line feel
- **Status Indicators**: Use 3px solid left border bars in lime green (#C6FF00), orange (#FF6D00), or cyan (#00E5FF) for system status and active states
- **High Contrast Priority**: Ensure minimum 7:1 contrast ratio for all text on black backgrounds for maximum readability
- **No Anti-aliasing Simulation**: Maintain crisp, pixel-perfect borders—avoid opacity tricks that soften technical precision
- **Monospace for Data**: Technical data, logs, and numerical values should consider monospace presentation for enhanced readability

<colors_extraction>
#C6FF00
#DEFF66
#9FCC00
#000000
#00E5FF
#0A0A0A
#121212
#1A1A1A
#FFFFFFF2
#FFFFFFB3
#FFFFFF80
#0000004D
#000000E6
#000000B3
#4CAF50
#66BB6A
#F44336
#EF5350
#FFC107
#FFD54F
#2196F3
#42A5F5
#FF6D00
#FF9100
#00E5FF80
#7DA300
#5B7A00
#3A5200
#1A2900
#00B8CC
</colors_extraction>
