CREATE TABLE public.vacancy
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    owner_email VARCHAR(255) NOT NULL, 
    title TEXT NOT NULL,
    description_offer TEXT NOT NULL,
    salary_cents bigint NOT NULL,
    PRIMARY KEY (id)
);
