# Calibre Collection Sync

This is a KUAL extension to build/update collections on Kindle devices from Calibre

The extension can generate collections from `Calibre Kindle Collection Plugin` or from Calibre Metadata (`metadata.calibre`)

## Installation
Unzip release packet into root folder of your device

## How to use
Start KUAL and take a look at ccsync menu.

The update menu will not delete empty collections, it will only update existing ones or create a new one.

The rebuild menu first will delete old collections and then will create a new one.

## Calibre Metadata
For this method, there is no need to install any additional plugins in Calibre.
When `Send to Device` is used, Calibre writes/updates the `metadata.calibre` file into your device, which is then used for collection creation.

Some customization can be made by editing `extensions/ccsync/meta.ini`

If `tags` in `[general]` section is set to `1`, a collection per tag will be generated, the default value is `2` which means a tags map will be used.

One could customize collection creation from book tags (tags can be Created/Updated from Calibre). A `tag_name -> Collection Name` map can be set in `[tags]` section of the meta.ini file.
One could Combine/Merge multiple tags into a single collection using this feature.

```ini
[general]
; Series
;   0 - Don't create collection
;   1 - Create collection
series = 1

; Author
;   0 - Don't create collection
;   1 - Always create collection
;   2 - Only create collection if not in series
author = 2

; Tags
;   0 - Don't create collection
;   1 - Create collection by all tags
;   2 - Create collection from tags map
tags = 2

[tags]
;science fiction = Sci-Fi
;Научна Фантастика = Sci-Fi
```

## Build from source
In order to build from source a `C cross compiler` must to be used(sqlite3 dependency).
Please look build_kindle.sh