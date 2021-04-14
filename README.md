# CLI xkcd
CLI tools to cache and show strips from xkcd.com. Exercise from the book "The GO Programming Language"

## Run

```bash
go get -u ./...

go install ./cmd/xkcd-cache
go install ./cmd/xkcd-cache-view

# If no 'out' argument is passed then the ouput is standard output
# 'w' (workers) has a default of 60
$(go env GOROOT)/bin/xkcd-cache -w 50 -out cache-$(date +%Y%m%d).json 
```

Commands can also be chained:

```bash
$(go env GOROOT)/bin/xkcd-cache | $(go env GOROOT)/bin/xkcd-cache-view
```