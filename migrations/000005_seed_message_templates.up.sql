-- translation placeholders
insert into message_template (language_id, code, text)
values ((select id from language where code = 'ru'), 'TRANSLATION_PLACEHOLDER', '⏳Перевожу...'),
       ((select id from language where code = 'en'), 'TRANSLATION_PLACEHOLDER', '⏳Translating...');

-- language titles
insert into message_template (language_id, code, text)
values ((select id from language where code = 'ru'), 'LANGUAGE_TITLE_RU', 'Русский'),
       ((select id from language where code = 'en'), 'LANGUAGE_TITLE_RU', 'Russian');
insert into message_template (language_id, code, text)
values ((select id from language where code = 'ru'), 'LANGUAGE_TITLE_EN', 'Английский'),
       ((select id from language where code = 'en'), 'LANGUAGE_TITLE_EN', 'English');

-- translation prompts
insert into message_template (language_id, code, text)
values ((select id from language where code = 'ru'), 'TRANSLATION_PROMPT_TEXT',
        'Переведи текст на русский язык, без отсебятины. В ответ верни только перевод предоставленного текста и ничего больше.'),
       ((select id from language where code = 'en'), 'TRANSLATION_PROMPT_TEXT',
        'Translate the text into English without adding anything extra. Return only the translation of the provided text and nothing else.');

-- registration messages
insert into message_template (language_id, code, text)
values ((select id from language where code = 'ru'), 'NOT_REGISTERED_MSG',
        'Вы ещё не зарегистрированы. Пожалуйста, используйте команду /start.'),
       ((select id from language where code = 'en'), 'NOT_REGISTERED_MSG',
        'You are not registered yet. Please use /start command.');
insert into message_template (language_id, code, text)
values ((select id from language where code = 'ru'), 'ALREADY_REGISTERED_MSG',
        'Вы уже зарегистрированы!'),
       ((select id from language where code = 'en'), 'ALREADY_REGISTERED_MSG',
        'You are already registered!');
insert into message_template (language_id, code, text)
values ((select id from language where code = 'ru'), 'SUCCESSFULLY_REGISTERED_MSG',
        'Регистрация прошла успешно!'),
       ((select id from language where code = 'en'), 'SUCCESSFULLY_REGISTERED_MSG',
        'You''ve successfully registered!');

-- language messages
insert into message_template (language_id, code, text)
values ((select id from language where code = 'ru'), 'UNSUPPORTED_LANGUAGE_MSG',
        'Неподдерживаемый язык: %s'),
       ((select id from language where code = 'en'), 'UNSUPPORTED_LANGUAGE_MSG',
        'Unsupported language: %s');
insert into message_template (language_id, code, text)
values ((select id from language where code = 'ru'), 'LANGUAGE_SET_MSG', 'Язык установлен: %s'),
       ((select id from language where code = 'en'), 'LANGUAGE_SET_MSG', 'Language is set: %s');


-- unexpected error message
insert into message_template (language_id, code, text)
values ((select id from language where code = 'ru'), 'UNEXPECTED_ERR_MSG',
        '⚠️ Произошла неожиданная ошибка. Повторите попытку позже.'),
       ((select id from language where code = 'en'), 'UNEXPECTED_ERR_MSG',
        '⚠️ Unexpected error. Try again later.');