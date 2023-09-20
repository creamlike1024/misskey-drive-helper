# misskey-drive-helper
batch delete misskey drive files and folders

## Components
- `getFolderId`: get folder ids from folder name, **only search the root directory**, recursive search will be implemented in the future
- `rmFiles`: delete all files in folders by folder ids without deleting folder
- `rmFolders`: batch delete folders by folder ids, folder must be empty

Note: `rmFiles` will delete all files in folder regardless of whether they were attached to notes.
Due to misskey api `drive/files/attached-notes`'s low performance, `rmFiles` will not check whether files were attached to notes.

To obtain the folder ID using getFolderId, it is necessary to move the folder to be operated on to the root directory

## Usage

Envirionment variables `MISSKEY_SITE` and `MISSKEY_TOKEN` are required.

Token must have `read:drive` and `write:drive` permissions.

```bash
$ export MISSKEY_SITE="https://m.isle.moe"
$ export MISSKEY_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
$ ./getFolderId folder1 folder2 ... folderN | ./rmFiles
$ ./getFolderId folder1 folder2 ... folderN | ./rmFolders
$ unset MISSKEY_SITE; unset MISSKEY_TOKEN
```

```powershell
> $Env:MISSKEY_SITE="https://m.isle.moe"
> $Env:MISSKEY_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
> .\getFolderId.exe folder1 folder2 ... folderN | .\rmFiles.exe
> .\getFolderId.exe folder1 folder2 ... folderN | .\rmFolders.exe
> Remove-Item -Path Env:\MISSKEY_SITE; Remove-Item -Path Env:\MISSKEY_TOKEN
```
