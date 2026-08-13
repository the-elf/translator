insert into button (language_id, group_name, data, title)
values ((select id from language where code = 'ru'), 'geo_translit', 'yes', 'Да'),
       ((select id from language where code = 'ru'), 'geo_translit', 'no', 'Нет'),
       ((select id from language where code = 'en'), 'geo_translit', 'yes', 'Yes'),
       ((select id from language where code = 'en'), 'geo_translit', 'no', 'No')