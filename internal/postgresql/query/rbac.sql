
-- name: SelectProfile :one
SELECT
  id,
  profile_picture,
  profile_background,
  first_name,
  last_name,
  mobile,
  email,
  created_at
FROM
  profiles
WHERE
  id = @id
LIMIT 1;

-- name: InsertProfile :one
INSERT INTO profiles (
  profile_picture,
  profile_background,
  first_name,
  last_name,
  mobile,
  email
)
VALUES (
  @profile_picture,
  @profile_background,
  @first_name,
  @last_name,
  @mobile,
  @email
)
RETURNING id;

-- name: UpdateProfile :exec
UPDATE profiles SET
  profile_picture       = @profile_picture,
  profile_background    = @profile_background,
  first_name            = @first_name,
  last_name             = @last_name,
  mobile                = @mobile,
  email                 = @email
WHERE id = @id;

-- name: DeleteProfile :exec
DELETE FROM profiles
WHERE id = @id;


-- name: SelectAccounts :one
SELECT
  id,
  username,
  hashedpassword,
  profile,
  is_blocked,
  created_at
FROM
  accounts
WHERE
  username = @username
LIMIT 1;

-- name: SelectAccountsById :one
SELECT
  id,
  username,
  hashedpassword,
  profile,
  is_blocked,
  created_at
FROM
  accounts
WHERE
  id = @id
LIMIT 1;


-- name: InsertAccounts :one
INSERT INTO accounts (
  username,
  hashedpassword,
  profile
)
VALUES (
  @username,
  @hashedpassword,
  @profile
)
RETURNING id;


-- name: ChangePassword :exec
UPDATE accounts SET
  hashedpassword = @hashedpassword
WHERE username = @username;

-- name: DeleteAccount :exec
UPDATE accounts SET
  is_blocked = true
WHERE username = @username;


-- name: SelectAccountRole :one
SELECT
  id,
  account_id,
  role_id,
  created_at
FROM
  account_roles
WHERE
  id = @id
LIMIT 1;

-- name: InsertAccountRole :one
INSERT INTO account_roles (
    account_id,
    role_id
)
VALUES (
  @account_id,
  @role_id
)
RETURNING id;

-- name: UpdateAccountRole :exec
UPDATE account_roles SET
  account_id = @account_id,
  role_id    = @role_id
WHERE id = @id;

-- name: DeleteAccountRole :exec
DELETE FROM account_roles
WHERE id = @id;


-- name: SelectHelpText :one
SELECT
  id,
  task_id,
  helptext,
  created_at
FROM
  helptext
WHERE
  id = @id
LIMIT 1;

-- name: SelectHelpTextByTasks :one
SELECT
  id,
  task_id,
  helptext,
  created_at
FROM
  helptext
WHERE
  task_id = @task_id
LIMIT 1;

-- name: InsertHelpText :one
INSERT INTO helptext (
    task_id,
    helptext
)
VALUES (
  @task_id,
  @helptext
)
RETURNING id;

-- name: UpdateHelpText :exec
UPDATE helptext SET
  task_id   = @task_id,
  helptext  = @helptext
WHERE id = @id;

-- name: DeleteHelpText :exec
DELETE FROM helptext
WHERE id = @id;

-- name: DeleteHelpTextByTask :exec
DELETE FROM helptext
WHERE task_id = @task_id;


-- name: SelectMenu :one
SELECT
  id,
  task_id,
  name,
  created_at
FROM
  menu
WHERE
  id = @id
LIMIT 1;

-- name: SelectMenuByTask :many
SELECT
  id,
  task_id,
  name,
  created_at
FROM
  menu
WHERE
  task_id = @task_id;

-- name: InsertMenu :one
INSERT INTO menu (
    task_id,
    name
)
VALUES (
  @task_id,
  @name
)
RETURNING id;

-- name: UpdateMenu :exec
UPDATE menu SET
  task_id   = @task_id,
  name  = @name
WHERE id = @id;

-- name: DeleteMenu :exec
DELETE FROM menu
WHERE id = @id;

-- name: DeleteMenuByTask :exec
DELETE FROM menu
WHERE task_id = @task_id;


-- name: SelectNavigation :one
SELECT
  id,
  task_id,
  name,
  created_at
FROM
  navigation
WHERE
  id = @id
LIMIT 1;

-- name: SelectNavigationByTask :many
SELECT
  id,
  task_id,
  name,
  created_at
FROM
  navigation
WHERE
  task_id = @task_id;

-- name: InsertNavigation :one
INSERT INTO navigation (
    task_id,
    name
)
VALUES (
  @task_id,
  @name
)
RETURNING id;

-- name: UpdateNavigation :exec
UPDATE navigation SET
  task_id   = @task_id,
  name  = @name
WHERE id = @id;

-- name: DeleteNavigation :exec
DELETE FROM navigation
WHERE id = @id;

-- name: DeleteNavigationByTask :exec
DELETE FROM navigation
WHERE task_id = @task_id;


-- name: SelectRole :one
SELECT
  id,
  role,
  created_at
FROM
  roles
WHERE
  id = @id
LIMIT 1;

-- name: InsertRole :one
INSERT INTO roles (
    role
)
VALUES (
  @role
)
RETURNING id;

-- name: UpdateRole :exec
UPDATE roles SET
  role = @role
WHERE id = @id;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = @id;



-- name: SelectRoleTask :one
SELECT
  id,
  task_id,
  role_id,
  created_at
FROM
  role_tasks
WHERE
  id = @id
LIMIT 1;

-- name: InsertRoleTask :one
INSERT INTO role_tasks (
    task_id,
    role_id
)
VALUES (
  @task_id,
  @role_id
)
RETURNING id;

-- name: UpdateRoleTask :exec
UPDATE role_tasks SET
  task_id = @task_id,
  role_id    = @role_id
WHERE id = @id;

-- name: DeleteRoleTask :exec
DELETE FROM role_tasks
WHERE id = @id;

-- name: DeleteRoleTaskByTask :exec
DELETE FROM role_tasks
WHERE task_id = @task_id;


-- name: SelectTask :one
SELECT
  id,
  task,
  created_at
FROM
  tasks
WHERE
  id = @id
LIMIT 1;

-- name: InsertTask :one
INSERT INTO tasks (
    task
)
VALUES (
  @task
)
RETURNING id;

-- name: UpdateTask :exec
UPDATE tasks SET
  task = @task
WHERE id = @id;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = @id;