package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// TechGlossary explains technical concepts in PM-friendly language.
func (h *Handlers) TechGlossary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	term := req.GetString("term", "")

	glossary := map[string]string{
		"api":            "API (Application Programming Interface): How different software systems talk to each other. Like a waiter between kitchen (backend) and customer (frontend). When devs say 'the API is down', it means systems can't communicate.",
		"ci/cd":          "CI/CD (Continuous Integration/Deployment): Automated pipeline that builds, tests, and deploys code. When devs push code, CI automatically checks it. CD auto-deploys approved code. Key metric: how often can we safely release.",
		"pr":             "PR (Pull Request) / MR (Merge Request): A developer's request to merge their code changes into the main codebase. Other devs review it before approval. Long PR review time = bottleneck.",
		"technical debt":  "Technical Debt: Shortcuts taken in code that work now but cost more later. Like credit card debt — small now, expensive if ignored. PM should allocate 15-20% sprint to pay it down.",
		"sprint":         "Sprint: Fixed time period (usually 2 weeks) where team commits to deliver specific work. Has a goal, a backlog of items, and ceremonies (planning, standup, review, retro).",
		"story points":   "Story Points: Relative complexity estimate (not time!). A 5-point story is roughly 2.5x harder than a 2-point story. Used for planning, NOT for measuring productivity. Never compare points between teams.",
		"blocker":        "Blocker: Something preventing a developer from making progress. Could be: waiting for another team, unclear requirement, infrastructure issue, or missing access.",
		"wip":            "WIP (Work In Progress): How many tasks someone is working on simultaneously. Research shows optimal is 1-2 items. More WIP = more context switching = slower delivery.",
		"cycle time":     "Cycle Time: Days from 'started working' to 'done'. Short cycle time (1-5 days) = healthy flow. Long cycle time (>10 days) = something is stuck or story is too big.",
		"deployment":     "Deployment: Releasing code to production (live users). Modern teams deploy multiple times per day. Failed deployments = change failure rate (DORA metric).",
		"rollback":       "Rollback: Reverting a deployment that caused problems. Quick rollback = team has safety nets. PM should ask: 'can we safely undo this if it breaks?'",
		"microservice":   "Microservice: Architecture where the system is split into small independent services. Each team owns a service. Dependency between services = coordination overhead for PM.",
		"monolith":       "Monolith: Single large application. Simpler to deploy but harder to scale. Teams step on each other's code. PM impact: merge conflicts, longer testing cycles.",
		"staging":        "Staging: A copy of production used for final testing before release. If staging breaks, production is safe. PM should ask: 'has this been tested in staging?'",
		"regression":     "Regression: When new code accidentally breaks something that used to work. Caught by automated tests. High regression rate = fragile codebase = tech debt signal.",
		"code review":    "Code Review: Process where other developers check code before it's merged. Catches bugs, shares knowledge, maintains standards. Average should be <24 hours.",
		"test coverage":  "Test Coverage: Percentage of code that has automated tests. Higher = more confidence when changing code. Industry standard: 60-80%. Below 40% = risky deployments.",
		"refactoring":    "Refactoring: Improving code structure without changing what it does. Makes future development faster. Not visible to users but critical for long-term velocity.",
		"hotfix":         "Hotfix: Emergency code fix deployed directly to production, skipping normal process. If frequent, signals quality problems in normal development.",
		"feature flag":   "Feature Flag: Toggle that enables/disables features without deployment. Allows: gradual rollouts, A/B testing, quick disable if broken. PM power tool.",
		"load balancer":  "Load Balancer: Distributes traffic across multiple servers. If one server fails, others handle it. PM relevance: system can handle traffic spikes without going down.",
		"latency":        "Latency: How long a request takes to complete (milliseconds). Users notice >200ms. PM should ask about latency impact when approving complex features.",
		"database migration": "Database Migration: Changing the structure of the database (adding/removing columns, tables). Risky if done wrong. PM should know: migrations can block releases.",
		"devops":         "DevOps: Culture/practice of collaboration between development and operations. Goal: faster, more reliable releases. PM relevance: understand deployment pipeline, not just code.",
		"sprint velocity": "Sprint Velocity: Average story points completed per sprint. Used ONLY for planning capacity, NOT for performance evaluation. Comparing velocity between teams is meaningless.",
	}

	if term != "" {
		lower := strings.ToLower(term)
		if def, ok := glossary[lower]; ok {
			return textResult(def), nil
		}
		// Try partial match
		for key, def := range glossary {
			if strings.Contains(key, lower) || strings.Contains(lower, key) {
				return textResult(def), nil
			}
		}
		// AI explain
		systemPrompt := "Explain this technical concept to a non-technical PM/Scrum Master in 2-3 sentences. Include: what it is, why it matters for the PM, and what question the PM should ask developers about it."
		result, err := h.AI.Complete(ctx, systemPrompt, "Explain: "+term)
		if err != nil {
			return textResult(fmt.Sprintf("Term '%s' not in glossary. Ask your developers to explain it — showing curiosity builds trust.", term)), nil
		}
		return textResult(result), nil
	}

	// List all terms
	var sb strings.Builder
	sb.WriteString("Technical Glossary for PMs:\n\n")
	sb.WriteString("Available terms: ")
	var terms []string
	for k := range glossary {
		terms = append(terms, k)
	}
	sb.WriteString(strings.Join(terms, ", "))
	sb.WriteString("\n\nUse term parameter to get definition. Or ask any tech term — AI will explain.")
	return textResult(sb.String()), nil
}

