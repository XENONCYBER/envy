# Themes

Theme gallery and customization guide for Envy.

## Default Theme

Envy ships with **Catppuccin Mocha** as the default theme - a soft, pastel color scheme that's easy on the eyes during long coding sessions.

## Built-in Theme Gallery

### Catppuccin Mocha (Default)

A soft, muted color palette inspired by coffee.

```lua
theme = {
  base = "#1e1e2e",      -- Dark background
  text = "#cdd6f4",      -- Soft white text
  accent = "#cba6f7",    -- Lavender purple
  surface0 = "#313244",  -- Elevated surfaces
  surface1 = "#45475a",  -- Higher elevation
  overlay0 = "#6c7086",  -- Dimmed elements

  success = "#a6e3a1",   -- Green
  warning = "#f9e2af",   -- Yellow
  error = "#f38ba8",     -- Red

  prod_bg = "#f38ba8",   -- Red for production
  dev_bg = "#a6e3a1",    -- Green for dev
  stage_bg = "#f9e2af",  -- Yellow for staging

  current_bg = "#a6e3a1",
  previous_bg = "#f9e2af"
}
```

**Characteristics:**
- Low contrast for reduced eye strain
- Pastel accents
- Popular among developers
- Good for dark rooms

### Catppuccin Latte

Light variant for bright environments.

```lua
theme = {
  base = "#eff1f5",
  text = "#4c4f69",
  accent = "#8839ef",
  surface0 = "#ccd0da",
  surface1 = "#bcc0cc",
  overlay0 = "#9ca0b0",

  success = "#40a02b",
  warning = "#df8e1d",
  error = "#d20f39",

  prod_bg = "#d20f39",
  dev_bg = "#40a02b",
  stage_bg = "#df8e1d",

  current_bg = "#40a02b",
  previous_bg = "#df8e1d"
}
```

### Gruvbox Dark

Retro groove color scheme.

```lua
theme = {
  base = "#282828",
  text = "#ebdbb2",
  accent = "#d79921",
  surface0 = "#3c3836",
  surface1 = "#504945",
  overlay0 = "#928374",

  success = "#b8bb26",
  warning = "#fabd2f",
  error = "#fb4934",

  prod_bg = "#fb4934",
  dev_bg = "#b8bb26",
  stage_bg = "#fabd2f",

  current_bg = "#b8bb26",
  previous_bg = "#fabd2f"
}
```

**Characteristics:**
- Retro aesthetic
- Warm colors
- Easy on the eyes
- Popular in Vim/Neovim community

### Gruvbox Light

```lua
theme = {
  base = "#fbf1c7",
  text = "#3c3836",
  accent = "#b57614",
  surface0 = "#ebdbb2",
  surface1 = "#d5c4a1",
  overlay0 = "#bdae93",

  success = "#79740e",
  warning = "#b57614",
  error = "#9d0006",

  prod_bg = "#9d0006",
  dev_bg = "#79740e",
  stage_bg = "#b57614",

  current_bg = "#79740e",
  previous_bg = "#b57614"
}
```

### Dracula

Dark theme with vibrant colors.

```lua
theme = {
  base = "#282a36",
  text = "#f8f8f2",
  accent = "#bd93f9",
  surface0 = "#44475a",
  surface1 = "#6272a4",
  overlay0 = "#6272a4",

  success = "#50fa7b",
  warning = "#f1fa8c",
  error = "#ff5555",

  prod_bg = "#ff5555",
  dev_bg = "#50fa7b",
  stage_bg = "#f1fa8c",

  current_bg = "#50fa7b",
  previous_bg = "#f1fa8c"
}
```

**Characteristics:**
- Vibrant, saturated colors
- High contrast
- Easy to distinguish elements
- Popular in VS Code

### Solarized Dark

Precision colors for machines and people.

```lua
theme = {
  base = "#002b36",
  text = "#839496",
  accent = "#268bd2",
  surface0 = "#073642",
  surface1 = "#586e75",
  overlay0 = "#657b83",

  success = "#859900",
  warning = "#b58900",
  error = "#dc322f",

  prod_bg = "#dc322f",
  dev_bg = "#859900",
  stage_bg = "#b58900",

  current_bg = "#859900",
  previous_bg = "#b58900"
}
```

