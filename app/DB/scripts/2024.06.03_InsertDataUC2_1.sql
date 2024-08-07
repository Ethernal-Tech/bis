USE [BIS1]
GO
SET IDENTITY_INSERT [dbo].[BankType] ON

INSERT [dbo].[BankType]
    ([Id], [Name])
VALUES
    (1, N'Commercial')
INSERT [dbo].[BankType]
    ([Id], [Name])
VALUES
    (2, N'Central')
SET IDENTITY_INSERT [dbo].[BankType] OFF
GO

INSERT [dbo].[Jurisdiction]
    ([Id], [Name])
VALUES
    (N'SG', N'Singapore')
INSERT [dbo].[Jurisdiction]
    ([Id], [Name])
VALUES
    (N'MY', N'Malaysia')
INSERT [dbo].[Jurisdiction]
    ([Id], [Name])
VALUES
    (N'AU', N'Australia')
INSERT [dbo].[Jurisdiction]
    ([Id], [Name])
VALUES
    (N'KR', N'South Korea')

GO
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'5493003GYPR7VI37GG77', N'Reserve Bank of Australia', N'8 Chifley Square', N'AU', 2)
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'RVHOHKPBCJ2GSJ37YH94', N'Bank of Korea', N'39 Namdaemun-ro', N'KR', 2)
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'xxxxxxxxxxxxxxxxxxxx', N'AU Commercial Bank', N'-------', N'AU', 1)
INSERT [dbo].[Bank]
    ([GlobalIdentifier], [Name], [Address], [JurisdictionId], [BankTypeId])
VALUES
    (N'yyyyyyyyyyyyyyyyyyyy', N'KR Commercial Bank', N'-------', N'KR', 1)
GO
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (N'1', N'BOK Admin', N'BOKadmin', N'admin', N'RVHOHKPBCJ2GSJ37YH94')
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (N'2', N'RBA Admin', N'RBAadmin', N'admin', N'5493003GYPR7VI37GG77')
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (N'3', N'AU Commercial Bank Admin', N'AUCadmin', N'admin', N'xxxxxxxxxxxxxxxxxxxx')
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (N'4', N'KR Commercial Admin', N'KRCadmin', N'admin', N'yyyyyyyyyyyyyyyyyyyy')
GO
SET IDENTITY_INSERT [dbo].[PolicyType] ON

INSERT [dbo].[PolicyType]
    ([Id], [Code], [Name])
VALUES
    (1, N'AML', N'Sanctions check')
INSERT [dbo].[PolicyType]
    ([Id], [Code], [Name])
VALUES
    (2, N'AMT', N'KR CFM - Payments (SECU)')
INSERT [dbo].[PolicyType]
    ([Id], [Code], [Name])
VALUES
    (3, N'NETT', N'Multilateral Netting Reporting')
SET IDENTITY_INSERT [dbo].[PolicyType] OFF
GO
SET IDENTITY_INSERT [dbo].[State] ON

INSERT [dbo].[State]
    ([Id], [Name])
VALUES
    (1, N'Compliance check created')
INSERT [dbo].[State]
    ([Id], [Name])
VALUES
    (2, N'Policies applied')
INSERT [dbo].[State]
    ([Id], [Name])
VALUES
    (3, N'Compliance proof requested')
INSERT [dbo].[State]
    ([Id], [Name])
VALUES
    (4, N'Compliance proof generation failed')
INSERT [dbo].[State]
    ([Id], [Name])
VALUES
    (5, N'Compliance proof generation succeeded')
INSERT [dbo].[State]
    ([Id], [Name])
VALUES
    (6, N'Compliance proof attached to the selected settlement asset')
INSERT [dbo].[State]
    ([Id], [Name])
VALUES
    (7, N'Settlement asset transferred to the beneficiary bank')
INSERT [dbo].[State]
    ([Id], [Name])
VALUES
    (8, N'Assets released to the client')
SET IDENTITY_INSERT [dbo].[State] OFF
GO
SET IDENTITY_INSERT [dbo].[TransactionType] ON

INSERT [dbo].[TransactionType]
    ([Id], [Code], [Name])
VALUES
    (1, N'SECU', N'Acquisition of unlisted securities')
SET IDENTITY_INSERT [dbo].[TransactionType] OFF
GO