CREATE TABLE public.responses (
    id SERIAL PRIMARY KEY,
    vacancy_id INT NOT NULL,
    email VARCHAR(225) NOT NULL UNIQUE,
    FOREIGN KEY (vacancy_id) REFERENCES public.vacancy(id)
);