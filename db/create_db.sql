CREATE TABLE IF NOT EXISTS users (
    user_id      int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    balance int CHECK(balance >= 0)  NOT NULL
);

INSERT INTO users(balance) VALUES('800');
INSERT INTO users(balance) VALUES('1200');
INSERT INTO users(balance) VALUES('500');
INSERT INTO users(balance) VALUES('100');
INSERT INTO users(balance) VALUES('3400');
INSERT INTO users(balance) VALUES('900');
INSERT INTO users(balance) VALUES('840');

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    sender int NOT NULL REFERENCES users(user_id),
    count int CHECK(count >= 0) NOT NULL,
    receiver int NOT NULL REFERENCES users(user_id)
);
