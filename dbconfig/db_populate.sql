\c commendeer

-- Sets up registered users
INSERT INTO AccessCode (Email, SystemUsername, Used)
VALUES ('01@email.com', 'user1', FALSE),
	('02@email.com', 'user2', FALSE),
	('03@email.com', 'user3', FALSE),
	('04@email.com', 'user4', FALSE),
	('05@email.com', 'user5', FALSE),
	('06@email.com', 'user6', FALSE),
	('07@email.com', 'user7', FALSE),
	('08@email.com', 'user8', FALSE),
	('09@email.com', 'user9', FALSE),
	('10@email.com', 'user10', FALSE),
	('11@email.com', 'user11', FALSE),
	('12@email.com', 'user12', FALSE),
	('13@email.com', 'user13', FALSE),
	('14@email.com', 'user14', FALSE),
	('15@email.com', 'user15', FALSE),
	('16@email.com', 'user16', FALSE),
	('17@email.com', 'user17', FALSE),
	('18@email.com', 'user18', FALSE),
	('19@email.com', 'user19', FALSE),
	('20@email.com', 'user20', FALSE),
	('21@email.com', 'user21', FALSE),
	('22@email.com', 'user22', FALSE),
	('23@email.com', 'user23', FALSE),
	('24@email.com', 'user24', FALSE),
	('25@email.com', 'user25', FALSE),
	('26@email.com', 'user26', FALSE),
	('27@email.com', 'user27', FALSE),
	('28@email.com', 'user28', FALSE),
	('29@email.com', 'user29', FALSE),
	('30@email.com', 'user30', FALSE),
	('31@email.com', 'user31', FALSE),
	('32@email.com', 'user32', FALSE),
	('33@email.com', 'user33', FALSE),
	('34@email.com', 'user34', FALSE),
	('35@email.com', 'user35', FALSE),
	('36@email.com', 'user36', FALSE),
	('37@email.com', 'user37', FALSE),
	('38@email.com', 'user38', FALSE),
	('39@email.com', 'user39', FALSE),
	('40@email.com', 'user40', FALSE),
	('41@email.com', 'user41', FALSE),
	('42@email.com', 'user42', FALSE),
	('43@email.com', 'user43', FALSE),
	('44@email.com', 'user44', FALSE),
	('45@email.com', 'user45', FALSE),
	('46@email.com', 'user46', FALSE),
	('47@email.com', 'user47', FALSE),
	('48@email.com', 'user48', FALSE),
	('49@email.com', 'user49', FALSE),
	('50@email.com', 'user50', FALSE);

-- Sets up registered users: admins / users
INSERT INTO UserInfo (Username, Pass, Administrator)
VALUES ('admin1', '$2a$14$LGaIbgBOpAXs0yNiMjYl9.rAnVpW3FC.f7zFRHoKkjH7YLdt/LD4a', TRUE), --"4dm1n123"
	('admin2', '$2a$14$UFzQXGDzUE38Qadgomu0neDG6gvj7QCraGpsZLwKdaWC2xc9P4/Hy', TRUE), --"aDmIn3Z1"
	('sysuser1', '$2a$14$vyJVUMtm6lc15LIAZ8VCAOHeX9FHsBfxM/7khFqlCtuKXXhlyUAYS', FALSE), --"adN1M231"
	('sysuser2', '$2a$14$p8dbFHcxQy3.Gmd8w6rnSeWBWTEdtWt6nlAlVfZeeajul2rWBI83m', FALSE); --"321n1md4"