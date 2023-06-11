CREATE TABLE users
(
    id serial primary key,
    balance float   not null
);

CREATE TABLE transactions
(
    id serial primary key,
    user_id   int     references users (id),
    amount    float       not null,
    operation varchar(40) not null,
    date      timestamp   not null
);

INSERT INTO users (balance) VALUES (4.13);
INSERT INTO users (balance) VALUES (32);
INSERT INTO users (balance) VALUES (11.321);
INSERT INTO users (balance) VALUES (41.12);
INSERT INTO users (balance) VALUES (1.32);
INSERT INTO users (balance) VALUES (541.32);
INSERT INTO users (balance) VALUES (339.012);

INSERT INTO transactions (user_id, amount, operation, date) VALUES (1, 30, '', NOW());
INSERT INTO transactions (user_id, amount, operation, date) VALUES (2, 101, '', NOW());
INSERT INTO transactions (user_id, amount, operation, date) VALUES (2, 32, '', NOW());
INSERT INTO transactions (user_id, amount, operation, date) VALUES (3, 62, '', NOW());
INSERT INTO transactions (user_id, amount, operation, date) VALUES (3, 32, '', NOW());
INSERT INTO transactions (user_id, amount, operation, date) VALUES (3, 190, '', NOW());
INSERT INTO transactions (user_id, amount, operation, date) VALUES (3, 305, '', NOW());
INSERT INTO transactions (user_id, amount, operation, date) VALUES (3, 101, '', NOW());
INSERT INTO transactions (user_id, amount, operation, date) VALUES (3, 103, '', NOW());