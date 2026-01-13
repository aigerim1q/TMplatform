package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"

	"silencendo/db"
	"silencendo/models"
	"silencendo/repository"
)

// CreateProjectInput holds chat-derived data for project creation.
type CreateProjectInput struct {
	Title       string
	Description string
	Status      string
}

// TaskPlan represents a task to create.
type TaskPlan struct {
	Title        string
	Description  string
	Status       string
	Priority     string
	AssigneeName *string
}

// StagePlan represents a stage with its tasks.
type StagePlan struct {
	Title           string
	Order           int
	ResponsibleName *string
	Tasks           []TaskPlan
}

// MemberPlan captures a user to add to a project.
type MemberPlan struct {
	Name  string
	Email *string
}

// ListFilters controls project listing.
type ListFilters struct {
	Search string
	Limit  int
	Offset int
}

// ProjectService orchestrates project-related operations with validation, idempotency, and ACL.
type ProjectService struct {
	db          *sql.DB
	users       *repository.UserRepository
	projects    *repository.ProjectRepository
	stages      *repository.StageRepository
	tasks       *repository.TaskRepository
	memberships *repository.MembershipRepository
	idempotency *repository.IdempotencyRepository
}

func NewProjectService(dbConn *sql.DB) *ProjectService {
	return &ProjectService{
		db:          dbConn,
		users:       repository.NewUserRepository(dbConn),
		projects:    repository.NewProjectRepository(dbConn),
		stages:      repository.NewStageRepository(dbConn),
		tasks:       repository.NewTaskRepository(dbConn),
		memberships: repository.NewMembershipRepository(dbConn),
		idempotency: repository.NewIdempotencyRepository(dbConn),
	}
}

func (s *ProjectService) CreateProjectFromChat(ctx context.Context, input CreateProjectInput, user models.User) (models.Project, bool, error) {
	title := strings.TrimSpace(input.Title)
	description := strings.TrimSpace(input.Description)
	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = "active"
	}

	if title == "" {
		return models.Project{}, false, ErrValidation
	}
	if len(title) > 200 {
		return models.Project{}, false, ErrValidation
	}

	normalized := normalize(title)
	key := hashKey(user.ID, "create_project", normalized, description, status)

	var created bool
	var result models.Project

	err := db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		actor, err := s.ensureUser(ctx, tx, user)
		if err != nil {
			return err
		}

		if resourceType, resourceID, found, err := s.idempotency.Lookup(ctx, tx, key, actor.ID); err != nil {
			return err
		} else if found && resourceType == "project" {
			result, err = s.projects.GetByID(ctx, tx, resourceID)
			return err
		}

		existing, err := s.projects.GetByOwnerAndNormalizedTitle(ctx, tx, actor.ID, normalized)
		if err == nil {
			result = existing
			return s.idempotency.Record(ctx, tx, key, actor.ID, "create_project", "project", existing.ID)
		}
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return err
		}

		project := models.Project{
			ID:              uuid.NewString(),
			OwnerID:         actor.ID,
			Title:           title,
			Description:     description,
			Status:          status,
			NormalizedTitle: normalized,
		}
		created = true

		saved, err := s.projects.Create(ctx, tx, project)
		if err != nil {
			return err
		}

		if err := s.memberships.AddIfMissing(ctx, tx, saved.ID, actor.ID, "owner"); err != nil {
			return err
		}

		if err := s.idempotency.Record(ctx, tx, key, actor.ID, "create_project", "project", saved.ID); err != nil {
			return err
		}

		result = saved
		return nil
	})

	return result, created, err
}

