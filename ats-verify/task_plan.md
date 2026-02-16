# Current Stack: Go (Backend) + React/Vite (Frontend) + PostgreSQL

## Phase 1: [B]lueprint (Logic & Requirements)
- [ ] Define DB Schema (Users, Roles, Parcels, TrackingEvents, Risks).
- [ ] Define API Contract (OpenAPI/Swagger).
- [ ] Define "Definition of Done".

## Phase 2: [L]ink (Interfaces & Environment)
- [ ] Setup PostgreSQL + pgvector (if needed for search).
- [ ] Setup Go Modules & Project Structure.
- [ ] Setup Frontend (Vite + Tailwind).

## Phase 3: [A]rchitect (Core Implementation)
- [ ] **Auth Service:** JWT, Roles (Admin, Paid, ATS, Customs, Marketplace).
- [ ] **Ingestion Service:** CSV Parsing, Deduplication Logic.
- [ ] **Tracking Service:** Kazpost/CDEK Integration.
- [ ] **Analysis Service:** PDF Parsing (IMEI), Risk Analysis (IIN/BIN).

## Phase 4: [S]tyle & Structure (Refinement)
- [ ] Apply "Premium" Design Tokens to Frontend.
- [ ] Implement Dashboard Widgets (Risk Colors, Usage Stats).

## Phase 5: [T]rigger (Deployment)
- [ ] Docker Compose Setup.
- [ ] CI/CD Pipeline.
