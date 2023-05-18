CREATE TABLE auth (
                      id SERIAL PRIMARY KEY,
                      login VARCHAR(255) NOT NULL,
                      password VARCHAR(255) NOT NULL,
                      is_authorized BOOLEAN NOT NULL DEFAULT false,
                      login_time TIMESTAMP WITHOUT TIME ZONE,
                      logout_time TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE variants (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL
);

CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,
                       variant_id INTEGER NOT NULL,
                       task VARCHAR(255) NOT NULL,
                       correct_answer VARCHAR(255) NOT NULL,
                       answer_1 VARCHAR(255) NOT NULL,
                       answer_2 VARCHAR(255) NOT NULL,
                       answer_3 VARCHAR(255) NOT NULL,
                       answer_4 VARCHAR(255) NOT NULL,
                       FOREIGN KEY (variant_id) REFERENCES variants(id) ON DELETE CASCADE
);

CREATE TABLE answers (
                         auth_id INTEGER NOT NULL,
                         variant_id INTEGER NOT NULL,
                         task_id INTEGER NOT NULL,
                         answer VARCHAR(255) NOT NULL,
                         answer_time TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE results (
                         id SERIAL PRIMARY KEY,
                         variant_id INTEGER NOT NULL,
                         user_id INTEGER NOT NULL,
                         task_id INTEGER NOT NULL,
                         is_correct BOOLEAN NOT NULL,
                         answer_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
                         FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
                         FOREIGN KEY (user_id) REFERENCES auth(id) ON DELETE CASCADE
);

INSERT INTO auth (login, password, is_authorized, login_time, logout_time) VALUES ('1', '1',false, null, null);

INSERT INTO variants (name) VALUES ('Test 1'), ('Test 2');

INSERT INTO tasks (variant_id, task, correct_answer, answer_1, answer_2, answer_3, answer_4) VALUES
                                                                                                 (1, 'What is the meaning of the word "dog"?', 'A domesticated carnivorous mammal', 'A domesticated herbivorous mammal', 'A wild carnivorous mammal', 'A wild herbivorous mammal', 'A domesticated carnivorous mammal'),
                                                                                                 (1, 'What is the meaning of the word "apple"?', 'A round fruit with red or green skin and firm white flesh', 'A long yellow fruit with a curved shape', 'A small yellow fruit with a fuzzy skin', 'A small red fruit with a juicy interior', 'A round fruit with red or green skin and firm white flesh'),
                                                                                                 (2, 'What is the meaning of the word "computer"?', 'An electronic device for storing and processing data', 'A mechanical device for performing mathematical calculations', 'A device for converting one form of energy to another', 'A machine for printing text or images', 'An electronic device for storing and processing data'),
                                                                                                 (2, 'What is the meaning of the word "table"?', 'A piece of furniture with a flat top and one or more legs', 'A piece of furniture for sleeping on', 'A piece of furniture for storing clothes', 'A piece of furniture for sitting on', 'A piece of furniture with a flat top and one or more legs');
