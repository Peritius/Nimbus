<div align="center">

# Nimbus

![Nimbus](nimbus1.png)

**Free and open-source cross-platform weather tray app. No ads, no trackers, no API keys.**

[![Windows & Linux](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-blue?logo=windows&logoColor=white)](#download)
[![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![License](https://img.shields.io/badge/license-GPL--3.0-blue)](LICENSE)

[![Release](https://img.shields.io/github/v/release/Peritius/Nimbus?include_prereleases&logo=github)](https://github.com/Peritius/Nimbus/releases)

</div>

[Features](#features) · [Download](#download) · [Build from source](#build-from-source) · [Configuration](#configuration) · [License](#license)

## About

Nimbus is a lightweight cross-platform weather tray app for Windows and Linux. It sits quietly in your system tray showing the current temperature and displays a 7-day forecast panel or settings window natively. Data comes from Open-Meteo, so there is no account to create and no key to paste.

## Features

- **Temperature tray icon**: current temperature number on a clean background with condition indicator
- **7-day Forecast**: native window (Windows, Linux GTK/Qt) or fallback web view
- **Settings GUI**: native configuration window with instant preview
- **Units**: temperature (°C / °F), pressure (hPa / mmHg / inHg), wind (m/s / km/h)
- **Theme**: Auto / Dark / Light system integration
- **Languages**: English, Ukrainian (switchable in-app)
- **Font scale**: adjustable 1–100% scale for tray text
- **No console window**: clean background execution on Windows

## Download

### Windows
Pre-built binaries are available on [Releases](https://github.com/Peritius/Nimbus/releases).

### Linux — Two builds
Pick the one that matches your desktop. **nimbus-gtk** draws its windows with GTK 3 and contains no Qt; **nimbus-qt** draws them with Qt 6 and contains no GTK. Both can be installed side by side.

#### Debian / Ubuntu
```bash
sudo apt install ./nimbus-gtk_1.0.0-1_amd64.deb    # GTK build
sudo apt install ./nimbus-qt_1.0.0-1_amd64.deb     # Qt build
```

#### RHEL / Rocky / Fedora
```bash
sudo dnf install nimbus-gtk-1.0.0-1.x86_64.rpm
sudo dnf install nimbus-qt-1.0.0-1.x86_64.rpm
```

#### openSUSE
```bash
sudo zypper install nimbus-gtk-1.0.0-1.x86_64.rpm
sudo zypper install nimbus-qt-1.0.0-1.x86_64.rpm
```

## Build from source

### Requirements
- Go (version specified in `go.mod`)
- Qt 6 development package (`qt6-base-dev` / `qt6-qtbase-devel`) **for `nimbus-qt` build only**

### Windows
```bash
go build -ldflags="-s -w -H windowsgui" -o nimbus.exe .
```

### Linux
```bash
make gtk        # -> build/nimbus-gtk
make qt         # -> build/nimbus-qt
```

### Package builds
```bash
make deb-gtk    # Debian packages
make deb-qt
make rpm-gtk    # RPM packages
make rpm-qt
```

## Configuration

Created automatically on first run:
- **Windows**: `%APPDATA%\Nimbus\config.json`
- **Linux**: `~/.config/Nimbus/config.json`

```json
{
  "latitude": 50.4501,
  "longitude": 30.5234,
  "city_name": "Kyiv",
  "update_interval": 10,
  "units": "celsius",
  "pressure_unit": "hpa",
  "wind_unit": "ms",
  "icon_theme": "auto",
  "language": "en",
  "font_scale": 100
}
```

| Field | Values | Description |
|---|---|---|
| `latitude`, `longitude` | float | Coordinates |
| `city_name` | string | Display name |
| `update_interval` | int (minutes) | Refresh interval |
| `units` | `celsius` / `fahrenheit` | Temperature unit |
| `pressure_unit` | `hpa` / `mmhg` / `inhg` | Pressure unit |
| `wind_unit` | `ms` / `kmh` | Wind unit |
| `icon_theme` | `auto` / `dark` / `light` | Window theme |
| `language` | `en` / `uk` | UI language |
| `font_scale` | int 1–100 | Tray font scale % |

## Weather API

Powered by [Open-Meteo](https://open-meteo.com/) — free of charge, no API key required. IP geolocation via [ip-api.com](http://ip-api.com/).

## License

Nimbus is released under the [GNU General Public License v3.0](LICENSE).

---

<details>
<summary>Українською</summary>

# Nimbus

![Nimbus](nimbus1.png)

**Безкоштовний кросплатформенный інформер погоди в системному треї. Без реклами, без трекерів, без ключів API.**

## Про застосунок

Nimbus — це легкий кросплатформенный інформер погоди для Windows і Linux. Він працює в системному треї, показуючи поточну температуру, та надає вікно прогнозу на 7 днів або налаштування. Дані надходять з Open-Meteo, тому не потрібно акаунтів чи ключів API.

## Можливості

- Температура в системному треї у вигляді піктограми
- Прогноз на 7 днів у рідному вікні (Windows, Linux GTK/Qt) або веб-вкладці
- Вікно налаштувань із миттєвим прев'ю
- Одиниці вимірювання: температура (°C / °F), тиск (гПа / мм рт. ст. / inHg), вітер (м/с / км/ч)
- Тема оформлення: Авто / Темна / Світла
- Мови: англійська, українська
- Масштаб шрифту в треї (1–100%)
- Без консольного вікна на Windows

## Завантажити

### Windows
Готові бінарники доступні на сторінці [Releases](https://github.com/Peritius/Nimbus/releases).

### Linux — Дві збірки
- **nimbus-gtk**: використовує GTK 3 (без Qt).
- **nimbus-qt**: використовує Qt 6 (без GTK).
Їх можна встановлювати паралельно.

#### Debian / Ubuntu
```bash
sudo apt install ./nimbus-gtk_1.0.0-1_amd64.deb
sudo apt install ./nimbus-qt_1.0.0-1_amd64.deb
```

#### RHEL / Rocky / Fedora
```bash
sudo dnf install nimbus-gtk-1.0.0-1.x86_64.rpm
sudo dnf install nimbus-qt-1.0.0-1.x86_64.rpm
```

## Збірка з вихідного коду

### Вимоги
- Go (версія з `go.mod`)
- Пакет розробки Qt 6 (`qt6-base-dev` / `qt6-qtbase-devel`) лише для збірки `nimbus-qt`.

### Збірка
```bash
go build -ldflags="-s -w -H windowsgui" -o nimbus.exe .  # Windows
make gtk                                                # Linux GTK
make qt                                                 # Linux Qt
```

## Конфігурація

Створюється автоматично при першому запуску:
- **Windows**: `%APPDATA%\Nimbus\config.json`
- **Linux**: `~/.config/Nimbus/config.json`

## API погоди

Використовує [Open-Meteo](https://open-meteo.com/) — безкоштовно, без ключів API.

## Ліцензія

GNU General Public License v3.0.

</details>