// QAHealthCheck analyzes sprint quality signals from Jira data.
func (h *Handlers) QAHealthCheck(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	project := req.GetString("project", "")

	// Count bugs vs total
	bugJQL := "issuetype = Bug AND resolution = Unresolved"
	allJQL := "resolution = Unresolved"
	if project != "" {
		bugJQL = fmt.Sprintf("project = %s AND %s", project, bugJQL)
		allJQL = fmt.Sprintf("project = %s AND %s", project, allJQL)
	}

	bugs, _ := h.Jira.SearchIssues(ctx, bugJQL, 200, 0)
	all, _ := h.Jira.SearchIssues(ctx, allJQL, 200, 0)

	bugCount := len(bugs.Issues)
	totalCount := len(all.Issues)

	var sb strings.Builder
	sb.WriteString("QA Health Check:\n\n")

	// Bug ratio
	bugRatio := 0.0
	if totalCount > 0 {
		bugRatio = float64(bugCount) / float64(totalCount) * 100
	}
	sb.WriteString(fmt.Sprintf("Open bugs: %d / %d total issues (%.0f%%)\n", bugCount, totalCount, bugRatio))

	status := "HEALTHY"
	if bugRatio > 30 {
		status = "CRITICAL — bug debt overwhelming"
	} else if bugRatio > 20 {
		status = "CONCERNING — prioritize bug reduction"
	}
	sb.WriteString(fmt.Sprintf("Status: %s\n\n", status))

	// Bug priority breakdown
	var critical, high, medium, low int
	for _, bug := range bugs.Issues {
		switch strings.ToLower(bug.Priority) {
		case "highest", "critical":
			critical++
		case "high":
			high++
		case "medium":
			medium++
		default:
			low++
		}
	}
	sb.WriteString("Bug Priority Breakdown:\n")
	sb.WriteString(fmt.Sprintf("  Critical: %d | High: %d | Medium: %d | Low: %d\n\n", critical, high, medium, low))

	// Stale bugs (>30 days old)
	staleBugs := 0
	for _, bug := range bugs.Issues {
		if time.Since(bug.Created).Hours() > 30*24 {
			staleBugs++
		}
	}
	if staleBugs > 0 {
		sb.WriteString(fmt.Sprintf("Stale bugs (>30 days): %d — these are rotting quality\n\n", staleBugs))
	}

	// Recommendations
	sb.WriteString("PM Actions:\n")
	if critical > 0 {
		sb.WriteString(fmt.Sprintf("  1. URGENT: %d critical bugs need immediate attention\n", critical))
	}
	if bugRatio > 20 {
		sb.WriteString("  2. Allocate 25-30% next sprint to bug fixes\n")
	}
	if staleBugs > 5 {
		sb.WriteString("  3. Schedule bug triage session — close or prioritize stale bugs\n")
	}
	sb.WriteString("  4. Ask devs: 'What's causing these bugs?' (root cause > symptom fix)\n")

	return textResult(sb.String()), nil
}

