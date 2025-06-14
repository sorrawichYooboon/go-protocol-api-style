CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    release_date DATE
);

INSERT INTO
    movies (
        title,
        description,
        release_date
    )
VALUES (
        'The Matrix',
        'A computer hacker learns about the true nature of reality and his role in the war against its controllers.',
        '1999-03-31'
    ),
    (
        'Inception',
        'A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea.',
        '2010-07-16'
    ),
    (
        'Interstellar',
        'A team of explorers travel through a wormhole in space in an attempt to ensure humanity''s survival.',
        '2014-11-07'
    );