create table todos (
  "todoId"      uuid        not null default gen_random_uuid(),
  "task"        text        not null,
  "isCompleted" boolean     not null default false,
  "createdAt"   timestamptz not null default now(),
  "updatedAt"   timestamptz not null default now()
);
