# pb-templ-htmx-tailwind

POC with [PocketBase](https://pocketbase.io/) + [Templ](https://templ.guide/),
[HTMX](https://htmx.org/) + [hyperscript](https://hyperscript.org/)
and [Tailwind](https://tailwindcss.com/) + [daisyUI](https://daisyui.com/)

## Roadmap

- [x] Integrate `templ` with PocketBase
- [x] Integrate `tailwind` & `daisyUI` with `templ`
- [x] Integrate `htmx` with `templ` and PocketBase
- [x] Auth
  - [x] Simple login with a Pocketbase user (email/username + password workflow)
  - [x] Registering a new user
  - [x] Display user, customize navbar when logged-in
  - [x] Logout
  - [ ] Handle OAuth2
- [x] Proper form handling w/ per-field error
- [ ] Error handling
  - [x] Display error when endpoint fails
  - [ ] Consistency
- [ ] Lists
  - [ ] CRUD operations
  - [x] Display user's lists
  - [ ] Archive
    - [x] Archive a list
    - [x] Unarchive
    - [ ] Toggle to only display not archived lists


## Development

### Requirements

- [Go](https://go.dev/)
- [Bun](https://bun.sh/)
  - Only used to generate CSS files with tailwind
- [Make](https://www.gnu.org/software/make/)
  - Optional but will *make* your life easier

Once these two requirements are met:

```sh
$ make deps
$ make
```

### VSCode debug

The `.vscode` folder is configured to add a Run and a Debug launch commands.

### Build workflow

- Generate templ Go files
- Generate minified CSS with tailwind & daisyUI based on templ files
- Embed assets in Go binary (CSS, favicon, htmx, hyperscript, etc)

### Live mode

```sh
$ go install github.com/bokwoon95/wgo@latest
$ make live
```

Live mode profits from `Make`'s multi-process abilities by running three
commands in parallel:
- `bun run tailwindcss -m -i ./assets/tailwind.css -o ./assets/dist/styles.min.css --watch`
- `wgo -file .go -xfile=_templ.go go run main.go serve :: wgo -file .css templ generate --notify-proxy`
- `templ generate --watch --proxy="http://127.0.0.1:8090" --open-browser=false`

This will start templ's hot reload server that allows to automatically refresh
the page whenever a template changes without restarting the server. Modifying
the templates also changes the CSS file and triggers a page reload by sending
a payload to templ's hot reload server. Whenever a go file that is not a templ
generated file is changed, the backend restarts too.

Live mode also disables browser caching by setting the `Cache-Control` header to
`"no-store"` for assets. This ensures the main CSS file is not cached when templ
hot reloads the page.

### Docker

```sh
$ docker build -t ptht:latest .
$ docker run --rm -p 8080:8080 -v $PWD/dockerdata/pb_data:/pb_data --name ptht ptht:latest
```

<details>
  <summary>Details & tricks</summary>

  #### Bun in Go docker image

  As shown in the build workflow section, this project requires both templ, bun
  and Go to be run in a sequential order. At first a four stage Dockerfile was
  created:

  - Using first `golang:alpine` to install templ and generate templ go files
  - Then using `oven/bun:alpine` to build the CSS based on templ go files
  - Then back to using `golang:alpine` for the rest of the build
    - Installing dependencies, embedding assets and building
  - Then using `gcr.io/distroless/static` to serve

  The main issue was that the whole build directory was copied over from one
  step to another and the cache was invalidated way too often.

  The Dockerfile was switched to using `golang:debian` (because Bun can't be
  installed on alpine distros without glibc) and the bun setup script was used,
  which greatly improved caching but removed the ability to automatically bump
  bun's version using Renovate.

  In the end the chosen solution was to use the `oven/bun:debian` image as a
  first step and copy over the bun binary to the build step. This way the bun
  version is pinned and can be upgraded using renovate.

</details>
