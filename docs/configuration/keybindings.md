# Keybindings

Complete guide to customizing keyboard shortcuts in Envy.

## Overview

Envy supports full keybinding customization through Lua configuration. You can:
- Change any keyboard shortcut
- Add Vim-style navigation
- Create Emacs-style bindings
- Disable keys you don't use

## Default Keybindings

### Navigation

| Action | Default | Vim Alternative |
|--------|---------|-----------------|
| Move Up | `↑` | `k` |
| Move Down | `↓` | `j` |
| Move Left | `←` | `h` |
| Move Right | `→` | `l` |

### Actions

| Action | Default | Context |
|--------|---------|---------|
| Select/Confirm | `Enter` | Everywhere |
| Go Back | `Esc` | Detail view, forms |
| Quit | `q` | Grid view |
| Search | `i` | Grid view |
| Copy | `y` | Detail view |
| Edit Key | `e` | Detail view |
| Edit Project | `E` | Detail view |
| View History | `H` | Detail view |
| Delete | `shift+d` | Grid (project), Detail (key) |
| Create Project | `shift+n` | Grid view |
| Add Key | `shift+a` | Create/Edit forms |
| Save | `shift+s` | Forms |
| Force Quit | `Ctrl+c` | Everywhere |

### Form Navigation

| Action | Default |
|--------|---------|
| Next Field | `Tab` |
| Previous Field | `Shift+Tab` |
| Toggle/Activate | `Space` |

## Keybinding Syntax

### Single Keys

```lua
keys = {
  quit = "q",
  yank = "y",
  edit = "e"
}
```

### Special Keys

```lua
keys = {
  enter = "enter",
  back = "esc",
  tab = "tab",
  space = " "
}
```

### Arrow Keys

```lua
keys = {
  up = "up",
  down = "down",
  left = "left",
  right = "right"
}
```

### Modifier Combinations

```lua
keys = {
  create = "shift+n",
  save = "shift+s",
  add = "shift+a",
  delete = "shift+d",
  force_quit = "ctrl+c",
  shift_tab = "shift+tab"
}
```

**Valid modifiers:**
- `ctrl+` - Control key
- `shift+` - Shift key
- `alt+` - Alt/Option key (limited support)

**Note:** Some terminals intercept certain combinations (like `Ctrl+s`).

### Case Sensitivity

Case matters for letters:

```lua
keys = {
  edit = "e",      -- Lowercase e
  edit_project = "E"  -- Uppercase E (Shift+e)
}
```

## Configuration Examples

### Vim-Only Mode

For Vim purists who don't use arrow keys:

```lua
return {
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

    -- Vim-style actions
    quit = "q",
    search = "/",
    yank = "y",
    create = "n",
    edit = "e",
    edit_project = "E",
    delete = "D",
    save = "S",
    add = "A",
    history = "H",

    -- Keep these unchanged
    enter = "enter",
    back = "esc",
    tab = "tab",
    shift_tab = "shift+tab",
    space = " ",
    force_quit = "ctrl+c"
  }
}
```

### Emacs Style

For Emacs enthusiasts:

```lua
return {
  keys = {
    -- Emacs navigation
    vim_up = "p",      -- Previous line
    vim_down = "n",    -- Next line
    vim_left = "b",    -- Backward
    vim_right = "f",   -- Forward

    -- Emacs actions
    search = "ctrl+s",
    yank = "ctrl+y",
    save = "ctrl+x ctrl+s",
    delete = "ctrl+d",
    quit = "ctrl+x",

    -- Standard bindings
    enter = "enter",
    back = "esc",
    create = "N",
    edit = "e",
    edit_project = "E",
    add = "A",
    history = "H",
    tab = "tab",
    shift_tab = "shift+tab",
    space = " ",
    force_quit = "ctrl+c"
  }
}
```

### Minimal Bindings

Only essential keys:

```lua
return {
  keys = {
    -- Navigation (arrows only)
    up = "up",
    down = "down",
    left = "left",
    right = "right",

    -- Minimal actions
    enter = "enter",
    back = "esc",
    quit = "q",
    yank = "c",

    -- No Vim navigation
    vim_up = "",
    vim_down = "",
    vim_left = "",
    vim_right = ""
  }
}
```

### Two-Handed Layout

Optimized for touch typing:

```lua
return {
  keys = {
    -- Navigation (left hand)
    up = "up",
    down = "down",
    left = "left",
    right = "right",
    vim_up = "e",     -- Left hand home row
    vim_down = "d",
    vim_left = "s",
    vim_right = "f",

    -- Actions (right hand)
    quit = "p",
    yank = "o",       -- Easy to reach
    edit = "i",
    delete = "u",
    search = "l",
    create = "ctrl+n",
    edit_project = "ctrl+e",
    add = "ctrl+a",
    save = "ctrl+s",
    history = "ctrl+h",

    -- Keep standard
    enter = "enter",
    back = "esc",
    tab = "tab",
    shift_tab = "shift+tab",
    space = " ",
    force_quit = "ctrl+c"
  }
}
```

### Safe Mode

Prevent accidental deletions:

