update message_template
set text = E'1) Проверь ВЕСЬ текст:
			- если он на грузинском языке, то сразу перейди к пункту 2
			- если он на грузинском транслите, то сразу верни сообщение "GEORGIAN_TRANSLITERATION"
			- в противном случае сразу верни сообщение "NOT_GEORGIAN_TEXT"
		2) Переведи текст на русский язык, без отсебятины. В ответ верни только перевод предоставленного текста и ничего больше.'
where code = 'TRANSLATION_PROMPT_TEXT'
  and language_id = (select id from language where code = 'ru');

update message_template
set text = E'1) Check the WHOLE text:
    		- if it is in Georgian (Mkhedruli), proceed to step 2
    		- if it is Georgian text written in Latin letters (transliteration), return "GEORGIAN_TRANSLITERATION"
    		- otherwise return "NOT_GEORGIAN_TEXT"
		2) Translate the text into English without adding anything extra. Return only the translation of the provided text and nothing else.'
where code = 'TRANSLATION_PROMPT_TEXT'
  and language_id = (select id from language where code = 'en');