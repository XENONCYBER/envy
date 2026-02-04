# TUI Guide

Complete guide to navigating Envy's Terminal User Interface.

## Launching the TUI

Start Envy without any arguments:

```bash
envy
```

If this is your first time, you'll create a master password. After that, the TUI opens directly to your project grid.

## Interface Overview

The TUI operates through two primary views:

**Grid View** — Your project dashboard showing all stored projects as cards you can browse and select from.

**Detail View** — Opens when you select a project, displaying all secrets within that project and providing tools to manage them.

## Grid View

The grid view serves as your home base. It displays all projects as cards arranged in a grid, with each card showing the project name, a preview of stored keys, and an environment badge indicating whether it's development, staging, or production.

When you first open Envy, you'll land here. The view adapts to your terminal size, automatically adjusting how many cards fit across the screen.

### Navigating the Grid

Use these keys to move around:

| Key | Alternative | Action |
|-----|-------------|--------|
| `↑` | `k` | Move up |
| `↓` | `j` | Move down |
| `←` | `h` | Move left |
| `→` | `l` | Move right |
| `Enter` | — | Open selected project |
| `i` | `/` | Activate search |
| `shift+n` | — | Create new project |
| `shift+d` | — | Delete selected project (with confirmation) |
| `q` | — | Quit Envy |

### Search

Press `i` to activate search. A search bar appears at the top of the grid. Start typing and projects filter in real-time based on your query.

Search has three modes you can cycle through by pressing `Tab`:

**All** — Searches both project names and the keys within them. Type "database" and you'll see projects with database-related keys even if the project name doesn't contain that word.

**Projects** — Limits search to project names only. Useful when you know the project name but not necessarily what's stored inside.

**Keys** — Searches only within key names. Helps find specific types of credentials across multiple projects.

Search works case-insensitively and supports partial matching. Typing `api` finds `API_KEY`, `my-api-service`, or anything containing those characters. Results update as you type, so you often don't need to finish typing the full name.

Press `Enter` or `Esc` to exit search mode and return to normal navigation.

### Environment Badges

Each project displays a colored badge indicating its environment. These use consistent colors across the entire interface:

**[DEV]** appears in green — This is your development environment where you experiment and build.

**[STAGE]** appears in yellow — Staging or pre-production where you test before going live.

**[PROD]** appears in red — Production environment handling live traffic and real data.

The color coding helps prevent mistakes. That red badge is a visual reminder that you're looking at production credentials.

## Detail View

Press `Enter` on any project in the grid to open its detail view. This shows all secrets stored in that project and provides tools to view, copy, edit, and manage them.

The layout presents your project name at the top with the environment badge, followed by a scrollable list of all keys. Each key displays its name and a masked representation of its value.

### Detail Navigation

Move through your secrets with these keys:

| Key | Action |
|-----|--------|
| `↑/↓` or `k/j` | Navigate between keys |
| `Enter` or `Space` | Reveal/hide secret value |
| `y` | Copy value to clipboard |
| `e` | Edit key value (opens sidebar) |
| `shit+e` | Edit project (add/remove keys) |
| `shift+h` | View version history (opens sidebar) |
| `shift+d` | Delete selected key (with confirmation) |
| `esc` or `q` | Return to grid view |

### Revealing Secrets

By default, all secret values appear masked as dots to prevent shoulder surfing and accidental exposure. The key name is visible but the value is hidden.

Press `Enter` or `Space` on any key to toggle between masked and revealed states. When revealed, you'll see the actual value — a database URL, API key, or password. Press the same key again to mask it.

Secrets automatically hide when you navigate to a different key, exit the detail view, or quit Envy. This automatic hiding protects you from leaving sensitive data visible on screen.

### Copying to Clipboard

Press `y` to copy the currently selected secret value to your system clipboard. This is often safer than revealing and manually copying, especially in shared spaces.

After copying, a status message appears confirming the action with a note that the clipboard clears automatically. After 30 seconds, Envy wipes the clipboard content, preventing accidental pastes later. This happens automatically in the background using your system's native clipboard functionality.

### Edit Sidebar

Press `e` on any key to open the edit sidebar. The screen splits, showing the project detail on the left and an editing panel on the right.

The edit panel displays which key you're modifying and provides an input field for the new value. Type your changes, then press `Enter` to save. The previous value automatically moves to history before being replaced, preserving your ability to revert if needed.

Press `Esc` to cancel editing and return to the detail view without making changes.

### History Sidebar

Press `H` to view version history for any key. The sidebar opens showing all previous values alongside the current one.

The current value appears with a green badge and shows when it was created. Below, previous values display with yellow badges, each showing the timestamp when that version was active. You can see up to five previous values directly; older versions exist in the vault but aren't displayed to keep the interface manageable.

