# Lua Configuration

Complete guide to configuring Envy using Lua.

## Configuration File

Envy reads configuration from:

| OS | Path |
|----|------|
| Linux | `~/.config/envy/config.lua` |
| macOS | `~/Library/Application Support/envy/config.lua` |
| Windows | `%APPDATA%\envy\config.lua` |

**Note:** The config file is optional. If it doesn't exist, Envy uses sensible defaults.

## Configuration Structure

The config file must return a Lua table:

```lua
return {
  -- Backend settings
  backend = {
    keys_path = "~/.envy/keys.json",
    lock_path = "~/.envy/.lock"
  },

  -- Keybindings
  keys = {
    quit = "q",
    yank = "y",
    edit = "e",
    -- ... more keys
  },

  -- Theme customization
  theme = {
    base = "#1e1e2e",
    text = "#cdd6f4",
    -- ... more colors
  }
}
```

## Backend Configuration

Control where Envy stores its data.

### Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `keys_path` | string | `"~/.envy/keys.json"` | Path to encrypted vault |
| `lock_path` | string | `"~/.envy/.lock"` | Path to lock file |

### Examples

```lua
-- Default configuration
backend = {
  keys_path = "~/.envy/keys.json",
  lock_path = "~/.envy/.lock"
}

-- Custom data directory
backend = {
  keys_path = "~/Dropbox/envy/keys.json",
  lock_path = "~/Dropbox/envy/.lock"
}

-- XDG Base Directory compliance
backend = {
  keys_path = os.getenv("XDG_DATA_HOME") .. "/envy/keys.json",
  lock_path = os.getenv("XDG_DATA_HOME") .. "/envy/.lock"
}
```

### Path Expansion

