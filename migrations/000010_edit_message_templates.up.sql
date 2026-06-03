update message_template
set text = E'Вы ещё не зарегистрированы.\nПожалуйста, используйте команду /start.'
where code = 'NOT_REGISTERED_MSG'
  and language_id = 2;

update message_template
set text = E'You are not registered yet.\nPlease use /start command.'
where code = 'NOT_REGISTERED_MSG'
  and language_id = 3;