func (s *ProjectService) CreateStagesAndTasksFromChat(ctx context.Context, projectID string, plans []StagePlan, user models.User) ([]models.Stage, []models.Task, error) {
	if projectID == "" {
		return nil, nil, ErrValidation
	}
	if len(plans) == 0 {
		return nil, nil, ErrValidation
	}

	// Flatten plan for idempotency key.
	var flatParts []string
	for _, p := range plans {
		flatParts = append(flatParts, normalize(p.Title))
		if p.ResponsibleName != nil {
			flatParts = append(flatParts, normalize(*p.ResponsibleName))
		}
		for _, t := range p.Tasks {
			flatParts = append(flatParts, normalize(t.Title), normalize(t.Priority), normalize(t.Status))
			if t.AssigneeName != nil {
				flatParts = append(flatParts, normalize(*t.AssigneeName))
			}
		}
	}
	key := hashKey(user.ID, "add_stages_tasks", projectID, strings.Join(flatParts, "|"))

	var stages []models.Stage
	var tasks []models.Task

	err := db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		actor, err := s.ensureUser(ctx, tx, user)
		if err != nil {
			return err
		}

		// Verify membership.
		isMember, err := s.memberships.IsMember(ctx, tx, projectID, actor.ID)
		if err != nil {
			return err
		}
		if !isMember {
			return ErrForbidden
		}

		// Idempotency check.
		if resourceType, resourceID, found, err := s.idempotency.Lookup(ctx, tx, key, actor.ID); err != nil {
			return err
		} else if found && resourceType == "project" && resourceID == projectID {
			// Return current state without creating duplicates.
			stages, err = s.stages.ListByProject(ctx, tx, projectID)
			if err != nil {
				return err
			}
			tasks, err = s.tasks.ListByProject(ctx, tx, projectID)
			return err
		}

		// Ensure project exists and actor can see it.
		if _, err := s.projects.GetByID(ctx, tx, projectID); err != nil {
			return err
		}

		for idx, plan := range plans {
			stageTitle := strings.TrimSpace(plan.Title)
			if stageTitle == "" {
				return ErrValidation
			}
			normalizedStage := normalize(stageTitle)

			var stage models.Stage
			stage, err = s.stages.GetByProjectAndNormalizedTitle(ctx, tx, projectID, normalizedStage)
			if err != nil {
				if !errors.Is(err, repository.ErrNotFound) {
					return err
				}
				stage = models.Stage{
					ID:              uuid.NewString(),
					ProjectID:       projectID,
					Title:           stageTitle,
					NormalizedTitle: normalizedStage,
					OrderIndex:      plan.Order,
				}
				if plan.Order == 0 {
					stage.OrderIndex = idx + 1
				}
				stage, err = s.stages.Create(ctx, tx, stage)
				if err != nil {
					return err
				}
			}

			// Assign responsible if provided.
			if plan.ResponsibleName != nil {
				responsibleName := strings.TrimSpace(*plan.ResponsibleName)
				if responsibleName != "" {
					respUser := models.User{ID: uuid.NewString(), Name: responsibleName}
					respUser, err = s.ensureUser(ctx, tx, respUser)
					if err != nil {
						return err
					}
					if err := s.memberships.AddIfMissing(ctx, tx, projectID, respUser.ID, "member"); err != nil {
						return err
					}
					if err := s.stages.UpdateResponsible(ctx, tx, stage.ID, &respUser.ID); err != nil {
						return err
					}
					stage.ResponsibleID = &respUser.ID
				}
			}

			stages = append(stages, stage)

			for _, taskPlan := range plan.Tasks {
				taskTitle := strings.TrimSpace(taskPlan.Title)
				if taskTitle == "" {
					return ErrValidation
				}
				normalizedTask := normalize(taskTitle)

				var task models.Task
				task, err = s.tasks.GetByProjectAndNormalizedTitle(ctx, tx, projectID, normalizedTask)
				if err != nil && !errors.Is(err, repository.ErrNotFound) {
					return err
				}

				// Only create if missing to preserve idempotency at task level.
				if errors.Is(err, repository.ErrNotFound) {
					task = models.Task{
						ID:              uuid.NewString(),
						ProjectID:       projectID,
						StageID:         &stage.ID,
						Title:           taskTitle,
						NormalizedTitle: normalizedTask,
						Description:     strings.TrimSpace(taskPlan.Description),
						Status:          defaultIfEmpty(taskPlan.Status, "open"),
						Priority:        defaultIfEmpty(taskPlan.Priority, "medium"),
					}

					if task.Status == "" {
						task.Status = "open"
					}
					if task.Priority == "" {
						task.Priority = "medium"
					}

					if taskPlan.AssigneeName != nil {
						assigneeName := strings.TrimSpace(*taskPlan.AssigneeName)
						if assigneeName != "" {
							assignee := models.User{ID: uuid.NewString(), Name: assigneeName}
							assignee, err = s.ensureUser(ctx, tx, assignee)
							if err != nil {
								return err
							}
							if err := s.memberships.AddIfMissing(ctx, tx, projectID, assignee.ID, "member"); err != nil {
								return err
							}
							task.AssigneeID = &assignee.ID
						}
					}

					task, err = s.tasks.Create(ctx, tx, task)
					if err != nil {
						return err
					}
				}

				tasks = append(tasks, task)
			}
		}

		return s.idempotency.Record(ctx, tx, key, actor.ID, "add_stages_tasks", "project", projectID)
	})

	return stages, tasks, err
}

