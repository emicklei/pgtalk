# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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