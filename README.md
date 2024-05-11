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
