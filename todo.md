# Project TODO

Two-person workflow. **Do the "Do Together First" section before splitting.**

---

## Do Together First (~20–30 min)

These must be agreed on before either person starts their track. Commit the results immediately so both people are unblocked.

- [x] Define all shared structs in `types.go` — `User` (with role field), `Classroom`, `CoinTransaction`, `ShopItem`, `AvatarConfig`
- [x] Decide on roles as constants: `student`, `teacher`, `admin`
- [x] Agree on JSON filenames: `users.json`, `classrooms.json`, `shop.json`, `attendance.json`, `transactions.json`
- [x] Sketch and write down every route, the HTTP method, and which role(s) can access it
- [x] Commit `types.go` and the agreed route list so both people can reference them

---

## Track A — Foundation (Person 1)

### Project Scaffold
- [x] `go mod init`
- [x] Set up folder structure: `main.go`, `store/`, `handlers/`, `templates/`, `static/`
- [x] `main.go` with a basic HTTP server and static file serving

### Persistence
- [x] Load application state from `data/data.json`
- [x] Save application state to `data/data.json`
- [x] Initialize nil maps after loading
- [x] Ensure `saveData()` is called after state changes
- [x] Create `data/` directory automatically if missing

### User Storage
- [x] Seed test users (student, teacher, admin)
- [x] Password hashing (`bcrypt`)
- [x] Password verification during login

### Auth Routes
- [x] `POST /login` — validate credentials, create session
- [x] `POST /logout` — destroy session
- [x] `POST /register` — if needed for admin setup

### Session Management
- [x] Cookie-based session helpers
- [x] In-memory session store
- [x] Secure session IDs
- [x] Session expiration

### Auth Middleware
- [x] `requireLogin` middleware — redirects to `/login` if no valid session
- [x] `requireRole(role string)` middleware — returns 403 if role doesn't match

**Track A Deliverable:** Running server where you can log in, get a session cookie, and log out. Middleware functions exported and ready.

---

### Feature Work (start after deliverable is done)

#### Coins
- [ ] `POST /coins/award` — teacher awards coins to a student
- [ ] Validate teacher is assigned to the student's classroom
- [ ] Record awards as `CoinTransaction` entries
- [ ] `GET /coins/balance/:userID` — return a student's current coin balance
- [ ] Compute balance from transaction history
- [ ] Double Day multiplier support

#### Shop
- [x] Add shop structs
- [ ] Seed shop items
- [ ] `GET /shop` — return available shop items
- [ ] `POST /shop/buy` — deduct coins and record purchase transaction
- [ ] Validate student has enough coins before purchase
- [ ] Persist purchases through transaction history

#### Attendance
- [x] Add attendance structs
- [ ] `POST /attendance` — teacher submits attendance for a classroom and date
- [ ] Validate teacher owns the classroom
- [ ] Prevent duplicate attendance submissions for the same day
- [ ] `GET /attendance/:classroomID` — return attendance records for a classroom
- [ ] `GET /attendance/export/:classroomID` — CSV export for reports

#### Admin / Classrooms
- [x] Add classroom structs
- [ ] `POST /classrooms` — create a classroom
- [ ] `POST /classrooms/:id/assign-teacher` — assign a teacher to a classroom
- [ ] `POST /classrooms/:id/assign-student` — assign a student to a classroom
- [ ] `GET /classrooms` — list all classrooms
- [ ] Validate users exist before assignment

#### Schedule / Double Days
- [x] Add schedule structs
- [ ] `POST /schedule/doubleday` — mark a date as a double day
- [ ] `GET /schedule` — return upcoming double days
- [ ] Prevent duplicate double-day entries
- [ ] Integrate double-day multiplier with coin awards

---

## Track B — UI Shell (Person 2)

### Setup
- [x] Static file serving confirmed working
- [x] `html/template` base layout with a `{{template "content" .}}` slot
- [x] CSS scaffold (can use a small utility CSS or write from scratch)
- [ ] Add a `fakeUser` variable at the top of the handlers file so you can swap it for a real session lookup in one line later

### Stub Pages — All Roles
- [x] `/login` — login form page (full page)
- [x] Logout — popup/modal component

### Stub Pages — Student
- [x] `/student/coins` — receive coins popup
- [x] `/student/shop` — shop full page
- [x] `/student/avatar` — avatar customization full page

### Stub Pages — Teacher
- [ ] `/teacher/attendance` — attendance full page
- [ ] `/teacher/reports` — reports full page
- [ ] `/teacher/schedule` — schedule / double days full page

### Stub Pages — Admin
- [ ] `/admin/reports` — all reports full page
- [ ] `/admin/classrooms` — classrooms full page

### Navigation
- [ ] Header component with role-aware nav links (use `fakeUser.Role` to show/hide links)
- [ ] Make sure nav links match the agreed route list

### Shared Components
- [ ] Form styles and button styles
- [ ] Table component (used in attendance, reports, classrooms)
- [ ] Popup/modal component (used for logout, receive coins)
- [ ] Flash/error message display for form feedback

**Track B Deliverable:** Every page reachable, looks roughly right, using the hardcoded fake user.

---

### Feature Work (start after deliverable is done)

#### Student Pages
- [ ] Shop page: fetch items from `/shop`, render item cards with buy buttons
- [ ] Coins popup: HTMX swap to show current balance and trigger award
- [ ] Avatar page: render avatar options, POST selection on change

#### Teacher Pages
- [ ] Attendance page: student list with checkboxes, POST on submit, HTMX swap for confirmation
- [ ] Reports page: table of student data, "Export CSV" button that hits the export route
- [ ] Schedule page: calendar or date list, button to mark a double day

#### Admin Pages
- [ ] Classrooms page: list of classrooms, form to create one, assign teacher/students UI
- [ ] All reports page: aggregate view across all classrooms

---

## Where the Tracks Merge

Once Track A has working middleware and Track B has all stub pages:

- [ ] Person 2 removes `fakeUser` and replaces with real session lookup from context
- [ ] Person 2 wraps every route with `requireLogin` and `requireRole` as appropriate
- [ ] Smoke test: log in as each role and confirm the right pages are accessible

---

## Final Packaging

- [ ] Test on a clean machine (no Go installed) to verify the binary is self-contained
- [ ] Embed templates and static files using Go's `embed` package so they're bundled into the binary
- [ ] Build for Windows: `GOOS=windows GOARCH=amd64 go build -o app.exe`
- [ ] Confirm the `.exe` opens the server and can be visited at `localhost:<port>`
- [ ] Optionally: auto-open the browser on launch using `os/exec` to call `start` (Windows) or `open` (Mac)
- [ ] Write a short README with instructions for running the `.exe`