CREATE TABLE users (
                         id SERIAL PRIMARY KEY,
                         name VARCHAR(100) NOT NULL,
                         key VARCHAR(100) NOT NULL,
                         secret VARCHAR(100) NOT NULL
);

CREATE TABLE books (
                       id SERIAL PRIMARY KEY,
                       isbn VARCHAR(50) NOT NULL UNIQUE,
                       title VARCHAR(100) NOT NULL,
                       author VARCHAR(100) NOT NULL,
                       published INTEGER NOT NULL,
                       pages INTEGER NOT NULL,
                       status INTEGER NOT NULL
);
