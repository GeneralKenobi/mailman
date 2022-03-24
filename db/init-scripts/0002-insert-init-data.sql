SET search_path TO mailmandb;

INSERT INTO customer(email)
VALUES ('john.smith@yahoo.com');
INSERT INTO customer(email)
VALUES ('anna@gmail.com');

-- Mails for john
INSERT INTO mailing_entry(customer_id, mailing_id, title, content, insert_time)
VALUES ((SELECT id FROM customer WHERE email = 'john.smith@yahoo.com'),
        1,
        'Welcome to mailman',
        'Hi John\n\n, Welcome to mailman!\n\n See you around',
        '2022-03-12T10:16:38.725412916Z');
INSERT INTO mailing_entry(customer_id, mailing_id, title, content, insert_time)
VALUES ((SELECT id FROM customer WHERE email = 'john.smith@yahoo.com'),
        2,
        'Terms of usage',
        'Hi John\n\n, Here are the terms of usage\n\n...',
        '2022-03-12T10:19:01.123456789Z');

-- Mails for Anna
INSERT INTO mailing_entry(customer_id, mailing_id, title, content, insert_time)
VALUES ((SELECT id FROM customer WHERE email = 'anna@gmail.com'),
        1,
        'Welcome to mailman',
        'Hi Anna\n\n, Welcome to mailman!\n\n See you around',
        '2022-03-12T10:19:01.123456789Z');
