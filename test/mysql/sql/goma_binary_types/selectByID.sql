select
  id
, binary_columns
, tinyblob_columns
, blob_columns
, mediumblob_columns
, longblob_columns
, varbinary_columns
FROM
  goma_binary_types
WHERE
  id = /* id */64
