# Changelog

All notable changes to this project will be documented in this file.

## [v1.9.2] - 2025-03-23

- fix handling join results with null result entities
- joins on join too

## [v1.9.0] - 2025-03-22

- add other join types for multiJoin
- add naming a multiJoin

## [v1.8.3] - 2025-02-11

- fix return types of boolean operators (SQLExpression)

## [v1.8.2] - 2025-02-10

- allow NewColumns with initial list

## [v1.8.1] - 2025-01-23

- fix int32Access.Equals to accept int32 too
- add NewColumns()

## [v1.8.0] - 2025-01-22

- add access to CommandTag in iterator to get RowsAffected value

## [v1.7.1] - 2024-12-24

- add TextArrayToStrings

## [v1.7.0] - 2024-12-20

- add ability to specify custom datatype mapping

## [v1.6.3] - 2024-12-19

- fix handling IN operator with empty values

## [v1.6.2, v1.6.1, v1.6.0] - 2024-12-19

- add GetParams to iterator for inspection.

## [v1.5.1] - 2024-12-16

- fix handling text[] datatype

## [v1.5.0] - 2024-09-25

- add FOR UPDATE,SHARE,...
- add SKIP LOCKED
- updated dependencies

## [v1.4.2] - 2024-01-23

- add RightOuterJoin and deprecate RightJoin
- add checks in JOIN SQL writing to prevent panics

## [v1.4.1] - 2024-01-18

- remove dependency on field ordinal of pgcolum

## [v1.4.0] - 2024-01-15

- add operand with google uuid.UUID type

## [v1.3.0] - 2023-10-02

- add FieldAccess.GreaterThan

## [v1.2.3] - 2023-08-07

- fix custom expression result handling

## [v1.2.2] - 2023-08-07

- (potential breaking) change of signature of FieldAccess[T].In() (from any -> T)

## [v1.2.1] - 2023-08-07 

- put SQL AS expression in brackets

## [v1.2.0] - 2023-06-21

- relax interface of ordery, add NewSQLConstant

## [v1.1.0] - 2023-06-16

- add ilike

## [v1.0.3] - 2023-05-24

- replace panic call on StringToUUID conversion

## [v1.0.2] - 2023-05-15
q
- indent non-comparision operators in pretty SQL

## [v1.0.1] - 2023-05-01

- SQL printing uses no tabs/line-ends; use IndentedSQL for a more-pretty form

## [v1.0.0] - 2023-04-30

- initial major version after one year using it in a commercial production environment.