// Internal helper to create stages/tasks assuming membership already checked.
func (s *ProjectService) createStagesTasksInternal(ctx context.Context, tx *sql.Tx, projectID string, plans []StagePlan) ([]models.Stage, []models.Task, error) {
	var stages []models.Stage
	var tasks []models.Task

	for idx, plan := range plans {
		stageTitle := strings.TrimSpace(plan.Title)
		if stageTitle == "" {
			return nil, nil, ErrValidation
		}
		normalizedStage := normalize(stageTitle)

		var stage models.Stage
		stage, err := s.stages.GetByProjectAndNormalizedTitle(ctx, tx, projectID, normalizedStage)
		if err != nil {
			if !errors.Is(err, repository.ErrNotFound) {
				return nil, nil, err
			}
			stage = models.Stage{
				ID:              uuid.NewString(),
				ProjectID:       projectID,
				Title:           stageTitle,
				NormalizedTitle: normalizedStage,
				OrderIndex:      plan.Order,
			}
			if plan.Order == 0 {
				stage.OrderIndex = idx + 1
			}
			stage, err = s.stages.Create(ctx, tx, stage)
			if err != nil {
				return nil, nil, err
			}
		}

		// Assign responsible if provided.
		if plan.ResponsibleName != nil {
			responsibleName := strings.TrimSpace(*plan.ResponsibleName)
			if responsibleName != "" {
				respUser := models.User{ID: uuid.NewString(), Name: responsibleName}
				respUser, err = s.ensureUser(ctx, tx, respUser)
				if err != nil {
					return nil, nil, err
				}
				if err := s.memberships.AddIfMissing(ctx, tx, projectID, respUser.ID, "member"); err != nil {
					return nil, nil, err
				}
				if err := s.stages.UpdateResponsible(ctx, tx, stage.ID, &respUser.ID); err != nil {
					return nil, nil, err
				}
				stage.ResponsibleID = &respUser.ID
			}
		}

		stages = append(stages, stage)

		for _, taskPlan := range plan.Tasks {
			taskTitle := strings.TrimSpace(taskPlan.Title)
			if taskTitle == "" {
				return nil, nil, ErrValidation
			}
			normalizedTask := normalize(taskTitle)

			var task models.Task
			task, err = s.tasks.GetByProjectAndNormalizedTitle(ctx, tx, projectID, normalizedTask)
			if err != nil && !errors.Is(err, repository.ErrNotFound) {
				return nil, nil, err
			}

			// Only create if missing to preserve idempotency at task level.
			if errors.Is(err, repository.ErrNotFound) {
				task = models.Task{
					ID:              uuid.NewString(),
					ProjectID:       projectID,
					StageID:         &stage.ID,
					Title:           taskTitle,
					NormalizedTitle: normalizedTask,
					Description:     strings.TrimSpace(taskPlan.Description),
					Status:          defaultIfEmpty(taskPlan.Status, "open"),
					Priority:        defaultIfEmpty(taskPlan.Priority, "medium"),
				}

				if task.Status == "" {
					task.Status = "open"
				}
				if task.Priority == "" {
					task.Priority = "medium"
				}

				if taskPlan.AssigneeName != nil {
					assigneeName := strings.TrimSpace(*taskPlan.AssigneeName)
					if assigneeName != "" {
						assignee := models.User{ID: uuid.NewString(), Name: assigneeName}
						assignee, err = s.ensureUser(ctx, tx, assignee)
						if err != nil {
							return nil, nil, err
						}
						if err := s.memberships.AddIfMissing(ctx, tx, projectID, assignee.ID, "member"); err != nil {
							return nil, nil, err
						}
						task.AssigneeID = &assignee.ID
					}
				}

				task, err = s.tasks.Create(ctx, tx, task)
				if err != nil {
					return nil, nil, err
				}
			}

			tasks = append(tasks, task)
		}
	}

	return stages, tasks, nil
}

