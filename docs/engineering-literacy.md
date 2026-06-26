# Understanding Your Engineering Team

Panduan buat PM/Scrum Master yang ingin lebih memahami cara kerja developer, QA, dan engineering team secara umum. Bukan buat jadi engineer, tapi buat komunikasi yang lebih baik dan decision-making yang lebih informed.

---

## Tools yang Bantu PM Lihat "Dunia Developer"

### 1. Apa yang Tim Sebenarnya Kerjakan

| Tool | Insight yang Didapat |
|------|---------------------|
| `pm_flow_metrics(board_id)` | WIP (work in progress), throughput, cycle time, lead time |
| `pm_github_activity(days:7)` | Commits, PRs merged, issues closed minggu ini |
| `pm_github_prs` | Open PRs: siapa review siapa, berapa lama |
| `pm_github_pr_metrics` | Avg PR age, stale PR count |
| `pm_time_report(days:7)` | Siapa kerja di mana, berapa lama |
| `jira_workload(project)` | Distribusi issue per orang |
| `pm_resource_utilization(board_id)` | Assigned vs done vs WIP vs blocked per member |

### 2. Kesehatan Proses Development

| Tool | Insight yang Didapat |
|------|---------------------|
| `pm_tech_debt` | Daftar tech debt: code, architecture, testing, infra |
| `pm_tech_debt_ratio` | Rasio bugs/debt vs fitur. > 20% = warning |
| `pm_tech_debt_budget(board_id)` | Berapa % sprint yang harusnya buat bayar debt |
| `pm_anti_patterns(board_id)` | Hero culture, zombie sprint, scope creep |
| `pm_commitment_check(board_id)` | Apakah tim overcommit? |
| `jira_trace_branch(key)` | Status implementasi: ada branch? Ada PR? Udah merged? |

### 3. Incident & Production Health

| Tool | Insight yang Didapat |
|------|---------------------|
| `pm_incidents` | Production incidents: severity, status, assignee |
| `pm_incident_summary` | Summary by status/urgency + avg resolution time |
| `pm_oncall` | Siapa yang on-call sekarang |

---

## Konsep Engineering yang PM Perlu Pahami

### Work In Progress (WIP)

**Apa:** Jumlah item yang sedang dikerjakan bersamaan.

**Kenapa penting:** WIP tinggi = context switching tinggi = quality turun + cycle time naik. Kalau tim punya 5 orang tapi 15 item In Progress, itu masalah.

**Tool:** `pm_flow_metrics(board_id)` — lihat WIP count dan deteksi flow problems.

**Tanda bahaya:**
- WIP > 2x jumlah developer = terlalu banyak paralel
- Cycle time naik padahal throughput turun = bottleneck

**Yang bisa SM lakukan:** "Tim, kita coba limit WIP ke 2 item per orang. Selesaikan yang ada sebelum mulai yang baru."

---

### Cycle Time vs Lead Time

**Cycle Time:** Dari mulai dikerjakan sampai selesai (In Progress → Done).
**Lead Time:** Dari dibuat sampai selesai (Created → Done).

**Kenapa penting:** Cycle time naik = ada yang stuck. Lead time tinggi padahal cycle time rendah = item terlalu lama di backlog sebelum dimulai.

**Tool:** `pm_flow_metrics(board_id)` — keduanya dihitung otomatis.

**Target sehat:**
- Cycle time < 3 hari buat task biasa
- Lead time < 1 sprint buat story

---

### Tech Debt

**Apa:** Shortcut teknis yang diambil buat kirim lebih cepat. Bukan selalu buruk, tapi harus dikelola.

**Kenapa PM harus care:** Tech debt yang nggak dikelola = velocity yang pelan-pelan turun. Feature yang harusnya 3 hari jadi 2 minggu karena codebase fragile.

**Tool:**
- `pm_tech_debt` — daftar semua debt yang tercatat
- `pm_tech_debt_ratio` — kalau > 20% sprint items adalah bugs/debt, itu sinyal buruk
- `pm_tech_debt_budget(board_id)` — rekomendasi berapa % sprint buat address debt

**Percakapan yang benar:**
- "Tim, tech debt ratio kita 35%. Kita alokasikan 20% sprint capacity buat bayar. Pilih yang impact-nya paling tinggi."
- Bukan: "Kita nggak punya waktu buat refactoring."

---

### Pull Request (PR) / Code Review

**Apa:** Developer submit kode, tim lain review sebelum masuk ke main branch. Quality gate.

**Kenapa PM harus care:** PR yang lama di-review = bottleneck yang invisible. Kalau avg PR age > 2 hari, itu blocking delivery.

**Tool:**
- `pm_github_prs` — list semua open PR dengan age
- `pm_github_pr_metrics` — avg age, stale count
- `jira_trace_branch(key:"PROJ-123")` — status implementasi dari Jira ticket ke code

