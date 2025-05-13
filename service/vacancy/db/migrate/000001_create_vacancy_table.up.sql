CREATE TABLE public.vacancies
(
    id SERIAL PRIMARY KEY,
    owner_email VARCHAR(255) NOT NULL UNIQUE, 
    title TEXT NOT NULL,
    description_offer TEXT NOT NULL,
    salary_cents BIGINT NOT NULL,
    responses BIGINT DEFAULT 0
);
