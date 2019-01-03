CREATE DATABASE Commendeer
	WITH
	ENCODING = 'UTF8'
	CONNECTION LIMIT = -1;

\c commendeer

CREATE TABLE IF NOT EXISTS AccessCode (
	CodeID SERIAL PRIMARY KEY,
	Email VARCHAR(150),
	SystemUsername VARCHAR(50) NOT NULL,
	Code VARCHAR(10),
	Used BOOL NOT NULL
);

CREATE TABLE IF NOT EXISTS UserInfo (
	UserID SERIAL PRIMARY KEY,
	Username VARCHAR(50) NOT NULL,
	Pass VARCHAR(500) NOT NULL,
	Administrator BOOL NOT NULL
);

CREATE TABLE IF NOT EXISTS QuestionType (
	QuestionTypeID SERIAL PRIMARY KEY,
	Description TEXT
);

CREATE TABLE IF NOT EXISTS Question (
	QuestionID SERIAL PRIMARY KEY,
	QuestionTypeID INTEGER REFERENCES QuestionType(QuestionTypeID),
	QuestionOrder INTEGER,
	Title TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS MultiChoiceQuestionOption (
	MultiChoiceQuestionOptionID SERIAL PRIMARY KEY,
	QuestionID INTEGER REFERENCES Question(QuestionID),
	OptionDescription TEXT
);

CREATE TABLE IF NOT EXISTS MultiChoiceQuestionOption_Result (
	MultiChoiceQuestionOption_ResultID SERIAL PRIMARY KEY,
	QuestionID INTEGER REFERENCES Question(QuestionID),
	MultiChoiceQuestionOptionID INTEGER REFERENCES MultiChoiceQuestionOption(MultiChoiceQuestionOptionID),
	CodeID INTEGER REFERENCES AccessCode(CodeID)
);

CREATE TABLE IF NOT EXISTS Question_Result (
	Question_ResultID SERIAL PRIMARY KEY,
	QuestionID INTEGER REFERENCES Question(QuestionID),
	CodeID INTEGER REFERENCES AccessCode(CodeID),
	Answer TEXT
);

CREATE TABLE IF NOT EXISTS QuantitativeResult (
	QuantitativeResultID SERIAL PRIMARY KEY,
	QuestionID INTEGER REFERENCES Question(QuestionID),
	Total NUMERIC(7, 2)
);