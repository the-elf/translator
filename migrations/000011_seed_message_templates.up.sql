insert into message_template (language_id, code, text)
values ((select id from language where code = 'ru'), 'TRANSLATION_ERROR_MSG',
        '⚠️ Не удалось выполнить перевод. Попробуйте ещё раз.'),
       ((select id from language where code = 'en'), 'TRANSLATION_ERROR_MSG',
        '⚠️ Translation failed. Please try again.');