// Package config handles application configuration via Lua
package config

import (
	"os"
	"path/filepath"
	"runtime"

	lua "github.com/yuin/gopher-lua"
)

type AppConfig struct {
	Backend BackendConfig

	Keys KeyMap

	Theme Theme
}

type BackendConfig struct {
	KeysPath string

	LockPath string
}

func DefaultBackendConfig() BackendConfig {
	return BackendConfig{
		KeysPath: GetDefaultKeysPath(),
		LockPath: GetDefaultLockPath(),
	}
}

func DefaultAppConfig() AppConfig {
	return AppConfig{
		Backend: DefaultBackendConfig(),
		Keys:    DefaultKeyMap(),
		Theme:   DefaultTheme(),
	}
}

func GetConfigDir() string {
	return GetDefaultConfigDir()
}

func GetLuaConfigPath() string {
	return GetDefaultConfigPath()
}

func LoadAppConfig() AppConfig {
	config := DefaultAppConfig()

	luaPath := GetLuaConfigPath()
	if _, err := os.Stat(luaPath); os.IsNotExist(err) {
		return config
	}

	L := lua.NewState()
	defer L.Close()

	registerConfigFunctions(L)

	if err := L.DoFile(luaPath); err != nil {
		return config
	}

	config.Backend = extractBackendConfig(L, config.Backend)
	config.Keys = extractKeyMap(L, config.Keys)
	config.Theme = extractTheme(L, config.Theme)

	return config
}

func registerConfigFunctions(L *lua.LState) {
	envyMod := L.NewTable()

	home, _ := os.UserHomeDir()
	L.SetField(envyMod, "home", lua.LString(home))

	L.SetField(envyMod, "os", lua.LString(runtime.GOOS))

	L.SetField(envyMod, "default_data_dir", lua.LString(GetDefaultDataDir()))
	L.SetField(envyMod, "default_config_dir", lua.LString(GetDefaultConfigDir()))
	L.SetField(envyMod, "default_keys_path", lua.LString(GetDefaultKeysPath()))

	L.SetField(envyMod, "expand_path", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		home, _ := os.UserHomeDir()
		if len(path) > 0 && path[0] == '~' {
			path = filepath.Join(home, path[1:])
		}
		L.Push(lua.LString(path))
		return 1
	}))

	L.SetGlobal("envy", envyMod)
}

func extractBackendConfig(L *lua.LState, defaults BackendConfig) BackendConfig {
	config := defaults

	backend := L.GetGlobal("backend")
	if backend.Type() != lua.LTTable {
		return config
	}

	tbl := backend.(*lua.LTable)

	if val := tbl.RawGetString("keys_path"); val.Type() == lua.LTString {
		path := string(val.(lua.LString))
		config.KeysPath = expandPath(path)
	}

	if val := tbl.RawGetString("lock_path"); val.Type() == lua.LTString {
		path := string(val.(lua.LString))
		config.LockPath = expandPath(path)
	}

	return config
}