**Characteristics:**
- Scientifically designed
- Consistent perceived brightness
- Works in varying light conditions
- Low eye strain

### Solarized Light

```lua
theme = {
  base = "#fdf6e3",
  text = "#657b83",
  accent = "#268bd2",
  surface0 = "#eee8d5",
  surface1 = "#93a1a1",
  overlay0 = "#839496",

  success = "#859900",
  warning = "#b58900",
  error = "#dc322f",

  prod_bg = "#dc322f",
  dev_bg = "#859900",
  stage_bg = "#b58900",

  current_bg = "#859900",
  previous_bg = "#b58900"
}
```

### Nord

Polar-inspired color palette.

```lua
theme = {
  base = "#2e3440",
  text = "#d8dee9",
  accent = "#88c0d0",
  surface0 = "#3b4252",
  surface1 = "#434c5e",
  overlay0 = "#4c566a",

  success = "#a3be8c",
  warning = "#ebcb8b",
  error = "#bf616a",

  prod_bg = "#bf616a",
  dev_bg = "#a3be8c",
  stage_bg = "#ebcb8b",

  current_bg = "#a3be8c",
  previous_bg = "#ebcb8b"
}
```

**Characteristics:**
- Arctic, north-bluish colors
- Clean and uncluttered
- Good for long sessions
- Matches Nordic aesthetic

### One Dark (Atom)

GitHub Atom's iconic theme.

```lua
theme = {
  base = "#282c34",
  text = "#abb2bf",
  accent = "#61afef",
  surface0 = "#3e4451",
  surface1 = "#4b5263",
  overlay0 = "#5c6370",

  success = "#98c379",
  warning = "#e5c07b",
  error = "#e06c75",

  prod_bg = "#e06c75",
  dev_bg = "#98c379",
  stage_bg = "#e5c07b",

  current_bg = "#98c379",
  previous_bg = "#e5c07b"
}
```

### Tokyo Night

Dark, stormy blues.

```lua
theme = {
  base = "#1a1b26",
  text = "#a9b1d6",
  accent = "#7aa2f7",
  surface0 = "#24283b",
  surface1 = "#414868",
  overlay0 = "#565f89",

  success = "#73daca",
  warning = "#e0af68",
  error = "#f7768e",

  prod_bg = "#f7768e",
  dev_bg = "#73daca",
  stage_bg = "#e0af68",

  current_bg = "#73daca",
  previous_bg = "#e0af68"
}
```

### Monokai Pro

Professional theme based on Monokai.

```lua
theme = {
  base = "#2d2a2e",
  text = "#fcfcfa",
  accent = "#ffd866",
  surface0 = "#403e41",
  surface1 = "#5b595c",
  overlay0 = "#727072",

  success = "#a9dc76",
  warning = "#ffd866",
  error = "#ff6188",

  prod_bg = "#ff6188",
  dev_bg = "#a9dc76",
  stage_bg = "#ffd866",

  current_bg = "#a9dc76",
  previous_bg = "#ffd866"
}
```

## Special Purpose Themes

### High Contrast

For accessibility and visibility.

```lua
theme = {
  base = "#000000",
  text = "#ffffff",
  accent = "#ffff00",
  surface0 = "#333333",
  surface1 = "#666666",
  overlay0 = "#999999",

  success = "#00ff00",
  warning = "#ffff00",
  error = "#ff0000",

  prod_bg = "#ff0000",
  dev_bg = "#00ff00",
  stage_bg = "#ffff00",

  current_bg = "#00ff00",
  previous_bg = "#ffff00",

  -- Larger elements for visibility
  grid_cols = 2,
  card_width = 50,
  card_height = 12
}
```

### OLED Dark

Pure black for OLED displays (saves battery).

```lua
theme = {
  base = "#000000",      -- Pure black
  text = "#e0e0e0",
  accent = "#bb86fc",
  surface0 = "#121212",
  surface1 = "#1e1e1e",
  overlay0 = "#666666",

  success = "#03dac6",
  warning = "#ffde03",
  error = "#cf6679",

  prod_bg = "#cf6679",
  dev_bg = "#03dac6",
  stage_bg = "#ffde03",

  current_bg = "#03dac6",
  previous_bg = "#ffde03"
}
```

