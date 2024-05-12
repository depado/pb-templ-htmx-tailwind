# pb-templ-htmx-tailwind

POC with [PocketBase](https://pocketbase.io/), [Templ](https://templ.guide/),
[HTMX](https://htmx.org/) and [Tailwind](https://tailwindcss.com/) +
[daisyUI](https://daisyui.com/)

## Roadmap

- [x] Integrate `templ` with PocketBase
- [x] Integrate `tailwind` & `daisyUI` with `templ`
- [x] Integrate `htmx` with `templ` and PocketBase

## Development

### Requirements

- [Go](https://go.dev/)
- [Bun](https://bun.sh/)
  - Only used to generate CSS files with tailwind
- [Make](https://www.gnu.org/software/make/)

Once these two requirements are met:

```sh
$ go install github.com/a-h/templ/cmd/templ@latest
$ go install github.com/bokwoon95/wgo@latest
$ bun install
$ make dev
```

### Build workflow

- Generate templ Go files
- Generate minified CSS with tailwind & daisyUI based on templ files
- Embed assets in Go binary (CSS, favicon, htmx, etc)

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
