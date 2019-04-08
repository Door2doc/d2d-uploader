USE upload;

DROP TABLE IF EXISTS opname_opname;
DROP TABLE IF EXISTS patient_patient;
DROP TABLE IF EXISTS seh_sehmut;
DROP TABLE IF EXISTS seh_sehreg;

CREATE TABLE opname_opname (
    inschrtijd NVARCHAR(5),
    afdeling   NVARCHAR(4),
    plannr     NVARCHAR(10),
    specialism NVARCHAR(5)
);

CREATE TABLE patient_patient (
    patientnr NVARCHAR(13),
    gebdat    DATETIME
);

CREATE TABLE seh_sehmut (
    sehmutid   NVARCHAR(10),
    sehid      NVARCHAR(10),
    behkamerco NVARCHAR(4),
    bednr      NVARCHAR(2),
    datum      DATETIME
);

CREATE TABLE seh_sehreg (
    sehid      NVARCHAR(10),
    patientnr  NVARCHAR(13),
    locatiecod NVARCHAR(3),
    datum      DATETIME,
    aanksdatum DATETIME,
    aankstijd  NVARCHAR(5),
    triagetijd NVARCHAR(5),
    artsbhtijd NVARCHAR(5),
    patgezt    NVARCHAR(5),
    arbehetijd NVARCHAR(5),
    klacht     NTEXT,
    specialism NVARCHAR(5),
    trianivcod NVARCHAR(2),
    vvcode     NVARCHAR(5),
    vervoertyp NVARCHAR(3),
    bestemming NVARCHAR(30),
    opnameid   NVARCHAR(10),
    vervall    BIT
);

INSERT INTO seh_sehreg(sehid, patientnr, locatiecod, datum, aanksdatum, aankstijd, patgezt,
                       klacht, specialism, vvcode, vervoertyp, bestemming, opnameid, vervall)
VALUES ('84192', '1', 'BLA', getdate(), getdate(), '8:54', '8:54', 'val van flat', 'CHI', 'AMB', 'EIG', 'OPN', '1', 0);

INSERT INTO seh_sehmut(sehmutid, sehid, behkamerco, bednr, datum)
VALUES ('368203', '84192', 'BED9', '1', getdate()),
       ('368204', '84192', NULL, NULL, getdate()),
       ('368205', '84192', 'BED9', '1', getdate());

insert into patient_patient(patientnr, gebdat) values (1, '1978-07-24 00:00:00');



