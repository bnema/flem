
## Global

- [x] Setup a new project with Turborepo
- [x] Added apps/api and apps/web
- [x] Test both apps in dev and build

## Backend

### API Global

- [ ] Setup Prisma + Postgres
- [x] Configure CORS
- [ ] Configure basic user auth with KeyCloak

### HTTP Server

- [x] Setup Fastify
- [x] Create /v1/healthcheck endpoint
- [ ] Create endpoints for CRUD operations on users, movies and tv shows
- [ ] Create security middleware for auth & roles check
- [ ] Create error handling middleware
- [ ] Create CRUD endpoints for GPT-3.5 responses (movies and tv shows)

### OPENAI API Client

- [ ] Setup communication with OpenAI API (GPT-3.5)
- [ ] Setup new prompt for movies and tv shows
- [ ] Handle responses from OpenAI API (Create a json and return it to the client with the HTTP Server)