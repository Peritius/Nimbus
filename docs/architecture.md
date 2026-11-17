# Архитектура Nimbus

Десктопный трей-информер погоды: один Go-процесс, три нативных рендерера,
две Linux-сборки и одна Windows. Ядро (конфиг, сеть, i18n, трей) общее для всех.

## Стек

| Компонент | Выбор |
|---|---|
| Язык | Go 1.26.5, статические бинарники |
| Сборка | `CGO_ENABLED=0` везде, кроме Qt-shim |
| Трей | `fyne.io/systray` v1.12.2 (на Linux — `org.kde.StatusNotifierItem` через D-Bus) |
| GTK | `internal/gtk` — purego v0.10.2, `dlopen` libgtk-3 в рантайме, без cgo и заголовков |
| Qt | `qtshim/` — C++/cgo, вложенный модуль, собирается в `libnimbusqt.so`, вшивается в бинарь |
| Windows MPI | `github.com/lxn/win`, GDI-рисование |
| Шрифты/иконки | `golang.org/x/image` (Weather Icons TTF, opentype), `internal/appicon` (PNG) |
| Конфиг | JSON в `os.UserConfigDir()/Nimbus/`, атомарная запись |
| Локализация | `internal/i18n`, EN и UK |

## Продукты

| Бинарь | Окна | Зависимости |
|---|---|---|
| `nimbus.exe` | Win32/GDI | нет |
| `nimbus-gtk` | GTK3 | libgtk-3 в системе (иначе — заглушка/ошибка) |
| `nimbus-qt` | Qt6 | libnimbusqt.so внутри бинаря, Qt6 в системе |

Взаимоисключение тулкитов держится тегами сборки: `internal/ui/*gtk*` и
`internal/tray/toolkit_gtk.go` под `!qt`, Qt-часть под `qt`. GTK-бинарь не может
загрузить Qt, Qt-бинарь — не тянет `internal/gtk` вовсе.

## Контракт бэкендов (`internal/gui`)

`Backend` описан на уровне задач, а не виджетов:

```go
type Backend interface {
    Name() string
    Settings(cfg *config.Config, onFontScale func(int)) *config.Config
    Forecast(req Forecast)   // возвращается сразу, данные тянет сам
    About(theme string)
    Error(title, message string)
}
```

Правила, которые нельзя нарушать:

- **Никаких координат, размеров и типов контролов через границу.** Иначе
  абстракция превращается в «ещё один тулкит», который надо реализовывать трижды.
- **Любой метод можно звать из любой горутины.** Бэкенд сам маршалит в свою
  очередь (`gtk.Invoke`, `qt.invoke`).
- **`Forecast` не блокирует**: вызывающий — единственный диспетчер меню трея.
- `Forecast.OnMove` зовётся только при реальном перетаскивании (см. комментарий в
  `gui.go`: нужен и handoff в move-loop, и отличие конечной позиции от начальной —
  иначе клик или отменённый Escape затирают позицию).
- `Looper` реализует только тот бэкенд, который владеет своим event loop: Qt.
  У GTK цикл принадлежит процессу (`internal/tray`), у Win32 — свой message loop
  в `tray`.

Выбор бэкенда — `Factory{Name, Rank, Probe, Open}` + `Register` из `init()` в
файле под build tag. `NIMBUS_GUI_BACKEND` форсирует имя и **отключает probe**:
если форсированный бэкенд не открылся, возвращается `nullBackend`, а не другой
рендерер — иначе баг-репорт «показывает не то» невозможно воспроизвести.

## Потоки

- GTK: `gtk_init` и `gtk_main` на одном OS-thread (`runtime.LockOSThread` до
  `Init`). Любой вызов GTK — через `gtk.Invoke`.
- Qt: то же самое для `QApplication`; вход — `nimbus_qt_invoke`.
- Трей живёт на своей D-Bus горутине systray.
- Панель и форма не держат мьютексов трея: `cfg` — указатель, который меняется
  целиком (копия + swap), а не мутируется на месте.

Известная слабость: `internal/tray` читает `cfg`, `lastData`, `lang` без
синхронизации (описано в шапке `tray.go`). Закрывается вместе с `panelctl`
(см. `todo.md`).

