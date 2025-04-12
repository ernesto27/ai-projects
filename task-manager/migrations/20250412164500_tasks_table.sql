-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) CHECK (type IN ('Story', 'Bug', 'Task', 'Epic')),
    status VARCHAR(50) DEFAULT 'To Do',
    priority VARCHAR(20) CHECK (priority IN ('Low', 'Medium', 'High', 'Critical')),
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    reporter_id INTEGER REFERENCES users(id),
    assignee_id INTEGER REFERENCES users(id),
    sprint_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMP,
    time_estimate INTEGER, -- in minutes
    time_spent INTEGER, -- in minutes
    CONSTRAINT fk_project FOREIGN KEY(project_id) REFERENCES projects(id),
    CONSTRAINT fk_reporter FOREIGN KEY(reporter_id) REFERENCES users(id),
    CONSTRAINT fk_assignee FOREIGN KEY(assignee_id) REFERENCES users(id)
);

CREATE INDEX idx_tasks_project ON tasks(project_id);
CREATE INDEX idx_tasks_assignee ON tasks(assignee_id);
CREATE INDEX idx_tasks_status ON tasks(status);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE tasks;