CREATE TABLE public.responses (
    id SERIAL PRIMARY KEY,
    vacancy_id INT NOT NULL,
    email VARCHAR(225) NOT NULL,
    FOREIGN KEY (vacancy_id) REFERENCES public.vacancies(id)
);