// AddMembersAndOptionalPlan adds members to a project and, optionally, creates a simple plan (one stage with tasks).
// It is idempotent on memberships and on the plan content when provided.
func (s *ProjectService) AddMembersAndOptionalPlan(ctx context.Context, projectID string, members []MemberPlan, plans []StagePlan, user models.User) ([]models.User, []models.Stage, []models.Task, error) {
	if projectID == "" {
		return nil, nil, nil, ErrValidation
	}

	var addedUsers []models.User
	var stages []models.Stage
	var tasks []models.Task

	err := db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		actor, err := s.ensureUser(ctx, tx, user)
		if err != nil {
			return err
		}

		// Ensure actor is member
		if err := s.membershipsCheck(ctx, tx, projectID, actor.ID); err != nil {
			return err
		}

		// Add members idempotently
		for _, m := range members {
			name := strings.TrimSpace(m.Name)
			if name == "" {
				continue
			}
			member := models.User{ID: uuid.NewString(), Name: name, Email: m.Email}
			member, err = s.ensureUser(ctx, tx, member)
			if err != nil {
				return err
			}
			if err := s.memberships.AddIfMissing(ctx, tx, projectID, member.ID, "member"); err != nil {
				return err
			}
			addedUsers = append(addedUsers, member)
		}

		// Optional plan creation
		if len(plans) > 0 {
			if _, err := s.projects.GetByID(ctx, tx, projectID); err != nil {
				return err
			}
			// Reuse existing idempotent flow from CreateStagesAndTasksFromChat with actor as user
			stages, tasks, err = s.createStagesTasksInternal(ctx, tx, projectID, plans)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return addedUsers, stages, tasks, err
}

func (s *ProjectService) AssignResponsible(ctx context.Context, entityType, entityID, assigneeID string, actor models.User) error {
	if entityType == "" || entityID == "" || assigneeID == "" {
		return ErrValidation
	}

	return db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		actor, err := s.ensureUser(ctx, tx, actor)
		if err != nil {
			return err
		}

		// Ensure assignee exists.
		assignee, err := s.users.GetByID(ctx, tx, assigneeID)
		if err != nil {
			return err
		}

		// Resolve entity and project.
		var projectID string
		switch strings.ToLower(entityType) {
		case "task":
			task, err := s.tasks.GetByID(ctx, tx, entityID)
			if err != nil {
				return err
			}
			projectID = task.ProjectID
			if err := s.membershipsCheck(ctx, tx, projectID, actor.ID); err != nil {
				return err
			}
			if err := s.memberships.AddIfMissing(ctx, tx, projectID, assignee.ID, "member"); err != nil {
				return err
			}
			return s.tasks.UpdateAssignee(ctx, tx, entityID, &assignee.ID)
		case "stage":
			stage, err := s.stages.GetByID(ctx, tx, entityID)
			if err != nil {
				return err
			}
			projectID = stage.ProjectID
			if err := s.membershipsCheck(ctx, tx, projectID, actor.ID); err != nil {
				return err
			}
			if err := s.memberships.AddIfMissing(ctx, tx, projectID, assignee.ID, "member"); err != nil {
				return err
			}
			return s.stages.UpdateResponsible(ctx, tx, entityID, &assignee.ID)
		default:
			return ErrValidation
		}
	})
}

