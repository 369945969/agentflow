# AgentFlow Tech Blue Glow Style Guide

**Style Overview**:
A refined dark glassmorphism theme featuring luminous sky blue accents with subtle glow effects, built on frosted-glass-like translucent containers with blur effects, thin luminous borders, and blue-tinted shadows over a soft multi-color gradient background, creating a futuristic, premium AI-powered tech interface aesthetic.

## Colors
### Primary Colors
  - **primary-base**: `text-[#3B9BFF]` or `bg-[#3B9BFF]`
  - **primary-lighter**: `text-[#5FB4FF]` or `bg-[#5FB4FF]`
  - **primary-darker**: `text-[#2A7FDB]` or `bg-[#2A7FDB]`
  - **primary-glow**: `shadow-[0_0_20px_rgba(59,155,255,0.4)]` - Luminous glow effect for emphasis

### Background Colors

#### Structural Backgrounds

Choose based on layout type:

**For Vertical Layout** (Top Header + Optional Side Panels):
- **bg-nav-primary**: `style="background: linear-gradient(180deg, rgba(15, 25, 45, 0.85) 0%, rgba(20, 35, 60, 0.85) 100%);"` - Top header with subtle gradient and glassmorphism
- **bg-nav-secondary**: `style="background: linear-gradient(180deg, rgba(20, 30, 50, 0.75) 0%, rgba(25, 40, 65, 0.75) 100%);"` - Inner Left sidebar (if present)
- **bg-page**: `style="background: radial-gradient(circle at 20% 30%, rgba(59, 155, 255, 0.15) 0%, transparent 50%), radial-gradient(circle at 80% 70%, rgba(135, 180, 255, 0.12) 0%, transparent 50%), #0F1928;"` - Page background with ambient orbs

**For Horizontal Layout** (Side Navigation + Optional Top Bar):
- **bg-nav-primary**: `style="background: linear-gradient(90deg, rgba(15, 25, 45, 0.85) 0%, rgba(20, 35, 60, 0.85) 100%);"` - Left main sidebar
- **bg-nav-secondary**: `style="background: linear-gradient(180deg, rgba(20, 30, 50, 0.75) 0%, rgba(25, 40, 65, 0.75) 100%);"` - Inner Top header (if present)
- **bg-page**: `style="background: radial-gradient(circle at 20% 30%, rgba(59, 155, 255, 0.15) 0%, transparent 50%), radial-gradient(circle at 80% 70%, rgba(135, 180, 255, 0.12) 0%, transparent 50%), #0F1928;"` - Page background

#### Container Backgrounds
For main content area. Glassmorphism effect achieved through backdrop blur and semi-transparent backgrounds.
- **bg-container-primary**: `bg-white/5 backdrop-blur-xl` - Primary glass containers
- **bg-container-secondary**: `bg-white/3 backdrop-blur-lg` - Secondary glass containers
- **bg-container-inset**: `bg-[#3B9BFF]/10 backdrop-blur-md` - Inset areas with primary tint
- **bg-container-inset-strong**: `bg-[#2A7FDB]/15 backdrop-blur-md` - Stronger inset emphasis

### Text Colors
- **color-text-primary**: `text-white/95`
- **color-text-secondary**: `text-white/70`
- **color-text-tertiary**: `text-white/50`
- **color-text-quaternary**: `text-white/30`
- **color-text-on-light-primary**: `text-[#0F1928]/90` - Text on light backgrounds
- **color-text-on-light-secondary**: `text-[#0F1928]/70` - Text on light backgrounds
- **color-text-link**: `text-[#5FB4FF]` - Links, text-only buttons, and clickable text

### Functional Colors
Use **sparingly** to maintain the refined dark aesthetic.
  - **color-success-default**: rgba(80, 200, 120, 0.25)
  - **color-success-light**: rgba(80, 200, 120, 0.15) - tag/label bg
  - **color-error-default**: rgba(255, 100, 100, 0.25) - alert banner bg
  - **color-error-light**: rgba(255, 100, 100, 0.15) - tag/label bg
  - **color-warning-default**: rgba(255, 180, 70, 0.25) - tag/label bg
  - **color-warning-light**: rgba(255, 180, 70, 0.15) - tag/label bg, alert banner bg
  - **color-function-default**: rgba(135, 180, 255, 0.3)
  - **color-function-light**: rgba(135, 180, 255, 0.15) - tag/label bg

