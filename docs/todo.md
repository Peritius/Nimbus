# TODO (открытые задачи и идеи)

Единственный список задач. Всё, что не делается сейчас, — сюда.
Контекст решений — `architecture.md`, пользовательское поведение — `product.md`.

## В работе

(пусто)

## Очередь

План унификации UI, вариант A: одна продуктовая модель в Go, три тонких
рендерера (GTK, Qt, Win32). Согласован 2026-11-17. Каждый шаг собирается и
проверяется отдельно; ни один не требует править два бэкенда сразу.

### P1 — паритет

- [ ] **Шаг 2. Позиция панели в Go.**
      `corner`/`clamp`/`inside` живут трижды: GTK (`placePanel`, тест в
      `internal/ui/settings_linux_test.go`), Win32 (`pointerAnchor`,
      `rememberedWork`), C++ (`Panel::place`). Третья копия — причина бага Qt.
      Сделать `panelspec.Origin(size, want, pointer, work) (Point, bool)`,
      добавить `nimbus_qt_forecast_size(&w,&h)` (с `adjustSize`) и
      `nimbus_qt_forecast_show(id, x, y, tramp)`, который только `move`.
      Тесты перенести в общий пакет, копии удалить.

- [ ] **Шаг 3. `internal/panelctl` — жизненный цикл панели.** Главный выигрыш.
      Три копии (ожидание кэша, 5-секундный fallback, toggle, подписка на
      `weather.UpdateCh`, отписка): `internal/ui/forecast_linux.go`,
      `internal/qt/forecast_linux.go`, `internal/ui/forecast_windows.go`.
      Интерфейс рендерера: `Up() bool`, `Show(spec)`, `Update(rows)` (на месте,
      без пересоздания), `Close()`.
      Здесь же: **`nimbus_qt_forecast_update(rows)`** — Qt пересоздаёт панель при
      автообновлении (`refreshPanel`: close → build), из-за чего окно мигает и
      отпрыгивает к сохранённой позиции, а при первом запуске — к курсору;
      удалить мёртвое `gui.Forecast.Interval`; закрыть гонки в `internal/tray`
      (`cfg`, `lastData`, `lang` читаются без синхронизации).
      Критерий: панель не мигает и не двигается при обновлении; `go test -race`
      на сценарии «открыл–закрыл–открыл».

- [ ] **Шаг 4. `panelctl` на Qt и Windows + About.**
      Перевести оставшиеся два бэкенда. Добавить логотип в Qt About
      (`nimbus_qt_about` рисует только текст; GTK и Win32 — логотип).
      Сверить `qtshim/shim.h` ↔ `internal/qt/load_linux.go`: 28 имён
      привязываются строками, опечатка ловится только в рантайме — нужен тест
      или `nimbus_qt_api_version()` с проверкой при загрузке.

- [ ] **Шаг 5. `internal/settingsform` — состав формы вместо четырёх копий.**
      Правила «какие поля, какие опции, как читать, как сохранять» написаны
      четыре раза: `internal/ui/settings_linux.go`,
      `internal/qt/settings_linux.go`, `internal/ui/settings_windows.go`,
      `internal/webui/webui.go`. Модель `Field{Key, Group, Kind, Label, Value,
      Options, Active, Action}` + `Build(cfg, lang)` + `Adopt(cfg, values,
      action)`. В модели нет координат, размеров и типов виджетов. Образец —
      Qt-сторона (`qtForm*`), переносится на GTK и Win32.
      Здесь же:
      - интервалы: один локализованный список (сейчас «5 min / 1 hour» не
        переводятся), **дефолт 5 минут** (`config.Default()` ставит 10, которых
        в списке нет — первое «Сохранить» молча меняет значение), тест
        `TestIntervalIndexUnknownValue` переписать, дубли в `webui` убрать;
      - **webui → заглушка** «поставьте nimbus-gtk или nimbus-qt» (один HTML во
        временной папке + `xdg-open`, ~50 строк вместо 312 + три шаблона);
      - **WM_CLASS/`StartupWMClass`**: Qt называет себя `nimbus`, GTK — по
        basename бинаря; в `dist/*.desktop` нет `StartupWMClass`, поэтому окна
        группируются по-разному и без иконки.

### P2 — продукт и UX

- [ ] **Иконка трея: реальные цвет и тема.** `icons.MakeIcon` игнорирует `code`
      и `theme`, хотя README и `tray.go` обещают цвет по температуре. Цифры
      рисуются `basicfont.Face7x13` с nearest-апскейлом; на Linux уходит кадр
      64×64 и уменьшается до 22–24 px — заметное мыло. Рисовать тем же
      opentype-путём, что `fonts.Glyph`. ~100 строк, единственное заметное
      визуальное улучшение за пределами плана A.
- [ ] **README:** интервал в примере конфига (10 → 5 после шага 5); GNOME без
      расширения AppIndicator иконки не покажет; на Wayland позиция панели не
      запоминается (обе сборки); Qt-сборка: темы нет на Qt < 6.8, sticky только
      в GTK.

### P3 — чистка и тесты

- [ ] Тесты `icons`/`fonts`/`i18n`: валидный PNG нужного размера, непустой
      глиф, все строки непустые для en/uk.
- [ ] Мёртвый i18n: 25 методов из 64 не вызываются (`AutoDetect`, `EditConfig`,
      `ResetConfig`, `Refresh`, `UnitLabel`, `ThemeLabel`, `WeatherLine`,
      `PressureUnit*`, `Language*`…), −150 строк.
- [ ] `weather.DetectLocation` — ходит по http, не вызывается, −25 строк.
- [ ] Общий `http.Client` + `io.LimitReader` вместо нового клиента на каждый
      запрос в `internal/weather/weather.go`.

## Не делаем (осознанно)

- **Пиксельная идентичность.** Гарантируется одинаковый состав, порядок, подписи
  и метрики; отступы рисует тулкит.
- **Тема на Qt < 6.8 через `QPalette`.** Дала бы паритет, но убила бы
  нативность, ради которой Qt-сборка существует.
- **Sticky в Qt.** X11 — только через xcb (новая зависимость), Wayland — правило
  KWin. Документируем как ограничение.
- **Переписывание на другой язык.** Дублирование создают три тулкита, а не Go.
- **Один Qt на всё.** Потеря нативности в GNOME и Qt в Windows-сборке.

## Готово

- 2026-11-17 — нулевой шаг: `test(weather): cover the forecast cache`
  (`514a8ec`), `chore(weather): drop the unused Emoji helper` (`f11b13a`),
  удалён устаревший `stash@{0}`.
- 2026-11-17 — CI (`cb2d760`): job `test-linux-qt` в `ubuntu:22.04`
  (`make shim`, `go vet/test/build -tags qt`, проверка glibc ≤ 2.34) вместо шага
  `go build -tags qt ... || echo "skipping"`; `go test -race` отдельным шагом.
  Прогон `-race` на Fedora VM чистый; Qt-job в Actions ещё не запускалась.
- 2026-11-17 — документация: `architecture.md`, `product.md`, `todo.md`,
  `history.md`.