func (s *ProjectService) ListProjects(ctx context.Context, user models.User, filters ListFilters) ([]models.Project, error) {
	limit := filters.Limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	offset := filters.Offset
	if offset < 0 {
		offset = 0
	}

	var projects []models.Project

	err := db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		actor, err := s.ensureUser(ctx, tx, user)
		if err != nil {
			return err
		}
		projects, err = s.projects.ListByUser(ctx, tx, actor.ID, filters.Search, limit, offset)
		return err
	})

	return projects, err
}

// SetActiveProject sets the project's business status to active. Session selection is managed by the caller.
func (s *ProjectService) SetActiveProject(ctx context.Context, projectID string, user models.User) (models.Project, error) {
	return s.SetProjectStatus(ctx, projectID, "active", user)
}

// SetProjectStatus updates the business status (active/inactive) for a project without affecting session routing.
func (s *ProjectService) SetProjectStatus(ctx context.Context, projectID, status string, user models.User) (models.Project, error) {
	var updated models.Project

	status = strings.ToLower(strings.TrimSpace(status))
	if projectID == "" || (status != "active" && status != "inactive") {
		return updated, ErrValidation
	}

	err := db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		actor, err := s.ensureUser(ctx, tx, user)
		if err != nil {
			return err
		}

		if err := s.membershipsCheck(ctx, tx, projectID, actor.ID); err != nil {
			return err
		}

		target, err := s.projects.GetByID(ctx, tx, projectID)
		if err != nil {
			return err
		}

		if err := s.projects.UpdateStatus(ctx, tx, projectID, status); err != nil {
			return err
		}
		target.Status = status
		updated = target
		return nil
	})

	return updated, err
}

// SwitchActiveProject sets the target project to active and, if provided, marks the previous project inactive in a single transaction.
// Routing remains a caller concern; this only updates business status.
func (s *ProjectService) SwitchActiveProject(ctx context.Context, newProjectID, previousProjectID string, user models.User) (activated models.Project, deactivated *models.Project, err error) {
	newProjectID = strings.TrimSpace(newProjectID)
	previousProjectID = strings.TrimSpace(previousProjectID)
	if newProjectID == "" {
		return activated, nil, ErrValidation
	}

	err = db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		actor, err := s.ensureUser(ctx, tx, user)
		if err != nil {
			return err
		}

		if err := s.membershipsCheck(ctx, tx, newProjectID, actor.ID); err != nil {
			return err
		}

		activated, err = s.projects.GetByID(ctx, tx, newProjectID)
		if err != nil {
			return err
		}
		if err := s.projects.UpdateStatus(ctx, tx, newProjectID, "active"); err != nil {
			return err
		}
		activated.Status = "active"

		if previousProjectID != "" && previousProjectID != newProjectID {
			if err := s.membershipsCheck(ctx, tx, previousProjectID, actor.ID); err != nil {
				return err
			}
			prev, err := s.projects.GetByID(ctx, tx, previousProjectID)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return nil
				}
				return err
			}
			if err := s.projects.UpdateStatus(ctx, tx, previousProjectID, "inactive"); err != nil {
				return err
			}
			prev.Status = "inactive"
			deactivated = &prev
		}

		return nil
	})

	return activated, deactivated, err
}

// SoftDeleteProject performs a soft delete and returns the deleted project.
func (s *ProjectService) SoftDeleteProject(ctx context.Context, projectID string, user models.User) (models.Project, error) {
	var deleted models.Project
	if strings.TrimSpace(projectID) == "" {
		return deleted, ErrValidation
	}

	err := db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		actor, err := s.ensureUser(ctx, tx, user)
		if err != nil {
			return err
		}

		if err := s.membershipsCheck(ctx, tx, projectID, actor.ID); err != nil {
			return err
		}

		deleted, err = s.projects.GetByID(ctx, tx, projectID)
		if err != nil {
			return err
		}

		if err := s.projects.SoftDelete(ctx, tx, projectID); err != nil {
			return err
		}
		return nil
	})

	return deleted, err
}

