CREATE TABLE public.profiles (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    resume TEXT DEFAULT 'The job resume is empty',
    password TEXT NOT NULL
)