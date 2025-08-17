# Chimp, a Go Web Frontend Toolkit

Chimp tries to provide as much as possible for the Hypermedia Templ stack. You can use the opinionated variant of this stack and it's CLI Tool or reuse the individual modules as building blocks for your app.

### Goal: make this tool fast to build frontend parts of apps with.

## Scaffolding CLI

Install using `go install github.com/s00500/chimp/cli/chimp@latest`


## Modules

- Broker
- SessionStore
- IconLibraries
- Common used JS files that can be included using handlers and script tags
- Common components like Modal, Notification and Loader
- a few helper functions and middleware like the **URLPathMiddleware** to get the current urlpath in the templates

- CLI to add templates and scaffold you rporoject
- CLI to install tools (like templ, tailwind and air)


