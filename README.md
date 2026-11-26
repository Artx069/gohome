# GoHome

Util to show current worktime and possible leave times with Matrix integration.

Fork of [sbreitf1/gohome](https://github.com/sbreitf1/gohome) with support for the new Matrix login page (`login.jsf`) and configurable weekly working hours.

![example view of current worktime](example.png)

## Getting started

To download and install the utility execute the following commands:

```
go install github.com/Artx069/gohome@latest
```

Now start it with `gohome` from command line.

You will probably be asked to enter a Matrix host. Paste the same host you visit in your browser (including protocol and path) here. Finally, you need to enter your Matrix credentials.

![Login example](login.png)

These values are stored in `~/.config/gohome` (XDG compatible) and are used in all following runs. You can enter an empty password here to only store host and username. You will be prompted for your password on every run.

Old configs in `~/.gohome` will be automatically migrated.

## User Config

You can edit your user settings in `~/.config/gohome/userconfig.json`. Following values are available:

| Key | Description |
| --- | ----------- |
| `WeeklyHours` | Contractual working hours per week, e.g. `36`. The daily target is derived as `WeeklyHours / WorkDays`. |
| `WorkDays` | Working days per week (default `5`). |
| `TargetTime` | A target time as provided by parameter `-t` in format `08:00`. Ignored when `WeeklyHours` is set. |

Set your weekly hours from the command line:

```
gohome set workhours=36            # 36h / 5 days -> daily target 07:12
gohome set workhours=30 workdays=4 # 30h / 4 days -> daily target 07:30
gohome set workhours=              # reset
```

Your overtime for the day and the resulting flexi-time balance are calculated against this daily target. The parameter `-t` still overrides it for a single run.

Use parameter `--save-config` to persist command line parameters in user config.

## Extensions

- Integrate with Gnome Desktop using the [Gnome Extension](https://gitlab.com/sebjung/gohome-gnome-extension) by [sebjung](https://gitlab.com/sebjung)

## Thanks

- Thanks to [danielb42](https://github.com/danielb42) for the [initial idea and cool project name](https://github.com/danielb42/gohome)
- Thanks to [fleaz](https://github.com/fleaz) for suggesting XDG config
- Thanks to [Nerzal](https://github.com/Nerzal) for Matrix v4.2.2 Support
