USE [BIS]
GO

SET IDENTITY_INSERT [dbo].[Country] ON
INSERT [dbo].[Country]
    ([Id], [Name], [CountryCode])
VALUES
    (1, N'Malaysia', N'+60')
INSERT [dbo].[Country]
    ([Id], [Name], [CountryCode])
VALUES
    (2, N'Singapore', N'+65')
SET IDENTITY_INSERT [dbo].[Country] OFF
GO

SET IDENTITY_INSERT [dbo].[Bank] ON
INSERT [dbo].[Bank]
    ([Id], [GlobalIdentifier], [Name], [Address], [CountryId])
VALUES
    (1, N'549300BUPYUQGB5BFX94', N'Hong Leong Bank', N'Jalan Dato Onn 50480 Kuala Lumpur', 1)
INSERT [dbo].[Bank]
    ([Id], [GlobalIdentifier], [Name], [Address], [CountryId])
VALUES
    (2, N'984500653R409CC5AB28', N'JPM Chase', N'10 Shenton Wy', 2)
SET IDENTITY_INSERT [dbo].[Bank] OFF
GO

SET IDENTITY_INSERT [dbo].[BankClient] ON
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (1, N'AAAA-BB-CC-321', N'Company A', N'Company A Address', 1)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (2, N'AAAA-BB-CC-321', N'Company A', N'Company A Address', 2)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (4, N'EEEE-NN-MM-321', N'Company B', N'Company B Address', 1)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (5, N'EEEE-NN-MM-321', N'Company B', N'Company B Address', 2)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (7, N'AA-BB-CC-000', N'Company C', N'Company C Address', 1)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (8, N'AA-BB-CC-000', N'Company C', N'Company C Address', 2)
SET IDENTITY_INSERT [dbo].[BankClient] OFF
GO

SET IDENTITY_INSERT [dbo].[BankEmployee] ON
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (1, N'Hong Leong Bank Admin', N'HLBadmin', N'admin', 1)
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (2, N'J.P. Admin', N'JPMadmin', N'admin', 2)
SET IDENTITY_INSERT [dbo].[BankEmployee] OFF
GO

SET IDENTITY_INSERT [dbo].[TransactionType] ON
INSERT [dbo].[TransactionType]
    ([Id], [Code], [Name])
VALUES
    (1, N'DDWN', 'Loan Drawdown')
INSERT [dbo].[TransactionType]
    ([Id], [Code], [Name])
VALUES
    (2, N'PPAY', 'Loan Repayment')
SET IDENTITY_INSERT [dbo].[TransactionType] OFF
GO

SET IDENTITY_INSERT [dbo].[Policy] ON
INSERT [dbo].[Policy]
    ([Id], [CountryId], [Code], [Name])
VALUES
    (1, 1, N'CFM', N'MY - Foreign Exchange Policy')
INSERT [dbo].[Policy]
    ([Id], [CountryId], [Code], [Name])
VALUES
    (2, 1, N'SCL', N'MY - Sanctions')
INSERT [dbo].[Policy]
    ([Id], [CountryId], [Code], [Name])
VALUES
    (3, 2, N'SCL', N'SG - Sanctions')
SET IDENTITY_INSERT [dbo].[Policy] OFF
GO

SET IDENTITY_INSERT [dbo].[TransactionTypePolicy] ON
-- DDWN, Malaysia, CFM+SCL
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [CountryId], [PolicyId], [Amount], [Checklist])
VALUES
    (1, 1, 1, 1, 100000000, N'')
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [CountryId], [PolicyId], [Amount], [Checklist])
VALUES
    (2, 1, 1, 2, 0, N'UN Sanction List')

-- DDWN Syngapore, SCL
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [CountryId], [PolicyId], [Amount], [Checklist])
VALUES
    (3, 2, 2, 3, 0, N'UN Sanction List')
SET IDENTITY_INSERT [dbo].[TransactionTypePolicy] OFF
GO

SET IDENTITY_INSERT [dbo].[Status] ON
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (1, N'TransactionCreated')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (2, N'PoliciesApplied')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (3, N'ComplianceProofRequested')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (4, N'ComplianceCheckPassed')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (5, N'ProofInvalid')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (6, N'AssetSent')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (7, N'TransactionCompleted')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (8, N'TransactionCanceled')
SET IDENTITY_INSERT [dbo].[Status] OFF
GO