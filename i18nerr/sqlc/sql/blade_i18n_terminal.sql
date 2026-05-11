-- name: SelectI18nByLangAndKey :one
select ilk.code as 'key',
    ilv.value,
    ilk.data_type
from blade_i18n_language ild
    join blade_i18n_value ilv
on ild.id = ilv.definition_id
    join blade_i18n_key ilk on ilk.id = ilv.key_id
    join blade_i18n_terminal ilt on ilk.terminal_id = ilt.id
where ild.is_deleted = 0
  and ild.status = 1
  and ilv.is_deleted = 0
  and ilv.status = 1
  and ilk.is_deleted = 0
  and ilk.status = 1
  and ilt.is_deleted = 0
  and ild.code = sqlc.arg('lang')
  and ilk.code = sqlc.arg('key')
  and ilt.code IN (sqlc.slice('terminal_codes'))
limit 1;

-- name: SelectLanguageTerminalByCode :one
SELECT id,
    parent_id,
    code,
    name,
    create_user,
    create_time,
    update_user,
    update_time,
    create_dept,
    status,
    is_deleted
FROM blade_i18n_terminal
WHERE code = ?
  AND is_deleted = 0
LIMIT 1;

-- name: SelectLanguageTerminalByParentId :many
SELECT id,
    parent_id,
    code,
    name,
    create_user,
    create_time,
    update_user,
    update_time,
    create_dept,
    status,
    is_deleted
FROM blade_i18n_terminal
WHERE parent_id = ?
  AND is_deleted = 0;
