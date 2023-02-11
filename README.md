# Gods Unchained User Grabber

Since GU currently allows users to view your name and that can be searched on [GUDecks](https://gudecks.com), this means users often have to obscure their name in ways such as using "???" or certain characters the GU font doesn't support. Because of this it creates an unfair advantage for users who want to have a meaningful name.

This tools scans the log files of GU and finds your opponents ID, as soon as the game knows who your opponent is. This information is already readily available in the log files. This just allows you to view the IDs of "???" and more, creating an equal playing field.

Once this tool detects the opponents ID, it will automatically open their GUDecks link in your browser.

## Building from source

For those who don't like running arbitrary binaries.

```sh
  git clone https://github.com/hysmio/gu-user-grabber
  cd gu-user-grabber
  go build .
  # the binary should now be available when running ./gu-user-grabber.exe or ./gu-user-grabber if not on windows
  # ps I haven't tested on linux or mac
```

## Configuration

You can create a `config.json` file, but running the application will create it for you.

**Note, if on Windows, your path will need double backslashes eg. `\\`, linux/osx can use `/`**

```json
{
  "file_path": "C:\\Users\\Username\\AppData\\LocalLow\\Immutable\\gods\\debug.json"
}
```