// DevWorkflowExplainer helps PM understand what developers do in their daily workflow.
func (h *Handlers) DevWorkflowExplainer(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	phase := req.GetString("phase", "overview")

	workflows := map[string]string{
		"overview": `Developer Daily Workflow:

1. PULL latest code (get team's changes)
2. PICK a ticket from sprint board
3. CREATE branch (isolated workspace)
4. CODE the solution (hours/days)
5. WRITE tests (verify it works)
6. COMMIT changes (save progress)
7. PUSH to remote (share with team)
8. CREATE PR/MR (request review)
9. ADDRESS feedback (fix review comments)
10. MERGE (code goes to main branch)
11. DEPLOY (code goes live)
12. VERIFY in production

Where PM can help:
- Clear requirements (step 4) → less rework
- Quick PR reviews (step 9) → faster delivery
- Remove blockers → unblock any step`,

		"planning": `What Developers Need from Sprint Planning:

MUST HAVE:
- Clear acceptance criteria (how do we know it's done?)
- Known dependencies (what's blocking this?)
- Access/credentials needed
- Design/UX specs if UI work
- API contracts if integration work

NICE TO HAVE:
- Priority within sprint (what first?)
- Who to ask questions
- Related past work to reference

RED FLAGS (dev frustration):
- "Figure it out" without context
- Changing requirements mid-sprint
- No acceptance criteria = no definition of done
- Underestimated complexity = overtime pressure`,

		"testing": `Testing Layers (what QA/devs do):

UNIT TESTS (developer writes):
  - Tests individual functions
  - Fast (seconds), runs on every commit
  - Coverage target: 60-80%

INTEGRATION TESTS:
  - Tests components working together
  - Tests API contracts, database queries
  - Slower but catches interaction bugs

E2E TESTS (QA writes):
  - Tests full user journeys
  - Slowest but most realistic
  - Catches what users will experience

PM Questions to Ask:
- "What's our test coverage for this feature?"
- "Are there automated regression tests?"
- "How will we know if this breaks something else?"`,

		"deployment": `Deployment Process Explained:

1. STAGING DEPLOY
   - Code goes to test environment
   - QA verifies in staging
   - PM can preview here before production

2. PRODUCTION DEPLOY
   - Code goes live to real users
   - Usually automated (CI/CD)
   - Monitoring kicks in immediately

3. MONITORING
   - Error rates, latency, user behavior
   - Alerts if something breaks
   - First 30 minutes are critical

4. ROLLBACK (if broken)
   - Revert to previous version
   - Should be <5 minutes
   - No data loss

PM Power Moves:
- Ask: "Can we rollback if needed?"
- Ask: "What does monitoring look like?"
- NEVER push for Friday deploys (weekend = no support)`,

		"code_review": `Code Review: What It Is & Why It Matters to PM

WHAT HAPPENS:
- Developer A writes code
- Developer B reads and comments on it
- They discuss approach, catch bugs, share knowledge
- Average: 30-60 minutes per review

WHY PM SHOULD CARE:
- Long review queues = delivery bottleneck
- No reviews = quality risk (bugs slip through)
- Single reviewer = knowledge silo

HEALTHY METRICS:
- Review turnaround: <24 hours
- Review participation: multiple reviewers per PR
- Comments: constructive, not blocking

PM ACTION:
- If PRs sit >48 hours: raise in standup
- Working agreement: "Reviews within 24h"
- Don't pressure devs to skip reviews for speed`,
	}

	if content, ok := workflows[phase]; ok {
		return textResult(content), nil
	}

	// List available phases
	var phases []string
	for k := range workflows {
		phases = append(phases, k)
	}
	return textResult(fmt.Sprintf("Available phases: %s\n\nUse phase parameter to learn about each.", strings.Join(phases, ", "))), nil
}

