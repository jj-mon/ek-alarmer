package configmanager

var RuleTmpl = `SELECT
	"{{.ProjectID}}" AS project_id,
	{{.StreamName}}.payload.event.sourceName AS source_name,		
    cast({{.StreamName}}.payload.event.readings[0].value, 'float') as v,
    lag(v) as prev_v,
  CASE
    WHEN cast({{.StreamName}}.payload.event.readings[0].value, 'float') < {{.LoLo}} AND lag(v) >= {{.LoLo}} THEN CRITICAL_LOW
    WHEN cast({{.StreamName}}.payload.event.readings[0].value, 'float') < {{.Lo}} AND lag(v) >= {{.Lo}} THEN WARNING_LOW
    WHEN cast({{.StreamName}}.payload.event.readings[0].value, 'float') > {{.Hi}} AND lag(v) <= {{.Hi}} THEN WARNING_HIGH
    WHEN cast({{.StreamName}}.payload.event.readings[0].value, 'float') > {{.HiHi}} AND lag(v) <= {{.HiHi}} THEN CRITICAL_HIGH
  ELSE IGNORE
  END as alarmLevel
FROM {{.StreamName}}
WHERE
	source_name = "{{.SourceName}}" AND
    (cast({{.StreamName}}.payload.event.readings[0].value, 'float') < {{.LoLo}} AND lag(v) >= {{.LoLo}}) OR
    (cast({{.StreamName}}.payload.event.readings[0].value, 'float') < {{.Lo}} AND lag(v) >= {{.Lo}}) OR
    (cast({{.StreamName}}.payload.event.readings[0].value, 'float') > {{.Hi}} AND lag(v) <= {{.Hi}}) OR
    (cast({{.StreamName}}.payload.event.readings[0].value, 'float') > {{.HiHi}} AND lag(v) <= {{.HiHi}})`