**Tanda bahaya:**
- PR open > 3 hari tanpa review
- Satu orang jadi bottleneck reviewer
- PR yang besar (>500 lines) — harder to review, more bugs

**Yang bisa SM lakukan:** "Kita punya 5 PR yang udah 4 hari tanpa review. Bisa kita prioritaskan review di pagi ini?"

---

### Overcommitment

**Apa:** Tim ambil lebih banyak dari yang bisa diselesaikan dalam sprint.

**Kenapa penting:** Overcommitment berulang = team burnout, carryover naik, predictability turun. Stakeholder kecewa karena expectation nggak ketemu reality.

**Tool:**
- `pm_commitment_check(board_id)` — compare sprint items vs historical completion rate
- `pm_commitment_report(board_id)` — delivery rate over time
- `pm_capacity_plan(board_id)` — rekomendasi berapa item based on velocity

**Percakapan yang benar:**
- "Historically kita deliver 80% dari commitment. Sprint ini kita ambil 30 items padahal avg kita 22. Kita reduce ke 24 supaya realistic."
- Bukan: "Kita harus selesaikan semuanya."

---

### Hero Culture

**Apa:** Satu orang menyelesaikan mayoritas work. Bus factor = 1.

**Kenapa berbahaya:** Kalau orang itu cuti/resign, tim collapse. Juga: burnout risk tinggi buat si hero.

**Tool:**
- `pm_anti_patterns(board_id)` — detects if one person does > 50% of completed items
- `pm_resource_utilization(board_id)` — workload distribution table
- `jira_workload(project)` — issue count per person

**Yang bisa SM lakukan:**
- Pair junior dengan hero buat knowledge transfer
- Redistribute tasks di sprint planning
- "Mas X, sprint ini kamu udah di 8 items. Bisa kita pindahkan 2 ke yang lain?"

---

### Definition of Done (DoD)

**Apa:** Checklist yang harus dipenuhi sebelum item dianggap "selesai". Termasuk: code reviewed, tests written, deployed to staging, documented.

**Kenapa PM harus care:** Tanpa DoD, "done" artinya beda buat setiap orang. Developer bilang done tapi belum di-test. QA belum bisa verify. Deployment belum jalan.

**Tool:**
- `pm_dod` — lihat/manage DoD checklist
- `pm_dod(action:"add", item:"Unit tests written", category:"testing")` — tambah item

**DoD yang baik untuk software team:**
- Code reviewed by at least 1 person
- Unit tests written and passing
- Integration tested
- No critical/high bugs open
- Deployed to staging
- PO accepted

---

### QA & Testing Vocabulary

| Istilah | Artinya | Implikasi buat PM |
|---------|---------|-------------------|
| Unit Test | Test per function/method | Fast feedback, developer write ini |
| Integration Test | Test interaksi antar komponen | Butuh waktu lebih, tapi catch bugs yang unit miss |
| E2E Test | Test full flow dari user perspective | Paling lambat, paling realistic |
| Regression | Bug lama muncul lagi | Sinyal test coverage kurang |
| Flaky Test | Test yang kadang pass kadang fail | Wasted time debugging false failures |
| Test Coverage | % kode yang di-cover test | 80%+ itu target sehat, 100% nggak realistis |
| Staging | Environment mirip production buat testing | Deploy sini dulu sebelum live |

**Tool yang relevan:**
- `pm_tech_debt_add(title:"Flaky tests in payment module", category:"testing", impact:"high")` — track testing issues as tech debt
- `pm_dod(action:"add", item:"QA sign-off on staging", category:"testing")` — enforce QA in DoD

---

### Deployment & Release Vocabulary

| Istilah | Artinya | Implikasi buat PM |
|---------|---------|-------------------|
| CI/CD | Continuous Integration/Deployment | Kode otomatis di-test dan di-deploy |
| Rollback | Kembalikan ke versi sebelumnya | Safety net kalau deploy gagal |
| Feature Flag | Toggle fitur on/off tanpa deploy ulang | Bisa ship code tapi belum "aktif" |
| Hotfix | Perbaikan darurat langsung ke production | Bypass normal sprint flow |
| Canary/Blue-Green | Deploy ke sebagian user dulu | Reduce blast radius |

**Tool yang relevan:**
- `pm_github_releases` — recent releases, correlate with sprint delivery
- `pm_incidents` — production issues post-deploy
- `pm_release_notes(board_id)` — auto-generate what shipped

---

## Metrik yang PM Harus Pantau (Bukan Velocity Saja)