// CreateRuleBasedPlan adds members (if provided) and ensures a simple default plan.
func (s *ProjectService) CreateRuleBasedPlan(ctx context.Context, projectID string, memberNames []string, user models.User) ([]models.User, []models.Stage, []models.Task, error) {
	memberNames = uniqueNonEmpty(memberNames)
	var members []MemberPlan
	for _, name := range memberNames {
		members = append(members, MemberPlan{Name: name})
	}

	projectTitle, err := s.fetchProjectTitleForUser(ctx, projectID, user)
	if err != nil {
		return nil, nil, nil, err
	}

	plan := buildPlanForProject(projectTitle, memberNames)
	return s.AddMembersAndOptionalPlan(ctx, projectID, members, plan, user)
}

func (s *ProjectService) fetchProjectTitleForUser(ctx context.Context, projectID string, user models.User) (string, error) {
	var title string

	err := db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		actor, err := s.ensureUser(ctx, tx, user)
		if err != nil {
			return err
		}

		if err := s.membershipsCheck(ctx, tx, projectID, actor.ID); err != nil {
			return err
		}

		project, err := s.projects.GetByID(ctx, tx, projectID)
		if err != nil {
			return err
		}

		title = project.Title
		return nil
	})

	return title, err
}

func buildPlanForProject(projectTitle string, memberNames []string) []StagePlan {
	picker := makeNamePicker(memberNames)
	if isMinecraftProject(projectTitle) {
		return buildMinecraftPlan(picker)
	}

	return []StagePlan{
		{Title: "Planning", Order: 1, ResponsibleName: picker(0), Tasks: []TaskPlan{{Title: "Clarify scope", AssigneeName: picker(0)}, {Title: "Define milestones", AssigneeName: picker(1)}}},
		{Title: "Execution", Order: 2, ResponsibleName: picker(1), Tasks: []TaskPlan{{Title: "Assign responsibilities", AssigneeName: picker(0)}, {Title: "Start core work", AssigneeName: picker(1)}}},
		{Title: "Review", Order: 3, ResponsibleName: picker(0), Tasks: []TaskPlan{{Title: "Review progress", AssigneeName: picker(0)}, {Title: "Plan next steps", AssigneeName: picker(1)}}},
	}
}

func buildMinecraftPlan(picker func(int) *string) []StagePlan {
	return []StagePlan{
		{Title: "Preparation", Order: 1, ResponsibleName: picker(0), Tasks: []TaskPlan{
			{Title: "Gather food, beds, and shields for the run", AssigneeName: picker(0)},
			{Title: "Craft iron gear (bucket, armor, tools, flint & steel)", AssigneeName: picker(1)},
			{Title: "Collect boats and blocks for fast overworld travel", AssigneeName: picker(0)},
		}},
		{Title: "Nether", Order: 2, ResponsibleName: picker(1), Tasks: []TaskPlan{
			{Title: "Find lava pool and enter the Nether quickly", AssigneeName: picker(1)},
			{Title: "Locate fortress and secure 12 blaze rods", AssigneeName: picker(0)},
			{Title: "Trade with piglins for 12-16 ender pearls", AssigneeName: picker(1)},
		}},
		{Title: "Endgame", Order: 3, ResponsibleName: picker(0), Tasks: []TaskPlan{
			{Title: "Craft eyes of ender and locate the stronghold", AssigneeName: picker(0)},
			{Title: "Set spawn, place beds, and prep water buckets in the End", AssigneeName: picker(1)},
			{Title: "Defeat the Ender Dragon with bed/axe strategy", AssigneeName: picker(0)},
		}},
	}
}

func isMinecraftProject(projectTitle string) bool {
	lower := strings.ToLower(strings.TrimSpace(projectTitle))
	return strings.Contains(lower, "minecraft") || strings.Contains(lower, "nether") || strings.Contains(lower, "ender") || strings.Contains(lower, "speedrun")
}

func makeNamePicker(names []string) func(int) *string {
	if len(names) == 0 {
		return func(int) *string { return nil }
	}
	return func(idx int) *string {
		if len(names) == 0 {
			return nil
		}
		name := strings.TrimSpace(names[idx%len(names)])
		if name == "" {
			return nil
		}
		copy := name
		return &copy
	}
}

