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
-- INSERT [dbo].[Country]
--     ([Id], [Name], [CountryCode])
-- VALUES
--     (3, N'Liechtenstein', N'+423')
-- INSERT [dbo].[Country]
--     ([Id], [Name], [CountryCode])
-- VALUES
--     (4, N'Switzerland', N'+41')
-- INSERT [dbo].[Country]
--     ([Id], [Name], [CountryCode])
-- VALUES
--     (5, N'Iraq', N'+964')
SET IDENTITY_INSERT [dbo].[Country] OFF
GO

SET IDENTITY_INSERT [dbo].[Bank] ON
INSERT [dbo].[Bank]
    ([Id], [GlobalIdentifier], [Name], [Address], [CountryId])
VALUES
    (1, N'AAAA-BB-CC-123', N'Hong Leong Bank', N'Jalan Dato Onn 50480 Kuala Lumpur', 1)
INSERT [dbo].[Bank]
    ([Id], [GlobalIdentifier], [Name], [Address], [CountryId])
VALUES
    (2, N'AAAA-BB-CC-234', N'Maybank', N'Jalan Dato Onn 50480 Kuala Lumpur', 1)
INSERT [dbo].[Bank]
    ([Id], [GlobalIdentifier], [Name], [Address], [CountryId])
VALUES
    (3, N'EEEE-NN-MM-123', N'J.P. Morgan', N'10 Shenton Wy', 2)
SET IDENTITY_INSERT [dbo].[Bank] OFF
GO

SET IDENTITY_INSERT [dbo].[BankClient] ON
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (1, N'AAAA-BB-CC-321', N'Ethernal', N'Jalan Dato Onn', 1)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (2, N'AAAA-BB-CC-321', N'Ethernal', N'Jalan Dato Onn', 2)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (3, N'AAAA-BB-CC-321', N'Ethernal', N'Jalan Dato Onn', 3)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (4, N'EEEE-NN-MM-321', N'J.P. Morgan Client', N'Jln. Bukit Ho Swee', 1)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (5, N'EEEE-NN-MM-321', N'J.P. Morgan Client', N'Jln. Bukit Ho Swee', 2)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (6, N'EEEE-NN-MM-321', N'J.P. Morgan Client', N'Jln. Bukit Ho Swee', 3)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (7, N'AA-BB-CC-000', N'M23', N'Address', 1)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (8, N'AA-BB-CC-000', N'M23', N'Address', 2)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (9, N'AA-BB-CC-000', N'M23', N'Address', 3)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (10, N'AA-BB-CC-111', N'John Doe', N'Address', 1)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (11, N'AA-BB-CC-111', N'John Doe', N'Address', 2)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (12, N'AA-BB-CC-111', N'John Doe', N'Address', 3)
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
    (2, N'Maybank Admin', N'MBadmin', N'admin', 2)
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (3, N'J.P. Admin', N'JPadmin', N'admin', 3)
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
    ([Id], [Code], [Name])
VALUES
    (1, N'CFM', N'MY - Foreign Exchange Policy')
INSERT [dbo].[Policy]
    ([Id], [Code], [Name])
VALUES
    (2, N'SCL', N'MY - Sanctions')
-- INSERT [dbo].[Policy]
--     ([Id], [Code], [Name])
-- VALUES
--     (3, N'CFM', N'SG - Foreign Exchange Policy')
INSERT [dbo].[Policy]
    ([Id], [Code], [Name])
VALUES
    (3, N'SCL', N'SG - Sanctions')
-- INSERT [dbo].[Policy]
--     ([Id], [Code], [Name])
-- VALUES
--     (5, N'SCL', N'Sanctions')
SET IDENTITY_INSERT [dbo].[Policy] OFF
GO

SET IDENTITY_INSERT [dbo].[TransactionTypePolicy] ON
-- DDWN, Malaysia, CFM+SCL
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [CountryId], [PolicyId], [Amount], [Checklist])
VALUES
    (1, 1, 1, 1, 150000000, N'')
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [CountryId], [PolicyId], [Amount], [Checklist])
VALUES
    (2, 1, 1, 2, 0, N'UN Sanction List')

