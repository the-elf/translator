insert into button (language_id, group_name, data, title)
values ((select id from language where code = 'ru'), 'set_lang', 'ru', 'Русский'),
       ((select id from language where code = 'ru'), 'set_lang', 'en', 'Английский'),
       ((select id from language where code = 'en'), 'set_lang', 'ru', 'Russian'),
       ((select id from language where code = 'en'), 'set_lang', 'en', 'English')