| Metrik | Tool | Sehat | Warning |
|--------|------|-------|---------|
| WIP | `pm_flow_metrics` | < 2x team size | > 3x team size |
| Cycle Time | `pm_flow_metrics` | < 3 hari | > 5 hari |
| PR Age | `pm_github_pr_metrics` | < 1 hari | > 3 hari |
| Tech Debt Ratio | `pm_tech_debt_ratio` | < 15% | > 25% |
| Completion Rate | `pm_sprint_health` | > 80% | < 60% |
| Carryover | `pm_sprint_compare` | < 15% | > 30% |
| Blocker Age | `pm_blocker_aging` | < 2 hari avg | > 5 hari avg |
| Team Balance | `pm_resource_utilization` | Even spread | 1 person > 40% |

---

## Pertanyaan yang Bagus buat SM Tanyakan di Standup

Bukan "ada blocker?" — terlalu generic. Coba yang lebih specific:

1. "PR kamu udah berapa lama open? Butuh reviewer?"
2. "WIP kita udah 12 items. Yang mana yang bisa kita finish hari ini?"
3. "Item ini udah In Progress 4 hari. Ada yang stuck? Perlu dibreak jadi smaller task?"
4. "Ada dependency ke tim lain yang belum resolve?"
5. "Tech debt mana yang bikin kerjaan ini jadi lebih lama?"

Tool buat prep: `pm_standup_prep(board_id)` — auto-generate these talking points.

---

## Cara Baca GitHub/GitLab Activity (untuk PM Non-Teknis)

```
pm_github_activity(days:7)
```

Output yang perlu diperhatikan:
- **Commits:** Banyak commit = aktif. Tapi commit banyak bukan berarti productive — bisa juga banyak fix kecil karena rushing.
- **PRs merged:** Ini yang paling relevan. PR merged = kode masuk ke main = fitur progress.
- **PRs open:** Banyak open PR yang lama = review bottleneck.
- **Issues closed:** Correlate dengan Jira items done.

---

## Red Flags yang PM Harus Escalate

| Signal | Tool untuk Verify | Action |
|--------|-------------------|--------|
| 1 orang selesaikan > 50% items | `pm_anti_patterns` | Redistribute, pair programming |
| Tech debt ratio > 30% | `pm_tech_debt_ratio` | Allocate 25% sprint buat debt |
| Avg cycle time > 5 hari | `pm_flow_metrics` | Cari bottleneck: review? testing? unclear spec? |
| PR age > 3 hari average | `pm_github_pr_metrics` | Prioritize reviews, add reviewers |
| Incidents naik | `pm_incident_summary` | Discuss quality in retro |
| Sprint completion < 60% | `pm_sprint_health` | Reduce commitment, fix estimation |
| Same items carried 3+ sprints | `pm_anti_patterns` (zombie) | Break down items, reassign, or descope |

---

## Learning Path: From PM to Technical PM

### Week 1-2: Observasi
- Run `pm_flow_metrics(board_id)` daily — learn to read WIP, cycle time
- Run `pm_github_prs` — understand PR review flow
- Ask engineer: "Walk me through how this feature goes from Jira ticket to production"

### Week 3-4: Vocabulary
- Read this guide's vocabulary sections
- Run `pm_tech_debt` — understand what debt looks like
- Ask engineer: "What's the biggest tech debt slowing us down?"

### Month 2: Data-Driven Questions
- Use `pm_github_pr_metrics` to spot review bottlenecks
- Use `pm_anti_patterns` monthly to detect dysfunction early
- Use `pm_commitment_check` before accepting sprint scope

### Month 3+: Proactive Engineering Partnership
- Use `pm_tech_debt_budget` to propose debt allocation
- Use `jira_trace_branch` to verify "done" means actually deployed
- Use `pm_resource_utilization` to balance load before it becomes burnout

---

## Quick Reference Card

**"Tim merasa overwhelmed"**
```
pm_commitment_check(board_id:X)  -> are we overcommitted?
pm_resource_utilization(board_id:X)  -> who is overloaded?
pm_flow_metrics(board_id:X)  -> is WIP too high?
```

**"Delivery makin lambat"**
```
pm_flow_metrics(board_id:X)  -> cycle time trend
pm_github_pr_metrics  -> review bottleneck?
pm_tech_debt_ratio  -> debt slowing us?
pm_anti_patterns(board_id:X)  -> structural issues?
```

**"Kualitas turun"**
```
pm_incidents  -> production issues
pm_tech_debt  -> testing gaps
pm_dod  -> is DoD being followed?
pm_github_prs  -> are PRs getting proper review?
```

**"Mau hire/scale team"**
```
pm_resource_utilization(board_id:X)  -> current load
pm_capacity_plan(board_id:X)  -> capacity vs velocity
pm_maturity_assessment(board_id:X)  -> team readiness
portfolio_workload  -> cross-project load
```