## Qt-shim

`qtshim/shim.cpp` — не биндинг, а сам UI: виджеты, layout, палитра, event loop.
`shim.h` объявляет C-API, `internal/qt/load_linux.go` привязывает 28 имён через
purego из вшитой `.so` (пишется во временный файл и грузится).

- Никакого `Q_OBJECT`/moc: сигналы подключаются лямбдами, иначе `go build` не
  смог бы собрать shim одним шагом.
- События идут обратно колбэком: `CLOSED`, `MOVED`, `ACTION`, `SEARCH`, `PICK`,
  `SLIDE`.
- Тема применяется `QStyleHints::colorScheme` — только Qt ≥ 6.8;
  `nimbus_qt_can_theme` сообщает, умеет ли сборка.
- Позиция панели считается в C++ (`Panel::place`) — это одна из копий логики,
  которую убирает шаг 2 плана.

## Трей

`systray` даёт меню: **7-day Forecast**, **Settings…**, **About**, **Quit**.
Обновление — не пункт меню, а внутренний тикер: `startRefreshLoop` +
`weather.UpdateCh`.

- Иконка — температура числом, перерисовывается при смене данных и при движении
  слайдера масштаба шрифта (с задержкой 150 мс).
- На Linux в трей отдаётся PNG (не ICO: StatusNotifierItem декодирует байты сам,
  а pure-Go реализации знают только PNG), на Windows — ICO из нескольких кадров.
- Тултип: температура, ощущается, влажность, ветер, давление, состояние.

## Данные

Open-Meteo, два запроса (current и daily 7 дней) + геокодинг для поиска города.
`weather.Store`/`Cached` — кэш по координатам, `UpdateCh` — широковещательный
сигнал «данные обновились» для панели и трея. `FetchedAt` пишется при сохранении.

`config.Config` — плоский JSON: координаты и имя города, `update_interval`
(минуты), единицы (`celsius`/`fahrenheit`), давление (`hpa`/`mmhg`/`inhg`),
ветер (`ms`/`kmh`), `icon_theme` (`auto`/`dark`/`light`), язык, масштаб шрифта,
позиция панели (`forecast_x`/`forecast_y`, nil — «никогда не двигали»).

Маппинг «значение в конфиге ↔ пункт в форме» — `internal/config/options.go`
(`Index`, `Pick`, `ParseCoord`), общий для всех трёх форм. Раньше он был скопирован
в каждую и разъехался (обрезка пробелов в координатах на Windows).

## Сборка

`Makefile`:

| Цель | Результат |
|---|---|
| `build` | `build/nimbus.exe` (Windows) |
| `gtk` / `build-linux` | `build/nimbus-gtk` |
| `shim` | `internal/qt/libnimbusqt.so` (нужен `qmake6` + Qt 6 dev) |
| `qt` | `build/nimbus-qt` (сначала `shim`) |
| `deb-gtk`, `deb-qt`, `rpm-gtk`, `rpm-qt` | пакеты в `dist/` |

`.so` не в git. `internal/qt` под тегом `qt`, поэтому обычный `go build` не
требует ни C++ компилятора, ни Qt.

Qt-сборка делается в `ubuntu:22.04`, чтобы floor был glibc 2.34 / GLIBCXX 3.4 —
это покрывает Ubuntu 22.04, Debian 12, RHEL 9 и Fedora. Release-job измеряет
floor через `objdump` и падает, если он вырос.

## Ограничения платформ

| Ограничение | Причина |
|---|---|
| GNOME без расширения AppIndicator иконки не покажет | systray говорит на `org.kde.StatusNotifierItem`; Ubuntu ставит расширение, Fedora GNOME — нет |
| Wayland: позиция панели не запоминается | клиент не знает координат своего toplevel, оба чтения дают (0,0) → «не двигали» |
| Sticky-панель только в GTK | `gtk_window_stick`; в Qt это `_NET_WM_DESKTOP` через xcb или правило KWin |
| Тема окон Qt недоступна на Qt < 6.8 | `QStyleHints::colorScheme` появился в 6.8; ручная `QPalette` убила бы нативность |