### Presentation Mode

High visibility for screen sharing.

```lua
theme = {
  base = "#1e1e2e",
  text = "#ffffff",      -- Brighter text
  accent = "#ff6b6b",    -- Vibrant accent
  surface0 = "#2d2d44",
  surface1 = "#3d3d5c",
  overlay0 = "#7a7a9a",

  success = "#51cf66",
  warning = "#ffd43b",
  error = "#ff6b6b",

  prod_bg = "#ff6b6b",
  dev_bg = "#51cf66",
  stage_bg = "#ffd43b",

  current_bg = "#51cf66",
  previous_bg = "#ffd43b",

  -- Larger cards
  grid_cols = 2,
  card_width = 45,
  card_height = 11
}
```

### Redacted Mode

Hide values for demonstrations.

```lua
theme = {
  base = "#1a1a1a",
  text = "#666666",      -- Muted text
  accent = "#444444",    -- Subdued accent
  surface0 = "#222222",
  surface1 = "#2a2a2a",
  overlay0 = "#333333",

  success = "#2d4a2d",   -- Dark green
  warning = "#4a4a2d",   -- Dark yellow
  error = "#4a2d2d",     -- Dark red

  prod_bg = "#4a2d2d",
  dev_bg = "#2d4a2d",
  stage_bg = "#4a4a2d",

  current_bg = "#2d4a2d",
  previous_bg = "#4a4a2d"
}
```

## Color Guidelines

### Choosing Accent Colors

The `accent` color is used for:
- Selected items
- Borders
- Active elements
- Logo

**Recommendations:**
- Use a color that stands out from base
- Avoid colors too similar to error/warning
- Consider color blindness (avoid red/green combos)

### Environment Badge Colors

**Critical:** Users rely on these colors to identify environment risk:

- **Prod = Red**: Highest risk, production systems
- **Dev = Green**: Safe to experiment
- **Stage = Yellow**: Caution, pre-production

Keep these consistent even with custom themes!

### Semantic Colors

- **Success (green)**: Operations completed, current values
- **Warning (yellow)**: Previous values, caution
- **Error (red)**: Errors, production environment

## Creating Custom Themes

### Step 1: Choose a Base Palette

Start with a known palette:
- Material Design colors
- Tailwind CSS palette
- Open Color
- Your brand colors

### Step 2: Assign Colors

```lua
theme = {
  -- Darkest to lightest
  base = "#0f0f0f",      -- Background
  surface0 = "#1a1a1a",  -- Cards
  surface1 = "#252525",  -- Elevated
  overlay0 = "#666666",  -- Muted text
  text = "#e0e0e0",      -- Primary text

  -- Accent (your brand color)
  accent = "#ff6b6b",

  -- Semantic (standard)
  success = "#4caf50",
  warning = "#ff9800",
  error = "#f44336",

  -- Environment (keep consistent!)
  prod_bg = "#f44336",
  dev_bg = "#4caf50",
  stage_bg = "#ff9800",

  -- History (same as semantic)
  current_bg = "#4caf50",
  previous_bg = "#ff9800",

  -- Layout
  grid_cols = 3,
  grid_visible_rows = 2,
  card_width = 38,
  card_height = 9
}
```

### Step 3: Test Contrast

Ensure sufficient contrast ratios:
- Text on base: 7:1 minimum
- Accent on base: 4.5:1 minimum
- Semantic colors should be distinguishable

### Step 4: Iterate

Test in real usage:
- Long sessions
- Different lighting
- Various terminal emulators

## Sharing Themes

Share your custom themes:

1. Create a gist or repo
2. Include screenshots
3. Document the inspiration
4. Share on GitHub Discussions

## Theme Switching

Currently, Envy loads one theme at startup. To switch:

1. Edit `~/.config/envy/config.lua`
2. Change theme values
3. Restart Envy

Future versions may support:
- Multiple theme profiles
- Automatic light/dark switching
- Time-based switching

---

**Next:** Learn about customizing [Keybindings](./keybindings.md).
