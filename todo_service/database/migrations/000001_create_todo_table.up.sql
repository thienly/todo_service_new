BEGIN ;
CREATE TABLE IF NOT EXISTS users(
                                    id serial PRIMARY KEY,
                                    username VARCHAR (50) UNIQUE NOT NULL,
                                    password VARCHAR (50) NOT NULL,
                                    email VARCHAR (300) UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS todos(
                                    id serial PRIMARY KEY,
                                    user_id serial NOT NULL,
                                    title VARCHAR (50) NOT NULL,
                                    description VARCHAR (200) ,
                                    done  BOOLEAN
);
COMMIT ;