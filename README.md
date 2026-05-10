# owlet

A lightweight terminal command snippet selector.

## Snippets

`owlet` reads snippet files from:

    ~/.owlet/snippets/

Example:

    ~/.owlet/
    └─ snippets/
       ├─ conda.toml
       └─ git.toml

Each file contains command snippets:

    [[snippets]]
    command = "conda env list"
    desc = "View all conda environments."

    [[snippets]]
    command = "conda activate myenv"
    desc = "Activate a conda environment."

`command` is required. `desc` is optional.

## Keys

| Key | Action |
|---|---|
| `↑` / `↓` | Move |
| `Tab` | Show details |
| `Enter` | Copy and print command |
| `Esc` / `Ctrl+C` | Quit |