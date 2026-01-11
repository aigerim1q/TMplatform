# Chatbot DB Test Cases

## Create Project (Happy Path)
- **Input:** "Create a project for a 5-story residential building"
- **Expected:** Project created with title captured; response includes project id and status.

## Duplicate Project Creation (Idempotent)
- **Input 1:** "Create project Skyline" (creates)
- **Input 2:** "Create project Skyline" (retry)
- **Expected:** Second call returns existing project without creating a new record.

## Incomplete Input → Clarification
- **Input:** "Create project" (no title)
- **Expected:** Bot asks for project title; no DB write occurs.

## Assign Responsible
- **Success:** "Assign Alex to task permits" (task exists and user has access) → task assignee updated.
- **No Permission:** Other user tries to assign in a project they do not belong to → forbidden message, no change.
- **User Not Found:** Unknown assignee name → bot asks to confirm/clarify.

## Add Stages & Tasks
- **Input:** "Add stages: Design; Construction | tasks: permits, drawings to project Skyline"
- **Expected:** Stages and tasks created idempotently under Skyline. Re-running does not duplicate rows.

## List Projects
- **Empty:** New user without membership → "no projects yet" message.
- **Many:** User with several projects → paginated list with ids and statuses.
- **Filtered:** "list projects skyline" → returns matching projects only.

## Show Project Details
- **Input:** "show project Skyline"
- **Expected:** Returns project summary with stages, tasks, statuses, and assignees.

## Database Error Handling
- Simulate DB unavailable (bad DSN) → startup fails with clear log; bot should not start.
- Transaction failure during stage/task creation → entire transaction rolls back, no partial records.
