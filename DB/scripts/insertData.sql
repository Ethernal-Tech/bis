USE [BIS]
GO

SET IDENTITY_INSERT [dbo].[Country] ON
INSERT [dbo].[Country]
    ([Id], [Name], [CountryCode])
VALUES
    (1, N'Malaysia', N'')
INSERT [dbo].[Country]
    ([Id], [Name], [CountryCode])
VALUES
    (2, N'Singapore', N'+65')
INSERT [dbo].[Country]
    ([Id], [Name], [CountryCode])
VALUES
    (3, N'Liechtenstein', N'+423')
INSERT [dbo].[Country]
    ([Id], [Name], [CountryCode])
VALUES
    (4, N'Switzerland', N'+41')
INSERT [dbo].[Country]
    ([Id], [Name], [CountryCode])
VALUES
    (5, N'Iraq', N'+964')
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
    (3, N'EEEE-NN-MM-123', N'DBS Bank Limited', N'10 Shenton Wy', 2)
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
    (3, N'EEEE-NN-MM-321', N'J.P. Morgan', N'Jln. Bukit Ho Swee', 3)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (4, N'AA-BB-CC-000', N'M23', N'Address', 1)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (5, N'AA-BB-CC-111', N'John Doe', N'Address', 1)
INSERT [dbo].[BankClient]
    ([Id], [GlobalIdentifier], [Name], [Address], [BankId])
VALUES
    (6, N'AA-BB-CC-222', N'John Doe', N'Address', 1)
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
    (3, N'DBS Admin', N'DBSadmin', N'admin', 3)
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (4, N'Aviatrans Admin', N'AAadmin', N'admin', 4)
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (5, N'Logarcheo Admin', N'LGadmin', N'admin', 5)
INSERT [dbo].[BankEmployee]
    ([Id], [Name], [Username], [Password], [BankId])
VALUES
    (6, N'Al-Arabi Admin', N'ALadmin', N'admin', 6)
SET IDENTITY_INSERT [dbo].[BankEmployee] OFF
GO

SET IDENTITY_INSERT [dbo].[TransactionType] ON
INSERT [dbo].[TransactionType]
    ([Id], [Code], [Name])
VALUES
    (1, N'31111', 'Borrowing 1')
INSERT [dbo].[TransactionType]
    ([Id], [Code], [Name])
VALUES
    (2, N'31121', 'Borrowing 2')
INSERT [dbo].[TransactionType]
    ([Id], [Code], [Name])
VALUES
    (3, N'30001', 'Loan Repayment')
SET IDENTITY_INSERT [dbo].[TransactionType] OFF
GO

SET IDENTITY_INSERT [dbo].[Policy] ON
INSERT [dbo].[Policy]
    ([Id], [Name])
VALUES
    (1, N'Capital Flow Management')
INSERT [dbo].[Policy]
    ([Id], [Name])
VALUES
    (2, N'Saction Check List')
SET IDENTITY_INSERT [dbo].[Policy] OFF
GO

SET IDENTITY_INSERT [dbo].[TransactionTypePolicy] ON
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [PolicyId], [CountryId], [Amount])
VALUES
    (1, 1, 1, 1, 150000000)
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [PolicyId], [CountryId], [Amount])
VALUES
    (2, 1, 2, 2, 0)
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [PolicyId], [CountryId], [Amount])
VALUES
    (3, 2, 1, 1, 150000000)
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [PolicyId], [CountryId], [Amount])
VALUES
    (4, 2, 2, 1, 0)
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [PolicyId], [CountryId], [Amount])
VALUES
    (5, 1, 1, 3, 500000000)
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [PolicyId], [CountryId], [Amount])
VALUES
    (6, 1, 2, 3, 0)
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [PolicyId], [CountryId], [Amount])
VALUES
    (7, 1, 1, 4, 500000000)
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [PolicyId], [CountryId], [Amount])
VALUES
    (8, 1, 2, 4, 0)
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [PolicyId], [CountryId], [Amount])
VALUES
    (9, 1, 1, 5, 500000000)
INSERT [dbo].[TransactionTypePolicy]
    ([Id], [TransactionTypeId], [PolicyId], [CountryId], [Amount])
VALUES
    (10, 1, 2, 5, 0)
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
