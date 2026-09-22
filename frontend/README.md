# Frontend

React and TypeScript single-page app that consumes the backend's
`/api/calculate` endpoint. Built with Vite; tested with Vitest and React
Testing Library.

## Setup

Requires Node 20 or newer.

```bash
cd frontend
npm install
```

By default the app calls the backend at `http://localhost:8080`. To change
that, copy the example env file and edit it:

```bash
cp .env.example .env
```

`VITE_API_URL` is the only variable it reads.

## Running

```bash
npm run dev        # dev server at http://localhost:5173
npm run build       # typecheck, then production build into dist/
npm run preview     # serve the production build locally
```

The backend must be running for calculations to succeed. If it is not
reachable, the app shows an explicit message instead of failing silently.

## Running with Docker

```bash
docker build --build-arg VITE_API_URL=http://localhost:8080 -t sezzle-frontend .
docker run -p 3000:80 sezzle-frontend
```

The `--build-arg` is required at build time, not run time: Vite reads
`VITE_API_URL` while bundling and inlines it into the JavaScript output, so
it cannot be changed later by setting an environment variable on the running
container. If it is omitted, the Dockerfile falls back to
`http://localhost:8080`.

The `Dockerfile` is a two-stage build: `npm ci` and `npm run build` run in a
`node:22-alpine` stage, then only the compiled `dist/` output is copied into
an `nginx:1.27-alpine` image, which serves the static files. `nginx.conf`
redirects unknown paths back to `index.html`, which is required for a
single-page app's client-side routing to work on a hard refresh.

## Project structure

```
src/
  api/client.ts             the only file that knows about HTTP
  components/
    Calculator.tsx          form state, validation, submission
    Display.tsx             result / error output
    OperationPad.tsx         operation selector
  operations.ts             operation catalogue (ids match the backend)
  validation.ts             operand parsing
  format.ts                 display formatting
```

## Notes on the implementation

The backend's `/api/calculate` is stateless: one operation, up to two
operands, one result. The UI mirrors that directly - pick an operation, fill
in the operands, submit - rather than simulating a running calculator tape,
which the API does not support.

Validation is split by responsibility. The frontend only checks that a
field is present and is a number, so it never sends a request that cannot
succeed. Rules like division by zero stay on the backend; the UI just
displays whatever message comes back. This avoids maintaining the same
rules in two places.

Operand inputs are `type="text"` with `inputMode="decimal"`, not
`type="number"`. A native number input reports invalid text as an empty
string, which makes "you typed letters" indistinguishable from "you typed
nothing" and produces a worse error message.

Square root hides the second field and omits `b` from the request entirely,
since the backend treats a present `b: 0` as an actual operand, not as "no
value".

Displayed results are rounded to 12 significant digits. Without that,
`0.1 + 0.2` would render as `0.30000000000000004`, a side effect of binary
floating-point representation rather than a bug in the arithmetic.

## Tests

```bash
npm test               # run once
npm run test:watch     # watch mode
npm run test:coverage   # coverage report (text + HTML in coverage/)
npm run typecheck      # tsc --noEmit
```

### What the tests check

`src/validation.test.ts` - operand parsing: integers, decimals, negatives,
whitespace, and rejected input like letters or multiple decimal points.

`src/format.test.ts` - result formatting, including the floating-point
rounding, and the human-readable expression labels shown next to a result.

`src/api/client.test.ts` - the API client function. `fetch` is replaced with
a mock (`vi.stubGlobal('fetch', ...)`), so these tests never make a real
network call. They check the request body shape, that `b` is omitted for
`sqrt`, that the backend's `{"error": "..."}` envelope is surfaced correctly,
and that an unreachable backend produces a clear error rather than an
unhandled exception.

`src/components/Calculator.test.tsx` - renders the actual `Calculator`
component with React Testing Library and simulates a user typing and
clicking with `@testing-library/user-event`. The API client is mocked here
too (`vi.mock('../api/client', ...)`), so these tests are only checking the
component's behavior: switching operations, showing results, validation
messages, error messages from a failed API call, and the Clear button.
Elements are queried the same way a user would find them - by label text and
button role - rather than by CSS class, so the tests keep working if styling
changes.

One setup detail worth calling out: `src/test/setup.ts` calls
`cleanup()` after every test. React Testing Library does this automatically
under Jest, but Vitest needs it registered explicitly; without it, each
test's rendered component stays in the DOM and the next test's queries
match leftover elements from previous tests.
