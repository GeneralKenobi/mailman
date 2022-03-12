SET search_path TO mailmandb;

insert into customer(email)
values ('john.smith@yahoo.com');
insert into customer(email)
values ('anna@gmail.com');

-- Mails for john
insert into mailing_entry(customer_id, mailing_id, title, content, insert_time)
values ((select id from customer where email = 'john.smith@yahoo.com'),
        1,
        'Welcome to mailman',
        'Hi John\n\n, Welcome to mailman!\n\n See you around',
        '2022-03-12T10:16:38.725412916Z');
insert into mailing_entry(customer_id, mailing_id, title, content, insert_time)
values ((select id from customer where email = 'john.smith@yahoo.com'),
        2,
        'Terms of usage',
        'Hi John\n\n, Here are the terms of usage\n\n...',
        '2022-03-12T10:19:01.123456789Z');

-- Mails for Anna
insert into mailing_entry(customer_id, mailing_id, title, content, insert_time)
values ((select id from customer where email = 'anna@gmail.com'),
        1,
        'Welcome to mailman',
        'Hi Anna\n\n, Welcome to mailman!\n\n See you around',
        '2022-03-12T10:19:01.123456789Z');