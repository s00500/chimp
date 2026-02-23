# DataTable

Server-driven data table with SSE pagination, sorting, and inline row editing.

## Setup

```go
@components.DataTable("users",
    components.WithColumns([]components.Column{
        {Field: "name", Header: "Name", Sortable: true},
        {Field: "email", Header: "Email", Sortable: true},
        {Field: "role", Header: "Role"},
    }),
    components.WithDataEndpoint("/users"),
    components.WithPageSize(25),
)
```

## Sending Rows

The SSE endpoint receives query params: `page`, `pageSize`, `sortField`, `sortDir`.

```go
func handleUsers(w http.ResponseWriter, r *http.Request) {
    sse := datastar.NewSSE(w, r)
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
    sortField := r.URL.Query().Get("sortField")
    sortDir := r.URL.Query().Get("sortDir")

    users, total := fetchUsers(page, pageSize, sortField, sortDir)
    components.SendDataTableRows(sse, "users", total, UsersRows(users))
}
```

### Row Template

```go
templ UsersRows(users []User) {
    for _, user := range users {
        <tr id={ components.DataTableRowID("users", user.ID) } class="table-row">
            <td class="table-cell">{ user.Name }</td>
            <td class="table-cell">{ user.Email }</td>
            <td class="table-cell">
                @components.DataTableRowActions() {
                    <button class="btn btn-ghost btn-sm"
                        data-on:click={ fmt.Sprintf("@get('/users/%s/edit')", user.ID) }>
                        Edit
                    </button>
                }
            </td>
        </tr>
    }
}
```

## Inline Row Editing

Single rows can be swapped in-place using `SendDataTableRow`. This relies on idiomorph matching the `<tr>` by its element ID — no explicit CSS selector or merge mode needed.

### Row ID Convention

Use `DataTableRowID(tableID, rowID)` for consistent IDs across templates and handlers:

```go
// Generates "users-row-123"
components.DataTableRowID("users", "123")
```

Always set this as the `id` on your `<tr>` elements so idiomorph can find them.

### Edit Row Template (manual)

You can write edit rows manually for full control:

```go
templ UserEditRow(user User) {
    <tr id={ components.DataTableRowID("users", user.ID) } class="table-row"
        data-signals={ fmt.Sprintf("{editName: '%s', editEmail: '%s'}", user.Name, user.Email) }>
        <td class="table-cell">
            <input class="input" data-bind="editName"/>
        </td>
        <td class="table-cell">
            <input class="input" data-bind="editEmail"/>
        </td>
        <td class="table-cell">
            @components.DataTableRowActions() {
                <button class="btn btn-primary btn-sm"
                    data-on:click={ fmt.Sprintf("@put('/users/%s')", user.ID) }>
                    Save
                </button>
                <button class="btn btn-ghost btn-sm"
                    data-on:click={ fmt.Sprintf("@get('/users/%s/row')", user.ID) }>
                    Cancel
                </button>
            }
        </td>
    </tr>
}
```

### Edit Row Template (using DataTableEditRow)

For columns with `EditType` set, use the built-in `DataTableEditRow` component. It automatically renders the correct input type per column and wires up Save/Cancel buttons for the actions column.

Define columns with edit types:

```go
var userColumns = []components.Column{
    {Field: "name", Header: "Name", Sortable: true, EditType: components.EditText},
    {Field: "email", Header: "Email", Sortable: true, EditType: components.EditText},
    {Field: "role", Header: "Role", EditType: components.EditSelect, EditOptions: []string{"admin", "user", "viewer"}},
    {Field: "actions", Header: "", Type: components.ColumnActions, Width: "100px"},
}
```

Then in your handler, build a values map and render:

```go
func handleUserEdit(w http.ResponseWriter, r *http.Request) {
    sse := datastar.NewSSE(w, r)
    id := chi.URLParam(r, "id")
    user := fetchUser(id)
    values := map[string]string{
        "name":  user.Name,
        "email": user.Email,
        "role":  user.Role,
    }
    components.SendDataTableRow(sse, "users", id,
        components.DataTableEditRow("users", id, values, userColumns, "/users"))
}
```

The `DataTableEditRow` component:
- Renders `<input type="text">` for `EditText` columns
- Renders `<select>` for `EditSelect` columns (using `EditOptions`)
- Renders read-only text for columns without `EditType`
- Renders Save (POST to `endpoint/rowID`) and Cancel (GET to `endpoint/rowID/cancel`) buttons for `ColumnActions` columns
- Automatically initializes edit signals via `data-signals`

