USE [BIS]
GO

SET IDENTITY_INSERT [dbo].[Bank] ON
INSERT [dbo].[Bank]
    ([Id], [GlobalIdentifier], [Name], [Address], [CountryId], [BankTypeId])
VALUES
    (3, N'MMMM-EE-AAA-123', N'Bank Negara Malaysia', N'Jalan Dato Onn', 1, 2)

INSERT [dbo].[Bank]
    ([Id], [GlobalIdentifier], [Name], [Address], [CountryId], [BankTypeId])
VALUES
    (4, N'SSSS-WW-PP-123', N'Monetary Authority of Singapore', N'10 Shenton Wy', 2, 2)

SET IDENTITY_INSERT [dbo].[Bank] OFF
GO