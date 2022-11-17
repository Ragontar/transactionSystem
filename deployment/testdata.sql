INSERT INTO public.users (user_id, name, surname)
VALUES
    ('dfe67bab-d676-49e0-ac0e-cc076118a7dd', 'Vasya', 'Pupkin'),
    ('bf2694fa-222c-47bf-8a53-a05c8112745f', 'Timoha', 'Erohin'),
    ('2825a71e-bfa8-4aeb-8198-567aa6e34bb9', 'Kalistrat', 'Sichov');

INSERT INTO public.accounts (account_id, owner_id, balance)
VALUES
    ('abfc7e41-7862-44e9-8106-3aeeb31601e9', 'dfe67bab-d676-49e0-ac0e-cc076118a7dd', 1000),
    ('b981d68f-323e-4980-b982-cfbbe2dfca2e', 'bf2694fa-222c-47bf-8a53-a05c8112745f', 10000),
    ('40504a71-dd71-4060-8351-5f8970761b43', '2825a71e-bfa8-4aeb-8198-567aa6e34bb9', 0);