func extractKeyMap(L *lua.LState, defaults KeyMap) KeyMap {
	config := defaults

	keys := L.GetGlobal("keys")
	if keys.Type() != lua.LTTable {
		return config
	}

	tbl := keys.(*lua.LTable)

	// Navigation
	if val := tbl.RawGetString("up"); val.Type() == lua.LTString {
		config.Up = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("down"); val.Type() == lua.LTString {
		config.Down = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("left"); val.Type() == lua.LTString {
		config.Left = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("right"); val.Type() == lua.LTString {
		config.Right = string(val.(lua.LString))
	}

	// Vim navigation
	if val := tbl.RawGetString("vim_up"); val.Type() == lua.LTString {
		config.VimUp = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("vim_down"); val.Type() == lua.LTString {
		config.VimDown = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("vim_left"); val.Type() == lua.LTString {
		config.VimLeft = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("vim_right"); val.Type() == lua.LTString {
		config.VimRight = string(val.(lua.LString))
	}

	// Actions
	if val := tbl.RawGetString("enter"); val.Type() == lua.LTString {
		config.Enter = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("back"); val.Type() == lua.LTString {
		config.Back = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("quit"); val.Type() == lua.LTString {
		config.Quit = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("search"); val.Type() == lua.LTString {
		config.Search = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("yank"); val.Type() == lua.LTString {
		config.Yank = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("create"); val.Type() == lua.LTString {
		config.Create = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("edit"); val.Type() == lua.LTString {
		config.Edit = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("edit_project"); val.Type() == lua.LTString {
		config.EditProject = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("delete"); val.Type() == lua.LTString {
		config.Delete = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("save"); val.Type() == lua.LTString {
		config.Save = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("add"); val.Type() == lua.LTString {
		config.Add = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("history"); val.Type() == lua.LTString {
		config.History = string(val.(lua.LString))
	}

	// Form navigation
	if val := tbl.RawGetString("tab"); val.Type() == lua.LTString {
		config.Tab = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("shift_tab"); val.Type() == lua.LTString {
		config.ShiftTab = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("space"); val.Type() == lua.LTString {
		config.Space = string(val.(lua.LString))
	}

	// Special
	if val := tbl.RawGetString("force_quit"); val.Type() == lua.LTString {
		config.ForceQuit = string(val.(lua.LString))
	}

	return config
}

func extractTheme(L *lua.LState, defaults Theme) Theme {
	config := defaults

	theme := L.GetGlobal("theme")
	if theme.Type() != lua.LTTable {
		return config
	}

	tbl := theme.(*lua.LTable)

	// Colors
	if val := tbl.RawGetString("base"); val.Type() == lua.LTString {
		config.Base = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("text"); val.Type() == lua.LTString {
		config.Text = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("accent"); val.Type() == lua.LTString {
		config.Accent = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("surface0"); val.Type() == lua.LTString {
		config.Surface0 = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("surface1"); val.Type() == lua.LTString {
		config.Surface1 = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("overlay0"); val.Type() == lua.LTString {
		config.Overlay0 = string(val.(lua.LString))
	}

	// Semantic colors
	if val := tbl.RawGetString("success"); val.Type() == lua.LTString {
		config.Success = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("warning"); val.Type() == lua.LTString {
		config.Warning = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("error"); val.Type() == lua.LTString {
		config.Error = string(val.(lua.LString))
	}

	// Environment badges
	if val := tbl.RawGetString("prod_bg"); val.Type() == lua.LTString {
		config.ProdBg = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("dev_bg"); val.Type() == lua.LTString {
		config.DevBg = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("stage_bg"); val.Type() == lua.LTString {
		config.StageBg = string(val.(lua.LString))
	}

	// History section colors
	if val := tbl.RawGetString("current_bg"); val.Type() == lua.LTString {
		config.CurrentBg = string(val.(lua.LString))
	}
	if val := tbl.RawGetString("previous_bg"); val.Type() == lua.LTString {
		config.PreviousBg = string(val.(lua.LString))
	}

	// Grid layout
	if val := tbl.RawGetString("grid_cols"); val.Type() == lua.LTNumber {
		config.GridCols = int(val.(lua.LNumber))
	}
	if val := tbl.RawGetString("grid_visible_rows"); val.Type() == lua.LTNumber {
		config.GridVisibleRows = int(val.(lua.LNumber))
	}

	// Card dimensions
	if val := tbl.RawGetString("card_width"); val.Type() == lua.LTNumber {
		config.CardWidth = int(val.(lua.LNumber))
	}
	if val := tbl.RawGetString("card_height"); val.Type() == lua.LTNumber {
		config.CardHeight = int(val.(lua.LNumber))
	}

	return config
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[1:])
	}
	return path
}

func EnsureDataDir(config BackendConfig) error {
	dir := filepath.Dir(config.KeysPath)
	return os.MkdirAll(dir, 0o700)
}
