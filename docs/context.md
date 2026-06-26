# Context: Kenapa Project Ini Ada

## Masalah yang Diselesaikan

### PM/SM Menghabiskan 40% Waktu untuk Hal yang Bisa Diotomasi

Riset PwC (2025): rata-rata PM kehilangan **7.4 jam/minggu** untuk task administratif — status reporting, data gathering, copy-paste antar tools. ZipDo: 78% PM report faster task completion setelah pakai AI, dan AI bisa automate 35% status report generation (hemat 5-7 jam/minggu).

PMI (2026) bahkan publish standar global pertama untuk AI in Project Work — sinyal bahwa ini bukan trend, ini shift permanen.

### Jira Itu Database, Bukan Brain

Jira menyimpan data. Tapi dia nggak:
- Ingat keputusan sprint lalu
- Deteksi bahwa 1 orang carry 60% work (hero culture)
- Predict kapan backlog selesai berdasarkan historical throughput
- Generate report yang different per audience
- Escalate otomatis kalau blocker udah 5 hari

PM harus jadi "glue" antara data Jira dan insight yang dibutuhkan stakeholder. Itu manual, repetitive, dan error-prone.

### AI Agents Sudah Drive, Bukan Cuma Suggest

Quire (2026): "AI is starting to drive: it reads the project state, takes multi-step action, and stops to check in only when something matters."

Gartner: 40% enterprise apps akan punya task-specific AI agents by 2026.

Beda 2023 vs 2026:
- 2023: AI suggest summary ticket → user copy-paste manual
- 2026: AI baca sprint state → run forecast → detect risk → generate report → kirim ke Slack → done

zara-jira-mcp adalah implementasi dari shift ini. PM bicara natural language, AI execute multi-step PM workflows.

---

## Data Points (Research-Backed)

| Finding | Source | Implication |
|---------|--------|-------------|
| PM save 7.4 hours/week with AI automation | PwC 2025 | Tool harus eliminate manual data gathering |
| 78% PM report faster task completion with AI | ZipDo 2025 | Adoption barrier is setup, not value |
| 25% AI adoption increase → 1.5% throughput decrease + 7.2% stability decrease | DORA 2025 (initial) | AI tanpa proper workflow = amplify dysfunction |
| Later DORA finding: positive relationship when teams learn where/when AI useful | DORA 2025 (mature) | Right tool + right workflow = real gains |
| Developers merged 98% more PRs with AI, org delivery flat | DORA/Typo 2025 | Individual productivity ≠ team delivery. Need system-level view |
| Devs using AI took 19% longer while believing 20% faster (39-point gap) | METR 2025 | Perception ≠ reality. Need objective metrics |
| 31% projects don't meet goals | PMI 2025 | Sprint goal tracking is critical, not optional |
| AI reduces admin load 30-50% for SM | Techademy coaching data | SM can shift time to coaching + impediment removal |
| Org waste 9.9% of every dollar due to poor project performance | PMI | PM tools that prevent waste pay for themselves |
| Only 52% teams achieve sprint goals | Scrum Alliance | Goal-focused PM tools have massive headroom |
| 15-20% sprint capacity should go to tech debt | 6 industry sources | Tools must track + recommend debt allocation |
| AI omissions reduce team trust → hurt performance | Frontiers in Psychology 2025 | Transparency in AI-generated reports is non-negotiable |

---

## Peran SM yang Berubah vs Tetap

### Yang AI Ambil Alih (SM Nggak Perlu Lagi)

- Generate sprint status manually
- Copy-paste ticket data ke slides
- Hitung velocity manual dari Jira
- Tulis meeting notes dari scratch
- Monitor sprint progress secara manual
- Track action items di spreadsheet

### Yang Tetap Butuh Manusia (SM Makin Penting)

- Coaching 1-on-1 (emotional intelligence)
- Conflict resolution (reading the room)
- Organizational change management
- Building trust + psychological safety
- Strategic decision-making
- Cross-team relationship building
- Facilitating difficult conversations

### SM Role Shift

**Before AI:** 60% admin, 40% people work
**After AI:** 10% admin (oversee AI output), 90% people work

Sumber: "I am not seeing traditional scrum masters replaced. What I am seeing is the role amplified, restructured, and re-prioritised." — Techademy, 2026

---

## Positioning zara-jira-mcp

**Bukan:** pengganti SM.
**Tapi:** amplifier yang eliminates 50% busywork SM, supaya SM bisa focus ke hal yang AI nggak bisa: people.

**Untuk siapa:**
1. SM yang mau shift dari "status reporter" ke "team coach"
2. PM yang overwhelmed dengan admin tapi punya skill people
3. Engineering manager yang butuh visibility tanpa micromanage
4. Solo dev / small team yang nggak punya dedicated SM

**Value proposition:**
> "Stop wasting 7 hours/week on reports. Let AI handle the data, you handle the humans."
