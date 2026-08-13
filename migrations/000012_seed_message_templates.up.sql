insert into message_template (language_id, code, text)
values ((select id from language where code = 'ru'), 'NOT_GEORGIAN_TEXT_ERROR_MSG',
        '⚠️ Сейчас бот принимает текст только на грузинском языке - мхедрули или латинская транслитерация.'),
       ((select id from language where code = 'en'), 'NOT_GEORGIAN_TEXT_ERROR_MSG',
        '⚠️ The bot currently accepts text in Georgian only - Mkhedruli script or Latin transliteration.');