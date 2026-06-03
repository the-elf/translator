insert into message_template (language_id, code, text)
values ((select id from language where code = 'ru'), 'PREFERRED_LANGUAGE_MSG',
        'Выберите предпочитаемый язык'),
       ((select id from language where code = 'en'), 'PREFERRED_LANGUAGE_MSG',
        'Choose preferred language');