-- DDWN Syngapore, CFM+SCL
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [CountryId], [PolicyId], [Amount], [Checklist])
VALUES
    (3, 1, 2, 3, 150000000, N'')
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [CountryId], [PolicyId], [Amount], [Checklist])
VALUES
    (4, 1, 2, 4, 0, N'UN Sanction List')

-- -- DDWN Liechtenstein, CFM
-- INSERT [dbo].[TransactionTypePolicy]
--     ([Id], [TransactionTypeId], [CountryId], [PolicyId], [Amount], [Checklist])
-- VALUES
--     (5, 1, 3, 1, 150000000, N'')

-- -- DDWN Switzerland, CFM
-- INSERT [dbo].[TransactionTypePolicy]
--     ([Id], [TransactionTypeId], [CountryId], [PolicyId], [Amount], [Checklist])
-- VALUES
--     (6, 1, 4, 1, 150000000, N'')

-- -- DDWN Iraq, CFM
-- INSERT [dbo].[TransactionTypePolicy]
--     ([Id], [TransactionTypeId], [CountryId], [PolicyId], [Amount], [Checklist])
-- VALUES
--     (7, 1, 5, 1, 150000000, N'')
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
    (3, N'ProofRequested')
INSERT [dbo].[Status]
    ([Id], [Name])
VALUES
    (4, N'ProofReceived')
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

-- SET IDENTITY_INSERT [dbo].[Transaction] ON 
-- INSERT [dbo].[Transaction] ([Id], [OriginatorBank], [BeneficiaryBank], [Sender], [Receiver], [Currency], [Amount], [TransactionTypeId]) VALUES (1, 1, 2, 1, 2, N'MYR', 100000, 1)
-- INSERT [dbo].[Transaction] ([Id], [OriginatorBank], [BeneficiaryBank], [Sender], [Receiver], [Currency], [Amount], [TransactionTypeId]) VALUES (2, 3, 1, 3, 1, N'MYR', 500000, 2)
-- INSERT [dbo].[Transaction] ([Id], [OriginatorBank], [BeneficiaryBank], [Sender], [Receiver], [Currency], [Amount], [TransactionTypeId]) VALUES (3, 4, 1, 4, 1, N'USD', 1000000, 1)
-- SET IDENTITY_INSERT [dbo].[Transaction] OFF
-- GO

-- SET IDENTITY_INSERT [dbo].[TransactionHistory] ON 
-- INSERT [dbo].[TransactionHistory] ([Id], [Transactionid], [StatusId], [Date]) VALUES (1, 1, 1, CAST(N'2023-10-10T00:00:00.0000000' AS DateTime2))
-- INSERT [dbo].[TransactionHistory] ([Id], [Transactionid], [StatusId], [Date]) VALUES (2, 1, 2, CAST(N'2023-10-11T00:00:00.0000000' AS DateTime2))
-- INSERT [dbo].[TransactionHistory] ([Id], [Transactionid], [StatusId], [Date]) VALUES (3, 1, 3, CAST(N'2023-10-12T00:00:00.0000000' AS DateTime2))
-- INSERT [dbo].[TransactionHistory] ([Id], [Transactionid], [StatusId], [Date]) VALUES (4, 1, 4, CAST(N'2023-10-13T00:00:00.0000000' AS DateTime2))
-- INSERT [dbo].[TransactionHistory] ([Id], [Transactionid], [StatusId], [Date]) VALUES (5, 1, 7, CAST(N'2023-10-13T15:00:00.0000000' AS DateTime2))
-- INSERT [dbo].[TransactionHistory] ([Id], [Transactionid], [StatusId], [Date]) VALUES (6, 2, 1, CAST(N'2023-10-14T00:00:00.0000000' AS DateTime2))
-- INSERT [dbo].[TransactionHistory] ([Id], [Transactionid], [StatusId], [Date]) VALUES (7, 3, 1, CAST(N'2023-10-15T00:00:00.0000000' AS DateTime2))
-- SET IDENTITY_INSERT [dbo].[TransactionHistory] OFF
-- GO
