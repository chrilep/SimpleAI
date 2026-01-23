# SimpleAI Installation for Linux

## Quick Start

SimpleAI runs as a portable AppImage file on all Linux distributions (Debian, Ubuntu, SUSE, Arch, Fedora, etc.).

### 1. Download

Download the latest version:

**[Download SimpleAI.AppImage](https://github.com/chrilep/SimpleAI/raw/main/automated-prereleases/SimpleAI.AppImage)**

The file will be saved to your `Downloads` folder.

### 2. Copy it to your ~/bin/ folder and Run it (same if first use or update)

```bash
# Move AppImage to ~/bin
mkdir -p ~/bin
mv ~/Downloads/SimpleAI.AppImage ~/bin/

# Make the file executable
chmod +x ~/bin/SimpleAI.AppImage

# Run from console to test if it works (see some logging if currently active)
~/bin/SimpleAI.AppImage
```

## 3. Install it for current user (optional, no root needed)

```bash
# Change to directory
cd ~/bin

# Extract the AppImage temporarily
./SimpleAI.AppImage --appimage-extract

# Run the integration script with path to AppImage
cd squashfs-root
./integrate.sh ~/bin/SimpleAI.AppImage

# Clean up extracted files
cd .. && rm -rf squashfs-root
```

**That's it!** SimpleAI now appears in your desktop environment's application menu (KDE, GNOME, XFCE, etc.) with an icon.

---

## Requirements

### For Window Positioning (Recommended)

SimpleAI saves window positions. This requires `xdotool`:

**Debian/Ubuntu:**

```bash
sudo apt-get install xdotool
```

**openSUSE/SUSE:**

```bash
sudo zypper install xdotool
```

**Fedora/RHEL:**

```bash
sudo dnf install xdotool
```

**Arch Linux:**

```bash
sudo pacman -S xdotool
```

Without `xdotool`, windows will open at the default position on every start.

---

## Updates

To install new versions:

```bash
# Delete old version
rm ~/bin/SimpleAI.AppImage

# Download new version and repeat steps 1-3
```

If you used system integration, you **don't** need to repeat it - the shortcut will continue to work automatically.

---

## Uninstallation

### Remove AppImage:

```bash
rm ~/bin/SimpleAI.AppImage
```

### Remove system integration:

```bash
rm ~/.local/share/applications/SimpleAI.desktop
rm ~/.local/share/icons/hicolor/*/apps/simpleai.png
```

### Delete settings (optional):

```bash
rm -rf ~/.config/SimpleAI
rm -rf ~/.cache/SimpleAI
```

---

## Troubleshooting

### "Permission denied" when starting

The file is not executable:

```bash
chmod +x ~/bin/SimpleAI.AppImage
```

### No icon in file manager

This is normal for portable AppImages. Use system integration (see above) to get an icon.

### Window positions are not saved

Install `xdotool` (see Requirements).

### SimpleAI doesn't appear in application menu

Run the integration script again or manually update the desktop database:

```bash
update-desktop-database ~/.local/share/applications/
```

### Error: "Program not found" or "squashfs-root not found"

This happens if you ran an older version of the integration script. Fix it by removing the old desktop file and running integration again:

```bash
# Remove old desktop file
rm ~/.local/share/applications/SimpleAI.desktop

# Extract AppImage and run integration with correct path
cd ~/bin
./SimpleAI.AppImage --appimage-extract
cd squashfs-root
./integrate.sh ~/bin/SimpleAI.AppImage
cd .. && rm -rf squashfs-root
```

---

## Notes

- **No sudo required**: Run all commands as a regular user
- **Portable**: The AppImage can also be launched from USB drives or other locations
- **Auto-updates**: Use the GitHub link above - the file is automatically updated with every commit

---

## Support

For issues or questions:

- [GitHub Issues](https://github.com/chrilep/SimpleAI/issues)
- E-Mail: christian@lepthien.info