### ReadEditSignals

Use `ReadEditSignals` in your save handler to read edited values from the Datastar request:

```go
func handleUserSave(w http.ResponseWriter, r *http.Request) {
    sse := datastar.NewSSE(w, r)
    id := chi.URLParam(r, "id")

    edited, err := components.ReadEditSignals(r, userColumns)
    // edited["name"], edited["email"], edited["role"] contain the new values

    // ... save logic ...
    user := fetchUser(id)
    components.SendDataTableRow(sse, "users", id, UserRow(user))
}
```

### Handlers

```go
// GET /users/{id}/edit — swap to edit mode
func handleUserEdit(w http.ResponseWriter, r *http.Request) {
    sse := datastar.NewSSE(w, r)
    id := chi.URLParam(r, "id")
    user := fetchUser(id)
    components.SendDataTableRow(sse, "users", id, UserEditRow(user))
}

// POST /users/{id} — save and swap back to display mode
func handleUserSave(w http.ResponseWriter, r *http.Request) {
    sse := datastar.NewSSE(w, r)
    id := chi.URLParam(r, "id")
    edited, _ := components.ReadEditSignals(r, userColumns)
    // ... save edited values ...
    user := fetchUser(id)
    components.SendDataTableRow(sse, "users", id, UserRow(user))
}

// GET /users/{id}/cancel — cancel edit, restore display row
func handleUserCancel(w http.ResponseWriter, r *http.Request) {
    sse := datastar.NewSSE(w, r)
    id := chi.URLParam(r, "id")
    user := fetchUser(id)
    components.SendDataTableRow(sse, "users", id, UserRow(user))
}
```

## How It Works (Idiomorph)

Datastar's default merge mode is **outer** with **idiomorph**. When no CSS selector is specified, idiomorph matches incoming HTML against the DOM by element `id` attributes.

This means:

- `SendDataTableRows` wraps your rows in a `<tbody id="users-body">` — idiomorph finds the matching tbody and morphs its contents.
- `SendDataTableRow` sends a `<tr id="users-row-123">` — idiomorph finds the matching row and morphs it in place.

No explicit `WithSelector()` or `WithModeOuter()` needed. Just give elements good IDs.

### When You Still Need Explicit Options

Not all SSE helpers can rely on idiomorph. Modes other than outer (the default) still need to be specified:

| Helper | Why it needs options |
|--------|-------------------|
| `SendNotification` | Uses **append** mode to stack multiple toasts into `#notifications` |
| `SendAutocompleteResults` | Uses **inner** mode to replace listbox children |
| `SendDataTableRows` | Idiomorph — no options needed |
| `SendDataTableRow` | Idiomorph — no options needed |

## Available Helpers

| Function | Description |
|----------|-------------|
| `DataTableRowID(tableID, rowID)` | Returns conventional row element ID (`"users-row-123"`) |
| `SendDataTableRows(sse, id, total, rows)` | Patch all rows + update pagination signals |
| `SendDataTableRow(sse, tableID, rowID, row)` | Patch a single row in place |
| `ReadEditSignals(r, columns)` | Read edited values from Datastar request signals |
| `EditSignalName(field)` | Returns signal name for a field (e.g. `"name"` → `"editName"`) |

## Templates

| Component | Description |
|-----------|-------------|
| `DataTable(id, options...)` | Main data table shell with pagination and sorting |
| `DataTableRowActions()` | Wrapper for action buttons in a row cell |
| `DataTableEditRow(tableID, rowID, values, columns, endpoint)` | Edit-mode row with inputs per column |

## Options

| Option | Description |
|--------|-------------|
| `WithColumns(cols)` | Column definitions (Field, Header, Sortable, Width, EditType, EditOptions) |
| `WithDataEndpoint(url)` | SSE endpoint for fetching rows |
| `WithPageSize(n)` | Rows per page (default: 25) |
| `WithSignalPrefix(p)` | Custom signal namespace (default: table ID) |
| `WithSelectable()` | Enable row selection |
| `WithHidePagination()` | Hide the pagination controls |

## Column EditType

| EditType | Description |
|----------|-------------|
| `EditNone` (default) | Column is read-only in edit mode |
| `EditText` | Renders `<input type="text">` |
| `EditSelect` | Renders `<select>` using `EditOptions` |