func (s *ProjectService) GetProjectDetails(ctx context.Context, projectID string, user models.User) (models.ProjectDetails, error) {
	var details models.ProjectDetails
	if projectID == "" {
		return details, ErrValidation
	}

	err := db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		actor, err := s.ensureUser(ctx, tx, user)
		if err != nil {
			return err
		}

		if err := s.membershipsCheck(ctx, tx, projectID, actor.ID); err != nil {
			return err
		}

		details.Project, err = s.projects.GetByID(ctx, tx, projectID)
		if err != nil {
			return err
		}

		details.Stages, err = s.stages.ListByProject(ctx, tx, projectID)
		if err != nil {
			return err
		}

		details.Tasks, err = s.tasks.ListByProject(ctx, tx, projectID)
		return err
	})

	return details, err
}

// FindProjectByTitle returns a project the user can access that matches the given title (case-insensitive).
func (s *ProjectService) FindProjectByTitle(ctx context.Context, title string, user models.User) (models.Project, error) {
	normalized := normalize(title)
	if normalized == "" {
		return models.Project{}, ErrValidation
	}

	var project models.Project
	err := db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		actor, err := s.ensureUser(ctx, tx, user)
		if err != nil {
			return err
		}
		project, err = s.projects.FindByUserAndNormalizedTitle(ctx, tx, actor.ID, normalized)
		return err
	})
	return project, err
}

func (s *ProjectService) FindStageForUserByTitle(ctx context.Context, title string, user models.User) (models.Stage, error) {
	normalized := normalize(title)
	if normalized == "" {
		return models.Stage{}, ErrValidation
	}
	var stage models.Stage
	err := db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		actor, err := s.ensureUser(ctx, tx, user)
		if err != nil {
			return err
		}
		stage, err = s.stages.FindForUserByTitle(ctx, tx, actor.ID, normalized)
		return err
	})
	return stage, err
}

func (s *ProjectService) FindTaskForUserByTitle(ctx context.Context, title string, user models.User) (models.Task, error) {
	normalized := normalize(title)
	if normalized == "" {
		return models.Task{}, ErrValidation
	}
	var task models.Task
	err := db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		actor, err := s.ensureUser(ctx, tx, user)
		if err != nil {
			return err
		}
		task, err = s.tasks.FindForUserByTitle(ctx, tx, actor.ID, normalized)
		return err
	})
	return task, err
}

func (s *ProjectService) ensureUser(ctx context.Context, tx *sql.Tx, user models.User) (models.User, error) {
	user.ID = strings.TrimSpace(user.ID)
	user.Name = strings.TrimSpace(user.Name)
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	if user.Name == "" {
		user.Name = "Chat User"
	}

	if user.Email != nil {
		email := strings.TrimSpace(*user.Email)
		user.Email = &email
		if email == "" {
			user.Email = nil
		}
	}

	return s.users.Upsert(ctx, tx, user)
}

func (s *ProjectService) membershipsCheck(ctx context.Context, tx *sql.Tx, projectID, userID string) error {
	ok, err := s.memberships.IsMember(ctx, tx, projectID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrForbidden
	}
	return nil
}

// EnsureUser creates or updates a user record and returns the persisted version.
func (s *ProjectService) EnsureUser(ctx context.Context, user models.User) (models.User, error) {
	var stored models.User
	err := db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		var err error
		stored, err = s.ensureUser(ctx, tx, user)
		return err
	})
	return stored, err
}

func defaultIfEmpty(value, fallback string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return fallback
	}
	return v
}

func uniqueNonEmpty(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, v := range values {
		n := strings.TrimSpace(v)
		if n == "" {
			continue
		}
		key := strings.ToLower(n)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, n)
	}
	return out
}
func (s *ProjectService) GetUserNameByID(ctx context.Context, userID string) (string, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return "", nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	u, err := s.users.GetByID(ctx, tx, userID)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	return strings.TrimSpace(u.Name), nil
}
