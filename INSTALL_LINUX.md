# SimpleAI Installation für Linux

## Schnellstart

SimpleAI läuft als portable AppImage-Datei auf allen Linux-Distributionen (Debian, Ubuntu, SUSE, Arch, Fedora, etc.).

### 1. Download

Lade die neueste Version herunter:

**[Download SimpleAI.AppImage](https://github.com/chrilep/SimpleAI/raw/main/automated-prereleases/SimpleAI.AppImage)**

Die Datei landet in deinem `Downloads`-Ordner.

### 2. Installation

```bash
# Verschiebe AppImage nach ~/bin
mkdir -p ~/bin
mv ~/Downloads/SimpleAI.AppImage ~/bin/

# Mache die Datei ausführbar
chmod +x ~/bin/SimpleAI.AppImage
```

### 3. Starten

```bash
~/bin/SimpleAI.AppImage
```

**Fertig!** SimpleAI läuft jetzt.

---

## Optional: System-Integration

Wenn du SimpleAI im Anwendungsmenü mit Icon haben möchtest:

```bash
# Wechsle ins Verzeichnis
cd ~/bin

# Extrahiere das AppImage
./SimpleAI.AppImage --appimage-extract

# Führe das Integrations-Script aus
cd squashfs-root
./integrate.sh

# Aufräumen
cd .. && rm -rf squashfs-root
```

**Das war's!** SimpleAI erscheint jetzt im Anwendungsmenü deiner Desktop-Umgebung (KDE, GNOME, XFCE, etc.) mit Icon.

---

## Anforderungen

### Für Window-Positionierung (empfohlen)

SimpleAI speichert Fenster-Positionen. Dafür wird `xdotool` benötigt:

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

Ohne `xdotool` öffnen sich Fenster bei jedem Start an der Standard-Position.

---

## Updates

Neue Versionen installieren:

```bash
# Alte Version löschen
rm ~/bin/SimpleAI.AppImage

# Neue Version herunterladen und Schritte 1-3 wiederholen
```

Falls du die System-Integration genutzt hast, musst du diese **nicht** erneut durchführen - die Verknüpfung funktioniert automatisch weiter.

---

## Deinstallation

### AppImage entfernen:

```bash
rm ~/bin/SimpleAI.AppImage
```

### System-Integration entfernen:

```bash
rm ~/.local/share/applications/SimpleAI.desktop
rm ~/.local/share/icons/hicolor/*/apps/simpleai.png
```

### Einstellungen löschen (optional):

```bash
rm -rf ~/.config/SimpleAI
rm -rf ~/.cache/SimpleAI
```

---

## Troubleshooting

### "Permission denied" beim Starten

Die Datei ist nicht ausführbar:

```bash
chmod +x ~/bin/SimpleAI.AppImage
```

### Kein Icon im Dateimanager

Das ist normal für portable AppImages. Nutze die System-Integration (siehe oben), um ein Icon zu erhalten.

### Fenster-Positionen werden nicht gespeichert

Installiere `xdotool` (siehe Anforderungen).

### SimpleAI erscheint nicht im Anwendungsmenü

Führe das Integrations-Script erneut aus oder aktualisiere die Desktop-Datenbank manuell:

```bash
update-desktop-database ~/.local/share/applications/
```

---

## Hinweise

- **Kein sudo nötig**: Führe alle Befehle als normaler Benutzer aus
- **Portable**: Die AppImage kann auch von USB-Sticks oder anderen Orten gestartet werden
- **Auto-Updates**: Nutze den GitHub-Link oben - die Datei wird bei jedem Commit automatisch aktualisiert

---

## Support

Bei Problemen oder Fragen:

- [GitHub Issues](https://github.com/chrilep/SimpleAI/issues)
- E-Mail: christian@lepthien.info