```lua
return {
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
    edit = "e",
    edit_project = "E",
    history = "H",  --Case sensitive(Shift+h)

    -- Safer delete - requires Ctrl
    delete = "ctrl+d",

    -- Create and add
    create = "N",
    add = "A",
    save = "S",

    tab = "tab",
    shift_tab = "shift+tab",
    space = " ",
    force_quit = "ctrl+c"
  }
}
```

### macOS Optimized

Better for Mac keyboards:

```lua
return {
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

    -- macOS-friendly shortcuts
    quit = "q",
    yank = "y",
    edit = "e",
    edit_project = "E",
    delete = "D",
    search = "i",
    history = "H",

    -- Cmd key would be better but limited terminal support
    -- Use Shift instead
    create = "N",
    save = "S",
    add = "A",

    enter = "enter",
    back = "esc",
    tab = "tab",
    shift_tab = "shift+tab",
    space = " ",
    force_quit = "ctrl+c"
  }
}
```

## Context-Specific Bindings

### Grid View Only

These only work in the project grid:

```lua
keys = {
  create = "N",  -- Create new project
  search = "i",       -- Activate search
  delete = "D"        -- Delete project (with confirmation)
}
```

### Detail View Only

These only work when viewing a project's secrets:

```lua
keys = {
  yank = "y",          -- Copy selected secret
  edit = "e",          -- Edit value of key(sidebar)
  edit_project = "E",  -- Edit project (full view)
  history = "H",       -- View history (sidebar)
  delete = "d"         -- Delete key (with confirmation)
}
```

### Form/Insert Mode

In forms and text input:

```lua
keys = {
  tab = "tab",         -- Next field
  shift_tab = "shift+tab",  -- Previous field
  space = " "          -- Toggle/activate buttons
}
```

### Global Bindings

Work everywhere:

```lua
keys = {
  force_quit = "ctrl+c",  -- Emergency exit
  enter = "enter",        -- Confirm
  back = "esc"            -- Cancel/Back
}
```

## Special Considerations

### Terminal Intercepted Keys

Some terminals intercept certain keys:

| Key | Issue | Solution |
|-----|-------|----------|
| `Ctrl+s` | Freezes terminal (XON/XOFF) | Press `Ctrl+q` to resume, or use `Ctrl+n` instead |
| `Ctrl+c` | Usually sends SIGINT | Envy intercepts this for force quit |
| `Ctrl+z` | Suspends process | Don't use in keybindings |
| `Ctrl+\` | Sends SIGQUIT | Avoid |

### SSH/Remote Sessions

When using Envy over SSH:

- Test your keybindings locally first
- Some terminal emulators handle keys differently
- `Ctrl+c` might behave unexpectedly

### TMUX/Screen

When using terminal multiplexers:

- You may need to configure the multiplexer to pass certain keys through
- `Ctrl+a` conflicts with GNU Screen prefix
- `Ctrl+b` conflicts with TMUX prefix

`Note: I use Ctrl+a as my tmux prefix, so most of the keybinds have been tested with that`

## Disabling Keys

Set a key to empty string to disable it:

```lua
keys = {
  -- Disable Vim navigation entirely
  vim_up = "",
  vim_down = "",
  vim_left = "",
  vim_right = "",

  -- Disable delete key (accident prevention)
  delete = ""
}
```

**Note:** Navigation keys (`up`, `down`, `left`, `right`) cannot be disabled as they're essential.

## Conflict Resolution

If two actions have the same key:

1. Last one in the config wins (undefined behavior)
2. Avoid duplicates in your config

Example of bad config:
```lua
-- DON'T DO THIS
keys = {
  quit = "q",
  search = "q"  -- Conflict! Which one wins?
}
```

## Testing Keybindings

After changing keybindings:

1. Save config file
2. Restart Envy
3. Test each binding in its context:
   - Grid navigation
   - Detail view actions
   - Form input
   - Search mode

## Keybinding Reference Table

| Config Key | Default | Description | Context |
|------------|---------|-------------|---------|
| `up` | `"up"` | Move up | Grid, Detail |
| `down` | `"down"` | Move down | Grid, Detail |
| `left` | `"left"` | Move left | Grid |
| `right` | `"right"` | Move right | Grid |
| `vim_up` | `"k"` | Vim up | Grid, Detail |
| `vim_down` | `"j"` | Vim down | Grid, Detail |
| `vim_left` | `"h"` | Vim left | Grid |
| `vim_right` | `"l"` | Vim right | Grid |
| `enter` | `"enter"` | Confirm | Everywhere |
| `back` | `"esc"` | Go back | Detail, Forms |
| `quit` | `"q"` | Quit | Grid |
| `search` | `"i"` | Search | Grid |
| `yank` | `"y"` | Copy | Detail |
| `create` | `"shift+n"` | Create project | Grid |
| `edit` | `"e"` | Edit key | Detail |
| `edit_project` | `"shift+e"` | Edit project | Detail |
| `delete` | `"d"` | Delete | Grid, Detail |
| `save` | `"shift+s"` | Save | Forms |
| `add` | `"shift+a"` | Add key | Forms |
| `history` | `"H"` | View history | Detail |
| `tab` | `"tab"` | Next field | Forms |
| `shift_tab` | `"shift+tab"` | Previous field | Forms |
| `space` | `" "` | Toggle | Forms |
| `force_quit` | `"ctrl+c"` | Force quit | Everywhere |

---

**Next:** Learn about [File Locations](./file-locations.md).
