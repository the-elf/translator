insert into message_template (language_id, code, text)
values ((select id from language where code = 'ru'), 'GEORGIAN_TRANSLIT_MSG',
        'Это грузинский транслит?'),
       ((select id from language where code = 'en'), 'GEORGIAN_TRANSLIT_MSG',
        'Is this Georgian transliteration?');