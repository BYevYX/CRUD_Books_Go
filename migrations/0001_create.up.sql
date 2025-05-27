CREATE TABLE authors (
  id SERIAL PRIMARY KEY,
  name VARCHAR(80) NOT NULL,
  birthdate DATE,
  death_date DATE
);

CREATE TABLE books (
  id SERIAL PRIMARY KEY,
  name VARCHAR(80) NOT NULL,
  pages_count INTEGER NOT NULL,
  publication_date DATE,
  author_id INTEGER NOT NULL,

  CONSTRAINT FK_author_id FOREIGN KEY (author_id) REFERENCES authors
);
