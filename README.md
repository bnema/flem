# Setup

Install pnpm globally:

```bash
npm i -g pnpm
```

Install typescript globally:

```bash
npm i -g typescript
```

Install dependencies (PNPM mandatory):

```bash
pnpm i
```

Try to run the turbo dev command:

```bash
pnpm dev

# or for specific app : 
pnpm dev --filter <app>
```

API should be accessible at http://localhost:3333/v1/healthcheck
Client should be accessible at http://localhost:3000


Add a package to a specific app :

```bash
pnpm add <package> --filter <app>

# example :
pnpm add @openai/openai-api --filter api
```