This history proves invaluable when you need to rollback after a bad configuration change or track when specific credentials were introduced.

### Project Edit Mode

Press `E` to enter project edit mode. This full-screen interface lets you modify project structure — rename the project, add new keys, or remove existing ones.

The form displays the current project name in an editable field at the top. Below, a scrollable list shows all existing keys. Navigate through fields using `Tab` to move forward or `Shift+Tab` to move backward. Arrow keys also work for navigation.

To add keys, fill in the "New Key Name" and "New Key Value" fields, then press `+ Add` or `ctrl+a`. The key immediately appears in the list above, and you can add another. Continue this process to bulk-add multiple keys before saving.

To delete a key, navigate to the keys list, highlight the key you want to remove, and press `d`. A confirmation dialog appears to prevent accidents.

When finished, press `Save` or `ctrl+s` to persist changes. Press `Esc` or `q` to cancel and discard all modifications.

## Create Project View

Press `ctrl+n` from the grid to create a new project. A form appears guiding you through project setup.

Start by entering a project name. This identifies your project in the grid, so choose something descriptive like "webapp", "api-service", or "payment-processor".

Next, select the environment. Three boxes appear labeled DEV, PROD, and STAGE. Use arrow keys to move between them and `Space` or `Enter` to select. The chosen environment gets highlighted with its characteristic color — green for dev, red for prod, yellow for stage.

Then add your first key. Enter the key name (like `DATABASE_URL` or `API_KEY`) and its value. You don't have to stop at one key — press `+ Add` or `ctrl+a` to add the key to a pending list and immediately start adding another. The interface shows how many keys you've added and lists them below.

Continue adding as many keys as needed. When ready, press `Save` or `ctrl+s` to create the project with all its keys. The form closes and you return to the grid, now showing your new project card.

## Modes

The TUI has two modes shown in the bottom bar:

### NORMAL Mode

- Navigate with arrow keys/Vim keys
- Press single keys for actions (`q`, `y`, `e`, etc.)
- Green indicator in status bar

### INSERT Mode

- Type text into input fields
- Press `Esc` or `Enter` to return to NORMAL
- Orange/yellow indicator in status bar

You automatically enter INSERT mode when:
- Searching (`i` in grid view)
- Editing text fields (in create/edit forms)
- Editing key values (in sidebar)

## Confirmation Dialogs

Destructive actions trigger confirmation dialogs to prevent accidents. When you press `d` to delete a project or key, the interface pauses and asks for explicit confirmation.

The dialog clearly states what you're about to delete, including the project name and environment for context. Press `y` to proceed with deletion or `n` to abort. Pressing `Esc` also cancels the operation.

Confirmation appears for:
- Deleting entire projects from the grid
- Removing individual keys from a project

This extra step has saved many users from accidentally wiping out production credentials or entire project configurations.

## Tips and Best Practices

### Keyboard Shortcuts

- Learn Vim navigation (`h/j/k/l`) for faster movement
- Use `Tab` to cycle through form fields quickly
- `ctrl+c` force quits from anywhere (emergency exit)

### Search Efficiency

- Use `Tab` to switch search modes based on what you're looking for
- Partial matching means you don't need to type full names
- Search is case-insensitive

### Security Habits

- Don't leave secrets revealed on screen
- Copy to clipboard instead of manually typing when possible
- Remember clipboard auto-clears after 30 seconds

### Navigation Patterns

Most TUI sessions follow similar patterns. You'll typically start in the grid, drill down into projects to copy or edit secrets, then return to the grid.

From the grid, press `Enter` on any project to view its secrets. Once in the detail view, press `y` to copy values, `e` to edit them, or `H` to check history. Press `E` if you need to restructure the project itself. When finished, press `Esc` to climb back up to the grid.

For creating new projects, press `shift+n` from anywhere in the grid. The create form opens, lets you build out the project, then drops you back at the grid with your new project visible.

Search from the grid with `i` to quickly find projects without scrolling. Tab through search modes to narrow results.

## Troubleshooting TUI Issues

### Display Problems

If the interface looks garbled:
- Ensure your terminal supports Unicode
- Try a different terminal emulator
- Check that your terminal has sufficient colors (256+ recommended)

### Key Bindings Not Working

Some terminals intercept certain keys:
- `ctrl+s` might freeze terminal (press `ctrl+q` to resume)
- Use `ctrl+n` instead of `ctrl+s` for creating projects in those terminals
- Customize keys in [config.lua](../configuration/lua-config.md) if needed

### TUI Won't Launch

```bash
# Check if vault exists
ls -la ~/.envy/keys.json

# If corrupted, restore from backup
cp ~/.envy/keys.json.backup ~/.envy/keys.json
```

---

**Next:** Learn about [CLI Commands](./cli-commands.md) for scripting and automation.
