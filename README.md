# lightning + diary = liary

## Overview

liary is fastest cli tool for create a diary.
You can write immediately when you want to create diary or add content to diary.

## Features

- Diary file is markdown
- Quickly open diary with your favorite editor
- Add content to diary from terminal. there is no need to open a diary in the editor
- Provide vim plugin

## Usage

### Open diary

```sh
# Make today diary and open with editor
liary

# Make diary of one day ago and open with editor
liary -b 1d

# Make Specified file and open with editor
liary -f ~/diary/2018/01/01.md
```

### Write diary

```sh
# Add content in today diary
liary -a "add content"

# Add code block in today diary
liary -c python -a "print('hello world')"
```