Envy automatically expands:
- `~` to your home directory
- Environment variables (using Lua's `os.getenv`)

### Built-in Helper Module

Envy provides an `envy` module in the Lua context:

```lua
-- Available in config.lua:
print(envy.home)           -- Home directory path
print(envy.os)             -- Operating system: "linux", "darwin", "windows"
print(envy.default_data_dir)    -- Default data directory
print(envy.default_config_dir)  -- Default config directory
print(envy.default_keys_path)   -- Default keys.json path

-- Path expansion helper
local expanded = envy.expand_path("~/myfolder")  -- "/home/user/myfolder"
```

### Advanced Backend Example

```lua
-- Platform-specific paths
local keys_path, lock_path

if envy.os == "windows" then
  keys_path = envy.home .. \\\"AppData\\\Local\\envy\\keys.json\"
  lock_path = envy.home .. \\\"AppData\\\Local\\envy\\.lock\"
elseif envy.os == "darwin" then
  keys_path = envy.home .. "/Library/Application Support/envy/keys.json"
  lock_path = envy.home .. "/Library/Application Support/envy/.lock"
else
  -- Linux and others
  keys_path = envy.home .. ".envy/keys.json"
  lock_path = envy.home .. ".envy/.lock"
end

return {
  backend = {
    keys_path = keys_path,
    lock_path = lock_path
  }
}
```

## Keybindings Configuration

Customize all keyboard shortcuts in the TUI.

### Available Keybinding Options

| Option | Default | Description |
|--------|---------|-------------|
| **Navigation** |||
| `up` | `"up"` | Move up (arrow key name) |
| `down` | `"down"` | Move down |
| `left` | `"left"` | Move left |
| `right` | `"right"` | Move right |
| `vim_up` | `"k"` | Vim-style up |
| `vim_down` | `"j"` | Vim-style down |
| `vim_left` | `"h"` | Vim-style left |
| `vim_right` | `"l"` | Vim-style right |
| **Actions** |||
| `enter` | `"enter"` | Confirm/select |
| `back` | `"esc"` | Go back/cancel |
| `quit` | `"q"` | Quit application |
| `search` | `"i"` | Activate search |
| `yank` | `"y"` | Copy to clipboard |
| `create` | `"ctrl+n"` | Create new project |
| `edit` | `"e"` | Edit key value |
| `edit_project` | `"E"` | Edit project |
| `delete` | `"d"` | Delete item |
| `save` | `"ctrl+s"` | Save changes |
| `add` | `"ctrl+a"` | Add key/item |
| `history` | `"H"` | View history |
| **Form Navigation** |||
| `tab` | `"tab"` | Next field |
| `shift_tab` | `"shift+tab"` | Previous field |
| `space` | `" "` | Space/activate |
| **Special** |||
| `force_quit` | `"ctrl+c"` | Force quit |

### Special Key Names

```lua
-- Modifier combinations
"ctrl+n"      -- Control + n
"ctrl+s"      -- Control + s
"ctrl+a"      -- Control + a
"ctrl+c"      -- Control + c
"shift+tab"   -- Shift + Tab

-- Special keys
"enter"       -- Enter/Return
"esc"         -- Escape
"tab"         -- Tab
"space"       -- Space (or use " ")
"up"          -- Up arrow
"down"        -- Down arrow
"left"        -- Left arrow
"right"       -- Right arrow

-- Letters and numbers
"a" through "z"
"A" through "Z"  (case matters!)
"0" through "9"
```

### Keybinding Examples

```lua
-- Default keybindings
keys = {
  -- Navigation
  up = "up",
  down = "down",
  left = "left",
  right = "right",

  -- Vim navigation (alternative)
  vim_up = "k",
  vim_down = "j",
  vim_left = "h",
  vim_right = "l",

  -- Actions
  enter = "enter",
  back = "esc",
  quit = "q",
  search = "i",
  yank = "y",
  create = "ctrl+n",
  edit = "e",
  edit_project = "E",
  delete = "d",
  save = "ctrl+s",
  add = "ctrl+a",
  history = "H",

  -- Form navigation
  tab = "tab",
  shift_tab = "shift+tab",
  space = " ",

  -- Special
  force_quit = "ctrl+c"
}
```

### Emacs-Style Keybindings

```lua
-- For Emacs enthusiasts
keys = {
  -- Navigation (Emacs-style)
  vim_up = "p",      -- Previous line
  vim_down = "n",    -- Next line
  vim_left = "b",    -- Backward char
  vim_right = "f",   -- Forward char

  -- Actions
  quit = "ctrl+x",
  save = "ctrl+x ctrl+s",
  search = "ctrl+s",
  yank = "ctrl+y",
  delete = "ctrl+d",

  -- Keep defaults for others
  enter = "enter",
  back = "esc",
  create = "ctrl+n",
  edit = "e",
  edit_project = "E",
  add = "ctrl+a",
  history = "H",
  tab = "tab",
  shift_tab = "shift+tab",
  space = " ",
  force_quit = "ctrl+c"
}
```

### Vim-Only Mode

```lua
-- Disable arrow keys, use Vim keys exclusively
keys = {
  -- Disable arrows
  up = "",
  down = "",
  left = "",
  right = "",

  -- Vim navigation
  vim_up = "k",
  vim_down = "j",
  vim_left = "h",
  vim_right = "l",

  -- Actions
  quit = "q",
  search = "/",      -- Vim-style search
  yank = "y",
  create = "n",      -- 'n' for new
  edit = "e",
  edit_project = "E",
  delete = "d",
  save = ":w",       -- Vim-style save
  add = "a",
  history = "H",
  enter = "enter",
  back = "esc",
  tab = "tab",
  shift_tab = "shift+tab",
  space = " ",
  force_quit = "ctrl+c"
}
```

## Theme Configuration

Customize the visual appearance of Envy.

### Color Options

| Option | Default | Description |
|--------|---------|-------------|
| **Base Colors** |||
| `base` | `"#1e1e2e"` | Background color |
| `text` | `"#cdd6f4"` | Primary text color |
| `accent` | `"#cba6f7"` | Accent/highlight color |
| `surface0` | `"#313244"` | Surface background |
| `surface1` | `"#45475a"` | Elevated surface |
| `overlay0` | `"#6c7086"` | Dimmed text/overlay |
| **Semantic Colors** |||
| `success` | `"#a6e3a1"` | Success/green |
| `warning` | `"#f9e2af"` | Warning/yellow |
| `error` | `"#f38ba8"` | Error/red |
| **Environment Badges** |||
| `prod_bg` | `"#f38ba8"` | Production badge background |
| `dev_bg` | `"#a6e3a1"` | Development badge background |
| `stage_bg` | `"#f9e2af"` | Staging badge background |
| **History Section** |||
| `current_bg` | `"#a6e3a1"` | Current value background |
| `previous_bg` | `"#f9e2af"` | Previous values background |
| **Layout** |||
| `grid_cols` | `3` | Number of grid columns |
| `grid_visible_rows` | `2` | Visible rows in grid |
| `card_width` | `38` | Project card width |
| `card_height` | `9` | Project card height |

### Default Theme (Catppuccin Mocha)

```lua
theme = {
  -- Base colors
  base = "#1e1e2e",
  text = "#cdd6f4",
  accent = "#cba6f7",
  surface0 = "#313244",
  surface1 = "#45475a",
  overlay0 = "#6c7086",

  -- Semantic colors
  success = "#a6e3a1",
  warning = "#f9e2af",
  error = "#f38ba8",

  -- Environment badges
  prod_bg = "#f38ba8",  -- Red for production
  dev_bg = "#a6e3a1",   -- Green for development
  stage_bg = "#f9e2af", -- Yellow for staging

  -- History colors
  current_bg = "#a6e3a1",
  previous_bg = "#f9e2af",

  -- Layout
  grid_cols = 3,
  grid_visible_rows = 2,
  card_width = 38,
  card_height = 9
}
```

### Theme Examples

#### Gruvbox Dark

```lua
theme = {
  -- Base colors (Gruvbox Dark)
  base = "#282828",
  text = "#ebdbb2",
  accent = "#d79921",
  surface0 = "#3c3836",
  surface1 = "#504945",
  overlay0 = "#928374",

  -- Semantic colors
  success = "#b8bb26",
  warning = "#fabd2f",
  error = "#fb4934",

  -- Environment badges
  prod_bg = "#fb4934",
  dev_bg = "#b8bb26",
  stage_bg = "#fabd2f",

  -- History colors
  current_bg = "#b8bb26",
  previous_bg = "#fabd2f",

  -- Layout
  grid_cols = 3,
  grid_visible_rows = 2,
  card_width = 38,
  card_height = 9
}
```

#### Dracula

```lua
theme = {
  -- Base colors (Dracula)
  base = "#282a36",
  text = "#f8f8f2",
  accent = "#bd93f9",
  surface0 = "#44475a",
  surface1 = "#6272a4",
  overlay0 = "#6272a4",

  -- Semantic colors
  success = "#50fa7b",
  warning = "#f1fa8c",
  error = "#ff5555",

  -- Environment badges
  prod_bg = "#ff5555",
  dev_bg = "#50fa7b",
  stage_bg = "#f1fa8c",

  -- History colors
  current_bg = "#50fa7b",
  previous_bg = "#f1fa8c",

  -- Layout
  grid_cols = 3,
  grid_visible_rows = 2,
  card_width = 38,
  card_height = 9
}
```

#### Solarized Dark

```lua
theme = {
  -- Base colors (Solarized Dark)
  base = "#002b36",
  text = "#839496",
  accent = "#268bd2",
  surface0 = "#073642",
  surface1 = "#586e75",
  overlay0 = "#657b83",

  -- Semantic colors
  success = "#859900",
  warning = "#b58900",
  error = "#dc322f",

  -- Environment badges
  prod_bg = "#dc322f",
  dev_bg = "#859900",
  stage_bg = "#b58900",

  -- History colors
  current_bg = "#859900",
  previous_bg = "#b58900",

  -- Layout
  grid_cols = 3,
  grid_visible_rows = 2,
  card_width = 38,
  card_height = 9
}
```

#### One Dark (Atom)

```lua
theme = {
  -- Base colors (One Dark)
  base = "#282c34",
  text = "#abb2bf",
  accent = "#61afef",
  surface0 = "#3e4451",
  surface1 = "#4b5263",
  overlay0 = "#5c6370",

  -- Semantic colors
  success = "#98c379",
  warning = "#e5c07b",
  error = "#e06c75",

  -- Environment badges
  prod_bg = "#e06c75",
  dev_bg = "#98c379",
  stage_bg = "#e5c07b",

  -- History colors
  current_bg = "#98c379",
  previous_bg = "#e5c07b",

  -- Layout
  grid_cols = 3,
  grid_visible_rows = 2,
  card_width = 38,
  card_height = 9
}
```

#### Nord

```lua
theme = {
  -- Base colors (Nord)
  base = "#2e3440",
  text = "#d8dee9",
  accent = "#88c0d0",
  surface0 = "#3b4252",
  surface1 = "#434c5e",
  overlay0 = "#4c566a",

  -- Semantic colors
  success = "#a3be8c",
  warning = "#ebcb8b",
  error = "#bf616a",

  -- Environment badges
  prod_bg = "#bf616a",
  dev_bg = "#a3be8c",
  stage_bg = "#ebcb8b",

  -- History colors
  current_bg = "#a3be8c",
  previous_bg = "#ebcb8b",

  -- Layout
  grid_cols = 3,
  grid_visible_rows = 2,
  card_width = 38,
  card_height = 9
}
```

### High Contrast Theme

For accessibility:

```lua
theme = {
  -- High contrast black/white
  base = "#000000",
  text = "#ffffff",
  accent = "#ffff00",
  surface0 = "#333333",
  surface1 = "#666666",
  overlay0 = "#999999",

  -- Bright semantic colors
  success = "#00ff00",
  warning = "#ffff00",
  error = "#ff0000",

  -- Environment badges
  prod_bg = "#ff0000",
  dev_bg = "#00ff00",
  stage_bg = "#ffff00",

  -- History colors
  current_bg = "#00ff00",
  previous_bg = "#ffff00",

  -- Larger cards for readability
  grid_cols = 2,
  grid_visible_rows = 2,
  card_width = 50,
  card_height = 12
}
```

### Compact Layout

For smaller screens:

```lua
theme = {
  -- Default colors
  base = "#1e1e2e",
  text = "#cdd6f4",
  accent = "#cba6f7",
  surface0 = "#313244",
  surface1 = "#45475a",
  overlay0 = "#6c7086",

  success = "#a6e3a1",
  warning = "#f9e2af",
  error = "#f38ba8",

  prod_bg = "#f38ba8",
  dev_bg = "#a6e3a1",
  stage_bg = "#f9e2af",

  current_bg = "#a6e3a1",
  previous_bg = "#f9e2af",

  -- Compact layout
  grid_cols = 4,
  grid_visible_rows = 3,
  card_width = 30,
  card_height = 7
}
```

## Complete Configuration Examples

### Minimal Config

```lua
-- Only change what you need
return {
  keys = {
    quit = "ctrl+q"  -- Only customize quit key
  }
}
```

### Developer Config

```lua
-- Optimized for development workflow
return {
  backend = {
    keys_path = "~/.envy/keys.json",
    lock_path = "~/.envy/.lock"
  },

  keys = {
    -- Navigation
    up = "up",
    down = "down",
    left = "left",
    right = "right",
    vim_up = "k",
    vim_down = "j",
    vim_left = "h",
    vim_right = "l",

    -- Quick actions
    enter = "enter",
    back = "esc",
    quit = "q",
    search = "/",      -- Vim-style
    yank = "y",
    create = "n",      -- Quick 'n' for new
    edit = "e",
    edit_project = "E",
    delete = "d",
    save = "ctrl+s",
    add = "a",
    history = "h",     -- Lowercase h

    tab = "tab",
    shift_tab = "shift+tab",
    space = " ",
    force_quit = "ctrl+c"
  },

  theme = {
    base = "#1e1e2e",
    text = "#cdd6f4",
    accent = "#89b4fa",  -- Blue accent for calmer look
    surface0 = "#313244",
    surface1 = "#45475a",
    overlay0 = "#6c7086",

    success = "#a6e3a1",
    warning = "#f9e2af",
    error = "#f38ba8",

    prod_bg = "#f38ba8",
    dev_bg = "#a6e3a1",
    stage_bg = "#f9e2af",

    current_bg = "#a6e3a1",
    previous_bg = "#f9e2af",

    grid_cols = 3,
    grid_visible_rows = 2,
    card_width = 38,
    card_height = 9
  }
}
```

### Team Config Template

```lua
-- Standardized team configuration
-- Share this file with team members

return {
  backend = {
    -- Use XDG directories for consistency
    keys_path = (os.getenv("XDG_DATA_HOME") or (envy.home .. "/.local/share")) .. "/envy/keys.json",
    lock_path = (os.getenv("XDG_DATA_HOME") or (envy.home .. "/.local/share")) .. "/envy/.lock"
  },

  keys = {
    -- Standard navigation
    up = "up",
    down = "down",
    left = "left",
    right = "right",
    vim_up = "k",
    vim_down = "j",
    vim_left = "h",
    vim_right = "l",

    -- Standard actions
    enter = "enter",
    back = "esc",
    quit = "q",
    search = "i",
    yank = "y",
    create = "ctrl+n",
    edit = "e",
    edit_project = "E",
    delete = "d",
    save = "ctrl+s",
    add = "ctrl+a",
    history = "H",

    tab = "tab",
    shift_tab = "shift+tab",
    space = " ",
    force_quit = "ctrl+c"
  },

  theme = {
    -- Company branding colors (customize these)
    base = "#1e1e2e",
    text = "#cdd6f4",
    accent = "#cba6f7",
    surface0 = "#313244",
    surface1 = "#45475a",
    overlay0 = "#6c7086",

    success = "#a6e3a1",
    warning = "#f9e2af",
    error = "#f38ba8",

    prod_bg = "#f38ba8",
    dev_bg = "#a6e3a1",
    stage_bg = "#f9e2af",

    current_bg = "#a6e3a1",
    previous_bg = "#f9e2af",

    grid_cols = 3,
    grid_visible_rows = 2,
    card_width = 38,
    card_height = 9
  }
}
```

## Configuration Reload

Envy reads the configuration file:
- At startup
- When explicitly reloaded (if future feature)

Changes take effect on the next Envy launch.

## Troubleshooting Config

### Config Not Loading

1. Check file path:
   ```bash
   # Linux/macOS
   ls -la ~/.config/envy/config.lua

   # macOS (alternate)
   ls -la ~/Library/Application\ Support/envy/config.lua
   ```

2. Validate Lua syntax:
   ```bash
   lua -c ~/.config/envy/config.lua
   ```

3. Check Envy is reading it:
   - Syntax errors are silently ignored
   - Envy uses defaults if config is invalid

### Common Mistakes

```lua
-- WRONG: Not returning a table
local config = {
  keys = { quit = "q" }
}

-- RIGHT: Must return the table
return {
  keys = { quit = "q" }
}
```

```lua
-- WRONG: Invalid key name
keys = {
  quit = "ctrl+q"  -- '+' not valid in key name
}

-- RIGHT: Valid key name
keys = {
  quit = "ctrl+q"  -- Actually this is fine, just no spaces
}
```

```lua
-- WRONG: Missing quotes around values
theme = {
  base = #1e1e2e  -- Syntax error
}

-- RIGHT: Properly quoted
theme = {
  base = "#1e1e2e"
}
```

---

**Next:** Learn about [Theme Gallery](./themes.md) for more visual customization options.