### Accent Colors
  - Secondary palette for highlights and categorization. **Avoid overuse** to maintain tech aesthetic.
  - **accent-cyan-bright**: `text-[#00E5FF]` or `bg-[#00E5FF]`
  - **accent-periwinkle-soft**: `text-[#87B4FF]` or `bg-[#87B4FF]`

### Data Visualization Charts
For data visualization charts only. Harmonious with the dark glassmorphism theme.
  - Standard data colors: rgba(255, 255, 255, 0.15), rgba(255, 255, 255, 0.25), rgba(255, 255, 255, 0.4), rgba(255, 255, 255, 0.6), rgba(255, 255, 255, 0.8), rgba(255, 255, 255, 0.95)
  - Important data can use small amounts of: #3B9BFF, #5FB4FF, #00E5FF, #87B4FF

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
  - **Small**: 12px — Elements inside cards (e.g., icons, small buttons)
  - **Medium**: 16px — Standard containers and cards
  - **Large**: 24px — Large cards and panels
  - **Full**: full — Toggles, avatars, small tags

## Layout & Spacing
  - **Tight**: 12px - For closely related small internal elements, such as icons and text within buttons
  - **Compact**: 16px - For small gaps between small containers, such as a line of tags
  - **Standard**: 24px - For gaps between medium containers like list items
  - **Relaxed**: 32px - For gaps between large containers and sections
  - **Section**: 40px - For major section divisions

## Create Boundaries (contrast of surface color, borders, shadows)
Refined dark glassmorphism using frosted-glass effect with luminous borders and blue-tinted shadows.

### Borders
  - **Default**: 1px solid rgba(59, 155, 255, 0.3). Used for glass containers. `border border-[#3B9BFF]/30`
  - **Stronger**: 1px solid rgba(59, 155, 255, 0.5). Used for active or focused states. `border border-[#3B9BFF]/50`
  - **Luminous**: 1px solid rgba(95, 180, 255, 0.6). Used for premium emphasis. `border border-[#5FB4FF]/60`

### Dividers
  - `border-t` or `border-b` `border-white/10` for subtle separation
  - `border-t` or `border-b` `border-[#3B9BFF]/20` for tech-emphasized dividers

### Shadows & Effects
  - **Case 1 (subtle glass)**: `shadow-[0_4px_16px_rgba(0,0,0,0.3)]` - Basic glassmorphism depth
  - **Case 2 (moderate glass)**: `shadow-[0_8px_24px_rgba(0,0,0,0.4),0_0_12px_rgba(59,155,255,0.15)]` - Glass with blue ambient glow
  - **Case 3 (pronounced glass)**: `shadow-[0_12px_32px_rgba(0,0,0,0.5),0_0_20px_rgba(59,155,255,0.25)]` - Strong glass with luminous blue glow
  - **Case 4 (premium luminous)**: `shadow-[0_16px_48px_rgba(0,0,0,0.6),0_0_32px_rgba(59,155,255,0.4)]` - Maximum depth with strong blue luminescence

## Visual Emphasis for Containers
When containers (tags, cards, list items, rows) need visual emphasis to indicate priority, status, or category, use the following techniques:

