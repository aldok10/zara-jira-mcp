# Reporting to Management & Cross-Team Communication

Panduan buat PM/Scrum Master yang perlu report ke atasan, Product Owner, atau stakeholder lintas divisi.

## Prinsip Utama

1. **Management nggak butuh detail teknis.** Mereka butuh: status, risiko, timeline, dan apa yang perlu mereka keputuskan.
2. **Data beats opinion.** Jangan bilang "sprint berjalan lancar" — tunjukkan health score 82/100.
3. **Escalate early, not late.** Blocker yang udah 5 hari itu bukan update lagi — itu escalation.
4. **Different audience, different report.** Jangan kirim sprint burndown ke VP.

---

## Tool Matrix: Siapa Butuh Apa

| Audience | Tool | Apa yang Didapat |
|----------|------|------------------|
| VP / C-Level | `pm_exec_report` | Status 1 baris, delivered value, risks, team health. Under 250 words, zero jargon |
| Product Owner | `pm_goal_check` | Apakah sprint goal on track? Progress vs key results |
| Product Owner | `pm_scope_creep` | Item yang masuk mid-sprint tanpa approval |
| Head of Engineering | `pm_sprint_health` | Health score 0-100 dengan breakdown |
| Head of Engineering | `pm_anti_patterns` | Deteksi dini: hero culture, zombie sprint, scope creep |
| Cross-team Dependencies | `pm_dependencies` | Dependency map: siapa nunggu siapa |
| Cross-team Dependencies | `portfolio_blockers` | Semua blocker lintas project |
| All Stakeholders | `pm_weekly_digest` | AI summary mingguan: wins, concerns, next focus |
| All Stakeholders | `pm_release_notes` | Apa yang shipped: features, bugs fixed, improvements |
| Board/Steering Committee | `portfolio_summary` | AI executive summary lintas project |

---

## Skenario: Lapor ke Atasan (VP/Director)

**Kapan:** Weekly atau bi-weekly

```
pm_exec_report(board_id:X)
```

Output berisi:
- **Status**: On Track / Watch / At Risk + alasan
- **Delivered This Sprint**: Value bisnis, bukan nomor tiket
- **Coming Next**: Apa yang stakeholder bisa expect
- **Risks & Blockers**: Apa yang bisa delay + apa yang butuh keputusan dari atas
- **Team Health**: 1 kalimat sinyal

Bisa langsung kirim ke Lark/Slack:
```
pm_exec_report(board_id:X, send_to_lark:true)
```

---

## Skenario: Update ke Product Owner

**Kapan:** Daily atau mid-sprint

### "Apakah sprint goal bakal tercapai?"
```
pm_goal_check(board_id:X)
```
AI evaluasi berdasarkan progress aktual vs key results yang di-set di planning.

### "Ada scope creep nggak?"
```
pm_scope_creep(board_id:X)
```
Mendeteksi item yang ditambah setelah sprint start. PO perlu tahu kalau scope berubah tanpa mereka sadar.

### "Forecast realitis kapan selesai?"
```
pm_forecast(board_id:X, remaining_items:12)
```
Monte Carlo 10,000 simulasi → "50% chance done Thursday, 85% chance done next Monday."

PO bisa pakai angka ini buat manage expectation stakeholder mereka.

---

## Skenario: Escalation ke Management

**Kapan:** Blocker > 3 hari, risk critical, health < 50

### Auto-escalate (set and forget)
```
pm_escalate(board_id:X)
```
Otomatis kirim alert ke Lark/Slack kalau:
- Risk critical/high sudah > 3 hari tanpa mitigation
- Blocker aktif > 3 hari
- Sprint health score < 50

### Manual escalation dengan context
```
pm_impediment_aging
```
Report berisi: semua blocker aktif + berapa hari, avg resolution time, mana yang chronic.

Kirim ini ke management dengan note: "These need executive action."

---

## Skenario: Cross-Team Coordination

### "Team lain blocking kita"
```
pm_record_dependency(from_issue:"TEAM-A-123", to_issue:"TEAM-B-456", type:"blocked_by", description:"Waiting for API v2 from Platform team")
```

Lalu saat meeting lintas tim:
```
pm_dependencies
portfolio_blockers
```

### "Portfolio-wide view buat steering committee"
```
portfolio_summary
```
AI-generated executive summary lintas semua project: mana yang sehat, mana yang butuh attention.

---

## Skenario: Report Periodik ke Stakeholder

### Weekly digest (auto-generated)
```
pm_weekly_digest(board_id:X, send_to_lark:true)
```
AI merangkum: meetings, decisions, risks resolved, blockers, sprint progress.

### End-of-sprint report
```
pm_scorecard(board_id:X)
pm_release_notes(board_id:X, send_to_lark:true)
```
Scorecard buat internal (team performance), release notes buat stakeholder (what shipped).

### Stakeholder satisfaction tracking
```
pm_stakeholder_pulse(stakeholder:"Product Owner", score:4, sprint_name:"Sprint 23", feedback:"Happy with velocity, concerned about quality")
pm_stakeholder_trend
```
Track apakah stakeholder makin puas atau makin frustrated over time.

---

## Skenario: Justify SM Value ke Management

"Apa sih yang SM bikin?" — pertanyaan klasik.

```
pm_sm_impact(sprint_name:"Sprint 23")
```

Output:
- Blockers resolved + avg resolution time
- Risks mitigated
- Pending action items (fewer = better follow-through)
- Retros facilitated
- Impact score

Kirim ini di performance review atau monthly report.

---

## Skenario: Team Maturity Assessment

Buat diskusi 1-on-1 dengan management tentang team readiness:

```
pm_maturity_assessment(board_id:X)
```

Based on data: velocity stability, blocker resolution speed, risk awareness, completion consistency. Output stage: Forming/Storming/Norming/Performing.

---

## Skenario: OKR Alignment Report

Hubungkan sprint work ke business objectives:

```
pm_outcome_map(board_id:X, objective:"Reduce churn by 20%", key_results:"Ship retention dashboard\nAutomate win-back emails")
pm_outcome_history
```

Saat ditanya "sprint ini kontribusi ke OKR mana?" — ada jawabannya.

---

## Communication Channels

Semua report bisa di-route ke channel yang tepat:

| Channel | Tool | Use Case |
|---------|------|----------|
| Lark | `jira_notify_lark` | Lark groups |
| Slack | `slack_send` / `slack_notify_team` | Slack channels |
| Discord | `discord_send` | Discord servers |
| Telegram | `telegram_send` | Telegram groups |
| Teams | `teams_send` | Microsoft Teams |
| Email | `email_send` | Formal reports ke stakeholder |
| Confluence | `confluence_create_page` | Documented reports |
| All at once | `broadcast` | Critical announcements |
| Smart routing | `notify_routed` | Auto-pick channel based on severity |

---

## Best Practices

1. **Set sprint goals with key results** (`pm_set_sprint_goal`) — tanpa ini, `pm_goal_check` dan `pm_exec_report` kurang context.
2. **Snapshot every sprint end** (`pm_snapshot_sprint`) — tanpa history, forecasting nggak jalan.
3. **Record decisions immediately** (`pm_record_decision`) — 3 bulan dari sekarang, nggak ada yang ingat kenapa keputusan itu diambil.
4. **Weekly risk scan** (`pm_auto_detect_risks`) — jangan tunggu masalah jadi krisis.
5. **Track stakeholder pulse** (`pm_stakeholder_pulse`) — early warning kalau relationship memburuk.
