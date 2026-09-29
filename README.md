# EXERCISE PRESCRIPTION APP

Stacks (Possible to changed in future):
* Backend: Go
* Template: stdlib
* Reactivity: Nomini
* Database: SQLite

## Install latest Go

Follow the instructions from the [site](https://go.dev/doc/install) on how to install.

After finish installing Go, install air for live reload:
```
go install github.com/air-verse/air@latest
```
Then copy _.air.toml.example_ file into the _.air.toml_ file.

More details on air in their [github page](https://github.com/air-verse/air).

To run the app in development server, simply run `make start` in the console.

It will automatically start the SQLite DB and the app.

## Templating tips with html/template

For templating, I use standard Go html/template library. I choose .gohtml as the file type to make it easier to deal with Go templating syntax. For VSCode, I install 'Go HTML Language Server' extension by Onurcan Sever. This extension enabled developer to autocomplete Go templating syntax such as 'define', 'block', etc as well as dot completions of Go structs inside the template.

I also run 'go install github.com/gohugoio/gotmplfmt@latest' and install the 'gotmplfmt' extension by GoHugo for auto indentation and formatting, and configure user settings for gohtml to use the 'gotmplfmt' as default formatter.

And because .gohtml file did not support autocompletion for html tags, configurations in user settings is made to include gohtml in emmet.