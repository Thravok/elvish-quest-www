# elvish-quest-www

Elvish Quest (not to be confused with Elvish Email) is a collection of open-source tools hosted for free.

This repository is a minimal static site served by a single Go binary with embedded HTML/CSS. No database, auth, or runtime dependencies beyond the Go standard library.

## Local development

Requires Go 1.23+.

```bash
make run
# or: go run .
```

Open [http://localhost:8080](http://localhost:8080). Health check: [http://localhost:8080/health](http://localhost:8080/health).

Build a binary:

```bash
make build
./elvish-quest-www
```

Set `PORT` to change the listen port (default `8080`).

## Docker

Build and run:

```bash
make docker-build
docker run --rm -p 8080:8080 elvish-quest-www
```

Or with Compose:

```bash
docker compose up --build
```

Configure your orchestrator (Coolify, Kubernetes, etc.) to probe `GET /health` on port `8080`.

## License

GNU GPL v3 — see [LICENSE](LICENSE).
