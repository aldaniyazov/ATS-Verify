# VISION
**ATS-Verify** is a comprehensive microservice platform designed to streamline logistics and customs verification processes. It provides distinct workflows for different stakeholders (ATS, Customs, Marketplaces) to manage parcel data, track shipments, analyze risks, and verify goods against declarations using intelligent text analysis.

# SYSTEM RULES (The 5 Rules of ATS-Verify)
> **Constraint:** This system relies on strict role-based logic and specific data validation rules.

## 1. Role Logic
*   **Admin:** Omni-view, Risk Management (IIN/BIN coloring), Global CSV Downloads.
*   **Paid Users:** Tracking + Invoice Matching (IMEI).
*   **ATS Staff:** Track & Trace internal check (Used/Not Used).
*   **Customs:** View Table, Mark "Used", Analysis (IMEI in PDF).
*   **Marketplace (uploader_{mp}):** Upload CSV. Suffix determines `marketplace` field (e.g., `uploader_wb` -> "Wildberries").

## 2. Ingestion Logic (Marketplace CSV)
**Input:** `marketplace | country | brand | name | track number | SNT | date`
**Deduplication Rules:**
1.  **New Track Number:** `INSERT` (status: `used = false`).
2.  **Existing Track (Used=False):** `UPDATE` (Overwrite data).
3.  **Existing Track (Used=True):** `ERROR` ("Track already used").
    *   *Exception:* If track exists but not used, notify user and allow overwrite.

## 3. Analysis Logic (Customs/Paid)
**IMEI PDF Verification:**
1.  **Input:** CSV (Imei1..4) + PDF (Declaration).
2.  **Process:**
    *   Extract unique IMEIs from CSV.
    *   Extract text from PDF (Graph 31).
    *   Search: Find 14-character IMEI substrings within 15-character sequences (ignoring check digit).
3.  **Output:** Report (Found X/Y, usage matrix).

**Risk Analysis (Admin):**
1.  **Input:** CSV (`Date | AppId | IIN/BIN | doc | User | Org | Status | Reject | Reason`).
2.  **Metrics:**
    *   Doc reuse count.
    *   Doc usage across different IIN/BIN.
    *   Flip-flop status (Accepted -> Rejected).
    *   High-frequency IIN/BIN.
3.  **Action:** Assign Risk Level (Red/Yellow/Green) to IIN/BIN.

## 4. Tracking Logic
*   **Providers:** Kazpost (API), CDEK (API).
*   **Scope:** Available to all users.

# USER STORIES
1.  **As a Marketplace Emp**, I want to upload my daily CSV so that parcels are registered.
2.  **As a Customs Officer**, I want to upload a Declaration PDF and a range of IMEIs to check which ones are validly declared.
3.  **As an Admin**, I want to color-code IINs (Red/Yellow/Green) based on suspicious activity patterns.
4.  **As an ATS Emp**, I want to check a track number to see if it has been scanned/used.

# DESIGN TOKENS (Modern/Premium)
*   **Palette:** Deep Navy (Background), Electric Blue (Accents), Glassmorphism (Cards), Semantic Red/Yellow/Green (Risks).
*   **Typography:** Inter / Outfit.
*   **Vibe:** Professional Dashboard, High Density Data, Smooth Transitions.