// EngineeringMetricsExplainer helps PM understand which metrics to look at and which to avoid.
func (h *Handlers) EngineeringMetricsExplainer(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return textResult(`Engineering Metrics — PM Guide

TRACK THESE (Team Health):
  Cycle Time         How long from start to done (days)
  Throughput         Items completed per sprint
  Sprint Goal Rate   Did we hit the goal? (yes/no per sprint)
  Escaped Defects    Bugs found in production
  PR Review Time     Hours from PR open to merged
  Deploy Frequency   How often we ship safely

AVOID THESE (Toxic Metrics):
  Lines of Code      More lines ≠ better. Less is often more.
  Individual Velocity   Destroys collaboration. NEVER compare people.
  Hours Worked       Overtime = burnout, not commitment.
  Commits per Day    Activity ≠ value. Quality > quantity.
  Bug Count per Dev  Creates blame culture, hides real issues.

UNDERSTAND THESE (Context Required):
  Velocity           Only for OUR planning. Not a target. Not comparable.
  Story Points       Relative only. They drift. They lie. Use for sizing only.
  Test Coverage      60-80% = good. 100% = probably wasting time.
  Bug Count          Matters in TREND, not absolute number.

THE ONE QUESTION THAT MATTERS:
  "Are we delivering value to users, sustainably, with improving quality?"
  If yes → metrics are secondary.
  If no → dig into cycle time + escaped defects + team satisfaction.

WHAT TO ASK IN STANDUP:
  NOT: "How many points did you do yesterday?"
  BUT: "What's blocking you?" / "Who needs help?"

WHAT TO ASK IN REVIEW:
  NOT: "Why didn't we finish everything?"
  BUT: "Did we meet the sprint goal? What did we learn?"
`), nil
}

// SprintQualityGate evaluates if sprint is release-ready based on quality signals.
func (h *Handlers) SprintQualityGate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) == 0 {
		return textResult("No active sprint."), nil
	}

	issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)

	var done, bugs, bugsInSprint, openBugs int
	for _, i := range issues {
		l := strings.ToLower(i.Status)
		if strings.Contains(l, "done") || strings.Contains(l, "closed") {
			done++
		}
		if strings.ToLower(i.Type) == "bug" {
			bugsInSprint++
			if !strings.Contains(l, "done") && !strings.Contains(l, "closed") {
				openBugs++
			}
		}
	}

	// Check for high-priority open bugs in project
	bugResult, _ := h.Jira.SearchIssues(ctx, "issuetype = Bug AND priority in (Highest, High) AND resolution = Unresolved", 50, 0)
	criticalBugs := len(bugResult.Issues)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sprint Quality Gate: %s\n\n", sprints[0].Name))

	gate1 := openBugs == 0
	gate2 := criticalBugs == 0
	gate3 := float64(done)/float64(len(issues)) >= 0.8
	gate4 := bugsInSprint <= 2

	sb.WriteString(fmt.Sprintf("[%s] No open bugs in sprint: %d open\n", gateIcon(gate1), openBugs))
	sb.WriteString(fmt.Sprintf("[%s] No critical/high bugs in project: %d found\n", gateIcon(gate2), criticalBugs))
	sb.WriteString(fmt.Sprintf("[%s] Sprint completion >= 80%%: %.0f%%\n", gateIcon(gate3), float64(done)/float64(len(issues))*100))
	sb.WriteString(fmt.Sprintf("[%s] Bugs in sprint <= 2: %d bugs\n", gateIcon(gate4), bugsInSprint))

	passed := 0
	if gate1 {
		passed++
	}
	if gate2 {
		passed++
	}
	if gate3 {
		passed++
	}
	if gate4 {
		passed++
	}

	sb.WriteString(fmt.Sprintf("\nGates passed: %d/4\n", passed))
	if passed == 4 {
		sb.WriteString("VERDICT: RELEASE READY\n")
	} else if passed >= 3 {
		sb.WriteString("VERDICT: CONDITIONAL — review failing gates before release\n")
	} else {
		sb.WriteString("VERDICT: NOT READY — address quality issues first\n")
	}

	_ = bugs // suppress unused
	return textResult(sb.String()), nil
}

func gateIcon(pass bool) string {
	if pass {
		return "PASS"
	}
	return "FAIL"
}
