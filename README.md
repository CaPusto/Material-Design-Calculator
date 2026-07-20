# Material Design Calculator

A simple cross-platform desktop calculator with Material Design dark/light theme, built with [Fyne](https://fyne.io/) and [expr-lang/expr](https://github.com/expr-lang/expr).

<p align="center">
  <table>
    <tr>
      <td align="center"><b>On Rocky Linux - dark theme</b></td>
      <td align="center"><b>On Windows 11 - light theme</b></td>
    </tr>
    <tr>
      <td><img src="Screenshot 2026-07-20 140420.png" alt="qr" width="400"></td>
      <td><img src="Screenshot 2026-07-20 135658.png" alt="qr" width="400"></td>
    </tr>
  </table>
</p>




## Features

- Basic arithmetic: `+`, `-`, `×`, `÷`
- Unary operations: `1/x`, `x²`, `√x`, `%`, `±`
- Parentheses and expression chaining
- Keyboard input support
- History panel with expression recovery
- Clipboard copy/paste
- Memory: MC, MR, M+, M-
- π constant
- Division-by-zero protection (AST-level)
- i18n: English / Русский (auto-detected, override with `-locale=ru`)
- Dark & light themes (follows system)

## Build

### Dependencies

**Linux:**
```bash
# Debian/Ubuntu
sudo apt install gcc pkg-config libgl1-mesa-dev libegl1-mesa-dev \
  libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev \
  libxkbcommon-dev libwayland-dev

# Rocky/RHEL
sudo dnf install gcc pkgconfig mesa-libGL-devel mesa-libEGL-devel \
  libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel \
  libxkbcommon-devel wayland-devel
```

**Windows (cross-compile from Linux):**
```bash
sudo dnf install mingw64-gcc
make build-windows
```

**Windows (native-compile from Windows):**
```bash
install mingw64-gcc
go mod tidy
go build -ldflags="-s -w -H=windowsgui" -o calculator.exe .
```

**macOS:** Xcode Command Line Tools.

### Build & run

```bash
make build     # Linux
make build-windows  # Windows .exe
make run       # build + launch
```

Or manually:
```bash
go build -ldflags="-s -w" -o calculator .
./calculator
```

## Usage

- Type expressions directly or use the on-screen buttons
- Press `=` or `Enter` to evaluate
- `Backspace` to delete last character, `Escape` to clear
- Click history entries to restore an expression
- `%` works as percentage of the preceding value (e.g. `100+10%` → `110`)
- Press `=` repeatedly to repeat the last operation

## Project structure

```
├── main.go       # GUI and calculator logic
├── theme.go      # Material Design theme (dark/light)
├── i18n.go       # Localization (EN/RU)
├── validator.go  # Zero-division AST checker
├── Makefile      # Build automation
```

## License

MIT