| Technique | Implementation Notes | Best For | Avoid |
|-----------|---------------------|----------|-------|
| Background Tint | Slightly adjust opacity or add subtle color tint to glass surfaces | Gentle emphasis within glassmorphism aesthetic | Heavy opacity that breaks translucency |
| Border Highlight | Use luminous borders with blue glow for emphasis | Active/selected states, premium features | - |
| Glow/Shadow Effect | Blue-tinted glow shadows for tech aesthetic | Hover states, AI-powered features, key actions | Overuse that diminishes premium feel |
| Status Tag/Label | Add colored translucent tag/label inside container | Larger glass containers | - |
| Side Accent Bar | **Left edge only**, with luminous blue gradient | Small list items (e.g., side nav tabs), task cards | Large cards, rounded containers |

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
  <iconify-icon icon="lucide:cpu" class="text-base"></iconify-icon>
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
  - **Graphic**: Use a simple, relevant icon (e.g., a `workflow` icon for orchestration platform, a `network` icon for multi-agent systems).

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
  - Example 1 (glass button with glow):
    - button: flex items-center bg-white/5 backdrop-blur-xl border border-[#3B9BFF]/30 hover:border-[#5FB4FF]/50 hover:shadow-[0_0_20px_rgba(59,155,255,0.3)] transition
      - span(button copy): whitespace-nowrap
  - Example 2 (icon button with glass):
    - button: flex items-center justify-center bg-white/5 backdrop-blur-xl border border-[#3B9BFF]/30 hover:border-[#5FB4FF]/50 transition
      - icon

- **Tag Group (Filter Tags)** (Note: `overflow-x-auto` and `whitespace-nowrap` are required)
  - container(scrollable): flex overflow-x-auto [&::-webkit-scrollbar]:hidden
    - label (Tag item 1):
      - input: type="radio" name="tag1" class="sr-only peer" checked
      - div: bg-white/5 backdrop-blur-md text-white/70 border border-white/10 peer-checked:bg-[#3B9BFF]/20 peer-checked:border-[#3B9BFF]/50 peer-checked:text-white/95 hover:opacity-80 transition whitespace-nowrap

### Data Entry
- **Progress bars/Slider**: h-3 bg-white/5 backdrop-blur-md
- **Checkbox**
  - label:
    - input: type="checkbox" class="sr-only peer"
    - div: bg-white/5 backdrop-blur-md border border-[#3B9BFF]/30 rounded-md flex items-center justify-center peer-checked:bg-[#3B9BFF]/25 peer-checked:border-[#3B9BFF]/60 text-transparent peer-checked:text-[#5FB4FF]
      - svg(Checkmark): stroke="currentColor" stroke-width="4"
    - span(text)
- **Radio button**
  - label:
    - input: type="radio" name="option1" class="sr-only peer"
    - div: bg-white/5 backdrop-blur-md border border-[#3B9BFF]/30 rounded-full flex items-center justify-center peer-checked:bg-[#3B9BFF]/25 peer-checked:border-[#3B9BFF]/60 text-transparent peer-checked:text-[#5FB4FF]
      - svg(dot indicator): fill="currentColor"
    - span(text)
- **Switch/Toggle**
  - label:
    - div: relative
      - input: type="checkbox" class="sr-only peer"
      - div(Toggle track): w-14 h-7 bg-white/5 backdrop-blur-md border border-[#3B9BFF]/30 peer-checked:bg-[#3B9BFF]/25 peer-checked:border-[#3B9BFF]/60 transition
      - div(Toggle thumb): absolute top-1 left-1 w-5 h-5 bg-[#5FB4FF] rounded-full peer-checked:translate-x-7 peer-checked:shadow-[0_0_12px_rgba(95,180,255,0.6)] transition
    - span(text)

- **Select/Dropdown**
  - Select container: flex items-center bg-white/5 backdrop-blur-xl border border-[#3B9BFF]/30
    - text
    - Dropdown icon(square container): flex items-center justify-center bg-transparent
      - icon

### Container
- **Navigation Menu - horizontal**
    - Navigation with sections/grouping:
        - Nav Container: flex items-center justify-between w-full
        - Left Section: flex items-center gap-10
          - Menu Item: flex items-center gap-3
        - Right Section: flex items-center gap-4
          - Menu Item: flex items-center gap-3
          - Notification (if applicable): relative flex items-center justify-center w-12 h-12
            - notification-icon: w-6 h-6
            - badge (if has unread): absolute -top-1 -right-1 w-6 h-6 rounded-full flex items-center justify-center bg-[#3B9BFF] shadow-[0_0_12px_rgba(59,155,255,0.6)]
              - badge-count:
          - Avatar(if applicable): flex items-center
            - avatar-image: w-10 h-10 rounded-full border-2 border-[#3B9BFF]/40
            - dropdown-icon (if applicable): w-6 h-6

- **Card**
    - Example 1 (Vertical glass card with image and text):
        - Card: bg-white/5 backdrop-blur-xl border border-[#3B9BFF]/30 rounded-2xl flex flex-col p-4 gap-4 shadow-[0_8px_24px_rgba(0,0,0,0.4),0_0_12px_rgba(59,155,255,0.15)]
        - Image: rounded-lg w-full
        - Text area: flex flex-col gap-3
          - card-title: text-base font-semibold
          - card-subtitle: text-xs font-normal
    - Example 2 (Horizontal glass card with image and text):
        - Card: bg-white/5 backdrop-blur-xl border border-[#3B9BFF]/30 rounded-2xl flex gap-4 shadow-[0_8px_24px_rgba(0,0,0,0.4),0_0_12px_rgba(59,155,255,0.15)]
        - Image: rounded-lg h-full
        - Text area: flex flex-col gap-4
          - card-title: text-base font-semibold
          - card-subtitle: text-xs font-normal
    - Example 3 (Image-focused glass card):
        - Card: flex flex-col gap-4
        - Image: rounded-2xl w-full border border-[#3B9BFF]/20
        - Text area: flex flex-col gap-3
          - card-title: text-base font-semibold
          - card-subtitle: text-xs font-normal
    - Example 4 (Text-only glass cards with premium luminous effect):
        - Card: bg-white/5 backdrop-blur-xl border border-[#3B9BFF]/40 rounded-2xl flex shadow-[0_8px_24px_rgba(0,0,0,0.4),0_0_16px_rgba(59,155,255,0.2)]

## Additional Notes
- **Glassmorphism Effect**: Always combine `backdrop-blur-*` with semi-transparent backgrounds (`bg-white/5`, `bg-white/3`, etc.) to achieve the frosted glass effect
- **Luminous Borders**: Use borders with partial opacity from the primary blue color (`border-[#3B9BFF]/30`) to create the luminous outline effect
- **Blue-Tinted Shadows**: Combine standard shadows with blue-tinted glow shadows for depth and tech aesthetic (`shadow-[0_8px_24px_rgba(0,0,0,0.4),0_0_12px_rgba(59,155,255,0.15)]`)
- **Ambient Background**: The multi-color gradient background with radial orbs should remain subtle to maintain readability while adding atmospheric depth
- **Glow Effects**: Use sparingly on interactive elements, active states, and key features to maintain the premium feel without overwhelming the interface
- **Contrast**: Ensure sufficient contrast between text and glass surfaces for accessibility while maintaining the translucent aesthetic

<colors_extraction>
#3B9BFF
#5FB4FF
#2A7FDB
linear-gradient(180deg, rgba(15, 25, 45, 0.85) 0%, rgba(20, 35, 60, 0.85) 100%)
linear-gradient(180deg, rgba(20, 30, 50, 0.75) 0%, rgba(25, 40, 65, 0.75) 100%)
radial-gradient(circle at 20% 30%, rgba(59, 155, 255, 0.15) 0%, transparent 50%), radial-gradient(circle at 80% 70%, rgba(135, 180, 255, 0.12) 0%, transparent 50%), #0F1928
linear-gradient(90deg, rgba(15, 25, 45, 0.85) 0%, rgba(20, 35, 60, 0.85) 100%)
#FFFFFF0D
#FFFFFF08
#3B9BFF1A
#2A7FDB26
#FFFFFFF2
#FFFFFFB3
#FFFFFF80
#FFFFFF4D
#0F1928E6
#0F1928B3
#50C87840
#50C87826
#FF646440
#FF646426
#FFB44640
#FFB44626
#87B4FF4D
#87B4FF26
#00E5FF
#87B4FF
#FFFFFF26
#FFFFFF40
#FFFFFF66
#FFFFFF99
#FFFFFFCC
#FFFFFFF2
#3B9BFF4D
#3B9BFF80
#3B9BFF99
#00000040
#00000066
#00000080
#3B9BFF66
</colors_extraction>
